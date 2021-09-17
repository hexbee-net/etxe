use crate::pos::{ByteOffset, BytePos, Span};

pub trait ParserSource {
  fn src(&self) -> &str;
  fn start_index(&self) -> BytePos;

  fn span(&self) -> Span<BytePos> {
    let start = self.start_index();
    Span::new(start, start + ByteOffset::from(self.src().len() as i64))
  }
}
