#[derive(Clone, Eq, PartialEq, Hash, Debug)]
pub enum Token<S> {
  IntLiteral(i64),

  StringDelimiter(StringDelimiter<S>),
  StringLiteral(StringLiteral<S>),

  EOF,
}

#[derive(Clone, PartialEq, Eq, Debug, Hash)]
pub enum StringLiteral<S> {
  Escaped(S),
  Raw(S),
}

#[derive(Clone, Copy, PartialEq, Eq, Debug, Hash)]
pub enum StringDelimiter<S> {
  SingleLine,
  MultiLine,
  HereDoc(S),
  Raw(S),
}
