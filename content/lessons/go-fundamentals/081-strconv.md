# The `strconv` Package

Converting data types (like turning a string into an integer) is one of the most common tasks in programming, especially when reading from databases, API URLs, or JSON payloads.

The `strconv` (String Conversions) package provides highly optimized functions for this exact purpose.

## 1. String to Integer (`Atoi`)

The most famous function in the package is `Atoi` (ASCII to Integer). Because the conversion can fail (e.g., trying to parse `"abc"` into a number), it returns an `(int, error)` pair.

```go
import (
    "fmt"
    "strconv"
)

func main() {
    num, err := strconv.Atoi("42")
    if err != nil {
        fmt.Println("Failed to parse:", err)
        return
    }
    fmt.Println(num * 2) // Outputs: 84
}
```

## 2. Integer to String (`Itoa`)

The reverse operation is `Itoa` (Integer to ASCII). Because converting an integer to a string can never fail, it only returns a single `string` value.

```go
age := 25
text := "I am " + strconv.Itoa(age) + " years old."
fmt.Println(text)
```

## 3. The `Parse` and `Format` Families

`Atoi` and `Itoa` are convenient wrappers, but under the hood, they use the highly robust `Parse` and `Format` families, which allow you to specify bases (like Hexadecimal) and bit sizes (like `int8` vs `int64`).

* **Parsing (String -> Type):**
  * `ParseBool("true")`
  * `ParseFloat("3.14", 64)`
  * `ParseInt("-42", 10, 64)` (Base 10, 64-bit size)

* **Formatting (Type -> String):**
  * `FormatBool(true)`
  * `FormatFloat(3.1415, 'f', 2, 64)` (Format as float, 2 decimal places)
  * `FormatInt(-42, 10)`

## 4. Performance vs `fmt.Sprintf`

You might be wondering: *Why use `strconv.Itoa(42)` when I can just use `fmt.Sprintf("%d", 42)`?*

**Performance.**

The `fmt` package uses **Reflection** (`reflect`) under the hood to dynamically inspect the data types of the arguments you pass it at runtime. Reflection is incredibly slow and allocates memory on the Heap.

The `strconv` package, on the other hand, is hardcoded for specific types. It bypasses reflection entirely. 

If you are formatting thousands of numbers inside a loop (like generating a massive CSV file), `strconv.Itoa` is approximately **400% faster** and uses significantly less memory than `fmt.Sprintf`!
