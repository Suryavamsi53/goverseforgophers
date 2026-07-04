# Output Formatting

Formatting output properly is essential for building readable CLIs and writing detailed logs. The `fmt` (format) package provides a family of functions to print data to the console and format strings dynamically.

## 1. Print vs Println

* `fmt.Print`: Prints arguments to the console as-is.
* `fmt.Println`: Prints arguments, inserts spaces between them, and appends a newline at the end.

```go
fmt.Print("Hello", "World")   // Output: HelloWorld
fmt.Println("Hello", "World") // Output: Hello World\n
```

## 2. Printf (Print Formatted)

`fmt.Printf` is the most powerful printing function. It allows you to inject variables into a string using format specifiers (called "verbs"). 

*Note: `Printf` does not append a newline automatically, so you usually need to add `\n` at the end.*

```go
name := "Alice"
age := 25
fmt.Printf("User %s is %d years old.\n", name, age)
```

### Essential Formatting Verbs

| Verb | Description | Example | Output |
| :---: | :--- | :--- | :--- |
| `%v` | The default value format (Works for any type) | `Printf("%v", 42)` | `42` |
| `%+v` | Adds field names when printing structs | `Printf("%+v", user)` | `{Name:Alice Age:25}` |
| `%#v` | Go-syntax representation (great for debugging) | `Printf("%#v", "test")` | `"test"` |
| `%T` | Prints the **Type** of the variable | `Printf("%T", 3.14)` | `float64` |
| `%t` | Boolean | `Printf("%t", true)` | `true` |
| `%d` | Base-10 Integer | `Printf("%d", 100)` | `100` |
| `%f` | Floating-point (decimal) | `Printf("%.2f", 3.1415)`| `3.14` |
| `%s` | String | `Printf("%s", "Gopher")` | `Gopher` |
| `%q` | Double-quoted string safely escaped | `Printf("%q", "Hi")` | `"Hi"` |

## 3. Sprintf (String Print Formatted)

Sometimes you want to format a string but **store it in a variable** instead of printing it to the console immediately. You can use `fmt.Sprintf` for this. It takes the exact same verbs as `Printf`.

```go
price := 19.99
currency := "USD"

// Constructs the string and returns it
receipt := fmt.Sprintf("Total: %.2f %s", price, currency)

// You can now use the 'receipt' string later
// fmt.Println(receipt)
```
