# Anonymous Structs

Sometimes you need to group data together temporarily, but you don't want to pollute your package namespace by declaring a formal `type MyStruct struct`.

Go allows you to declare and instantiate **Anonymous Structs** inline.

## 1. Syntax

You define the struct shape and provide its values in a single, combined block.

```go
// Define and instantiate immediately
serverConfig := struct {
    Host string
    Port int
}{
    Host: "localhost",
    Port: 8080,
}

fmt.Println(serverConfig.Host)
```

## 2. Real-World Use Case A: JSON Parsing

When interacting with external APIs, you often receive massive JSON payloads, but you only care about extracting two or three specific fields. 

Instead of creating a permanent struct for a one-off API call, use an anonymous struct:

```go
jsonPayload := `{"name": "Alice", "age": 25, "secret_token": "xyz", "metadata": {}}`

// We only care about name and age
var user struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

// json.Unmarshal will safely ignore "secret_token" and "metadata"
json.Unmarshal([]byte(jsonPayload), &user)

fmt.Println(user.Name) // "Alice"
```

## 3. Real-World Use Case B: Table-Driven Tests

If you write tests in Go, you will use anonymous structs constantly. "Table-Driven Testing" is the standard idiomatic way to write unit tests in Go.

You create a slice of anonymous structs, where each struct represents a single test case (input and expected output).

```go
func TestAdd(t *testing.T) {
    // A slice of anonymous structs!
    tests := []struct {
        name     string
        inputA   int
        inputB   int
        expected int
    }{
        {"Positive numbers", 2, 2, 4},
        {"Negative numbers", -1, -1, -2},
        {"Zero values", 0, 0, 0},
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            result := add(tc.inputA, tc.inputB)
            if result != tc.expected {
                t.Errorf("expected %d, got %d", tc.expected, result)
            }
        })
    }
}
```
