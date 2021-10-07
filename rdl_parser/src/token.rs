use ordered_float::NotNan;

#[derive(Clone, Copy, Eq, PartialEq, Hash, Debug)]
pub enum Token<S> {
  IntLiteral(i64),
  FloatLiteral(NotNan<f64>),
  BoolLiteral(bool),
  CharLiteral(char),

  SingleLineStringDelimiter,
  MultiLineStringDelimiter,
  HereDocStringDelimiter,

  StringLiteral(StringLiteral<S>),

  Comment(Comment<S>),
  DocComment(Comment<S>),

  // Operators
  LogicalNot,        // !
  LogicalAnd,        // &&
  LogicalOr,         // ||
  BitwiseNot,        // ~
  BitwiseAnd,        // &
  BitwiseOr,         // |
  BitwiseXOr,        // ^
  BitwiseShiftLeft,  // <:
  BitwiseShiftRight, // :>
  Multiplication,    // *
  Division,          // /
  Modulo,            // %
  Addition,          // +
  Subtraction,       // -
  Equal,             // ==
  NotEqual,          // !=
  Less,              // <
  LessOrEqual,       // <=
  Greater,           // >
  GreaterOrEqual,    // >=
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

  Let,
  If,
  Else,
  Forall,
  In,
  Do,
  Break,
  Continue,
  Match,
  Return,
  Identifier(S),

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
