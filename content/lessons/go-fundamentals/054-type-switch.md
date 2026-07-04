# Type Switch

If you have an `any` variable that could be one of five different types (e.g., when parsing an unknown JSON payload), writing five separate `if val, ok := data.(Type)` blocks gets extremely messy.

Go provides a cleaner construct: the **Type Switch**.

## 1. Syntax

A type switch looks exactly like a standard `switch` statement, but instead of comparing values, it compares **types**. You use the special syntax `.(type)` inside the switch condition.

```go
func processData(data any) {
    // 'v' takes on the concrete type and value of the matched case
    switch v := data.(type) {
    case int:
        fmt.Printf("It's an integer! Double it: %d\n", v*2)
    case string:
        fmt.Printf("It's a string! Length: %d\n", len(v))
    case bool:
        fmt.Printf("It's a boolean! Inverted: %t\n", !v)
    default:
        fmt.Printf("Unknown type: %T\n", v)
    }
}

func main() {
    processData(42)
    processData("Hello")
    processData(true)
}
```

## 2. Multiple Types per Case

Just like a normal switch, you can check for multiple types in a single case by separating them with commas.

*Warning:* If you match multiple types in one case, the variable `v` remains an `any` interface, because the compiler doesn't know which concrete type it actually ended up being!

```go
func check(data any) {
    switch v := data.(type) {
    case int, float64:
        // 'v' is still of type 'any' here!
        // We cannot do v * 2 because the compiler doesn't know if it's an int or float64.
        fmt.Println("It is a number")
    case string:
        // 'v' is safely a string here
        fmt.Println(len(v))
    }
}
```

## 3. Real-World Use Case: Custom JSON Unmarshaling

Type switches are the backbone of writing dynamic parsers. If you read an unstructured JSON document into a `map[string]any`, you use a type switch to dynamically route the nested data to the correct processing pipelines based on whether it is a string, a nested array `[]any`, or a nested object `map[string]any`.
