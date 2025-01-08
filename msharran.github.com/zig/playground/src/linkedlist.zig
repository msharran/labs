const std = @import("std");
const stdout = std.io.getStdOut().writer();
const User = struct {
    name: []const u8,
    next: ?*User = null,

    pub fn init(name: []const u8) User {
        return User{ .name = name };
    }

    pub fn insertTop(self: *User, user: *User) void {
        user.next = self.next;
        self.next = user;
        std.debug.print("inserted  {s} {any} on top of {s} {any}\n", .{ user.name, user.name, self.name, self.name });
    }
};

test "user manual insertion" {
    var user = User.init("pedro");
    var user2 = User.init("juan");
    var user3 = User.init("pablo");

    std.debug.print("user : {s} {}\n", .{ user.name, user });
    user.insertTop(&user2);
    std.debug.print("user : {s} {}\n", .{ user.name, user });
    user.insertTop(&user3);
    std.debug.print("user : {s} {}\n", .{ user.name, user });

    std.debug.print("1: {s}\n", .{user.name});
    std.debug.print("2: {s}\n", .{user.next.?.name});
    std.debug.print("3: {s}\n", .{user.next.?.next.?.name});

    // print all users

    var u: ?*User = &user;

    while (u) |next| {
        try stdout.print("{s}\n", .{next.name});
        u = next.next;
    }
}

test "user insertion from a while loop" {
    const users = [_]User{ User.init("pedro"), User.init("juan"), User.init("pablo") };

    var head: ?*User = null;
    var prev: ?*User = null;

    var i: usize = 0;
    while (i < users.len) : (i += 1) {
        if (head == null) {
            head = &users[i];
        } else {
            prev.?.*.insertTop(&users[i]);
        }
        prev = &users[i];
    }

    var u: ?*User = head;

    while (u) |next| {
        try stdout.print("{s}\n", .{next.name});
        u = next.next;
    }
}
