# Multiple Return Values

One of Go's most distinctive and defining features is the ability to return multiple values from a single function natively. 

In languages like C or Java, returning multiple values requires returning complex arrays or creating temporary "wrapper" objects. Go eliminates this boilerplate entirely.

## 1. Syntax

To return multiple values, wrap the return types in parentheses `()`.

```go
// Returns both a string and an integer
func getProfile() (string, int) {
    return "Alice", 25
}

func main() {
    name, age := getProfile()
    fmt.Printf("Name: %s, Age: %d\n", name, age)
}
```

## 2. The `(value, error)` Paradigm

The absolute most common use case for multiple returns in Go is error handling. Because Go does not have `try/catch` exception blocks, functions that can fail return the result as the first value, and an `error` interface as the second value.

```go
import (
    "fmt"
    "strconv"
)

func main() {
    // Atoi returns (int, error)
    num, err := strconv.Atoi("100abc")
    
    if err != nil {
        fmt.Println("Failed to parse number:", err)
        return
    }
    
    fmt.Println("Parsed successfully:", num)
}
```
*Architecture Insight: This forces developers to acknowledge and handle errors immediately where they occur, rather than letting a hidden exception crash the system 10 layers up the call stack.*

## 3. The Blank Identifier (`_`)

If a function returns multiple values, but you only care about one of them, you **must** ignore the unused values using the blank identifier (`_`). 

Because Go refuses to compile unused variables, the `_` tells the compiler to discard the value immediately.

```go
// We only want the integer, we don't care about the error
num, _ := strconv.Atoi("42")
```
*Warning: Ignoring errors using `_` in production code is generally a terrible idea. It is often referred to as "swallowing" the error and leads to silent failures.*
