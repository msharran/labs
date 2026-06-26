//go:build ignore

#include <linux/bpf.h>
#include <linux/in.h>
#include <linux/in6.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_endian.h>

#ifndef TASK_COMM_LEN
#define TASK_COMM_LEN 16
#endif

#ifndef AF_INET
#define AF_INET 2
#endif

#ifndef AF_INET6
#define AF_INET6 10
#endif

struct sockaddr_header {
    unsigned short sa_family;
    char sa_data[14];
};

// This struct matches the fields shown by:
//   sudo cat /sys/kernel/debug/tracing/events/syscalls/sys_enter_connect/format
struct sys_enter_connect_ctx {
    unsigned short common_type;
    unsigned char common_flags;
    unsigned char common_preempt_count;
    int common_pid;

    int syscall_nr;
    long fd;
    long uservaddr;
    long addrlen;
};

struct event {
    __u32 pid;
    __u32 tgid;
    __u16 family;
    __u16 dport;
    __u8 ip_version;
    __u8 pad1;
    char comm[TASK_COMM_LEN];
    __u8 dst[16];
    __u8 pad2[2];
};

// Ring buffer map: kernel writes events, Go reads them.
struct {
    __uint(type, BPF_MAP_TYPE_RINGBUF);
    __uint(max_entries, 1 << 24);
} events SEC(".maps");

SEC("tracepoint/syscalls/sys_enter_connect")
int handle_enter(struct sys_enter_connect_ctx *ctx) {
    __u64 pid_tgid = bpf_get_current_pid_tgid();
    struct event *e;
    struct sockaddr_header sa = {};

    if (ctx->uservaddr == 0) {
        return 0;
    }

    // First read just enough to learn the address family.
    if (bpf_probe_read_user(&sa, sizeof(sa), (void *)ctx->uservaddr) < 0) {
        return 0;
    }

    // Only emit events for IPv4 / IPv6 connects.
    if (sa.sa_family != AF_INET && sa.sa_family != AF_INET6) {
        return 0;
    }

    e = bpf_ringbuf_reserve(&events, sizeof(*e), 0);
    if (!e) {
        return 0;
    }

    __builtin_memset(e, 0, sizeof(*e));

    e->pid = pid_tgid >> 32;
    e->tgid = (__u32)pid_tgid;
    e->family = sa.sa_family;
    bpf_get_current_comm(&e->comm, sizeof(e->comm));

    if (sa.sa_family == AF_INET) {
        struct sockaddr_in in4 = {};
        __u32 addr;

        if (bpf_probe_read_user(&in4, sizeof(in4), (void *)ctx->uservaddr) < 0) {
            bpf_ringbuf_discard(e, 0);
            return 0;
        }

        e->ip_version = 4;
        e->dport = bpf_ntohs(in4.sin_port);

        addr = bpf_ntohl(in4.sin_addr.s_addr);
        e->dst[0] = (__u8)(addr >> 24);
        e->dst[1] = (__u8)(addr >> 16);
        e->dst[2] = (__u8)(addr >> 8);
        e->dst[3] = (__u8)addr;
    }

    if (sa.sa_family == AF_INET6) {
        struct sockaddr_in6 in6 = {};

        if (bpf_probe_read_user(&in6, sizeof(in6), (void *)ctx->uservaddr) < 0) {
            bpf_ringbuf_discard(e, 0);
            return 0;
        }

        e->ip_version = 6;
        e->dport = bpf_ntohs(in6.sin6_port);
        __builtin_memcpy(e->dst, &in6.sin6_addr, 16);
    }

    bpf_ringbuf_submit(e, 0);
    return 0;
}

char LICENSE[] SEC("license") = "GPL";
