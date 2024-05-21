use std::env;
use std::fs::File;
use std::io::Read;
use std::path::Path;
use std::process::exit;

fn main() {
    let args: Vec<String> = env::args().collect();

    if args.len() < 2 {
        fatal(format!("usage: {} [file ...]", args[0]).as_str());
    }

    args[1..].iter().for_each(|f|{
        // Create a path to the desired file
        let path = Path::new(&f);

        // Open the path in read-only mode, returns `io::Result<File>`
        let mut file = match File::open(&path) {
            Err(why) => fatal(format!("couldn't open {}: {}", path.display(), why).as_str()),
            Ok(file) => file,
        };

        // Read the file contents into a string, returns `io::Result<usize>`
        let mut s = String::new();
        match file.read_to_string(&mut s) {
            Err(why) => fatal(format!("couldn't read {}: {}", path.display(), why).as_str()),
            Ok(_) => {print!("{}", s);},
        }
    });

    // `file` goes out of scope, and the "hello.txt" file gets closed
}

fn fatal(msg: &str) -> ! {
    eprintln!("{}", msg);
    exit(1)
}
