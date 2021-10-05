use std::str::Chars;

use crate::core::error::Errors;
use crate::pos::{self, Location, Spanned};
use crate::token::Token;
use crate::ParserSource;
use itertools::{Itertools, MultiPeek};
use ordered_float::NotNan;

quick_error! {
  #[derive(Clone, Debug, PartialEq, Eq)]
  pub enum Error {
    NonParseableInt(err: std::num::ParseIntError) {
      display("cannot parse integer")
    }
    NonParseableFloat(err: std::num::ParseFloatError) {
      display("cannot parse float")
    }
    LiteralIncomplete {
      display("cannot parse literal, incomplete")
    }
    HexLiteralOverflow {
      display("cannot parse hex literal, overflow")
    }
    HexLiteralUnderflow {
      display("cannot parse hex literal, underflow")
    }
    HexLiteralWrongPrefix {
      display("wrong hex literal prefix, should start as '0x' or '-0x'")
    }
    BinLiteralWrongPrefix {
      display("wrong bin literal prefix, should start as '0b' or '-0b'")
    }
    OctLiteralWrongPrefix {
      display("wrong oct literal prefix, should start as '0o' or '-0o'")
    }
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
    while let Some((start, end, ch)) = self.bump() {
      return match ch {
        ',' => Some(Ok(pos::spanned(start, end, Token::Comma))),

        '{' => Some(Ok(pos::spanned(start, end, Token::LBrace))),
        '[' => Some(Ok(pos::spanned(start, end, Token::LBracket))),
        '(' => Some(Ok(pos::spanned(start, end, Token::LParen))),
        '}' => Some(Ok(pos::spanned(start, end, Token::RBrace))),
        ']' => Some(Ok(pos::spanned(start, end, Token::RBracket))),
        ')' => Some(Ok(pos::spanned(start, end, Token::RParen))),

        '/' if self.peek_is('/') => Some(self.line_comment(start)),
        '/' if self.peek_is('*') => Some(self.block_comment(start)),

        '"' => Some(self.string_start(start)),
        '<' if self.peek_is('<') => Some(self.heredoc_start(start)), // :>

        ch if is_dec(ch) || (ch == '-' && self.test_peek_reset(is_dec)) || (ch == '+' && self.test_peek_reset(is_dec)) => {
          Some(self.numeric_literal(start))
        }

        ch if is_operator_char(ch) => Some(self.operator(start)),

        _ => None,
      };
    }

    // TODO: this shouldn't be necessary
    let next_loc = self.next_loc()?;
    Some(Ok(pos::spanned(next_loc, next_loc, Token::EOF)))
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
    // Skip the second '<' of the redirection operator.
    self.bump();

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

  fn numeric_literal(&mut self, start: Location) -> Result<SpannedToken<'input>, SpannedError> {
    let (end, int) = self.take_while(start, |ch| is_dec(ch) || ch == '_');

    match self.peek() {
      Some((_, '.')) => {
        self.catchup();
        let (end, float) = self.take_while(start, is_dec);

        let (end, float) = if self.test_peek(|ch| ch == 'e' || ch == 'E') {
          self.catchup();

          if self.test_peek(|ch| ch == '-' || ch == '+') {
            self.catchup();
          }

          self.take_while(start, |ch| is_dec(ch) || ch == '_')
        } else {
          (end, float)
        };

        self.parse_float(start, end, float)
      }
      Some((_, 'e')) | Some((_, 'E')) => {
        self.catchup();

        if self.test_peek(|ch| ch == '-' || ch == '+') {
          self.catchup();
        }

        let (end, float) = self.take_while(start, |ch| is_dec(ch) || ch == '_');
        self.parse_float(start, end, float)
      }

      Some((_, 'x')) => {
        let hex_start = self.catchup();
        let (end, hex) = self.take_while(hex_start, |ch| is_hex(ch) || ch == '_');

        if int != "0" && int != "-0" {
          self.errors.push(pos::spanned(start, end, Error::HexLiteralWrongPrefix));
        }

        if hex.is_empty() {
          return self.recover(start, end, Error::LiteralIncomplete, Token::IntLiteral(0));
        }

        let is_positive = int == "0";
        match i64_from_str_radix(hex, is_positive, 16) {
          Ok(val) => Ok(pos::spanned(start, end, Token::IntLiteral(val))),
          Err(err) => self.recover(start, end, err, Token::IntLiteral(0)),
        }
      }

      Some((_, 'b')) => {
        let bin_start = self.catchup();
        let (end, bin) = self.take_while(bin_start, |ch| is_bin(ch) || ch == '_');

        if int != "0" && int != "-0" {
          self.errors.push(pos::spanned(start, end, Error::BinLiteralWrongPrefix));
        }

        if bin.is_empty() {
          return self.recover(start, end, Error::LiteralIncomplete, Token::IntLiteral(0));
        }

        let is_positive = int == "0";
        match i64_from_str_radix(bin, is_positive, 2) {
          Ok(val) => Ok(pos::spanned(start, end, Token::IntLiteral(val))),
          Err(err) => self.recover(start, end, err, Token::IntLiteral(0)),
        }
      }

      Some((_, 'o')) => {
        let oct_start = self.catchup();
        let (end, oct) = self.take_while(oct_start, |ch| is_oct(ch) || ch == '_');

        if int != "0" && int != "-0" {
          self.errors.push(pos::spanned(start, end, Error::OctLiteralWrongPrefix));
        }

        if oct.is_empty() {
          return self.recover(start, end, Error::LiteralIncomplete, Token::IntLiteral(0));
        }

        let is_positive = int == "0";
        match i64_from_str_radix(oct, is_positive, 8) {
          Ok(val) => Ok(pos::spanned(start, end, Token::IntLiteral(val))),
          Err(err) => self.recover(start, end, err, Token::IntLiteral(0)),
        }
      }

      None | Some(_) => self.parse_int(start, end, int),
    }
  }

  fn line_comment(&mut self, start: Location) -> Result<SpannedToken<'input>, SpannedError> {
    todo!()
  }

  fn block_comment(&mut self, start: Location) -> Result<SpannedToken<'input>, SpannedError> {
    todo!()
  }

  fn operator(&mut self, start: Location) -> Result<SpannedToken<'input>, SpannedError> {
    let (end, op) = self.take_while(start, is_operator_char);

    let token = match op {
      "!" => Token::LogicalNot,
      "&&" => Token::LogicalAnd,
      "||" => Token::LogicalOr,
      "~" => Token::BitwiseNot,
      "&" => Token::BitwiseAnd,
      "|" => Token::BitwiseOr,
      "<:" => Token::BitwiseShiftLeft,
      ":>" => Token::BitwiseShiftRight,
      "*" => Token::Multiplication,
      "/" => Token::Division,
      "%" => Token::Modulo,
      "+" => Token::Addition,
      "-" => Token::Subtraction,
      "==" => Token::Equal,
      "!=" => Token::NotEqual,
      "." => Token::Dot,
      ".." => Token::DotDot,
      "=" => Token::Assign,
      "->" => Token::RArrow,

      _ => todo!(),
    };

    Ok(pos::spanned(start, end, token))
  }

  fn parse_int(&mut self, start: Location, end: Location, v: &str) -> Result<SpannedToken<'input>, SpannedError> {
    let v = v.replace('_', "");
    match v.parse::<i64>() {
      Ok(val) => Ok(pos::spanned(start, end, Token::IntLiteral(val))),
      Err(e) => self.recover(start, end, Error::NonParseableInt(e), Token::IntLiteral(0)),
    }
  }

  fn parse_float(&mut self, start: Location, end: Location, v: &str) -> Result<SpannedToken<'input>, SpannedError> {
    let v = v.replace('_', "");
    match v.parse::<f64>() {
      Ok(val) => Ok(pos::spanned(start, end, Token::FloatLiteral(NotNan::new(val).unwrap()))),
      Err(e) => self.recover(
        start,
        end,
        Error::NonParseableFloat(e),
        Token::FloatLiteral(NotNan::new(0.0).unwrap()),
      ),
    }
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

    while let Some((_, ch)) = self.peek() {
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

  fn catchup(&mut self) -> Location {
    while self.loc < self.peek_loc {
      self.chars.next().map(|ch| self.loc.shift(ch));
    }
    self.loc
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

  fn test_peek_reset<F: FnMut(char) -> bool>(&mut self, mut test: F) -> bool {
    let r = self.peek().map_or(false, |(_, ch)| test(ch));
    self.reset_peek();
    r
  }

  fn peek_is(&mut self, ch: char) -> bool {
    let r = self.peek().map_or(false, |(_, c)| c == ch);
    self.reset_peek();
    r
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

  fn recover(
    &mut self,
    start: Location,
    end: Location,
    code: Error,
    value: Token<&'input str>,
  ) -> Result<SpannedToken<'input>, SpannedError> {
    self.errors.push(pos::spanned(start, end, code));
    Ok(pos::spanned(start, end, value))
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

fn is_dec(ch: char) -> bool {
  ch.is_digit(10)
}

fn is_hex(ch: char) -> bool {
  ch.is_digit(16)
}

fn is_bin(ch: char) -> bool {
  ch.is_digit(2)
}

fn is_oct(ch: char) -> bool {
  ch.is_digit(8)
}

fn is_operator_char(ch: char) -> bool {
  match ch {
    '!' => true,
    '%' => true,
    '&' => true,
    '*' => true,
    '+' => true,
    '-' => true,
    '.' => true,
    '/' => true,
    ':' => true,
    '<' => true,
    '=' => true,
    '>' => true,
    '|' => true,
    '~' => true,

    _ => false,
  }
}

fn i64_from_str_radix(hex: &str, is_positive: bool, radix: u32) -> Result<i64, Error> {
  let sign: i64 = if is_positive { 1 } else { -1 };
  let mut result = 0i64;

  for c in hex.chars() {
    if c == '_' {
      continue;
    }

    let x = c.to_digit(radix).expect("invalid literal");
    result = result
      .checked_mul(radix as i64)
      .and_then(|result| result.checked_add((x as i64) * sign))
      .ok_or_else(|| {
        if is_positive {
          Error::HexLiteralOverflow
        } else {
          Error::HexLiteralUnderflow
        }
      })?;
  }

  Ok(result)
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
      let s = input.replace("➖", "");
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
        r#"➖"foo"➖"#,
        Ok(pos::spanned(loc(0), loc(1), Token::SingleLineStringDelimiter)),
        Context::String(Token::SingleLineStringDelimiter),
      ),
      (
        r#"➖""➖"#,
        Ok(pos::spanned(loc(0), loc(1), Token::SingleLineStringDelimiter)),
        Context::String(Token::SingleLineStringDelimiter),
      ),
      (
        r#"➖"➖"#,
        Ok(pos::spanned(loc(0), loc(1), Token::SingleLineStringDelimiter)),
        Context::String(Token::SingleLineStringDelimiter),
      ),
      (
        r#"➖"""foo"""➖"#,
        Ok(pos::spanned(loc(0), loc(3), Token::MultiLineStringDelimiter)),
        Context::String(Token::MultiLineStringDelimiter),
      ),
      (
        r#"➖"""➖"""➖"#,
        Ok(pos::spanned(loc(0), loc(3), Token::MultiLineStringDelimiter)),
        Context::String(Token::MultiLineStringDelimiter),
      ),
      (
        r#"➖"""➖"#,
        Ok(pos::spanned(loc(0), loc(3), Token::MultiLineStringDelimiter)),
        Context::String(Token::MultiLineStringDelimiter),
      ),
      (
        r#"➖""""➖"#,
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

  #[test]
  fn numeric() {
    let tests = vec![
      ("1234", Ok(pos::spanned(loc(0), loc(4), Token::IntLiteral(1234))), Context::General),
      (
        "-1234",
        Ok(pos::spanned(loc(0), loc(5), Token::IntLiteral(-1234))),
        Context::General,
      ),
      ("12_34", Ok(pos::spanned(loc(0), loc(5), Token::IntLiteral(1234))), Context::General),
      (
        "12.34",
        Ok(pos::spanned(loc(0), loc(5), Token::FloatLiteral(NotNan::new(12.34).unwrap()))),
        Context::General,
      ),
      (
        "-12.34",
        Ok(pos::spanned(loc(0), loc(6), Token::FloatLiteral(NotNan::new(-12.34).unwrap()))),
        Context::General,
      ),
      (
        "123_456.0",
        Ok(pos::spanned(loc(0), loc(9), Token::FloatLiteral(NotNan::new(123456.0).unwrap()))),
        Context::General,
      ),
      (
        "12.34e5",
        Ok(pos::spanned(loc(0), loc(7), Token::FloatLiteral(NotNan::new(1234000.0).unwrap()))),
        Context::General,
      ),
      (
        "12.34e0_5",
        Ok(pos::spanned(loc(0), loc(9), Token::FloatLiteral(NotNan::new(1234000.0).unwrap()))),
        Context::General,
      ),
      (
        "1234e5",
        Ok(pos::spanned(loc(0), loc(6), Token::FloatLiteral(NotNan::new(123400000.0).unwrap()))),
        Context::General,
      ),
      (
        "1234e0_5",
        Ok(pos::spanned(loc(0), loc(8), Token::FloatLiteral(NotNan::new(123400000.0).unwrap()))),
        Context::General,
      ),
      (
        "1234e-2",
        Ok(pos::spanned(loc(0), loc(7), Token::FloatLiteral(NotNan::new(12.34).unwrap()))),
        Context::General,
      ),
      (
        "1234e+2",
        Ok(pos::spanned(loc(0), loc(7), Token::FloatLiteral(NotNan::new(123400.0).unwrap()))),
        Context::General,
      ),
      (
        "1234.56e-2",
        Ok(pos::spanned(loc(0), loc(10), Token::FloatLiteral(NotNan::new(12.3456).unwrap()))),
        Context::General,
      ),
      (
        "1234.56e+2",
        Ok(pos::spanned(loc(0), loc(10), Token::FloatLiteral(NotNan::new(123456.0).unwrap()))),
        Context::General,
      ),
      (
        "0x1234",
        Ok(pos::spanned(loc(0), loc(6), Token::IntLiteral(0x1234))),
        Context::General,
      ),
      (
        "-0x1234",
        Ok(pos::spanned(loc(0), loc(7), Token::IntLiteral(-0x1234))),
        Context::General,
      ),
      (
        "0x12_34",
        Ok(pos::spanned(loc(0), loc(7), Token::IntLiteral(0x1234))),
        Context::General,
      ),
      (
        "0b1010",
        Ok(pos::spanned(loc(0), loc(6), Token::IntLiteral(0b1010))),
        Context::General,
      ),
      (
        "-0b1010",
        Ok(pos::spanned(loc(0), loc(7), Token::IntLiteral(-0b1010))),
        Context::General,
      ),
      (
        "0b10_10",
        Ok(pos::spanned(loc(0), loc(7), Token::IntLiteral(0b1010))),
        Context::General,
      ),
      (
        "0o1234",
        Ok(pos::spanned(loc(0), loc(6), Token::IntLiteral(0o1234))),
        Context::General,
      ),
      (
        "-0o1234",
        Ok(pos::spanned(loc(0), loc(7), Token::IntLiteral(-0o1234))),
        Context::General,
      ),
      (
        "0o12_34",
        Ok(pos::spanned(loc(0), loc(7), Token::IntLiteral(0o1234))),
        Context::General,
      ),
    ];

    test(tests);
  }

  #[test]
  fn operators() {
    let tests = vec![
      ("!", Ok(pos::spanned(loc(0), loc(1), Token::LogicalNot)), Context::General),
      ("&&", Ok(pos::spanned(loc(0), loc(2), Token::LogicalAnd)), Context::General),
      ("||", Ok(pos::spanned(loc(0), loc(2), Token::LogicalOr)), Context::General),
      ("~", Ok(pos::spanned(loc(0), loc(1), Token::BitwiseNot)), Context::General),
      ("&", Ok(pos::spanned(loc(0), loc(1), Token::BitwiseAnd)), Context::General),
      ("|", Ok(pos::spanned(loc(0), loc(1), Token::BitwiseOr)), Context::General),
      ("<:", Ok(pos::spanned(loc(0), loc(2), Token::BitwiseShiftLeft)), Context::General),
      (":>", Ok(pos::spanned(loc(0), loc(2), Token::BitwiseShiftRight)), Context::General),
      ("*", Ok(pos::spanned(loc(0), loc(1), Token::Multiplication)), Context::General),
      ("/", Ok(pos::spanned(loc(0), loc(1), Token::Division)), Context::General),
      ("%", Ok(pos::spanned(loc(0), loc(1), Token::Modulo)), Context::General),
      ("+", Ok(pos::spanned(loc(0), loc(1), Token::Addition)), Context::General),
      ("-", Ok(pos::spanned(loc(0), loc(1), Token::Subtraction)), Context::General),
      ("==", Ok(pos::spanned(loc(0), loc(2), Token::Equal)), Context::General),
      ("!=", Ok(pos::spanned(loc(0), loc(2), Token::NotEqual)), Context::General),
      (".", Ok(pos::spanned(loc(0), loc(1), Token::Dot)), Context::General),
      ("..", Ok(pos::spanned(loc(0), loc(2), Token::DotDot)), Context::General),
      ("=", Ok(pos::spanned(loc(0), loc(1), Token::Assign)), Context::General),
      ("->", Ok(pos::spanned(loc(0), loc(2), Token::RArrow)), Context::General),

      (",", Ok(pos::spanned(loc(0), loc(1), Token::Comma)), Context::General),
      ("{", Ok(pos::spanned(loc(0), loc(1), Token::LBrace)), Context::General),
      ("[", Ok(pos::spanned(loc(0), loc(1), Token::LBracket)), Context::General),
      ("(", Ok(pos::spanned(loc(0), loc(1), Token::LParen)), Context::General),
      ("}", Ok(pos::spanned(loc(0), loc(1), Token::RBrace)), Context::General),
      ("]", Ok(pos::spanned(loc(0), loc(1), Token::RBracket)), Context::General),
      (")", Ok(pos::spanned(loc(0), loc(1), Token::RParen)), Context::General),
    ];

    test(tests);
  }
}
