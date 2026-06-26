---
title: eBPF Network Egress Monitor
type: source
created: 2026-06-26
updated: 2026-06-26
sources: []
tags: [coding-challenge, ebpf, go, sre, observability, containers, kubernetes, ecs, eks]
---

# eBPF Network Egress Monitor

Build a Go CLI that you can SSH into an ECS or EKS node and run with sudo to see which pod/task/container and process is making outbound network connections to which DNS name or IP and port.

The goal is a simple dependency view for incident response:

```text
WORKLOAD        CONTAINER  PROCESS  PID   DEST                  PROTO  CONNS  BYTES_SENT  BYTES_RECV
payments/api    app        node     991   api.stripe.com:443    TCP    2      8KB         24KB
checkout-task   app        curl     1234  example.com:443       TCP    3      12KB        81KB
worker/default  app        python   812   postgres:5432         TCP    9      44KB        512KB
```

Use best-effort DNS names when available, and fall back to IP:port when not.

## Scope

This challenge is for Linux worker nodes only:

- EKS worker nodes
- ECS on EC2/container instances

It is not aimed at Fargate, because the intended workflow is to SSH into the node and run a host-level eBPF tool there.

## What to build

Create a CLI named `egress-monitor` with a `top` command:

```bash
sudo egress-monitor top
```

It should show a live table grouped by:

- workload / namespace / task
- container
- process name
- PID
- destination DNS name or IP
- destination port
- protocol

Track:

- connection count
- bytes sent
- bytes received

Optional later additions like latency, failures, filters, and TUI can come after the MVP.

## Step 0 — Set up the Linux environment

You will need:

- Go
- `clang` / `llvm`
- Linux kernel headers
- `bpftool`
- `github.com/cilium/ebpf`
- root access or required capabilities to load eBPF programs

If you are on macOS, use a Linux VM or cloud instance. eBPF depends on Linux kernel features and cannot be developed directly against the macOS kernel.

A small test workload helps:

```bash
curl https://example.com
curl https://api.github.com
psql -h <some-host> -p 5432
```

## Step 1 — Observe outbound connects

Start by tracing outbound TCP `connect` calls.

Hook one of:

- `sys_enter_connect` / `sys_exit_connect`
- `tcp_v4_connect` / `tcp_v6_connect`
- a relevant tracepoint or kprobe path

For each connection attempt, emit:

- PID / TGID
- process name (`comm`)
- destination IP
- destination port
- address family
- result code

At this stage, raw events are enough:

```text
PID=1234 COMM=curl DEST=93.184.216.34:443 RESULT=0
```

## Step 2 — Build the Go event reader

Use `github.com/cilium/ebpf` to:

- load the compiled eBPF object
- attach programs
- read events from a ring buffer or perf event buffer
- decode C structs into Go structs

A basic event command is useful while developing:

```bash
sudo egress-monitor events
```

Example output:

```text
TIME                  PID    PROCESS   DEST              RESULT
2026-06-26T10:01:00   1234   curl      93.184.216.34:443  OK
```

## Step 3 — Add bytes sent and received

Track traffic per observed connection or per process/destination pair.

Possible hook points:

- `tcp_sendmsg`
- `tcp_recvmsg`
- socket-level tracepoints
- TCP close events to flush counters

The output should now include:

- protocol
- connection count
- bytes sent
- bytes received

Example:

```text
WORKLOAD        CONTAINER  PROCESS  PID   DEST                  PROTO  CONNS  BYTES_SENT  BYTES_RECV
checkout-task   app        curl     1234  example.com:443       TCP    3      12KB        81KB
```

## Step 4 — Map PID to container and workload

Make the tool container-aware for ECS and EKS.

Capture or derive:

- cgroup ID
- container ID
- pod namespace / pod name on Kubernetes
- ECS task/container metadata on ECS

Possible approaches:

- capture cgroup ID in eBPF and resolve it in Go
- read `/proc/<pid>/cgroup`
- use `/proc`, `/sys/fs/cgroup`, kubelet data, or container runtime metadata
- use ECS agent metadata on EC2-backed ECS nodes

The key is to answer: **which workload/container made the connection?**

## Step 5 — Correlate DNS names

DNS names matter more than raw IPs in incident response.

Add best-effort DNS correlation by observing DNS queries/responses on port 53.

Track:

- process / PID making the DNS query
- queried name
- returned IPs

Then prefer `example.com:443` over `93.184.216.34:443` when you can correlate it.

Fallback to IP:port when DNS is unavailable.

## Step 6 — Render the live table

Turn the event stream into a `top`-style view grouped by:

- workload / namespace / task
- container
- process
- destination
- protocol

The final table should answer:

- which container is talking out?
- which process is doing it?
- what DNS name or IP is it calling?
- how many connections?
- how much data sent and received?

## Step 7 — Test on ECS or EKS

Validate on a real node or a local equivalent.

For EKS-style testing:

- a small Kubernetes cluster
- a few test services
- one service calling another
- one service calling an external HTTPS endpoint
- one service calling a database endpoint

For ECS-style testing:

- an EC2-backed ECS cluster, or a Linux host with Docker/containerd
- multiple containers making outbound calls
- task/container metadata available to the host

## Extensions

If you want to go further:

- add filters
- add JSON output
- add Prometheus metrics
- add a TUI
- track connect latency
- track failed connects and resets
- support IPv6
- export dependency snapshots
- detect unexpected egress from a namespace or task

## What you should learn

By completing this challenge, you should learn:

- how Linux applications create outbound TCP connections
- how to attach eBPF programs from Go
- how to pass events from kernel space to user space
- how to use eBPF maps and ring buffers
- how to correlate kernel events with process metadata
- how cgroups connect Linux processes to containers
- how ECS and EKS nodes expose workload context
- how to turn low-level signals into a useful operational view
