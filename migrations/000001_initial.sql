PRAGMA foreign_keys = ON;

CREATE TABLE schema_migrations (
    version INTEGER PRIMARY KEY,
    applied_at TEXT NOT NULL
) STRICT;

CREATE TABLE devices (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL CHECK(length(name) BETWEEN 1 AND 80),
    state TEXT NOT NULL CHECK(state IN ('stopped','provisioning','starting','booting','ready','stopping','snapshotting','restoring','failed','deleting')),
    image_package TEXT NOT NULL,
    image_revision TEXT NOT NULL,
    runtime_revision TEXT NOT NULL,
    hardware_profile TEXT NOT NULL,
    avd_path TEXT NOT NULL UNIQUE,
    adb_serial TEXT NOT NULL DEFAULT '',
    process_id INTEGER NOT NULL DEFAULT 0 CHECK(process_id >= 0),
    is_main INTEGER NOT NULL DEFAULT 0 CHECK(is_main IN (0,1)),
    selected INTEGER NOT NULL DEFAULT 0 CHECK(selected IN (0,1)),
    version INTEGER NOT NULL DEFAULT 1 CHECK(version > 0),
    last_error_code TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
) STRICT;

CREATE UNIQUE INDEX one_main_device ON devices(is_main) WHERE is_main = 1;
CREATE INDEX devices_state ON devices(state);

CREATE TABLE commands (
    id TEXT PRIMARY KEY,
    device_id TEXT REFERENCES devices(id) ON DELETE SET NULL,
    parent_command_id TEXT REFERENCES commands(id) ON DELETE SET NULL,
    type TEXT NOT NULL,
    state TEXT NOT NULL CHECK(state IN ('queued','running','succeeded','failed','cancelled','timed-out')),
    parameters_json TEXT NOT NULL DEFAULT '{}' CHECK(json_valid(parameters_json)),
    idempotency_key TEXT UNIQUE,
    correlation_id TEXT NOT NULL,
    error_code TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL,
    started_at TEXT,
    finished_at TEXT
) STRICT;

CREATE INDEX commands_device_created ON commands(device_id, created_at DESC);
CREATE INDEX commands_state ON commands(state);

CREATE TABLE profiles (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    schema_version TEXT NOT NULL,
    seed TEXT NOT NULL,
    fields_json TEXT NOT NULL CHECK(json_valid(fields_json)),
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
) STRICT;

CREATE TABLE snapshots (
    id TEXT PRIMARY KEY,
    source_device_id TEXT REFERENCES devices(id) ON DELETE SET NULL,
    name TEXT NOT NULL,
    type TEXT NOT NULL CHECK(type IN ('fast','managed')),
    state TEXT NOT NULL CHECK(state IN ('creating','complete','incomplete','deleting')),
    path TEXT NOT NULL UNIQUE,
    manifest_json TEXT NOT NULL CHECK(json_valid(manifest_json)),
    total_bytes INTEGER NOT NULL DEFAULT 0 CHECK(total_bytes >= 0),
    created_at TEXT NOT NULL,
    completed_at TEXT
) STRICT;

CREATE INDEX snapshots_source_created ON snapshots(source_device_id, created_at DESC);

CREATE TABLE settings (
    key TEXT PRIMARY KEY,
    value_json TEXT NOT NULL CHECK(json_valid(value_json)),
    updated_at TEXT NOT NULL
) STRICT;

CREATE TABLE browser_sessions (
    id_hash BLOB PRIMARY KEY,
    csrf_hash BLOB NOT NULL,
    created_at TEXT NOT NULL,
    expires_at TEXT NOT NULL,
    last_seen_at TEXT NOT NULL
) STRICT;

CREATE INDEX browser_sessions_expiry ON browser_sessions(expires_at);

CREATE TABLE bootstrap_tokens (
    token_hash BLOB PRIMARY KEY,
    created_at TEXT NOT NULL,
    expires_at TEXT NOT NULL,
    consumed_at TEXT
) STRICT;

CREATE TABLE cached_entitlements (
    installation_id TEXT PRIMARY KEY,
    compact_token TEXT NOT NULL,
    plan TEXT NOT NULL,
    claims_json TEXT NOT NULL CHECK(json_valid(claims_json)),
    issued_at TEXT NOT NULL,
    expires_at TEXT NOT NULL,
    offline_deadline TEXT NOT NULL,
    last_validated_at TEXT NOT NULL,
    highest_observed_time TEXT NOT NULL
) STRICT;

INSERT INTO schema_migrations(version, applied_at)
VALUES (1, strftime('%Y-%m-%dT%H:%M:%fZ', 'now'));
