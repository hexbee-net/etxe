use std::fmt;

/// A region of code in a source file
#[derive(Clone, Copy, Default, PartialEq, Eq, Ord, PartialOrd)]
pub struct Span<I> {
  start: I,
  end: I,
}

impl<I: fmt::Debug> fmt::Debug for Span<I> {
  fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
    write!(f, "{:?}..{:?}", self.start, self.end)
  }
}

impl<I: fmt::Display> fmt::Display for Span<I> {
  fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
    self.start.fmt(f)?;
    write!(f, "..")?;
    self.end.fmt(f)?;
    Ok(())
  }
}

impl<I: Ord> Span<I> {
  /// Create a new span
  ///
  /// ```rust
  /// use etxe_rdl_parser::pos::Span;
  ///
  /// let span = Span::new(3, 6);
  /// assert_eq!(span.start(), 3);
  /// assert_eq!(span.end(), 6);
  /// ```
  ///
  /// `start` and `end` are reordered to maintain the invariant that `start <= end`
  ///
  /// ```rust
  /// use etxe_rdl_parser::pos::Span;
  ///
  /// let span = Span::new(6, 3);
  /// assert_eq!(span.start(), 3);
  /// assert_eq!(span.end(), 6);
  /// ```
  pub fn new(start: I, end: I) -> Self {
    if start <= end {
      Span { start, end }
    } else {
      Span { start: end, end: start }
    }
  }

  pub fn map<F, J>(self, mut f: F) -> Span<J>
  where
    F: FnMut(I) -> J,
    J: Ord,
  {
    Span::new(f(self.start), f(self.end))
  }
}

impl<I> Span<I> {
  /// Create a span like `new` but does not check that `start <= end`
  pub const fn new_unchecked(start: I, end: I) -> Self {
    Span { start, end }
  }

  /// Get the start index
  pub fn start(self) -> I {
    self.start
  }

  /// Get the end index
  pub fn end(self) -> I {
    self.end
  }
}
