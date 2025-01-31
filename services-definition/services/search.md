# [![Tribal Worldwide](../assets/header.png?style=centerme)](https://tribalworldwide.gt/)

## Search Music

Use this service to search music by title, artist, album, or genre.

> This service requires HEADER -> `Authorization` with value -> `Bearer Token` of the request and response body.

### **Endpoint**

```curl
POST 'http://localhost:8080/api/search?q=query'
```

---

### **Request**

#### Headers

| HEADER          | TYPE               | RULE           | VALUE          |
|-----------------|--------------------|----------------|----------------|
| `Content-Type`  | `Application/json` | **`Required`** |                |
| `Authorization` | `Bearer Token`     | **`Required`** |                |

#### Url Parameters

| PARAMETER | TYPE      | RULE              | VALUE                 |
|-----------|-----------|-------------------|-----------------------|
| `q`       | `String`  | **`Required`**    | The search query.     |

### **Response**

- On success, the HTTP status code in the response header is `200 OK` and the response body contains a simple success response.

On error, the header status code is an [`error code`](../README.md).

#### Search Music response object definition

| PARAMETER             | TYPE     | RULE           | VALUE                                         |
|-----------------------|----------|----------------|-----------------------------------------------|
| `songs`               | `Array`  | **`Required`** | array containing the list of songs finded     |
| `songs[].name`        | `String` | **`Required`** | song name                                     |
| `songs[].artist`      | `String` | **`Required`** | song artist                                   |
| `songs[].duration`    | `String` | **`Required`** | song duration                                 |
| `songs[].album`       | `String` | **`Required`** | song album                                    |
| `songs[].artwork`     | `String` | **`Required`** | song artwork                                  |
| `songs[].price`       | `String` | **`Required`** | song price                                    |
| `songs[].origin`      | `String` | **`Required`** | song origin                                   |
| `songs[].provider`    | `String` | **`Required`** | song provider                                 |
| `songs[].preview_url` | `String` | **`Required`** | song preview url                              |

---

### **Samples**

#### Success

```json
GET 'http://localhost:8080/api/search?q=alex'
--header 'Content-Type: application/json'
--header 'Authorization: Bearer Token'

HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 31 Jan 2025 08:33:00 GMT

{
    "songs": [
        {
            "id": "",
            "name": "Isn't Life Wonderful? (Alex's original vocal mix)",
            "artist": "Alex D'Elia vs. Ebop Allstars",
            "duration": "",
            "album": "",
            "artwork": "",
            "price": 0,
            "origin": "ChartLyrics",
            "provider": "ChartLyrics",
            "preview_url": "http://www.chartlyrics.com/app/add.aspx?a=108081&t=1835396"
        },
        {
            "id": "",
            "name": "Alex Party (Saturday Night Party)",
            "artist": "Alex Party",
            "duration": "",
            "album": "",
            "artwork": "",
            "price": 0,
            "origin": "ChartLyrics",
            "provider": "ChartLyrics",
            "preview_url": "http://www.chartlyrics.com/app/add.aspx?a=44334&t=9622699"
        },
        {
            "id": "",
            "name": "Don't Give Me Your Life (Classic Alex Party mix)",
            "artist": "Alex Party",
            "duration": "",
            "album": "",
            "artwork": "",
            "price": 0,
            "origin": "ChartLyrics",
            "provider": "ChartLyrics",
            "preview_url": "http://www.chartlyrics.com/app/add.aspx?a=44334&t=10544439"
        }, .....
    ]
}
```

---

[API Overview](../README.md)

[![Tribal Worldwide](../assets/footer.png?style=centerme)](https://tribalworldwide.gt/)
