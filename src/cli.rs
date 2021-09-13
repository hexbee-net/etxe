use clap::{App, Arg, ArgMatches, AppSettings};
// use tracing::{debug, error, info, span, warn, Level};
use tracing_subscriber;

use etxe::parser;

// const VERSION: &'static str = env!("CARGO_PKG_VERSION");

const SUB_CHECK: &str = "check";
const SUB_VALIDATE: &str = "validate";
const SUB_INSPECT: &str = "inspect";

fn main() {
    // install global collector configured based on RUST_LOG env var.
    tracing_subscriber::fmt::init();

    let matches = App::new("Etxe")
        .version(clap::crate_version!())
        .about("IaC done right")
        .setting(AppSettings::SubcommandRequiredElseHelp)
        .arg(
            Arg::with_name("config")
                .short("c")
                .long("config")
                .takes_value(true)
                .value_name("FILE")
                .help("Provides a config file to Etxe"),
        )
        .arg(
            Arg::with_name("debug")
                .short("d")
                .long("debug")
                .multiple(true)
                .help("Turn debugging information on"),
        )
        .arg(
            Arg::with_name("input")
                .short("i")
                .long("input")
                .multiple(true)
                .global(true)
                .takes_value(true)
                .default_value(".")
                .value_name("PATH")
                .help("")
        )
        .subcommand(
            App::new(SUB_CHECK)
                .about("check the syntax of configuration files")
        )
        .subcommand(
            App::new(SUB_VALIDATE)
                .about("TODO")
        )
        .subcommand(
            App::new(SUB_INSPECT)
                .about("TODO")
        ).get_matches();

    let input: Vec<_> = matches.values_of("input").unwrap().collect();

    match matches.subcommand() {
        (SUB_CHECK, Some(subcommand)) => {
            check_input(&input, subcommand)
        }
        (_, None) => {
            println!("{}", matches.usage());
        }
        _ => {
            println!("unrecognized command");
            println!("{}", matches.usage());
        }
    }
}

#[tracing::instrument]
fn check_input<'a>(input: &Vec<&str>, subcommand: &ArgMatches<'a>) {
    println!("Using input: {}", input.join(", "));

    for src in input.iter() {
        parser::parse(src)
    }

    // todo!()
}
