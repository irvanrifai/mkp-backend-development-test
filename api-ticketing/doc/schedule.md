# Schedule API Documentation

## List Schedules
- Method: `GET`
- Path: `/api/schedules`
- Auth: Protected (JWT)

### Success Response
- Status: `200 OK`
- Body:
```json
[
  {
    "id": 1,
    "movie_id": 1,
    "studio_id": 1,
    "show_time": "2026-05-29T15:00:00Z",
    "price_per_ticket": 125000,
    "status": "ACTIVE",
    "created_at": "2026-05-29T00:00:00Z",
    "updated_at": "2026-05-29T00:00:00Z",
    "deleted_at": null
  }
]
```

### Error Responses
- `500 Internal Server Error`
  - Example:
  ```json
  {"error":"database error"}
  ```

## Create Schedule
- Method: `POST`
- Path: `/api/schedules`
- Auth: Protected (JWT)

### Request Body
```json
{
  "movie_id": 1,
  "studio_id": 1,
  "show_time": "2026-05-29T15:00:00Z",
  "price_per_ticket": 125000,
  "status": "ACTIVE"
}
```

### Success Response
- Status: `201 Created`
- Body:
```json
{
  "id": 1,
  "movie_id": 1,
  "studio_id": 1,
  "show_time": "2026-05-29T15:00:00Z",
  "price_per_ticket": 125000,
  "status": "ACTIVE",
  "created_at": "2026-05-29T00:00:00Z",
  "updated_at": "2026-05-29T00:00:00Z",
  "deleted_at": null
}
```

### Error Responses
- `400 Bad Request`
  - Invalid JSON or missing body
  - Example:
  ```json
  {"error":"Bad Request"}
  ```
- `422 Unprocessable Entity`
  - Validation error
  - Example:
  ```json
  {"errors":{"movie_id":"Field ini wajib diisi"}}
  ```
- `500 Internal Server Error`
  - Example:
  ```json
  {"error":"database error"}
  ```

## Get Schedule Detail
- Method: `GET`
- Path: `/api/schedules/:id`
- Auth: Protected (JWT)

### Success Response
- Status: `200 OK`
- Body:
```json
{
  "id": 1,
  "movie_id": 1,
  "studio_id": 1,
  "show_time": "2026-05-29T15:00:00Z",
  "price_per_ticket": 125000,
  "status": "ACTIVE",
  "created_at": "2026-05-29T00:00:00Z",
  "updated_at": "2026-05-29T00:00:00Z",
  "deleted_at": null
}
```

### Error Responses
- `400 Bad Request`
  - Invalid schedule ID
  - Example:
  ```json
  {"error":"Invalid schedule ID"}
  ```
- `500 Internal Server Error`
  - Example:
  ```json
  {"error":"database error"}
  ```

## Update Schedule
- Method: `PUT`
- Path: `/api/schedules/:id`
- Auth: Protected (JWT)

### Request Body
```json
{
  "movie_id": 2,
  "studio_id": 3,
  "show_time": "2026-05-29T18:00:00Z",
  "price_per_ticket": 130000,
  "status": "INACTIVE"
}
```

### Success Response
- Status: `200 OK`
- Body:
```json
{
  "id": 1,
  "movie_id": 2,
  "studio_id": 3,
  "show_time": "2026-05-29T18:00:00Z",
  "price_per_ticket": 130000,
  "status": "INACTIVE",
  "created_at": "2026-05-29T00:00:00Z",
  "updated_at": "2026-05-29T00:00:00Z",
  "deleted_at": null
}
```

### Error Responses
- `400 Bad Request`
  - Invalid schedule ID or invalid payload
  - Example:
  ```json
  {"error":"Invalid schedule ID"}
  ```
- `422 Unprocessable Entity`
  - Validation error
  - Example:
  ```json
  {"errors":{"price_per_ticket":"Field ini wajib diisi"}}
  ```
- `500 Internal Server Error`
  - Example:
  ```json
  {"error":"database error"}
  ```

## Delete Schedule
- Method: `DELETE`
- Path: `/api/schedules/:id`
- Auth: Protected (JWT)

### Success Response
- Status: `200 OK`
- Body:
```json
{
  "message": "ok"
}
```

### Error Responses
- `400 Bad Request`
  - Invalid schedule ID
  - Example:
  ```json
  {"error":"Invalid schedule ID"}
  ```
- `500 Internal Server Error`
  - Example:
  ```json
  {"error":"database error"}
  ```
