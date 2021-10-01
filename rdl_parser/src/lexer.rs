use std::str::Chars;

use crate::core::error::Errors;
use crate::pos::{self, Location, Spanned};
use crate::token::Token;
use crate::ParserSource;
use itertools::{Itertools, MultiPeek};

quick_error! {
  #[derive(Clone, Debug, PartialEq, Eq, Hash)]
  pub enum Error {
    UnterminatedHeredocDelimiter {
      display("unterminated heredoc delimiter")
    }
    MissingHeredocDelimiter {
      display("missing heredoc delimiter")
    }
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
enum Context<'input> {
  General,
  String(Token<&'input str>),
  HereDocString {
    delimiter: &'input str,
    skip_leading_tabs: bool,
    quoted_delimiter: bool,
  },
}

pub struct Lexer<'input> {
  input: &'input str,
  // chars: MultiPeek<CharIndices<'input>>,
  chars: MultiPeek<Chars<'input>>,
  loc: Location,
  peek_loc: Location,

  ctx: Option<Context<'input>>,
  ctx_stack: Vec<Context<'input>>,

  pub errors: Errors<SpannedError>,
}

impl<'input> Iterator for Lexer<'input> {
  type Item = LexerResult<'input>;

  fn next(&mut self) -> Option<LexerResult<'input>> {
    match *self.context() {
      Context::General => self.iter_base(),
      Context::String(Token::SingleLineStringDelimiter) => self.iter_single_line_string(),
      Context::String(Token::MultiLineStringDelimiter) => self.iter_multi_line_string(),
      Context::HereDocString {
        delimiter,
        skip_leading_tabs,
        quoted_delimiter,
      } => self.iter_heredoc_string(delimiter, skip_leading_tabs, quoted_delimiter),
      _ => unreachable!(),
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

      ctx: None,
      ctx_stack: Vec::new(),

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

  fn iter_string(&self, _string_type: &Token<&'input str>) -> Option<LexerResult<'input>> {
    todo!()
  }

  fn iter_single_line_string(&mut self) -> Option<LexerResult<'input>> {
    todo!()
  }

  fn iter_multi_line_string(&mut self) -> Option<LexerResult<'input>> {
    todo!()
  }

  fn iter_heredoc_string(&mut self, _delimiter: &str, _skip_leading_tabs: bool, _quoted_delimiter: bool) -> Option<LexerResult<'input>> {
    todo!()
  }

  // Token Lexers //////////////////////

  fn string_start(&mut self, start: Location) -> Result<SpannedToken<'input>, SpannedError> {
    let section = self.peek_while(start, |c| c == '"');
    self.reset_peek();

    let (end, token) = match section {
      (end, r#"""#) => (end, Token::SingleLineStringDelimiter),
      (_end, r#""""#) => (start + '"', Token::SingleLineStringDelimiter),
      (_end, r#"""""#) | (_end, r#""""""""#) => (self.advance_by(2).unwrap(), Token::MultiLineStringDelimiter),

      (end, _) => {
        return Err(pos::spanned(start, end, Error::UnterminatedStringLiteral));
      }
    };

    self.push_context(Context::String(token));

    Ok(pos::spanned(start, end, token))
  }

  fn heredoc_start(&mut self, start: Location) -> Result<SpannedToken<'input>, SpannedError> {
    if !self.test_peek(|ch| ch == '<') {
      self.reset_peek();
      println!("not a heredoc");
      todo!()
    }

    // Skip the second '<' of the redirection operator.
    self.catchup();

    // Check for a minus sign after the redirection operator.
    let skip_leading_tabs = self.test_peek(|ch| ch == '-');
    if skip_leading_tabs {
      self.catchup();
    }

    // Skip the possible whitespaces before the delimiter.
    self.bump_while(|ch| ch.is_whitespace() && ch != '\n');

    // Check for an opening quote before the delimiter.
    let quote_char = if self.test_peek(|ch| ch == '"' || ch == '\'') {
      self.bump().map(|(_, _, q)| q)
    } else {
      None
    };

    // Read until the end of the delimiter.
    let (mut end, delimiter) = self.take_while(self.loc, |ch| ch != '\n' && quote_char.map_or(true, |q| ch != q));
    if delimiter.is_empty() {
      return Err(pos::spanned(start, end, Error::MissingHeredocDelimiter));
    }

    // If we have a quoted delimiter, check that the closing and the opening quotes match.
    if let Some(q) = quote_char {
      match self.bump() {
        Some((_, l, ch)) if ch == q => {
          end = l;
        }
        _ => {
          return Err(pos::spanned(start, end, Error::UnterminatedHeredocDelimiter));
        }
      }
    }

    // Push the new lexing context.
    self.push_context(Context::HereDocString {
      delimiter,
      skip_leading_tabs,
      quoted_delimiter: quote_char.is_some(),
    });

    // Return the spanned token.
    Ok(pos::spanned(start, end, Token::HereDocStringDelimiter))
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

  fn take_while<F: FnMut(char) -> bool>(&mut self, start: Location, keep_going: F) -> (Location, &'input str) {
    let end = self.bump_while(keep_going);
    (end, self.slice(start, end))
  }

  fn take_until<F: FnMut(char) -> bool>(&mut self, start: Location, terminate: F) -> (Location, &'input str) {
    let end = self.bump_until(terminate);
    (end, self.slice(start, end))
  }

  fn bump_while<F: FnMut(char) -> bool>(&mut self, mut keep_going: F) -> Location {
    self.bump_until(|c| !keep_going(c))
  }

  fn bump_until<F: FnMut(char) -> bool>(&mut self, mut terminate: F) -> Location {
    let mut end = self.loc;
    self.reset_peek();

    while let Some((l, ch)) = self.peek() {
      if terminate(ch) {
        self.reset_peek();
        return end;
      }

      end = self.bump().unwrap().1;
    }

    end
  }

  fn advance_by(&mut self, n: usize) -> Result<Location, usize> {
    let mut res = Ok(self.loc);
    for i in 0..n {
      res = self.bump().map(|(_, l, _)| l).ok_or(i);
    }
    res
  }

  fn catchup(&mut self) {
    while self.loc < self.peek_loc {
      self.chars.next().map(|ch| self.loc.shift(ch));
    }
  }

  fn peek(&mut self) -> Option<(Location, char)> {
    self.chars.peek().cloned().map(|ch| {
      self.peek_loc.shift(ch);
      (self.peek_loc, ch)
    })
  }

  fn test_peek<F: FnMut(char) -> bool>(&mut self, mut test: F) -> bool {
    self.peek().map_or(false, |(_, ch)| test(ch))
  }

  fn reset_peek(&mut self) {
    self.chars.reset_peek();
    self.peek_loc = self.loc;
  }

  fn peek_while<F: FnMut(char) -> bool>(&mut self, start: Location, mut keep_going: F) -> (Location, &'input str) {
    self.peek_until(start, |c| !keep_going(c))
  }

  fn peek_until<F: FnMut(char) -> bool>(&mut self, start: Location, mut terminate: F) -> (Location, &'input str) {
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

  // Context manipulation //////////////

  fn context(&mut self) -> &Context<'input> {
    self.ctx.as_ref().unwrap_or(&Context::General)
  }

  fn push_context(&mut self, context: Context<'input>) {
    self.ctx.map(|c| self.ctx_stack.push(c));
    self.ctx = Some(context);
  }

  fn pop_context(&mut self) -> Option<Context> {
    let ctx = self.ctx;
    self.ctx = self.ctx_stack.pop();
    ctx
  }
}

#[cfg(test)]
mod test {
  use super::*;

  fn loc(pos: u32) -> Location {
    Location::new(0, pos, pos)
  }

  #[test]
  fn bump() {
    let input = "こんにちは";
    let mut lexer = Lexer::new(input);

    assert_eq!(Some((Location::new(0, 0, 0), Location::new(0, 1, 3), 'こ')), lexer.bump());
    assert_eq!(Some((Location::new(0, 1, 3), Location::new(0, 2, 6), 'ん')), lexer.bump());
    assert_eq!(Some((Location::new(0, 2, 6), Location::new(0, 3, 9), 'に')), lexer.bump());
  }

  #[test]
  fn slice() {
    let input = "こんにちは";
    let mut lexer = Lexer::new(input);

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
  fn take_until() {
    let input = "##foo##";
    let mut lexer = Lexer::new(input);

    let res = lexer.take_until(loc(0), |c| c != '#');

    assert_eq!(res, (loc(2), "##"))
  }

  #[test]
  fn peek() {
    let input = "##foo##";
    let mut lexer = Lexer::new(input);

    let res: Option<(Location, char)> = lexer.peek();

    assert_eq!(res, Some((loc(1), '#')))
  }

  #[test]
  fn peek_while() {
    let input = "##foo##";
    let mut lexer = Lexer::new(input);

    let res = lexer.peek_while(loc(0), |c| c == '#');

    assert_eq!(res, (loc(2), "##"))
  }

  #[test]
  fn peek_until() {
    let input = "##foo##";
    let mut lexer = Lexer::new(input);

    let res = lexer.peek_until(loc(0), |c| c != '#');

    assert_eq!(res, (loc(2), "##"))
  }

  type TestCase<'a> = (&'a str, Result<SpannedToken<'a>, SpannedError>, Context<'a>);

  fn test(tests: Vec<TestCase>) {
    for (input, expected, expected_context) in tests {
      let s = input.replace("~", "");
      let input = &*s;
      let mut lexer = Lexer::new(input);

      let res = lexer.next().unwrap();

      assert_eq!(expected, res);
      assert_eq!(&expected_context, lexer.context());
    }
  }

  #[test]
  fn string_start() {
    let tests = vec![
      (
        r#"~"foo"~"#,
        Ok(pos::spanned(loc(0), loc(1), Token::SingleLineStringDelimiter)),
        Context::String(Token::SingleLineStringDelimiter),
      ),
      (
        r#"~""~"#,
        Ok(pos::spanned(loc(0), loc(1), Token::SingleLineStringDelimiter)),
        Context::String(Token::SingleLineStringDelimiter),
      ),
      (
        r#"~"~"#,
        Ok(pos::spanned(loc(0), loc(1), Token::SingleLineStringDelimiter)),
        Context::String(Token::SingleLineStringDelimiter),
      ),
      (
        r#"~"""foo"""~"#,
        Ok(pos::spanned(loc(0), loc(3), Token::MultiLineStringDelimiter)),
        Context::String(Token::MultiLineStringDelimiter),
      ),
      (
        r#"~"""~"""~"#,
        Ok(pos::spanned(loc(0), loc(3), Token::MultiLineStringDelimiter)),
        Context::String(Token::MultiLineStringDelimiter),
      ),
      (
        r#"~"""~"#,
        Ok(pos::spanned(loc(0), loc(3), Token::MultiLineStringDelimiter)),
        Context::String(Token::MultiLineStringDelimiter),
      ),
      (
        r#"~""""~"#,
        Err(pos::spanned(loc(0), loc(4), Error::UnterminatedStringLiteral)),
        Context::General,
      ),
    ];

    test(tests)
  }

  #[test]
  fn heredoc_start() {
    let tests = vec![
      (
        "<< EOF\nfoo\nEOF",
        Ok(pos::spanned(loc(0), loc(6), Token::HereDocStringDelimiter)),
        Context::HereDocString {
          delimiter: "EOF",
          skip_leading_tabs: false,
          quoted_delimiter: false,
        },
      ),
      (
        "<<EOF\nfoo\nEOF",
        Ok(pos::spanned(loc(0), loc(5), Token::HereDocStringDelimiter)),
        Context::HereDocString {
          delimiter: "EOF",
          skip_leading_tabs: false,
          quoted_delimiter: false,
        },
      ),
      (
        "<<   EOF\nfoo\nEOF",
        Ok(pos::spanned(loc(0), loc(8), Token::HereDocStringDelimiter)),
        Context::HereDocString {
          delimiter: "EOF",
          skip_leading_tabs: false,
          quoted_delimiter: false,
        },
      ),
      (
        "<<- EOF\nfoo\nEOF",
        Ok(pos::spanned(loc(0), loc(7), Token::HereDocStringDelimiter)),
        Context::HereDocString {
          delimiter: "EOF",
          skip_leading_tabs: true,
          quoted_delimiter: false,
        },
      ),
      (
        "<<-EOF\nfoo\nEOF",
        Ok(pos::spanned(loc(0), loc(6), Token::HereDocStringDelimiter)),
        Context::HereDocString {
          delimiter: "EOF",
          skip_leading_tabs: true,
          quoted_delimiter: false,
        },
      ),
      (
        "<<-   EOF\nfoo\nEOF",
        Ok(pos::spanned(loc(0), loc(9), Token::HereDocStringDelimiter)),
        Context::HereDocString {
          delimiter: "EOF",
          skip_leading_tabs: true,
          quoted_delimiter: false,
        },
      ),
      (
        "<< - EOF\nfoo\nEOF",
        Ok(pos::spanned(loc(0), loc(8), Token::HereDocStringDelimiter)),
        Context::HereDocString {
          delimiter: "- EOF",
          skip_leading_tabs: false,
          quoted_delimiter: false,
        },
      ),
      (
        "<<   - EOF\nfoo\nEOF",
        Ok(pos::spanned(loc(0), loc(10), Token::HereDocStringDelimiter)),
        Context::HereDocString {
          delimiter: "- EOF",
          skip_leading_tabs: false,
          quoted_delimiter: false,
        },
      ),
      (
        "<< -   EOF\nfoo\nEOF",
        Ok(pos::spanned(loc(0), loc(10), Token::HereDocStringDelimiter)),
        Context::HereDocString {
          delimiter: "-   EOF",
          skip_leading_tabs: false,
          quoted_delimiter: false,
        },
      ),
      (
        "<< \"EOF\"\nfoo\nEOF",
        Ok(pos::spanned(loc(0), loc(8), Token::HereDocStringDelimiter)),
        Context::HereDocString {
          delimiter: "EOF",
          skip_leading_tabs: false,
          quoted_delimiter: true,
        },
      ),
      (
        "<<\"EOF\"\nfoo\nEOF",
        Ok(pos::spanned(loc(0), loc(7), Token::HereDocStringDelimiter)),
        Context::HereDocString {
          delimiter: "EOF",
          skip_leading_tabs: false,
          quoted_delimiter: true,
        },
      ),
      (
        "<<   \"EOF\"\nfoo\nEOF",
        Ok(pos::spanned(loc(0), loc(10), Token::HereDocStringDelimiter)),
        Context::HereDocString {
          delimiter: "EOF",
          skip_leading_tabs: false,
          quoted_delimiter: true,
        },
      ),
      (
        "<<- \"EOF\"\nfoo\nEOF",
        Ok(pos::spanned(loc(0), loc(9), Token::HereDocStringDelimiter)),
        Context::HereDocString {
          delimiter: "EOF",
          skip_leading_tabs: true,
          quoted_delimiter: true,
        },
      ),
      (
        "<<-\"EOF\"\nfoo\nEOF",
        Ok(pos::spanned(loc(0), loc(8), Token::HereDocStringDelimiter)),
        Context::HereDocString {
          delimiter: "EOF",
          skip_leading_tabs: true,
          quoted_delimiter: true,
        },
      ),
      (
        "<<-   \"EOF\"\nfoo\nEOF",
        Ok(pos::spanned(loc(0), loc(11), Token::HereDocStringDelimiter)),
        Context::HereDocString {
          delimiter: "EOF",
          skip_leading_tabs: true,
          quoted_delimiter: true,
        },
      ),
      (
        "<<\nfoo\nEOF",
        Err(pos::spanned(loc(0), loc(2), Error::MissingHeredocDelimiter)),
        Context::General,
      ),
      (
        "<< \"FOO\nfoo\nEOF",
        Err(pos::spanned(loc(0), loc(7), Error::UnterminatedHeredocDelimiter)),
        Context::General,
      ),
      (
        "<< 'FOO\nfoo\nEOF",
        Err(pos::spanned(loc(0), loc(7), Error::UnterminatedHeredocDelimiter)),
        Context::General,
      ),
      (
        "<< \"FOO'\nfoo\nEOF",
        Err(pos::spanned(loc(0), loc(8), Error::UnterminatedHeredocDelimiter)),
        Context::General,
      ),
      (
        "<< 'FOO\"\nfoo\nEOF",
        Err(pos::spanned(loc(0), loc(8), Error::UnterminatedHeredocDelimiter)),
        Context::General,
      ),
    ];

    test(tests)
  }
}
