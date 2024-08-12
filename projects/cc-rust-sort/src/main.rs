fn main() {
    greet();
}

fn basic_datatypes() {
    let i1 = 1;
    let i2 = i1;

    println!("{}", i1);
}

fn greet() {
    let name = String::from("Sharran");
    let n = name.clone();
    println!("inner name, {} {:p}", &name, &name);
    println!("inner n, {} {:p}", &n, &n);
}
