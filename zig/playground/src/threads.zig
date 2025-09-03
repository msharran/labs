const std = @import("std");

const SomeStruct = struct {
    pool: *std.Thread.Pool,

    pub fn init(allocator: std.mem.Allocator) !SomeStruct {
        var pool = try allocator.create(std.Thread.Pool);
        try pool.init(std.Thread.Pool.Options{ .allocator = allocator, .n_jobs = 5 });
        return .{ .pool = pool };
    }

    pub fn deinit(self: *SomeStruct) void {
        self.pool.deinit();
    }
};

pub fn main() !void {
    const allocator = std.heap.page_allocator;

    var some_struct = try SomeStruct.init(allocator);
    defer some_struct.deinit();

    try some_struct.pool.spawn(work, .{3});
    try some_struct.pool.spawn(work, .{5});
    try some_struct.pool.spawn(work, .{7});
}

fn work(inc: u32) void {
    std.debug.print("Start Inc = {d}\n", .{inc});
    var total: u32 = 0;
    var i: u32 = 0;
    while (i < 100000) : (i += 1) {
        total += inc;
    }
    std.debug.print("Total = {d}, Inc = {d}\n", .{ total, inc });
}
