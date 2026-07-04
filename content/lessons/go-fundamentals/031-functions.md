# Functions

Functions in Go are first-class citizens. They can be assigned to variables, passed as arguments, and returned from other functions. 

## 1. Syntax

A function is declared using the `func` keyword, followed by the name, parameters, return types, and the body.

```go
// func name(params) returnType { body }
func add(x int, y int) int {
    return x + y
}
```

### Parameter Shorthand
If consecutive parameters share the exact same type, you can omit the type for all but the last one.

```go
// x and y are both integers
func multiply(x, y, z int) int {
    return x * y * z
}
```

## 2. No Function Overloading

In Java or C++, you can have multiple functions with the same name as long as their parameters are different (Function Overloading). 

```java
// Java
int add(int a, int b) { ... }
float add(float a, float b) { ... }
```

**Go absolutely forbids function overloading.** 

```go
// Go
func addInts(a, b int) int { ... }
func addFloats(a, b float64) float64 { ... }
```
Why? The Go designers believed that overloading makes code confusing to read and drastically slows down compiler speeds. By forcing unique names for every function, the code remains explicitly clear. *(Note: Generics, introduced in Go 1.18, solve the problem of writing duplicate logic for different types).*

## 3. Functions as Types

Because functions are first-class citizens, you can declare them as a Type. This is heavily used in Go's standard library (like the `http.HandlerFunc` type used for web servers).

```go
// Define a function signature as a Type
type MathOperation func(int, int) int

func execute(op MathOperation, a int, b int) {
    result := op(a, b)
    fmt.Println("Result:", result)
}

func main() {
    execute(add, 10, 5) // Passes the 'add' function as a variable!
}
```

### 🧠 Architecture Insight: Strategy Pattern
By passing functions as variables into other functions, you are implementing the Strategy Pattern with zero boilerplate. This allows you to dynamically swap out algorithms or logic at runtime without writing complex interfaces or class hierarchies.
