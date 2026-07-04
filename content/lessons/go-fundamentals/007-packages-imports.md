# Packages and Imports

Go programs are constructed by linking together packages. A package is simply a directory containing one or more Go source files (`.go`) that are compiled together. 

Packages serve two primary purposes:
1. **Code Reusability**: They allow you to modularize your code into logical, reusable components.
2. **Namespace Management**: They prevent naming collisions (e.g., `math.Max()` vs `custom.Max()`).

---

## 1. Importing Packages

To use code from another package, you must import it. Go supports importing packages from the standard library, external modules, and your own local project.

### Factored Imports
The convention in Go is to use a single, factored `import` statement to group all imports together:

```go
import (
    "fmt"                           // Standard Library
    "math/rand"                     // Standard Library (nested)
    
    "github.com/google/uuid"        // External Module
    
    "github.com/yourusername/app/db" // Local Package
)
```

### The Unused Import Error
Go is famously strict about unused imports. If you import a package (like `"fmt"`) but do not use any of its exported functions or variables in your code, the compiler will throw an error and refuse to compile. 

This guarantees that Go binaries do not bloat over time with unused dependencies.

---

## 2. Advanced Import Techniques

Sometimes, standard imports aren't enough. Go provides three special import syntaxes for edge cases.

### Aliasing Imports
If two packages have the same name (e.g., `math/rand` and `crypto/rand`), or if a package name is too long, you can alias it by providing a name immediately before the import path.

```go
import (
    "crypto/rand"
    mrand "math/rand" // Aliased to 'mrand'
)

func main() {
    // Usage:
    // rand.Read(...)
    // mrand.Intn(10)
}
```

### Dot Imports (Use with Caution)
You can use a period (`.`) to import all exported identifiers from a package directly into your current namespace. 

```go
import (
    . "fmt"
)

func main() {
    Println("Hello without fmt prefix!") // Usually requires fmt.Println
}
```
*Warning: Dot imports are generally discouraged in production code because they make it difficult to determine where a function originated, destroying the namespace benefits of packages. They are mostly used in testing.*

### Blank Imports
Sometimes you need to import a package *only* for its side effects (like executing its `init()` function or registering a database driver), but you don't intend to call any of its functions directly. 

Because the Go compiler rejects unused imports, you must use the blank identifier (`_`) to tell the compiler you are intentionally importing it for side effects.

```go
import (
    "database/sql"
    _ "github.com/lib/pq" // Postgres driver registers itself with database/sql
)
```

---

## 3. Package Initialization (`init`)

Every package can optionally define one or more `init()` functions. 

The `init()` function is special: it takes no arguments, returns no values, and is **executed automatically** before the `main()` function starts.

```go
package config

import "fmt"

var DefaultPort string

func init() {
    // This runs automatically when the package is imported
    DefaultPort = "8080"
    fmt.Println("Config package initialized!")
}
```

**Order of Execution:**
1. If package `A` imports package `B`, the `init()` functions in package `B` run first.
2. Next, package level variables in `A` are evaluated.
3. Next, the `init()` functions in `A` run.
4. Finally, `main()` is executed.
