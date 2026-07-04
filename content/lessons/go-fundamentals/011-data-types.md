# Data Types

Go is a strongly and statically typed language. Every variable has a specific type that dictates what kind of data it can hold and what operations can be performed on it.

## 1. Basic Types

* **`bool`**: Represents a boolean value, either `true` or `false`.
* **`string`**: Represents an immutable sequence of characters (text).

```go
var isActive bool = true
var greeting string = "Hello, Gophers!"
```

## 2. Integer Types

Go provides both signed (can be negative) and unsigned (strictly positive) integers at various memory sizes.

* **Architecture-dependent**:
  * `int`: 32 or 64 bits (depends on your OS architecture). This is the default and most commonly used integer type.
  * `uint`: Unsigned 32 or 64 bits.
* **Specific sizes (Signed)**: `int8`, `int16`, `int32`, `int64`
* **Specific sizes (Unsigned)**: `uint8`, `uint16`, `uint32`, `uint64`

```go
var age int = 30
var byteCount uint64 = 18446744073709551615
```

## 3. Floating-Point & Complex Types

For numbers with decimal points, Go provides two sizes:
* **`float32`**: 32-bit floating-point number.
* **`float64`**: 64-bit floating-point number (the default and most precise).

```go
var pi float64 = 3.14159265359
```

Go also has built-in support for complex numbers (numbers with real and imaginary parts):
* **`complex64`**, **`complex128`**

## 4. Aliases: `byte` and `rune`

Go has two special type aliases used heavily in text processing:

* **`byte`**: An exact alias for `uint8`. It is used to emphasize that a value is a piece of raw data rather than a small number. Strings in Go are essentially slices of bytes.
* **`rune`**: An exact alias for `int32`. It represents a single Unicode character (a Unicode code point).

```go
var a byte = 'A' // ASCII value 65
var heart rune = '❤' // Unicode value 10084
```

## Summary of Default Types
If you use the short declaration `:=`, Go infers the following default types:
* Whole numbers `:= 10` become `int`.
* Decimals `:= 3.14` become `float64`.
* Text `:= "hello"` becomes `string`.
