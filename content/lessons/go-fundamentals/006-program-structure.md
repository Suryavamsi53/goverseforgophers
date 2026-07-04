# Program Structure

Every Go source file (`.go`) strictly follows a specific structural order. Understanding this order is essential for writing readable and compilable Go code.

## The Anatomy of a Go File

A Go source file consists of three main parts, and they must appear in this exact order:

1. **Package Declaration** (Required)
2. **Import Declarations** (Optional)
3. **Top-Level Declarations** (Variables, Constants, Types, and Functions)

---

### 1. Package Declaration
Every Go file must start with a package declaration. It tells the compiler which package the file belongs to. 

```go
package main
```
* If you are building an executable program, the package must be named `main`.
* If you are building a shared library, the package can be named anything (e.g., `package math` or `package stringutil`). All files in the same directory must share the same package name.

### 2. Import Declarations
Immediately following the package declaration are the `import` statements. These declare which other packages your code depends on.

```go
import (
    "fmt"
    "math"
)
```
* Best Practice: Go groups imports into a single block using parentheses. This is called a "factored" import statement.
* **Strict Compiler:** If you import a package but do not use it anywhere in your file, the Go compiler will throw an error and refuse to compile. This keeps Go codebases incredibly clean and free of unused dependencies.

### 3. Top-Level Declarations
Below the imports, you define the actual logic of your package. This includes functions, variables, constants, and custom types. 

These declarations can appear in any order. A function at the top of the file can freely call a function at the bottom of the file; Go does not require forward declarations.

```go
// 1. Constants
const Pi = 3.14159

// 2. Variables
var IsReady bool = true

// 3. Types
type User struct {
    Name string
    Age  int
}

// 4. Functions
func Initialize() {
    fmt.Println("Starting up...")
}
```

---

## Visibility (Exported vs. Unexported)

One of Go's most unique architectural features is how it handles visibility (public vs. private). Go does not use keywords like `public`, `private`, or `protected`.

Instead, Go uses **Capitalization**:

* **Exported (Public):** If a top-level declaration (variable, function, type) starts with a **capital letter**, it is exported. It can be accessed by *other packages*.
  * Example: `fmt.Println()` is capitalized, so we can use it outside the `fmt` package.
  * Example: `type User struct` is exported.

* **Unexported (Private):** If a top-level declaration starts with a **lowercase letter**, it is unexported. It can only be used *within the same package*.
  * Example: `func calculateTax()` is only visible inside its own package.
  * Example: `var activeConnections int` is hidden from the rest of the application.

This simple rule forces developers to think carefully about the public API surface of their packages without cluttering the code with access modifiers.
