const std = @import("std");

// https://zig.news/xq/zig-build-explained-part-3-1ima
const fio = @cImport({
    @cInclude("fio.h");
});

pub fn main() !void {
    const args = fio.fio_start_args{
        .threads = 4,
        .workers = 4,
    };

    std.debug.print("foo\n", .{});
    // fio.fio_start(args);
    fio.fio_start(args);
}
