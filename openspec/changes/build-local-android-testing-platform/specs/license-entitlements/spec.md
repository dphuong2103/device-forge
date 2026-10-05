## ADDED Requirements

### Requirement: Account and subscription validation
The cloud service SHALL authenticate users, derive subscription state from the authoritative billing record, and issue signed entitlements containing plan features and limits.

#### Scenario: Active subscriber signs in
- **WHEN** an authenticated account has an active subscription and an activation is permitted
- **THEN** the cloud SHALL issue a signed entitlement for that account and installation with its effective limits and expiry policy

### Requirement: Local entitlement verification
The Agent SHALL verify entitlement signature, audience, installation binding, issue time, expiry, and feature limits before enabling entitled operations.

#### Scenario: Entitlement signature is invalid
- **WHEN** the cached or received entitlement fails signature verification
- **THEN** the Agent SHALL reject it and SHALL not enable paid operations from its claims

### Requirement: Bounded offline operation
The Agent SHALL continue entitled local operation during cloud unavailability only while the last valid entitlement remains within its signed offline-grace policy.

#### Scenario: Cloud is unavailable within grace
- **WHEN** refresh fails for a network reason and the cached entitlement remains within its offline window
- **THEN** the Agent SHALL continue allowed operations and SHALL show the next required validation time

#### Scenario: Offline grace expires
- **WHEN** cloud validation cannot be completed before the signed offline window ends
- **THEN** the Agent SHALL block newly initiated entitled operations while preserving access to diagnostics, export, stop, and deletion actions

### Requirement: Installation activation limits
The cloud SHALL enforce subscription-specific installation activation limits and SHALL support revoking an activation.

#### Scenario: Activation limit is reached
- **WHEN** another installation requests activation beyond the plan limit
- **THEN** the cloud SHALL deny activation and SHALL return a non-secret explanation suitable for remediation

### Requirement: Minimal cloud data
The cloud SHALL not require or accept APKs, virtual-device contents, snapshots, proxy credentials, location history, screen content, or interaction logs for licensing.

#### Scenario: Entitlement telemetry is recorded
- **WHEN** the cloud records an entitlement refresh
- **THEN** it SHALL retain only account, installation, plan, version, timing, and necessary security audit metadata under the applicable retention policy
