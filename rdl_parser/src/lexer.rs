use std::str::Chars;

use crate::core::error::Errors;
use crate::pos::{self, ByteIndex, ByteOffset, Location, RawIndex, Spanned};
use crate::token::{self, Token};
use crate::ParserSource;
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

pub type BorrowedToken<'input> = Token<&'input str>;

pub type SpannedError = Spanned<Error, Location>;
pub type SpannedToken<'input> = Spanned<Token<&'input str>, Location>;

type LexerResult<'a> = Result<SpannedToken<'a>, SpannedError>;

#[derive(Clone, Copy, Debug, PartialEq)]
enum Context<S> {
  General,
  String(token::StringDelimiter<S>),
  // SingleLineString,
  // MultiLineString,
  // HereDocString(S),
}

pub struct Lexer<'input> {
  input: &'input str,
  // chars: MultiPeek<CharIndices<'input>>,
  chars: MultiPeek<Chars<'input>>,
  loc: Location,
  peek_loc: Location,

  context: Option<Context<&'input str>>,
  context_stack: Vec<Context<&'input str>>,

  pub errors: Errors<SpannedError>,
}

impl<'input> Iterator for Lexer<'input> {
  type Item = LexerResult<'input>;

  fn next(&mut self) -> Option<LexerResult<'input>> {
    match self.context.as_ref() {
      Some(Context::General) => self.iter_base(),
      Some(Context::String(start_sequence)) => self.iter_string(start_sequence),
      None => self.iter_base(),
    }
  }
}

impl<'input> Lexer<'input> {
  pub fn new<S>(input: &'input S) -> Self
  where
    S: ?Sized + ParserSource,
  {
    let src = input.src();
    Lexer {
      input: src,
      chars: src.chars().multipeek(),
      loc: Location::new(0, 0, 0),
      peek_loc: Location::new(0, 0, 0),

      context: None,
      context_stack: Vec::new(),

      errors: Errors::new(),
    }
  }

  // Iterator contexts /////////////////

  pub fn iter_base(&mut self) -> Option<LexerResult<'input>> {
    while let Some((prev_loc, _, ch)) = self.bump() {
      return match ch {
        '"' => Some(self.string_start(prev_loc)),
        '<' => Some(self.heredoc_start(prev_loc)),
        _ => None,
      };
    }

    let next_loc = self.next_loc()?;
    Some(Ok(pos::spanned(next_loc, next_loc, Token::EOF)))
  }

  fn iter_string(&self, _start_sequence: &token::StringDelimiter<&str>) -> Option<LexerResult<'input>> {
    todo!()
  }

  // Token Lexers //////////////////////

  fn string_start(&mut self, start: Location) -> Result<SpannedToken<'input>, SpannedError> {
    let section = self.peek_while(start, |c| c == '"');
    self.reset_peek();

    let (end, delimiter) = match section {
      (end, r#"""#) => (end, token::StringDelimiter::SingleLine),
      (_end, r#""""#) => (start + '"', token::StringDelimiter::SingleLine),
      (_end, r#"""""#) | (_end, r#""""""""#) => (self.advance_by(2).unwrap(), token::StringDelimiter::MultiLine),

      (end, _) => {
        return Err(pos::spanned(start, end, Error::UnterminatedStringLiteral));
      }
    };

    self.context.map(|c| self.context_stack.push(c));
    self.context = Some(Context::String(delimiter));

    Ok(pos::spanned(start, end, Token::StringDelimiter(delimiter)))
  }

  fn heredoc_start(&mut self, _start: Location) -> Result<SpannedToken<'input>, SpannedError> {
    todo!()
  }

  // Utilities /////////////////////////

  fn bump(&mut self) -> Option<(Location, Location, char)> {
    let prev_loc = self.loc;
    self.chars.next().map(|ch| {
      self.loc.shift(ch);
      self.peek_loc = self.loc;
      (prev_loc, self.loc, ch)
    })
  }

  fn advance_by(&mut self, n: usize) -> Result<Location, usize> {
    let mut res = Ok(self.loc);
    for i in 0..n {
      res = self.bump().map(|(_, l, _)| l).ok_or(i);
    }
    res
  }

  fn peek(&mut self) -> Option<(Location, char)> {
    self.chars.peek().cloned().map(|ch| {
      self.peek_loc.shift(ch);
      (self.peek_loc, ch)
    })
  }

  fn reset_peek(&mut self) {
    self.chars.reset_peek();
    self.peek_loc = self.loc;
  }

  fn peek_while<F>(&mut self, start: Location, mut keep_going: F) -> (Location, &'input str)
  where
    F: FnMut(char) -> bool,
  {
    self.peek_until(start, |c| !keep_going(c))
  }

  fn peek_until<F>(&mut self, start: Location, mut terminate: F) -> (Location, &'input str)
  where
    F: FnMut(char) -> bool,
  {
    let mut end = self.peek_loc;
    while let Some((l, ch)) = self.peek() {
      if terminate(ch) {
        return (end, self.slice(start, end));
      }
      end = l;
    }

    (end, self.slice(start, end))
  }

  fn next_loc(&mut self) -> Option<Location> {
    let loc = self.peek().map(|(l, _)| l);
    self.reset_peek();
    loc
  }

  fn slice(&self, start: Location, end: Location) -> &'input str {
    let start = start.absolute.to_usize();
    let end = end.absolute.to_usize();
    &self.input[start..end]
  }
}

#[cfg(test)]
mod test {
  use super::*;
  // use crate::token::StringDelimiter;

  fn lexer(input: &str) -> Lexer {
    Lexer::new(input)
  }

  #[test]
  fn bump() {
    let input = "こんにちは";
    let start = 2;
    let end = 3;
    let mut lexer = lexer(input);

    assert_eq!(Some((Location::new(0, 0, 0), Location::new(0, 1, 3), 'こ')), lexer.bump());
    assert_eq!(Some((Location::new(0, 1, 3), Location::new(0, 2, 6), 'ん')), lexer.bump());
    assert_eq!(Some((Location::new(0, 2, 6), Location::new(0, 3, 9), 'に')), lexer.bump());
  }

  #[test]
  fn slice() {
    let input = "こんにちは";
    let start = 2;
    let end = 3;
    let mut lexer = lexer(input);

    lexer.bump();
    lexer.bump();

    if let Some((start, end, _)) = lexer.bump() {
      let res = lexer.slice(start, end);
      assert_eq!(res, "に");
    } else {
      panic!()
    }
  }

  #[test]
  fn peek() {
    let input = "##foo##";
    let mut lexer = lexer(input);

    let res: Option<(Location, char)> = lexer.peek();

    assert_eq!(res, Some((Location::new(0, 1, 1), '#')))
  }

  #[test]
  fn peek_until() {
    let input = "##foo##";
    let mut lexer = lexer(input);

    let res = lexer.peek_until(Location::new(0, 0, 0), |c| c != '#');

    assert_eq!(res, (Location::new(0, 2, 2), "##"))
  }

  #[test]
  fn peek_while() {
    let input = "##foo##";
    let mut lexer = lexer(input);

    let res = lexer.peek_while(Location::new(0, 0, 0), |c| c == '#');

    assert_eq!(res, (Location::new(0, 2, 2), "##"))
  }

  #[test]
  fn string_start() {
    let tests: Vec<(&str, Result<SpannedToken, SpannedError>, Option<Context<&str>>)> = vec![
      (
        r#"~"foo"~"#,
        Ok(pos::spanned(
          Location::new(0, 0, 0),
          Location::new(0, 1, 1),
          Token::StringDelimiter(token::StringDelimiter::SingleLine),
        )),
        Some(Context::String(token::StringDelimiter::SingleLine)),
      ),
      (
        r#"~""~"#,
        Ok(pos::spanned(
          Location::new(0, 0, 0),
          Location::new(0, 1, 1),
          Token::StringDelimiter(token::StringDelimiter::SingleLine),
        )),
        Some(Context::String(token::StringDelimiter::SingleLine)),
      ),
      (
        r#"~"~"#,
        Ok(pos::spanned(
          Location::new(0, 0, 0),
          Location::new(0, 1, 1),
          Token::StringDelimiter(token::StringDelimiter::SingleLine),
        )),
        Some(Context::String(token::StringDelimiter::SingleLine)),
      ),
      (
        r#"~"""foo"""~"#,
        Ok(pos::spanned(
          Location::new(0, 0, 0),
          Location::new(0, 3, 3),
          Token::StringDelimiter(token::StringDelimiter::MultiLine),
        )),
        Some(Context::String(token::StringDelimiter::MultiLine)),
      ),
      (
        r#"~"""~"""~"#,
        Ok(pos::spanned(
          Location::new(0, 0, 0),
          Location::new(0, 3, 3),
          Token::StringDelimiter(token::StringDelimiter::MultiLine),
        )),
        Some(Context::String(token::StringDelimiter::MultiLine)),
      ),
      (
        r#"~"""~"#,
        Ok(pos::spanned(
          Location::new(0, 0, 0),
          Location::new(0, 3, 3),
          Token::StringDelimiter(token::StringDelimiter::MultiLine),
        )),
        Some(Context::String(token::StringDelimiter::MultiLine)),
      ),
      (
        r#"~""""~"#,
        Err(pos::spanned(
          Location::new(0, 0, 0),
          Location::new(0, 4, 4),
          Error::UnterminatedStringLiteral,
        )),
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
