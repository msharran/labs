const std = @import("std");
const debug = std.debug;
const assert = debug.assert;

const Proto = @This();

const CRLF: []const u8 = "\\r\\n";
const CRLF_LEN: usize = CRLF.len;
comptime {
    // assert that the length of CRLF is 4
    std.debug.assert(CRLF_LEN == 4);
}

arena: *std.heap.ArenaAllocator,

pub fn init(arena: *std.heap.ArenaAllocator) Proto {
    return Proto{ .arena = arena };
}

pub const DataType = enum {
    SimpleString,
    Error,
    Integer,
    BulkString,
    Array,

    fn yes(char: u8) bool {
        return switch (char) {
            '+', '-', ':', '$', '*' => true,
            else => false,
        };
    }

    fn toChar(self: DataType) !u8 {
        return switch (self) {
            DataType.SimpleString => @intCast('+'),
            DataType.Error => @intCast('-'),
            DataType.Integer => @intCast(':'),
            DataType.BulkString => @intCast('$'),
            DataType.Array => @intCast('*'),
        };
    }

    fn fromChar(data_type: u8) !DataType {
        return switch (data_type) {
            '+' => DataType.SimpleString,
            '-' => DataType.Error,
            ':' => DataType.Integer,
            '$' => DataType.BulkString,
            '*' => DataType.Array,
            else => return error.InvalidDataType,
        };
    }
};

pub const ValueTag = enum {
    list,
    single,
};

pub const Value = union(ValueTag) {
    list: std.ArrayList(Message),
    single: []const u8,
};

pub const Message = struct {
    type: DataType,
    value: Value,

    pub fn init(t: DataType, v: []const u8) Message {
        return Message{ .type = t, .value = Value{ .single = v } };
    }

    pub fn initList(t: DataType, v: std.ArrayList(Message)) Message {
        return Message{ .type = t, .value = Value{ .list = v } };
    }
};

pub fn deserialise(self: Proto, raw: []const u8) !Message {
    if (raw.len == 0) {
        return error.EmptyRequest;
    }

    const data_type = try DataType.fromChar(raw[0]);
    switch (data_type) {
        DataType.SimpleString, DataType.Error, DataType.Integer => {
            const value = raw[1 .. raw.len - CRLF_LEN];
            const last_2_bytes = raw[raw.len - CRLF_LEN ..];
            if (!std.mem.eql(u8, last_2_bytes, CRLF)) {
                return error.InvalidTerminator_ShouldBeCRLF;
            }
            return Message.init(data_type, value);
        },
        DataType.BulkString => {
            var parts = std.mem.splitSequence(u8, raw, CRLF);
            const length_part = parts.first();

            // first char is "$"
            // rest is the length of the string in len_bytes
            // get the length of the string and parse to int
            const len = try std.fmt.parseInt(usize, length_part[1..], 10);

            const string_part = parts.next() orelse {
                return error.MissingContent;
            };

            if (string_part.len != len) {
                return error.ContentLengthMismatch;
            }

            return Message.init(data_type, string_part);
        },
        DataType.Array => {
            const raw_msgs = try self.toOwnedMessages(raw);

            // init an array list of messages
            var list = std.ArrayList(Message).init(self.arena.allocator());

            for (raw_msgs) |msg| {
                const m = try self.deserialise(msg);
                try list.append(m);
            }

            return Message.initList(data_type, list);
        },
    }
}

/// merges merges parts of items in an array
/// into a single item
/// e.g. "*2\r\n$4\r\nECHO\r\n$5\r\nhello\r\n"
/// item 1: "$4\r\nECHO\r\n" => "ECHO\r\n"
/// item 2: "$5\r\nhello\r\n" => "hello\r\n"
fn toOwnedMessages(self: Proto, raw: []const u8) ![][]u8 {
    var parts = std.mem.splitSequence(u8, raw, CRLF);
    const arrlenbytes = parts.first();
    if (arrlenbytes.len != 2) {
        return error.InvalidArrayLength;
    }

    const arr_len = try std.fmt.parseInt(usize, arrlenbytes[1..], 10);

    var commands = std.ArrayList([]u8).init(self.arena.allocator());

    var cmd: ?std.ArrayList(u8) = null;
    while (parts.next()) |part| {
        if (part.len == 0) {
            continue;
        }
        if (DataType.yes(part[0])) {
            if (cmd) |*c| {
                try commands.append(try c.toOwnedSlice());
            }

            // init new raw_content with the data_type part
            // no need deinit since we are converting to owned slice
            cmd = std.ArrayList(u8).init(self.arena.allocator());

            try cmd.?.appendSlice(part);
            try cmd.?.appendSlice(CRLF);
        } else {
            if (cmd == null) {
                return error.MissingDataType;
            }
            try cmd.?.appendSlice(part);
            try cmd.?.appendSlice(CRLF);
        }
    }

    if (cmd) |*c| {
        try commands.append(try c.toOwnedSlice());
    }

    if (arr_len != commands.items.len) {
        return error.ArrayLengthMismatch;
    }
    const s = try commands.toOwnedSlice();
    return s;
}

pub fn serialise(allocator: std.mem.Allocator, m: Message) ![]u8 {
    const data_type = try m.type.toChar();
    const content = m.value;

    var buf: []u8 = undefined;
    switch (m.type) {
        DataType.SimpleString, DataType.Error, DataType.Integer => {
            buf = try allocator.alloc(u8, 1 + content.len + CRLF_LEN); // data_type + content + crlf
            errdefer allocator.free(buf);
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
    const m = Message{ .type = DataType.SimpleString, .value = "hello" };
    const got = try serialise(allocator, m);
    defer allocator.free(got);
    const want = "+hello\\r\\n";
    try std.testing.expectEqualStrings(got, want);
}

test "serialise error" {
    const allocator = std.heap.page_allocator;
    const m = Message{ .type = DataType.Error, .value = "error" };
    const got = try serialise(allocator, m);
    defer allocator.free(got);
    const want = "-error\\r\\n";
    try std.testing.expectEqualStrings(got, want);
}

test "serialise integer" {
    const allocator = std.heap.page_allocator;
    const m = Message{ .type = DataType.Integer, .value = "123" };
    const got = try serialise(allocator, m);
    defer allocator.free(got);
    const want = ":123\\r\\n";
    try std.testing.expectEqualStrings(got, want);
}

test "deserialise bulk_string" {
    const raw = "$4\\r\\nECHO\\r\\n";
    const got = try deserialise(raw);
    const want = Message{ .type = DataType.BulkString, .value = "ECHO" };

    try std.testing.expectEqualStrings(got.value, want.value);
}

test "deserialise first bulk_string" {
    const raw = "$4\\r\\nECHO\\r\\n$5\\r\\nhello\\r\\n";
    const got = try deserialise(raw);

    // should only deserialise the first bulk string
    const want = Message{ .type = DataType.BulkString, .value = "ECHO" };

    try std.testing.expectEqualStrings(got.value, want.value);
}

test "deserialise array" {
    const raw = "*3\\r\\n$4\\r\\nECHO\\r\\n$5\\r\\nhello\\r\\n:123\\r\\n";

    var arena = std.heap.ArenaAllocator.init(std.testing.allocator);
    defer arena.deinit();

    const resp = Proto{ .arena = &arena };
    const got = try resp.deserialise(raw);

    assert(got.type == DataType.Array);
    assert(got.value.list.items.len == 3);
    assert(got.value.list.items[0].type == DataType.BulkString);
    const gotval0 = got.value.list.items[0].value.single;
    try std.testing.expectEqualStrings("ECHO", gotval0);
    assert(got.value.list.items[1].type == DataType.BulkString);
    try std.testing.expectEqualStrings(got.value.list.items[1].value.single, "hello");
    assert(got.value.list.items[2].type == DataType.Integer);
    try std.testing.expectEqualStrings(got.value.list.items[2].value.single, "123");
}
