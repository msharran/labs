const std = @import("std");
const http = std.http;
const print = std.debug.print;

pub fn main() !void {
    // Create an allocator.
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer std.debug.assert(gpa.deinit() == .ok);
    const allocator = gpa.allocator();

    // Create an HTTP client.
    var client = http.Client{ .allocator = allocator };
    defer client.deinit();

    const uri = try std.Uri.parse("https://jsonplaceholder.typicode.com/todos/1");

    // allocate a byte slice buffer to store the server response headers
    // and pass it to the request options
    const buffer = try allocator.alloc(u8, 2048);
    defer allocator.free(buffer);

    var req = try client.open(.GET, uri, .{ .server_header_buffer = buffer });
    defer req.deinit();

    try req.send();
    try req.wait();

    print("Status: {}\n", .{req.response.status});
    // parse the body as a string

    const body_buf = try allocator.alloc(u8, req.response.content_length orelse 1024);
    defer allocator.free(body_buf);

    _ = try req.read(body_buf);
    print("Body: {s}\n", .{body_buf});
}
