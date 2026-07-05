# HTTP Methods and Status Codes

Because REST uses the URL to define the *Noun* (`/users`), it relies entirely on HTTP Methods to define the *Verb* (the Action).

## 1. The 5 Core HTTP Methods

1. **`GET` (Read)**: Fetches a resource. It must be **Safe** (it never modifies data in the database) and **Idempotent** (calling it 1 time or 100 times returns the same result).
   * `GET /users` -> Returns an array of users.
   * `GET /users/42` -> Returns a single user.

2. **`POST` (Create)**: Creates a brand new resource. It is **NOT Idempotent**. If you `POST /users` twice, you will create two different users in the database!

3. **`PUT` (Full Update)**: Replaces an entire resource. It is **Idempotent**. If you `PUT /users/42` sending `{name: "Bob"}`, the database overwrites User 42. If you do it 10 times, the end state of the database is exactly the same!

4. **`PATCH` (Partial Update)**: Updates a specific field. If User 42 has an email and an age, and you `PATCH /users/42` sending only `{age: 30}`, the email remains untouched.

5. **`DELETE` (Destroy)**: Removes a resource. It is **Idempotent**. `DELETE /users/42` deletes the user. If you run it a second time, it usually returns a 404, but the state of the database (the user is gone) remains unchanged.

## 2. HTTP Status Codes (The API Vocabulary)

A major anti-pattern in Go API development is returning a `200 OK` for every request, and putting the real error inside the JSON (`{"status": "error", "message": "Not Found"}`).

You MUST use native HTTP Status Codes. They are divided into 5 classes:

### 2xx (Success)
* **`200 OK`**: The standard success code for `GET`, `PUT`, and `PATCH`.
* **`201 Created`**: The standard success code for a `POST`. (You should also return a `Location` header with the URL of the newly created resource!).
* **`204 No Content`**: The standard success code for a `DELETE`. It means "Success, and I have no JSON body to send back to you."

### 4xx (Client Error - The User messed up)
* **`400 Bad Request`**: The JSON payload was malformed, or a validation failed (e.g., "Email is required").
* **`401 Unauthorized`**: "Who are you?" (The JWT token is missing or expired).
* **`403 Forbidden`**: "I know who you are, but you aren't allowed to do this." (RBAC failure).
* **`404 Not Found`**: The resource (`/users/999`) does not exist.
* **`409 Conflict`**: Trying to `POST` a user with an email that already exists in the database.
* **`429 Too Many Requests`**: The user hit the Rate Limiter.

### 5xx (Server Error - You messed up)
* **`500 Internal Server Error`**: Your Go code suffered a `panic`, or the database query failed. The user did nothing wrong.
* **`502 Bad Gateway`**: The Go API is trying to talk to the Billing Service, but the Billing Service is offline.
* **`504 Gateway Timeout`**: The Billing Service is online, but it took too long to respond.
