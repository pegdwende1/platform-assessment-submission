# Mid-Project Change: Acquisition Integration Strategy

## Context

Six months into the platform initiative, Acme acquires a company that introduces:

- ~10 additional engineers
- A separate source control platform
- Different CI/CD tooling
- A different cloud provider
- Different monitoring and observability platforms
- Several customer-facing applications that must remain operational

This document describes how the original platform strategy adapts to this change.

---

## Guiding Principle

**Do not break what already works.** The acquired team has shipped production software successfully — their tooling has value even if it differs from ours. The integration goal is not to "fix" their setup but to find the right convergence point that serves both teams without disrupting either.

---

## What Remains Unchanged

These recommendations from the original proposal hold regardless of the acquisition:

| Recommendation | Why It Still Applies |
|----------------|---------------------|
| **GitOps as the deployment model** | Declarative, auditable deployments work regardless of cloud provider or CI tooling |
| **Progressive delivery (canary rollouts)** | Risk reduction during deployment is universal — doesn't matter what's inside the container |
| **Image signing and supply chain security** | Cosign/Sigstore is cloud-agnostic; the verification model doesn't change |
| **Policy-as-code (Kyverno or equivalent)** | Security standards apply to all workloads regardless of origin |
| **DORA metrics** | Measuring deployment health applies to both teams; gives objective data for integration decisions |
| **QA self-service model** | QA empowerment is an organizational capability, not a tooling one |
| **Cadence tiers** | The acquired team's services get their own tier classification based on risk profile |
| **Rollback strategy (instant, ArgoCD, git revert)** | Mechanism is the same even if the underlying infrastructure differs |

**The principles are platform-agnostic. The implementations may differ across clouds, but the operating model remains consistent.**

---

## What Shifts in Priority

| Original Priority | New Priority | Reason |
|-------------------|--------------|--------|
| Single-cluster ArgoCD deployment | **Multi-cluster, multi-cloud ArgoCD** | Must manage workloads on two cloud providers |
| Org-wide Helm standardization | **Abstraction layer that supports multiple packaging formats** | Acquired team may use Kustomize, Terraform, or their own approach |
| GHCR as single registry | **Multi-registry strategy** (GHCR + acquired team's registry) | Can't force immediate registry migration for running apps |
| Single observability stack | **Observability federation** — unified view across platforms | Need cross-platform visibility before consolidation |
| Service catalog (nice-to-have) | **Service catalog becomes critical** — must inventory what exists across both orgs | You can't integrate what you haven't mapped |
| Aggressive onboarding to new platform | **Slower, trust-building onboarding** — prove value before asking teams to change | Forcing tools on an acquired team destroys morale and productivity |

---

## Immediate Standardization vs. Phased Approach

**I would pursue a phased approach. Immediate standardization would be a mistake.**

### Why Not Standardize Immediately

1. **Customer-facing applications must keep running.** Migrations introduce risk. Forcing an infrastructure change on apps serving live customers — with a team that doesn't know the new platform yet — is asking for incidents.

2. **10 engineers just joined a new company.** They're already dealing with new colleagues, new processes, new culture. Telling them "and also rewrite your deployment pipelines by next month" is a recipe for attrition.

3. **We don't know what they have yet.** Their tooling might be better than ours in specific areas. Standardizing before evaluating means potentially discarding good solutions.

4. **The original team is only 6 months into their own migration.** The platform isn't mature enough to absorb 10 more engineers and their workloads simultaneously.

### The Phased Approach

| Phase | Timeline | Focus |
|-------|----------|-------|
| **Phase 1: Observe & Map** | Months 1-2 | Inventory their systems, understand their workflows, build relationships |
| **Phase 2: Bridge** | Months 3-4 | Connect observability, establish shared standards for new work only, cross-team visibility |
| **Phase 3: Converge** | Months 5-8 | Migrate services one-by-one to unified platform, acquired team drives their own migration |
| **Phase 4: Consolidate** | Months 9-12 | Retire legacy tooling, single operational model (with multi-cloud support) |

---

## Evaluating Existing Tooling: Framework

Before deciding to replace or retain anything, I'd evaluate the acquired team's tooling against these criteria:

### Evaluation Matrix

| Criterion | Questions to Answer | Weight |
|-----------|---------------------|--------|
| **Operational health** | Is it working? What's their deployment failure rate? MTTR? | High |
| **Team expertise** | How well does the team know it? Is it maintained by one person or shared knowledge? | High |
| **Scalability** | Can it support org growth (20+ engineers, 30+ services)? | Medium |
| **Security posture** | Does it meet our compliance requirements? Supply chain controls? | High |
| **Integration cost** | How hard is it to connect to our platform (observability, catalog, policies)? | Medium |
| **Maintenance burden** | Is it self-hosted? Who operates it? What's the operational cost? | Medium |
| **Vendor/community support** | Is it actively maintained? End-of-life risk? | Low |
| **Overlap with our platform** | Does it solve the same problem we already solved? Better or worse? | Medium |

### Decision Framework

```
┌──────────────────────────────────────────────────────────┐
│              Is it working well for them?                  │
└────────────────────────┬─────────────────────────────────┘
                         │
                    ┌────┴────┐
                    │  YES    │                    ┌─── NO ──→ Replace (but not urgently)
                    └────┬────┘
                         │
         ┌───────────────┴───────────────────────┐
         │  Does it integrate with our platform?  │
         └───────────────┬───────────────────────┘
                         │
                    ┌────┴────┐
                    │  YES    │                    ┌─── NO ──→ Can we build a bridge?
                    └────┬────┘                                    │
                         │                                    ┌────┴────┐
                         ▼                                    │  YES    │──→ Bridge, evaluate later
                    RETAIN IT                                  └─────────┘
                    (at least for now)                              │ NO
                                                                   ▼
                                                          Plan migration
                                                          (team drives timeline)
```

### Specific Evaluation by Area

**Source Control (they use something else, we use GitHub):**
- **Short-term:** Keep both. Set up mirroring for visibility if needed.
- **Medium-term:** Converge to one platform. GitHub is likely the target (broader ecosystem, Actions integration), but evaluate what they'd lose in migration.
- **Key question:** Do they have extensive automation tied to their SCM (bots, integrations, review workflows)? If yes, migration is more work than just moving repos.

**CI/CD Tooling (different from GitHub Actions):**
- **Short-term:** Keep running. Their pipelines deploy their apps; don't touch what works.
- **Medium-term:** For new services and new features, use the shared CI templates. Existing services migrate when they need significant pipeline changes anyway (natural migration).
- **Key question:** Is their CI self-hosted (Jenkins, Drone)? If yes, who maintains it? If it's one person, that's a bus factor risk that accelerates migration priority.

**Cloud Provider (different from ours):**
- **Short-term:** Run multi-cloud. ArgoCD can manage clusters on any provider.
- **Medium-term:** Evaluate consolidation based on cost, compliance, and capability. Some services may stay on the original cloud permanently if migration cost exceeds benefit.
- **Long-term:** Converge where it makes sense, but accept that multi-cloud may be permanent for specific workloads.
- **Key question:** Are their services architected for portability (containers, K8s)? Or are they tightly coupled to managed services (Lambda, Cloud Functions, proprietary databases)?

**Monitoring & Observability (different platforms):**
- **Short-term:** Federation. Set up a unified Grafana instance that queries both Prometheus instances (theirs and ours). Cross-platform dashboards give leadership visibility immediately.
- **Medium-term:** Standardize on the exporters and instrumentation format (OpenTelemetry). This allows backend consolidation later without re-instrumenting applications.
- **Key question:** Are they using a commercial platform (Datadog, New Relic)? If yes, do we want that contract? It might be better than our open-source stack.

---

## Integrating Two Engineering Organizations

### Technical Integration

| Area | Approach |
|------|----------|
| **Shared service catalog** | Backstage catalogs services from both orgs. This is the first integration point — everyone sees everything. |
| **Unified observability view** | Federated dashboards before backend consolidation. People need to see cross-org health immediately. |
| **Common security baseline** | Agree on minimum standards (image signing, vulnerability SLAs, no root). Enforce via policy-as-code. Apply to all new deployments; retrofit existing ones gradually. |
| **Shared CI templates (opt-in)** | Offer the templates. Don't mandate. When their team sees it saves time, they'll adopt. |
| **Single GitOps control plane** | ArgoCD manages both clusters. One place to see all deployment status across both clouds. |

### Organizational Integration

| Principle | Action |
|-----------|--------|
| **Embed, don't impose** | Assign a platform engineer to work alongside the acquired team for 2-3 months. Learn their patterns, help them learn ours. |
| **Joint architecture decisions** | Include acquired engineers in platform RFC/ADR processes immediately. They may have solved problems we haven't yet. |
| **Shared on-call** | Start with shadow rotations. Acquired team shadows our on-call to learn the platform; our team shadows theirs to learn their services. |
| **Respect existing expertise** | Their monitoring might be better. Their deployment scripts might be more sophisticated. Evaluate honestly — "not invented here" is a trap. |
| **Common communication channels** | Shared Slack channels for releases, incidents, and platform updates from day one. |

### What I Would NOT Do

- **Force an immediate "rip and replace"** — destroys trust and risks customer-facing outages
- **Treat their setup as inferior** — they shipped working software; respect that
- **Let two parallel platforms exist indefinitely** — set a convergence timeline (6-12 months) with checkpoints
- **Make the acquired team solely responsible for their own migration** — this is a shared effort; platform team provides support
- **Delay all integration** — waiting too long creates two entrenched camps that become harder to merge

---

## Revised Implementation Timeline

The original 12-month roadmap now shifts to accommodate integration:

| Month | Original Plan | Revised Plan |
|-------|--------------|--------------|
| 7 | Multi-cluster management | **Discovery & inventory of acquired systems** |
| 8 | Feature flags integration | **Observability federation (unified dashboards)** |
| 9 | Advanced canary strategies | **ArgoCD managing both clusters, shared security policies** |
| 10 | Backstage developer portal | **Backstage with both orgs' services cataloged** |
| 11 | SLO-based automation | **First acquired services migrated to shared platform** |
| 12 | Chaos engineering | **70% of acquired services on unified platform, legacy tooling decommission plan** |
| 13-15 | — | **Complete migration, retire legacy CI/CD and observability** |

**Net impact: the original 12-month roadmap extends to ~15 months.** This is acceptable. Rushing integration to maintain the original timeline would create more risk than the 3-month delay costs.

---

## Summary

| Question | Answer |
|----------|--------|
| Which recommendations remain unchanged? | GitOps, progressive delivery, image signing, policy-as-code, cadence tiers, DORA metrics, QA self-service |
| Which priorities shift? | Multi-cloud becomes essential, service catalog becomes critical, observability federation is urgent, onboarding pace slows |
| Immediate standardization or phased? | **Phased.** Observe → Bridge → Converge → Consolidate over 6-9 months |
| How to evaluate existing tooling? | Operational health + team expertise + integration cost + scalability. Retain what works, bridge what can connect, replace only what fails the criteria. |
| How to integrate two orgs? | Embed engineers cross-team, shared observability and catalog first, joint architectural decisions, opt-in adoption of shared tooling, convergence timeline with checkpoints |

The fundamental shift is from "build the platform for our org" to "build the platform that serves both orgs." The principles don't change. The scope, timeline, and empathy required all increase.
