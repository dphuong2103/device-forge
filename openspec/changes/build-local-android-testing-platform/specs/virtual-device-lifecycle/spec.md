## ADDED Requirements

### Requirement: Host preflight
The system SHALL validate supported Windows version, virtualization acceleration, SDK/runtime components, storage, memory, graphics mode, and port availability before creating or starting a virtual device.

#### Scenario: Virtualization is unavailable
- **WHEN** a user attempts to create or start a device without usable hardware acceleration
- **THEN** the system SHALL block the operation and SHALL report the failed check and an actionable remediation

### Requirement: Reproducible device provisioning
The system SHALL create each virtual device from a recorded system-image identity, runtime version, hardware profile, and configuration version.

#### Scenario: Device creation succeeds
- **WHEN** all requested components are available and compatible
- **THEN** the system SHALL create the device, persist its reproducibility metadata, and leave it in a known stopped state

### Requirement: Managed lifecycle
The system SHALL support starting, observing, gracefully stopping, force-stopping, and deleting managed virtual devices through explicit lifecycle states.

#### Scenario: Device reaches readiness
- **WHEN** the emulator process starts and Android completes boot
- **THEN** the system SHALL mark the device ready only after its control and observation channels pass health checks

#### Scenario: Startup fails
- **WHEN** the process exits, boot times out, or health checks fail
- **THEN** the system SHALL mark the device failed, release reserved resources, and retain diagnostic information

### Requirement: Destructive deletion protection
The system SHALL require explicit confirmation before deleting a device and SHALL identify snapshots or local artifacts affected by the deletion.

#### Scenario: User cancels deletion
- **WHEN** the user declines the deletion confirmation
- **THEN** the system SHALL preserve the device and all associated artifacts unchanged

### Requirement: Resource-aware concurrency
The system SHALL enforce configured device concurrency and host resource limits.

#### Scenario: Starting another device exceeds a limit
- **WHEN** the requested start would exceed an applicable entitlement or host resource limit
- **THEN** the system SHALL leave the device stopped and SHALL identify the limiting resource or entitlement
