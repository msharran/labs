const std = @import("std");
const net = std.net;
const posix = std.posix;
const resp = @import("RESP.zig");
const CmdHandler = @import("CmdHandler.zig");

const Server = @This();

const ServeMux = struct {};

/// cmd_handler is responsible for handling the commands received by the server.
cmd_handler: CmdHandler,

/// allocator is used to allocate memory for the server.
allocator: std.mem.Allocator,

/// init allocates memory for the server.
/// Caller should call self.deinit to free the memory.
pub fn init(allocator: std.mem.Allocator) Server {
    const cmd_handler = CmdHandler.init(allocator);
    return Server{
        .cmd_handler = cmd_handler,
        .allocator = allocator,
    };
}

pub fn deinit(self: *Server) void {
    self.cmd_handler.deinit();
}

pub fn listen_and_serve(self: *Server, address: std.net.Address) !void {
    const tpe: u32 = posix.SOCK.STREAM;
    const protocol = posix.IPPROTO.TCP;
    const listener = try posix.socket(address.any.family, tpe, protocol);
    defer posix.close(listener);

    try posix.setsockopt(listener, posix.SOL.SOCKET, posix.SO.REUSEADDR, &std.mem.toBytes(@as(c_int, 1)));
    try posix.bind(listener, &address.any, address.getOsSockLen());
    try posix.listen(listener, 128);

    std.debug.print("Server listening on {}\n", .{address});

    var buf: [1024]u8 = undefined;
    while (true) {
        var client_address: net.Address = undefined;
        var client_address_len: posix.socklen_t = @sizeOf(net.Address);

        const socket = posix.accept(listener, &client_address.any, &client_address_len, 0) catch |err| {
            std.debug.print("error accept: {}\n", .{err});
            continue;
        };
        defer posix.close(socket);

        std.debug.print("client {} connected\n", .{client_address});

        const read = read_all(socket, &buf) catch |err| {
            std.debug.print("error reading: {}\n", .{err});
            continue;
        };

        if (read == 0) {
            continue;
        }

        const msg = resp.deserialise(buf[0..read]) catch |err| {
            std.debug.print("error serialising: {}\n", .{err});
            continue;
        };

        const resp_msg = try self.cmd_handler.handle(msg);

        const allocator = self.allocator;
        const msg_serialised = resp.serialise(allocator, resp_msg) catch |err| {
            std.debug.print("error serialising: {}\n", .{err});
            continue;
        };
        defer allocator.free(msg_serialised);

        write_all(socket, msg_serialised) catch |err| {
            std.debug.print("error writing: {}\n", .{err});
        };
    }
}

fn write_all(socket: posix.socket_t, msg: []const u8) !void {
    var pos: usize = 0;
    while (pos < msg.len) {
        const written = try posix.write(socket, msg[pos..]);
        if (written == 0) {
            return error.Closed;
        }
        pos += written;
    }
}

fn read_all(socket: posix.socket_t, buf: []u8) !usize {
    var pos: usize = 0;
    while (pos < buf.len) {
        const read = try posix.read(socket, buf[pos..]);
        if (read == 0) { // 0 means EOF
            return pos;
        }
        pos += read;
    }
    return pos;
}
