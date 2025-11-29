<div align="center">

# zGate Platform
### Zero Trust Ephemeral Database Access

Minimal gateway providing short‑lived database credentials, policy‑based authorization, and audited proxy sessions.

</div>

---

## Overview
The current repository hosts the core server only (no CLI helper, no YAML seed loader). All metadata (users, roles, databases, permissions, refresh tokens) is stored in a SQLite database initialized automatically on first run. Users, roles, and databases must be inserted programmatically (future admin endpoints or migrations can be added). Each CONNECT request creates an ephemeral DB principal via the vendor protocol layer and starts a dynamic proxy port for the session.

## Features
| Feature | Status | Notes |
|---------|--------|-------|
| Ephemeral DB Users | ✔ | Created at `POST /api/connect`, destroyed at disconnect.
| JWT Auth + Refresh | ✔ | Access 15m, refresh 7d with rotation.
| Role + Custom Perms | ✔ | Stored in SQLite tables (`roles`, `role_permissions`, `user_roles`).
| Multi-DB (MSSQL/MySQL) | ✔ | Vendor handlers under `internal/protocol/`.
| Dynamic Proxy Ports | ✔ | Allocated per session by `proxy.Manager`.
| Session Revocation | ✔ | `DELETE /api/sessions/{id}` and `POST /api/logout`.
| Audited Events | ✔ | Structured logs via `utils.Logger`.
| SQLite Metadata Store | ✔ | Auto schema creation; encrypted sensitive fields via 32‑byte key.

## Directory Layout (Current)
```
.
├── data/                # Default location for SQLite file
├── internal/
│   ├── api/             # HTTP handlers (login, refresh, databases, connect, sessions)
│   ├── auth/            # JWT creation/validation & refresh token utils
│   ├── conn/            # Connection abstractions
│   ├── gateway/         # Session orchestration helpers
│   ├── policy/          # Policy engine (permission resolution)
│   ├── protocol/        # DB protocol managers & temp principal logic
│   │   ├── mssql/
│   │   └── mysql/
│   ├── proxy/           # Session proxy lifecycle & credential generation
│   ├── store/           # SQLite store, schema, CRUD for users/roles/tokens
│   └── utils/           # Logger initialization
├── main.go              # Server entrypoint
└── README.md
```

## Data Model (SQLite)
Tables: `databases`, `roles`, `role_permissions`, `users`, `user_roles`, `user_custom_permissions`, `refresh_tokens`.
Important constraints:
- Role + database permission uniqueness enforced (`UNIQUE(role_name, database_name)`).
- Refresh tokens tracked with revocation + rotation metadata.
- User custom permissions supplement role permissions.

## Security Model
| Aspect | Detail |
|--------|-------|
| Zero Credentials at Rest | Temp DB users created only inside session lifecycle.
| Ephemeral Principals | Naming pattern started in `protocol/manager.go` (`zgate_<base>_<suffix>`).
| Token Strategy | Access: 15m; Refresh: 7d; rotated on refresh, old revoked.
| Store Encryption | Sensitive admin passwords stored encrypted using 32‑byte key.
| Audit Logging | Structured logs for auth, session, proxy lifecycle, token actions.

## Configuration
Environment variables (server will exit if required ones missing):
| Variable | Required | Description |
|----------|----------|-------------|
| `ZGATE_PORT` | Yes | Port number (e.g. `8080`). Colon auto-added if absent.
| `ZGATE_STORE_KEY` | Yes | 64 hex chars (32 bytes) key; AES-256 for sensitive fields.
| `ZGATE_JWT_SECRET` | Yes | HMAC secret for signing access tokens.
| `ZGATE_STORE_PATH` | No | Path to SQLite file; defaults to `data/zgate.db`.

`.env` is loaded automatically (via `godotenv`).

## Running the Server
```bash
export ZGATE_PORT=8080
export ZGATE_STORE_KEY="$(openssl rand -hex 32)"  # must decode to 32 bytes
export ZGATE_JWT_SECRET="$(openssl rand -hex 32)"
go run main.go --api-addr :8080
```
On startup the schema is ensured; if the database file is absent it is created.

## API Endpoints
Public:
- `POST /api/login` {username,password} → access + refresh tokens
- `POST /api/refresh` {refresh_token} → rotated tokens
- `POST /api/logout` {refresh_token} → revokes token

Authenticated (Bearer access token):
- `GET /api/databases` → list databases user can access (via policy engine)
- `POST /api/connect` {database_name} → starts proxy, returns port + temp creds
- `POST /api/disconnect` {database_name} → stops session, drops temp user
- `GET /api/sessions` → enumerate active refresh token sessions
- `DELETE /api/sessions/{id}` → revoke specific session

## Example Flow (cURL)
```bash
# Login
curl -s -X POST http://localhost:8080/api/login \
  -d '{"username":"alice","password":"secret"}' \
  -H 'Content-Type: application/json' | jq .

# Use access token to list databases
ACCESS=... # fill from login response
curl -s -H "Authorization: Bearer $ACCESS" http://localhost:8080/api/databases | jq .

# Connect (returns temp DB credentials + proxy port)
curl -s -X POST http://localhost:8080/api/connect \
  -H "Authorization: Bearer $ACCESS" \
  -d '{"database_name":"my_mssql"}' | jq .
```

## Development
```bash
go run main.go --api-addr :8080
Logging uses `utils.InitLogger()`; adjust implementation for levels/format as needed.


