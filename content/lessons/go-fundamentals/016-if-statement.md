# if Statement

The `if` statement allows you to execute blocks of code conditionally. Go's `if` statements are very similar to other languages, but with a few strictly enforced formatting rules to keep code clean.

## 1. Basic Syntax

In Go, you **do not** put parentheses `()` around the condition, but the curly braces `{}` are **mandatory**, even if the block contains only one line of code.

```go
age := 20

if age >= 18 {
    fmt.Println("You are an adult.")
}
```

## 2. if-else and else-if

You can chain multiple conditions together. 
*CRITICAL RULE:* The `else` keyword **must** be on the exact same line as the closing `}` of the previous block. Go will not compile if you put `else` on a new line!

```go
score := 85

if score >= 90 {
    fmt.Println("Grade: A")
} else if score >= 80 {
    fmt.Println("Grade: B")
} else {
    fmt.Println("Grade: C")
}
```

## 3. `if` with a Short Statement

Go allows you to execute a short statement immediately before the condition evaluates. This is an extremely common, idiomatic pattern used heavily for error handling or map lookups.

Any variables declared in this short statement are **only** in scope inside the `if` and `else` blocks. They disappear once the statement finishes.

```go
// syntax: if [initialization]; [condition] { ... }

if user, isActive := getUser(); isActive {
    // Both 'user' and 'isActive' exist here
    fmt.Println("Welcome back,", user.Name)
} else {
    // They also exist here!
    fmt.Println("Account is inactive.")
}

// fmt.Println(user.Name) // ERROR: user is undefined here
```

This pattern keeps variable scopes as tight and clean as possible, minimizing memory usage and preventing accidental misuse of variables later in the function.
