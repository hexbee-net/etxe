use ordered_float::NotNan;

#[derive(Clone, Copy, Eq, PartialEq, Hash, Debug)]
pub enum Token<S> {
  IntLiteral(i64),
  FloatLiteral(NotNan<f64>),

  SingleLineStringDelimiter,
  MultiLineStringDelimiter,
  HereDocStringDelimiter,

  StringLiteral(StringLiteral<S>),

  Comment(Comment<S>),
  DocComment(Comment<S>),

  LogicalNot,        // !
  LogicalAnd,        // &&
  LogicalOr,         // ||
  BitwiseNot,        // ~
  BitwiseAnd,        // &
  BitwiseOr,         // |
  BitwiseShiftLeft,  // <:
  BitwiseShiftRight, // :>
  Multiplication,    // *
  Division,          // /
  Modulo,            // %
  Addition,          // +
  Subtraction,       // -
  Equal,             // ==
  NotEqual,          // !=
  Dot,               // .
  DotDot,            // ..
  Assign,            // =
  RArrow,            // ->

  Comma,
  LBrace,
  LBracket,
  LParen,
  RBrace,
  RBracket,
  RParen,

  EOF,
}

#[derive(Clone, Copy, PartialEq, Eq, Debug, Hash)]
pub enum StringLiteral<S> {
  Escaped(S),
  Raw(S),
}

#[derive(Clone, Copy, Eq, PartialEq, Debug, Hash)]
pub struct Comment<S = String> {
  pub typ: CommentType,
  pub content: S,
}

#[derive(Clone, Copy, Eq, PartialEq, Debug, Hash)]
pub enum CommentType {
  Block,
  Line,
}
