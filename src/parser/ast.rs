
#[derive(Clone, PartialEq, Eq, Hash, Debug)]
pub enum Literal {
    Byte(u8),
    Int(i64),
    Float(NotNan<f64>),
    String(String),
    Char(char),
}
