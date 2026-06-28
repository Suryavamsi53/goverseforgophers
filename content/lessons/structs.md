# Structs in Go

A struct is a collection of fields. It is useful for grouping data together to form records.

## Declaring a Struct

```go
package main

import "fmt"

type User struct {
    Username string
    Email    string
    Age      int
}

func main() {
    u := User{
        Username: "gopher",
        Email:    "gopher@golang.org",
        Age:      10,
    }
    
    fmt.Println(u.Username)
}
```

Structs are the foundation of object-oriented programming in Go.
