# Type Conversion

In many languages (like C or Java), the compiler will automatically convert an `int` to a `float` if you try to add them together. This is called *implicit type conversion* or *coercion*.

**Go does NOT support implicit type conversion.** 

Go's design prioritizes clarity and predictability over saving a few keystrokes. If you want to mix types, you must explicitly convert them.

## 1. Basic Type Conversion

To convert a value to a different type, wrap the value in the target type's name like a function call: `Type(value)`.

```go
var i int = 42
var f float64 = float64(i)  // Convert int to float64
var u uint = uint(f)        // Convert float64 to uint
```

If you try to perform math on mismatched types without conversion, the compiler will catch it:

```go
var x int = 10
var y float64 = 2.5

// sum := x + y // ERROR: invalid operation: mismatched types int and float64
sum := float64(x) + y // Correct!
```

## 2. Converting Strings to Bytes (and vice versa)

Strings in Go are essentially read-only slices of bytes. Converting between a `string` and a `[]byte` (byte slice) is extremely common, especially when dealing with files or network requests.

```go
// String to Byte Slice
greeting := "Hello"
bytes := []byte(greeting) // [72 101 108 108 111]

// Byte Slice to String
str := string(bytes) // "Hello"
```

## 3. Advanced Conversions (The `strconv` package)

You cannot convert an `int` to a `string` using `string(42)`. Doing so will try to find the Unicode character at position 42 (which is `*`), rather than returning the text `"42"`.

To convert between strings and numeric types, you must use the `strconv` (String Conversions) package from the standard library.

```go
import "strconv"

func main() {
    // String to Int (Atoi = ASCII to Integer)
    num, err := strconv.Atoi("100")
    
    // Int to String (Itoa = Integer to ASCII)
    text := strconv.Itoa(100)
    
    // String to Float
    pi, err := strconv.ParseFloat("3.14", 64)
}
```
