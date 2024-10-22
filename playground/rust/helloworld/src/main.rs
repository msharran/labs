use std::fs;

fn main() {
    fs::read_dir(path)
}



// // #[derive(Debug)]
// // struct Rectangle {
// //     width: u32,
// //     height: u32,
// // }
// //
// // impl Rectangle {
// //     fn area(&self) -> u32 {
// //         self.width * self.height
// //     }
// // }
// //
// // fn main() {
// //     println!("staring main");
// //     let rect = Rectangle {
// //         width: 10,
// //         height: 5,
// //     };
// //     println!("Area of rect = {:?}", rect.area());
// //     dbg!(&rect);
// // }
// //
//
// use std::error::Error;
//
// fn main() -> Result<(), Box<dyn Error>> {
//     let mins = Minutes(5);
//     let secs = Seconds::from(&mins);
//
//     dbg!(&mins);
//     dbg!(&secs);
//
//     let secs = Seconds(120);
//     let mins: Minutes = (&secs).into();
//
//     dbg!(&mins);
//     dbg!(&secs);
//
//     Ok(())
// }
//
// #[derive(Debug)]
// struct Seconds(u32);
//
// #[derive(Debug)]
// struct Minutes(u32);
//
// impl From<&Minutes> for Seconds {
//     fn from(value: &Minutes) -> Self {
//         Self(value.0 * 60)
//     }
// }
//
// impl From<&Seconds> for Minutes {
//     fn from(value: &Seconds) -> Self {
//         Self(value.0 / 60)
//     }
// }
