use etxe_core::error::Errors;
use etxe_rdl_parser::lexer::{Context, Error, Lexer};
use etxe_rdl_parser::pos::{self, Location};
use etxe_rdl_parser::token::{StringLiteral, Token};

fn loc(pos: u32) -> Location {
  Location::new(0, pos, pos)
}

struct TestCase<'a> {
  description: &'a str,
  input: &'a str,
  steps: Vec<TestStep<'a>>,
}

struct TestStep<'a> {
  description: &'a str,
  token: Result<(Location, Location, Token<&'a str>), (Location, Location, Error)>,
  context: Context<'a>,
  errors: Option<Errors<&'a str>>,
}

fn test(tests: Vec<TestCase>) {
  for case in tests {
    let input = &*case.input.replace("➖", "");
    let mut lexer = Lexer::new(input);

    for step in case.steps {
      let exp_token = match step.token {
        Ok((start, end, token)) => Ok(pos::spanned(start, end, token)),
        Err((start, end, err)) => Err(pos::spanned(start, end, err)),
      };

      let res = lexer.next().unwrap();

      assert_eq!(exp_token, res, "[{}]/[{}]: token", case.description, step.description);
      assert_eq!(
        &step.context,
        lexer.context(),
        "[{}]/[{}]: context",
        case.description,
        step.description
      );
    }
  }
}

#[test]
fn string_start() {
  let tests = vec![
    TestCase {
      description: "simple string",
      input: r#"➖"123"➖"#,
      steps: vec![TestStep {
        description: "opening delimiter",
        token: Ok((loc(0), loc(1), Token::StringDelimiter)),
        context: Context::String,
        errors: None,
      }],
    },
    TestCase {
      description: "empty string",
      input: r#"➖""➖"#,
      steps: vec![TestStep {
        description: "opening delimiter",
        token: Ok((loc(0), loc(1), Token::StringDelimiter)),
        context: Context::String,
        errors: None,
      }],
    },
    TestCase {
      description: "single token string",
      input: r#"➖"➖"#,
      steps: vec![TestStep {
        description: "opening delimiter",
        token: Ok((loc(0), loc(1), Token::StringDelimiter)),
        context: Context::String,
        errors: None,
      }],
    },
  ];

  test(tests)
}

#[test]
fn heredoc_start() {
  let tests = vec![
    TestCase {
      description: "no-quote no-tab spaced",
      input: "<< EOF\n123\nEOF",
      steps: vec![TestStep {
        description: "opening delimiter",
        token: Ok((loc(0), loc(6), Token::StringDelimiter)),
        context: Context::HereDocString {
          delimiter: "EOF",
          skip_leading_tabs: false,
          quoted_delimiter: false,
        },
        errors: None,
      }],
    },
    TestCase {
      description: "no-quote no-tab not-spaced",
      input: "<<EOF\n123\nEOF",
      steps: vec![TestStep {
        description: "opening delimiter",
        token: Ok((loc(0), loc(5), Token::StringDelimiter)),
        context: Context::HereDocString {
          delimiter: "EOF",
          skip_leading_tabs: false,
          quoted_delimiter: false,
        },
        errors: None,
      }],
    },
    TestCase {
      description: "no-quote no-tab multi-spaced",
      input: "<<   EOF\n123\nEOF",
      steps: vec![TestStep {
        description: "opening delimiter",
        token: Ok((loc(0), loc(8), Token::StringDelimiter)),
        context: Context::HereDocString {
          delimiter: "EOF",
          skip_leading_tabs: false,
          quoted_delimiter: false,
        },
        errors: None,
      }],
    },
    TestCase {
      description: "no-quote tab spaced",
      input: "<<- EOF\n123\nEOF",
      steps: vec![TestStep {
        description: "opening delimiter",
        token: Ok((loc(0), loc(7), Token::StringDelimiter)),
        context: Context::HereDocString {
          delimiter: "EOF",
          skip_leading_tabs: true,
          quoted_delimiter: false,
        },
        errors: None,
      }],
    },
    TestCase {
      description: "no-quote tab not-spaced",
      input: "<<-EOF\n123\nEOF",
      steps: vec![TestStep {
        description: "opening delimiter",
        token: Ok((loc(0), loc(6), Token::StringDelimiter)),
        context: Context::HereDocString {
          delimiter: "EOF",
          skip_leading_tabs: true,
          quoted_delimiter: false,
        },
        errors: None,
      }],
    },
    TestCase {
      description: "no-quote tab multi-spaced",
      input: "<<-   EOF\n123\nEOF",
      steps: vec![TestStep {
        description: "opening delimiter",
        token: Ok((loc(0), loc(9), Token::StringDelimiter)),
        context: Context::HereDocString {
          delimiter: "EOF",
          skip_leading_tabs: true,
          quoted_delimiter: false,
        },
        errors: None,
      }],
    },
    TestCase {
      description: "spaced tab",
      input: "<< - EOF\n123\nEOF",
      steps: vec![TestStep {
        description: "opening delimiter",
        token: Ok((loc(0), loc(8), Token::StringDelimiter)),
        context: Context::HereDocString {
          delimiter: "- EOF",
          skip_leading_tabs: false,
          quoted_delimiter: false,
        },
        errors: None,
      }],
    },
    TestCase {
      description: "multi-spaced tab",
      input: "<<   - EOF\n123\nEOF",
      steps: vec![TestStep {
        description: "opening delimiter",
        token: Ok((loc(0), loc(10), Token::StringDelimiter)),
        context: Context::HereDocString {
          delimiter: "- EOF",
          skip_leading_tabs: false,
          quoted_delimiter: false,
        },
        errors: None,
      }],
    },
    TestCase {
      description: "post-spaced tab",
      input: "<< -   EOF\n123\nEOF",
      steps: vec![TestStep {
        description: "opening delimiter",
        token: Ok((loc(0), loc(10), Token::StringDelimiter)),
        context: Context::HereDocString {
          delimiter: "-   EOF",
          skip_leading_tabs: false,
          quoted_delimiter: false,
        },
        errors: None,
      }],
    },
    TestCase {
      description: "quote no-tab spaced",
      input: "<< \"EOF\"\n123\nEOF",
      steps: vec![TestStep {
        description: "opening delimiter",
        token: Ok((loc(0), loc(8), Token::StringDelimiter)),
        context: Context::HereDocString {
          delimiter: "EOF",
          skip_leading_tabs: false,
          quoted_delimiter: true,
        },
        errors: None,
      }],
    },
    TestCase {
      description: "quote no-tab not-spaced",
      input: "<<\"EOF\"\n123\nEOF",
      steps: vec![TestStep {
        description: "opening delimiter",
        token: Ok((loc(0), loc(7), Token::StringDelimiter)),
        context: Context::HereDocString {
          delimiter: "EOF",
          skip_leading_tabs: false,
          quoted_delimiter: true,
        },
        errors: None,
      }],
    },
    TestCase {
      description: "quote no-tab multi-spaced",
      input: "<<   \"EOF\"\n123\nEOF",
      steps: vec![TestStep {
        description: "opening delimiter",
        token: Ok((loc(0), loc(10), Token::StringDelimiter)),
        context: Context::HereDocString {
          delimiter: "EOF",
          skip_leading_tabs: false,
          quoted_delimiter: true,
        },
        errors: None,
      }],
    },
    TestCase {
      description: "quote tab spaced",
      input: "<<- \"EOF\"\n123\nEOF",
      steps: vec![TestStep {
        description: "opening delimiter",
        token: Ok((loc(0), loc(9), Token::StringDelimiter)),
        context: Context::HereDocString {
          delimiter: "EOF",
          skip_leading_tabs: true,
          quoted_delimiter: true,
        },
        errors: None,
      }],
    },
    TestCase {
      description: "quote tab not-spaced",
      input: "<<-\"EOF\"\n123\nEOF",
      steps: vec![TestStep {
        description: "opening delimiter",
        token: Ok((loc(0), loc(8), Token::StringDelimiter)),
        context: Context::HereDocString {
          delimiter: "EOF",
          skip_leading_tabs: true,
          quoted_delimiter: true,
        },
        errors: None,
      }],
    },
    TestCase {
      description: "quote tab multi-spaced",
      input: "<<-   \"EOF\"\n123\nEOF",
      steps: vec![TestStep {
        description: "opening delimiter",
        token: Ok((loc(0), loc(11), Token::StringDelimiter)),
        context: Context::HereDocString {
          delimiter: "EOF",
          skip_leading_tabs: true,
          quoted_delimiter: true,
        },
        errors: None,
      }],
    },
    TestCase {
      description: "err missing delimiter",
      input: "<<\n123\nEOF",
      steps: vec![TestStep {
        description: "opening delimiter",
        token: Err((loc(0), loc(2), Error::MissingHeredocDelimiter)),
        context: Context::General,
        errors: None,
      }],
    },
    TestCase {
      description: "unterminated double-quoted delimiter",
      input: "<< \"EOF\n123\nEOF",
      steps: vec![TestStep {
        description: "opening delimiter",
        token: Err((loc(0), loc(7), Error::UnterminatedHeredocDelimiter)),
        context: Context::General,
        errors: None,
      }],
    },
    TestCase {
      description: "unterminated single-quoted delimiter",
      input: "<< 'EOF\n123\nEOF",
      steps: vec![TestStep {
        description: "opening delimiter",
        token: Err((loc(0), loc(7), Error::UnterminatedHeredocDelimiter)),
        context: Context::General,
        errors: None,
      }],
    },
    TestCase {
      description: "quote mismatch single-double",
      input: "<< \"EOF'\n123\nEOF",
      steps: vec![TestStep {
        description: "opening delimiter",
        token: Err((loc(0), loc(8), Error::UnterminatedHeredocDelimiter)),
        context: Context::General,
        errors: None,
      }],
    },
    TestCase {
      description: "quote mismatch double-single",
      input: "<< 'EOF\"\nfoo\nEOF",
      steps: vec![TestStep {
        description: "opening delimiter",
        token: Err((loc(0), loc(8), Error::UnterminatedHeredocDelimiter)),
        context: Context::General,
        errors: None,
      }],
    },
  ];

  test(tests);
}

#[test]
fn string_simple() {
  let tests = vec![
    TestCase {
      description: "single line string",
      input: r#"➖"123"➖"#,
      steps: vec![
        TestStep {
          description: "string start",
          token: Ok((loc(0), loc(1), Token::StringDelimiter)),
          context: Context::String,
          errors: None,
        },
        TestStep {
          description: "string literal",
          token: Ok((loc(1), loc(4), Token::String(StringLiteral::Escaped("123")))),
          context: Context::String,
          errors: None,
        },
        TestStep {
          description: "string end",
          token: Ok((loc(4), loc(5), Token::StringDelimiter)),
          context: Context::General,
          errors: None,
        },
      ],
    },
    TestCase {
      description: "escaped single line string",
      input: r#"➖"123\n"➖"#,
      steps: vec![
        TestStep {
          description: "string start",
          token: Ok((loc(0), loc(1), Token::StringDelimiter)),
          context: Context::String,
          errors: None,
        },
        TestStep {
          description: "string literal",
          token: Ok((loc(1), loc(6), Token::String(StringLiteral::Escaped("123\\n")))),
          context: Context::String,
          errors: None,
        },
        TestStep {
          description: "string end",
          token: Ok((loc(6), loc(7), Token::StringDelimiter)),
          context: Context::General,
          errors: None,
        },
      ],
    },
  ];

  test(tests)
}

#[test]
fn string_interpolation() {
  let tests = vec![
    TestCase {
      description: "simple interpolation",
      input: r#"➖"123${}"➖"#,
      steps: vec![
        TestStep {
          description: "string start",
          token: Ok((loc(0), loc(1), Token::StringDelimiter)),
          context: Context::String,
          errors: None,
        },
        TestStep {
          description: "string literal",
          token: Ok((loc(1), loc(4), Token::String(StringLiteral::Escaped("123")))),
          context: Context::String,
          errors: None,
        },
        TestStep {
          description: "opening interpolation delimiter",
          token: Ok((loc(4), loc(6), Token::StringInterpolation)),
          context: Context::StringInterpolation,
          errors: None,
        },
        TestStep {
          description: "closing interpolation delimiter",
          token: Ok((loc(6), loc(7), Token::StringInterpolation)),
          context: Context::String,
          errors: None,
        },
      ],
    },
    TestCase {
      description: "simple directive",
      input: r#"➖"123%{}"➖"#,
      steps: vec![
        TestStep {
          description: "string start",
          token: Ok((loc(0), loc(1), Token::StringDelimiter)),
          context: Context::String,
          errors: None,
        },
        TestStep {
          description: "string literal",
          token: Ok((loc(1), loc(4), Token::String(StringLiteral::Escaped("123")))),
          context: Context::String,
          errors: None,
        },
        TestStep {
          description: "opening directive delimiter",
          token: Ok((loc(4), loc(6), Token::StringDirective)),
          context: Context::StringDirective,
          errors: None,
        },
        TestStep {
          description: "closing directive delimiter",
          token: Ok((loc(6), loc(7), Token::StringDirective)),
          context: Context::String,
          errors: None,
        },
      ],
    },
  ];

  test(tests)
}

#[test]
fn heredoc() {
  let tests = vec![
    TestCase {
      description: "not leading tabs",
      input: "<<EOF\n\n  123\n    456\nEOF",
      steps: vec![
        TestStep {
          description: "opening delimiter",
          token: Ok((loc(0), loc(5), Token::StringDelimiter)),
          context: Context::HereDocString {
            delimiter: "EOF",
            skip_leading_tabs: false,
            quoted_delimiter: false,
          },
          errors: None,
        },
        TestStep {
          description: "string literal",
          token: Ok((
            Location::new(1, 0, 6),
            Location::new(4, 0, 21),
            Token::HeredocString(vec!["\n", "  123\n", "    456\n"]),
          )),
          context: Context::HereDocStringEnd {
            start: Location::new(4, 0, 21),
            end: Location::new(4, 3, 24),
          },
          errors: None,
        },
        TestStep {
          description: "closing delimiter",
          token: Ok((Location::new(4, 0, 21), Location::new(4, 3, 24), Token::StringDelimiter)),
          context: Context::General,
          errors: None,
        },
      ],
    },
    TestCase {
      description: "leading tabs, empty line",
      input: "<<-EOF\n\n  123\n    456\nEOF",
      steps: vec![
        TestStep {
          description: "opening delimiter",
          token: Ok((loc(0), loc(6), Token::StringDelimiter)),
          context: Context::HereDocString {
            delimiter: "EOF",
            skip_leading_tabs: true,
            quoted_delimiter: false,
          },
          errors: None,
        },
        TestStep {
          description: "string literal",
          token: Ok((
            Location::new(1, 0, 7),
            Location::new(4, 0, 22),
            Token::HeredocString(vec!["\n", "123\n", "  456\n"]),
          )),
          context: Context::HereDocStringEnd {
            start: Location::new(4, 0, 22),
            end: Location::new(4, 3, 25),
          },
          errors: None,
        },
        TestStep {
          description: "closing delimiter",
          token: Ok((Location::new(4, 0, 22), Location::new(4, 3, 25), Token::StringDelimiter)),
          context: Context::General,
          errors: None,
        },
      ],
    },
    TestCase {
      description: "quoted interpolation",
      input: "<<'EOF'\n\n  ${123}\n    456\nEOF",
      steps: vec![
        TestStep {
          description: "opening delimiter",
          token: Ok((loc(0), loc(7), Token::StringDelimiter)),
          context: Context::HereDocString {
            delimiter: "EOF",
            skip_leading_tabs: false,
            quoted_delimiter: true,
          },
          errors: None,
        },
        TestStep {
          description: "string literal",
          token: Ok((
            Location::new(1, 0, 8),
            Location::new(4, 0, 26),
            Token::HeredocString(vec!["\n", "  ${123}\n", "    456\n"]),
          )),
          context: Context::HereDocStringEnd {
            start: Location::new(4, 0, 26),
            end: Location::new(4, 3, 29),
          },
          errors: None,
        },
        TestStep {
          description: "closing delimiter",
          token: Ok((Location::new(4, 0, 26), Location::new(4, 3, 29), Token::StringDelimiter)),
          context: Context::General,
          errors: None,
        },
      ],
    },
    TestCase {
      description: "interpolation",
      input: "<<EOF\n\n  ${}\n    456\nEOF",
      steps: vec![
        TestStep {
          description: "opening delimiter",
          token: Ok((loc(0), loc(5), Token::StringDelimiter)),
          context: Context::HereDocString {
            delimiter: "EOF",
            skip_leading_tabs: false,
            quoted_delimiter: false,
          },
          errors: None,
        },
        TestStep {
          description: "string literal",
          token: Ok((
            Location::new(1, 0, 6),
            Location::new(2, 2, 9),
            Token::HeredocString(vec!["\n", "  "]),
          )),
          context: Context::HereDocString {
            delimiter: "EOF",
            skip_leading_tabs: false,
            quoted_delimiter: false,
          },
          errors: None,
        },
        TestStep {
          description: "interpolation start",
          token: Ok((Location::new(2, 2, 9), Location::new(2, 4, 11), Token::StringInterpolation)),
          context: Context::StringInterpolation,
          errors: None,
        },
        TestStep {
          description: "interpolation end",
          token: Ok((Location::new(2, 4, 11), Location::new(2, 5, 12), Token::StringInterpolation)),
          context: Context::HereDocString {
            delimiter: "EOF",
            skip_leading_tabs: false,
            quoted_delimiter: false,
          },
          errors: None,
        },
        TestStep {
          description: "string literal",
          token: Ok((
            Location::new(2, 5, 12),
            Location::new(4, 0, 21),
            Token::HeredocString(vec!["\n", "    456\n"]),
          )),
          context: Context::HereDocStringEnd {
            start: Location::new(4, 0, 21),
            end: Location::new(4, 3, 24),
          },
          errors: None,
        },
        TestStep {
          description: "closing delimiter",
          token: Ok((Location::new(4, 0, 21), Location::new(4, 3, 24), Token::StringDelimiter)),
          context: Context::General,
          errors: None,
        },
      ],
    },
    TestCase {
      description: "directive",
      input: "<<EOF\n\n  %{}\n    456\nEOF",
      steps: vec![
        TestStep {
          description: "opening delimiter",
          token: Ok((loc(0), loc(5), Token::StringDelimiter)),
          context: Context::HereDocString {
            delimiter: "EOF",
            skip_leading_tabs: false,
            quoted_delimiter: false,
          },
          errors: None,
        },
        TestStep {
          description: "string literal",
          token: Ok((
            Location::new(1, 0, 6),
            Location::new(2, 2, 9),
            Token::HeredocString(vec!["\n", "  "]),
          )),
          context: Context::HereDocString {
            delimiter: "EOF",
            skip_leading_tabs: false,
            quoted_delimiter: false,
          },
          errors: None,
        },
        TestStep {
          description: "interpolation start",
          token: Ok((Location::new(2, 2, 9), Location::new(2, 4, 11), Token::StringDirective)),
          context: Context::StringDirective,
          errors: None,
        },
        TestStep {
          description: "interpolation end",
          token: Ok((Location::new(2, 4, 11), Location::new(2, 5, 12), Token::StringDirective)),
          context: Context::HereDocString {
            delimiter: "EOF",
            skip_leading_tabs: false,
            quoted_delimiter: false,
          },
          errors: None,
        },
        TestStep {
          description: "string literal",
          token: Ok((
            Location::new(2, 5, 12),
            Location::new(4, 0, 21),
            Token::HeredocString(vec!["\n", "    456\n"]),
          )),
          context: Context::HereDocStringEnd {
            start: Location::new(4, 0, 21),
            end: Location::new(4, 3, 24),
          },
          errors: None,
        },
        TestStep {
          description: "closing delimiter",
          token: Ok((Location::new(4, 0, 21), Location::new(4, 3, 24), Token::StringDelimiter)),
          context: Context::General,
          errors: None,
        },
      ],
    },
  ];

  test(tests);
}
