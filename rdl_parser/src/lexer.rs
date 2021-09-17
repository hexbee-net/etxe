use std::str::CharIndices;

use crate::core::error::Errors;
use crate::parser_source::ParserSource;
use crate::pos::{self, Location, Spanned};
use crate::token::{StringLiteral, Token};
use itertools::Itertools;

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

pub struct Lexer<'input> {
  input: &'input str,
  chars: std::iter::Peekable<CharIndices<'input>>,
  loc: Location,
  pub errors: Errors<SpannedError>,
}

impl<'input> Iterator for Lexer<'input> {
  type Item = Result<SpannedToken<'input>, SpannedError>;

  fn next(&mut self) -> Option<Self::Item> {
    while let Some((start, ch)) = self.bump() {
      return match ch {
        '"' => Some(self.string_literal(start)),
        ch => None,
      };
    }

    let next_loc = self.next_loc()?;
    Some(Ok(pos::spanned(next_loc, next_loc, Token::EOF)))
  }
}

impl<'input> Lexer<'input> {
  pub fn new(input: &'input str) -> Self {
    Lexer {
      input,
      chars: input.char_indices().peekable(),
      // start_index: input.start_index(),
      loc: 0,
      errors: Errors::new(),
    }
  }

  fn bump(&mut self) -> Option<(Location, char)> {
    self.chars.next().map(|(loc, ch)| {
      self.loc = loc;
      (loc, ch)
    })
  }

  fn lookahead(&mut self) -> Option<(Location, char)> {
    self.chars.peek().cloned()
  }

  fn skip_to_end(&mut self) {
    while let Some(_) = self.bump() {}
  }

  fn curr_loc(&self) -> Location {
    self.loc
  }

  fn next_loc(&mut self) -> Option<Location> {
    self.lookahead().map(|l| l.0)
  }

  fn slice(&self, start: Location, end: Location) -> &'input str {
    //   let start = start.absolute - ByteOffset::from(self.start_index.to_usize() as i64);
    //   let end = end.absolute - ByteOffset::from(self.start_index.to_usize() as i64);
    // &self.input[start.to_usize()..end.to_usize()]
    &self.input[start..end]
  }

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
    while let Some((end, ch)) = self.lookahead() {
      if terminate(ch) {
        return Some((end, self.slice(start, end)));
      } else {
        self.bump();
      }
    }

    self.next_loc().map(|l| (l, self.slice(start, l)))
  }

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

  fn escape_code(&mut self, start: Location) -> Result<char, SpannedError> {
    match self.bump() {
      Some((_, 'n')) => Ok('\n'),
      Some((_, 'r')) => Ok('\r'),
      Some((_, 't')) => Ok('\t'),
      Some((_, '"')) => Ok('"'),
      Some((_, '\\')) => Ok('\\'),
      Some((p, 'u')) => {
        todo!("take unicode char code")
      }

      Some((end, ch)) => self
        .recover(start, end, Error::UnexpectedEscapeCode(ch), ch)
        .map(|s| s.value),

      None => self.eof_recover('\0').map(|s| s.value),
    }
  }

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
      self.take_until(scan_start, |b| b == '"' || b == '\\');

      match self.bump() {
        Some((start, '\\')) => {
          self.escape_code(start)?;
        }

        Some((_, '"')) => {
          let end = self.curr_loc();
          let mut content_end = end;
          // content_end.absolute.0 -= 1;

          let token = Token::StringLiteral(StringLiteral::Escaped(
            self.slice(content_start, content_end),
          ));

          return Ok(pos::spanned(start, end, token));
        }

        _ => break,
      }
    }

    let end = self.curr_loc();

    let token = Token::StringLiteral(StringLiteral::Escaped(self.slice(content_start, end + 1)));
    self.recover(start, end, Error::UnterminatedStringLiteral, token)
  }

  fn heredoc_literal(&mut self, start: Location) -> Result<SpannedToken<'input>, SpannedError> {
    todo!()
  }
}

#[cfg(test)]
mod test {
  use super::*;

  fn lexer(input: &str) -> Lexer {
    Lexer::new(input)
  }

  fn test(input: &str, expected: BorrowedToken<'_>, err: Option<Error>) {
    let mut lexer = lexer(input);
    let end = input.len() - 1;
    let expected = pos::spanned(0, end, expected);
    println!("{:?}", expected);

    for token in lexer.by_ref() {
      match token {
        Ok(token) => {
          assert_eq!(token, expected)
        }
        Err(err) => {
          panic!(err)
        }
      }
    }

    assert_eq!(err.map(|e| pos::spanned(0, end, e)), lexer.errors.pop());
  }

  #[test]
  fn string_literals() {
    let inputs = vec![
      (
        r#""foo""#,
        Token::StringLiteral(StringLiteral::Escaped("foo")),
        None,
      ),
      (
        r#""foo"#,
        Token::StringLiteral(StringLiteral::Escaped("foo")),
        Some(Error::UnterminatedStringLiteral),
      ),
    ];

    for (input, expected, err) in inputs {
      test(input, expected, err);
    }
  }
}
