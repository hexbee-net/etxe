use std::str::CharIndices;

use crate::core::error::Errors;
use crate::pos::{self, Location, Spanned};
use crate::token::{self, Token};
use itertools::{Itertools, MultiPeek};

quick_error! {
  #[derive(Clone, Debug, PartialEq, Eq, Hash)]
  pub enum Error {
    UnexpectedChar(ch: char) {
      display("unexpected character")
    }
    UnexpectedEof {
      display("unexpected end of file")
    }
    UnterminatedStringLiteral {
      display("unterminated string literal")
    }
    UnexpectedEscapeCode(ch: char) {
      display("unexpected escape code")
    }
  }
}

pub type SpannedError = Spanned<Error, Location>;

pub type BorrowedToken<'input> = Token<&'input str>;
pub type SpannedToken<'input> = Spanned<Token<&'input str>, Location>;

type LexerResult<'a> = Result<SpannedToken<'a>, SpannedError>;

#[derive(Clone, Copy, Debug, PartialEq)]
enum Context<S> {
  String(token::StringDelimiter<S>),
  // SingleLineString,
  // MultiLineString,
  // HereDocString(S),
}

pub struct Lexer<'input> {
  input: &'input str,
  chars: MultiPeek<CharIndices<'input>>,
  loc: Location,

  context: Option<Context<&'input str>>,
  context_stack: Vec<Context<&'input str>>,

  pub errors: Errors<SpannedError>,
}

impl<'input> Iterator for Lexer<'input> {
  type Item = LexerResult<'input>;

  fn next(&mut self) -> Option<LexerResult<'input>> {
    match self.context.as_ref() {
      Some(Context::String(start_sequence)) => self.iter_string(start_sequence),
      None => self.iter_base(),
    }
  }
}

impl<'input> Lexer<'input> {
  pub fn new(input: &'input str) -> Self {
    Lexer {
      input,
      chars: input.char_indices().multipeek(),
      loc: 0,

      context: None,
      context_stack: Vec::new(),

      errors: Errors::new(),
    }
  }

  // Iterator contexts /////////////////

  pub fn iter_base(&mut self) -> Option<LexerResult<'input>> {
    while let Some((start, ch)) = self.bump() {
      return match ch {
        '"' => Some(self.string_start(self.curr_loc())),
        '<' => Some(self.heredoc_start(start)),
        _ => None,
      };
    }

    let next_loc = self.next_loc()?;
    Some(Ok(pos::spanned(next_loc, next_loc, Token::EOF)))
  }

  fn iter_string(
    &self,
    _start_sequence: &token::StringDelimiter<&str>,
  ) -> Option<LexerResult<'input>> {
    todo!()
  }

  // Token Lexers //////////////////////

  fn string_start(&mut self, start: Location) -> Result<SpannedToken<'input>, SpannedError> {
    let (end, delimiter) = match self.peek_while(start + 1, |c| c == '"') {
      None => (start + 1, token::StringDelimiter::SingleLine),
      Some((end, r#""#)) => (end, token::StringDelimiter::SingleLine),
      Some((_end, r#"""#)) => (start + 1, token::StringDelimiter::SingleLine),
      Some((_end, r#""""#)) | Some((_end, r#"""""""#)) => (
        self.advance_by(2).unwrap(),
        token::StringDelimiter::MultiLine,
      ),

      Some((end, _)) => {
        self.reset_peek();
        return Err(pos::spanned(start, end, Error::UnterminatedStringLiteral));
      }
    };

    self.reset_peek();

    self.context.map(|c| self.context_stack.push(c));
    self.context = Some(Context::String(delimiter));

    Ok(pos::spanned(start, end, Token::StringDelimiter(delimiter)))
  }

  fn heredoc_start(&mut self, _start: Location) -> Result<SpannedToken<'input>, SpannedError> {
    todo!()
  }

  #[allow(dead_code)]
  fn escape_code(&mut self, start: Location) -> Result<char, SpannedError> {
    match self.bump() {
      Some((_, 'n')) => Ok('\n'),
      Some((_, 'r')) => Ok('\r'),
      Some((_, 't')) => Ok('\t'),
      Some((_, '"')) => Ok('"'),
      Some((_, '\\')) => Ok('\\'),
      Some((_, 'u')) => {
        todo!("take unicode char code")
      }

      Some((end, ch)) => self
        .recover(start, end, Error::UnexpectedEscapeCode(ch), ch)
        .map(|s| s.value),

      None => self.eof_recover('\0').map(|s| s.value),
    }
  }

  #[allow(dead_code)]
  fn string_literal(&mut self, start: Location) -> Result<SpannedToken<'input>, SpannedError> {
    let content_start =
      self
        .next_loc()
        .ok_or(pos::spanned(start, start, Error::UnterminatedStringLiteral))?;

    loop {
      let scan_start =
        self
          .next_loc()
          .ok_or(pos::spanned(start, start, Error::UnterminatedStringLiteral))?;
      self.take_until(scan_start, |c| c == '"' || c == '\\');

      match self.bump() {
        Some((start, '\\')) => {
          self.escape_code(start)?;
        }

        Some((_, '"')) => {
          let end = self.curr_loc();
          let content_end = end;

          let token = Token::StringLiteral(token::StringLiteral::Escaped(
            self.slice(content_start, content_end),
          ));

          return Ok(pos::spanned(start, end, token));
        }

        _ => break,
      }
    }

    let end = self.curr_loc();

    let token = Token::StringLiteral(token::StringLiteral::Escaped(
      self.slice(content_start, end + 1),
    ));
    self.recover(start, end, Error::UnterminatedStringLiteral, token)
  }

  #[allow(dead_code)]
  fn heredoc_literal(&mut self, _start: Location) -> Result<SpannedToken<'input>, SpannedError> {
    todo!()
  }

  // Utilities /////////////////////////

  fn bump(&mut self) -> Option<(Location, char)> {
    self.chars.next().map(|(loc, ch)| {
      self.loc = loc;
      (loc, ch)
    })
  }

  fn peek(&mut self) -> Option<(usize, char)> {
    self.chars.peek().cloned()
  }

  fn reset_peek(&mut self) {
    self.chars.reset_peek();
  }

  #[allow(dead_code)]
  fn skip_to_end(&mut self) {
    while let Some(_) = self.bump() {}
  }

  fn advance_by(&mut self, n: usize) -> Result<Location, usize> {
    let mut loc = self.curr_loc();
    for i in 0..n {
      match self.bump() {
        Some((l, _)) => {
          loc = l;
        }
        None => {
          return Err(i);
        }
      }
    }

    Ok(loc)
  }

  #[allow(dead_code)]
  fn advance_to(&mut self, loc: Location) -> Result<Location, usize> {
    let mut i = 0;
    while let Some((l, _)) = self.bump() {
      if l == loc {
        return Ok(l);
      }
      i += 1
    }

    Err(i)
  }

  fn curr_loc(&self) -> Location {
    self.loc
  }

  fn next_loc(&mut self) -> Option<Location> {
    let loc = self.peek().map(|l| l.0);
    self.reset_peek();
    return loc;
  }

  fn slice(&self, start: Location, end: Location) -> &'input str {
    //   let start = start.absolute - ByteOffset::from(self.start_index.to_usize() as i64);
    //   let end = end.absolute - ByteOffset::from(self.start_index.to_usize() as i64);
    // &self.input[start.to_usize()..end.to_usize()]
    &self.input[start..end]
  }

  #[allow(dead_code)]
  fn take_while<F>(&mut self, start: Location, mut keep_going: F) -> Option<(Location, &'input str)>
  where
    F: FnMut(char) -> bool,
  {
    self.take_until(start, |c| !keep_going(c))
  }

  fn take_until<F>(&mut self, start: Location, mut terminate: F) -> Option<(Location, &'input str)>
  where
    F: FnMut(char) -> bool,
  {
    while let Some((end, ch)) = self.peek() {
      if terminate(ch) {
        self.reset_peek();
        return Some((end, self.slice(start, end)));
      } else {
        self.bump();
      }
    }

    self.next_loc().map(|l| (l, self.slice(start, l)))
  }

  fn peek_while<F>(&mut self, start: Location, mut keep_going: F) -> Option<(Location, &'input str)>
  where
    F: FnMut(char) -> bool,
  {
    self.peek_until(start, |c| !keep_going(c))
  }

  fn peek_until<F>(&mut self, start: Location, mut terminate: F) -> Option<(Location, &'input str)>
  where
    F: FnMut(char) -> bool,
  {
    let mut last_loc = None;
    while let Some((end, ch)) = self.peek() {
      last_loc = Some(end);
      if terminate(ch) {
        return Some((end, self.slice(start, end)));
      }
    }

    last_loc.map(|l| (l, self.slice(start, l + 1)))
  }

  #[allow(dead_code)]
  fn error<T>(&mut self, location: Location, code: Error) -> Result<T, SpannedError> {
    self.skip_to_end();
    Err(pos::spanned(location, location, code))
  }

  fn recover<T>(
    &mut self,
    start: Location,
    end: Location,
    code: Error,
    value: T,
  ) -> Result<Spanned<T, Location>, SpannedError> {
    self.errors.push(pos::spanned(start, end, code));
    Ok(pos::spanned(start, end, value))
  }

  fn eof_recover<T>(&mut self, value: T) -> Result<Spanned<T, Location>, SpannedError> {
    let end = self.curr_loc();
    self.recover(end, end, Error::UnexpectedEof, value)
  }
}

#[cfg(test)]
mod test {
  use super::*;
  use crate::token::StringDelimiter;

  fn lexer(input: &str) -> Lexer {
    Lexer::new(input)
  }

  #[test]
  fn string_start() {
    let tests: Vec<(
      &str,
      Result<SpannedToken, SpannedError>,
      Option<Context<&str>>,
    )> = vec![
      (
        r#"~"foo"~"#,
        Ok(pos::spanned(
          0,
          1,
          Token::StringDelimiter(token::StringDelimiter::SingleLine),
        )),
        Some(Context::String(token::StringDelimiter::SingleLine)),
      ),
      (
        r#"~""~"#,
        Ok(pos::spanned(
          0,
          1,
          Token::StringDelimiter(token::StringDelimiter::SingleLine),
        )),
        Some(Context::String(token::StringDelimiter::SingleLine)),
      ),
      (
        r#"~"~"#,
        Ok(pos::spanned(
          0,
          1,
          Token::StringDelimiter(token::StringDelimiter::SingleLine),
        )),
        Some(Context::String(token::StringDelimiter::SingleLine)),
      ),
      (
        r#"~"""foo"""~"#,
        Ok(pos::spanned(
          0,
          2,
          Token::StringDelimiter(token::StringDelimiter::MultiLine),
        )),
        Some(Context::String(token::StringDelimiter::MultiLine)),
      ),
      (
        r#"~"""~"""~"#,
        Ok(pos::spanned(
          0,
          2,
          Token::StringDelimiter(token::StringDelimiter::MultiLine),
        )),
        Some(Context::String(token::StringDelimiter::MultiLine)),
      ),
      (
        r#"~"""~"#,
        Ok(pos::spanned(
          0,
          2,
          Token::StringDelimiter(token::StringDelimiter::MultiLine),
        )),
        Some(Context::String(token::StringDelimiter::MultiLine)),
      ),
      (
        r#"~""""~"#,
        Err(pos::spanned(0, 3, Error::UnterminatedStringLiteral)),
        None,
      ),
    ];

    for (input, expected, expected_context) in tests {
      let s = input.replace("~", "");
      let mut lexer = lexer(&*s);
      let res = lexer.next().unwrap();

      assert_eq!(res, expected);

      match res {
        Ok(_) => {
          assert_eq!(lexer.context, expected_context);
        }
        Err(_) => {
          assert!(lexer.context.is_none());
        }
      }
    }
  }
}
