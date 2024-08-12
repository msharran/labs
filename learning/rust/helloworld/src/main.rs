#[derive(Debug)]
struct Rectangle {
    width: u32,
    height: u32,
}

impl Rectangle {
    fn area(&self) -> u32 {
        self.width * self.height
    }
}

fn main() {
    println!("staring main");
    let rect = Rectangle {
        width: 10,
        height: 5,
    };
    println!("Area of rect = {:?}", rect.area());
    dbg!(&rect);
}

