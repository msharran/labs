const std = @import("std");
const expect = std.testing.expect;

pub fn main() !void {
    std.debug.print("Hello, world!\n", .{});
}

test "always succeeds" {
    try expect(true);
}

test "while with continue statement" {
    var i: u32 = 0;
    while (i < 10) : (i += 1) {
        std.debug.print("i = {}\n", .{i});
    }

    try expect(i == 10);
}

test "for with array" {
    var arr = [_]u8{ 'a', 'b', 'c' };
    for (arr, 0..) |c, i| {
        std.debug.print("c, i = {}, {}\n", .{ c, i });
    }
}

fn addFive(x: i32) i32 {
    return x + 5;
}

test "addFive" {
    try expect(addFive(5) == 10);
}
