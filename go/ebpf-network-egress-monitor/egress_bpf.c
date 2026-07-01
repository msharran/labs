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

#define EVT_CONNECT 1
#define EVT_DNS_QUERY 2
#define EVT_DNS_RESPONSE 3
#define DNS_PORT 53
#define DNS_PAYLOAD_LEN 512

struct sockaddr_header {
    unsigned short sa_family;
    char sa_data[14];
};

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

struct sys_enter_sendto_ctx {
    unsigned short common_type;
    unsigned char common_flags;
    unsigned char common_preempt_count;
    int common_pid;
    int syscall_nr;
    long fd;
    long buff;
    long len;
    long flags;
    long addr;
    long addr_len;
};

struct sys_enter_recvfrom_ctx {
    unsigned short common_type;
    unsigned char common_flags;
    unsigned char common_preempt_count;
    int common_pid;
    int syscall_nr;
    long fd;
    long ubuf;
    long size;
    long flags;
    long addr;
    long addr_len;
};

struct sys_exit_ctx {
    unsigned short common_type;
    unsigned char common_flags;
    unsigned char common_preempt_count;
    int common_pid;
    int syscall_nr;
    long ret;
};

struct event {
    __u32 type;
    __u32 pid;
    __u32 tgid;
    __u16 family;
    __u16 dport;
    __u16 data_len;
    __u8 ip_version;
    __u8 pad1;
    char comm[TASK_COMM_LEN];
    __u8 dst[16];
    __u8 data[DNS_PAYLOAD_LEN];
};

struct recv_args {
    __u64 buf;
    __u64 size;
};

struct {
    __uint(type, BPF_MAP_TYPE_RINGBUF);
    __uint(max_entries, 1 << 24);
} events SEC(".maps");

struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 8192);
    __type(key, __u64);
    __type(value, struct recv_args);
} recvfrom_args SEC(".maps");

static __always_inline void fill_task(struct event *e, __u32 type) {
    __u64 pid_tgid = bpf_get_current_pid_tgid();
    e->type = type;
    e->pid = pid_tgid >> 32;
    e->tgid = (__u32)pid_tgid;
    bpf_get_current_comm(&e->comm, sizeof(e->comm));
}

static __always_inline void copy_dns_payload(struct event *e, void *buf, long len) {
    __u16 n = len;
    if (n > DNS_PAYLOAD_LEN) {
        n = DNS_PAYLOAD_LEN;
    }
    e->data_len = n;
    bpf_probe_read_user(e->data, DNS_PAYLOAD_LEN, buf);
}

SEC("tracepoint/syscalls/sys_enter_connect")
int handle_enter(struct sys_enter_connect_ctx *ctx) {
    struct event *e;
    struct sockaddr_header sa = {};

    if (ctx->uservaddr == 0) {
        return 0;
    }
    if (bpf_probe_read_user(&sa, sizeof(sa), (void *)ctx->uservaddr) < 0) {
        return 0;
    }
    if (sa.sa_family != AF_INET && sa.sa_family != AF_INET6) {
        return 0;
    }

    e = bpf_ringbuf_reserve(&events, sizeof(*e), 0);
    if (!e) {
        return 0;
    }
    __builtin_memset(e, 0, sizeof(*e));
    fill_task(e, EVT_CONNECT);
    e->family = sa.sa_family;

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
    } else {
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

SEC("tracepoint/syscalls/sys_enter_sendto")
int handle_sendto(struct sys_enter_sendto_ctx *ctx) {
    struct sockaddr_header sa = {};
    __u16 dport = 0;
    struct event *e;

    if (ctx->buff == 0 || ctx->len < 12 || ctx->addr == 0) {
        return 0;
    }
    if (bpf_probe_read_user(&sa, sizeof(sa), (void *)ctx->addr) < 0) {
        return 0;
    }
    if (sa.sa_family == AF_INET) {
        struct sockaddr_in in4 = {};
        if (bpf_probe_read_user(&in4, sizeof(in4), (void *)ctx->addr) < 0) {
            return 0;
        }
        dport = bpf_ntohs(in4.sin_port);
    } else if (sa.sa_family == AF_INET6) {
        struct sockaddr_in6 in6 = {};
        if (bpf_probe_read_user(&in6, sizeof(in6), (void *)ctx->addr) < 0) {
            return 0;
        }
        dport = bpf_ntohs(in6.sin6_port);
    } else {
        return 0;
    }
    if (dport != DNS_PORT) {
        return 0;
    }

    e = bpf_ringbuf_reserve(&events, sizeof(*e), 0);
    if (!e) {
        return 0;
    }
    __builtin_memset(e, 0, sizeof(*e));
    fill_task(e, EVT_DNS_QUERY);
    e->family = sa.sa_family;
    e->dport = DNS_PORT;
    copy_dns_payload(e, (void *)ctx->buff, ctx->len);
    bpf_ringbuf_submit(e, 0);
    return 0;
}

SEC("tracepoint/syscalls/sys_enter_recvfrom")
int handle_recvfrom_enter(struct sys_enter_recvfrom_ctx *ctx) {
    __u64 pid_tgid = bpf_get_current_pid_tgid();
    struct recv_args args = {};
    if (ctx->ubuf == 0 || ctx->size < 12) {
        return 0;
    }
    args.buf = (__u64)ctx->ubuf;
    args.size = ctx->size;
    bpf_map_update_elem(&recvfrom_args, &pid_tgid, &args, BPF_ANY);
    return 0;
}

SEC("tracepoint/syscalls/sys_exit_recvfrom")
int handle_recvfrom_exit(struct sys_exit_ctx *ctx) {
    __u64 pid_tgid = bpf_get_current_pid_tgid();
    struct recv_args *args;
    struct event *e;
    long len = ctx->ret;

    args = bpf_map_lookup_elem(&recvfrom_args, &pid_tgid);
    if (!args) {
        return 0;
    }
    bpf_map_delete_elem(&recvfrom_args, &pid_tgid);
    if (len < 12) {
        return 0;
    }

    e = bpf_ringbuf_reserve(&events, sizeof(*e), 0);
    if (!e) {
        return 0;
    }
    __builtin_memset(e, 0, sizeof(*e));
    fill_task(e, EVT_DNS_RESPONSE);
    copy_dns_payload(e, (void *)args->buf, len);
    bpf_ringbuf_submit(e, 0);
    return 0;
}

char LICENSE[] SEC("license") = "GPL";
