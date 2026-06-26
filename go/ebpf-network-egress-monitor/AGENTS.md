# AGENTS.md

Guidance for AI coding agents working in this repo.

## Repo scope

This repository is a single Go/eBPF experiment: build a CLI that passively observes outbound network egress on Linux and turns it into a service dependency map.

Primary reference:
- `CHALLENGE.md`

## Execution environment

- Development happens in an OrbStack Ubuntu VM.
- From this repo root, prefix commands with `orb` so they run inside the VM.
- Examples:
  - `orb go test ./...`
  - `orb go run .`
  - `orb make`
- For eBPF work, prefer the Linux VM over the macOS host.

## Implementation direction

- Use Go.
- Use the eBPF Go library from https://ebpf-go.dev/guides/getting-started (via `github.com/cilium/ebpf`).
- Keep changes narrowly scoped to this challenge.
- Preserve the learning-oriented nature of the repo and the challenge spec.

## Working conventions

- Inspect the repo before editing; do not assume a build system exists yet.
- Prefer small, local changes.
- Do not touch unrelated experiments or parent directories.
- If you add build/test commands, document them here or in `README.md`.
