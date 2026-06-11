# Handoff Summary — dbt Labs Interview Prep

## User Request

The user asked:

> See my recent emails from dbt labs. For the wed and fri interviews, what all should i prepare?  
> Give me proper challenges that i can do, and research about the coding platform they use in interview, it says ai assistance available. What level of assistance is it? Agentic or autocomplete?

## Context About User

- User: Sharran M
- Role/background: Staff Engineer / Platform Engineering / SRE
- Strong areas: Go, Kubernetes/EKS, AWS, Pulumi, Docker, Linux, CI/CD, platform engineering
- Current interview: dbt Labs, Sr Software Engineer II — Infrastructure
- User prefers practical, high-signal preparation, not vague “likely” answers.

## Emails Checked

Gmail was searched for recent dbt Labs / ModernLoop / Greenhouse / interview-related emails.

### Key email found

**Subject:** `dbt Labs | Coding Challenge Interview Confirmed for Sr Software Engineer II (Infrastructure)`  
**From:** Ankita Agarwal `<ankita.agarwal@dbtlabs.com>`  
**To:** Sharran M  
**Date:** June 8, 2026

### Confirmed interviews

#### Wednesday Interview

- **Date:** Wednesday, June 10, 2026
- **Time:** 7:30 PM – 8:30 PM IST
- **Type:** Software Engineering - Coding Interview
- **Interviewer:** Dashiel Lopez Mendez, Senior Software Engineer II
- **Platform:** CoderPad
- **CoderPad link:** Present in email

#### Friday Interview

- **Date:** Friday, June 12, 2026
- **Time:** 7:00 PM – 8:00 PM IST
- **Type:** Terraform Coding Challenge in AWS or Azure
- **Interviewer:** Norm Dole, Senior Infrastructure Engineer II
- **Platform:** CoderPad
- **CoderPad link:** Present in email

### Other email found

**BrightHire email**

- dbt Labs uses BrightHire to record interviews.
- The email says the user can ask the interviewer to stop recording at any time.
- Purpose stated: help interviewers focus on the conversation and ensure a fairer, more consistent hiring process.

## Important Exact Interview Instructions From Email

The email says:

- This is a **60-minute live coding challenge**.
- They will use **CoderPad**.
- They encourage collaboration and communication.
- They want the candidate to walk interviewers through the approach.
- They encourage asking for feedback and questions along the way.
- If stuck, candidate should ask the interviewer for help.
- It is an interactive format.
- The email explicitly says:

> You may use the AI function in CoderPad to help you throughout the interview. It is encouraged but not required.

### Language options in the email

The email asks which language the user prefers:

- Python
- Java
- Go
- Rust
- Kotlin
- JavaScript / TypeScript

Recommendation given: **use Go** unless the user is much faster in Python.

## Research Summary — CoderPad AI Assistance

Public CoderPad documentation was checked.

### CoderPad Interview

CoderPad Interview is a browser-based IDE where candidates can:

- write code
- run code
- view output
- collaborate live with interviewer
- use many languages, including Go, Python, Rust, TypeScript, Bash, PostgreSQL, and Terraform

### CoderPad AI Assist

CoderPad’s AI Assist supports an AI Assist tab with model options such as:

- GPT
- Claude
- Gemini
- Llama

It has at least two modes:

#### Ask mode

- Chat-style AI assistance.
- AI answers prompts.
- Does **not** directly edit the candidate’s code.

#### Edit mode

- AI can change code and file structure, if enabled.
- Depends on the company/interview setup.

### Agentic or autocomplete?

Answer given:

- It is **not just autocomplete**.
- It is primarily **chat-style AI assistance inside CoderPad**.
- It may include edit mode depending on dbt Labs’ configuration.
- Fully agentic assistance like Claude Code / Codex may exist only if dbt’s CoderPad Enterprise setup enables it.
- CoderPad docs mention Claude Code and Codex availability on Enterprise plans and in multi-file supported environments.

### Important warning

CoderPad saves AI prompts and AI outputs for review/playback. So the user should not secretly ask AI to solve the whole interview problem.

Recommended positioning:

> “I saw AI Assist is enabled. I’m comfortable solving directly. I may use it for syntax checks or test cases, and I’ll keep you in the loop before I do.”

Good AI uses:

```text
Check this Terraform for syntax issues only. Do not redesign it.
```

## Updated Prep Plan — Two Practice Problems

The user is unsure how much AI assistance will actually be available during the dbt Labs CoderPad interview. They want to prepare with two separate practice problems:

1. **Without AI assistance** — solve a graph/data-structure problem fully by hand.
2. **With AI assistance** — decide later; likely use AI for syntax checks, test cases, or code review rather than full solution generation.

This handoff currently defines only Problem 1.

## Problem 1 — Without AI Assistance: Service Dependency Ordering

### Online basis / why this is valid

This is a practical topological-sort problem, similar to:

- Docker Compose startup ordering: Docker docs state that Compose creates services in dependency order using `depends_on`, and `db`/`redis` are created before a dependent `web` service.
  - Reference: https://docs.docker.com/compose/how-tos/startup-order/
- Course Schedule II / prerequisite ordering: a canonical interview problem where prerequisites are modeled as directed edges and solved with topological sort.
  - Reference: https://www.geeksforgeeks.org/dsa/find-course-schedule-ii/
  - Also canonical LeetCode problem: `Course Schedule II`.

The custom practice problem below keeps the real infrastructure flavor: order services/resources so dependencies are created first, like Terraform/Pulumi/Docker Compose ordering engines.

### Question

You are given a YAML service configuration. Each service may declare a `depends_on` list containing other services that must be started before it.

Write a program that reads this configuration and prints a valid startup order such that:

1. Services with no dependencies appear before services that depend on them.
2. If service `api` depends on `db`, then `db` must appear before `api`.
3. If multiple services are available to start at the same time, output them in alphabetical order for deterministic results.
4. If the configuration contains a dependency cycle, report an error instead of printing an order.
5. If a service depends on a service name that is not defined in the YAML, report an error.

For learning focus, solve this as a graph problem:

- Each service is a graph node.
- Each dependency relationship is a directed edge from dependency to dependent.
  - Example: `api depends_on db` becomes `db -> api`.
- Services with zero incoming edges are independent startup candidates.
- Use topological sorting to produce the order.

### Input example

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

### Expected output

If printing as startup waves/stages, where services in the same stage can start in parallel:

```text
stage 0: cache, db, queue
stage 1: api, migrate, worker
stage 2: frontend
```

If printing as one flattened order, one valid deterministic output is:

```text
cache, db, queue, api, migrate, worker, frontend
```

### Invalid input example — cycle

```yaml
services:
  api:
    depends_on: [worker]
  worker:
    depends_on: [api]
```

Expected behavior:

```text
error: dependency cycle detected
```

### Invalid input example — missing dependency

```yaml
services:
  api:
    depends_on: [db]
```

Expected behavior:

```text
error: service api depends on undefined service db
```

### Suggested function shape

Use any interview language, but Go is preferred for the user.

```go
func OrderServices(services map[string][]string) (stages [][]string, order []string, err error)
```

Where `services` maps a service name to the list of services it depends on.

Example parsed input:

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

### Learning goals

This problem should help the user practice:

- Building a directed graph from real-world dependency data.
- Choosing edge direction correctly: dependency -> dependent.
- Maintaining an adjacency list.
- Maintaining indegree counts.
- Finding independent nodes with indegree `0`.
- Performing topological sort with Kahn's algorithm/BFS.
- Producing startup stages/levels, not just a single list.
- Detecting cycles when processed node count is less than total node count.
- Understanding why topological sort only works for DAGs: directed acyclic graphs.
- Optionally implementing a DFS-based topological sort and comparing it with Kahn's algorithm.

### Stretch variants

After solving the basic version, try these variants:

1. **Destroy order:** print shutdown/delete order, which is usually the reverse topological order.
2. **Parallel batches:** return only stages, since all services in a stage can run concurrently.
3. **Critical path:** if each service has a startup duration, compute the minimum total startup time with unlimited parallelism.
4. **Cycle path:** instead of only saying a cycle exists, print one cycle path.
5. **External dependencies:** allow dependencies not present in the YAML and treat them as already satisfied.
6. **Stable input order:** preserve YAML declaration order instead of alphabetical order.

### Notes for interview practice

Do not start by coding immediately. First explain:

1. This is a directed graph problem.
2. The graph must be a DAG to have a valid startup order.
3. A topological sort gives the required order.
4. Kahn's algorithm is a good fit because it naturally starts from independent services and can produce startup stages.
5. Complexity is `O(V + E)`, where `V` is number of services and `E` is number of dependency edges.
