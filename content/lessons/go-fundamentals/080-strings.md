# The `strings` Package

Strings in Go are immutable slices of bytes. Because they are immutable, manipulating them requires allocating new memory. 

The `strings` package in the standard library provides highly optimized routines for string manipulation.

## 1. Common Utility Functions

The package contains dozens of functions for searching, replacing, and mutating text.

```go
import (
    "fmt"
    "strings"
)

func main() {
    text := "Go is fast, Go is safe."

    // Contains
    fmt.Println(strings.Contains(text, "fast")) // true
    
    // Count
    fmt.Println(strings.Count(text, "Go")) // 2
    
    // Replace (The '2' means replace up to 2 instances. -1 means replace all)
    fmt.Println(strings.Replace(text, "Go", "Golang", 2))
    
    // Split and Join
    parts := strings.Split("a,b,c", ",") // ["a", "b", "c"]
    joined := strings.Join(parts, "-")   // "a-b-c"
    
    // Casing
    fmt.Println(strings.ToUpper("hello")) // "HELLO"
}
```

## 2. The String Concatenation Trap

Because strings are immutable, adding two strings together (`a + b`) creates a brand new string in memory and copies the bytes over.

If you are building a large string inside a loop, using `+=` is a catastrophic performance trap.

**❌ BAD (O(N²) Time Complexity):**
```go
var result string
for i := 0; i < 10000; i++ {
    // Every loop allocates a new, larger string in RAM and copies everything!
    result += "A" 
}
```

## 3. Enter `strings.Builder`

To efficiently build strings, you must use `strings.Builder`. 

Under the hood, `strings.Builder` uses a dynamically growing `[]byte` slice (just like `append`). It writes data directly into the byte slice, avoiding memory reallocation overhead, and finally converts it to a string at the very end with zero-copy magic.

**✅ GOOD (O(N) Time Complexity):**
```go
import (
    "fmt"
    "strings"
)

func main() {
    var builder strings.Builder
    
    // Optional: Pre-allocate memory if you know the final size!
    builder.Grow(10000) 

    for i := 0; i < 10000; i++ {
        // Writes directly to the underlying byte slice
        builder.WriteString("A") 
    }
    
    // Zero-allocation conversion to string
    finalString := builder.String() 
}
```
**Architecture Insight**: Whenever you are generating dynamic HTML, constructing large SQL queries, or formatting log outputs, **always** use `strings.Builder`.
