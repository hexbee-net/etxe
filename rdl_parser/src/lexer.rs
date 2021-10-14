use std::str::Chars;

use crate::core::error::Errors;
use crate::pos::{self, Location, Spanned};
use crate::token::{Comment, CommentType, StringLiteral, Token};
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
    EmptyCharLiteral {
      display("empty char literal")
    }
    UnterminatedCharLiteral {
      display("unterminated character literal")
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
    EmptyUnicodeEscape {
      display("a unicode escape must have at least 1 hex digit")
    }
    MalformedUnicodeEscape {
      display("format of unicode escape sequences uses braces")
    }
    InvalidUnicodeEscape {
      display("invalid character in unicode escape")
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

  String,
  HereDocString {
    delimiter: &'input str,
    skip_leading_tabs: bool,
    quoted_delimiter: bool,
  },
  HereDocStringEnd {
    start: Location,
    end: Location,
  },

  StringInterpolation,
  StringDirective,
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
      Context::String => self.iter_single_line_string(),
      Context::HereDocString {
        delimiter,
        skip_leading_tabs,
        quoted_delimiter,
      } => self.iter_heredoc_string(delimiter, skip_leading_tabs, quoted_delimiter),
      Context::HereDocStringEnd { start, end } => self.iter_heredoc_end(start, end),
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

        '/' if self.test_peek_reset(|ch| ch == '/') => Some(self.line_comment(start)),
        '/' if self.test_peek_reset(|ch| ch == '*') => Some(self.block_comment(start)),

        '"' => {
          self.push_context(Context::String);
          Some(Ok(pos::spanned(start, end, Token::StringDelimiter)))
        }
        '<' if self.test_peek_reset(|ch| ch == '<') => Some(self.heredoc_start(start)),
        '\'' => Some(self.char_literal(start)),

        ch if is_dec(ch) || (ch == '-' && self.test_peek_reset(is_dec)) || (ch == '+' && self.test_peek_reset(is_dec)) => {
          Some(self.numeric_literal(start))
        }

        ch if is_operator_char(ch) => Some(self.operator(start)),

        ch if is_ident_start(ch) => Some(self.identifier(start)),

        ch if ch.is_whitespace() && ch != '\n' => continue,

        _ => None,
      };
    }

    let next_loc = self.next_loc()?;
    Some(Ok(pos::spanned(next_loc, next_loc, Token::EOF)))
  }

  fn iter_single_line_string(&mut self) -> Option<LexerResult<'input>> {
    let start = self.loc;

    if self.test_peek_reset(|ch| ch == '"') {
      let end = self.catchup();
      self.pop_context();

      return Some(Ok(pos::spanned(start, end, Token::StringDelimiter)));
    }

    match (self.peek(), self.peek()) {
      (Some((_, '$')), Some((end, '{'))) => {
        self.catchup();
        self.push_context(Context::StringInterpolation);
        return Some(Ok(pos::spanned(start, end, Token::StringInterpolation)));
      }
      (Some((_, '%')), Some((end, '{'))) => {
        self.catchup();
        self.push_context(Context::StringDirective);
        return Some(Ok(pos::spanned(start, end, Token::StringDirective)));
      }
      _ => {}
    };
    self.reset_peek();

    let mut string_end = start;
    while let Some((_, ch)) = self.peek() {
      if ch == '"' {
        self.reset_peek();
        return Some(Ok(pos::spanned(
          start,
          string_end,
          Token::String(StringLiteral::Escaped(self.slice(start, string_end))),
        )));
      }

      let next = match self.peek() {
        Some((_, ch)) => ch,
        _ => {
          return Some(Err(pos::spanned(start, self.loc, Error::UnexpectedEof)));
        }
      };

      match ch {
        '$' if next == '$' => {
          self.bump();
          continue;
        }
        '%' if next == '%' => {
          self.bump();
          continue;
        }
        '\\' if next == '\\' => {
          self.bump();
          continue;
        }

        '\\' => {
          match self.parse_escape_code(start, '"') {
            Err(err) => self.errors.push(err),
            Ok(_) => {}
          }

          continue;
        }

        '$' if next == '{' => {
          self.reset_peek();
          return Some(Ok(pos::spanned(
            start,
            string_end,
            Token::String(StringLiteral::Escaped(self.slice(start, string_end))),
          )));
        }

        '%' if next == '{' => {
          self.reset_peek();
          return Some(Ok(pos::spanned(
            start,
            string_end,
            Token::String(StringLiteral::Escaped(self.slice(start, string_end))),
          )));
        }

        _ => {
          string_end = self.bump().unwrap().1;
          continue;
        }
      }
    }

    Some(Err(pos::spanned(start, self.loc, Error::UnexpectedEof)))
  }

  fn iter_heredoc_string(&mut self, delimiter: &str, skip_leading_tabs: bool, quoted_delimiter: bool) -> Option<LexerResult<'input>> {
    let start = self.loc;

    let mut tab_size = None;
    let mut line_start = self.loc;
    let mut line_end = self.loc;
    let mut lines = Vec::new();
    let end;

    loop {
      let prev_loc = self.loc;
      match self.bump() {
        Some((_, loc, '\n')) => {
          let line = self.slice(line_start, prev_loc);
          line_end = prev_loc;

          if line == delimiter {
            end = loc;
            break;
          }

          // Keep track of leading tabs count if needed.
          if skip_leading_tabs && !line.is_empty() {
            let tab_count = line.find(|ch| !char::is_whitespace(ch)).unwrap_or(0);
            if tab_size.is_none() || tab_count < tab_size.unwrap() {
              tab_size = Some(tab_count);
            }
          }

          lines.push(line);
          line_start = loc;
        }

        // "$${" escapes to literal "${".
        Some((_, _, '$')) if !quoted_delimiter && self.test_peek_reset(|ch| ch == '$') => {
          self.bump();
        }

        // "%%{" escapes to literal "%{".
        Some((_, _, '%')) if !quoted_delimiter && self.test_peek_reset(|ch| ch == '%') => {
          self.bump();
        }

        // Interpolation start.
        Some((_, _loc, '$')) if !quoted_delimiter && self.test_peek_reset(|ch| ch == '{') => {
          todo!();
        }

        // Directive start.
        Some((_, _loc, '%')) if !quoted_delimiter && self.test_peek_reset(|ch| ch == '{') => {
          todo!();
        }

        Some((_, _, _)) => {}

        // We reached end of file, check if the last line is the delimiter.
        None => {
          if line_start != prev_loc {
            if self.slice(line_start, prev_loc) == delimiter {
              end = prev_loc;
              break;
            }
          }

          return Some(Err(pos::spanned(start, prev_loc, Error::UnexpectedEof)));
        }
      }
    }

    self.push_context(Context::HereDocStringEnd { start: line_start, end });

    if lines.is_empty() {
      return Some(Ok(pos::spanned(start, end, Token::HeredocString(lines))));
    }

    // Remove leading tabs if needed.
    tab_size.map(|tab_size| {
      for i in 0..lines.len() {
        if !lines[i].is_empty() {
          lines[i] = &lines[i][tab_size..];
        }
      }
    });

    Some(Ok(pos::spanned(start, line_end, Token::HeredocString(lines))))
  }

  fn iter_heredoc_end(&mut self, start: Location, end: Location) -> Option<LexerResult<'input>> {
    self.pop_context();

    if let Some(Context::HereDocString { .. }) = self.pop_context() {
    } else {
      panic!("invalid context stack");
    }

    Some(Ok(pos::spanned(start, end, Token::StringDelimiter)))
  }

  // Token Lexers //////////////////////

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
    self.bump();

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
    Ok(pos::spanned(start, end, Token::StringDelimiter))
  }

  fn char_literal(&mut self, start: Location) -> Result<SpannedToken<'input>, SpannedError> {
    let ch = match self.bump() {
      Some((_, _, '\\')) => self.parse_escape_code(start, '\'').unwrap_or_else(|e| {
        self.errors.push(e);
        self.bump_until(|ch| ch == '\'');
        '\0'
      }),
      Some((_, end, '\'')) => return self.recover(start, end, Error::EmptyCharLiteral, Token::CharLiteral('\0')),
      Some((_, _, ch)) => ch,
      None => return self.recover(start, start, Error::UnexpectedEof, Token::CharLiteral('\0')),
    };

    match self.bump() {
      Some((_, end, '\'')) => Ok(pos::spanned(start, end, Token::CharLiteral(ch))),
      Some((_, end, _)) => self.recover(start, end, Error::UnterminatedCharLiteral, Token::CharLiteral(ch)),
      None => return self.recover(start, start, Error::UnexpectedEof, Token::CharLiteral('\0')),
    }
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

        self.float(start, end, float)
      }
      Some((_, 'e')) | Some((_, 'E')) => {
        self.catchup();

        if self.test_peek(|ch| ch == '-' || ch == '+') {
          self.catchup();
        }

        let (end, float) = self.take_while(start, |ch| is_dec(ch) || ch == '_');
        self.float(start, end, float)
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

      None | Some(_) => self.int(start, end, int),
    }
  }

  fn line_comment(&mut self, start: Location) -> Result<SpannedToken<'input>, SpannedError> {
    let (end, comment) = self.take_until(start, |ch| ch == '\n');

    let token = if comment.starts_with("///") {
      Token::DocComment(Comment {
        typ: CommentType::Line,
        content: comment[3..].trim(),
      })
    } else {
      Token::Comment(Comment {
        typ: CommentType::Line,
        content: comment[2..].trim(),
      })
    };

    Ok(pos::spanned(start, end, token))
  }

  fn block_comment(&mut self, start: Location) -> Result<SpannedToken<'input>, SpannedError> {
    self.bump(); // Skip first '*'

    loop {
      let (_, comment) = self.take_until(start, |ch| ch == '*');
      self.bump(); // Skip next b'*'

      match self.peek() {
        Some((_, '/')) => {
          let (_, end, _) = self.bump().unwrap();

          let token = if comment.starts_with("/**") && comment != "/**" {
            Token::DocComment(Comment {
              typ: CommentType::Block,
              content: comment[3..].trim(),
            })
          } else {
            Token::Comment(Comment {
              typ: CommentType::Block,
              content: comment[2..].trim(),
            })
          };

          return Ok(pos::spanned(start, end, token));
        }
        Some((_, _)) => continue,
        None => return self.error(start, self.loc, Error::UnexpectedEof),
      }
    }
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
      "^" => Token::BitwiseXOr,
      "<:" => Token::BitwiseShiftLeft,
      ":>" => Token::BitwiseShiftRight,
      "*" => Token::Multiplication,
      "/" => Token::Division,
      "%" => Token::Modulo,
      "+" => Token::Addition,
      "-" => Token::Subtraction,
      "==" => Token::Equal,
      "!=" => Token::NotEqual,
      "<" => Token::Less,
      "<=" => Token::LessOrEqual,
      ">" => Token::Greater,
      ">=" => Token::GreaterOrEqual,
      "." => Token::Dot,
      ".." => Token::DotDot,
      "=" => Token::Assign,
      "->" => Token::RArrow,

      _ => todo!(),
    };

    Ok(pos::spanned(start, end, token))
  }

  fn identifier(&mut self, start: Location) -> Result<SpannedToken<'input>, SpannedError> {
    let (end, ident) = self.take_while(start, is_ident_continue);

    let token = match ident {
      "true" => Token::BoolLiteral(true),
      "false" => Token::BoolLiteral(false),

      "resource" => Token::Resource,
      "data" => Token::Data,
      "provider" => Token::Provider,
      "module" => Token::Module,
      "let" => Token::Let,
      "if" => Token::If,
      "else" => Token::Else,
      "for" => Token::Forall,
      "in" => Token::In,
      "while" => Token::Do,
      "break" => Token::Break,
      "continue" => Token::Continue,
      "match" => Token::Match,
      "return" => Token::Return,

      src => Token::Identifier(src),
    };

    Ok(pos::spanned(start, end, token))
  }

  fn int(&mut self, start: Location, end: Location, v: &str) -> Result<SpannedToken<'input>, SpannedError> {
    let v = v.replace('_', "");
    match v.parse::<i64>() {
      Ok(val) => Ok(pos::spanned(start, end, Token::IntLiteral(val))),
      Err(e) => self.recover(start, end, Error::NonParseableInt(e), Token::IntLiteral(0)),
    }
  }

  fn float(&mut self, start: Location, end: Location, v: &str) -> Result<SpannedToken<'input>, SpannedError> {
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

  fn parse_escape_code(&mut self, start: Location, delimiter: char) -> Result<char, SpannedError> {
    match self.bump() {
      Some((_, _, 'n')) => Ok('\n'),
      Some((_, _, 'r')) => Ok('\r'),
      Some((_, _, 't')) => Ok('\t'),
      Some((_, _, '\\')) => Ok('\\'),
      Some((_, _, ch)) if ch == delimiter => Ok(ch),
      Some((_, _, 'u')) => self.parse_codepoint(start),

      Some((_, end, ch)) => Err(pos::spanned(start, end, Error::UnexpectedEscapeCode(ch))),
      None => Err(pos::spanned(start, start, Error::UnexpectedEof)),
    }
  }

  fn parse_codepoint(&mut self, start: Location) -> Result<char, SpannedError> {
    if self.test_peek(|ch| ch != '{') {
      return Err(pos::spanned(start, start, Error::MalformedUnicodeEscape));
    }
    let code_start = self.catchup();

    let (end, codepoint) = Lexer::take_while(self, code_start, is_hex);

    if self.test_peek(|ch| ch != '}') {
      return Err(pos::spanned(start, end, Error::MalformedUnicodeEscape));
    }
    self.catchup();

    if codepoint.is_empty() {
      return Err(pos::spanned(start, end, Error::EmptyUnicodeEscape));
    }

    match u32::from_str_radix(codepoint, 16) {
      Ok(v) => match char::from_u32(v) {
        None => Err(pos::spanned(start, end, Error::InvalidUnicodeEscape)),
        Some(v) => Ok(v),
      },
      Err(_) => Err(pos::spanned(start, end, Error::InvalidUnicodeEscape)),
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

  fn skip_to_end(&mut self) {
    while let Some(_) = self.bump() {}
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

  fn error<T>(&mut self, start: Location, end: Location, code: Error) -> Result<T, SpannedError> {
    self.skip_to_end();
    Err(pos::spanned(start, end, code))
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
    '^' => true,
    '|' => true,
    '~' => true,

    _ => false,
  }
}

fn is_ident_start(ch: char) -> bool {
  unic_ucd_ident::is_xid_start(ch)
}

fn is_ident_continue(ch: char) -> bool {
  unic_ucd_ident::is_xid_continue(ch) || ch == '-'
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

  #[test]
  fn codepoint() {
    let input = "{2764}";
    let mut lexer = Lexer::new(input);

    let res = lexer.parse_codepoint(loc(0));

    assert_eq!(res, Ok('❤'))
  }

  #[test]
  fn ident_start() {
    assert!(is_ident_start('a'));
    assert!(is_ident_start('A'));
    assert!(!is_ident_start('1'));
    assert!(!is_ident_start('-'));
    assert!(!is_ident_start('_'));
  }

  #[test]
  fn ident_continue() {
    assert!(is_ident_continue('a'));
    assert!(is_ident_continue('A'));
    assert!(is_ident_continue('1'));
    assert!(is_ident_continue('-'));
    assert!(is_ident_continue('_'));
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
        Ok(pos::spanned(loc(0), loc(1), Token::StringDelimiter)),
        Context::String,
      ),
      (
        r#"➖""➖"#,
        Ok(pos::spanned(loc(0), loc(1), Token::StringDelimiter)),
        Context::String,
      ),
      (
        r#"➖"➖"#,
        Ok(pos::spanned(loc(0), loc(1), Token::StringDelimiter)),
        Context::String,
      ),
    ];

    test(tests)
  }

  #[test]
  fn heredoc_start() {
    let tests = vec![
      (
        "<< EOF\nfoo\nEOF",
        Ok(pos::spanned(loc(0), loc(6), Token::StringDelimiter)),
        Context::HereDocString {
          delimiter: "EOF",
          skip_leading_tabs: false,
          quoted_delimiter: false,
        },
      ),
      (
        "<<EOF\nfoo\nEOF",
        Ok(pos::spanned(loc(0), loc(5), Token::StringDelimiter)),
        Context::HereDocString {
          delimiter: "EOF",
          skip_leading_tabs: false,
          quoted_delimiter: false,
        },
      ),
      (
        "<<   EOF\nfoo\nEOF",
        Ok(pos::spanned(loc(0), loc(8), Token::StringDelimiter)),
        Context::HereDocString {
          delimiter: "EOF",
          skip_leading_tabs: false,
          quoted_delimiter: false,
        },
      ),
      (
        "<<- EOF\nfoo\nEOF",
        Ok(pos::spanned(loc(0), loc(7), Token::StringDelimiter)),
        Context::HereDocString {
          delimiter: "EOF",
          skip_leading_tabs: true,
          quoted_delimiter: false,
        },
      ),
      (
        "<<-EOF\nfoo\nEOF",
        Ok(pos::spanned(loc(0), loc(6), Token::StringDelimiter)),
        Context::HereDocString {
          delimiter: "EOF",
          skip_leading_tabs: true,
          quoted_delimiter: false,
        },
      ),
      (
        "<<-   EOF\nfoo\nEOF",
        Ok(pos::spanned(loc(0), loc(9), Token::StringDelimiter)),
        Context::HereDocString {
          delimiter: "EOF",
          skip_leading_tabs: true,
          quoted_delimiter: false,
        },
      ),
      (
        "<< - EOF\nfoo\nEOF",
        Ok(pos::spanned(loc(0), loc(8), Token::StringDelimiter)),
        Context::HereDocString {
          delimiter: "- EOF",
          skip_leading_tabs: false,
          quoted_delimiter: false,
        },
      ),
      (
        "<<   - EOF\nfoo\nEOF",
        Ok(pos::spanned(loc(0), loc(10), Token::StringDelimiter)),
        Context::HereDocString {
          delimiter: "- EOF",
          skip_leading_tabs: false,
          quoted_delimiter: false,
        },
      ),
      (
        "<< -   EOF\nfoo\nEOF",
        Ok(pos::spanned(loc(0), loc(10), Token::StringDelimiter)),
        Context::HereDocString {
          delimiter: "-   EOF",
          skip_leading_tabs: false,
          quoted_delimiter: false,
        },
      ),
      (
        "<< \"EOF\"\nfoo\nEOF",
        Ok(pos::spanned(loc(0), loc(8), Token::StringDelimiter)),
        Context::HereDocString {
          delimiter: "EOF",
          skip_leading_tabs: false,
          quoted_delimiter: true,
        },
      ),
      (
        "<<\"EOF\"\nfoo\nEOF",
        Ok(pos::spanned(loc(0), loc(7), Token::StringDelimiter)),
        Context::HereDocString {
          delimiter: "EOF",
          skip_leading_tabs: false,
          quoted_delimiter: true,
        },
      ),
      (
        "<<   \"EOF\"\nfoo\nEOF",
        Ok(pos::spanned(loc(0), loc(10), Token::StringDelimiter)),
        Context::HereDocString {
          delimiter: "EOF",
          skip_leading_tabs: false,
          quoted_delimiter: true,
        },
      ),
      (
        "<<- \"EOF\"\nfoo\nEOF",
        Ok(pos::spanned(loc(0), loc(9), Token::StringDelimiter)),
        Context::HereDocString {
          delimiter: "EOF",
          skip_leading_tabs: true,
          quoted_delimiter: true,
        },
      ),
      (
        "<<-\"EOF\"\nfoo\nEOF",
        Ok(pos::spanned(loc(0), loc(8), Token::StringDelimiter)),
        Context::HereDocString {
          delimiter: "EOF",
          skip_leading_tabs: true,
          quoted_delimiter: true,
        },
      ),
      (
        "<<-   \"EOF\"\nfoo\nEOF",
        Ok(pos::spanned(loc(0), loc(11), Token::StringDelimiter)),
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
  fn char_literal() {
    let tests = vec![
      ("'a'", Ok(pos::spanned(loc(0), loc(3), Token::CharLiteral('a'))), Context::General),
      (
        r#"'\n'"#,
        Ok(pos::spanned(loc(0), loc(4), Token::CharLiteral('\n'))),
        Context::General,
      ),
      (
        r#"'\r'"#,
        Ok(pos::spanned(loc(0), loc(4), Token::CharLiteral('\r'))),
        Context::General,
      ),
      (
        r#"'\t'"#,
        Ok(pos::spanned(loc(0), loc(4), Token::CharLiteral('\t'))),
        Context::General,
      ),
      (
        r#"'\''"#,
        Ok(pos::spanned(loc(0), loc(4), Token::CharLiteral('\''))),
        Context::General,
      ),
      (
        r#"'\u{2764}'"#,
        Ok(pos::spanned(loc(0), loc(10), Token::CharLiteral('❤'))),
        Context::General,
      ),
      (
        r#"'\u{}'"#,
        Ok(pos::spanned(loc(0), loc(6), Token::CharLiteral('\0'))),
        Context::General,
      ),
      (
        r#"'\u{k}'"#,
        Ok(pos::spanned(loc(0), loc(7), Token::CharLiteral('\0'))),
        Context::General,
      ),
      (
        r#"'\"'"#,
        Ok(pos::spanned(loc(0), loc(4), Token::CharLiteral('\0'))),
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
      ("^", Ok(pos::spanned(loc(0), loc(1), Token::BitwiseXOr)), Context::General),
      ("<:", Ok(pos::spanned(loc(0), loc(2), Token::BitwiseShiftLeft)), Context::General),
      (":>", Ok(pos::spanned(loc(0), loc(2), Token::BitwiseShiftRight)), Context::General),
      ("*", Ok(pos::spanned(loc(0), loc(1), Token::Multiplication)), Context::General),
      ("/", Ok(pos::spanned(loc(0), loc(1), Token::Division)), Context::General),
      ("%", Ok(pos::spanned(loc(0), loc(1), Token::Modulo)), Context::General),
      ("+", Ok(pos::spanned(loc(0), loc(1), Token::Addition)), Context::General),
      ("-", Ok(pos::spanned(loc(0), loc(1), Token::Subtraction)), Context::General),
      ("==", Ok(pos::spanned(loc(0), loc(2), Token::Equal)), Context::General),
      ("!=", Ok(pos::spanned(loc(0), loc(2), Token::NotEqual)), Context::General),
      ("<", Ok(pos::spanned(loc(0), loc(1), Token::Less)), Context::General),
      ("<=", Ok(pos::spanned(loc(0), loc(2), Token::LessOrEqual)), Context::General),
      (">", Ok(pos::spanned(loc(0), loc(1), Token::Greater)), Context::General),
      (">=", Ok(pos::spanned(loc(0), loc(2), Token::GreaterOrEqual)), Context::General),
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

  #[test]
  fn line_comment() {
    let tests = vec![
      (
        "// foo",
        Ok(pos::spanned(
          loc(0),
          loc(6),
          Token::Comment(Comment {
            typ: CommentType::Line,
            content: "foo",
          }),
        )),
        Context::General,
      ),
      (
        "//foo",
        Ok(pos::spanned(
          loc(0),
          loc(5),
          Token::Comment(Comment {
            typ: CommentType::Line,
            content: "foo",
          }),
        )),
        Context::General,
      ),
      (
        "//   foo",
        Ok(pos::spanned(
          loc(0),
          loc(8),
          Token::Comment(Comment {
            typ: CommentType::Line,
            content: "foo",
          }),
        )),
        Context::General,
      ),
      (
        "/// foo",
        Ok(pos::spanned(
          loc(0),
          loc(7),
          Token::DocComment(Comment {
            typ: CommentType::Line,
            content: "foo",
          }),
        )),
        Context::General,
      ),
      (
        "///foo",
        Ok(pos::spanned(
          loc(0),
          loc(6),
          Token::DocComment(Comment {
            typ: CommentType::Line,
            content: "foo",
          }),
        )),
        Context::General,
      ),
      (
        "///   foo",
        Ok(pos::spanned(
          loc(0),
          loc(9),
          Token::DocComment(Comment {
            typ: CommentType::Line,
            content: "foo",
          }),
        )),
        Context::General,
      ),
    ];

    test(tests);
  }

  #[test]
  fn block_comment() {
    let tests = vec![
      (
        "/* foo */",
        Ok(pos::spanned(
          loc(0),
          loc(9),
          Token::Comment(Comment {
            typ: CommentType::Block,
            content: "foo",
          }),
        )),
        Context::General,
      ),
      (
        "/*foo*/",
        Ok(pos::spanned(
          loc(0),
          loc(7),
          Token::Comment(Comment {
            typ: CommentType::Block,
            content: "foo",
          }),
        )),
        Context::General,
      ),
      (
        "/*   foo   */",
        Ok(pos::spanned(
          loc(0),
          loc(13),
          Token::Comment(Comment {
            typ: CommentType::Block,
            content: "foo",
          }),
        )),
        Context::General,
      ),
      (
        "/** foo */",
        Ok(pos::spanned(
          loc(0),
          loc(10),
          Token::DocComment(Comment {
            typ: CommentType::Block,
            content: "foo",
          }),
        )),
        Context::General,
      ),
      (
        "/**foo*/",
        Ok(pos::spanned(
          loc(0),
          loc(8),
          Token::DocComment(Comment {
            typ: CommentType::Block,
            content: "foo",
          }),
        )),
        Context::General,
      ),
      (
        "/**   foo   */",
        Ok(pos::spanned(
          loc(0),
          loc(14),
          Token::DocComment(Comment {
            typ: CommentType::Block,
            content: "foo",
          }),
        )),
        Context::General,
      ),
      (
        "/* foo *",
        Err(pos::spanned(loc(0), loc(8), Error::UnexpectedEof)),
        Context::General,
      ),
    ];

    test(tests);
  }

  #[test]
  fn identifier() {
    let tests = vec![
      ("true", Ok(pos::spanned(loc(0), loc(4), Token::BoolLiteral(true))), Context::General),
      (
        "false",
        Ok(pos::spanned(loc(0), loc(5), Token::BoolLiteral(false))),
        Context::General,
      ),
      ("resource", Ok(pos::spanned(loc(0), loc(8), Token::Resource)), Context::General),
      ("data", Ok(pos::spanned(loc(0), loc(4), Token::Data)), Context::General),
      ("provider", Ok(pos::spanned(loc(0), loc(8), Token::Provider)), Context::General),
      ("module", Ok(pos::spanned(loc(0), loc(6), Token::Module)), Context::General),
      ("let", Ok(pos::spanned(loc(0), loc(3), Token::Let)), Context::General),
      ("if", Ok(pos::spanned(loc(0), loc(2), Token::If)), Context::General),
      ("else", Ok(pos::spanned(loc(0), loc(4), Token::Else)), Context::General),
      ("for", Ok(pos::spanned(loc(0), loc(3), Token::Forall)), Context::General),
      ("in", Ok(pos::spanned(loc(0), loc(2), Token::In)), Context::General),
      ("while", Ok(pos::spanned(loc(0), loc(5), Token::Do)), Context::General),
      ("break", Ok(pos::spanned(loc(0), loc(5), Token::Break)), Context::General),
      ("continue", Ok(pos::spanned(loc(0), loc(8), Token::Continue)), Context::General),
      ("match", Ok(pos::spanned(loc(0), loc(5), Token::Match)), Context::General),
      ("return", Ok(pos::spanned(loc(0), loc(6), Token::Return)), Context::General),
      ("foo", Ok(pos::spanned(loc(0), loc(3), Token::Identifier("foo"))), Context::General),
    ];

    test(tests);
  }

  #[test]
  fn string_simple() {
    let tests = vec![
      (
        r#"➖"foo"➖"#,
        Ok(pos::spanned(loc(1), loc(4), Token::String(StringLiteral::Escaped("foo")))),
      ),
      (
        r#"➖"foo\n"➖"#,
        Ok(pos::spanned(loc(1), loc(6), Token::String(StringLiteral::Escaped("foo\\n")))),
      ),
    ];

    for (input, expected) in tests {
      let s = input.replace("➖", "");
      let input = &*s;
      let mut lexer = Lexer::new(input);

      let res = lexer.next().unwrap();
      assert_eq!(Ok(pos::spanned(loc(0), loc(1), Token::StringDelimiter)), res);
      assert_eq!(&Context::String, lexer.context());

      let res = lexer.next().unwrap();
      assert_eq!(expected, res);
      assert_eq!(&Context::String, lexer.context());

      let res = lexer.next().unwrap();
      if let Ok(Spanned { span: _, value }) = res {
        assert_eq!(Token::StringDelimiter, value)
      } else {
        panic!("{:?}", res);
      }
    }
  }

  #[test]
  fn string_interpolation() {
    let tests = vec![
      (r#"➖"123${}"➖"#, Token::StringInterpolation),
      (r#"➖"123%{}"➖"#, Token::StringDirective),
    ];

    for (input, expected) in tests {
      let s = input.replace("➖", "");
      let input = &*s;
      let mut lexer = Lexer::new(input);

      let res = lexer.next().unwrap();
      assert_eq!(Ok(pos::spanned(loc(0), loc(1), Token::StringDelimiter)), res);
      assert_eq!(&Context::String, lexer.context());

      let res = lexer.next().unwrap();
      assert_eq!(Ok(pos::spanned(loc(1), loc(4), Token::String(StringLiteral::Escaped("123")))), res);
      assert_eq!(&Context::String, lexer.context());

      let res = lexer.next().unwrap();
      if let Ok(Spanned { span: _, value }) = res {
        assert_eq!(expected, value)
      } else {
        panic!("{:?}", res);
      }
    }
  }

  #[test]
  fn heredoc() {
    let tests = vec![(
      // "<<EOF\n  foo\n    bar\nEOF",
      // "<<EOF\n\n  foo\n    bar\nEOF\n",
      "<<-EOF\n\n  foo\n    bar\nEOF",
      Ok(pos::spanned(
        Location::new(1, 0, 7),
        Location::new(3, 7, 21),
        Token::HeredocString(vec!["", "foo", "  bar"]),
      )),
      Context::HereDocString {
        delimiter: "EOF",
        skip_leading_tabs: true,
        quoted_delimiter: false,
      },
      Context::HereDocStringEnd {
        start: Location::new(4, 0, 22),
        end: Location::new(4, 3, 25),
      },
    )];

    for (input, expected, exp_ctx_inside, exp_ctx_after) in tests {
      let s = input.replace("➖", "");
      let input = &*s;
      let mut lexer = Lexer::new(input);

      println!("input:\n➖\n{}\n➖", input);

      let res = lexer.next().unwrap();
      assert_eq!(Ok(pos::spanned(loc(0), loc(6), Token::StringDelimiter)), res);
      assert_eq!(&exp_ctx_inside, lexer.context());

      let res = lexer.next().unwrap();
      assert_eq!(expected, res);
      assert_eq!(&exp_ctx_after, lexer.context());

      let res = lexer.next().unwrap();
      if let Ok(Spanned { span: _, value }) = res {
        assert_eq!(Token::StringDelimiter, value)
      } else {
        panic!("{:?}", res);
      }
    }
  }
}
