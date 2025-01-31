# [![Tribal Worldwide](../assets/header.png?style=centerme)](https://tribalworldwide.gt/)

## Login Service

Use this service to login to the system

### **Endpoint**

```curl
POST 'http://localhost:8080/api/auth/login'
```

---

### **Request**

#### Headers

| HEADER          | TYPE     | RULE           | VALUE            |
|-----------------|----------|----------------|------------------|
| `Content-Type`  | `String` | **`Required`** | application/json |

#### JSON Body

`Login Request` object wrapped in a [Request Wrapper object](../../general-object-models.md)

##### Login request object definition

| PARAMETER     | TYPE      | RULE           | VALUE         |
|---------------|-----------|----------------|---------------|
| `username`    | `String`  | **`Required`** | username      |
| `password`    | `String`  | **`Required`** | password      |

### **Response**

- On success, the HTTP status code in the response header is `200 OK` and the response body contains a simple success response.

On error, the header status code is an [`error code`](../README.md).

- If an unexpected error occurs, the request will return `500 INTERNAL SERVER ERROR` HTTP response code with `internal_server_error` API response.

#### Login response object definition

| PARAMETER            | TYPE     | RULE           | VALUE              |
|----------------------|----------|----------------|--------------------|
| `token`              | `String` | **`Required`** | token              |
| `user`               | `Json`   | **`Required`** | user               |
| `user.id`            | `String` | **`Required`** | user id            |
| `user.username`      | `String` | **`Required`** | user username      |
| `user.email`         | `String` | **`Required`** | user email         |
| `user.created_at`    | `String` | **`Required`** | user created_at    |
| `user.updated_at`    | `String` | **`Required`** | user updated_at    |

---

### **Samples**

#### Success

```json
POST 'http://localhost:8080/api/auth/login'
--header 'Content-Type: application/json'
--data '{
    "username": "testuser2",
    "password": "password1234"
}'


HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 31 Jan 2025 08:33:00 GMT

{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzg0MjMyMDMsInVzZXJuYW1lIjoidGVzdHVzZXIyIn0.khb66Yg6dLNB2mvba_rdDnL7SkujiLuroFBs8dqqBW0",
    "user": {
        "id": "678eaff323fb2b1c153c1bf6",
        "username": "testuser2",
        "email": "test+2@example.com",
        "created_at": "2025-01-20T20:20:03.73Z",
        "updated_at": "2025-01-20T20:20:03.73Z"
    }
}

```

---

[API Overview](../README.md)

[![Tribal Worldwide](../assets/footer.png?style=centerme)](https://tribalworldwide.gt/)
