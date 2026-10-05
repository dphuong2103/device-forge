## 1. Resolve Product and Runtime Baselines

- [x] 1.1 Validate Android Emulator, system-image, and scrcpy-compatible component distribution terms and record the approved installation model
- [x] 1.2 Define and document the minimum Windows, CPU virtualization, RAM, GPU, disk, and Chromium support matrix
- [x] 1.3 Decide the supported Android API levels and AOSP, Google APIs, or Google Play image families
- [x] 1.4 Define SOCKS5 TCP, UDP, and DNS behavior for the first supported release
- [x] 1.5 Define plan-specific activation limits, concurrent-device limits, entitlement lifetimes, and offline-grace durations
- [x] 1.6 Record cross-host snapshot transfer and live input mirroring as deferred capabilities unless their requirements are approved

## 2. Establish Repository and Contracts

- [x] 2.1 Create Go Agent, Next.js frontend, cloud service, shared contract, packaging, and test workspace structure
- [x] 2.2 Configure repeatable builds, formatting, linting, unit tests, dependency checks, and Windows CI artifacts
- [x] 2.3 Define versioned HTTP, WebSocket event, command, error, and WebRTC signaling contracts
- [x] 2.4 Define Agent interfaces for runtime lifecycle, ADB control, media, profiles, proxy, location, snapshots, secrets, and entitlements
- [x] 2.5 Define the SQLite schema and migrations for devices, commands, profiles, snapshots, settings, sessions, and cached entitlements
- [x] 2.6 Add structured logging, correlation IDs, diagnostic event types, and centralized secret-redaction tests

## 3. Build the Local Agent Control Plane

- [x] 3.1 Implement Agent startup, single-instance coordination, graceful shutdown, and loopback-only IPv4 and IPv6 listeners
- [x] 3.2 Build the Next.js static export and serve its versioned assets and API from the Agent origin
- [x] 3.3 Implement expiring single-use bootstrap tokens and HttpOnly SameSite browser sessions
- [x] 3.4 Enforce host, origin, session, CSRF, and request-size checks on HTTP and WebSocket entry points
- [x] 3.5 Implement Windows-backed encryption for refresh tokens, proxy credentials, and other local secrets
- [ ] 3.6 Implement the tray application's Agent discovery, start, authenticated browser-open, status, and exit actions
- [x] 3.7 Implement the narrowly scoped privileged-helper protocol or document that no MVP operation requires the helper
- [ ] 3.8 Implement startup reconciliation between SQLite records, emulator processes, reserved ports, and ADB-visible devices
- [ ] 3.9 Add control-plane tests for token replay, hostile origins, DNS rebinding headers, CSRF, restart recovery, and log redaction

## 4. Implement Virtual Device Lifecycle

- [ ] 4.1 Implement host preflight checks with actionable results for OS, WHPX, SDK components, storage, memory, graphics, and ports
- [ ] 4.2 Implement approved Android SDK component discovery and on-device installation with version and license records
- [ ] 4.3 Implement AVD provisioning from a versioned system image and hardware profile into Agent-managed directories
- [ ] 4.4 Implement the per-device actor, lifecycle state machine, durable command status, deadlines, and safe cancellation
- [ ] 4.5 Implement emulator launch with assigned ports, process supervision, boot-completion detection, ADB health checks, and diagnostics
- [ ] 4.6 Implement graceful stop, force-stop fallback, resource release, and failed-start cleanup
- [ ] 4.7 Implement confirmed device deletion while preserving snapshots unless separately selected
- [ ] 4.8 Implement host and entitlement concurrency limits with resource reservation and clear rejection reasons
- [ ] 4.9 Add adapter and Windows integration tests for create, start, ready, stop, crash, restart reconciliation, and deletion flows

## 5. Build the Local Device Management UI

- [ ] 5.1 Implement the application shell, Agent connection state, session-expiry handling, and accessible error presentation
- [ ] 5.2 Implement the device dashboard with lifecycle, health, resource, Main Device, and selection indicators
- [ ] 5.3 Implement the preflight and create-device workflow using supported images and hardware profiles
- [ ] 5.4 Implement start, stop, force-stop, and confirmed deletion actions with live command progress
- [ ] 5.5 Implement device detail and diagnostic views without exposing protected credentials
- [ ] 5.6 Add frontend contract tests and end-to-end tests for bootstrap, reconnect, device creation, lifecycle, and destructive confirmation

## 6. Add Live Streaming and Input

- [ ] 6.1 Implement and pin the Android-side media/control adapter with version and license metadata
- [ ] 6.2 Implement Agent-side encoded video ingestion, stream lifecycle, and cleanup independent of the public API
- [ ] 6.3 Implement authenticated WebRTC signaling and local media delivery for one focused device
- [ ] 6.4 Implement browser video rendering with loading, connected, reconnecting, failed, rotation, and resize states
- [ ] 6.5 Implement coordinate transformation and injection for pointer, touch, keyboard, text, and Android navigation input
- [ ] 6.6 Implement bounded stream reconnection without restarting a healthy device
- [ ] 6.7 Implement focused-stream quality policy and reduced frame rate, resolution, or bitrate for background streams
- [ ] 6.8 Add latency, input-accuracy, rotation, disconnect-recovery, and concurrent-stream performance tests

## 7. Add Device Environment Configuration

- [ ] 7.1 Define the versioned allowlist of supported profile fields and coherent catalog relationships for each supported image
- [ ] 7.2 Implement deterministic profile generation from a recorded seed and schema version
- [ ] 7.3 Implement profile validation and per-field applied, deferred, unsupported, and failed results
- [ ] 7.4 Implement cold-boot or reprovision handling for immutable profile fields and runtime application for mutable fields
- [ ] 7.5 Implement per-device HTTP proxy validation, application, connectivity checks, removal, and limitation reporting
- [ ] 7.6 Implement the approved per-device SOCKS5 routing companion and protected credential exchange
- [ ] 7.7 Implement fixed simulated location, saved locations, route playback, speed, interval, stop, and final-position reporting
- [ ] 7.8 Implement UI workflows for profile randomization, profile results, proxy settings, and location controls
- [ ] 7.9 Add tests for seed reproducibility, unsupported identity fields, secret redaction, proxy isolation, and location routes

## 8. Add Multi-Device Orchestration

- [ ] 8.1 Implement atomic Main Device assignment and persisted device selection
- [ ] 8.2 Implement explicit target resolution for Main Device, selected devices, and their union with empty-target rejection
- [ ] 8.3 Implement bulk command fan-out through per-device actors with independent progress, timeout, cancellation, and results
- [ ] 8.4 Implement aggregate partial-success reporting without rolling back successful device commands
- [ ] 8.5 Implement UI controls for Main Device, multi-selection, target preview, bulk confirmation, progress, and per-device outcomes
- [ ] 8.6 Add concurrency tests proving mutations serialize per device while different devices progress independently
- [ ] 8.7 Add tests proving ordinary Main Device input is not implicitly mirrored to other devices

## 9. Add Snapshot and Restore

- [ ] 9.1 Implement native fast-snapshot creation, listing, loading, failure cleanup, and runtime-version recording
- [ ] 9.2 Implement managed snapshot quiescence and capture of required AVD state into Agent-managed storage
- [ ] 9.3 Implement snapshot manifests with source, state categories, versions, profile data, size, status, and checksums
- [ ] 9.4 Implement snapshot integrity and runtime, image, host, and graphics compatibility validation before restore
- [ ] 9.5 Implement staged restore, boot and health validation, atomic promotion, and pre-restore rollback preservation
- [ ] 9.6 Implement independent confirmed deletion and retention behavior for devices and snapshots
- [ ] 9.7 Implement snapshot list, detail, compatibility, progress, restore, rollback, and deletion UI
- [ ] 9.8 Add tests for interrupted creation, checksum failure, incompatible restore, failed boot rollback, and source-device deletion

## 10. Implement Cloud Accounts and Entitlements

- [ ] 10.1 Implement the cloud account, installation, subscription, plan, activation, entitlement issuance, and audit data model
- [ ] 10.2 Integrate the selected identity and billing providers with idempotent subscription-state synchronization
- [ ] 10.3 Implement installation activation, plan-limit enforcement, activation listing, and revocation
- [ ] 10.4 Implement asymmetric signing and key rotation for audience-bound, installation-bound, expiring entitlements
- [ ] 10.5 Implement Agent sign-in, refresh credential protection, entitlement retrieval, signature verification, and feature-limit evaluation
- [ ] 10.6 Implement cached offline-grace evaluation, visible validation deadlines, clock-anomaly handling, and expiry restrictions
- [ ] 10.7 Preserve diagnostics, export, device stop, and deletion operations after entitlement expiry
- [ ] 10.8 Implement account, subscription, activation, limit, offline status, and remediation views in the local UI
- [ ] 10.9 Add contract and security tests for inactive subscriptions, activation limits, revocation, invalid signatures, key rotation, outages, expiry, and clock rollback
- [ ] 10.10 Verify through automated schemas and telemetry tests that prohibited local test data is neither accepted nor emitted by cloud licensing flows

## 11. Package, Harden, and Release

- [ ] 11.1 Build a signed Windows installer that installs the Agent, tray application, UI assets, and only approved dependencies
- [ ] 11.2 Implement signed updates with staged replacement, previous-version retention, database backup, and rollback
- [ ] 11.3 Run the supported Windows, CPU, GPU, Android image, proxy, and snapshot compatibility matrix and publish known limitations
- [ ] 11.4 Run local API threat tests covering hostile websites, session theft, privilege boundaries, filesystem access, and secret handling
- [ ] 11.5 Run single-device and multi-device endurance tests with CPU, RAM, GPU, disk, startup, stream-latency, and cleanup measurements
- [ ] 11.6 Run end-to-end acceptance tests mapping every specification scenario to an automated test or documented manual verification
- [ ] 11.7 Produce installation, virtualization setup, privacy boundary, proxy limitation, snapshot compatibility, diagnostics, and recovery documentation
- [ ] 11.8 Gate entitlement enforcement behind audit-only telemetry until activation, outage, grace, revocation, and rollback behavior pass release criteria
