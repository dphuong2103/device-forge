# Contract v1

- `openapi.yaml` defines loopback HTTP resources and the authenticated event upgrade.
- `events.schema.json` defines event envelopes. Consumers ignore unknown event types only after negotiating a compatible minor version.
- `signaling.schema.json` defines bounded WebRTC offer, answer, and ICE messages.

Mutation requests use optimistic device versions and optional idempotency keys. Long-running operations return `202` and a command resource. Errors use stable machine codes plus a correlation ID; secret or environment-specific values never appear in `detail`.
