use ordered_float::NotNan;

#[derive(Clone, Eq, PartialEq, Hash, Debug)]
pub enum Token<S> {
  IntLiteral(i64),
  FloatLiteral(NotNan<f64>),
  BoolLiteral(bool),
  CharLiteral(char),

  StringDelimiter,
  String(StringLiteral<S>),
  HeredocString(Vec<S>),

  StringInterpolation,
  StringDirective,

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

  // Keywords
  Resource,
  Data,
  Provider,
  Module,
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

#[derive(Clone, PartialEq, Eq, Debug, Hash)]
pub enum StringLiteral<S> {
  Escaped(S),
  Raw(S),
}

#[derive(Clone, Eq, PartialEq, Debug, Hash)]
pub struct Comment<S = String> {
  pub typ: CommentType,
  pub content: S,
}

#[derive(Clone, Eq, PartialEq, Debug, Hash)]
pub enum CommentType {
  Block,
  Line,
}
