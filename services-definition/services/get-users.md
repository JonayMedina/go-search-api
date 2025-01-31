# [![Tribal Worldwide](../assets/header.png?style=centerme)](https://tribalworldwide.gt/)

## Get Users

Use this service to get the users of the system.

> This service requires Authentication of the request.

### **Endpoint**

```curl
GET 'http://localhost:8080/api/get-users'
```

---

### **Request**

#### Headers

| HEADER          | TYPE     | RULE           | VALUE               |
|-----------------|----------|----------------|---------------------|
| `Content-Type`  | `String` | **`Required`** | `application/json`  |
| `Authorization` | `String` | **`Required`** | `Bearer Token`      |

#### JSON

`Get Users Request`

##### Get Users request object definition

| PARAMETER | TYPE  | RULE | VALUE |
|-----------|-------|------|-------|

### **Response**

- On success, the HTTP status code in the response header is `200 OK` and the response body contains a simple success response.

On error, the header status code is an [`error code`](../README.md).

- If an unexpected error occurs, the request will return `500 INTERNAL SERVER ERROR` HTTP response code with `internal_server_error` API response.

#### Get App Area Sections response object definition

| PARAMETER             | TYPE      | RULE           | VALUE                                            |
|-----------------------|-----------|----------------|--------------------------------------------------|
| `users[]`             | `Array`   | **`Required`** | array cotaining the users of the system          |
| `users[].id`          | `String`  | **`Required`** | user id                                          |
| `users[].username`    | `String`  | **`Required`** | user username                                    |
| `users[].email`       | `String`  | **`Required`** | user email                                       |
| `users[].created_at`  | `String`  | **`Required`** | user created_at                                  |
| `users[].updated_at`  | `String`  | **`Required`** | user updated_at                                  |

---

### **Samples**

#### Success

```json
GET 'http://localhost:8080/api/get-users'
--header 'Content-Type: application/json'
--header 'Authorization: Bearer Token'

HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 31 Jan 2025 08:33:00 GMT

{
    "users": [
        {
            "id": "678eaff323fb2b1c153c1bf6",
            "username": "testuser2",
            "email": "test+2@example.com",
            "created_at": "2025-01-20T20:20:03.73Z",
            "updated_at": "2025-01-20T20:20:03.73Z"
        }
    ]
}

```

---

[API Overview](../README.md)

[![Tribal Worldwide](../assets/footer.png?style=centerme)](https://tribalworldwide.gt/)
