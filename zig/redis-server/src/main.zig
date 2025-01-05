const std = @import("std");
const server = @import("Server.zig");
const CmdHandler = @import("CmdHandler.zig");

fn pong_handler(_: server.Message) server.Message {
    return server.Message{ .data_type = server.DataType.SimpleString, .content = "PONG" };
}

pub fn main() !void {
    const address = try std.net.Address.parseIp("127.0.0.1", 6397);

    const gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer _ = gpa.deinit();

    const allocator = gpa.allocator();
    var handler: CmdHandler = CmdHandler.init(allocator);
    defer handler.deinit();

    handler.register_handler("PING", pong_handler);

    server.listen_and_serve(address) catch |err| {
        std.debug.print("Failed to start server: {}\n", .{err});
    };
}
