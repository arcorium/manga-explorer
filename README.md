# MANGA EXPLORER

Simple manga reader backend app

## Features

### User

- Role Based Access Control (RBAC) `LIMITED`
- Auth by Access Token
- Forget Password
- Email Verification
- Multi-Login Support
- CRUD User
- Login using external services (OAuth) **(TODO)**

### Manga

- Advance Search
- Get Random Manga
- Hierarchical Comments (Deep Nesting Reply Support)
- Bookmark
- History
- Rating
- CRUD Manga (Genre, Cover, Volume, Translation, Chapter, Page)
- Recommendation By History **(TODO)**
- Popular Manga **(TODO)**

## Quick Start

Prepare your `.env` file and place it in root project. the example is provided on `.env.example`.
Then run this :

```cmd
> docker-compose -f .\compose.yml -p manga-explorer up -d
```

Or for newer docker compose you can do this instead:

```cmd
> docker compose -f .\compose.yml -p manga-explorer up -d
```

You need to have docker-compose installed on your system

**NOTE:** If the backend failed to run, you can rerun the command, this is due to the backend
running first when the
database service is not ready yet.

Default admin details is:

```text
username: admin
email: admin@manga-explorer.com
password: adminadmin
```

You can change that on `docker-compose.yml` or set environment variable
for `ME_ADMIN_USERNAME`, `ME_ADMIN_PASSWORD`, `ME_ADMIN_EMAIL`.

Default endpoint is:

```text
Ip: ~
Port: 9999 (Both backend and docker)
Prefix: /api
```

You can change the exposed ip on `docker-compose.yml` to `80:9999` or others.

## Building

### Prerequisite

- [Golang v1.21](https://go.dev/dl/) `Optional`
- Postgresql `Optional`
- Docker`Optional`

### Docker

Build the image by:

```cmd
> docker build -t arc/manga-explorer .
```

Run the image by:

```cmd
> docker run -p 9999:9999 -d arc/manga-explorer
```

Environment Variables:

```dotenv
ME_ADMIN_USERNAME
ME_ADMIN_PASSWORD
ME_ADMIN_EMAIL
```

**NOTE:** Make sure you have Postgres database running and the details is set on `.env`

### Local

```cmd
> make build
> make migrate
> make run
```

## Response

Response is always start with the response message as key, that is `success` or `error` and followed
by the data object.

### Error Response

| Key     | Type                |
|---------|---------------------|
| code    | `number`            |
| message | `string`            |
| details | `object` `OPTIONAL` |

For Example:

```json
{
  "error": {
    "code": 28,
    "message": "Token you provide is expired"
  }
}
```

### Success Response

| Key  | Type                |
|------|---------------------|
| code | `number`            |
| data | `object`            |
| page | `object` `OPTIONAL` |

For Example:

```json
{
  "success": {
    "code": 0,
    "data": {
      "token_type": "Bearer",
      "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTUxNDkwNTUsImlhdCI6MTcxNTE0NTQ1NSwiaWQiOiJjM2Y5ZjdkMS0xNWVkLTQwMDktOTE3Yi0wMmY3ZWNmYTE3M2YiLCJpc3MiOiJtYW5nYS1leHBsb3JlciIsIm5hbWUiOiJhZG1pbiIsIm5iZiI6MTcxNTE0NTQ1NSwicm9sZSI6ImFkbWluIiwidWlkIjoiYTg2YWNhZmItYTNjYi00MjhhLWJkNDgtMjMzZDE5Yjg2ZWU3In0.TAA0s1jNCi9NtmLWecKHvGLPYi8htXjK8fKEoyLH34I"
    }
  }
}
```

## Usage

All endpoint documentations are defined on `docs/docs.go` and can be viewed
at `localhost:9999/api/docs/` (Change the ip and port)

### User

Prefix: `/v{VERSION}/users/`

Example: `localhost:9999/api/v1/users/`

**(TODO)**

### Auth

Prefix: `/v{VERSION}/auth/`

Example: `localhost:9999/api/v1/auth/`

**(TODO)**

### Manga

Prefix: `/v{VERSION}/mangas/`

Example: `localhost:9999/api/v1/mangas/`

**(TODO)**
