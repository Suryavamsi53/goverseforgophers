# Go Basics: Variables and Types

Go is a statically typed language, which means variables always have a specific type and that type cannot change.

## Variables

You can declare variables using the `var` keyword:

```go
package main

import "fmt"

func main() {
    var name string = "GoVerse"
    var age int = 1
    
    fmt.Printf("Welcome to %s, year %d\n", name, age)
}
```

## Short Variable Declaration

Inside a function, you can use the `:=` short assignment statement:

```go
func main() {
    name := "GoVerse" // Type inferred as string
    age := 1          // Type inferred as int
}
```
