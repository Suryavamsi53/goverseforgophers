# Variables

Go is a statically typed language, meaning variables must have a specific type, and that type cannot change over time. However, Go provides several ways to declare variables, ranging from explicit to highly concise.

## 1. The `var` Keyword

The most explicit way to declare a variable is using the `var` keyword, followed by the variable name, and then the type.

```go
var age int = 25
var name string = "Alice"
```

If you provide an initial value, Go can infer the type. You can omit the type declaration:

```go
var age = 25       // Go infers 'int'
var name = "Alice" // Go infers 'string'
```

## 2. Short Variable Declaration (`:=`)

Inside a function, you can use the short variable declaration operator `:=`. This replaces the `var` keyword and the explicit type. 

This is the most common way to declare variables in Go.

```go
func main() {
    age := 25
    name := "Alice"
    isActive := true
}
```
*Note: `:=` can only be used inside functions. Top-level (package-level) variables must use the `var` keyword.*

## 3. Multiple Declarations

You can declare multiple variables on the same line to keep your code compact.

```go
// Using var
var x, y, z int = 1, 2, 3

// Using short declaration
width, height := 100, 200
```

You can also group `var` declarations into a block, which is common for package-level variables:

```go
var (
    appVersion = "1.0.0"
    maxUsers   = 5000
    debugMode  = false
)
```

## 4. Zero Values

One of Go's safest features is that **variables are never uninitialized**. If you declare a variable without giving it a value, Go automatically assigns it a "zero value" based on its type.

```go
var i int     // 0
var f float64 // 0.0
var b bool    // false
var s string  // "" (empty string)
```

Because of this, you will never encounter "undefined" variables or garbage memory values in Go.
