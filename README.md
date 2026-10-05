# Device Force

Device Force is a local-first Android application testing platform for Windows. A per-user Go Agent manages Android Emulator instances and serves a statically exported Next.js interface on loopback. The cloud boundary is limited to accounts, subscriptions, installation activations, and signed entitlements.

## Workspace

- `cmd/agent`: local Windows Agent executable
- `cmd/cloud`: minimal account and entitlement service
- `cmd/tray`: per-user Windows tray entry point
- `internal`: Agent and cloud implementation packages
- `apps/web`: local Next.js App Router interface
- `contracts/v1`: versioned HTTP, event, command, and signaling contracts
- `migrations`: local SQLite migrations
- `packaging`: Windows packaging and update assets
- `docs`: decisions, security notes, support matrix, and operator guidance
- `tests`: cross-component and platform test fixtures

## Development

Requirements: Go 1.27+, Node.js 24 LTS+, and npm 11+.

```sh
npm ci
npm test
npm run build
go test ./cmd/... ./internal/...
go vet ./cmd/... ./internal/...
```

The Android Emulator execution path is Windows-only. Portable unit tests and static builds run on macOS, Linux, and Windows; WHPX, tray, DPAPI, installer, and emulator acceptance tests run on Windows workers.
