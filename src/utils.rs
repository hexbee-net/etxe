#[allow(dead_code)]
fn print_error(reason: &str, description: &str) {
    let indent: &str = &*format!("{}", "│ ".red());

    let width = textwrap::termwidth() - 4;
    let options = textwrap::Options::new(width)
        .initial_indent(indent)
        .subsequent_indent(indent);

    println!("{}", "╷".red());
    eprintln!("{} {}", "│ Error:".red(), reason);
    eprintln!("{}", "│".red());
    eprintln!("{}", textwrap::fill(description, options));
    eprintln!("{}", "╵".red());
    eprintln!();
}
