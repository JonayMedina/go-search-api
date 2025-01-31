# [![Tribal Worldwide](../assets/header.png?style=centerme)](https://tribalworldwide.gt/)

## Register

Use this service to register a new user

> This service requires encription of the request and response body.  For futher information see the [API Encription Guide](../../encription-guide.md)

### **Endpoint**

```curl
POST 'http://localhost:8080/api/auth/register'
```

---

### **Request**

#### Headers

| HEADER          | TYPE     | RULE           | VALUE               |
|-----------------|----------|----------------|---------------------|
| `Content-Type`  | `String` | **`Required`** | `application/json`  |

#### JSON Body

`Register Request`

##### Register request object definition

| PARAMETER     | TYPE     | RULE           | VALUE     |
|---------------|----------|----------------|-----------|
| `username`    | `String` | **`Required`** | username  |
| `email`       | `String` | **`Required`** | email     |
| `password`    | `String` | **`Required`** | password  |

### **Response**

- On success, the HTTP status code in the response header is `200 OK` and the response body contains a simple success response.

On error, the header status code is an [`error code`](../README.md).

- If an unexpected error occurs, the request will return `500 INTERNAL SERVER ERROR` HTTP response code with `internal_server_error` API response.

#### Register response object definition

| PARAMETER           | TYPE      | RULE           | VALUE               |
|---------------------|-----------|----------------|---------------------|
| `message`           | `String`  | **`Required`** | message             |
| `user`              | `Json`    | **`Required`** | user                |
| `user.id`           | `String`  | **`Required`** | user id             |
| `user.username`     | `String`  | **`Required`** | user username       |
| `user.email`        | `String`  | **`Required`** | user email          |
| `user.created_at`   | `String`  | **`Required`** | user created_at     |
| `user.updated_at`   | `String`  | **`Required`** | user updated_at     |

---

### **Samples**

#### Success

```json
POST 'http://localhost:8080/api/auth/register'
--header 'Content-Type: application/json'
--data '{
    "username": "testuser2",
    "email": "test+2@example.com",
    "password": "password"
}'

HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 31 Jan 2025 08:33:00 GMT

{
    "message": "Usuario registrado exitosamente",
    "user": {
        "id": "678eaff323fb2b1c153c1bf6",
        "username": "testuser2",
        "email": "test+2@example.com"
    }
}
```

---

[API Overview](../README.md)

[![Tribal Worldwide](../assets/footer.png?style=centerme)](https://tribalworldwide.gt/)
