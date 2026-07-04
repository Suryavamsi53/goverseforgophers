# REST API Basics

Building a REST API in Go requires combining the `net/http` package with the `encoding/json` package to process HTTP Methods, Status Codes, and JSON payloads.

## 1. Routing and HTTP Methods

In Go 1.22, the `ServeMux` was massively upgraded to support HTTP Methods and Path Variables natively (previously, developers had to use third-party libraries like `chi` or `gorilla/mux`).

```go
func main() {
    mux := http.NewServeMux()
    
    // Exact Method Matches
    mux.HandleFunc("GET /api/users", handleGetUsers)
    mux.HandleFunc("POST /api/users", handleCreateUser)
    
    // Path Variables (e.g., /api/users/123)
    mux.HandleFunc("GET /api/users/{id}", handleGetUserByID)

    http.ListenAndServe(":8080", mux)
}
```
Inside the handler, you can extract the wildcard variable using `r.PathValue()`:
```go
func handleGetUserByID(w http.ResponseWriter, r *http.Request) {
    userID := r.PathValue("id")
    fmt.Fprintf(w, "Fetching details for User %s", userID)
}
```

## 2. Parsing JSON Payloads (POST)

When a client sends data to create a user, we must parse the `r.Body`. Remember to use `json.NewDecoder` to stream the data efficiently without massive memory allocations!

```go
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

func handleCreateUser(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    
    // 1. Decode the request body
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
        return
    }
    
    // ... Save to database ...

    // 2. Return a success status
    w.WriteHeader(http.StatusCreated) // 201
    fmt.Fprintf(w, "User %s created successfully!", req.Name)
}
```

## 3. Returning JSON Responses (GET)

To return JSON to the client, we must set the correct `Content-Type` header, and then use `json.NewEncoder` to write directly to the `ResponseWriter`.

```go
type UserResponse struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
}

func handleGetUsers(w http.ResponseWriter, r *http.Request) {
    users := []UserResponse{
        {ID: 1, Name: "Alice"},
        {ID: 2, Name: "Bob"},
    }

    // 1. Set the correct headers
    w.Header().Set("Content-Type", "application/json")
    
    // 2. Set the status code (must happen AFTER setting headers!)
    w.WriteHeader(http.StatusOK) // 200

    // 3. Encode the slice directly to the network socket
    json.NewEncoder(w).Encode(users)
}
```
