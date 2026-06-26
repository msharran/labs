# eBPF Network Egress Monitor

## TL;DR

A Go/eBPF CLI for Linux that watches outbound network connections and shows which processes or containers are talking to which remote services.

See `CHALLENGE.md` for the full project brief.

## Run

```sh
make run-trace
```

The Makefile builds inside the OrbStack Linux VM and runs the monitor with `sudo`, which is required for loading and attaching eBPF programs.
