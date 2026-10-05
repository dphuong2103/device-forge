## ADDED Requirements

### Requirement: Main Device designation
The system SHALL allow at most one managed device to be designated as the Main Device and SHALL persist that designation independently of device runtime state.

#### Scenario: User designates a new Main Device
- **WHEN** another managed device is already the Main Device
- **THEN** the system SHALL atomically remove the old designation and assign it to the selected device

### Requirement: Explicit command targets
The system SHALL require every multi-device command to resolve to an explicit target set consisting of the Main Device, selected devices, or both.

#### Scenario: Target selection resolves to no devices
- **WHEN** the user submits a bulk action with no resolvable targets
- **THEN** the system SHALL reject the action without changing any device

### Requirement: Per-device bulk outcomes
The system SHALL execute bulk actions as independently tracked per-device commands and SHALL preserve each target's result.

#### Scenario: One target fails
- **WHEN** a bulk action succeeds on some devices and fails on another
- **THEN** the aggregate result SHALL report partial completion and SHALL show the outcome and error for every target

### Requirement: Serialized device mutation
The system SHALL serialize mutating commands for each individual device while allowing different devices to make progress independently.

#### Scenario: Conflicting commands target one device
- **WHEN** a device is already executing an incompatible mutation
- **THEN** the later command SHALL be queued or rejected according to command policy and SHALL not interleave the mutations

### Requirement: No implicit live input mirroring
The system SHALL treat Main Device designation and bulk actions as command-routing concepts and SHALL not mirror arbitrary interactive input unless the user enters a separately identified mirroring mode.

#### Scenario: User taps the Main Device
- **WHEN** no live mirroring mode is active
- **THEN** the tap SHALL affect only the Main Device
