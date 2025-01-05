const Message = @import("RESP.zig").Message;
const std = @import("std");

const CmdHandler = @This();

// command -> handler
// hash map of string to fn(Message) Message
// Example:
// - "PING" -> fn(Message) Message { return Message{ .data_type = DataType.SimpleString, .content = "PONG" }; }
// - "ECHO" -> fn(Message) Message { return Message{ .data_type = DataType.SimpleString, .content = message.content }; }

handlers: std.HashMap([]const u8, fn (Message) Message),

// init allocates memory for the hashmap
// and returns a CommandRouter instance.
// Caller should call self.deinit to free the memory.
pub fn init(allocator: std.mem.Allocator) CmdHandler {
    const handlers = std.AutoHashMap([]const u8, fn (Message) Message).init(allocator);
    return CmdHandler{ .handlers = handlers };
}

pub fn deinit(self: *CmdHandler) void {
    self.handlers.deinit();
}

pub fn register_handler(self: *CmdHandler, command: []const u8, handler_fn: fn (Message) Message) void {
    self.handlers.put(command, handler_fn);
}

pub fn handle(self: *CmdHandler, message: Message) !Message {
    const command = message.value_raw;
    const handler = try self.find_handler(command);
    return handler(message);
}

fn find_handler(self: *CmdHandler, command: []const u8) ?fn (Message) Message {
    return self.handlers.get(command);
}
