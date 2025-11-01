# Chirpy API Reference

Base URL (development): http://localhost:8080

Authentication
- Most endpoints require a Bearer JWT in the `Authorization` header: `Authorization: Bearer <token>`.
- Some webhook/admin endpoints may use an API key (`Authorization: ApiKey <key>`). See handler implementations for details.

Common response shapes
- ChirpResponse

```json
{
  "id": "<uuid>",
  "created_at": "RFC3339 timestamp",
  "updated_at": "RFC3339 timestamp",
  "body": "text up to 140 chars",
  "user_id": "<uuid>"
}
```

- CreateUser / Login responses include `token` and `refresh_token` where applicable.

HTTP endpoints

1) Health
- Method: GET
- Path: /api/healthz
- Auth: none
- Response: 200 OK (plain text "OK")

2) Create user
- Method: POST
- Path: /api/users
- Auth: none
- Request JSON:

```json
{
  "email": "user@example.com",
  "password": "secret"
}
```

- Success: 201 Created with created user object (includes id, timestamps). Handler: `internal/handlers/user.go`

3) Login
- Method: POST
- Path: /api/login
- Auth: none
- Request JSON (email + password)
- Success: 200 OK with a JSON payload containing `token` (access JWT) and `refresh_token`.

4) Refresh access token
- Method: POST
- Path: /api/refresh
- Auth: Bearer refresh token
- Success: 200 OK with new access token JSON.

5) Revoke refresh token
- Method: POST
- Path: /api/revoke
- Auth: Bearer refresh token
- Success: 204 No Content

6) Create chirp
- Method: POST
- Path: /api/chirps
- Auth: Bearer access token
- Request JSON (models.ChirpRequest):

```json
{
  "body": "Hello world",
  "user_id": "<uuid>" // optional; server resolves user from JWT
}
```

- Rules: body max 140 characters; profanity filtered automatically using `internal/utils/profanity.go`. If profanity is removed, server will still accept chirp but log filtering.
- Success: 201 Created with `ChirpResponse`.

7) List chirps
- Method: GET
- Path: /api/chirps
- Query params:
  - `author_id` (optional UUID) — when provided, filters to chirps by that author
  - `sort` (optional) — `asc` (default) or `desc` to control order by creation time
- Success: 200 OK with JSON array of `ChirpResponse` objects.

8) Get chirp by ID
- Method: GET
- Path: /api/chirps/{chirpID}
- Success: 200 OK with `ChirpResponse`
- Error: 404 Not Found if ID does not exist.

9) Delete chirp
- Method: DELETE
- Path: /api/chirps/{chirpID}
- Auth: Bearer access token
- Authorization: only the chirp author may delete their chirps.
- Success: 204 No Content
- Errors: 403 Forbidden when authenticated user is not the author.

Admin & Webhooks

- GET /admin/metrics — returns simple metrics; see `internal/handlers/admin.go`.
- POST /admin/reset — development-only reset that wipes test data (dangerous!).
- POST /api/polka/webhooks — expects `Authorization: ApiKey <key>`; used for Polka webhook handling.

Error handling summary
- 400 Bad Request — invalid input, invalid UUID, too long chirp body
- 401 Unauthorized — missing/invalid/expired token
- 403 Forbidden — insufficient permissions (e.g., deleting another's chirp)
- 404 Not Found — resource not found
- 500 Internal Server Error — unexpected server/db error

Examples

- Create chirp (curl):

```sh
curl -X POST http://localhost:8080/api/chirps \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{"body":"Hello from curl"}'
```

- List chirps by author sorted desc:

```sh
curl 'http://localhost:8080/api/chirps?author_id=<uuid>&sort=desc'
```

DB & generated code

- SQL migrations: `database/schema`
- Query files: `database/queries` (e.g., `chirps.sql`, `user.sql`, `refresh_token.sql`)
- sqlc generates typed DB wrappers into `internal/database`

If you'd like an OpenAPI spec (swagger) or a Postman collection, I can add one next.
