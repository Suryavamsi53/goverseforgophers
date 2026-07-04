# Constants

Constants represent fixed values that cannot be changed after they are declared. They are evaluated at compile time, not runtime, making them extremely efficient.

## 1. Declaring Constants

Constants are declared using the `const` keyword. They can be character, string, boolean, or numeric values.

```go
const Pi = 3.14159
const AppName = "GoVerse"
const IsProduction = true
```

Unlike variables, you **cannot** use the short declaration syntax (`:=`) with constants.

## 2. Typed vs. Untyped Constants

Go constants are unique because they can be "untyped." 

```go
const MaxUsers = 1000       // Untyped integer constant
const MaxUsersTyped int = 1000 // Typed integer constant
```

An untyped constant is highly flexible. It doesn't have a strict type until it is actually used in a context that requires one. This allows you to mix untyped constants with different numeric types without needing explicit conversions.

```go
const N = 100 // Untyped

var i int = N
var f float64 = N // Works perfectly!
```

## 3. Grouping Constants

Like variables, you can group constants into a block:

```go
const (
    StatusOk       = 200
    StatusNotFound = 404
    StatusError    = 500
)
```

## 4. `iota` (Enumerations)

Go does not have a dedicated `enum` keyword. Instead, it provides `iota`, a special identifier used in `const` blocks to generate auto-incrementing numbers.

`iota` starts at `0` and increments by 1 for each line in the block.

```go
const (
    Sunday = iota    // 0
    Monday           // 1
    Tuesday          // 2
    Wednesday        // 3
    Thursday         // 4
    Friday           // 5
    Saturday         // 6
)
```

You can use math with `iota` to create complex bitmasks or specific sequences:

```go
const (
    ReadPermission = 1 << iota // 1 (1 << 0)
    WritePermission            // 2 (1 << 1)
    ExecutePermission          // 4 (1 << 2)
)
```
