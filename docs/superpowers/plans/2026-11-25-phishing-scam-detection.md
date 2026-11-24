# Phishing & Scam Detection Engine Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a zero-contact heuristic phishing/scam/malware detection engine with safe shortener unrolling, on-demand sandboxed deep-dive analysis, and adaptive security dossier UI in JevBeez.

**Architecture:** Progressive Escalation architecture:
1. Zero-Contact Heuristic Layer (<20ms): Analyzes URL structure, typosquatting, dangerous extensions (.apk/.exe), high-risk TLDs, IP hosts, and homoglyphs. If an unshortener is detected, unroll redirect hops safely without downloading bodies. If fatal, terminates immediately with a threat verdict (PHISHING/MALWARE/SCAM) without contacting the target server.
2. Sandboxed Deep Dive Layer (On-Demand): When requested with `deep_dive: true`, fetches max 100KB HTML under strict sandbox rules, inspects forms for credential/PIN/OTP harvesting and external exfiltration, and evaluates social engineering indicators via TypeSafe Jev.
3. Adaptive UI Dossier: Dynamically renders security threat dossiers with actionable mitigation advice, 1-click WhatsApp emergency warning broadcasts, and deep-dive inspection triggers.

**Tech Stack:** Go 1.23 (Fiber, net/http, x/net/html), Nuxt 4 (Vue 3, Tailwind CSS, TypeScript), Vitest, TypeSafe Jev API.

**Spec:** Documented in conversation history: Progressive Escalation approach with Safe Shortener Resolver, Heuristic Threat Engine, Sandbox Form Inspector, and Dynamic Threat Dossier.

## Global Constraints

- Zero hardcoding: all high-risk TLDs, malicious extensions, shortener domains, and brand profiles must reside in `internal/config/constants.go`.
- Zero-Risk & SSRF Safety: Every network hop (shortener unroll, sandbox fetch) MUST strictly block private/loopback/cloud-metadata IP ranges.
- Maximum payload limit for sandbox: 100KB with 3s hard timeout to prevent slowloris/decompression attacks.
- Clean code: Self-documenting naming, zero markdown bloat, and unit tests for every new module.

## Review Focus

1. Redirect chain leading to private IP (SSRF via shortener redirect) -> must be blocked immediately with `ErrPrivateIPBlocked`.
2. Unshortener infinite redirect loop -> must abort cleanly at max 5 hops without hanging.
3. URL with `@` symbol trick (e.g. `http://bca.co.id@phishing.xyz/login`) -> must flag as phishing impersonating BCA on phishing.xyz.
4. Subdomain trap (e.g. `bca.co.id.promo-hadiah.xyz`) -> must flag as phishing impersonating BCA.
5. Direct link to `.apk` or `.exe` file -> must flag as MALWARE without fetching the binary.

---

### Task 1: Constants and Configuration for Phishing/Scam Engine

**Files:**
- Modify: `apps/api/internal/config/constants.go`
- Test: `apps/api/internal/config/constants_test.go`

**Interfaces:**
- Produces:
  - `VerdictPhishing`, `VerdictMalware`, `VerdictScam` constants
  - `DangerousExtensions []string`
  - `HighRiskTLDs []string`
  - `ShortenerDomains []string`
  - `BrandProfiles map[string]string` (brand keyword -> official domain)
  - `MaxShortenerHops = 5`
  - `MaxSandboxPayloadBytes = 100 * 1024`

- [ ] **Step 1: Write the failing test**
Create `apps/api/internal/config/constants_test.go` verifying that threat verdicts and list constants are non-empty and accessible.

- [ ] **Step 2: Run test to verify it fails**
Run: `go test ./internal/config/...` in `apps/api`
Expected: FAIL

- [ ] **Step 3: Implement constants in `apps/api/internal/config/constants.go`**
Add constants for verdicts, file extensions, TLDs, shortener domains, and brand profiles.

- [ ] **Step 4: Run test to verify it passes**
Run: `go test ./internal/config/...` in `apps/api`
Expected: PASS

---

### Task 2: Safe Shortener Resolver

**Files:**
- Create: `apps/api/internal/security/shortener.go`
- Create: `apps/api/internal/security/shortener_test.go`

**Interfaces:**
- Consumes: `config.ShortenerDomains`, `config.MaxShortenerHops`, `security.IsPrivateIP`
- Produces:
  - `IsShortener(urlStr string) bool`
  - `ResolveShortener(ctx context.Context, rawURL string, maxHops int) (string, []string, error)`

- [ ] **Step 1: Write the failing test**
Create `shortener_test.go` testing shortener detection and mocked safe HTTP unshortening (including loop detection and private IP redirection rejection).

- [ ] **Step 2: Run test to verify it fails**
Run: `go test ./internal/security -run TestShortener` in `apps/api`
Expected: FAIL

- [ ] **Step 3: Implement `shortener.go`**
Implement safe HTTP client using `HEAD`/`GET` (without reading body), SSRF validation on every redirect location, and maximum hop limits.

- [ ] **Step 4: Run test to verify it passes**
Run: `go test ./internal/security -run TestShortener` in `apps/api`
Expected: PASS

---

### Task 3: Zero-Contact URL Heuristic Threat Engine

**Files:**
- Create: `apps/api/internal/security/heuristic.go`
- Create: `apps/api/internal/security/heuristic_test.go`

**Interfaces:**
- Consumes: `config.DangerousExtensions`, `config.HighRiskTLDs`, `config.BrandProfiles`, `model.ThreatCheckResult`
- Produces:
  - `CheckURLHeuristics(rawURL string) *ThreatCheckResult`

- [ ] **Step 1: Write the failing test**
Create `heuristic_test.go` testing:
  - Direct `.apk` download URL -> MALWARE
  - `@` credential trick (`bca.co.id@evil.com`) -> PHISHING
  - Subdomain trap (`bca.co.id.evil-domain.xyz`) -> PHISHING
  - Brand keyword in non-official domain on high-risk TLD -> PHISHING
  - Raw public IP host with sensitive path -> SCAM/PHISHING
  - Legitimate news URL (`kompas.com/artikel`) -> Not a threat

- [ ] **Step 2: Run test to verify it fails**
Run: `go test ./internal/security -run TestHeuristic` in `apps/api`
Expected: FAIL

- [ ] **Step 3: Implement `heuristic.go`**
Parse URL hostname and path, extract apex domain, check brand impersonation against official domains, detect dangerous file extensions and anomalies.

- [ ] **Step 4: Run test to verify it passes**
Run: `go test ./internal/security -run TestHeuristic` in `apps/api`
Expected: PASS

---

### Task 4: Sandboxed Scraper & Form Threat Inspector

**Files:**
- Create: `apps/api/internal/scraper/sandbox.go`
- Create: `apps/api/internal/scraper/sandbox_test.go`

**Interfaces:**
- Consumes: `config.MaxSandboxPayloadBytes`, `security.SafeDialer`
- Produces:
  - `InspectFormThreats(ctx context.Context, targetURL string) ([]string, string, error)`

- [ ] **Step 1: Write the failing test**
Create `sandbox_test.go` using `httptest.Server` with mock HTML containing login/OTP/PIN inputs and external exfiltration actions.

- [ ] **Step 2: Run test to verify it fails**
Run: `go test ./internal/scraper -run TestSandbox` in `apps/api`
Expected: FAIL

- [ ] **Step 3: Implement `sandbox.go`**
Stream max 100KB with 3s timeout, parse HTML token tree, extract forms, identify sensitive input fields (password, pin, otp, card), and analyze form destination targets.

- [ ] **Step 4: Run test to verify it passes**
Run: `go test ./internal/scraper -run TestSandbox` in `apps/api`
Expected: PASS

---

### Task 5: Model & Decision Types Update

**Files:**
- Modify: `apps/api/internal/model/types.go`
- Modify: `apps/api/internal/jev/questions.go`
- Modify: `apps/api/internal/decision/evaluator.go`
- Modify: `apps/api/internal/decision/evaluator_test.go`

**Interfaces:**
- Produces updated `AnalyzeRequest` (with `DeepDive bool`), `VerificationResult` (with `IsThreat`, `ThreatType`, `ThreatReasons`, `ResolvedURL`, `FormThreats`), and threat-aware evaluation rules.

- [ ] **Step 1: Write failing test in evaluator_test.go**
Add tests for threat handling in `evaluator_test.go`.

- [ ] **Step 2: Update `types.go`, `questions.go`, and `evaluator.go`**
Add fields, Jev threat questions, and threat evaluation logic.

- [ ] **Step 3: Run test to verify it passes**
Run: `go test ./internal/decision/...` in `apps/api`
Expected: PASS

---

### Task 6: API Handler Pipeline Integration

**Files:**
- Modify: `apps/api/internal/handler/analyze.go`
- Modify: `apps/api/cmd/server/main.go`
- Test: `apps/api/internal/handler/analyze_test.go`

**Interfaces:**
- Integrates Shortener Resolver -> Heuristic Engine -> Zero-Contact Fast Verdict vs Sandboxed Deep Dive.

- [ ] **Step 1: Write tests for threat analysis pipeline in `analyze_test.go`**
Verify that a phishing URL returns immediately without scraping, and `deep_dive: true` performs form inspection.

- [ ] **Step 2: Implement threat handling in `analyze.go`**
Execute unshortening, run heuristic check, return early for threats on default mode, or trigger sandbox inspection on `deep_dive: true`.

- [ ] **Step 3: Run tests in `apps/api`**
Run: `go test ./...` in `apps/api`
Expected: ALL PASS

---

### Task 7: Frontend Adaptive Security Dossier & Emergency Response

**Files:**
- Modify: `apps/web/types/index.ts`
- Modify: `apps/web/locales/id.ts` and `locales/en.ts`
- Modify: `apps/web/components/VerdictCard.vue`
- Modify: `apps/web/composables/useAnalyze.ts`
- Test: `apps/web/test/threat-verdict.test.ts`

**Interfaces:**
- Renders Threat Banner, Threat Reasons, Dangerous Form alerts, "Buka Analisis Mendalam" button, and 1-Click WhatsApp Warning Share.

- [ ] **Step 1: Write frontend test in `apps/web/test/threat-verdict.test.ts`**
Verify threat rendering, action recommendations, and deep dive trigger.

- [ ] **Step 2: Update `types/index.ts`, `locales/`, `VerdictCard.vue`, and `useAnalyze.ts`**
Implement the dynamic threat dossier UI, deep dive trigger, and emergency WhatsApp broadcast message.

- [ ] **Step 3: Run frontend tests**
Run: `bun run test` in `apps/web`
Expected: ALL PASS

---

### Task 8: End-to-End Verification & Documentation

- [ ] **Step 1: Test full backend & frontend compilation and test suites**
- [ ] **Step 2: Verify zero-contact fast verdict with real scam link samples and ensure no regression on valid news articles**
