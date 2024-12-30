const std = @import("std");
const debug = std.debug;
const net = std.net;
const posix = std.posix;

const CRLF: []const u8 = "\\r\\n";
const CRLF_LEN: usize = CRLF.len;
comptime {
    // assert that the length of CRLF is 4
    std.debug.assert(CRLF_LEN == 4);
}

const DataType = enum { SimpleString, Error, Integer, BulkString };

const Message = struct {
    data_type: DataType,
    data: []const u8,
};

fn parse_data_type(first_byte: u8) !DataType {
    return switch (first_byte) {
        '+' => DataType.SimpleString,
        '-' => DataType.Error,
        ':' => DataType.Integer,
        '$' => DataType.BulkString,
        else => error.InvalidDataTypeInFirstByte,
    };
}

fn get_value(data_type: DataType, request: []const u8) !Message {
    switch (data_type) {
        DataType.SimpleString, DataType.Error, DataType.Integer => {
            const value = request[1 .. request.len - CRLF_LEN];
            const last_2_bytes = request[request.len - CRLF_LEN ..];
            if (!std.mem.eql(u8, last_2_bytes, CRLF)) {
                return error.InvalidTerminator_ShouldBeCRLF;
            }
            return Message{ .data_type = data_type, .data = value };
        },
        DataType.BulkString => {
            const value = request[1 .. request.len - CRLF_LEN];
            const terminator = request[request.len - CRLF_LEN ..];
            if (!std.mem.eql(u8, terminator, CRLF)) {
                return error.InvalidTerminator_ShouldBeCRLF;
            }

            // value is split by CRLF
            // first is the length of the string
            // second is the string itself
            // get the length, parse to int
            // get the string using the length
            const backslash = std.mem.indexOf(u8, value, "\\");
            if (backslash == null) {
                return error.InvalidBulkString_CRLFMissing;
            }

            const crlf = value[backslash.? .. backslash.? + CRLF_LEN];
            if (!std.mem.eql(u8, crlf, CRLF)) {
                return error.InvalidBulkString_InvalidCRLF;
            }

            const len_str = value[0..backslash.?];
            const len = try std.fmt.parseInt(usize, len_str, 10);

            const crlf_end = backslash.? + CRLF_LEN;
            const data = value[crlf_end .. crlf_end + len];
            if (data.len != len) {
                return error.InvalidBulkString_LengthIncorrect;
            }

            return Message{ .data_type = data_type, .data = data };
        },
    }
}

pub fn serialise(request: []const u8) !Message {
    if (request.len == 0) {
        return error.EmptyRequest;
    }

    const first_byte = request[0];

    const data_type = try parse_data_type(first_byte);
    return try get_value(data_type, request);
}

pub fn main() !void {
    // var request: []const u8 = "+foo\r\n";
    // var data = try serialise(request);
    // debug.print("data_type: {}\n", .{data.data_type});
    // debug.print("val_str: {?s}\n", .{data.val_str});
    // debug.print("val_int: {?d}\n----\n", .{data.val_int});
    //
    // request = "-ERR bar\r\n";
    // data = try serialise(request);
    // debug.print("data_type: {}\n", .{data.data_type});
    // debug.print("val_str: {?s}\n", .{data.val_str});
    // debug.print("val_int: {?d}\n----\n", .{data.val_int});
    //
    // request = ":9876\r\n";
    // data = try serialise(request);
    // debug.print("data_type: {}\n", .{data.data_type});
    // debug.print("val_str: {?s}\n", .{data.val_str});
    // debug.print("val_int: {?d}\n----\n", .{data.val_int});
    //
    // request = "$15\r\nZig is Amazing!\r\n";
    // data = try serialise(request);
    // debug.print("data_type: {}\n", .{data.data_type});
    // debug.print("val_str: {?s}\n", .{data.val_str});
    // debug.print("val_int: {?d}\n----\n", .{data.val_int});

    const address = try std.net.Address.parseIp("127.0.0.1", 6397);
    const tpe: u32 = posix.SOCK.STREAM;
    const protocol = posix.IPPROTO.TCP;
    const listener = try posix.socket(address.any.family, tpe, protocol);
    defer posix.close(listener);

    try posix.setsockopt(listener, posix.SOL.SOCKET, posix.SO.REUSEADDR, &std.mem.toBytes(@as(c_int, 1)));
    try posix.bind(listener, &address.any, address.getOsSockLen());
    try posix.listen(listener, 128);

    var buf: [1024]u8 = undefined;
    while (true) {
        var client_address: net.Address = undefined;
        var client_address_len: posix.socklen_t = @sizeOf(net.Address);

        const socket = posix.accept(listener, &client_address.any, &client_address_len, 0) catch |err| {
            std.debug.print("error accept: {}\n", .{err});
            continue;
        };
        defer posix.close(socket);

        std.debug.print("{} connected\n", .{client_address});

        const read = posix.read(socket, &buf) catch |err| {
            std.debug.print("error reading: {}\n", .{err});
            continue;
        };

        if (read == 0) {
            continue;
        }

        const resp = serialise(buf[0..read]) catch |err| {
            std.debug.print("error serialising: {}\n", .{err});
            continue;
        };

        std.debug.print("read data_type: {}\n", .{resp.data_type});

        write(socket, resp.data) catch |err| {
            std.debug.print("error writing: {}\n", .{err});
        };
    }
}

fn write(socket: posix.socket_t, msg: []const u8) !void {
    var pos: usize = 0;
    while (pos < msg.len) {
        const written = try posix.write(socket, msg[pos..]);
        if (written == 0) {
            return error.Closed;
        }
        pos += written;
    }
}
