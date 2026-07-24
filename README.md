# CI/CD Demo: Kubernetes-Native Deployment Pipeline

## Overview

This repository demonstrates a production-grade CI/CD workflow for a containerized Go API deployed to Kubernetes using GitOps principles. The focus is on the **delivery process** — security gates, progressive deployment, approval workflows, and rollback mechanisms — not the application complexity.

---

## Application

A lightweight Go REST API with two endpoints:

| Endpoint | Purpose |
|----------|---------|
| `GET /` | Returns service info (name, version, timestamp, hostname) |
| `GET /healthz` | Health check for Kubernetes liveness/readiness probes |

The app compiles to a single static binary (~15MB container image using distroless), starts in milliseconds, and has zero runtime dependencies — ideal for demonstrating Kubernetes deployment patterns.

---

## Architecture

```
┌──────────────┐     ┌────────────────────────────────────┐     ┌─────────────────────┐
│  GitHub Repo │────▶│  GitHub Actions (CI Pipeline)       │────▶│  GHCR (Container    │
│  (Source)    │     │  Lint → Test → SAST → SCA → Build  │     │  Registry)          │
└──────────────┘     └────────────────────────────────────┘     └─────────────────────┘
                                                                          │
                         ┌── Cosign signs image                           │
                         ├── Syft generates SBOM                          │
                         └── Trivy + Grype scan image                     │
                                                                          ▼
┌──────────────┐     ┌──────────────────┐        ┌──────────────────────────────────┐
│  Grafana     │◀────│  Kubernetes      │◀───────│  ArgoCD (GitOps Sync)            │
│  (Observe)   │     │  (Argo Rollouts) │        │  values-staging.yaml updated     │
└──────────────┘     └──────────────────┘        └──────────────────────────────────┘
                              │                                    ▲
                              │ Canary analysis                    │
                              ▼                                    │
                     ┌──────────────────┐        ┌──────────────────────────────────┐
                     │  Prometheus      │        │  Approval Gate → values-prod.yaml │
                     │  (AnalysisRun)   │        │  updated → ArgoCD syncs prod     │
                     └──────────────────┘        └──────────────────────────────────┘
```

---

## Design Decisions

### Why Go?

| Factor | Benefit |
|--------|---------|
| Static binary | ~15MB distroless image, fast registry pulls during rollouts |
| No runtime deps | No interpreter, no pip/npm — smaller attack surface |
| Millisecond startup | Readiness probes pass immediately, canary pods serve traffic fast |
| Built-in testing | `go test` ships with the language, no test framework to install |
| Ecosystem fit | K8s, ArgoCD, Prometheus — all written in Go; signals domain familiarity |

### Why GitOps (ArgoCD)?

- **Git is the single source of truth** — you can see what's deployed by reading the repo, and the commit history is the deployment audit trail.
- **Drift detection** — ArgoCD alerts when cluster state diverges from desired state.
- **Declarative rollback** — `git revert` on a promotion commit rolls back production.
- **Separation of concerns** — CI builds and pushes images; CD (ArgoCD) handles deployment to clusters. They're decoupled.

### Why Argo Rollouts (Canary)?

Standard Kubernetes Deployments do rolling updates that are hard to observe and impossible to pause mid-rollout. Argo Rollouts provides:

- **Progressive traffic shifting** (10% → 30% → 50% → 100%)
- **Automated analysis** — Prometheus queries evaluate error rate at each step
- **Instant abort** — if metrics breach SLOs, traffic returns to stable pods in seconds
- **Pause steps** — bake time between weight increases to observe behavior

### Why Environment Approval Gates?

GitHub Environments with required reviewers create a human checkpoint between automated stages. It's important to understand what's actually being approved in a GitOps model:

```
Build & Scan → [⏸️ Staging Approval] → Update values-staging.yaml → ArgoCD auto-syncs
                                                                           ↓
                                        [⏸️ Production Approval] → Update values-production.yaml → QA clicks Sync in ArgoCD
```

**What the GitHub approval controls:** whether the image tag gets committed to the Helm values file. It gates the *intent to deploy*, not the deployment itself.

**What ArgoCD controls:** the actual cluster state change. For staging, ArgoCD auto-syncs (deploys immediately when values change). For production, ArgoCD requires a manual Sync action — this is the second approval layer.

| Gate | Controls | Who approves |
|------|----------|--------------|
| GitHub Environment (staging) | Writing image tag to `values-staging.yaml` | CI reviewer |
| ArgoCD auto-sync (staging) | Deploying to staging cluster | Automatic |
| GitHub Environment (production) | Writing image tag to `values-production.yaml` | Engineering lead |
| ArgoCD manual Sync (production) | Deploying to production cluster | QA engineer via ArgoCD UI |

This two-layer approach means production requires both "yes, this image is ready" (GitHub) AND "yes, deploy it now" (ArgoCD). QA controls the timing of the actual deployment without needing kubectl access.

### Why Dual Container Scanners (Trivy + Grype)?

No single scanner catches everything. Running both increases coverage:
- **Trivy** — broad CVE database, fast, good for OS packages
- **Grype** — strong on language-specific dependencies, different vulnerability feeds

### Why Cosign Image Signing?

Supply chain security. The promotion workflow **verifies the image signature** before deploying to production, ensuring:
- The image was built by this CI pipeline (not tampered with)
- It hasn't been modified in the registry after build
- Only images that passed all security gates can reach production

### Why Kyverno Policies?

Admission control as code. Policies are tested in CI before reaching the cluster:
- No `:latest` tags (prevents non-reproducible deploys)
- Resource limits required (prevents noisy-neighbor issues)
- Non-root, read-only rootfs (container hardening)
- Probes required (ensures K8s can manage pod lifecycle)

---

## Repository Structure

```
platform-engineering-assessment/
├── .github/
│   └── workflows/
│       ├── ci.yaml                   # Full CI/CD pipeline with approval gates
│       ├── promote.yaml              # Manual promotion (hotfix/re-deploy)
│       └── k8s-manifest-scan.yaml    # Helm lint + kubeconform + Kyverno + Checkov
├── app/                              # Application source code
│   ├── main.go                       # API server
│   ├── main_test.go                  # Unit tests
│   ├── go.mod                        # Go module
│   └── .golangci.yml                 # Linter configuration
├── helm/                             # Helm chart for deployment
│   ├── Chart.yaml
│   ├── values.yaml                   # Base values
│   ├── values-staging.yaml           # Staging overrides (auto-updated by CI)
│   ├── values-production.yaml        # Production overrides (updated on promotion)
│   └── templates/
│       ├── rollout.yaml              # Argo Rollout (canary strategy)
│       ├── service.yaml              # Stable service
│       ├── service-canary.yaml       # Canary service (traffic splitting)
│       ├── hpa.yaml                  # Horizontal Pod Autoscaler
│       ├── pdb.yaml                  # Pod Disruption Budget
│       └── analysis-template.yaml    # Prometheus-based canary analysis
├── kyverno/                          # Policy-as-code
│   └── policies/
│       ├── disallow-latest-tag.yaml
│       ├── require-labels.yaml
│       ├── require-nonroot.yaml
│       ├── require-probes.yaml
│       ├── require-resource-limits.yaml
│       └── restrict-privilege-escalation.yaml
├── argocd/                           # ArgoCD Application definitions
│   ├── app-staging.yaml              # Auto-sync enabled
│   └── app-production.yaml           # Manual sync (QA triggers)
├── platform-proposal/                # Written proposal documents
│   ├── release_engineering_proposal.md
│   ├── platform_strategy_proposal.md
│   ├── acquisition_integration_strategy.md
│   └── team_adoption_operations.md
├── Dockerfile                        # Multi-stage build (distroless)
├── assessment_brief.md               # Original assessment requirements
└── README.md
```

---

## CI/CD Pipeline Flow

### Main Pipeline (`ci.yaml`)

```
┌─────────────────────────────────────────────────────────────────────┐
│                        ON PUSH TO MAIN                               │
├─────────────────────────────────────────────────────────────────────┤
│                                                                       │
│  ┌──────────┐   ┌──────────────┐   ┌──────────────┐                │
│  │  Lint &  │   │  SAST        │   │  SCA         │                │
│  │  Test    │   │  CodeQL +    │   │  Govulncheck │   (parallel)   │
│  │          │   │  Semgrep     │   │  + Trivy FS  │                │
│  └────┬─────┘   └──────┬───────┘   └──────┬───────┘                │
│       │                 │                   │                         │
│       └─────────────────┼───────────────────┘                        │
│                         ▼                                             │
│              ┌─────────────────────┐                                 │
│              │  Build, Sign & Scan │                                 │
│              │  • Docker build     │                                 │
│              │  • Cosign sign      │                                 │
│              │  • Trivy scan       │                                 │
│              │  • Grype scan       │                                 │
│              └──────────┬──────────┘                                 │
│                         ▼                                             │
│              ┌─────────────────────┐                                 │
│              │  SBOM & Dep-Track   │                                 │
│              │  • Syft generate    │                                 │
│              │  • Upload to DTrack │                                 │
│              │  • Analyze results  │                                 │
│              └──────────┬──────────┘                                 │
│                         ▼                                             │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │  ⏸️  APPROVAL GATE: staging environment (QA team reviews)   │    │
│  └──────────────────────────────┬──────────────────────────────┘    │
│                                 ▼                                     │
│              ┌─────────────────────┐                                 │
│              │  Deploy to Staging  │                                 │
│              │  (update values →   │                                 │
│              │   ArgoCD auto-sync) │                                 │
│              └──────────┬──────────┘                                 │
│                         ▼                                             │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │  ⏸️  APPROVAL GATE: production environment (lead reviews)   │    │
│  └──────────────────────────────┬──────────────────────────────┘    │
│                                 ▼                                     │
│              ┌─────────────────────┐                                 │
│              │ Deploy to Production│                                 │
│              │ • Verify signature  │                                 │
│              │ • Update values     │                                 │
│              │ • ArgoCD → Rollout  │                                 │
│              └─────────────────────┘                                 │
│                                                                       │
├─────────────────────────────────────────────────────────────────────┤
│                    📋 PIPELINE SUMMARY                                │
└─────────────────────────────────────────────────────────────────────┘
```

### Manifest Validation (`k8s-manifest-scan.yaml`)

Triggered on changes to `helm/`, `kyverno/`, or `argocd/` paths:

1. **Helm lint** — validates chart syntax
2. **Helm template** — renders manifests for both environments
3. **kubeconform** — validates against Kubernetes API schemas (including CRDs)
4. **Kyverno CLI** — tests rendered manifests against all policies
5. **Checkov** — security best practices scan on K8s resources

### Manual Promotion (`promote.yaml`)

For hotfixes or re-promotions outside the normal flow:
- Verifies image signature (Cosign)
- Checks image exists and was deployed to staging
- Requires `production` environment approval
- Updates production values file

---

## Security Controls

| Layer | Tool | Purpose |
|-------|------|---------|
| Source code | CodeQL + Semgrep | SAST — finds vulnerabilities in application code |
| Dependencies | Govulncheck + Trivy FS | SCA — finds known CVEs in Go modules |
| Container image | Trivy + Grype | Scans OS packages and app dependencies in built image |
| Supply chain | Cosign (Sigstore) | Image signing and verification — tamper detection |
| SBOM | Syft + Dependency-Track | Component inventory, continuous vulnerability monitoring |
| K8s manifests | Kyverno + Checkov | Policy enforcement before deployment |
| Runtime | Argo Rollouts AnalysisTemplate | Prometheus-based canary health validation |

---

## Rollback Strategy

Three mechanisms, ordered by speed:

| Method | Time | How | When to use |
|--------|------|-----|-------------|
| **Argo Rollouts Abort** | ~5 seconds | Abort canary → all traffic returns to stable pods | During canary rollout, metrics look bad |
| **ArgoCD Rollback** | ~30 seconds | Click "Rollback" in ArgoCD UI → redeploys previous revision | After full rollout, need immediate revert |
| **Git Revert** | ~2 minutes | `git revert <promotion-commit>` → ArgoCD auto-syncs | Full audit trail needed, or ArgoCD UI unavailable |

All three methods are available to QA engineers through either the ArgoCD UI or a simple git command.

---

## Running Locally

```bash
# Run the app
cd app
go run main.go

# Run tests
go test -v ./...

# Build container
docker build -t cicd-demo:local .

# Run container
docker run -p 8080:8080 cicd-demo:local

# Test endpoints
curl http://localhost:8080/
curl http://localhost:8080/healthz

# Lint Helm chart
helm lint helm/

# Test Kyverno policies against rendered manifests
helm template cicd-demo helm/ -f helm/values-production.yaml | \
  kyverno apply kyverno/policies/ --resource /dev/stdin
```

---

## Scaling This Approach

For a larger engineering organization:

| Challenge | Solution |
|-----------|----------|
| Many services, same patterns | **Shared Helm library chart** — base chart with org defaults, services only supply values |
| New services need pipelines | **Reusable GitHub Actions workflows** — `workflow_call` templates consumed by all repos |
| Auto-discovery of services | **ArgoCD ApplicationSets** — new directory in GitOps repo = new ArgoCD Application |
| QA needs self-service | **ArgoCD RBAC** — namespace-scoped sync permissions per team |
| Policy consistency | **Kyverno as admission controller** — policies enforce standards cluster-wide at deploy time |
| Multi-cluster | **ArgoCD multi-cluster** — single control plane manages dev/staging/prod clusters |
| Measuring improvement | **DORA metrics** from GitOps data — deploy frequency, lead time, failure rate, MTTR |

---

## Prerequisites (for full deployment)

| Component | Purpose |
|-----------|---------|
| Kubernetes cluster (EKS/GKE/AKS) | Runtime |
| ArgoCD | GitOps-based continuous delivery |
| Argo Rollouts controller | Progressive canary deployments |
| Prometheus | Metrics for canary analysis |
| Grafana | Dashboards with deploy annotations |
| Dependency-Track instance | SBOM analysis and monitoring |
| GitHub repo with Actions enabled | CI pipeline |

### GitHub Repository Settings Required

1. **Environments:** Create `staging` and `production` with required reviewers
2. **Secrets:** `DEPENDENCY_TRACK_URL`, `DEPENDENCY_TRACK_API_KEY`, `SEMGREP_APP_TOKEN` (optional)
3. **Permissions:** Enable `id-token: write` for Cosign keyless signing

---

## Tradeoffs & Discussion Points

| Decision | Tradeoff | Rationale |
|----------|----------|-----------|
| Monorepo (app + infra together) | Less separation of concerns | Simpler for demo; production would likely split into app repo + GitOps repo |
| Keyless Cosign (Sigstore) | Requires GitHub OIDC, tied to GitHub | No key management, no secret rotation, verifiable provenance |
| Dual scanners (Trivy + Grype) | Slower pipeline | Better coverage; different vulnerability databases catch different issues |
| Kyverno over OPA/Gatekeeper | Less powerful policy language | Simpler YAML-based policies, easier for teams to write and understand |
| Environment gates in CI | Blocks the entire pipeline | Clear promotion path; alternative would be separate deployment workflows |
| Canary over blue-green | More complex, needs service mesh for traffic splitting | Gradual rollout with real traffic; catches issues blue-green would miss |
