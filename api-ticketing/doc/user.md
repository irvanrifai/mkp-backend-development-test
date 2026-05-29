# User API Documentation

## Register
- Method: `POST`
- Path: `/api/auth/register`

### Request Body
```json
{
  "name": "string",
  "username": "string",
  "email": "string",
  "password": "string"
}
```

### Success Response
- Status: `201 Created`
- Body:
```json
{
  "id": 1,
  "name": "John Doe",
  "username": "johndoe",
  "email": "john@example.com",
  "phone": null,
  "created_at": "2026-05-29T00:00:00Z",
  "updated_at": "2026-05-29T00:00:00Z"
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
  {"errors":{"email":"Field ini wajib diisi"}}
  ```
- `500 Internal Server Error`
  - If usecase repository fails
  - Example:
  ```json
  {"error":"database error"}
  ```

## Login
- Method: `POST`
- Path: `/api/auth/login`

### Request Body
```json
{
  "email": "string",
  "password": "string"
}
```

### Success Response
- Status: `200 OK`
- Body:
```json
{
  "token": "string"
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
  {"errors":{"password":"Field ini wajib diisi"}}
  ```
- `401 Unauthorized`
  - Wrong credentials
  - Example:
  ```json
  {"error":"Wrong credentials"}
  ```
