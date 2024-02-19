struct Cli {
    pattern: String,
    file: std::path::PathBuf,
}

fn main() {
    let pattern = std::env::args().nth(1).expect("no pattern given");
    let file = std::env::args().nth(2).expect("no file given");

    let args = Cli {
        pattern,
        file: std::path::PathBuf::from(file),
    };

    let content = std::fs::read_to_string(&args.file)
        .expect("could not read file");

    for line in content.lines() {
        if line.contains(&args.pattern) {
            println!("{}", line);
        }
    }

}
