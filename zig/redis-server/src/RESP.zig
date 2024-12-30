const std = @import("std");
const debug = std.debug;

const CRLF: []const u8 = "\\r\\n";
const CRLF_LEN: usize = CRLF.len;
comptime {
    // assert that the length of CRLF is 4
    std.debug.assert(CRLF_LEN == 4);
}

pub const DataType = enum { SimpleString, Error, Integer, BulkString };

pub const Message = struct {
    data_type: DataType,
    content: []const u8,
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
            return Message{ .data_type = data_type, .content = value };
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

            return Message{ .data_type = data_type, .content = data };
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
