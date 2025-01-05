const std = @import("std");
const debug = std.debug;

const CRLF: []const u8 = "\\r\\n";
const CRLF_LEN: usize = CRLF.len;
comptime {
    // assert that the length of CRLF is 4
    std.debug.assert(CRLF_LEN == 4);
}

pub const Command = enum {
    Ping,
    Echo,
};

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

arena: *std.heap.ArenaAllocator,

pub const Message = struct { type: DataType, value_raw: ?[]const u8 = null, next: ?*Message = null };

pub fn deserialise(self: @This(), raw: []const u8) !Message {
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
            return Message{ .type = data_type, .value_raw = value };
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

            return Message{ .type = data_type, .value_raw = string_part };
        },
        DataType.Array => {
            const msg = try self.toOwnedMessages(raw);
            debug.print("commands: {s}\n", .{msg});
            // var msg = Message{ .type = DataType.Array };
            // for (messages) |cmd| {
            // debug.print("cmd: {s}\n", .{cmd});
            // var m = try deserialise(cmd);
            // msg.next = &m;
            // msg = m;
            // }

            return error.UnsupportedDataType;
        },
    }
}

/// mergeItemParts merges parts of items in an array
/// into a single item
/// e.g. "*2\r\n$4\r\nECHO\r\n$5\r\nhello\r\n"
/// item 1: "$4\r\nECHO\r\n" => "ECHO\r\n"
/// item 2: "$5\r\nhello\r\n" => "hello\r\n"
fn toOwnedMessages(self: @This(), raw: []const u8) ![][]u8 {
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
    debug.print("commands: {s}\n", .{s});
    return s;
}

pub fn serialise(allocator: std.mem.Allocator, m: Message) ![]u8 {
    const data_type = try m.type.toChar();
    const content = m.value_raw;

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
    const m = Message{ .type = DataType.SimpleString, .value_raw = "hello" };
    const got = try serialise(allocator, m);
    defer allocator.free(got);
    const want = "+hello\\r\\n";
    try std.testing.expectEqualStrings(got, want);
}

test "serialise error" {
    const allocator = std.heap.page_allocator;
    const m = Message{ .type = DataType.Error, .value_raw = "error" };
    const got = try serialise(allocator, m);
    defer allocator.free(got);
    const want = "-error\\r\\n";
    try std.testing.expectEqualStrings(got, want);
}

test "serialise integer" {
    const allocator = std.heap.page_allocator;
    const m = Message{ .type = DataType.Integer, .value_raw = "123" };
    const got = try serialise(allocator, m);
    defer allocator.free(got);
    const want = ":123\\r\\n";
    try std.testing.expectEqualStrings(got, want);
}

test "deserialise bulk_string" {
    const raw = "$4\\r\\nECHO\\r\\n";
    const got = try deserialise(raw);
    const want = Message{ .type = DataType.BulkString, .value_raw = "ECHO" };

    try std.testing.expectEqualStrings(got.value_raw, want.value_raw);
}

test "deserialise first bulk_string" {
    const raw = "$4\\r\\nECHO\\r\\n$5\\r\\nhello\\r\\n";
    const got = try deserialise(raw);

    // should only deserialise the first bulk string
    const want = Message{ .type = DataType.BulkString, .value_raw = "ECHO" };

    try std.testing.expectEqualStrings(got.value_raw, want.value_raw);
}

test "deserialise array" {
    const raw = "*3\\r\\n$4\\r\\nECHO\\r\\n$5\\r\\nhello\\r\\n$4\\r\\nPING\\r\\n";

    var arena = std.heap.ArenaAllocator.init(std.testing.allocator);
    defer arena.deinit();

    const resp = @This(){ .arena = &arena };
    const got = try resp.deserialise(raw);

    try std.testing.expectEqual(DataType.Array, got.type);
    try std.testing.expect(got.value_raw == null);
    try std.testing.expect(got.next != null);

    // test first child
    // var child = got.next;
    // try std.testing.expectEqual(DataType.BulkString, child.?.type);
    // try std.testing.expect(child.?.value_raw != null);
    // debug.print("child: {?}\n", .{child});

    // test second child
    // child = child.?.next;
}
