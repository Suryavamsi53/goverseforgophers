# Table-Driven Tests in Go

Testing is not an afterthought in Go; it is a first-class citizen built directly into the language (`go test`). You do not need to install massive 3rd-party assertion frameworks (like Jest or JUnit) to write production-grade tests.

However, writing `if got != want` over and over again for 10 different scenarios leads to incredibly bloated test files. The idiomatic Go solution is the **Table-Driven Test** pattern.

## 1. The Naive Approach

Imagine we are testing a simple discount calculator.

```go
func CalculateDiscount(price float64, isVIP bool) float64 {
    if isVIP { return price * 0.8 } // 20% off
    return price
}
```

The naive way to test this involves duplicating the test logic for every scenario:

```go
func TestCalculateDiscount(t *testing.T) {
    // Scenario 1
    got1 := CalculateDiscount(100.0, true)
    if got1 != 80.0 { t.Errorf("Expected 80, got %v", got1) }

    // Scenario 2
    got2 := CalculateDiscount(100.0, false)
    if got2 != 100.0 { t.Errorf("Expected 100, got %v", got2) }
}
```
If you need to test 15 different edge cases, your test file becomes hundreds of lines long and impossible to read.

## 2. The Table-Driven Approach

A Table-Driven test separates the *Data* (the scenarios) from the *Execution Logic*.
We define a slice of anonymous structs (the "Table") and iterate through it.

```go
func TestCalculateDiscount_Table(t *testing.T) {
    // 1. Define the Table
    tests := []struct {
        name     string  // Description of the scenario
        price    float64 // Input
        isVIP    bool    // Input
        expected float64 // Output
    }{
        {"VIP user gets 20% off", 100.0, true, 80.0},
        {"Standard user pays full", 100.0, false, 100.0},
        {"Free item remains free", 0.0, true, 0.0},
        {"Negative price returns 0", -50.0, false, 0.0}, // Edge case!
    }

    // 2. The Execution Loop
    for _, tt := range tests {
        // t.Run creates a distinct sub-test in the terminal!
        t.Run(tt.name, func(t *testing.T) {
            got := CalculateDiscount(tt.price, tt.isVIP)
            if got != tt.expected {
                // Instantly readable error message
                t.Errorf("got %v, want %v", got, tt.expected)
            }
        })
    }
}
```

## 3. Why Table-Driven Tests are Elite

1. **Adding tests takes 1 second**: If QA reports a new edge case, you do not write any new logic. You literally add a single line to the struct slice.
2. **Sub-tests (`t.Run`)**: By wrapping the execution in `t.Run()`, Go treats every item in the table as an independent test. If scenario #2 fails, scenario #3 and #4 will still execute! 
3. **Targeted Execution**: Because of `t.Run(tt.name)`, you can execute a *single* scenario from the terminal if you are debugging: `go test -run TestCalculateDiscount_Table/VIP_user_gets_20%_off`.

## 4. Map-Driven Tests (Alternative)

Instead of a Slice, you can use a Map where the Key is the test name.

```go
tests := map[string]struct{
    price    float64
    expected float64
}{
    "Standard": {100.0, 100.0},
    "Free":     {0.0, 0.0},
}

for name, tt := range tests {
    t.Run(name, func(t *testing.T) { ... })
}
```
*Warning: Maps in Go are iterated randomly! If your tests rely on a specific order (they shouldn't!), a Map-driven test will fail sporadically. Always use Slices for guaranteed order.*
