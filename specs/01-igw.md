# Internet Gateway (IGW) — Learning Spec

> **Status:** Draft (collaborative)  
> **Series:** `01-igw` → `02-natgw` (planned)  
> **Last updated:** 2026-08-28

## 1. Why build this?

We already model AWS-style VPC networking in Terraform (`terraform/floci/modules/aws-network`) and diagrams, but those resources are *managed for us*. This lab series builds miniature, observable versions of the same primitives so we can answer questions like:

- What actually happens to a packet when a route table points `0.0.0.0/0` at an Internet Gateway?
- Where does NAT happen, and who owns the connection state?
- What is the gateway responsible for vs. what subnets / security policy still own?

**Learning outcomes for the IGW milestone:**

1. Trace outbound and inbound flows end-to-end with `tcpdump` and routing tables.
2. Implement 1:1 public↔private address translation for instances with assigned public IPs.
3. Separate **control plane** (attach IGW, install routes, allocate public IPs) from **data plane** (forward/NAT packets).
4. Produce a topology we can reuse as the attachment point for NAT Gateway in `02-natgw.md`.

---

## 2. Reference model: AWS Internet Gateway

This spec mirrors AWS behavior where practical, but intentionally simplifies HA, scale, and EC2 integration.

| AWS concept | Behavior we care about |
|---|---|
| Attachment | One IGW per VPC; IGW is associated with the VPC edge |
| Public subnet routing | Route table entry: `0.0.0.0/0 → igw` |
| Outbound | Instance sends with private source IP; IGW rewrites to assigned public IP |
| Inbound | Internet sends to public IP; IGW rewrites destination to private IP and delivers inside VPC |
| Security | IGW does **not** filter — NACLs / security groups (out of scope v1) |
| Stateful return | Return traffic for established flows must work without extra routing on instances |

**Explicit non-goals for v1 (AWS features we defer):**

- IPv6
- Horizontal scaling / HA pairs
- BGP or dynamic routing
- Integration with real cloud APIs
- Bandwidth accounting or rate limiting

---

## 3. Proposed lab topology

Simulate a small VPC on a single Linux host (VM or bare metal) using **network namespaces** and **veth pairs**. This keeps the lab cheap, scriptable, and easy to reset.

```text
                         [ upstream / lab "internet" ]
                                    |
                           +--------+--------+
                           |  ns: internet   |  203.0.113.0/24 (example)
                           |  203.0.113.1    |
                           +--------+--------+
                                    | veth
                           +--------+--------+
                           |  ns: igw        |  ← our Internet Gateway
                           |  vpc-side: 10.0.0.1/24
                           |  inet-side: 203.0.113.254/24
                           +--------+--------+
                                    | veth (vpc link)
                           +--------+--------+
                           |  ns: vpc        |  10.0.0.0/16 (simulated)
                           |                 |
                           |  +-----------+  |
                           |  | ns: pub-a |  |  10.0.1.0/24 public subnet
                           |  | instance |  |  10.0.1.10 (+ public 203.0.113.10)
                           |  +-----------+  |
                           +-----------------+
```

### Addressing plan (initial)

| Namespace | Role | CIDR / addresses |
|---|---|---|
| `internet` | Upstream / simulated internet | `203.0.113.0/24` |
| `igw` | Internet Gateway data plane | `10.0.0.1/24` (VPC side), `203.0.113.254/24` (internet side) |
| `vpc` | L2/L3 aggregation (optional hop) | `10.0.0.0/16` |
| `pub-a` | Public subnet A | `10.0.1.0/24` |
| `instance-a` | Workload with public IP mapping | private `10.0.1.10`, public `203.0.113.10` |

> **Open question:** Do we need a separate `vpc` namespace, or connect `pub-a` directly to `igw` via veth? Direct attachment is simpler; an extra hop better mimics "VPC router" semantics.

---

## 4. Responsibilities

### 4.1 Data plane (`igw` namespace)

The data plane forwards IP traffic between the VPC-facing interface and the internet-facing interface and performs **1:1 static NAT** for instances with allocated public addresses.

| Direction | Action |
|---|---|
| Egress (VPC → internet) | If source private IP has a public mapping, SNAT source to public IP; forward to internet |
| Ingress (internet → VPC) | If destination public IP has a mapping, DNAT to private IP; forward into VPC |
| No mapping | Drop (or optionally ICMP unreachable — decide in v1) |
| Non-NAT traffic | Not handled by IGW in this lab (no generic SNAT pool) |

**Implementation sketch (v1):**

- Enable `net.ipv4.ip_forward=1` in `igw` namespace
- Use `nftables` (preferred) or `iptables` for DNAT/SNAT rules driven by a mapping table
- Connection tracking (`nf_conntrack`) enabled so return traffic is un-NATed correctly

### 4.2 Control plane (user-space)

A small controller (language TBD) manages **desired state**:

```yaml
vpc:
  id: vpc-local
  cidr: 10.0.0.0/16

internet_gateway:
  id: igw-1
  vpc_id: vpc-local
  state: attached   # attached | detached

public_ip_allocations:
  - allocation_id: eipalloc-1
    public_ip: 203.0.113.10
    private_ip: 10.0.1.10
    instance: instance-a

route_tables:
  - id: rtb-public-a
    associations: [pub-a]
    routes:
      - dst: 10.0.0.0/16
        target: local
      - dst: 0.0.0.0/0
        target: igw-1
```

The controller:

1. Creates namespaces / veth / addresses (or shells out to setup scripts)
2. Programs routes in workload namespaces (`default via 10.0.0.1` in `pub-a`)
3. Renders NAT rules in `igw` from `public_ip_allocations`
4. Supports `attach` / `detach` IGW (add/remove default route in associated route tables)

> **Open question:** YAML file + CLI for v1, or jump straight to a minimal HTTP API?

---

## 5. Packet flows

### 5.1 Outbound (instance → internet)

Example: `instance-a` (`10.0.1.10`, public `203.0.113.10`) pings `203.0.113.1`.

```text
1. instance-a: src 10.0.1.10 → dst 203.0.113.1
2. route in pub-a: default via 10.0.0.1 (igw vpc-side)
3. igw: match SNAT 10.0.1.10 → 203.0.113.10
4. igw: forward to internet namespace
5. internet: sees src 203.0.113.10 → dst 203.0.113.1
```

### 5.2 Inbound (internet → instance)

Example: external host sends to `203.0.113.10:8080`.

```text
1. internet: dst 203.0.113.10
2. route to igw internet-side interface
3. igw: DNAT 203.0.113.10 → 10.0.1.10
4. igw: forward into pub-a / instance-a
5. instance-a: receives dst 10.0.1.10
```

### 5.3 Detached IGW

When `state: detached`:

- Remove `0.0.0.0/0 → igw` routes from associated route tables
- Instances retain private addresses but lose outbound internet path
- Inbound public IP traffic is dropped at igw (no route / no mapping applied)

---

## 6. Routing contract

### In public subnet namespaces (`pub-a`)

| Destination | Next hop / device |
|---|---|
| `10.0.0.0/16` | local / connected |
| `0.0.0.0/0` | `10.0.0.1` (IGW VPC-side) **only when IGW attached** |

### In `igw` namespace

| Destination | Next hop / device |
|---|---|
| `10.0.0.0/16` | VPC-facing interface |
| `0.0.0.0/0` | Internet-facing interface |
| Mapped public IPs | Resolved via NAT prerouting/postrouting |

### In `internet` namespace

| Destination | Next hop / device |
|---|---|
| `203.0.113.0/24` | local |
| `203.0.113.10` (mapped public IPs) | via IGW internet-side address |

---

## 7. Repository layout (proposed)

```text
specs/
  01-igw.md          # this document
  02-natgw.md        # next milestone

netlab/              # implementation root (name TBD)
  README.md
  topology/
    setup.sh         # create namespaces + veth
    teardown.sh
  igw/
    mappings.yaml    # desired state
    apply.sh         # render nftables / ip route
    controller/      # optional Go/Rust daemon (phase 2)
  tests/
    outbound_ping.sh
    inbound_curl.sh
```

> **Open question:** Top-level directory name — `netlab/`, `gateway-lab/`, or under `vm/`?

---

## 8. Acceptance criteria (v1)

- [ ] **Attach/detach:** Toggling IGW attachment adds/removes default route in `pub-a`.
- [ ] **Outbound NAT:** From `instance-a`, traffic to `internet` appears with source `203.0.113.10`.
- [ ] **Inbound NAT:** A listener on `instance-a:8080` is reachable at `203.0.113.10:8080` from `internet`.
- [ ] **No public IP, no internet:** An instance without mapping cannot reach `internet` through the IGW.
- [ ] **Observability:** Documented `tcpdump` attachment points for each hop in §5.
- [ ] **Reset:** `teardown.sh` returns host to clean state.

---

## 9. Implementation phases

| Phase | Goal | Deliverables |
|---|---|---|
| **0 — Manual plumbing** | Prove topology and routing | `setup.sh`, ping across namespaces without NAT |
| **1 — Static 1:1 NAT** | Data plane only | `nftables` rules, hard-coded mapping, acceptance tests |
| **2 — Control plane** | Declarative mappings + routes | YAML + `apply`, attach/detach |
| **3 — Hardening** | Operability | idempotent setup, better errors, packet capture helpers |

NAT Gateway (`02-natgw.md`) should **reuse** this topology: private subnets appear, default route targets NAT GW in a public subnet, and public subnet still uses IGW for `0.0.0.0/0`.

---

## 10. Decisions to make together

Please react to these and we'll fold answers back into the spec.

1. **Host environment:** Linux VM only (OrbStack), or also support containerized setup?
2. **Topology shape:** `pub-a` ↔ `igw` direct veth, or include an explicit `vpc` router namespace?
3. **Control plane language:** Bash + `nft` first, or Go/Rust from the start?
4. **State interface:** YAML file, CLI subcommands (`igw attach`, `igw allocate-eip`), or HTTP API?
5. **Project location:** New top-level `netlab/` vs. extend `vm/` automation?
6. **Public IP pool:** Single `/24` lab network (`203.0.113.0/24` TEST-NET-3) — OK?
7. **Drop vs. reject:** For unmapped traffic at IGW, should we drop silently or return ICMP unreachable?

---

## 11. Glossary

| Term | Meaning in this lab |
|---|---|
| **VPC** | Isolated L3 network built from namespaces |
| **Public subnet** | Subnet whose route table sends `0.0.0.0/0` to the IGW |
| **IGW** | Edge gateway performing 1:1 NAT for allocated public IPs |
| **Public IP allocation** | Binding `{public_ip → private_ip}` managed by control plane |
| **Upstream / internet** | Simulated external network namespace |

---

## 12. References

- [AWS: Internet gateways](https://docs.aws.amazon.com/vpc/latest/userguide/VPC_Internet_Gateway.html)
- Existing repo baseline: `terraform/floci/modules/aws-network/main.tf`
- Diagram: `terraform/floci/modules/aws-network/examples/diagram/network_architecture.py`
- Linux: network namespaces, veth, nftables NAT
