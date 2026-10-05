## ADDED Requirements

### Requirement: Per-user Agent availability
The system SHALL run the primary Agent in the signed-in Windows user's session and SHALL provide a tray entry point that starts or discovers the Agent and opens the local user interface.

#### Scenario: User opens the product
- **WHEN** the user selects the product from the tray application
- **THEN** the tray application SHALL ensure the Agent is available and SHALL open its authenticated loopback user interface

### Requirement: Secure loopback interface
The Agent SHALL bind its UI and APIs only to loopback interfaces and SHALL require an authenticated local session for protected HTTP and WebSocket operations.

#### Scenario: Untrusted origin calls the Agent
- **WHEN** a request or WebSocket upgrade has an unapproved origin, host, or session
- **THEN** the Agent SHALL reject it without executing the requested operation

### Requirement: One-time browser bootstrap
The Agent SHALL use a single-use, short-lived bootstrap credential to establish a browser session and SHALL not expose durable credentials in the opened URL.

#### Scenario: Bootstrap token is replayed
- **WHEN** a previously exchanged bootstrap credential is presented again
- **THEN** the Agent SHALL reject the credential and SHALL not create another session

### Requirement: Local persistence and reconciliation
The Agent SHALL persist configuration and operation state locally and SHALL reconcile that state with emulator processes and connected devices after restart.

#### Scenario: Agent restarts while a device is running
- **WHEN** the Agent starts and finds a managed emulator process that is recorded locally
- **THEN** it SHALL reconnect to the device, update its observed lifecycle state, and expose the reconciled state before accepting mutations

### Requirement: Local test-data boundary
The system SHALL keep APKs, device contents, snapshots, proxy credentials, location history, screen content, and interaction logs off the product cloud.

#### Scenario: Agent refreshes an entitlement
- **WHEN** the Agent contacts the cloud for subscription validation
- **THEN** the request SHALL contain only account, installation, entitlement, version, and necessary security metadata
