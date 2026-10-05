# Product and Runtime Baseline v1

Status: Approved for implementation

Last reviewed: 2026-10-05

This baseline turns the initial product assumptions into versioned implementation inputs. Changes require a new baseline revision and compatibility review.

## Android component installation and licensing

The installer SHALL NOT redistribute Google Android SDK, Emulator, platform-tools, or system-image packages.

The Agent will:

1. Discover an existing user-selected Android SDK, or install Google's command-line tools after an explicit user action.
2. Show the exact package IDs, versions, download source, disk estimate, and applicable license before installation.
3. Invoke `sdkmanager` to download packages directly from Google's repositories into an Agent-managed SDK directory.
4. Require the user to accept each applicable SDK package license. The Agent will record the accepted package ID, license hash, and acceptance time, but will not pre-accept licenses on the user's behalf.
5. Verify installed package revisions and file presence before using them.

This model follows the Android SDK License Agreement's restrictions on copying and redistribution while using Google's supported `sdkmanager` installation and license flow:

- https://developer.android.com/studio/terms
- https://developer.android.com/tools/sdkmanager

The initial media/control adapter may bundle a pinned, unmodified scrcpy server artifact or build it from pinned source. scrcpy is Apache-2.0 licensed. Any distribution SHALL include its license, required notices, exact upstream version/commit, locally applied patches, and artifact checksums:

- https://github.com/Genymobile/scrcpy
- https://github.com/Genymobile/scrcpy/blob/master/LICENSE

This engineering decision is not a substitute for final legal review before public distribution.

## Supported Windows host

Minimum supported host:

- Windows 11 25H2, 64-bit x86; Windows on ARM is not supported.
- Current security updates installed.
- Windows Hypervisor Platform enabled and usable.
- Intel VT-x with EPT or AMD-V with RVI enabled in firmware.
- Intel 8th-generation Core or AMD Zen 1 class CPU or newer.
- 16 GB RAM for one active device.
- 40 GB free SSD space for the product, SDK, and one device; each additional device or retained snapshot requires its estimated space to be available before creation.
- DirectX 11-capable GPU with current driver and at least 4 GB dedicated or shared graphics memory.
- Current or previous stable Microsoft Edge or Google Chrome major release.
- Display resolution of at least 1280 x 800.

Recommended for four active devices:

- Modern 8-core Intel Core or AMD Ryzen CPU.
- 32 GB RAM or more.
- 8 GB VRAM and 100 GB or more free NVMe SSD space.
- 1920 x 1080 display or larger.

Nested virtualization is unsupported. The product performs resource preflight and may set a lower concurrency limit than the subscription allows.

References:

- https://developer.android.com/studio/install
- https://developer.android.com/studio/run/emulator-acceleration
- https://learn.microsoft.com/windows/release-health/supported-versions-windows-client

## Android image support

The v1 compatibility matrix is deliberately finite:

| Android | API | ABI | Image families |
|---|---:|---|---|
| 11 | 30 | x86_64 | AOSP, Google APIs |
| 13 | 33 | x86_64 | AOSP, Google APIs |
| 15 | 35 | x86_64 | AOSP, Google APIs |
| 16 | 36 | x86_64 | AOSP, Google APIs |

Only stable-channel image revisions explicitly present in the product compatibility catalog are supported. Preview, beta, Automotive, TV, Wear, and Google Play Store images are excluded from v1.

Google Play Store images are deferred because they add a distinct license and compatibility surface, are release-key signed, and do not permit the same troubleshooting privileges as AOSP images. Google APIs images provide Play services but not the Play Store. The UI must show the distinction before provisioning.

## SOCKS5 behavior

The v1 per-device SOCKS5 feature provides:

- TCP `CONNECT` routing through a managed Android `VpnService` companion.
- DNS resolution through the proxy path when the configured SOCKS5 server supports remote DNS.
- Per-device isolation and explicit connectivity verification.

The following are unsupported in v1 and must be shown before applying the configuration:

- UDP relay and UDP-based application protocols, including QUIC/HTTP3.
- ICMP forwarding.
- Transparent TLS interception or certificate installation.
- Per-domain bypass rules.

Applications must be tested over TCP while SOCKS5 mode is active. Unsupported traffic is blocked rather than silently bypassing the configured proxy.

## Entitlement defaults

Entitlement values are server-issued configuration, not hard-coded UI behavior:

| Plan | Active installations | Concurrent devices per installation | Entitlement lifetime | Offline grace |
|---|---:|---:|---:|---:|
| Starter | 1 | 1 | 24 hours | 72 hours |
| Professional | 2 | 4 | 24 hours | 7 days |
| Team | 2 per assigned seat | 8 | 24 hours | 7 days |

An online refresh is attempted at startup, at half the entitlement lifetime, and with bounded backoff after failures. Offline grace begins at the last successful validation time, not at token expiry. Plan configuration may reduce limits immediately after online validation but cannot retroactively shorten an already signed offline deadline.

After grace expires, creation, start, restore, environment mutation, and streaming are blocked. Stop, force-stop, diagnostics, local export, snapshot export, and deletion remain available.

## Deferred capabilities

The following are explicitly outside v1:

- Moving or restoring snapshots on a different host.
- Live mirroring of arbitrary Main Device taps, gestures, and keyboard input to other devices.

The Main Device is a command-routing and UI-selection concept in v1. Snapshots are restorable only on their originating host within the recorded runtime compatibility constraints.
