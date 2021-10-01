#[derive(Clone, Copy, Eq, PartialEq, Hash, Debug)]
pub enum Token<S> {
  IntLiteral(i64),

  SingleLineStringDelimiter,
  MultiLineStringDelimiter,
  HereDocStringDelimiter,

  StringLiteral(StringLiteral<S>),

  EOF,
}

#[derive(Clone, Copy, PartialEq, Eq, Debug, Hash)]
pub enum StringLiteral<S> {
  Escaped(S),
  Raw(S),
}
