const std = @import("std");
const server = @import("server.zig");

pub fn main() !void {
    const address = try std.net.Address.parseIp("127.0.0.1", 6397);

    server.listen_and_serve(address) catch |err| {
        std.debug.print("Failed to start server: {}\n", .{err});
    };
}
