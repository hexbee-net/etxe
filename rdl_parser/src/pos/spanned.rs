use std::fmt;

use crate::pos::Span;

#[derive(Copy, Clone, Debug, Eq, Default)]
pub struct Spanned<T, Pos> {
  pub span: Span<Pos>,
  pub value: T,
}

impl<T: fmt::Display, Pos: fmt::Display + Copy> fmt::Display for Spanned<T, Pos> {
  fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
    write!(f, "{}: {}", self.span.start(), self.value)
  }
}

impl<T, Pos> From<(T, Span<Pos>)> for Spanned<T, Pos> {
  fn from((value, span): (T, Span<Pos>)) -> Self {
    Spanned { span, value }
  }
}

impl<T, Pos: Default> From<T> for Spanned<T, Pos> {
  fn from(value: T) -> Self {
    Spanned {
      span: Span::default(),
      value,
    }
  }
}

impl<T: PartialEq, Pos> PartialEq<T> for Spanned<T, Pos> {
  fn eq(&self, other: &T) -> bool {
    self.value == *other
  }
}

impl<T, Pos> std::ops::Deref for Spanned<T, Pos> {
  type Target = T;
  fn deref(&self) -> &T {
    &self.value
  }
}

impl<T, Pos> std::ops::DerefMut for Spanned<T, Pos> {
  fn deref_mut(&mut self) -> &mut T {
    &mut self.value
  }
}

impl<T: AsRef<U>, U: ?Sized, Pos> AsRef<U> for Spanned<T, Pos> {
  fn as_ref(&self) -> &U {
    self.value.as_ref()
  }
}

impl<T, Pos> std::hash::Hash for Spanned<T, Pos>
where
  T: std::hash::Hash,
  Pos: std::hash::Hash + Copy,
{
  fn hash<H>(&self, state: &mut H)
  where
    H: std::hash::Hasher,
  {
    self.span.start().hash(state);
    self.span.end().hash(state);
    self.value.hash(state);
  }
}

impl<T, Pos> PartialEq for Spanned<T, Pos>
where
  T: std::cmp::PartialEq,
  Pos: std::cmp::PartialEq,
{
  fn eq(&self, other: &Self) -> bool {
    self.span == other.span && self.value == other.value
  }
}

impl<T, Pos> Spanned<T, Pos> {
  pub fn map<U, F>(self, mut f: F) -> Spanned<U, Pos>
  where
    F: FnMut(T) -> U,
  {
    Spanned {
      span: self.span,
      value: f(self.value),
    }
  }
}

pub fn spanned<T, Pos>(start: Pos, end: Pos, value: T) -> Spanned<T, Pos>
where
  Pos: Ord,
{
  Spanned {
    // span: span(start, end),
    span: Span::new(start, end),
    value,
  }
}
