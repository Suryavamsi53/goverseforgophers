# switch Statement

The `switch` statement is a cleaner way to write a long sequence of `if-else if` statements. Go's switch is much safer and more flexible than switch statements in languages like C or Java.

## 1. Basic Syntax and No Fallthrough

A fundamental difference in Go is that **cases do not fall through by default**. In C or Java, if you forget a `break` statement at the end of a case, execution bleeds into the next case. Go automatically breaks out of the switch as soon as a case succeeds.

```go
os := "darwin"

switch os {
case "darwin":
    fmt.Println("macOS") // Code stops here. No 'break' needed.
case "linux":
    fmt.Println("Linux")
default:
    fmt.Println("Unknown OS")
}
```

## 2. Multiple Values in One Case

You can test a variable against multiple values on the exact same line by separating them with commas.

```go
day := "Saturday"

switch day {
case "Saturday", "Sunday":
    fmt.Println("It's the weekend!")
case "Monday", "Tuesday", "Wednesday", "Thursday", "Friday":
    fmt.Println("It's a weekday.")
}
```

## 3. Switch with No Condition

If you omit the variable after the `switch` keyword, the switch evaluates as `switch true`. 

This allows you to evaluate completely different boolean conditions in each case, making it a very clean replacement for long, messy `if-else` chains.

```go
score := 85

switch {
case score >= 90:
    fmt.Println("A")
case score >= 80:
    fmt.Println("B")
case score >= 70:
    fmt.Println("C")
default:
    fmt.Println("F")
}
```

## 4. The `fallthrough` Keyword

If you specifically *want* the C-style behavior where execution spills over into the next case block regardless of whether it evaluates to true, you can use the `fallthrough` keyword at the end of a case block.

```go
num := 2

switch num {
case 1:
    fmt.Println("One")
case 2:
    fmt.Println("Two")
    fallthrough // Forces execution of the next case
case 3:
    fmt.Println("Three")
}

// Output:
// Two
// Three
```
