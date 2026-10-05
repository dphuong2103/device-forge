package domain

import (
	"context"
	"io"
	"time"
)

type Runtime interface {
	Preflight(context.Context) ([]PreflightCheck, error)
	Provision(context.Context, Device) error
	Start(context.Context, Device) (Process, error)
	Stop(context.Context, Device, time.Duration) error
	ForceStop(context.Context, Device) error
	Delete(context.Context, Device) error
	Observe(context.Context, Device) (RuntimeObservation, error)
}

type Process struct {
	PID         int
	ADBSerial   string
	ConsolePort int
}

type RuntimeObservation struct {
	ProcessRunning bool
	BootCompleted  bool
	ADBHealthy     bool
	PID            int
	ADBSerial      string
}

type ADBController interface {
	InstallAPK(context.Context, DeviceID, io.Reader) error
	InjectInput(context.Context, DeviceID, InputEvent) error
	ApplyProperty(context.Context, DeviceID, string, string) error
}

type InputEvent struct {
	Type     string
	X        float64
	Y        float64
	Key      string
	Text     string
	Rotation int
}

type Media interface {
	Start(context.Context, DeviceID, MediaQuality) (MediaSession, error)
	Stop(context.Context, string) error
	SetQuality(context.Context, string, MediaQuality) error
}

type MediaQuality struct {
	MaxWidth int
	MaxFPS   int
	Bitrate  int
}

type MediaSession struct {
	ID       string
	DeviceID DeviceID
	OfferSDP string
}

type ProfileManager interface {
	Generate(context.Context, string, string, string) (DeviceProfile, error)
	Validate(context.Context, Device, DeviceProfile) ([]FieldResult, error)
	Apply(context.Context, Device, DeviceProfile) ([]FieldResult, error)
}

type DeviceProfile struct {
	SchemaVersion string
	Seed          string
	Fields        map[string]string
}

type FieldResult struct {
	Field   string
	Status  string
	Message string
}

type ProxyManager interface {
	Apply(context.Context, DeviceID, ProxyConfiguration) error
	Verify(context.Context, DeviceID) error
	Remove(context.Context, DeviceID) error
}

type ProxyConfiguration struct {
	Protocol  string
	Host      string
	Port      int
	Username  string
	SecretRef string
	RemoteDNS bool
}

type LocationManager interface {
	Set(context.Context, DeviceID, GeoPoint) error
	PlayRoute(context.Context, DeviceID, Route) error
	StopRoute(context.Context, DeviceID) (GeoPoint, error)
}

type GeoPoint struct {
	Latitude  float64
	Longitude float64
	Altitude  float64
}

type Route struct {
	Points          []GeoPoint
	MetersPerSecond float64
	Interval        time.Duration
}

type Snapshots interface {
	CreateFast(context.Context, Device, string) (SnapshotManifest, error)
	CreateManaged(context.Context, Device) (SnapshotManifest, error)
	Validate(context.Context, SnapshotManifest) error
	Restore(context.Context, Device, SnapshotManifest) error
	Delete(context.Context, SnapshotID) error
}

type SecretStore interface {
	Put(context.Context, string, []byte) error
	Get(context.Context, string) ([]byte, error)
	Delete(context.Context, string) error
}

type Entitlements interface {
	Refresh(context.Context) (Entitlement, error)
	Current(context.Context) (Entitlement, error)
	Allows(Entitlement, string, time.Time) error
}

type DeviceRepository interface {
	List(context.Context) ([]Device, error)
	Get(context.Context, DeviceID) (Device, error)
	Create(context.Context, Device) error
	Update(context.Context, Device, uint64) error
	Delete(context.Context, DeviceID) error
	SetMain(context.Context, DeviceID) error
}

type CommandRepository interface {
	Create(context.Context, Command) error
	Update(context.Context, Command) error
	FindByIdempotencyKey(context.Context, string) (Command, error)
}
