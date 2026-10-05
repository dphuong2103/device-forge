## ADDED Requirements

### Requirement: Consistent snapshot creation
The system SHALL quiesce or stop a device as required by the selected snapshot method before declaring a snapshot complete.

#### Scenario: Snapshot succeeds
- **WHEN** all device state is captured and its manifest and checksums are durably stored
- **THEN** the system SHALL mark the snapshot complete and SHALL make it available for restore

#### Scenario: Snapshot capture fails
- **WHEN** any required state or manifest write fails
- **THEN** the system SHALL mark the snapshot incomplete, SHALL not offer it for normal restore, and SHALL preserve the source device

### Requirement: Snapshot contents and manifest
The system SHALL capture identity configuration, installed applications, application data, Android settings, and user data supported by the runtime, and SHALL record runtime, system-image, profile, configuration, and checksum metadata in a manifest.

#### Scenario: User inspects a snapshot
- **WHEN** the user opens snapshot details
- **THEN** the system SHALL display its creation status, creation time, source device, size, compatibility metadata, and included state categories

### Requirement: Compatibility validation
The system SHALL validate snapshot integrity and host/runtime compatibility before altering a restore target.

#### Scenario: Snapshot is incompatible
- **WHEN** required runtime or system-image compatibility conditions are not met
- **THEN** the system SHALL block restore and SHALL identify the unmet conditions

### Requirement: Atomic restore
The system SHALL stage restored data separately and SHALL preserve the previous target state until restored state passes integrity and startup validation.

#### Scenario: Restored device fails validation
- **WHEN** staged restore data cannot boot or pass required health checks
- **THEN** the system SHALL retain or reinstate the pre-restore target state and SHALL report the failure

### Requirement: Explicit snapshot deletion
The system SHALL require confirmation before deleting a snapshot and SHALL make deletion independent from deletion of the source device unless the user explicitly selects both.

#### Scenario: Source device is deleted
- **WHEN** retained snapshots exist and the user confirms deletion of only the source device
- **THEN** the system SHALL preserve those snapshots and their manifests
