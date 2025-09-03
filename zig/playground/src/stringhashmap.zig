const std = @import("std");
const assert = std.debug.assert;
const expect = std.testing.expect;

pub const ObjectStore = struct {
    objects: std.StringHashMap(Value),
    allocator: std.mem.Allocator,

    pub const ValueTag = enum {
        string,
        int,
    };

    pub const Value = union(ValueTag) {
        string: []const u8,
        int: i32,
    };

    pub fn init(allocator: std.mem.Allocator) !*ObjectStore {
        const store = ObjectStore{
            .objects = std.StringHashMap(Value).init(allocator),
            .allocator = allocator,
        };
        const kv_p = try allocator.create(ObjectStore);
        kv_p.* = store;
        return kv_p;
    }

    pub fn deinit(self: *ObjectStore) void {
        var iter = self.objects.iterator();
        while (iter.next()) |kv| {
            const k_p = kv.key_ptr.*;
            self.allocator.free(k_p);
        }
        self.objects.deinit();
        self.allocator.destroy(self);
    }

    pub fn putString(self: *ObjectStore, key: []const u8, value: []const u8) !void {
        const k_p = try self.allocator.dupe(u8, key);
        const v = Value{ .string = value };
        try self.objects.put(k_p, v);
    }

    pub fn getString(self: *ObjectStore, key: []const u8) !?[]const u8 {
        const v = self.objects.get(key);
        if (v == null) {
            return null;
        }

        switch (v.?) {
            .string => return v.?.string,
            else => return error.NotAStringValue,
        }
    }
};

test "inside-struct" {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer {
        const check = gpa.deinit();
        std.debug.print("check: {}\n", .{check});
    }

    const allocator = gpa.allocator();
    var object_store = try ObjectStore.init(allocator);
    defer object_store.deinit();

    const key1 = "key1";
    const key2 = "key2";
    try object_store.putString(key1, "value1");
    try object_store.putString(key2, "value2");

    const got1 = try object_store.getString("key1");
    const got2 = try object_store.getString("key2");

    std.debug.print("got1: {?s}\n", .{got1});
    std.debug.print("got2: {?s}\n", .{got2});
}

test "2" {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer {
        const check = gpa.deinit();
        std.debug.print("check: {}\n", .{check});
    }

    const allocator = gpa.allocator();

    const store = try ObjectStore.init(allocator);
    defer store.deinit();

    var capacity = store.*.objects.capacity();
    var count = store.*.objects.count();
    std.debug.print("count: {d} cap: {d}\n", .{ count, capacity });

    try store.putString("key1", "value1");
    try store.putString("key2", "value2");

    const value1 = try store.getString("key1");
    const value2 = try store.getString("key2");

    std.debug.print("value1: {?s}\n", .{value1});
    std.debug.print("value2: {?s}\n", .{value2});

    capacity = store.*.objects.capacity();
    count = store.*.objects.count();
    std.debug.print("count: {d} cap: {d}\n", .{ count, capacity });
}
