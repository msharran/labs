# Internet Gateway (IGW) — Learning Spec

> **Status:** Draft (collaborative)  
> **Series:** `01-igw` → `02-natgw` (planned)  
> **Last updated:** 2026-08-28

## 1. Goal

Understand **what an Internet Gateway actually does** — not build a VPC.

We already use AWS IGWs in Terraform (`terraform/floci/modules/aws-network`), but they are opaque managed resources. This lab isolates the gateway itself: a box at the network edge that forwards traffic between an internal host and the internet, and translates addresses when a host has a public IP.

**What we want to be able to explain after this lab:**

1. Why a host needs a default route pointing at the IGW to reach the internet.
2. What changes in a packet as it crosses the IGW (source/destination rewrite).
3. Why return traffic works without extra rules on the host (connection tracking).
4. What "attach" and "detach" mean in practice (routes on / routes off).

**What we are not building:**

- A VPC, subnets, route tables, or control-plane API
- HA, IPv6, security groups, or cloud integration

---

## 2. What an IGW does (the mental model)

An Internet Gateway is the **edge router between your internal network and the public internet**.

```text
   [ host ]          [ IGW ]          [ internet ]
  10.0.1.10  ←──→  translates  ←──→  203.0.113.10
  (private)        + forwards        (public)
```

| Direction | IGW action |
|---|---|
| **Outbound** | Host sends packet with private source IP → IGW rewrites source to the host's public IP → forwards to internet |
| **Inbound** | Internet sends to public IP → IGW rewrites destination to private IP → forwards to host |
| **No public IP** | IGW does not NAT this host; it cannot reach the internet through the gateway |

The IGW does **not** filter traffic. It routes and translates. Security policy is someone else's job.

In AWS, you attach an IGW to a VPC and add a route `0.0.0.0/0 → igw` in a route table. **Attach/detach is really about whether that route exists.** Detached = no path to the internet, even if the gateway process still exists.

---

## 3. Minimal lab topology

Three network namespaces on one Linux host. No VPC abstraction — just enough plumbing to see the IGW in action.

```text
                    [ ns: internet ]
                    203.0.113.1/24
                           |
                      veth pair
                           |
                    [ ns: igw ]  ← the thing we're learning
              inet: 203.0.113.254/24
              host: 10.0.1.1/24
                           |
                      veth pair
                           |
                    [ ns: host ]
                    10.0.1.10/24
                    public mapping: 203.0.113.10
```

| Namespace | Role | Address |
|---|---|---|
| `host` | Internal machine with a public IP allocation | `10.0.1.10`, mapped to `203.0.113.10` |
| `igw` | Internet Gateway — forwards + 1:1 NAT | `10.0.1.1` (host side), `203.0.113.254` (internet side) |
| `internet` | Simulated upstream | `203.0.113.1` |

That's it. One host, one gateway, one internet peer.

---

## 4. The three mechanisms to learn

### 4.1 Routing

The host needs a default route through the IGW:

```text
host:  0.0.0.0/0 via 10.0.1.1
igw:   10.0.1.0/24 dev host-side
       0.0.0.0/0    dev internet-side
```

**Exercise:** ping `203.0.113.1` from `host` with forwarding enabled but **no NAT**. The ping fails (or replies go to the wrong place). Observe why routing alone is not enough.

### 4.2 1:1 NAT

The IGW maps one public IP to one private IP:

```text
Outbound:  src 10.0.1.10  →  src 203.0.113.10
Inbound:   dst 203.0.113.10  →  dst 10.0.1.10
```

Implemented with `nftables` (or `iptables`) in the `igw` namespace, plus `net.ipv4.ip_forward=1`.

**Exercise:** add the NAT rule, ping again, `tcpdump` on all three namespaces and compare addresses at each hop.

### 4.3 Connection tracking

Return traffic for an outbound flow must be un-NATed automatically. Linux `nf_conntrack` handles this — the IGW does not need a per-flow rule for replies.

**Exercise:** start a TCP listener on `host:8080`, connect from `internet` to `203.0.113.10:8080`, watch the DNAT on ingress and the reverse on the reply.

---

## 5. Packet walks

### Outbound ping: `host` → `203.0.113.1`

```text
1. host:      src 10.0.1.10      dst 203.0.113.1
2. host routes via 10.0.1.1 (igw)
3. igw SNAT:  src 203.0.113.10   dst 203.0.113.1
4. internet:  sees public source, replies to 203.0.113.10
5. igw unsNAT reply back to 10.0.1.10
```

### Inbound TCP: `internet` → `host:8080`

```text
1. internet:  dst 203.0.113.10:8080
2. igw DNAT:  dst 10.0.1.10:8080
3. host:      receives connection on :8080
```

### Detached IGW

Remove the default route on `host`. The IGW namespace and NAT rules can still exist, but the host has no path to it. **This is what AWS "detach" means at the packet level.**

---

## 6. Lab scripts (proposed)

```text
netlab/igw/
  README.md           # concepts + tcpdump cheat sheet
  setup.sh            # create 3 namespaces, veth, addresses, routes
  nat.sh              # apply/remove the 1:1 NAT rule
  attach.sh           # add default route on host
  detach.sh           # remove default route on host
  teardown.sh         # delete namespaces
  tests/
    ping_outbound.sh
    curl_inbound.sh
```

No YAML control plane. No route table abstraction. Shell scripts you can read top to bottom.

---

## 7. Acceptance criteria

- [ ] Can draw the three-namespace topology from memory.
- [ ] Can explain why ping fails without NAT but works with it.
- [ ] Outbound traffic from `host` appears as `203.0.113.10` in `internet`.
- [ ] Inbound connection to `203.0.113.10` reaches `host`.
- [ ] Detaching (removing default route) blocks outbound internet access.
- [ ] Have `tcpdump` captures for at least one outbound and one inbound flow.

---

## 8. Implementation phases

| Phase | Focus | Done when |
|---|---|---|
| **0 — Plumbing** | Namespaces, veth, routes, forwarding | `host` can reach `igw` interface, not yet internet |
| **1 — NAT** | Single `nftables` 1:1 rule | Outbound ping to `internet` works |
| **2 — Inbound** | DNAT for incoming connections | `curl` from `internet` to `host` works |
| **3 — Attach/detach** | Route toggle scripts | Can demonstrate blocked vs. allowed path |

NAT Gateway (`specs/natgw/02-natgw.md`) comes later and adds a **second** gateway role (SNAT for hosts without public IPs). The IGW lab should be done first so that difference is obvious.

---

## 9. Open questions

1. **Host environment:** Linux VM (OrbStack) only?
2. **Project location:** `netlab/igw/` at repo root?
3. **nftables vs iptables:** `nft` preferred on modern Linux — OK?
4. **Address pool:** TEST-NET-3 `203.0.113.0/24` for the simulated internet?

---

## 10. References

- [AWS: Internet gateways](https://docs.aws.amazon.com/vpc/latest/userguide/VPC_Internet_Gateway.html)
- Repo baseline: `terraform/floci/modules/aws-network/main.tf`
- Linux: `ip netns`, veth pairs, `nftables` NAT, `nf_conntrack`
