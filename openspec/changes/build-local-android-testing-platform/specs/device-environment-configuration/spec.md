## ADDED Requirements

### Requirement: Reproducible QA profiles
The system SHALL generate or select versioned device profiles whose supported fields form a coherent configuration for the chosen Android image and whose generated values are reproducible from a recorded seed.

#### Scenario: User randomizes a device profile
- **WHEN** the user requests a randomized identity for a compatible device
- **THEN** the system SHALL produce a coherent profile, record its seed and schema version, and report every field as applied, deferred, unsupported, or failed

### Requirement: Honest identity boundaries
The system SHALL not claim to modify hardware-backed identifiers, integrity verdicts, or fields unsupported by the selected runtime.

#### Scenario: Requested field is not supported
- **WHEN** a profile contains a field that the current runtime cannot apply
- **THEN** the system SHALL mark the field unsupported and SHALL not report the profile as fully applied

### Requirement: Per-device proxy configuration
The system SHALL support per-device HTTP proxy settings and SHALL support SOCKS5 through an explicitly identified compatible routing mechanism.

#### Scenario: Proxy configuration is applied
- **WHEN** the user supplies a valid proxy endpoint and optional credentials
- **THEN** the system SHALL validate the settings, apply them only to the targeted device, verify connectivity where possible, and report protocol limitations

### Requirement: Proxy secret protection
The system SHALL protect stored proxy credentials with Windows-backed secret protection and SHALL redact them from logs, events, exports, and snapshot manifests.

#### Scenario: Diagnostic bundle is created
- **WHEN** the user exports diagnostics for a device with an authenticated proxy
- **THEN** the bundle SHALL identify the proxy mode without containing the proxy password or reusable credential

### Requirement: Simulated location controls
The system SHALL support fixed coordinates and route playback with explicit speed and update interval for compatible devices.

#### Scenario: Route playback is stopped
- **WHEN** the user stops an active simulated route
- **THEN** the system SHALL stop sending route updates and SHALL report the final simulated location
