# Unit Testing

Unlike most languages that rely on third-party frameworks like `JUnit` or `pytest`, Go has a world-class testing framework built directly into the standard library: the `testing` package.

## 1. The Rules of Testing

To write a test in Go, you must follow three strict rules:
1. The file name must end in `_test.go` (e.g., `math_test.go`). The Go compiler will automatically ignore these files when building the production binary.
2. The function name must start with `Test` followed by a capitalized word (e.g., `TestAdd`).
3. The function must accept a single argument: `(t *testing.T)`.

```go
// math.go
package math

func Add(a, b int) int {
    return a + b
}
```

```go
// math_test.go
package math

import "testing"

func TestAdd(t *testing.T) {
    result := Add(2, 3)
    expected := 5
    
    if result != expected {
        // t.Errorf prints the error and marks the test as Failed
        t.Errorf("expected %d, got %d", expected, result)
    }
}
```

Run your tests using the command line: 
`go test -v` (`-v` enables verbose output).

## 2. Table-Driven Tests (The Go Standard)

If you want to test `Add()` with positive numbers, negative numbers, and zeroes, writing three separate test functions is considered bad practice.

In Go, the idiomatic architecture is **Table-Driven Testing**. You create a slice of anonymous structs (the "table"), where each struct represents one test case, and loop through them.

```go
func TestAddTable(t *testing.T) {
    // 1. Define the Table
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"Positive", 2, 3, 5},
        {"Negative", -1, -1, -2},
        {"Zeroes", 0, 0, 0},
    }

    // 2. Loop through the test cases
    for _, tc := range tests {
        // t.Run executes a sub-test, ensuring a crash in one 
        // doesn't stop the others from running!
        t.Run(tc.name, func(t *testing.T) {
            
            result := Add(tc.a, tc.b)
            
            if result != tc.expected {
                t.Errorf("expected %d, got %d", tc.expected, result)
            }
        })
    }
}
```

## 3. Code Coverage

How much of your codebase is actually covered by unit tests? Go has this built-in too!

Run: `go test -coverprofile=coverage.out`

This will execute your tests and output the exact percentage of lines covered.
To visualize exactly which lines of code you missed, you can generate a beautiful HTML report:

`go tool cover -html=coverage.out`

It will open a browser window highlighting your covered code in green, and untested code in red!
