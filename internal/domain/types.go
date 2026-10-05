package domain

import "time"

type DeviceID string
type CommandID string
type SnapshotID string

type DeviceState string

const (
	DeviceStopped      DeviceState = "stopped"
	DeviceProvisioning DeviceState = "provisioning"
	DeviceStarting     DeviceState = "starting"
	DeviceBooting      DeviceState = "booting"
	DeviceReady        DeviceState = "ready"
	DeviceStopping     DeviceState = "stopping"
	DeviceSnapshotting DeviceState = "snapshotting"
	DeviceRestoring    DeviceState = "restoring"
	DeviceFailed       DeviceState = "failed"
	DeviceDeleting     DeviceState = "deleting"
)

type Device struct {
	ID              DeviceID
	Name            string
	State           DeviceState
	ImagePackage    string
	ImageRevision   string
	RuntimeRevision string
	HardwareProfile string
	AVDPath         string
	ADBSerial       string
	ProcessID       int
	IsMain          bool
	Selected        bool
	Version         uint64
	LastErrorCode   string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type CommandState string

const (
	CommandQueued    CommandState = "queued"
	CommandRunning   CommandState = "running"
	CommandSucceeded CommandState = "succeeded"
	CommandFailed    CommandState = "failed"
	CommandCancelled CommandState = "cancelled"
	CommandTimedOut  CommandState = "timed-out"
)

type Command struct {
	ID             CommandID
	DeviceID       DeviceID
	Type           string
	State          CommandState
	ParametersJSON []byte
	IdempotencyKey string
	CorrelationID  string
	ErrorCode      string
	CreatedAt      time.Time
	StartedAt      *time.Time
	FinishedAt     *time.Time
}

type CheckStatus string

const (
	CheckPass CheckStatus = "pass"
	CheckWarn CheckStatus = "warn"
	CheckFail CheckStatus = "fail"
)

type PreflightCheck struct {
	ID          string
	Status      CheckStatus
	Message     string
	Remediation string
}

type SnapshotManifest struct {
	ID              SnapshotID        `json:"id"`
	SourceDeviceID  DeviceID          `json:"sourceDeviceId"`
	CreatedAt       time.Time         `json:"createdAt"`
	RuntimeRevision string            `json:"runtimeRevision"`
	ImagePackage    string            `json:"imagePackage"`
	ImageRevision   string            `json:"imageRevision"`
	GraphicsMode    string            `json:"graphicsMode"`
	ProfileSchema   string            `json:"profileSchema"`
	IncludedState   []string          `json:"includedState"`
	Files           map[string]string `json:"files"`
	TotalBytes      int64             `json:"totalBytes"`
}

type Entitlement struct {
	AccountID         string
	InstallationID    string
	Plan              string
	ConcurrentDevices int
	Features          map[string]bool
	IssuedAt          time.Time
	ExpiresAt         time.Time
	OfflineDeadline   time.Time
}
