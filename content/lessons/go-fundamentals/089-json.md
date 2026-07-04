# JSON Processing

JSON is the universal language of the web. Go provides the `encoding/json` package to serialize and deserialize data.

## 1. Struct Tags and Marshalling

To convert a Go struct into a JSON string, use `json.Marshal()`.

You use "Struct Tags" (the string literals enclosed in backticks) to tell the JSON encoder exactly what to name the keys, or when to ignore them.

```go
import (
    "encoding/json"
    "fmt"
)

type User struct {
    // Mapped to exactly "user_id"
    ID int `json:"user_id"` 
    
    // Mapped to "name". If empty, omit the key entirely from the JSON
    Name string `json:"name,omitempty"` 
    
    // The "-" tells the encoder to COMPLETELY ignore this field (Great for passwords!)
    Password string `json:"-"` 
    
    // Unexported fields (lowercase) are automatically ignored by the JSON encoder
    internalState string 
}

func main() {
    u := User{ID: 1, Name: "Alice", Password: "super_secret"}
    
    // Marshal returns a []byte
    jsonData, _ := json.Marshal(u)
    
    fmt.Println(string(jsonData)) 
    // Output: {"user_id":1,"name":"Alice"}
}
```

## 2. Unmarshalling (JSON to Struct)

To convert JSON back into a Go struct, use `json.Unmarshal()`. You must pass a **pointer** to the struct so the package can inject the data into it.

```go
func main() {
    payload := `{"user_id": 99, "name": "Bob"}`
    
    var u User
    
    // Pass the pointer &u!
    err := json.Unmarshal([]byte(payload), &u)
    if err != nil {
        fmt.Println("Invalid JSON:", err)
    }
    
    fmt.Println(u.Name) // Bob
}
```

## 3. Dynamic JSON (Unknown Structures)

What if you receive a JSON payload, but you don't know the schema? You can unmarshal it directly into a map of `any`.

```go
payload := `{"status": "ok", "retries": 3}`

var data map[string]any
json.Unmarshal([]byte(payload), &data)

// You must use Type Assertions to extract the 'any' values
status := data["status"].(string)
retries := data["retries"].(float64) // JSON numbers always unmarshal to float64!
```

## 4. The Enterprise Way: Stream Encoding/Decoding

`json.Marshal` and `Unmarshal` require you to hold the entire JSON string in memory. If a client sends a 50MB JSON payload to your web server, your server's RAM will spike heavily.

Enterprise applications use `json.NewEncoder` and `json.NewDecoder`. These plug directly into the `io.Reader` and `io.Writer` interfaces!

```go
func handler(w http.ResponseWriter, r *http.Request) {
    var user User
    
    // DECODE: Reads the JSON stream straight off the TCP network socket 
    // directly into the struct without allocating a massive string in memory!
    json.NewDecoder(r.Body).Decode(&user)

    // ... process user ...

    // ENCODE: Streams the JSON response straight back to the client browser!
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
```
