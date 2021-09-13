mod grammar;

use std::fs;
// use crate::parser::grammar::RDLParser;

pub fn parse(src: &str) {
  let unparsed_file = fs::read_to_string(src).expect("cannot read file");
  println!("{}", unparsed_file)

  // let file = RDLParser

  // let file = INIParser::parse(Rule::file, &unparsed_file)
  //     .expect("unsuccessful parse") // unwrap the parse result
  //     .next().unwrap();
}
