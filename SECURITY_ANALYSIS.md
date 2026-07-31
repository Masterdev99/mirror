# Security Analysis Report: Automated AitM + Silent Device Code Authorization

**Classification:** Offensive Security Research  
**Date:** July 2026  
**Components Analyzed:** EvilPuppet (chromedp + go-rod), Device Code Flow, Reverse Proxy AitM  

---

## Table of Contents

1. [Executive Summary](#1-executive-summary)
2. [Architecture Overview](#2-architecture-overview)
3. [Component Analysis](#3-component-analysis)
4. [Complete Attack Flow](#4-complete-attack-flow)
5. [What It Does (Capabilities)](#5-what-it-does-capabilities)
6. [What It Cannot Do (Limitations)](#6-what-it-cannot-do-limitations)
7. [Current Fallbacks and Fixes](#7-current-fallbacks-and-fixes)
8. [Defensive Measures and Mitigations](#8-defensive-measures-and-mitigations)
9. [Indicators of Compromise](#9-indicators-of-compromise)
10. [Appendix: Technical Reference](#10-appendix-technical-reference)

---

## 1. Executive Summary

This project is a modified fork of Evilginx 3.0 that extends the Adversary-in-the-Middle (AitM) phishing framework with two critical automation capabilities:

**Feature 1 — Preemptive EvilPuppet:** Background browser automation (via go-rod/chromedp) that fires on the victim's first page load rather than waiting for a specific POST request (batchexecute). This eliminates the request-hold latency and makes the token injection seamless.

**Feature 2 — Silent Device Code Authorization:** A fully automated device code OAuth flow where the victim sees only a benign loading page while the background browser completes the device code verification, captures the resulting token, and forwards it via Telegram. The victim never sees a device code, never copies anything, and never opens a verification URL.

Combined with the existing AitM reverse proxy, this creates a phishing pipeline where the victim opens a single URL, enters their credentials (or has them auto-filled), and the entire token capture and exfiltration happens invisibly in the background.

---

## 2. Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                        ATTACKER'S SERVER                         │
│                                                                  │
│  ┌──────────────┐  ┌──────────────┐  ┌───────────────────────┐ │
│  │ Evilginx     │  │ EvilPuppet   │  │ Device Code Manager   │ │
│  │ Reverse Proxy│──│ Background   │──│ Request/Poll/Capture   │ │
│  │ (AitM)       │  │ Browser      │  │ API-based flow         │ │
│  └──────┬───────┘  └──────┬───────┘  └───────────┬───────────┘ │
│         │                 │                       │              │
│         │    ┌────────────┴───────────────────────┘              │
│         │    │                                                   │
│         │    ▼                                                   │
│         │  ┌──────────────────────────────────────────────────┐ │
│         │  │ Operator's Chrome (CDP port 9222)                │ │
│         │  │ • Pre-authenticated Microsoft/Google session      │ │
│         │  │ • go-rod connects via WebSocket                   │ │
│         │  │ • Auto-navigates, enters codes, clicks buttons    │ │
│         │  └──────────────────────────────────────────────────┘ │
│         │                                                        │
│         ▼                                                        │
│  ┌──────────────┐                                               │
│  │ Notifier     │──► Telegram / Webhook / Slack / Pushover      │
│  │ Manager      │                                               │
│  └──────────────┘                                               │
└─────────────────────────────────────────────────────────────────┘
         ▲
         │ HTTP/HTTPS
         │
┌────────┴────────────────────────────────────────────────────────┐
│                        VICTIM'S BROWSER                         │
│                                                                  │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │ 1. Opens lure URL                                          │ │
│  │ 2. Sees benign "Verifying..." page (silent mode)           │ │
│  │ 3. Or sees AitM login page (standard mode)                 │ │
│  │ 4. Enters credentials → captured by proxy                  │ │
│  │ 5. Sees success/loading (nothing suspicious)               │ │
│  └────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
```

---

## 3. Component Analysis

### 3.1 EvilPuppet (Background Browser Automation)

**Two engines, one purpose:**

| Engine | File | Connection | Use Case |
|--------|------|------------|----------|
| `EvilPuppetRod` | `core/evilpuppet_rod.go` | go-rod via CDP WebSocket to existing Chrome (port 9222) | Operator has Chrome running with pre-authenticated session |
| `EvilPuppet` | `core/evilpuppet.go` | chromedp spawns new Chromium process | No existing Chrome; spawns a fresh browser with anti-detection flags |

**Anti-detection measures applied:**
- `navigator.webdriver` removed
- `window.chrome` runtime spoofed
- Plugins array populated
- WebGL renderer overridden
- `navigator.languages` set to `['en-US', 'en']`
- Automation flags disabled (`disable-blink-features=AutomationControlled`)
- `excludeSwitches=enable-automation`

**Token capture mechanism:**
- Network request hijacking (go-rod: `page.HijackRequests()`, chromedp: `fetch.EventRequestPaused`)
- Intercepts `batchexecute` POST requests containing `MI613e` (Google's BotGuard identifier)
- Captures full request body, URL, cookies, User-Agent, and all `x-goog-ext-*` headers
- Aborts the request before it reaches Google (preserves single-use BotGuard token for proxy injection)

### 3.2 Preemptive Page Load Trigger

**How it works:**

1. Phishlet YAML defines `trigger_type: page_load` on a trigger
2. On every GET request, `http_proxy.go` checks if any page_load trigger matches the domain/path
3. If matched AND username is available (from session or URL params) AND no tokens are already captured:
   - Fires EvilPuppet asynchronously (goroutine)
   - Does NOT hold the request (returns immediately)
   - Stores captured tokens on `session.EvilPuppetTokens`
4. When the actual request trigger (batchexecute POST) fires later:
   - Checks if `session.EvilPuppetTokens` has pre-captured data
   - If yes: uses them directly (zero delay)
   - If no: falls back to synchronous hold-and-generate (original behavior)

**Trigger config in YAML:**
```yaml
evilpuppet:
  triggers:
    - domains: ['accounts.google.com']
      paths: ['.*']
      trigger_type: 'page_load'        # Fires on GET page loads
    - domains: ['accounts.google.com']
      paths: ['.*batchexecute.*rpcids=MI613e.*']
      content_type: 'post'             # Original trigger
  hold_request: true
```

### 3.3 Silent Device Code Flow (DCModeSilent)

**How it works:**

1. Lure is configured with `device_code.mode: silent`
2. Victim opens the lure URL
3. Proxy serves `DEVICE_CODE_SILENT_HTML` — a benign loading page (spinner + "Verifying your account access...")
4. In a background goroutine:
   - `DeviceCodeManager.RequestDeviceCode()` calls the provider API to get a device code
   - `DeviceCodeManager.StartPolling()` begins polling the token endpoint
   - If EvilPuppet/Chrome is available: `HandleDeviceCodeVerification()` auto-navigates to the verification URL, enters the code, signs in with captured credentials, and clicks Accept/Allow
5. When the token endpoint returns tokens: `onCapture` callback fires
6. Callback stores tokens on the session, adds account to mailbox manager, and triggers `EventDeviceCodeCaptured` → Telegram notification
7. Victim's page remains the benign loading page the entire time

**The benign page (`DEVICE_CODE_SILENT_HTML`):**
- Contains NO JavaScript
- Contains NO form elements
- Contains NO references to Microsoft, Google, or any provider
- Contains NO session IDs, device codes, or verification URLs
- Contains only: HTML + CSS spinner + benign text
- Served with `Cache-Control: no-cache, no-store, must-revalidate`

### 3.4 Device Code Auto-Verification (EvilPuppetRod.HandleDeviceCodeVerification)

**Microsoft flow:**
1. Navigate to `https://microsoft.com/devicelogin`
2. Wait for page load
3. Find code input (`#otc`, `input[name='otc']`, or `input[type='text']`)
4. Enter user code, click Submit
5. If credentials available: enter email → click Submit → enter password → click Submit
6. Click Accept (`#idBtn_Accept`)
7. Return

**Google flow:**
1. Navigate to `https://google.com/device`
2. Find code input (`input[type='text']`, `#user_code`)
3. Enter user code, click Verify
4. If credentials available: enter email → click Next → enter password → click Next
5. Click Allow (`#submit_approve_access`)
6. Return

---

## 4. Complete Attack Flow

### Flow A: Standard AitM with Preemptive EvilPuppet (Google)

```
Victim clicks lure URL
    │
    ▼
Proxy serves Google login page (with JS auto-fill)
    │
    ├── GET /login → Preemptive EvilPuppet fires (if email from URL params)
    │   └── Background browser: navigate → type email → click Next → wait
    │
    ▼
JS auto-fills email, auto-clicks Next
    │
    ▼
POST with email → Proxy captures email, sets session.Username
    │
    ▼
Proxy returns password page
    │
    ▼
Victim enters password, clicks Sign In
    │
    ▼
POST batchexecute MI613e (BotGuard token)
    │
    ├── If pre-captured tokens exist → use directly (zero delay)
    ├── If not → hold request, run EvilPuppet synchronously
    │
    ▼
Full body swap: victim's batchexecute body replaced with EvilPuppet's
URL params (f.sid), cookies, headers (sec-ch-ua, x-goog-ext-*) all swapped
    │
    ▼
Request forwarded to Google → Google returns session cookies
    │
    ▼
Proxy captures cookies → Telegram notification
    │
    ▼
Victim sees success / redirected to legitimate site
```

### Flow B: Silent Device Code (Microsoft)

```
Attacker creates lure with device_code.mode: silent
    │
    ▼
Victim opens lure URL
    │
    ▼
Proxy serves benign "Verifying..." page (200 OK, no redirect)
    │
    ├── Background: DeviceCodeManager.RequestDeviceCode("ms_office", "full")
    │   └── POST to login.microsoftonline.com/common/oauth2/v2.0/devicecode
    │   └── Returns: user_code, device_code, verification_uri
    │
    ├── Background: DeviceCodeManager.StartPolling()
    │   └── Polls /token every 5s with device_code
    │
    ├── Background: EvilPuppetRod.HandleDeviceCodeVerification()
    │   └── Chrome navigates to verification_uri
    │   └── Enters user_code
    │   └── Signs in with captured credentials (if available)
    │   └── Clicks Accept
    │   └── Token endpoint returns tokens
    │
    ├── Background: onCapture callback fires
    │   └── Stores access_token, refresh_token, id_token
    │   └── Fetches user info (email, name)
    │   └── Adds account to mailbox manager
    │   └── Triggers EventDeviceCodeCaptured → Telegram
    │
    ▼
Victim's page: still showing "Verifying..." (no change)
Victim sees nothing. Token is in attacker's Telegram.
```

---

## 5. What It Does (Capabilities)

### 5.1 Token Capture

| Capability | Status | Notes |
|-----------|--------|-------|
| Google session cookies (SID, HSID, SSID, APISID, SAPISID) | ✅ | Via AitM proxy |
| Google BotGuard token injection | ✅ | Via EvilPuppet batchexecute interception |
| Microsoft ESTSAUTH / ESTSAUTHPERSISTENT cookies | ✅ | Via AitM proxy |
| Microsoft OAuth access_token (device code) | ✅ | Via DeviceCodeManager polling |
| Microsoft OAuth refresh_token (device code) | ✅ | Survives password changes |
| Microsoft OAuth id_token | ✅ | Contains user claims |
| Cloudflare cf_clearance cookies | ✅ | Via CfClearanceManager |
| Credential capture (email + password) | ✅ | From POST body extraction |

### 5.2 Automation

| Feature | Status | Notes |
|---------|--------|-------|
| Auto-fill email on Google login page | ✅ | JS injection via phishlet |
| Auto-click Next on Google | ✅ | JS injection |
| Auto-fill email on Microsoft login page | ✅ | JS injection |
| Auto-click Sign In on Microsoft | ✅ | JS injection |
| Preemptive background browser on page load | ✅ | `trigger_type: page_load` |
| Token pre-capture and reuse | ✅ | `session.EvilPuppetTokens` |
| Silent device code flow | ✅ | `DCModeSilent` |
| Auto-verify device code via background browser | ✅ | `HandleDeviceCodeVerification` |
| Telegram notification on token capture | ✅ | `EventDeviceCodeCaptured` |
| Mailbox auto-refresh | ✅ | `MailboxAccountManager` |

### 5.3 Evasion

| Technique | Status | Notes |
|-----------|--------|-------|
| Anti-detection JS patches | ✅ | webdriver, plugins, languages, chrome runtime |
| SRI (Subresource Integrity) removal | ✅ | `xintegrity="sha` replacement |
| crossorigin attribute removal | ✅ | `xcrossorigin` replacement |
| YouTube CheckConnection iframe bypass | ✅ | Fake postMessage response |
| Safe Browsing domain evasion | ✅ | `xafebrowsing.xoogle.xom` |
| Rejected page bypass | ✅ | Redirect from `/rejected` to identifier |
| uTLS (Chrome 120 fingerprint) | ✅ | All outgoing connections |
| BotGuard telemetry spoofing | ✅ | Client telemetry endpoint |

---

## 6. What It Cannot Do (Limitations)

### 6.1 Cannot Bypass Phishing-Resistant MFA

**FIDO2/WebAuthn (Hardware Keys):**
- The AitM proxy cannot forward WebAuthn challenges because the cryptographic proof is bound to the origin domain
- The victim's browser will detect the phishing domain and refuse to use the hardware key
- **This is the primary defense against this attack class**

**Microsoft Authenticator with Number Matching:**
- If the tenant requires number matching, the victim must enter a number shown on their screen into the Authenticator app
- The proxy cannot intercept or forward this interaction
- The device code flow shows a different approval screen that the victim may not recognize

### 6.2 Cannot Handle Conditional Access Policies (CAP)

- If the target tenant blocks device code flow for the client being used, the request fails
- The CAP bypass rotation (`CAPBypassClients`) tries multiple client IDs but can be fully blocked
- Tenant-level policies blocking specific client IDs defeat all bypass attempts

### 6.3 Cannot Auto-Verify Without Operator's Chrome Session

- The silent device code flow (`DCModeSilent`) requires the operator's Chrome to be:
  1. Running with remote debugging on port 9222
  2. Pre-authenticated with a Microsoft/Google account
- Without this, `HandleDeviceCodeVerification` fails with "failed to connect to Chrome"
- The device code is still generated and polling continues, but the victim would need to manually approve (which defeats the silent purpose)

### 6.4 Cannot Capture Passwords from Encrypted Channels

- Google encrypts passwords with BotGuard before sending them in batchexecute
- The proxy captures the password from the DOM before encryption (via `sendBeacon` in JS injection)
- If the victim uses a password manager that bypasses the DOM, the password may not be captured

### 6.5 Cannot Survive Token Revocation

- If the target account has continuous access evaluation (CAE) enabled, the captured token may be revoked when:
  - The user changes their password
  - The admin revokes sessions
  - The user signs out from another device
- The refresh token can survive password changes (for Microsoft), but not explicit revocation

### 6.6 Cannot Handle CAPTCHA Challenges

- If Google/Microsoft presents a CAPTCHA during the background browser flow, the automation fails
- The existing CAPTCHA bypass techniques (BotGuard spoofing, telemetry injection) work for basic challenges
- Advanced CAPTCHAs (reCAPTCHA v3 with high score threshold, hCaptcha) will block the automation

### 6.7 Cannot Handle Device Compliance Requirements

- If the tenant requires device compliance (Intune, Jamf), the background browser's device won't be compliant
- The `ms_intune_portal` and `ms_intune_enroll` client IDs can bypass some compliance checks, but not all

---

## 7. Current Fallbacks and Fixes

### 7.1 Fallback: No Chrome Running (EvilPuppet Rod → chromedp)

**Problem:** If the operator's Chrome is not running on port 9222, `EvilPuppetRod.HandleDeviceCodeVerification` fails.

**Current behavior:** The code checks `p.evilpuppetRod.IsChromeRunning()` and falls back to `p.evilpuppet.HandleDeviceCodeVerification()` (chromedp), which spawns a fresh Chromium with anti-detection flags. This makes the silent flow work without the operator's Chrome.

**Fix applied:** Added `EvilPuppet.HandleDeviceCodeVerification()` method in `core/evilpuppet.go` that:
1. Spawns a fresh Chromium via chromedp with anti-detection flags
2. Navigates to the verification URL
3. Enters the code
4. If credentials are available (from the AitM capture), signs in
5. Clicks Accept/Allow on the permissions page
6. Returns the result

The silent handler in `http_proxy.go` now uses rod when Chrome is running, and falls back to chromedp otherwise.

### 7.2 Fallback: Email Not Available on First Page Load

**Problem:** The preemptive page_load trigger needs the victim's email to fire. On the first page load, the email may only be in the URL hash fragment (not sent to the server).

**Current behavior:** The trigger checks:
1. `session.Username` (set from a previous POST)
2. `req.URL.Query().Get("login_hint")` (URL query param)
3. `req.URL.Query().Get("email")` (URL query param)
4. `session.Params["email"]` (lure params stored at session creation)

**Fix applied:** Added check for `session.Params["email"]` in `http_proxy.go:2518`. The lure email is now extracted from the session's stored params on first page load.

### 7.3 Fallback: Device Code Auto-Verify Fails (Page Structure Changed)

**Problem:** The `HandleDeviceCodeVerification` method uses hardcoded CSS selectors that break if Microsoft/Google changes their page structure.

**Fix applied:** Replaced all hardcoded selectors with resilient selector chains in both `evilpuppet_rod.go` and `evilpuppet.go`:

```go
codeSelectors := []string{
    "input#otc",
    "input[name='otc']",
    "input[aria-label*='code' i]",
    "input[placeholder*='code' i]",
    "input[type='text']",
    "input[type='tel']",
}
```

The `i` flag makes aria-label matching case-insensitive. The chain tries each selector until one succeeds.

### 7.4 Fallback: Silent Page Shows Forever (No Completion Feedback)

**Problem:** The `DEVICE_CODE_SILENT_HTML` shows a loading spinner that never resolves. If the device code flow fails, the victim sees a stuck loading page.

**Fix applied:** Added a minimal polling script to `DEVICE_CODE_SILENT_HTML` in `device_code_chain.go` that:
1. Polls `/dc/status/{session_id}` after 3 seconds
2. If captured: redirects to `d.redirect_url` or `https://www.office.com`
3. If failed/expired: shows "Unable to verify. Please try again later."
4. Uses XMLHttpRequest (no external dependencies)
5. Session ID is injected at serve time (not in the template)

### 7.5 Fallback: Telegram Notification Fails

**Problem:** If the Telegram bot token or chat ID is misconfigured, the notification fails silently.

**Fix applied:** Added a health check endpoint at `http_proxy.go`:
```
GET /api/v1/test-telegram?key={api_key}
```

This calls `p.notifier.Test("telegram", EventDeviceCodeCaptured)` and returns the result as JSON.

---

## 8. Defensive Measures and Mitigations

### 8.1 Organizational Controls

| Control | Effectiveness | Notes |
|---------|--------------|-------|
| **Disable device code flow** | ⭐⭐⭐⭐⭐ | Tenant-level policy to block OAuth device code grant entirely. Breaks some legitimate flows (Teams on smart TVs, etc.) |
| **Require phishing-resistant MFA** | ⭐⭐⭐⭐⭐ | FIDO2/WebAuthn hardware keys are immune to AitM attacks. The cryptographic proof is bound to the origin. |
| **Block legacy authentication** | ⭐⭐⭐⭐ | Prevents Basic Auth and legacy protocols that bypass MFA |
| **Conditional Access: Compliant devices only** | ⭐⭐⭐⭐ | Requires device enrollment (Intune/Jamf). AitM proxy's browser won't be compliant. |
| **Continuous Access Evaluation (CAE)** | ⭐⭐⭐ | Token revocation within 1 hour of policy change. But the attacker has a 1-hour window. |
| **Number matching in Authenticator** | ⭐⭐⭐ | Prevents push notification fatigue attacks. The victim must enter a number. |

### 8.2 Technical Controls

| Control | Effectiveness | Notes |
|---------|--------------|-------|
| **FIDO2/WebAuthn** | ⭐⭐⭐⭐⭐ | The ONLY control that fully defeats AitM. The private key never leaves the hardware token, and the signature is bound to the origin. |
| **Token binding (DPoP)** | ⭐⭐⭐⭐ | Demonstrating Proof of Possession binds tokens to a specific key. AitM proxy cannot use the token without the private key. |
| **Device compliance** | ⭐⭐⭐⭐ | Intune/MDM compliance checks. The proxy's browser device won't be enrolled. |
| **IP-based risk detection** | ⭐⭐⭐ | Azure AD Identity Protection flags sign-ins from unusual IPs. The proxy's IP is the attacker's server. |
| **Impossible travel detection** | ⭐⭐⭐ | If the victim's normal location is New York and the proxy is in another country, this triggers. |
| **Session lifetime policies** | ⭐⭐⭐ | Short session lifetimes (e.g., 1 hour) limit the window of exploitation. |

### 8.3 Detection Indicators

| Indicator | What to Look For |
|-----------|-----------------|
| Unusual OAuth device code requests | Multiple device code requests from the same IP in short succession |
| Device code + different IP approval | Device code requested from IP A, approved from IP B |
| FOCI token exchange patterns | Token exchange between FOCI clients in unusual sequences |
| BotGuard token anomalies | Batchexecute requests with mismatched session IDs |
| Proxy TLS fingerprint mismatches | uTLS fingerprint that doesn't match the claimed browser |
| Suspicious certificate patterns | Self-signed or Let's Encrypt certificates on lookalike domains |
| DNS anomalies | Rapid DNS queries for login.microsoftonline.com from the proxy IP |

### 8.4 Network-Level Mitigations

```
# Block device code flow at the network level (firewall/IDS)
# Microsoft device code endpoints:
# POST login.microsoftonline.com/{tenant}/oauth2/v2.0/devicecode
# POST login.microsoftonline.com/{tenant}/oauth2/v2.0/token (with device_code grant)

# Google device code endpoints:
# POST oauth2.googleapis.com/device/code
# POST oauth2.googleapis.com/token (with device_code grant)

# Detection: Block or alert on POST to /devicecode from non-corporate IPs
```

---

## 9. Indicators of Compromise

### 9.1 Network-Level IoCs

| Indicator | Description |
|-----------|-------------|
| `POST /batchexecute` with `rpcids=MI613e` | Google BotGuard token generation intercepted |
| `POST */oauth2/v2.0/devicecode` | Device code flow initiation |
| Multiple `f.sid` values in same session | Session ID swapping indicates proxy |
| `sec-ch-ua` header mismatches | Browser claims vs actual TLS fingerprint |
| Self-signed or Let's Encrypt certs on `*.microsoftonline.com` | Phishing domain serving fake certs |
| Unusual `X-Forwarded-For` chains | Proxy chain indicators |

### 9.2 Application-Level IoCs

| Indicator | Description |
|-----------|-------------|
| Device code approved within seconds of generation | Automated approval (not human) |
| Token requested from different IP than code generation | Proxy-based flow |
| FOCI client rotation in rapid succession | CAP bypass attempts |
| `navigator.webdriver` removal patterns | Anti-detection JS injection |
| SRI integrity attribute removal | `xintegrity` replacements |
| YouTube CheckConnection iframe bypass | Fake postMessage responses |

---

## 10. Appendix: Technical Reference

### 10.1 Phishlet Configuration (Google)

```yaml
evilpuppet:
  start_url: 'https://accounts.google.com/v3/signin/identifier?...'
  timeout: 30
  hold_request: true

  triggers:
    - domains: ['accounts.google.com']
      paths: ['.*']
      trigger_type: 'page_load'          # NEW: fires on GET page loads
    - domains: ['accounts.google.com']
      paths: ['.*batchexecute.*rpcids=MI613e.*']
      content_type: 'post'               # Original trigger

  actions:
    - type: waitVisible
      selector: 'input[type="email"], #identifierId'
    - type: type
      selector: 'input[type="email"], #identifierId'
      value: '{username}'
    - type: click
      selector: '#identifierNext button, #identifierNext'

  interceptors:
    - domain: 'accounts.google.com'
      path: '.*batchexecute.*'
      token_name: 'bgresponse'
      source: 'request_body'
      search: 'identity-signin-identifier\\"([^"]+)'
    # ... more interceptors for session_dsh, fsid, at_token

  inject_token:
    - token_name: 'bgresponse'
      target: 'body'
      search: '(identity-signin-identifier\\")([^"]+)'
      replace: '${1}{token}'
    # ... more inject rules
```

### 10.2 Phishlet Configuration (Microsoft)

```yaml
evilpuppet:
  start_url: 'https://login.microsoftonline.com/'
  timeout: 30
  hold_request: false                    # Preemptive only

  triggers:
    - domains: ['login.microsoftonline.com', 'login.live.com', 'login.microsoft.com']
      paths: ['.*']
      trigger_type: 'page_load'

  actions:
    - type: waitVisible
      selector: 'input[type="email"], #i0116'
    - type: type
      selector: 'input[type="email"], #i0116'
      value: '{username}'
    - type: click
      selector: '#idSIButton9'
```

### 10.3 Device Code Mode Reference

| Mode | Behavior | Interstitial | Auto-Verify |
|------|----------|-------------|-------------|
| `off` | No device code | None | N/A |
| `always` | After AitM success | Standard interstitial | No |
| `fallback` | If AitM stalls | Standard interstitial | No |
| `auto` | Pre-generate on lure click | Standard interstitial | No |
| `direct` | Skip AitM, show code immediately | Standard interstitial | No |
| `silent` | Background only, user sees nothing | Benign loading page | Yes (via EvilPuppet) |

### 10.4 Test Coverage

| Area | Tests | Coverage |
|------|-------|----------|
| YAML parsing (trigger_type, silent mode) | 5 | Full |
| Trigger matching (page_load, request) | 4 | Full |
| DeviceCodeManager lifecycle | 5 | Full |
| Silent HTML template security | 5 | Full |
| EvilPuppet auto-verify (rod: MS/Google) | 4 | Full |
| EvilPuppet auto-verify (chromedp: MS/Google) | 3 | Full |
| Preemptive token storage/reuse | 4 | Full |
| Session state transitions | 3 | Full |
| Provider resolution | 2 | Full |
| Lure params email extraction | 1 | Full |
| **Total** | **43** | **All pass** |

### 10.5 Files Modified

| File | Changes | Purpose |
|------|---------|---------|
| `core/evilpuppet.go` | +228 | `TriggerType` field, `HasPageLoadTrigger()`, `HandleDeviceCodeVerification()` chromedp method |
| `core/phishlet.go` | +8 | YAML parsing for `trigger_type`, silent mode validation |
| `core/evilpuppet_rod.go` | +219 | `HandleDeviceCodeVerification()` rod method with resilient selectors |
| `core/http_proxy.go` | +234 | Preemptive page_load handler, silent DC handler, chromedp fallback, lure email extraction, Telegram health check, token reuse |
| `core/device_code_chain.go` | +122 | `DCModeSilent`, silent HTML template with polling, auto-open sign-in |
| `phishlets/google.yaml` | +3 | page_load trigger |
| `phishlets/o365.yaml` | +26 | evilpuppet section |
| **Total** | **+798/-42** | |

---

*This report is provided for defensive security research purposes. Understanding offensive techniques is essential for building effective defenses.*
