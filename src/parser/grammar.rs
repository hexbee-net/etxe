// etxe. IaC done right
// Copyright (c) 2021 Xavier Basty Kjellberg
//
// Licensed under the Apache License, Version 2.0
// <LICENSE-APACHE or http://www.apache.org/licenses/LICENSE-2.0> or the MIT
// license <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. All files in the project carrying such notice may not be copied,
// modified, or distributed except according to those terms.

// use pest::Parser;

#[derive(Parser)]
#[grammar = "parser/grammar.pest"]
pub struct RDLParser;
