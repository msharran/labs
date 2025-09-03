const std = @import("std");

pub fn main() !void {
    // const request = "$3\r\nfoo\r\n";
    const request = "$15\r\nZig is Amazing!\r\n";
    const index = std.mem.indexOf(u8, request, "\r");
    std.debug.print("Index: {?}\n", .{index});

    std.debug.print("Request: {s}\n", .{request[5 .. 5 + 15]});
    std.debug.print("Char: {c}\n", .{request[6]});
}
