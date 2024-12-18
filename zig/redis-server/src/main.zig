const std = @import("std");
const debug = std.debug;

const DataType = enum { SimpleString, Error, Integer, BulkString };

const Message = struct { data_type: DataType, val_str: ?[]const u8 = null, val_int: ?i64 = null };

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
        DataType.SimpleString, DataType.Error => {
            const value = request[1 .. request.len - 2];
            const last_2_bytes = request[request.len - 2 ..];
            if (!std.mem.eql(u8, last_2_bytes, "\r\n")) {
                return error.InvalidTerminator;
            }
            return Message{ .data_type = data_type, .val_str = value };
        },
        DataType.Integer => {
            const value = request[1 .. request.len - 2];
            const last_2_bytes = request[request.len - 2 ..];
            if (!std.mem.eql(u8, last_2_bytes, "\r\n")) {
                return error.InvalidTerminator;
            }
            // parse value as integer
            const val_int = try std.fmt.parseInt(i64, value, 10);
            const msg = Message{ .data_type = data_type, .val_int = val_int };
            return msg;
        },
        DataType.BulkString => {
            const value = request[1 .. request.len - 2];
            const last_2_bytes = request[request.len - 2 ..];
            if (!std.mem.eql(u8, last_2_bytes, "\r\n")) {
                return error.InvalidTerminator;
            }

            // value is split by CRLF
            // first is the length of the string
            // second is the string itself

            // get the length, parse to int
            // get the string using the length

            const cr = std.mem.indexOf(u8, value, "\r");
            if (cr == null) {
                return error.InvalidBulkString_CRLF;
            }

            const len_str = value[0..cr.?];
            const len = try std.fmt.parseInt(usize, len_str, 10);

            const data = value[cr.? + 2 .. cr.? + 2 + len];
            if (data.len != len) {
                return error.InvalidBulkString_LengthIncorrect;
            }

            return Message{ .data_type = data_type, .val_str = data };
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
    var request: []const u8 = "+foo\r\n";
    var data = try serialise(request);
    debug.print("data_type: {}\n", .{data.data_type});
    debug.print("val_str: {?s}\n", .{data.val_str});
    debug.print("val_int: {?d}\n----\n", .{data.val_int});

    request = "-ERR bar\r\n";
    data = try serialise(request);
    debug.print("data_type: {}\n", .{data.data_type});
    debug.print("val_str: {?s}\n", .{data.val_str});
    debug.print("val_int: {?d}\n----\n", .{data.val_int});

    request = ":9876\r\n";
    data = try serialise(request);
    debug.print("data_type: {}\n", .{data.data_type});
    debug.print("val_str: {?s}\n", .{data.val_str});
    debug.print("val_int: {?d}\n----\n", .{data.val_int});

    request = "$15\r\nZig is Amazing!\r\n";
    data = try serialise(request);
    debug.print("data_type: {}\n", .{data.data_type});
    debug.print("val_str: {?s}\n", .{data.val_str});
    debug.print("val_int: {?d}\n----\n", .{data.val_int});
}
