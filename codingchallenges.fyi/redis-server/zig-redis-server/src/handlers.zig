const std = @import("std");
const proto = @import("Proto.zig");

pub fn ping(_: proto.Message) proto.Message {
    return proto.Message{ .type = proto.DataType.SimpleString, .value = .{ .single = "PONG" } };
}
