## Why

Developers and QA teams need a reproducible way to create, operate, and restore multiple Android test environments on their own Windows computers without sending application or device data to a hosted execution service. This change establishes the local control plane, virtual-device capabilities, and minimal cloud entitlement boundary needed to deliver that product safely and incrementally.

## What Changes

- Add a Windows desktop Agent with a tray entry point and a browser-based localhost user interface.
- Add local creation, startup, shutdown, monitoring, and deletion of Android virtual devices.
- Add low-latency device screen streaming and browser-originated touch and keyboard control.
- Add coherent, reproducible QA device profiles plus proxy and simulated-location configuration.
- Add a Main Device selection model and bulk command execution across selected devices with per-device outcomes.
- Add local device snapshot and restore, including compatibility metadata and integrity checks.
- Add account authentication, subscription validation, signed entitlements, activation limits, and bounded offline use through a minimal cloud service.
- Keep APKs, device contents, snapshots, proxy credentials, and testing activity on the user's computer.

## Capabilities

### New Capabilities

- `local-agent-control-plane`: Windows Agent lifecycle, embedded localhost UI, authenticated local API, persistent local state, and secure component boundaries.
- `virtual-device-lifecycle`: Provisioning, starting, stopping, observing, and removing Android virtual devices with resource and compatibility checks.
- `device-streaming-control`: Live device display delivery and browser-based touch and keyboard input.
- `device-environment-configuration`: Reproducible QA profiles, HTTP and SOCKS5 proxy configuration, and simulated GPS location controls.
- `multi-device-orchestration`: Main Device designation, device selection, bulk commands, and per-device execution results.
- `device-snapshot-restore`: Consistent local snapshots, restore workflows, manifests, compatibility validation, and retention behavior.
- `license-entitlements`: Cloud account and subscription validation, device activation, signed feature entitlements, and offline grace behavior.

### Modified Capabilities

None. This repository has no existing capability specifications.

## Impact

- Introduces a Go-based Windows Agent, a Next.js App Router frontend packaged for local serving, Android SDK/Emulator integration, ADB-based control, and local SQLite persistence.
- Introduces real-time control and event protocols over localhost WebSocket plus a dedicated encoded-media transport.
- Introduces a minimal cloud API and billing-provider integration for accounts, subscriptions, activations, and entitlements.
- Requires Windows virtualization and graphics compatibility checks, Android runtime distribution/licensing validation, secure updater design, and performance testing across supported hardware.
- Establishes privacy and security boundaries between browser, Agent, emulated devices, privileged helper, and cloud services.
