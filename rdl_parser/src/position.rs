// //! Source code locations (borrowed from rustc's [libsyntax_pos])
// //!
// //! [libsyntax_pos]: https://github.com/rust-lang/rust/blob/master/src/libsyntax_pos/lib.rs
//
// use crate::source::CodeMap;
// pub use crate::source::Error;
// use std::{cmp, cmp::Ordering, fmt};
//
// pub use codespan::{
//   ByteIndex, ByteIndex as BytePos, ByteOffset, ColumnIndex as Column, ColumnOffset, Index,
//   LineIndex as Line, LineOffset, RawIndex,
// };
//
// ////////////////////////////////////////////////////////////////////////////////
// // Location
//
// /// A location in a source file
// #[derive(Copy, Clone, Default, Eq, PartialEq, Hash, Ord, PartialOrd)]
// pub struct Location {
//   pub line: Line,
//   pub column: Column,
//   pub absolute: BytePos,
// }
//
// impl fmt::Debug for Location {
//   fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
//     write!(
//       f,
//       "{}:{}/{}",
//       self.line.number(),
//       self.column.number(),
//       self.absolute.to_usize()
//     )
//   }
// }
//
// impl fmt::Display for Location {
//   fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
//     write!(
//       f,
//       "Line: {}, Column: {}",
//       self.line.number(),
//       self.column.number()
//     )
//   }
// }
//
// impl Location {
//   /// Update `self` by consuming the specified character.
//   ///
//   /// ```rust
//   /// use etxe_rdl_parser::position::{Location, Line, Column, BytePos};
//   ///
//   /// let mut loc = Location { line: Line(1), column: Column(1), absolute: BytePos(0) };
//   /// loc.shift(b' ');
//   /// assert_eq!(loc, Location { line: Line(1), column: Column(2), absolute: BytePos(1) });
//   /// loc.shift(b'\n');
//   /// assert_eq!(loc, Location { line: Line(2), column: Column(1), absolute: BytePos(2) });
//   /// ```
//   pub fn shift(&mut self, ch: u8) {
//     if ch == b'\n' {
//       self.line += LineOffset(1);
//       self.column = Column(1);
//     } else {
//       self.column += ColumnOffset(1);
//     }
//     self.absolute += ByteOffset(1);
//   }
// }
//
// ////////////////////////////////////////////////////////////////////////////////
// // Span
//
//// impl<I: Index> Span<I> {
//   /// Create a new span from a byte start and an offset
//   pub fn from_offset(start: I, off: I::Offset) -> Span<I> {
//     Span::new(start, start + off)
//   }
//
//   /// Return a new span with the start position replaced with the supplied byte position
//   ///
//   /// ```rust
//   /// use etxe_rdl_parser::position::{ByteIndex, Span};
//   ///
//   /// let span = Span::new(ByteIndex(3), ByteIndex(6));
//   /// assert_eq!(span.with_start(ByteIndex(2)), Span::new(ByteIndex(2), ByteIndex(6)));
//   /// assert_eq!(span.with_start(ByteIndex(5)), Span::new(ByteIndex(5), ByteIndex(6)));
//   /// assert_eq!(span.with_start(ByteIndex(7)), Span::new(ByteIndex(6), ByteIndex(7)));
//   /// ```
//   pub fn with_start(self, start: I) -> Span<I> {
//     Span::new(start, self.end())
//   }
//
//   /// Return a new span with the end position replaced with the supplied byte position
//   ///
//   /// ```rust
//   /// use etxe_rdl_parser::position::{ByteIndex, Span};
//   ///
//   /// let span = Span::new(ByteIndex(3), ByteIndex(6));
//   /// assert_eq!(span.with_end(ByteIndex(7)), Span::new(ByteIndex(3), ByteIndex(7)));
//   /// assert_eq!(span.with_end(ByteIndex(5)), Span::new(ByteIndex(3), ByteIndex(5)));
//   /// assert_eq!(span.with_end(ByteIndex(2)), Span::new(ByteIndex(2), ByteIndex(3)));
//   /// ```
//   pub fn with_end(self, end: I) -> Span<I> {
//     Span::new(self.start(), end)
//   }
//
//   /// Return a `Span` that would enclose both `self` and `end`.
//   ///
//   /// ```plain
//   /// self     ~~~~~~~
//   /// end                     ~~~~~~~~
//   /// returns  ~~~~~~~~~~~~~~~~~~~~~~~
//   /// ```
//   ///
//   /// ```rust
//   /// use etxe_rdl_parser::position::{ByteIndex, Span};
//   ///
//   /// let a = Span::new(ByteIndex(2), ByteIndex(5));
//   /// let b = Span::new(ByteIndex(10), ByteIndex(14));
//   ///
//   /// assert_eq!(a.to(b), Span::new(ByteIndex(2), ByteIndex(14)));
//   /// ```
//   pub fn to(self, end: Span<I>) -> Span<I> {
//     Span::new(
//       cmp::min(self.start(), end.start()),
//       cmp::max(self.end(), end.end()),
//     )
//   }
//
//   /// Return a `Span` between the end of `self` to the beginning of `end`.
//   ///
//   /// ```plain
//   /// self     ~~~~~~~
//   /// end                     ~~~~~~~~
//   /// returns         ~~~~~~~~~
//   /// ```
//   ///
//   /// ```rust
//   /// use etxe_rdl_parser::position::{ByteIndex, Span};
//   ///
//   /// let a = Span::new(ByteIndex(2), ByteIndex(5));
//   /// let b = Span::new(ByteIndex(10), ByteIndex(14));
//   ///
//   /// assert_eq!(a.between(b), Span::new(ByteIndex(5), ByteIndex(10)));
//   /// ```
//   pub fn between(self, end: Span<I>) -> Span<I> {
//     Span::new(self.end(), end.start())
//   }
//
//   /// Return a `Span` between the beginning of `self` to the beginning of `end`.
//   ///
//   /// ```plain
//   /// self     ~~~~~~~
//   /// end                     ~~~~~~~~
//   /// returns  ~~~~~~~~~~~~~~~~
//   /// ```
//   ///
//   /// ```rust
//   /// use etxe_rdl_parser::position::{ByteIndex, Span};
//   ///
//   /// let a = Span::new(ByteIndex(2), ByteIndex(5));
//   /// let b = Span::new(ByteIndex(10), ByteIndex(14));
//   ///
//   /// assert_eq!(a.until(b), Span::new(ByteIndex(2), ByteIndex(10)));
//   /// ```
//   pub fn until(self, end: Span<I>) -> Span<I> {
//     Span::new(self.start(), end.start())
//   }
//
//   /// Makes a `Span` from offsets relative to the start of this span.
//   pub fn sub_span(&self, begin: I::Offset, end: I::Offset) -> Span<I> {
//     assert!(end >= begin);
//     assert!(self.start() + end <= self.end());
//     Span {
//       start: self.start() + begin,
//       end: self.start() + end,
//     }
//   }
//
//   /// Return true if `self` fully encloses `other`.
//   ///
//   /// ```rust
//   /// use etxe_rdl_parser::position::{ByteIndex, Span};
//   ///
//   /// let a = Span::new(ByteIndex(5), ByteIndex(8));
//   ///
//   /// assert_eq!(a.contains(a), true);
//   /// assert_eq!(a.contains(Span::new(ByteIndex(6), ByteIndex(7))), true);
//   /// assert_eq!(a.contains(Span::new(ByteIndex(6), ByteIndex(10))), false);
//   /// assert_eq!(a.contains(Span::new(ByteIndex(3), ByteIndex(6))), false);
//   /// ```
//   pub fn contains(self, other: Span<I>) -> bool {
//     self.start() <= other.start() && other.end() <= self.end()
//   }
//
//   /// Return true if the position is within `self`.
//   ///
//   /// ```rust
//   /// use etxe_rdl_parser::position::{ByteIndex, Span};
//   ///
//   /// let a = Span::new(ByteIndex(5), ByteIndex(7));
//   /// assert_eq!(a.contains_pos(ByteIndex(5)), true);
//   /// assert_eq!(a.contains_pos(ByteIndex(6)), true);
//   /// assert_eq!(a.contains_pos(ByteIndex(7)), true);
//   /// assert_eq!(a.contains_pos(ByteIndex(4)), false);
//   /// assert_eq!(a.contains_pos(ByteIndex(8)), false);
//   /// ```
//   pub fn contains_pos(self, other: I) -> bool {
//     self.start() <= other && other <= self.end()
//   }
//
//   /// Return `Equal` if `self` contains `pos`, otherwise it returns `Less` if `pos` is before
//   /// `start` or `Greater` if `pos` is after or at `end`.
//   ///
//   /// ```rust
//   /// use etxe_rdl_parser::position::{ByteIndex, Span};
//   /// use std::cmp::Ordering::*;
//   ///
//   /// let a = Span::new(ByteIndex(5), ByteIndex(8));
//   ///
//   /// assert_eq!(a.containment(ByteIndex(4)), Less);
//   /// assert_eq!(a.containment(ByteIndex(5)), Equal);
//   /// assert_eq!(a.containment(ByteIndex(6)), Equal);
//   /// assert_eq!(a.containment(ByteIndex(8)), Equal);
//   /// assert_eq!(a.containment(ByteIndex(9)), Greater);
//   /// ```
//   pub fn containment(self, pos: I) -> Ordering {
//     use std::cmp::Ordering::*;
//
//     match (pos.cmp(&self.start), pos.cmp(&self.end)) {
//       (Equal, _) | (_, Equal) | (Greater, Less) => Equal,
//       (Less, _) => Less,
//       (_, Greater) => Greater,
//     }
//   }
//
//   /// Return `Equal` if `self` contains `pos`, otherwise it returns `Less` if `pos` is before
//   /// `start` or `Greater` if `pos` is *strictly* after `end`.
//   ///
//   /// ```rust
//   /// use etxe_rdl_parser::position::{ByteIndex, Span};
//   /// use std::cmp::Ordering::*;
//   ///
//   /// let a = Span::new(ByteIndex(5), ByteIndex(8));
//   ///
//   /// assert_eq!(a.containment_exclusive(ByteIndex(4)), Less);
//   /// assert_eq!(a.containment_exclusive(ByteIndex(5)), Equal);
//   /// assert_eq!(a.containment_exclusive(ByteIndex(6)), Equal);
//   /// assert_eq!(a.containment_exclusive(ByteIndex(8)), Greater);
//   /// assert_eq!(a.containment_exclusive(ByteIndex(9)), Greater);
//   /// ```
//   pub fn containment_exclusive(self, pos: I) -> Ordering {
//     if self.end == pos {
//       Ordering::Greater
//     } else {
//       self.containment(pos)
//     }
//   }
// }
//
// impl Span<BytePos> {
//   pub fn to_range(self, source: &CodeMap) -> Option<std::ops::Range<usize>> {
//     let start = source.pos_to_usize(self.start())?;
//     let end = source.pos_to_usize(self.end())?;
//     Some(start..end)
//   }
// }
//

// ////////////////////////////////////////////////////////////////////////////////
// // Public Functions
//
// pub fn span<Pos>(start: Pos, end: Pos) -> Span<Pos>
// where
//   Pos: Ord,
// {
//   Span::new(start, end)
// }
//
// // pub fn spanned<T, Pos>(span: Span<Pos>, value: T) -> Spanned<T, Pos> {
// //   Spanned { span, value }
// // }
//
// pub fn spanned<T, Pos>(start: Pos, end: Pos, value: T) -> Spanned<T, Pos>
// where
//   Pos: Ord,
// {
//   Spanned {
//     span: span(start, end),
//     value,
//   }
// }
//
// pub trait HasSpan {
//   fn span(&self) -> Span<BytePos>;
// }
