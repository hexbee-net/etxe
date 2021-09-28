use std::ops;
use std::cmp;
use std::{fmt};

pub use codespan::{ByteIndex, ByteOffset, ColumnIndex, ColumnOffset, Index, LineIndex, LineOffset, RawIndex, RawOffset};
use std::cmp::Ordering;
use std::fmt::Formatter;

/// A location in a source file
#[derive(Copy, Clone, Default, Eq, PartialEq, Hash, Ord)]
pub struct Location {
  pub line: LineIndex,
  pub column: ColumnIndex,
  pub absolute: ByteIndex,
}

impl fmt::Display for Location {
  fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
    write!(f, "Line: {}, Column: {}", self.line.number(), self.column.number())
  }
}

impl fmt::Debug for Location {
  fn fmt(&self, f: &mut Formatter<'_>) -> fmt::Result {
    write!(f, "L{}C{}({})", self.line, self.column, self.absolute)
  }
}

impl Location {
  pub fn new(line: RawIndex, column: RawIndex, absolute: RawIndex)  -> Location {
    Location{
      line: LineIndex::from(line),
      column: ColumnIndex::from(column),
      absolute: ByteIndex::from(absolute),
    }
  }

  pub fn shift(&mut self, ch: char) {
    match ch {
      '\n' => {
        self.line += LineOffset(1);
        self.column = ColumnIndex(0);
      }
      _ => {
        self.column += ColumnOffset(1);
      }
    }

    let l = ch.len_utf8() as RawOffset;
    self.absolute += ByteOffset(l);
  }
}

impl ops::Add<char> for Location {
  type Output = Location;

  fn add(self, rhs: char) -> Self::Output {
    let mut loc = self;
    loc.shift(rhs);
    loc
  }
}

impl cmp::PartialOrd for Location {
  fn partial_cmp(&self, other: &Self) -> Option<Ordering> {
    Some(self.absolute.cmp(&other.absolute))
  }
}
