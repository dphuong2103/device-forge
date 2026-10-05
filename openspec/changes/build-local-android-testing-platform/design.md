## Context

The repository is greenfield. The product runs Android test environments on a customer's Windows computer, exposes their controls through a local browser interface, and uses the cloud only for accounts and subscription entitlements. The design must accommodate hardware-dependent virtualization, GPU-backed emulator processes, multiple concurrent devices, low-latency media, sensitive local test data, and temporary loss of cloud connectivity.

The primary stakeholders are developers and QA engineers using the desktop product, administrators managing subscriptions, and engineers operating the minimal entitlement service. The first release targets a single interactive Windows user and a Chromium-based browser.

## Goals / Non-Goals

**Goals:**

- Deliver a secure local control plane for reproducible Android application testing.
- Support one device first without preventing multi-device operation, bulk commands, or snapshots later.
- Keep application binaries, device contents, test activity, proxy secrets, and snapshots local.
- Make device and bulk operations observable, cancellable where safe, and resilient to Agent or emulator restarts.
- Define clear compatibility and entitlement contracts so failures are explainable to users.

**Non-Goals:**

- Hosting or remotely operating customers' Android devices in the product cloud.
- Modifying hardware-backed identifiers, bypassing application integrity controls, or impersonating unsupported physical hardware.
- Guaranteeing snapshot portability across arbitrary hosts, system images, emulator versions, or graphics backends.
- Synchronously mirroring arbitrary taps and keystrokes from the Main Device to every other device in the initial implementation.
- Supporting nested virtualization, non-Windows hosts, or every browser in the first release.

## Decisions

### Use a per-user Agent with a small optional privileged helper

The Go Agent and emulator processes run in the signed-in user's Windows session so GPU acceleration, user-owned files, browser launch, and tray interactions work normally. A tray application starts or discovers the Agent and opens its UI. An optional narrowly scoped Windows service performs only operations that require elevation, such as installation, update replacement, or host feature checks.

A service-only design was rejected because Windows services run in Session 0 and create avoidable UI, GPU, and ownership problems. A single fully elevated desktop process was rejected because it unnecessarily expands the privileged attack surface.

### Serve a static Next.js application from the Agent

The Next.js App Router frontend is built as a static export and embedded or installed beside the Agent. The Agent serves the UI and versioned local API from one loopback origin. Server-rendered Next.js and a bundled Node.js runtime are excluded until a concrete server-side UI need appears.

The tray launches a one-time bootstrap URL. The Agent exchanges that token for a short-lived, HttpOnly, SameSite session and then invalidates the token. Every state-changing request and WebSocket upgrade validates the origin, host, session, and anti-CSRF data. The listener binds only to IPv4 and IPv6 loopback.

### Manage official Android Emulator AVDs through adapters

The runtime adapter provisions Android Virtual Devices, launches one emulator process per device, assigns ports and ADB serials, detects boot completion, and records the exact system-image and emulator versions. Windows Hypervisor Platform is the preferred accelerator. A preflight service verifies OS, virtualization, storage, memory, GPU mode, SDK components, and port availability before provisioning.

Runtime-specific behavior stays behind interfaces for lifecycle, console control, ADB, streaming, and snapshots. This permits isolated testing and future runtime changes without leaking command-line details into the API.

### Serialize mutations through per-device actors

Each device has a single command actor and explicit lifecycle state. Commands have IDs, target versions, deadlines, and durable status. Read-only observations can run concurrently, while lifecycle, configuration, and snapshot mutations are serialized. On restart, the Agent reconciles database state with actual processes and ADB-visible devices before accepting new mutations.

Bulk operations fan out into individual device commands. The aggregate result retains success, failure, timeout, cancellation, and error detail per target; it does not pretend the group is transactional.

### Separate control events from encoded media

HTTP handles resource-oriented commands and queries. WebSocket carries device state, command progress, logs, and low-volume interactive events. The media gateway receives encoded frames from the Android side and uses WebRTC for browser delivery, congestion handling, and recovery. Input returns through the authenticated control channel and is translated to device coordinates before injection.

An initial implementation may use a scrcpy-compatible server/protocol for capture and input. The integration remains isolated behind a media adapter so upstream protocol changes do not affect the public local API. Raw media over WebSocket is retained only as a diagnostic or constrained fallback.

### Treat test profiles as coherent and reproducible configurations

A device profile is a versioned document generated from a seed and a supported catalog. It combines only settings the selected runtime can reproduce: compatible model properties, display, API level, ABI, locale, timezone, network configuration, and initial location. Applying immutable boot properties creates or cold-boots an AVD; mutable properties can be changed at runtime.

The system reports each field as applied, deferred, unsupported, or failed. It does not claim to alter hardware-backed identifiers or application integrity verdicts.

### Implement proxy modes according to runtime capability

HTTP proxy configuration uses the emulator's supported launch/runtime mechanisms. Per-device SOCKS5 and advanced routing use a managed Android VPN companion where feasible, with credentials delivered through a protected local channel and never written to general logs. HTTPS interception and certificate installation are separate, explicit testing features and are not implied by proxy configuration.

### Maintain fast snapshots and compatibility-described backups

Fast snapshots use native emulator snapshot facilities. Managed backups quiesce or stop the device and persist the required AVD state with a manifest containing checksums, runtime version, system-image identity, profile version, configuration, and creation status. Restore is staged into a temporary target, validated, and atomically promoted only after integrity and compatibility checks pass.

Snapshots are local and compatibility-bound by default. Cross-computer portability is not promised until validated by a dedicated compatibility matrix.

### Keep operational data local and minimize cloud records

SQLite stores devices, profiles, command history, snapshot manifests, local settings, and cached signed entitlements. Large artifacts remain in managed filesystem directories. Secrets use Windows-protected credential storage or DPAPI-backed encryption rather than plaintext database columns.

The cloud stores account, subscription, plan, activation, entitlement issuance, and security audit metadata. It does not receive APKs, app data, snapshots, proxy credentials, GPS history, screen content, or device interaction logs.

### Use signed, expiring entitlements with bounded offline operation

After authentication, the Agent activates the installation and receives a signed entitlement containing account/installation identifiers, plan, feature limits, issue and expiry times, and offline policy. The Agent verifies it locally and refreshes it over outbound HTTPS. During a cloud outage, the last valid entitlement remains usable only within its signed grace policy. Clock rollback and corrupted entitlement state fail closed while preserving access to local export and cleanup operations.

## Risks / Trade-offs

- [Android runtime redistribution or system-image terms prevent the preferred packaging] -> Validate distribution terms before bundling; default to an on-device SDK component installer with recorded licenses.
- [Concurrent emulators exhaust RAM, CPU, GPU, or disk] -> Add preflight estimates, configurable concurrency limits, lower-quality background streams, quotas, and clear resource diagnostics.
- [Media integration depends on an unstable upstream protocol] -> Pin and test a known version behind an adapter; retain a replacement path and comply with upstream licensing.
- [Localhost APIs can be attacked by malicious websites] -> Use loopback binding, one-origin serving, one-time bootstrap, authenticated sessions, strict origin/host validation, CSRF protection, and request limits.
- [SOCKS5 behavior differs across Android versions and traffic types] -> Define supported TCP/UDP/DNS behavior explicitly and test the VPN companion against every supported image.
- [Snapshot restore appears successful but produces corrupted state] -> Quiesce devices, checksum artifacts, restore to a staging location, retain the prior state until validation, and surface compatibility failures before mutation.
- [Identity profiles imply fidelity the emulator cannot provide] -> Maintain an allowlisted field catalog and expose per-field application status instead of silently spoofing values.
- [Offline licensing blocks legitimate work during outages] -> Use a visible signed grace window and preserve non-entitled access to export, diagnostics, and deletion.

## Migration Plan

This is a new system with no user-data migration. Delivery proceeds behind capability flags: first install the Agent and single-device runtime, then enable media control, environment settings, multi-device commands, snapshots, and paid entitlement enforcement. Local database migrations are versioned and backed up before upgrade. Agent updates retain the previous signed binary and schema-compatible database backup for rollback. Cloud entitlement enforcement begins in audit-only mode before it becomes blocking.

## Open Questions

- Which Android system-image families, including Google APIs or Google Play images, are supported and distributable?
- What minimum Windows version and hardware profile define the support matrix and default concurrency limits?
- Does the first SOCKS5 release require UDP and DNS routing, or only application TCP traffic?
- What signed offline-grace duration and installation activation limit apply to each subscription plan?
- Is cross-host snapshot transfer a future requirement, and if so, which host differences must be supported?
- Is live input mirroring from the Main Device required after command-based bulk actions ship?
