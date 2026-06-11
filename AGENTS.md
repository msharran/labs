# AGENTS.md

Guidance for AI coding agents working in this repository.

## Repository overview

This is a polyglot labs/experiments repository. Each top-level language or platform directory contains independent experiments, challenges, or proofs of concept. There is no single repository-wide build or test command.

Top-level areas include:

- `c/` - C experiments, with a local `Makefile`.
- `go/` - many independent Go modules; most subdirectories with `go.mod` are standalone projects.
- `k8s/` - Kubernetes examples and experiments, including some Go modules and Makefile-based POCs.
- `kotlin/` - Gradle/Kotlin sample projects.
- `pulumi/` - Pulumi/IaC experiments.
- `python/` - Python experiments and scripts.
- `rust/` - independent Cargo crates.
- `terraform/` - Terraform experiments.
- `vm/` - VM-related automation, with a local `Makefile`.
- `zig/` - independent Zig projects.
- `scripts/` - standalone utility scripts.

## Working conventions

- Treat each experiment as independent. Before editing, identify the nearest project root by looking for files such as `go.mod`, `Cargo.toml`, `build.gradle`, `build.zig`, `Makefile`, or README files.
- Prefer project-local commands over repository-wide commands.
- Keep changes narrowly scoped to the requested experiment or directory.
- Do not reorganize top-level folders or normalize style across unrelated experiments unless explicitly asked.
- Do not overwrite generated, sample, or learning code unless the task specifically targets it.
- Preserve existing naming, formatting, and learning-oriented comments.

## Common commands by project type

Run these from the relevant subproject directory, not necessarily from the repository root.

### Go

```sh
go test ./...
go run .
go build ./...
gofmt -w <files>
```

If a Go subproject has a `Makefile`, prefer documented `make` targets.

### Rust

```sh
cargo test
cargo build
cargo fmt
cargo clippy --all-targets --all-features
```

### Kotlin / Gradle

```sh
./gradlew test
./gradlew build
```

Use the wrapper if present in the subproject; otherwise inspect the project before choosing a command.

### C / Makefile projects

```sh
make
make test
```

Inspect the local `Makefile` for supported targets before running destructive targets such as `clean`.

### Zig

```sh
zig build
zig build test
zig fmt <files>
```

### Terraform

```sh
terraform fmt
terraform validate
terraform plan
```

Do not apply infrastructure changes unless explicitly requested.

### Pulumi

```sh
pulumi preview
```

Do not run `pulumi up` unless explicitly requested.

## Safety notes

- This repository may contain experiments that contact local services, cloud APIs, or Kubernetes clusters. Inspect code and commands before running them.
- Avoid commands that mutate external infrastructure (`terraform apply`, `pulumi up`, `kubectl apply`, deploy scripts) unless the user explicitly asks.
- Do not commit or push unless asked.
- Leave unrelated untracked files alone.
