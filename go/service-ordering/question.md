# Service Startup Order

**Difficulty:** Medium  
**Tags:** Graph, Directed Graph, DAG, Topological Sort, BFS, Queue, Dependency Resolution

## Problem

You are given a set of services. Each service may depend on one or more other services. A service can only start after all services it depends on have already started.

Your task is to return a valid startup order where independent services come first, followed by services that depend on them.

This is similar to how tools like Docker Compose, Terraform, Pulumi, and build systems decide the order in which resources or services should be created.

## Input Format

For the coding interview version, assume the YAML is already parsed into this shape:

```go
map[string][]string
```

Where:

- the key is the service name
- the value is the list of services it depends on

Example YAML:

```yaml
services:
  api:
    depends_on: [db, cache]
  worker:
    depends_on: [db, queue]
  frontend:
    depends_on: [api]
  db:
    depends_on: []
  cache:
    depends_on: []
  queue:
    depends_on: []
  migrate:
    depends_on: [db]
```

Equivalent Go input:

```go
services := map[string][]string{
    "api":      {"db", "cache"},
    "worker":   {"db", "queue"},
    "frontend": {"api"},
    "db":       {},
    "cache":    {},
    "queue":    {},
    "migrate":  {"db"},
}
```

## Function Signature

```go
func OrderServices(services map[string][]string) ([][]string, error)
```

Return startup stages.

Each inner slice represents services that can start in parallel at that stage.

If the services cannot be ordered, return an error.

## Rules

1. If service `A` depends on service `B`, then `B` must appear before `A`.
2. Services with no dependencies should appear in the earliest possible stage.
3. Services in the same stage may be started in parallel.
4. Services inside each stage should be sorted alphabetically for deterministic output.
5. If a service depends on an undefined service, return an error.
6. If there is a circular dependency, return an error.

## Example 1

### Input

```go
services := map[string][]string{
    "api":      {"db", "cache"},
    "worker":   {"db", "queue"},
    "frontend": {"api"},
    "db":       {},
    "cache":    {},
    "queue":    {},
    "migrate":  {"db"},
}
```

### Output

```text
[
  ["cache", "db", "queue"],
  ["api", "migrate", "worker"],
  ["frontend"]
]
```

### Explanation

- `cache`, `db`, and `queue` have no dependencies, so they can start first.
- `api` can start after `db` and `cache`.
- `worker` can start after `db` and `queue`.
- `migrate` can start after `db`.
- `frontend` can start only after `api`.

A flattened startup order would be:

```text
cache, db, queue, api, migrate, worker, frontend
```

## Example 2

### Input

```go
services := map[string][]string{
    "web":   {"api"},
    "api":   {"db"},
    "db":    {},
    "logger": {},
}
```

### Output

```text
[
  ["db", "logger"],
  ["api"],
  ["web"]
]
```

### Explanation

`db` and `logger` are independent.  
`api` depends on `db`.  
`web` depends on `api`.

## Example 3: Missing Dependency

### Input

```go
services := map[string][]string{
    "api": {"db"},
}
```

### Output

```text
error
```

### Explanation

`api` depends on `db`, but `db` is not defined as a service.

## Example 4: Cycle

### Input

```go
services := map[string][]string{
    "api":    {"worker"},
    "worker": {"api"},
}
```

### Output

```text
error
```

### Explanation

`api` depends on `worker`, and `worker` depends on `api`. This creates a cycle, so no valid startup order exists.

## Constraints

Assume:

```text
1 <= number of services <= 10^4
0 <= total number of dependency edges <= 10^5
service names are non-empty strings
service names are unique
```

Expected complexity:

```text
Time:  O(V + E), ignoring alphabetical sorting cost
Space: O(V + E)
```

Where:

- `V` is the number of services
- `E` is the number of dependency relationships

## Learning Checkpoints

Try solving this in stages:

### Checkpoint 1: Model the graph

Convert this:

```text
api depends on db
```

Into this directed edge:

```text
db -> api
```

The edge points from dependency to dependent.

### Checkpoint 2: Track indegree

For each service, count how many dependencies must come before it.

Services with indegree `0` are ready to start.

### Checkpoint 3: Process ready services

Repeatedly:

1. take all currently ready services
2. add them as the next startup stage
3. remove their outgoing edges
4. reduce indegree of their dependents
5. newly indegree-`0` services become ready for the next stage

### Checkpoint 4: Detect invalid graphs

If not all services are processed, there must be a cycle.

### Checkpoint 5: Think about variants

After solving the base problem, try:

1. Return one flattened order instead of stages.
2. Return shutdown order.
3. Print the actual cycle path.
4. Allow missing external dependencies and treat them as already satisfied.
5. Add startup durations and compute minimum total startup time with unlimited parallelism.

## Clarifying Questions You Can Ask in an Interview

1. Should missing dependencies be errors or treated as external services?
2. If multiple services are ready at the same time, should output order be deterministic?
3. Do you want a flat order or parallel startup stages?
4. Should I parse YAML, or can I assume the input is already parsed?
5. Should the error include the cycle path or just indicate that a cycle exists?
