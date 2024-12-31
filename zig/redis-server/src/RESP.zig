const std = @import("std");
const debug = std.debug;

const CRLF: []const u8 = "\\r\\n";
const CRLF_LEN: usize = CRLF.len;
comptime {
    // assert that the length of CRLF is 4
    std.debug.assert(CRLF_LEN == 4);
}

pub const DataType = enum {
    SimpleString,
    Error,
    Integer,
    BulkString,

    fn to_string(self: DataType) !u8 {
        return switch (self) {
            DataType.SimpleString => @intCast('+'),
            DataType.Error => @intCast('-'),
            DataType.Integer => @intCast(':'),
            DataType.BulkString => @intCast('$'),
        };
    }

    fn from_string(data_type: u8) !DataType {
        return switch (data_type) {
            '+' => DataType.SimpleString,
            '-' => DataType.Error,
            ':' => DataType.Integer,
            '$' => DataType.BulkString,
        };
    }
};

pub const Message = struct {
    data_type: DataType,
    content: []const u8,
};

pub fn deserialise(raw: []u8) !Message {
    if (raw.len == 0) {
        return error.EmptyRequest;
    }

    const data_type = try DataType.from_string(raw[0]);
    switch (data_type) {
        DataType.SimpleString, DataType.Error, DataType.Integer => {
            const value = raw[1 .. raw.len - CRLF_LEN];
            const last_2_bytes = raw[raw.len - CRLF_LEN ..];
            if (!std.mem.eql(u8, last_2_bytes, CRLF)) {
                return error.InvalidTerminator_ShouldBeCRLF;
            }
            return Message{ .data_type = data_type, .content = value };
        },
        DataType.BulkString => {
            const value = raw[1 .. raw.len - CRLF_LEN];
            const terminator = raw[raw.len - CRLF_LEN ..];
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

            return Message{ .data_type = data_type, .content = data };
        },
    }
}

fn serialise(allocator: std.mem.Allocator, m: Message) ![]u8 {
    const data_type = try m.data_type.to_string();
    const content = m.content;

    var buf: []u8 = undefined;
    switch (m.data_type) {
        DataType.SimpleString, DataType.Error, DataType.Integer => {
            buf = try allocator.alloc(u8, 1 + content.len + CRLF_LEN); // data_type + content + crlf
            _ = try std.fmt.bufPrint(buf, "{c}{s}{s}", .{ data_type, content, CRLF });
        },
        else => {
            return error.UnsupportedDataType;
        },
    }
    return buf;
}

// write test for serialise fn

test "serialise simple string" {
    const allocator = std.heap.page_allocator;
    const m = Message{ .data_type = DataType.SimpleString, .content = "hello" };
    const got = try serialise(allocator, m);
    const want = "+hello\\r\\n";
    try std.testing.expectEqualStrings(got, want);
}

test "serialise error" {
    const allocator = std.heap.page_allocator;
    const m = Message{ .data_type = DataType.Error, .content = "error" };
    const got = try serialise(allocator, m);
    const want = "-error\\r\\n";
    try std.testing.expectEqualStrings(got, want);
}

test "serialise integer" {
    const allocator = std.heap.page_allocator;
    const m = Message{ .data_type = DataType.Integer, .content = "123" };
    const got = try serialise(allocator, m);
    const want = ":123\\r\\n";
    try std.testing.expectEqualStrings(got, want);
}
