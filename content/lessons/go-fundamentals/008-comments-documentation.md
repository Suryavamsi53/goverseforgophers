# Comments and Documentation

Go provides standard C-style comments, but with a strong emphasis on writing self-documenting code and generating documentation automatically.

## 1. Single-Line and Multi-Line Comments

You can write comments in two ways:

```go
// This is a single-line comment.
// It is the most common way to write comments in Go.
var name = "Gopher"

/*
This is a multi-line comment.
It's mostly used for commenting out large blocks of code during debugging.
*/
var version = 1.0
```

## 2. Doc Comments (GoDoc)

In Go, comments aren't just for developers reading the source code; they are parsed by the `go doc` tool to automatically generate documentation for your packages.

To document a type, variable, constant, or function, simply write a regular single-line comment immediately preceding its declaration, with no blank lines in between. 

By convention, the comment should be a complete sentence that begins with the name of the element being declared.

```go
// User represents a registered member of the system.
type User struct {
    Name string
}

// CalculateTax computes the standard tax rate for a given amount.
// It returns the tax as a float64.
func CalculateTax(amount float64) float64 {
    return amount * 0.20
}
```

If you run `go doc CalculateTax` in your terminal, Go will print out your comment and the function signature.

## 3. Package Documentation

To document a whole package, place a comment directly above the `package` clause. For multi-file packages, this only needs to be present in one file (often conventionally named `doc.go`).

```go
// Package mathutil provides utility functions for advanced trigonometry
// and statistics. It is designed to be highly concurrent.
package mathutil
```

## 4. Best Practices

* **Explain WHY, not WHAT**: Code tells you *what* it does. Comments should explain *why* it does it.
* **Keep it updated**: An incorrect comment is worse than no comment.
* **Capitalization**: Go tools expect doc comments to be proper English sentences, starting with the identifier name and ending with a period.
