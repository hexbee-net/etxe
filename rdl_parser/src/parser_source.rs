use crate::pos::{ByteIndex, ByteOffset, Span};

pub trait ParserSource {
  fn src(&self) -> &str;
  fn start_index(&self) -> ByteIndex;

  fn span(&self) -> Span<ByteIndex> {
    let start = self.start_index();
    Span::new(start, start + ByteOffset::from(self.src().len() as i64))
  }
}

impl<'a, S> ParserSource for &'a S
where
  S: ?Sized + ParserSource,
{
  fn src(&self) -> &str {
    (**self).src()
  }
  fn start_index(&self) -> ByteIndex {
    (**self).start_index()
  }
}

impl ParserSource for str {
  fn src(&self) -> &str {
    self
  }
  fn start_index(&self) -> ByteIndex {
    ByteIndex::from(0)
  }
}
