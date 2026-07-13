# Go Programming Curriculum — BEGINNER TIER
### Covers Level 1–6: Basics & Syntax, Control Flow, Functions, Arrays/Slices/Maps, Structs & Methods, Interfaces
### Format: Question → Answer shown immediately after

---

# LEVEL 1: Basics & Syntax

**1** What is the zero value of an `int` in Go?
a) null  b) 0  c) undefined  d) Compile error
**Answer: b) 0**

**2** What is the zero value of a `string` in Go?
a) `""`  b) `null`  c) `nil`  d) undefined
**Answer: a) `""` (empty string)**

**3** What is the zero value of a `bool` in Go?
a) `true`  b) `false`  c) `nil`  d) `0`
**Answer: b) `false`**

**4** Which is invalid inside a function body?
a) `var x int = 10`  b) `x := 10`  c) `x := 10` then `x := 20` in same scope  d) `var x = 10`
**Answer: c) — `:=` requires at least one new variable on the left; redeclaring the same single variable with `:=` in the same scope is a compile error**

**5** What does `const Pi = 3.14` create?
a) A typed float64 constant  b) An untyped constant that adapts to context  c) A reassignable variable  d) Compile error
**Answer: b) Untyped constant — takes the type required by context (float32, float64, etc.)**

**6** What happens if you declare `var x int` and never use it?
a) Nothing  b) Compile error  c) Runtime panic  d) Warning only
**Answer: b) Compile error — Go disallows unused local variables**

**7** What is `7 / 2` when both operands are `int`?
a) 3.5  b) 3  c) 4  d) Compile error
**Answer: b) 3 — integer division truncates**

**8** What does `iota` do inside a `const` block?
a) Generates random numbers  b) Auto-increments from 0 per line in the block  c) Marks immutability  d) Invalid syntax
**Answer: b) Auto-increments starting at 0 for each ConstSpec line**

**9** What is the type of `x := 5.0 / 2`?
a) int  b) float64  c) Compile error  d) float32
**Answer: b) float64 — untyped constant `5.0` forces floating-point division**

**10** What does `:=` do that `var` cannot?
a) Declare a constant  b) Infer type, only valid inside functions  c) Allow global declarations  d) Nothing, identical
**Answer: b) Short variable declaration — infers type, can only be used inside function bodies**

**11** Which statement about Go type conversion is true?
a) Implicit numeric conversion is allowed  b) Explicit conversion required, e.g. `float64(x)`  c) String+int auto-converts  d) Only pointers need conversion
**Answer: b) Explicit conversion required — Go has no implicit numeric widening**

**12** What does `x := 10; x := 20` (new scope, e.g. inside an if-block) do?
a) Compile error, x already declared  b) Shadows outer x within the block  c) Reassigns outer x  d) Undefined behavior
**Answer: b) Shadows — a new `x` is created scoped to the inner block**

**13** Which is a valid multiple-variable declaration with different types?
a) `var a, b int, string = 1, "x"`  b) `var ( a int = 1; b string = "x" )`  c) `var a, b = 1, "x"`  d) both b and c
**Answer: d) Both are valid — grouped var block, or `var a, b = 1, "x"` with inferred types**

**14** What is the output of `fmt.Println(10 % 3)`?
a) 3  b) 1  c) 0  d) 3.33
**Answer: b) 1 — modulo remainder**

**15** What does `const MaxRetries = 5` followed by `MaxRetries = 10` do?
a) Reassigns fine  b) Compile error, cannot assign to constant  c) Runtime panic  d) Creates shadow variable
**Answer: b) Compile error — constants are immutable**

**16** In `var b byte = 300`, what happens?
a) b becomes 300  b) Compile error — constant 300 overflows byte  c) b wraps to 44  d) Runtime panic
**Answer: b) Compile error — byte is uint8, max value 255, and this overflow is caught at compile time for constant expressions**

**17** What's the difference between `rune` and `byte` in Go?
a) No difference  b) `rune` is alias for int32 (Unicode code point), `byte` is alias for uint8  c) `byte` is signed, `rune` unsigned  d) `rune` is deprecated
**Answer: b) `rune` = int32 (Unicode code point), `byte` = uint8 (raw byte)**

**18** What does this print? `fmt.Println("go" + "lang")`
a) Compile error  b) "golang"  c) "go lang"  d) 5
**Answer: b) "golang" — string concatenation with `+`**

**19** Predict output:
```go
x := 5
y := &x
*y = 10
fmt.Println(x)
```
a) 5  b) 10  c) Compile error  d) nil
**Answer: b) 10 — `y` points to `x`'s address; dereferencing and assigning mutates `x`**

**20** What's wrong with this code?
```go
func main() {
	var timeout int
	fmt.Println("started")
}
```
a) Nothing, compiles fine  b) `timeout` declared and unused → compile error  c) Missing import  d) `fmt` undefined
**Answer: b) `timeout` is declared but never used**

---

### Level 1 — Coding Problems

**21** Which of the following is the correct way to format and print `maxConnections` (int), `timeoutSeconds` (float64), and `serviceName` (string) using `fmt.Printf`?
a) `fmt.Printf("Service: %s | MaxConns: %d | Timeout: %.1fs\n", serviceName, maxConnections, timeoutSeconds)`
b) `fmt.Printf("Service: %d | MaxConns: %s | Timeout: %.1f\n", serviceName, maxConnections, timeoutSeconds)`
c) `fmt.Printf("Service: %v | MaxConns: %f | Timeout: %d\n", serviceName, maxConnections, timeoutSeconds)`
d) `fmt.Printf("Service: %s | MaxConns: %i | Timeout: %d\n", serviceName, maxConnections, timeoutSeconds)`
**Answer: a) `fmt.Printf("Service: %s | MaxConns: %d | Timeout: %.1fs\n", ...)`**
```go
package main
import "fmt"
func main() {
	maxConnections := 100
	timeoutSeconds := 30.5
	serviceName := "auth-service"
	fmt.Printf("Service: %s | MaxConns: %d | Timeout: %.1fs\n", serviceName, maxConnections, timeoutSeconds)
}
```

**22** How do you properly use `iota` to define byte-size constants (`KB`, `MB`, `GB`) as powers of 1024?
a) `KB = iota * 1024`, `MB`, `GB`
b) `KB = 1 << (10 * iota)`, `MB`, `GB`
c) `KB = iota << 10`, `MB`, `GB`
d) `KB = 1024 ^ iota`, `MB`, `GB`
**Answer: b) `KB = 1 << (10 * iota)`, `MB`, `GB`**
```go
package main
import "fmt"
const (
	_  = iota
	KB = 1 << (10 * iota)
	MB
	GB
)
func main() {
	fmt.Println(KB, MB, GB)
}
```

**23** How do you convert `memoryMB` (int) to GB (float64) using proper float division to avoid integer truncation?
a) `memoryGB := memoryMB / 1024.0`
b) `memoryGB := float64(memoryMB / 1024)`
c) `memoryGB := float64(memoryMB) / 1024`
d) `memoryGB := memoryMB / float32(1024)`
**Answer: c) `memoryGB := float64(memoryMB) / 1024`**
```go
package main
import "fmt"
func main() {
	memoryMB := 2560
	memoryGB := float64(memoryMB) / 1024
	fmt.Printf("%d MB = %.2f GB\n", memoryMB, memoryGB)
}
```

**24** Given `maxRequestsPerMinute = 600`, which correctly computes `requestsPerSecond` as a float64?
a) `requestsPerSecond := float64(maxRequestsPerMinute) / 60.0`
b) `requestsPerSecond := maxRequestsPerMinute / 60`
c) `requestsPerSecond := float64(maxRequestsPerMinute / 60)`
d) `requestsPerSecond := float(maxRequestsPerMinute) / 60`
**Answer: a) `requestsPerSecond := float64(maxRequestsPerMinute) / 60.0`**
```go
package main
import "fmt"
func main() {
	maxRequestsPerMinute := 600
	requestsPerSecond := float64(maxRequestsPerMinute) / 60.0
	fmt.Printf("Requests/sec: %.2f\n", requestsPerSecond)
}
```

**25** Why can't a plain sequential `iota` be used to declare `StatusOK = 200`, `StatusNotFound = 404`, `StatusServerError = 500` directly?
a) `iota` only works with strings, not integers.
b) `iota` generates strictly sequential values (0, 1, 2...) and cannot generate arbitrary non-sequential numbers directly.
c) `iota` resets on every line, making it impossible.
d) `iota` cannot be used inside a `const` block.
**Answer: b) `iota` generates strictly sequential values (0, 1, 2...) and cannot generate arbitrary non-sequential numbers directly.**
```go
package main
import "fmt"
const (
	StatusOK          = 200
	StatusNotFound    = 404
	StatusServerError = 500
)
func main() {
	fmt.Println(StatusOK, StatusNotFound, StatusServerError)
}
```
Explanation: `iota` only auto-increments by a fixed step per line (usually +1). These status codes have non-uniform gaps (200 → 404 → 500), so they must be assigned explicitly rather than derived from `iota`.

---

# LEVEL 2: Control Flow

**26** Does Go require parentheses around `if` conditions?
a) Yes, always  b) No, and braces are mandatory  c) No, but parentheses are also disallowed  d) Optional both ways
**Answer: b) No parens needed, but `{ }` braces are mandatory even for single statements**

**27** What does Go's `switch` do differently from C/Java by default?
a) Nothing, identical  b) No fallthrough by default — each case breaks automatically  c) Requires `break` explicitly  d) Only works on integers
**Answer: b) Cases don't fall through unless you use the `fallthrough` keyword**

**28** Which loop construct does Go NOT have?
a) `for`  b) `while`  c) `do-while`  d) both b and c
**Answer: d) Go only has `for` — no dedicated `while` or `do-while` keywords (achieved via `for` variants)**

**29** What does `for i := 0; i < 5; i++ {}` with empty body do?
a) Compile error  b) Infinite loop  c) Runs 5 times doing nothing, valid  d) Runs once
**Answer: c) Valid — loop runs 5 times with an empty body**

**30** What does `for { }` (no condition) do?
a) Compile error  b) Infinite loop  c) Runs zero times  d) Runs once
**Answer: b) Infinite loop, equivalent to `while(true)`**

**31** In a `switch` statement, what does a case with multiple values look like?
a) `case 1, 2, 3:`  b) `case 1 | 2 | 3:`  c) `case (1,2,3):`  d) Not supported
**Answer: a) `case 1, 2, 3:` — comma-separated values in one case**

**32** What is a "switch with no expression" used for?
a) Invalid syntax  b) Acts like an if-else chain, each case is a boolean condition  c) Only works on strings  d) Same as regular switch
**Answer: b) `switch { case x > 10: ... }` — clean alternative to long if-else chains**

**33** What does `continue` do inside a nested loop by default?
a) Breaks all loops  b) Skips to next iteration of innermost enclosing loop  c) Compile error in nested loops  d) Skips to next iteration of outer loop
**Answer: b) Affects only the innermost loop unless a label is used**

**34** How do you break out of an outer loop from within a nested inner loop?
a) `break 2`  b) Using a labeled break: `break OuterLoop`  c) Not possible in Go  d) `return`
**Answer: b) Labeled break, e.g. `OuterLoop: for {...}` then `break OuterLoop`**

**35** What does `range` return when iterating over a slice?
a) Only the value  b) Only the index  c) Index and value  d) A pointer to each element
**Answer: c) Index, value (in that order) — `for i, v := range slice`**

**36** What happens if you modify a slice element via `for _, v := range slice { v = 100 }`?
a) Modifies the original slice  b) `v` is a copy; original slice unchanged  c) Compile error  d) Panic
**Answer: b) `v` is a copy of the value — mutating it does not affect the underlying slice**

**37** Predict the output:
```go
switch x := 5; {
case x > 10:
	fmt.Println("big")
case x > 3:
	fmt.Println("medium")
default:
	fmt.Println("small")
}
```
a) big  b) medium  c) small  d) Compile error
**Answer: b) medium — first matching case wins, no fallthrough**

**38** What does `fallthrough` do in a switch case?
a) Nothing, invalid  b) Forces execution to continue into the next case block regardless of its condition  c) Skips remaining cases  d) Restarts the switch
**Answer: b) Forces the next case's body to execute unconditionally**

**39** Predict output:
```go
for i := 0; i < 3; i++ {
	if i == 1 {
		continue
	}
	fmt.Println(i)
}
```
a) 0 1 2  b) 0 2  c) 0 1  d) 1 2
**Answer: b) 0 2 — iteration 1 is skipped via continue**

**40** What's the idiomatic Go way to iterate over a map's keys and values?
a) `for k, v := range myMap`  b) `for k in myMap`  c) `myMap.each(...)`  d) Maps can't be iterated
**Answer: a) `for k, v := range myMap`**

**41** Is map iteration order guaranteed in Go?
a) Yes, insertion order  b) Yes, sorted key order  c) No, deliberately randomized  d) Yes, but only for small maps
**Answer: c) No — Go intentionally randomizes map iteration order to prevent reliance on it**

**42** What does this print?
```go
count := 0
for count < 5 {
	count++
}
fmt.Println(count)
```
a) 4  b) 5  c) Infinite loop  d) 0
**Answer: b) 5 — this is Go's "while-style" for loop (condition-only)**

**43** What is wrong (if anything) with:
```go
if x := getValue(); x > 0 {
	fmt.Println(x)
}
fmt.Println(x)
```
a) Nothing wrong  b) Compile error — `x` is scoped to the if statement and unavailable outside  c) Runtime panic  d) x is 0 outside
**Answer: b) Compile error — `x` declared in the if-statement's initializer is only in scope within the if/else block**

**44** In a health-check retry loop, what's idiomatic for "retry up to N times with early exit on success"?
a) Recursive function only  b) `for i := 0; i < maxRetries; i++ { if success { break } }`  c) `while` loop  d) `goto` only
**Answer: b) A bounded for-loop with a `break` on success is the idiomatic pattern**

**45** What does `goto` do in Go, and when is it discouraged?
a) Jumps to a labeled statement; discouraged for readability except rare cases like breaking nested loops/cleanup  b) Not supported in Go  c) Only works in switch statements  d) Same as break
**Answer: a) Go supports `goto` but idiomatic Go rarely uses it — labeled break/continue or restructuring is usually preferred**

---

### Level 2 — Coding Problems

**46** Which of the following functions correctly classifies an HTTP status code using a `switch` statement with ranges?
a) `switch { case code >= 200 && code < 300: return "Success" }`
b) `switch code { case 200..299: return "Success" }`
c) `switch code >= 200 { case true: return "Success" }`
d) `switch (code) { case 200-299: return "Success" }`
**Answer: a) `switch { case code >= 200 && code < 300: return "Success" }`**
```go
package main
import "fmt"
func classify(code int) string {
	switch {
	case code >= 200 && code < 300:
		return "Success"
	case code >= 400 && code < 500:
		return "Client Error"
	case code >= 500 && code < 600:
		return "Server Error"
	default:
		return "Unknown"
	}
}
func main() {
	fmt.Println(classify(404))
}
```

**47** How do you correctly structure a retry loop in Go that attempts a call up to 3 times, breaking early on success?
a) `for attempt := 1; attempt <= 3; attempt++ { if success() { break } }`
b) `for (int attempt = 1; attempt <= 3; attempt++) { if (success()) break; }`
c) `loop attempt := 1 to 3 { if success { stop } }`
d) `while attempt <= 3 { attempt++; if success() { exit } }`
**Answer: a) `for attempt := 1; attempt <= 3; attempt++ { if success() { break } }`**
```go
package main
import "fmt"
func callAPI(attempt int) bool { return attempt == 2 }
func main() {
	for attempt := 1; attempt <= 3; attempt++ {
		if callAPI(attempt) {
			break
		}
	}
}
```

**48** How do you use a labeled break to completely exit a nested 2D loop when a target is found?
a) `break 2`
b) `exit Search`
c) `break Search` (where `Search:` is the label for the outer loop)
d) `return outer`
**Answer: c) `break Search` (where `Search:` is the label for the outer loop)**
```go
package main
import "fmt"
func main() {
	grid := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	target := 5
Search:
	for row := range grid {
		for col := range grid[row] {
			if grid[row][col] == target {
				break Search
			}
		}
	}
}
```

**49** How does the `fallthrough` keyword behave in a `switch` statement in Go?
a) It causes the `switch` statement to exit immediately.
b) It automatically transfers control to the very next `case` block, executing its code regardless of the next case's condition.
c) It acts as a `default` case if no other conditions match.
d) It re-evaluates the condition of the next `case` block before executing it.
**Answer: b) It automatically transfers control to the very next `case` block, executing its code regardless of the next case's condition.**
```go
package main
import "fmt"
func permissions(level int) []string {
	var perms []string
	switch {
	case level >= 3:
		perms = append(perms, "admin")
		fallthrough
	case level >= 2:
		perms = append(perms, "write")
		fallthrough
	case level >= 1:
		perms = append(perms, "read")
	}
	return perms
}
func main() {
	fmt.Println(permissions(3)) // Output: [admin write read]
}
```

---

# LEVEL 3: Functions

**50** How many values can a Go function return?
a) Only 1  b) Up to 2  c) Any number  d) Up to 4
**Answer: c) Any number of return values**

**51** What is a "named return"?
a) A function name that's a keyword  b) Return values declared in the function signature that act as local variables  c) A return value with a comment  d) Not valid Go
**Answer: b) e.g. `func div(a, b int) (result int, err error)` — `result` and `err` are pre-declared**

**52** When does `defer` execute?
a) Immediately  b) Just before the enclosing function returns, LIFO order  c) At program exit only  d) Before the function starts
**Answer: b) Deferred calls run in LIFO order right before the surrounding function returns**

**53** Predict output:
```go
func main() {
	defer fmt.Println("1")
	defer fmt.Println("2")
	defer fmt.Println("3")
}
```
a) 1 2 3  b) 3 2 1  c) Compile error  d) Only 3 prints
**Answer: b) 3 2 1 — LIFO (stack) order**

**54** What is a variadic function parameter?
a) A parameter with default value  b) `...T` — accepts zero or more arguments of type T as a slice  c) A pointer parameter  d) Not supported in Go
**Answer: b) `func sum(nums ...int)` accepts any number of int arguments, accessible as a slice inside**

**55** What is a closure in Go?
a) A function that closes files  b) A function value that references variables from outside its own body  c) Same as a method  d) A deprecated feature
**Answer: b) An anonymous function capturing variables from its enclosing scope**

**56** Predict output:
```go
func counter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}
func main() {
	c := counter()
	fmt.Println(c())
	fmt.Println(c())
}
```
a) 1 1  b) 1 2  c) 0 1  d) Compile error
**Answer: b) 1 2 — the closure retains its own `count` across calls**

**57** Can you pass a slice `...int` variadic call using an existing slice?
a) No, must list elements individually  b) Yes, using `mySlice...` spread syntax  c) Yes, using `*mySlice`  d) Only with arrays
**Answer: b) `sum(mySlice...)` spreads the slice into variadic args**

**58** What happens if a deferred function's arguments reference a variable that changes later?
a) Deferred call sees the latest value  b) Arguments are evaluated at defer-time, not execution-time  c) Compile error  d) Runtime panic
**Answer: b) Arguments to deferred calls are evaluated immediately when `defer` executes, not when the call actually runs**

**59** Predict output:
```go
func main() {
	x := 1
	defer fmt.Println(x)
	x = 2
	fmt.Println(x)
}
```
a) 2 2  b) 1 1  c) 2 1  d) 1 2
**Answer: c) 2 1 — prints "2" immediately, then deferred `fmt.Println(x)` uses x's value (1) captured at defer time**

**60** What's a common use of `defer` in real backend code?
a) Looping  b) Closing files/DB connections/unlocking mutexes reliably  c) Declaring constants  d) String formatting
**Answer: b) Resource cleanup — e.g. `defer file.Close()`, `defer mu.Unlock()`**

**61** Are Go functions first-class values?
a) No  b) Yes — can be assigned to variables, passed as arguments, returned from functions  c) Only named functions  d) Only methods
**Answer: b) Yes, functions are first-class citizens in Go**

**62** What does this function signature mean? `func process(data []byte) (result string, err error)`
a) Two required args  b) Takes bytes, returns a string and an error (named returns)  c) Invalid syntax  d) Takes no arguments
**Answer: b) Takes a byte slice, returns a named string result and named error**

**63** Can a Go function be recursive?
a) No  b) Yes, functions can call themselves  c) Only methods can  d) Only with special syntax
**Answer: b) Yes, standard recursion is supported**

**64** What's the zero value returned by named returns if you just call `return` with no values?
a) Compile error  b) The current values of the named return variables  c) Always nil  d) Always zero regardless of prior assignment
**Answer: b) Whatever the named return variables currently hold**

**65** Predict output:
```go
func mightPanic() (result int) {
	defer func() {
		if r := recover(); r != nil {
			result = -1
		}
	}()
	panic("boom")
}
func main() {
	fmt.Println(mightPanic())
}
```
a) Program crashes  b) -1  c) 0  d) "boom"
**Answer: b) -1 — deferred recover() catches the panic and sets the named return before returning**

**66** What is the difference between a function and a method in Go?
a) No difference  b) A method has a receiver argument, associating it with a type  c) Methods can't return values  d) Functions can't take structs
**Answer: b) Methods are functions with a receiver: `func (r ReceiverType) Name(...)`**

**67** Can you have a variadic parameter combined with regular parameters?
a) No  b) Yes, but variadic must be last: `func f(a int, b ...string)`  c) Yes, variadic must be first  d) Only in generics
**Answer: b) Variadic parameter must always be the last parameter**

**68** What does calling a variadic function with zero arguments do?
a) Compile error  b) The variadic parameter is an empty (nil) slice  c) Runtime panic  d) Uses default values
**Answer: b) It becomes a nil/empty slice, safely iterable with zero length**

**69** Can deferred functions modify named return values even after a panic is recovered?
a) True, this is exactly how idiomatic panic-recovery-with-cleanup patterns work in Go.
b) False, named returns are frozen the moment a panic occurs.
c) True, but only if the function uses pointer return types.
d) False, deferred functions execute, but cannot mutate the outer function's scope.
**Answer: a) True, this is exactly how idiomatic panic-recovery-with-cleanup patterns work in Go.**

---

### Level 3 — Coding Problems

**70** Which of the following is the correct syntax to define and call a variadic function that takes any number of integers?
a) `func sum(nums ...int) { ... }` called with `sum([]int{1, 2}...)`
b) `func sum(nums []int...) { ... }` called with `sum(1, 2, 3)`
c) `func sum(nums ...[]int) { ... }` called with `sum(1, 2)`
d) `func sum(...nums int) { ... }` called with `sum(1, 2)`
**Answer: a) `func sum(nums ...int) { ... }` called with `sum([]int{1, 2}...)`**
```go
package main
import "fmt"
func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}
func main() {
	fmt.Println(sum(1, 2, 3))
	nums := []int{10, 20, 30}
	fmt.Println(sum(nums...))
}
```

**71** Which implementation correctly uses named returns to return an error instead of panicking on division by zero?
a) `func safeDivide(a, b float64) (result float64, err error) { if b == 0 { err = errors.New("zero"); return }; result = a / b; return }`
b) `func safeDivide(a, b float64) (float64, error) { if b == 0 { return 0, errors.New("zero") }; return a / b, nil }`
c) `func safeDivide(a, b float64) (res float64) { if b == 0 { panic("zero") } return a / b }`
d) `func safeDivide(a, b float64) (result float64, err error) { if b == 0 { return nil, "error" }; result = a / b; return }`
**Answer: a) `func safeDivide(a, b float64) (result float64, err error) { if b == 0 { err = errors.New("zero"); return }; result = a / b; return }`**
```go
package main
import (
	"errors"
	"fmt"
)
func safeDivide(a, b float64) (result float64, err error) {
	if b == 0 {
		err = errors.New("division by zero")
		return
	}
	result = a / b
	return
}
func main() {
	r, err := safeDivide(10, 0)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", r)
	}
}
```

**72** How do you create a closure-based rate limiter in Go that tracks state (like a call count) internally?
a) By returning a function that references a variable declared in the outer function's scope.
b) By using global variables that the returned function accesses.
c) By passing a pointer to a struct every time the function is called.
d) By using the `static` keyword on a variable inside the function.
**Answer: a) By returning a function that references a variable declared in the outer function's scope.**
```go
package main
import "fmt"
func makeLimiter(max int) func() bool {
	count := 0
	return func() bool {
		if count >= max {
			return false
		}
		count++
		return true
	}
}
func main() {
	allow := makeLimiter(2)
	fmt.Println(allow()) // true
	fmt.Println(allow()) // true
	fmt.Println(allow()) // false
}
```

**73** How can you use `defer` to accurately log the entry and exit timing (duration) of a function?
a) `start := time.Now(); defer func() { fmt.Println(time.Since(start)) }()`
b) `defer fmt.Println(time.Since(time.Now()))`
c) `defer time.Since(time.Now())`
d) `start := time.Now(); defer fmt.Println(time.Since(start))`
**Answer: a) `start := time.Now(); defer func() { fmt.Println(time.Since(start)) }()`**
```go
package main
import (
	"fmt"
	"time"
)
func queryDB() {
	start := time.Now()
	defer func() {
		fmt.Printf("queryDB took %v\n", time.Since(start))
	}()
	time.Sleep(50 * time.Millisecond)
}
func main() {
	queryDB()
}
```

**74** How do you correctly use `recover()` inside a deferred function to prevent a panic from crashing the program?
a) `defer func() { if r := recover(); r != nil { fmt.Println(r) } }()`
b) `defer recover()`
c) `if err := recover(); err != nil { defer fmt.Println(err) }`
d) `defer func() { recover(func(err error) { fmt.Println(err) }) }()`
**Answer: a) `defer func() { if r := recover(); r != nil { fmt.Println(r) } }()`**
```go
package main
import "fmt"
func safeWorker(job func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered from panic:", r)
		}
	}()
	job()
}
func main() {
	safeWorker(func() {
		panic("job failed unexpectedly")
	})
	fmt.Println("main continues running")
}
```

---

# LEVEL 4: Arrays, Slices, and Maps

**75** What's the key difference between an array and a slice in Go?
a) No difference  b) Arrays have fixed size at compile time; slices are dynamic, backed by an array  c) Slices are always faster  d) Arrays can't hold structs
**Answer: b) Arrays are fixed-length value types; slices are flexible, reference-like views over an underlying array**

**76** What does `len()` vs `cap()` return for a slice?
a) Same thing  b) `len` = number of elements currently in use; `cap` = size of underlying array from the slice's start  c) `cap` is always double `len`  d) `len` is for arrays only
**Answer: b) len = current element count, cap = max elements before reallocation is needed**

**77** What happens when you append beyond a slice's capacity?
a) Panic  b) Go allocates a new, larger underlying array and copies data  c) Silently truncates  d) Compile error
**Answer: b) A new array is allocated (typically growth factor ~2x for small slices) and the old data is copied over**

**78** Predict output:
```go
a := []int{1, 2, 3}
b := a
b[0] = 100
fmt.Println(a[0])
```
a) 1  b) 100  c) Compile error  d) 0
**Answer: b) 100 — slices share the same underlying array; `b := a` copies the slice header, not the data**

**79** What's the zero value of a slice?
a) Empty slice `[]int{}`  b) `nil`  c) Compile error, must initialize  d) Array of zeros
**Answer: b) `nil` — a nil slice has len=0, cap=0, and is usable (e.g., append works on it)**

**80** How do you create a slice with initial length 5 and capacity 10?
a) `make([]int, 5, 10)`  b) `make([]int, 10, 5)`  c) `[]int{5, 10}`  d) `new([]int, 5, 10)`
**Answer: a) `make([]int, len, cap)`**

**81** What does map lookup return for a missing key?
a) Panic  b) Zero value of the value type, and `false` if using the "comma ok" form  c) nil always  d) Compile error
**Answer: b) e.g. `v, ok := myMap["missing"]` — v is zero value, ok is false**

**82** What happens if you write to a nil map?
a) Silently no-ops  b) Panic: "assignment to entry in nil map"  c) Auto-initializes  d) Compile error
**Answer: b) Panics at runtime — nil maps must be initialized with `make` or a literal before writing**

**83** Can you read from a nil map?
a) No, panics  b) Yes, returns zero value  c) Compile error  d) Only with comma-ok
**Answer: b) Reading from a nil map is safe and returns the zero value (unlike writing)**

**84** What does `a[1:3]` do for `a := []int{10, 20, 30, 40}`?
a) `[20, 30]`  b) `[10, 20, 30]`  c) `[20, 30, 40]`  d) `[30, 40]`
**Answer: a) `[20, 30]` — slicing is [low:high), high exclusive**

**85** What is the danger of slicing a large array/slice and keeping only a small sub-slice long-term?
a) None  b) The sub-slice keeps the entire original backing array alive, preventing GC of the rest  c) It copies data unnecessarily  d) Not possible in Go
**Answer: b) Memory leak risk — small slices can pin large backing arrays in memory unless explicitly copied**

**86** How do you safely copy a slice to avoid backing-array sharing issues?
a) `b := a`  b) `b := make([]int, len(a)); copy(b, a)`  c) `b := a[:]`  d) Not possible
**Answer: b) `make` + `copy` creates an independent backing array**

**87** What does `delete(myMap, "key")` do if "key" doesn't exist?
a) Panics  b) No-op, safe  c) Returns an error  d) Compile error
**Answer: b) Safe no-op — deleting a non-existent key does nothing**

**88** Are Go maps safe for concurrent read/write from multiple goroutines?
a) Yes, always  b) No — concurrent map writes (or write+read) cause a runtime panic/race  c) Only reads are safe, writes need locking too but reads are always fine  d) Yes, if buffered
**Answer: b) No — maps are not safe for concurrent use without external synchronization (sync.Mutex or sync.Map)**

**89** Predict output:
```go
arr := [3]int{1, 2, 3}
modify := func(a [3]int) {
	a[0] = 999
}
modify(arr)
fmt.Println(arr[0])
```
a) 999  b) 1  c) Compile error  d) 0
**Answer: b) 1 — arrays are value types; passing to a function copies the entire array**

**90** What's the idiomatic way to check if a key exists in a map without caring about its value?
a) `if myMap["key"] != nil`  b) `if _, ok := myMap["key"]; ok`  c) `if myMap.has("key")`  d) `if len(myMap["key"]) > 0`
**Answer: b) The "comma ok" idiom**

**91** What does `append(a, b...)` do when both `a` and `b` are `[]int`?
a) Compile error  b) Appends all elements of b onto a  c) Appends b as a single nested element  d) Only works with `copy`
**Answer: b) Spreads b's elements and appends them individually**

**92** What is the output?
```go
m := map[string]int{"a": 1, "b": 2}
m["c"] = 3
fmt.Println(len(m))
```
a) 2  b) 3  c) Compile error  d) 0
**Answer: b) 3 — map now has three key-value pairs**

**93** What's a struct-keyed map used for in real systems?
a) Not allowed in Go  b) Using composite keys (e.g., struct{UserID, ResourceID}) for caching/dedup logic  c) Only string keys are allowed  d) Structs can't be map keys unless they implement an interface
**Answer: b) Structs with only comparable fields can be map keys — common for composite-key caches**

**94** What does `make([]int, 0, 100)` optimize for in code that appends heavily?
a) Nothing  b) Pre-allocates capacity to avoid repeated reallocation during appends  c) Creates a fixed array  d) Wastes memory always
**Answer: b) Pre-allocating capacity avoids multiple reallocations/copies as elements are appended — common perf optimization**

---

### Level 4 — Coding Problems

**95** Which of the following is the standard, idiomatic way to deduplicate a slice of strings in Go using a map as a set?
a) By iterating over the slice and deleting duplicates from it directly using `delete()`.
b) By using a `map[string]bool` to track seen elements and appending unseen elements to a new slice.
c) By casting the slice to a `set` type: `set(items)`.
d) By using `slices.Unique(items)` which modifies the slice in-place without memory allocation.
**Answer: b) By using a `map[string]bool` to track seen elements and appending unseen elements to a new slice.**
```go
package main
import "fmt"
func dedupe(items []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}
func main() {
	fmt.Println(dedupe([]string{"a", "b", "a", "c", "b"}))
}
```

**96** How do you accurately count the frequency of each word in a slice of strings and return the counts?
a) `freq := make(map[string]int); for _, w := range words { freq[w]++ }`
b) `freq := map[string]int{}; for _, w := range words { freq[w] = freq[w] + 1 }`
c) Both a and b are correct and idiomatic in Go.
d) `freq := make(map[string]int); for w in words { freq[w]++ }`
**Answer: c) Both a and b are correct and idiomatic in Go.**
```go
package main
import "fmt"
func wordFreq(words []string) map[string]int {
	freq := make(map[string]int)
	for _, w := range words {
		freq[w]++
	}
	return freq
}
func main() {
	logs := []string{"ERROR", "INFO", "ERROR", "WARN", "ERROR"}
	fmt.Println(wordFreq(logs))
}
```

**97** How do you properly split a slice `items` into smaller chunks (batches) of size `size`?
a) By using a `for` loop incrementing by `size` and slicing `items[i:i+size]`, ensuring the upper bound doesn't exceed `len(items)`.
b) By using `slices.Chunk(items, size)`.
c) By copying elements into a 2D array iteratively using a `while` loop.
d) By dividing `len(items) / size` and dynamically sizing a matrix.
**Answer: a) By using a `for` loop incrementing by `size` and slicing `items[i:i+size]`, ensuring the upper bound doesn't exceed `len(items)`.**
```go
package main
import "fmt"
func batch(items []int, size int) [][]int {
	var batches [][]int
	for i := 0; i < len(items); i += size {
		end := i + size
		if end > len(items) {
			end = len(items)
		}
		batches = append(batches, items[i:end])
	}
	return batches
}
func main() {
	items := []int{1, 2, 3, 4, 5, 6, 7}
	fmt.Println(batch(items, 3))
}
```

**98** How do you safely merge two `map[string]int` config maps so that the `override` map overrides the `base` map values?
a) `merged := append(base, override...)`
b) Create a new map, loop over `base` to assign values, then loop over `override` to assign/overwrite values.
c) `merged := map.Merge(base, override)`
d) `for k, v := range override { base[k] = v }; merged := base` (This mutates the base map, which may be unsafe).
**Answer: b) Create a new map, loop over `base` to assign values, then loop over `override` to assign/overwrite values.**
```go
package main
import "fmt"
func mergeConfigs(base, override map[string]int) map[string]int {
	merged := make(map[string]int, len(base))
	for k, v := range base {
		merged[k] = v
	}
	for k, v := range override {
		merged[k] = v
	}
	return merged
}
func main() {
	base := map[string]int{"timeout": 30, "retries": 3}
	override := map[string]int{"timeout": 60}
	fmt.Println(mergeConfigs(base, override)) // {timeout: 60, retries: 3}
}
```

**99** What happens if you slice an array `a := original[:2]` and then `append(a, 999)` when the original array had a capacity of 5?
a) A panic occurs because `a` is full.
b) `a` receives a newly allocated backing array with the value `999` appended, leaving `original` untouched.
c) The value `999` overwrites the element at `original[2]` because `a` still shares the same backing array and has spare capacity.
d) The compiler throws an error about appending to a sub-slice.
**Answer: c) The value `999` overwrites the element at `original[2]` because `a` still shares the same backing array and has spare capacity.**
```go
package main
import "fmt"
func main() {
	original := make([]int, 3, 5)
	original[0], original[1], original[2] = 1, 2, 3

	a := original[:2] // len=2, cap=5 (shares backing array)
	a = append(a, 999) // overwrites original[2] since cap allows it in-place

	fmt.Println("original:", original) // [1 2 999] <- unexpectedly mutated
	fmt.Println("a:", a)               // [1 2 999]
}
```

---

# LEVEL 5: Structs and Methods

**100** How do you define a struct in Go?
a) `struct MyStruct { ... }`  b) `type MyStruct struct { Field1 Type1; Field2 Type2 }`  c) `class MyStruct { ... }`  d) `interface MyStruct { ... }`
**Answer: b)**

**101** What's the difference between a value receiver and a pointer receiver on a method?
a) No difference  b) Value receiver operates on a copy; pointer receiver can mutate the original  c) Pointer receivers are always faster  d) Value receivers can't call other methods
**Answer: b) `func (s Struct) Method()` copies; `func (s *Struct) Method()` operates on original via pointer**

**102** Predict output:
```go
type Counter struct{ count int }
func (c Counter) IncValue() { c.count++ }
func (c *Counter) IncPointer() { c.count++ }

func main() {
	c := Counter{}
	c.IncValue()
	c.IncPointer()
	fmt.Println(c.count)
}
```
a) 0  b) 1  c) 2  d) Compile error
**Answer: b) 1 — IncValue mutates a copy (no effect); IncPointer mutates the original**

**103** Can you embed one struct inside another in Go (composition)?
a) No, Go has no inheritance-like feature  b) Yes, via anonymous struct fields (embedding)  c) Only with interfaces  d) Only pointers can be embedded
**Answer: b) Yes — embedding promotes the embedded struct's fields/methods to the outer struct**

**104** What does struct embedding provide that's similar to inheritance?
a) True polymorphic override  b) Field and method promotion — outer struct can access embedded struct's members directly  c) Nothing, purely cosmetic  d) Multiple dispatch
**Answer: b) Promotion, not true inheritance — there's no dynamic override mechanism**

**105** What is a struct tag, e.g. `json:"name"`, used for?
a) Comments only  b) Metadata read via reflection, commonly for encoding/decoding (JSON, DB, etc.)  c) Compiler directives  d) Not valid Go syntax
**Answer: b) Struct tags provide metadata used by packages like encoding/json via reflection**

**106** What's the zero value of a struct?
a) nil  b) All fields set to their respective zero values  c) Compile error  d) Empty struct{}
**Answer: b) Every field takes its own zero value (0, "", false, nil, etc.)**

**107** Can structs be compared with `==`?
a) Never  b) Yes, if all fields are comparable types  c) Only pointer comparisons  d) Only with reflect.DeepEqual
**Answer: b) Structs are comparable if every field's type is comparable (no slices/maps/funcs as fields)**

**108** Predict output:
```go
type Point struct{ X, Y int }
p1 := Point{1, 2}
p2 := Point{1, 2}
fmt.Println(p1 == p2)
```
a) false  b) true  c) Compile error  d) panic
**Answer: b) true — struct equality compares all fields**

**109** What does `new(MyStruct)` return?
a) A value of type MyStruct  b) A pointer `*MyStruct` to a zero-valued struct  c) nil  d) Compile error
**Answer: b) `new()` allocates and returns a pointer to a zeroed value**

**110** How do you create a struct instance with named fields?
a) `Point{1, 2}` positional only  b) `Point{X: 1, Y: 2}`  c) Both a and b are valid  d) Neither is valid
**Answer: c) Both positional and named-field struct literals are valid Go**

**111** What happens when you pass a struct to a function by value (not pointer)?
a) Struct is passed by reference automatically  b) A full copy of the struct is made  c) Compile error for large structs  d) Only the first field is copied
**Answer: b) Go copies the entire struct when passed by value — can be costly for large structs**

**112** What does `func (s *Service) Start() error` typically indicate about design intent?
a) Nothing special  b) The method likely mutates the receiver's state (e.g., connection setup) so uses a pointer receiver  c) It's always faster than a value receiver  d) It must be called on a nil Service
**Answer: b) Pointer receivers are conventionally used when the method needs to mutate state or the struct is large**

**113** Can an embedded struct's method be "overridden" by defining a same-named method on the outer struct?
a) No, not possible  b) Yes — the outer struct's own method shadows the promoted one when called directly on the outer type  c) Compile error, ambiguous  d) Only for pointer receivers
**Answer: b) Yes, the outer type's own method takes precedence (shadowing), though the inner one can still be called explicitly via the field name**

**114** What's the convention for whether all methods on a type should use consistent receiver types (all value or all pointer)?
a) No convention, mix freely  b) Idiomatic Go generally keeps receiver type consistent across a type's methods to avoid confusion  c) Must always use pointer  d) Must always use value
**Answer: b) Consistency is the idiomatic Go convention, though not enforced by the compiler**

**115** What does `type ID int` (a defined type) let you do that a plain `int` doesn't?
a) Nothing different  b) Attach methods to it, gaining type safety distinct from raw int  c) Store it in maps only  d) Nothing, syntax error
**Answer: b) Named/defined types can have their own methods and provide type-safety (e.g., `UserID` vs `ProductID` won't mix accidentally)**

**116** Predict output:
```go
type Animal struct{ Name string }
type Dog struct {
	Animal
	Breed string
}
d := Dog{Animal{"Rex"}, "Labrador"}
fmt.Println(d.Name)
```
a) Compile error  b) "Rex"  c) "" (empty)  d) panic
**Answer: b) "Rex" — field promotion via embedding lets you access `d.Name` directly**

**117** What's the difference between `Dog{Animal{"Rex"}, "Labrador"}` and `Dog{Animal: Animal{"Rex"}, Breed: "Labrador"}`?
a) Different behavior  b) Same result, positional vs named-field initialization  c) First is invalid  d) Second is invalid
**Answer: b) Functionally identical — named fields are just clearer/safer against field-order changes**

**118** Can methods be defined on non-struct named types, e.g. `type Celsius float64`?
a) No, only structs  b) Yes, any named type can have methods  c) Only if it embeds a struct  d) Only pointer types
**Answer: b) Yes — Go allows methods on any user-defined (named) type, not just structs**

**119** What happens if you call a pointer-receiver method on a value that is not addressable (e.g., a map value)?
a) Works fine always  b) Compile error — map values aren't addressable, so pointer methods can't be called directly on them  c) Panic  d) Auto-converts
**Answer: b) Compile error — you can't take the address of a map element, so pointer-receiver methods fail there directly (would need to extract, modify, reassign)**

---

### Level 5 — Coding Problems

**120** How do you correctly define a `Server` struct with a pointer-receiver method `MarkUnhealthy()` that modifies its internal state?
a) `func (s Server) MarkUnhealthy() { s.Healthy = false }`
b) `func (s *Server) MarkUnhealthy() { s.Healthy = false }`
c) `func MarkUnhealthy(s *Server) { s.Healthy = false }` (This is a function, not a method)
d) `func (s *Server) MarkUnhealthy() { *s.Healthy = false }`
**Answer: b) `func (s *Server) MarkUnhealthy() { s.Healthy = false }`**
```go
package main
import "fmt"
type Server struct {
	Name    string
	Port    int
	Healthy bool
}
func (s *Server) MarkUnhealthy() {
	s.Healthy = false
}
func main() {
	s := Server{Name: "api-1", Port: 8080, Healthy: true}
	s.MarkUnhealthy()
	fmt.Printf("%+v\n", s)
}
```

**121** In Go, how do you use struct embedding to include a `BaseHandler` (which has a `Log` method) inside an `AuthHandler`, and call `Log` on an instance of `AuthHandler`?
a) `type AuthHandler struct { Base BaseHandler }; h.Base.Log("msg")`
b) `type AuthHandler struct { BaseHandler }; h.BaseHandler.Log("msg")`
c) `type AuthHandler struct { BaseHandler }; h.Log("msg")`
d) Go does not support calling embedded methods directly on the outer struct.
**Answer: c) `type AuthHandler struct { BaseHandler }; h.Log("msg")`**
```go
package main
import "fmt"
type BaseHandler struct{}
func (b BaseHandler) Log(msg string) {
	fmt.Println("[LOG]", msg)
}
type AuthHandler struct {
	BaseHandler
	Realm string
}
func main() {
	h := AuthHandler{Realm: "admin"}
	h.Log("authenticating request") // Promoted method!
}
```

**122** How do you define a `Config` struct with fields `APIKey` and `MaxRetries` that marshal to JSON as `api_key` and `max_retries`?
a) By naming the fields in lowercase: `api_key string` (This makes them unexported).
b) By using JSON tags: `APIKey string `+"`"+`json:"api_key"`+"`"
c) By implementing a custom `MarshalJSON` interface for every field.
d) By passing a formatting option to `json.Marshal(c, "snake_case")`.
**Answer: b) By using JSON tags: `APIKey string `+"`"+`json:"api_key"`+"`"**
```go
package main
import (
	"encoding/json"
	"fmt"
)
type Config struct {
	APIKey     string `json:"api_key"`
	MaxRetries int    `json:"max_retries"`
}
func main() {
	c := Config{APIKey: "secret123", MaxRetries: 5}
	data, _ := json.Marshal(c)
	fmt.Println(string(data))
}
```

**123** How do you write a method `Status()` on a `HealthCheck` struct that checks if `errCount` has exceeded a `maxErrThreshold` constant?
a) `func (h HealthCheck) Status() string { if h.errCount >= maxErrThreshold { return "unhealthy" }; return "healthy" }`
b) `func Status(h HealthCheck) string { if h.errCount >= maxErrThreshold { return "unhealthy" }; return "healthy" }`
c) `func (h *HealthCheck) string Status() { if h.errCount >= maxErrThreshold { return "unhealthy" }; return "healthy" }`
d) Constants cannot be used inside struct methods.
**Answer: a) `func (h HealthCheck) Status() string { if h.errCount >= maxErrThreshold { return "unhealthy" }; return "healthy" }`**
```go
package main
import "fmt"
const maxErrThreshold = 3
type HealthCheck struct {
	errCount int
}
func (h HealthCheck) Status() string {
	if h.errCount >= maxErrThreshold {
		return "unhealthy"
	}
	return "healthy"
}
func main() {
	h := HealthCheck{errCount: 4}
	fmt.Println(h.Status())
}
```

**124** Can you define methods on non-struct named types in Go, like creating a `Minutes()` method on a `type Duration int`?
a) No, methods can only be attached to structs in Go.
b) Yes, methods can be defined on any user-defined named type (like `type Duration int`) in the same package.
c) Yes, but only if the named type is an interface.
d) No, primitive types like `int` can never have methods, even if aliased.
**Answer: b) Yes, methods can be defined on any user-defined named type (like `type Duration int`) in the same package.**
```go
package main
import "fmt"
type Duration int // seconds
func (d Duration) Minutes() float64 {
	return float64(d) / 60.0
}
func main() {
	d := Duration(150)
	fmt.Printf("%.2f minutes\n", d.Minutes())
}
```

---

# LEVEL 6: Interfaces

**125** How are interfaces implemented in Go?
a) Explicit `implements` keyword like Java  b) Implicitly — any type with matching methods satisfies the interface automatically  c) Interfaces must be declared inside the struct  d) Not supported in Go
**Answer: b) Implicit satisfaction — no explicit declaration needed**

**126** What is the empty interface `interface{}` (or `any` in modern Go)?
a) An interface with no methods, satisfied by every type  b) Invalid syntax  c) Only satisfied by structs  d) Equivalent to `nil`
**Answer: a) Zero-method interface — every type satisfies it, useful for generic-ish containers pre-generics**

**127** What is a type assertion, e.g. `v, ok := i.(string)`?
a) Type conversion  b) Extracts the concrete value from an interface if it holds that type, with `ok` indicating success  c) Compile-time only check  d) Not valid syntax
**Answer: b) Runtime check + extraction; the "comma ok" form avoids panicking on mismatch**

**128** What happens with `v := i.(string)` (no comma-ok) if `i` doesn't hold a string?
a) Returns zero value  b) Panics at runtime  c) Compile error  d) Returns nil
**Answer: b) Panics — single-value type assertion panics on failure**

**129** What is a type switch used for?
a) Same as regular switch  b) Branching logic based on the dynamic/concrete type stored in an interface  c) Not valid Go syntax  d) Only works with structs
**Answer: b) `switch v := i.(type) { case int: ... case string: ... }`**

**130** Can a nil concrete type stored in an interface make `interface == nil` false?
a) No, always true if nil  b) Yes — a non-nil interface can wrap a nil pointer, making the interface itself non-nil (classic Go gotcha)  c) Compile error  d) Always panics
**Answer: b) Yes — this is one of Go's most famous gotchas: an interface holding a typed nil is NOT equal to nil**

**131** What must a type do to satisfy `io.Writer`?
a) Nothing special  b) Implement `Write(p []byte) (n int, err error)`  c) Embed io.Writer  d) Be a struct only
**Answer: b) Implement exactly that method signature**

**132** What's the idiomatic Go interface design philosophy?
a) Large interfaces with many methods, defined upfront  b) Small, focused interfaces (often 1-2 methods), often defined at point of use by the consumer  c) One giant interface per package  d) Interfaces should mirror structs 1:1
**Answer: b) "Accept interfaces, return structs" — small interfaces like io.Reader, io.Writer are the Go idiom**

**133** Where should interfaces ideally be defined in Go — producer or consumer package?
a) Always in the producer/implementation package  b) Often in the consumer package that needs the behavior — Go encourages defining interfaces where they're used  c) Doesn't matter  d) Only in a shared "interfaces" package
**Answer: b) Consumer-side interface definition is idiomatic — avoids unnecessary coupling**

**134** Predict output:
```go
type Shape interface{ Area() float64 }
type Circle struct{ R float64 }
func (c Circle) Area() float64 { return 3.14 * c.R * c.R }

func main() {
	var s Shape = Circle{R: 2}
	fmt.Println(s.Area())
}
```
a) Compile error  b) 12.56  c) 0  d) panic
**Answer: b) 12.56 — Circle implicitly satisfies Shape**

**135** What does the `error` interface look like in Go's standard library?
a) `type error interface { Error() string }`  b) A struct  c) A built-in type with no methods  d) `type error interface { Message() string }`
**Answer: a) The error interface requires exactly one method: `Error() string`**

**136** Can you assign a `nil` value of a concrete pointer type to an interface variable and get a "non-nil interface"?
a) No  b) Yes, this is the classic "nil interface vs nil concrete type" trap  c) Compile error  d) Only for structs
**Answer: b) Yes — `var p *MyType = nil; var i MyInterface = p` makes `i != nil` even though `p == nil`**

**137** What is a "duck typing" language characteristic, and does Go have it (loosely)?
a) N/A  b) "If it walks like a duck and implements the methods, it satisfies the interface" — Go's implicit interfaces resemble this at compile time  c) Go has no such concept  d) Only applies to reflection
**Answer: b) Go's implicit, structural interface satisfaction is often compared to duck typing, but it's still statically type-checked**

**138** How many methods can an interface have?
a) Exactly 1  b) Any number, including zero  c) Max 5  d) Must be at least 2
**Answer: b) Zero or more — zero-method interfaces (empty interface) are valid and useful**

**139** What does interface embedding look like, e.g. `io.ReadWriter`?
a) Not supported  b) `type ReadWriter interface { Reader; Writer }` — composes multiple interfaces  c) Only structs can embed  d) Requires explicit method redeclaration
**Answer: b) Interfaces can embed other interfaces to compose larger interfaces from small ones**

**140** What's the output?
```go
var i interface{} = 42
switch v := i.(type) {
case string:
	fmt.Println("string:", v)
case int:
	fmt.Println("int:", v)
default:
	fmt.Println("other")
}
```
a) string: 42  b) int: 42  c) other  d) Compile error
**Answer: b) int: 42 — type switch matches the concrete int type**

**141** Why is over-mocking with large interfaces considered an anti-pattern in Go?
a) It's not an anti-pattern  b) Large interfaces are harder to implement/mock and violate the "small interfaces" idiom, increasing coupling  c) Go doesn't support mocking  d) Interfaces can't be used in tests
**Answer: b) Smaller, focused interfaces are easier to implement, test, and mock — Go favors minimal interfaces**

**142** What does `var _ Shape = Circle{}` (blank identifier assignment) accomplish?
a) Nothing, invalid  b) A compile-time check that Circle satisfies the Shape interface, without creating a usable variable  c) Runtime assertion  d) Declares Circle as abstract
**Answer: b) A common idiom to statically verify interface satisfaction at compile time**

**143** Can a struct satisfy multiple interfaces simultaneously?
a) No, only one interface per type  b) Yes — implicit satisfaction means a type can satisfy any number of interfaces it happens to match  c) Only with embedding  d) Requires explicit declaration per interface
**Answer: b) Yes — as long as the method set matches, a type can satisfy many interfaces at once**

**144** What's the practical difference between `interface{}` / `any` and Go generics (1.18+) for writing reusable functions?
a) No difference  b) Empty interface loses type information (needs assertions) at runtime; generics preserve compile-time type safety  c) Generics are slower always  d) Generics replace interfaces entirely
**Answer: b) Generics give compile-time type checking and avoid the runtime type-assertion overhead/unsafety of `interface{}`**

---

### Level 6 — Coding Problems

**145** How do you correctly implement an interface in Go, such as a `Notifier` interface with `Send(msg string) error` on an `EmailNotifier` struct?
a) By explicitly declaring `type EmailNotifier implements Notifier struct { ... }`.
b) By implicitly providing a matching method `Send(msg string) error` with `EmailNotifier` as the receiver.
c) By embedding `Notifier` inside `EmailNotifier` and overriding the method.
d) By passing a function pointer to a `Notifier` constructor.
**Answer: b) By implicitly providing a matching method `Send(msg string) error` with `EmailNotifier` as the receiver.**
```go
package main
import "fmt"
type Notifier interface {
	Send(msg string) error
}
type EmailNotifier struct{ Address string }
func (e EmailNotifier) Send(msg string) error {
	fmt.Printf("Emailing %s: %s\n", e.Address, msg)
	return nil
}
func main() {
	var n Notifier = EmailNotifier{Address: "ops@company.com"}
	n.Send("Server down!")
}
```

**146** What is the correct syntax for a type switch in Go that determines the underlying type of an `interface{}` value `i`?
a) `switch i.(type) { case int: ... }`
b) `switch typeof(i) { case int: ... }`
c) `switch v := i.(type) { case int: ... }`
d) Both A and C are correct.
**Answer: d) Both A and C are correct.**
```go
package main
import "fmt"
func describe(i interface{}) {
	switch v := i.(type) {
	case int:
		fmt.Println("integer:", v)
	case string:
		fmt.Println("string:", v)
	case bool:
		fmt.Println("boolean:", v)
	default:
		fmt.Printf("unknown type: %T\n", v)
	}
}
func main() {
	describe(42)
	describe("hello")
	describe(true)
}
```

**147** How do you safely assert that a returned `error` is a specific custom error type (e.g., `*NotFoundError`) in modern Go?
a) By using a type assertion: `err.(*NotFoundError)`
b) By using `errors.Is(err, NotFoundError)`
c) By using `errors.As(err, &nfErr)` where `nfErr` is a pointer to `*NotFoundError`.
d) By checking `if err.Error() == "*NotFoundError"`
**Answer: c) By using `errors.As(err, &nfErr)` where `nfErr` is a pointer to `*NotFoundError`.**
```go
package main
import (
	"errors"
	"fmt"
)
type NotFoundError struct { Resource string }
func (e *NotFoundError) Error() string { return fmt.Sprintf("%s not found", e.Resource) }

func findUser(id int) error { return &NotFoundError{Resource: "user"} }

func main() {
	err := findUser(99)
	var nfErr *NotFoundError
	if errors.As(err, &nfErr) {
		fmt.Println("Not found error:", nfErr.Resource)
	}
}
```

**148** What happens when a function returns a typed nil pointer (like `*MyError(nil)`) that is assigned to an `error` interface variable?
a) The `error` interface variable becomes `nil`.
b) The `error` interface variable is non-nil because the interface value contains type information (`*MyError`), even though the value itself is nil.
c) The compiler prevents this assignment.
d) It panics at runtime when checked against `nil`.
**Answer: b) The `error` interface variable is non-nil because the interface value contains type information (`*MyError`), even though the value itself is nil.**
```go
package main
import "fmt"
type MyError struct{ msg string }
func (e *MyError) Error() string { return e.msg }

func doWork(fail bool) *MyError {
	return nil // typed nil *MyError
}
func main() {
	var err error = doWork(false) // wraps a nil *MyError into a non-nil error interface
	if err != nil {
		fmt.Println("err is non-nil, even though the underlying pointer is nil!")
	}
}
```

**149** When implementing a `Cache` interface with an in-memory map on a struct, why do the implementation methods typically use a pointer receiver (`*MemCache`)?
a) Because interfaces can only hold pointers.
b) Because map types in Go require pointer receivers to work.
c) To avoid copying the struct on every method call, and to ensure any internal state changes affect the original struct.
d) It's just a style preference; value receivers work exactly the same way for state mutation.
**Answer: c) To avoid copying the struct on every method call, and to ensure any internal state changes affect the original struct.**
```go
package main
import "fmt"
type Cache interface {
	Get(key string) (string, bool)
	Set(key, value string)
}
type MemCache struct {
	data map[string]string
}
func NewMemCache() *MemCache {
	return &MemCache{data: make(map[string]string)}
}
func (m *MemCache) Get(key string) (string, bool) {
	v, ok := m.data[key]
	return v, ok
}
func (m *MemCache) Set(key, value string) {
	m.data[key] = value
}
func main() {
	var c Cache = NewMemCache()
	c.Set("user:1", "alice")
	v, ok := c.Get("user:1")
	fmt.Println(v, ok)
}
```

---

# Go Programming Curriculum — INTERMEDIATE TIER
### Covers Level 7–8: Error Handling and Pointers

---

# LEVEL 7: Error Handling

**150** What is the `error` type in Go fundamentally?
a) A concrete struct  b) A built-in interface with a single `Error() string` method  c) An alias for string  d) A special keyword
**Answer: b) An interface with `Error() string`**

**151** What is the idiomatic way to create a basic error in Go?
a) `new Error("msg")`  b) `errors.New("msg")`  c) `panic("msg")`  d) `err("msg")`
**Answer: b) `errors.New("msg")` or `fmt.Errorf("...")`**

**152** How do you check if an error `err` is specifically an `io.EOF` error?
a) `if err == io.EOF` or `if errors.Is(err, io.EOF)`  b) `if err.type == io.EOF`  c) `if err.contains("EOF")`  d) You cannot
**Answer: a) `if errors.Is(err, io.EOF)` is the modern, robust way**

**153** What does `fmt.Errorf("failed to open: %w", err)` do?
a) Formats a string only  b) Wraps the original `err` so it can be extracted later via `errors.Unwrap` or `errors.Is`  c) Panics if err is nil  d) Compile error
**Answer: b) The `%w` verb wraps the error, retaining the original error for inspection**

**154** When should you use `panic` in a standard Go application?
a) For all errors  b) For database connection failures  c) For truly unrecoverable programming errors (e.g., nil pointer dereference, out of bounds)  d) Instead of returning `error`
**Answer: c) Only for unrecoverable errors/bugs, not for expected runtime failures**

**155** What function is used to stop a panic and regain control of the program?
a) `catch()`  b) `recover()`  c) `rescue()`  d) `stopPanic()`
**Answer: b) `recover()`, which must be called inside a deferred function**

**156** Predict the output:
```go
err1 := errors.New("error")
err2 := errors.New("error")
fmt.Println(err1 == err2)
```
a) true  b) false  c) Compile error  d) Panic
**Answer: b) false — `errors.New` returns a pointer to a struct, so each call creates a distinct pointer**

**157** What is the result of `recover()` if the program is NOT panicking?
a) Panics  b) Returns `nil`  c) Compile error  d) Blocks forever
**Answer: b) Returns `nil`**

**158** Which package provides `Is` and `As` functions for error handling?
a) `fmt`  b) `log`  c) `errors`  d) `runtime`
**Answer: c) The `errors` package (introduced in Go 1.13)**

**159** If a function returns `(int, error)`, what is the recommended way to name the error return variable if it's named?
a) `e`  b) `error`  c) `err`  d) `Exception`
**Answer: c) `err` is the standard idiomatic name for an error variable**

---

### Level 7 — Coding Problems

**160** How do you correctly implement a custom error type `HTTPError` containing a `Code` and `Message` that satisfies the `error` interface?
a) `type HTTPError struct{ Code int; Message string }; func (e HTTPError) Error() string { return e.Message }`
b) `type HTTPError interface { ... }`
c) `type HTTPError struct{ Code int; Message string }; func (e *HTTPError) Error() string { return fmt.Sprintf("HTTP %d: %s", e.Code, e.Message) }`
d) `func HTTPError(Code int, Message string) error { return fmt.Errorf("%d %s", Code, Message) }`
**Answer: c) `type HTTPError struct{ Code int; Message string }; func (e *HTTPError) Error() string { return fmt.Sprintf("HTTP %d: %s", e.Code, e.Message) }`**
```go
package main
import "fmt"
type HTTPError struct {
	Code    int
	Message string
}
func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.Code, e.Message)
}
func doRequest() error {
	return &HTTPError{Code: 404, Message: "Not Found"}
}
func main() {
	err := doRequest()
	if err != nil {
		fmt.Println(err)
	}
}
```

**161** How do you properly wrap an existing error `ErrInvalidFormat` with additional context so that it can be later checked with `errors.Is`?
a) `return fmt.Errorf("parsing failed: %w", ErrInvalidFormat)`
b) `return fmt.Errorf("parsing failed: %v", ErrInvalidFormat)`
c) `return errors.Wrap(ErrInvalidFormat, "parsing failed")`
d) `return ErrInvalidFormat + "parsing failed"`
**Answer: a) `return fmt.Errorf("parsing failed: %w", ErrInvalidFormat)`**
```go
package main
import (
	"errors"
	"fmt"
)
var ErrInvalidFormat = errors.New("invalid format")
func parseConfig(data string) error {
	if data == "" {
		return fmt.Errorf("parsing config failed: %w", ErrInvalidFormat)
	}
	return nil
}
func main() {
	err := parseConfig("")
	if errors.Is(err, ErrInvalidFormat) {
		fmt.Println("Matched wrapped error!")
	}
}
```

---

# LEVEL 8: Pointers

**162** What does the `&` operator do in Go?
a) Logical AND  b) Bitwise AND  c) Takes the memory address of a variable  d) Both b and c
**Answer: d) Bitwise AND, and Address-of operator**

**163** What does the `*` operator do when placed before a pointer variable (e.g., `*ptr`)?
a) Multiplies it by zero  b) Dereferences the pointer to access or mutate the underlying value  c) Declares a pointer type  d) Nothing
**Answer: b) Dereferences the pointer**

**164** Does Go support pointer arithmetic (e.g., `ptr++`) like C/C++?
a) Yes  b) No, not by default (only via the unsafe package)  c) Only on arrays  d) Yes, but only for ints
**Answer: b) No, pointer arithmetic is not allowed in safe Go**

**165** What happens if you dereference a `nil` pointer?
a) Returns zero value  b) Compile error  c) Runtime panic  d) Silently ignored
**Answer: c) Runtime panic (invalid memory address or nil pointer dereference)**

**166** Why would you pass a pointer to a struct into a function instead of passing the struct by value?
a) To allow the function to mutate the original struct  b) To avoid copying a large struct (performance)  c) Both a and b  d) You shouldn't
**Answer: c) Both a and b — mutation and performance**

**167** What is the output of this code?
```go
func modify(x *int) {
	*x = 20
}
func main() {
	a := 10
	modify(&a)
	fmt.Println(a)
}
```
a) 10  b) 20  c) Compile error  d) 0
**Answer: b) 20**

**168** What does `new(int)` return?
a) `int` initialized to 0  b) `*int` pointing to a newly allocated zeroed integer  c) Compile error  d) `nil`
**Answer: b) `*int` pointing to an integer with value 0**

**169** Can you take the address of a literal value directly (e.g., `&5`)?
a) Yes  b) No, you can only take the address of an addressable value like a variable  c) Yes, but only for strings  d) Only inside structs
**Answer: b) No, literals are not addressable**

**170** What is the output of this code?
```go
func changePtr(p *int) {
	newVal := 100
	p = &newVal
}
func main() {
	a := 50
	changePtr(&a)
	fmt.Println(a)
}
```
a) 50  b) 100  c) Compile error  d) Panic
**Answer: a) 50 — The function modifies its local copy of the pointer `p`, not the value `a` points to.**

**171** Which built-in function is often used interchangeably with struct initialization to get a pointer, e.g., `&MyStruct{}`?
a) `make`  b) `new`  c) `alloc`  d) `ptr`
**Answer: b) `new(MyStruct)` is equivalent to `&MyStruct{}`**

---

### Level 8 — Coding Problems

**172** Which of the following functions correctly swaps the values of two integers using their pointers?
a) `func swap(a, b *int) { a, b = b, a }`
b) `func swap(a, b *int) { *a, *b = *b, *a }`
c) `func swap(a, b *int) { temp := a; a = b; b = temp }`
d) `func swap(a, b int) { a, b = b, a }`
**Answer: b) `func swap(a, b *int) { *a, *b = *b, *a }`**
```go
package main
import "fmt"
func swap(a, b *int) {
	temp := *a
	*a = *b
	*b = temp
}
func main() {
	x, y := 1, 2
	swap(&x, &y)
	fmt.Println(x, y) // Output: 2 1
}
```

**173** Why is it safe to return a pointer to a local variable (`&count`) from a function in Go, unlike in C/C++?
a) Go functions do not have stack frames.
b) Go performs escape analysis at compile time, moving the local variable to the heap if its reference escapes the function.
c) Go garbage collects the pointer immediately, avoiding memory leaks.
d) It is not safe; this will cause a runtime panic.
**Answer: b) Go performs escape analysis at compile time, moving the local variable to the heap if its reference escapes the function.**
```go
package main
import "fmt"
func createCounter() *int {
	count := 10
	return &count // Perfectly safe, `count` escapes to the heap
}
func main() {
	c := createCounter()
	fmt.Println(*c)
}
```---

# Fundamental Coding Scenarios: Data Types
### Hands-on scenario questions covering Go's fundamental data types.

---

# Strings

**TYPE.1 String Parsing for Configs**
**Question:** You are building a CLI tool. A user provides a configuration string in the format `"HOST:PORT"`. Write a function `ParseConfig(config string) (string, string)` that safely extracts the host and port. If the string is malformed, return `"localhost"` and `"8080"`.
**Answer:**
```go
package main

import (
	"fmt"
	"strings"
)

func ParseConfig(config string) (string, string) {
	parts := strings.Split(config, ":")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "localhost", "8080"
	}
	return parts[0], parts[1]
}

func main() {
	host, port := ParseConfig("api.goverse.com:443")
	fmt.Printf("Connecting to %s on port %s\n", host, port)
}
```

# Integers & Floats

**TYPE.2 Financial Rounding (Floats & Integers)**
**Question:** Floating-point math can be imprecise. You are calculating a 15% tax on a shopping cart total of `$45.67`. Write a function `CalculateTax(amount float64) int` that calculates the tax and returns the result in **cents** (integer) to avoid float precision loss in financial systems.
**Answer:**
```go
package main

import (
	"fmt"
	"math"
)

func CalculateTax(amount float64) int {
	tax := amount * 0.15
	// Convert to cents and round properly
	cents := int(math.Round(tax * 100))
	return cents
}

func main() {
	cents := CalculateTax(45.67)
	fmt.Printf("Tax in cents: %d\n", cents) // Output: 685
}
```

# Booleans

**TYPE.3 Flag Toggling System**
**Question:** You have three system flags: `isReady`, `hasData`, and `isError`. Write a function `SystemStatus(isReady, hasData, isError bool) string` that returns "ONLINE" if the system is ready and has data but no errors, "IDLE" if ready but no data and no errors, and "OFFLINE" otherwise.
**Answer:**
```go
package main

import "fmt"

func SystemStatus(isReady, hasData, isError bool) string {
	if isError {
		return "OFFLINE"
	}
	if isReady && hasData {
		return "ONLINE"
	}
	if isReady && !hasData {
		return "IDLE"
	}
	return "OFFLINE"
}

func main() {
	fmt.Println("Status:", SystemStatus(true, true, false)) // Output: ONLINE
}
```

# Complex Numbers

**TYPE.4 Signal Processing (Complex Types)**
**Question:** Go natively supports complex numbers (`complex64`, `complex128`). You are simulating a basic signal phase shift. Write a function `PhaseShift(signal complex128) complex128` that takes a signal and rotates it by 90 degrees (multiplying it by the imaginary number `1i`).
**Answer:**
```go
package main

import "fmt"

func PhaseShift(signal complex128) complex128 {
	// Multiply by 0 + 1i
	return signal * 1i
}

func main() {
	// Signal with real part 5, imaginary part 2
	var signal complex128 = 5 + 2i
	shifted := PhaseShift(signal)
	fmt.Printf("Original: %v, Shifted: %v\n", signal, shifted)
	// Output: Original: (5+2i), Shifted: (-2+5i)
}
```

---

# Frequently Asked Interview MCQs

**174** What is the default behavior of a channel when you send data to it without a receiver?
a) It buffers the data up to 10 elements.
b) It blocks the sending goroutine forever or until a receiver is ready.
c) It returns an error.
d) It drops the data.
**Answer: b) It blocks.** Unbuffered channels block until the other side is ready.

**175** How does Go achieve inheritance?
a) Using the `extends` keyword.
b) Using struct embedding (composition).
c) Go does not support object-oriented programming.
d) Using abstract classes.
**Answer: b) Struct embedding.** Go uses composition over inheritance.

**176** Which of the following statements about `defer` is true?
a) Deferred functions are executed in FIFO order.
b) Deferred functions are executed when the program exits.
c) Deferred functions are executed in LIFO order just before the surrounding function returns.
d) Arguments to deferred functions are evaluated at execution time.
**Answer: c) Executed in LIFO order.** Also note arguments are evaluated at defer-time.

**177** What happens if you run a Go program with a data race?
a) It deterministically panics.
b) It prints a warning but continues.
c) It has undefined behavior unless compiled with `-race` which will panic.
d) Go prevents data races at compile time.
**Answer: c) Undefined behavior.** The race detector (`-race`) helps find them.

**178** What does `sync.WaitGroup` do?
a) Limits the number of goroutines running concurrently.
b) Waits for a collection of goroutines to finish executing.
c) Blocks a channel until data is available.
d) Replaces the need for a Mutex.
**Answer: b) Waits for a collection of goroutines to finish.**

**179** Can you return a pointer to a local variable safely in Go?
a) No, it will cause a segmentation fault.
b) Yes, Go's escape analysis moves the variable to the heap.
c) Yes, but only for primitive types.
d) No, it creates a memory leak.
**Answer: b) Yes, thanks to escape analysis.**

**180** What is the correct way to initialize an empty slice with a specific capacity but zero length?
a) `s := make([]int, 0, 10)`
b) `s := make([]int, 10)`
c) `s := []int{10}`
d) `s := new([]int)`
**Answer: a) `make([]int, 0, 10)`**

**181** What happens when you read from a closed channel?
a) Panic.
b) It blocks forever.
c) It yields the zero value of the channel's type immediately.
d) It returns a compiler error.
**Answer: c) Yields the zero value.** (You can also use the comma-ok idiom to check if it's open).

**182** What is an empty interface (`interface{}`) used for in Go?
a) To represent a lack of data (null).
b) To define a struct with no fields.
c) To hold values of any type.
d) To define functions with no arguments.
**Answer: c) To hold values of any type.** Since every type implements zero methods, every type satisfies the empty interface.

**183** How are map keys compared in Go?
a) Using a hash function provided by the developer.
b) They must implement the `Comparable` interface.
c) Map keys can be of any type, including slices.
d) Map keys must be of a comparable type (e.g., int, string, pointer, struct without slices/maps).
**Answer: d) Keys must be comparable.** You cannot use slices, maps, or functions as map keys.

**184** What is the primary difference between a value receiver and a pointer receiver in a Go method?
a) Value receivers can mutate the original struct; pointer receivers cannot.
b) Pointer receivers avoid copying the struct and allow mutation of the original struct.
c) They are completely interchangeable with no performance or behavioral difference.
d) Value receivers are required for interfaces, pointer receivers are not.
**Answer: b) Pointer receivers allow mutation and avoid copying.** Use pointer receivers if the method needs to mutate the receiver or if the struct is large.

**185** What does the `init()` function do in a Go package?
a) It acts as the main entry point of the application.
b) It is executed once per file when the package is initialized, before `main()`.
c) It must be called manually to initialize variables.
d) It runs concurrently in a separate goroutine.
**Answer: b) Executed automatically during package initialization.** You can have multiple `init()` functions in a single package or even a single file.

**186** How does the `select` statement behave if multiple `case` channels are ready at the same time?
a) It executes the first one listed top-to-bottom.
b) It panics.
c) It chooses one pseudo-randomly.
d) It executes all of them concurrently.
**Answer: c) It chooses one pseudo-randomly.** This prevents starvation of cases further down the list.

**187** When passing a `map` to a function, what is actually being passed?
a) A deep copy of the entire map.
b) A pointer to the map descriptor, meaning modifications inside the function affect the original map.
c) A read-only copy of the map.
d) Maps cannot be passed to functions.
**Answer: b) A pointer to the map descriptor.** Maps (like slices and channels) act as reference types, so mutating a map inside a function mutates the original.

**188** What is the initial stack size of a new goroutine in modern Go (>= 1.4)?
a) 2 KB
b) 8 KB
c) 1 MB
d) 2 MB
**Answer: a) 2 KB.** This incredibly small footprint allows Go programs to spawn hundreds of thousands of goroutines easily. The stack grows and shrinks dynamically as needed.

**189** What is a common way to cause a memory leak in Go using slices?
a) Appending to a slice in a loop.
b) Slicing a small portion of a massive array/slice and keeping it in memory.
c) Passing a slice to a function by value.
d) Using `make` instead of `new`.
**Answer: b) Slicing a small portion of a massive array.** The small slice retains a reference to the *entire* underlying array, preventing the garbage collector from freeing the massive array.

**190** What is the difference between `new(T)` and `make(T)`?
a) `new` allocates memory and returns a pointer; `make` initializes slices, maps, and channels and returns the value itself.
b) `new` is for primitives; `make` is for structs.
c) `new` returns an initialized object; `make` returns a zeroed object.
d) There is no difference; they are aliases.
**Answer: a) `new` allocates zeroed memory and returns a pointer; `make` initializes internal data structures for slices/maps/channels and returns the value.**

**191** How do you explicitly check if an interface value holds a specific underlying type?
a) By using a regular `if` statement like `if val == type`.
b) Using a Type Assertion: `v, ok := val.(SpecificType)`.
c) Using the `typeof()` function.
d) Interfaces cannot be checked at runtime.
**Answer: b) Using a Type Assertion.**

**192** What does `errors.Is(err, targetErr)` do differently from `err == targetErr`?
a) It compares the string values of the errors.
b) It panics if the errors are not equal.
c) It unwraps the error chain to see if `targetErr` exists anywhere in the chain.
d) It casts the error to a struct.
**Answer: c) It unwraps the error chain.** This was introduced in Go 1.13 and is the standard way to check wrapped errors.

**193** Which of the following is true about strings in Go?
a) They are mutable arrays of bytes.
b) They are immutable slices of bytes.
c) They are mutable slices of runes.
d) They are essentially linked lists of characters.
**Answer: b) They are immutable slices of bytes.** Once created, a string's contents cannot be changed.

**194** In the Go scheduler's G-P-M model, what does the 'P' stand for?
a) Process
b) Pointer
c) Logical Processor
d) Program Counter
**Answer: c) Logical Processor.** 'P' represents a logical processor (context). 'M' is an OS thread, and 'G' is a Goroutine.

**195** What is the "blank identifier" (`_`) used for in Go?
a) To define private variables.
b) To discard return values or avoid unused variable errors.
c) To act as a wildcard in regular expressions.
d) To define an untyped constant.
**Answer: b) To discard return values or avoid "unused variable/import" compile errors.**

**196** How can you enforce that a struct `MyStruct` implements an interface `MyInterface` at compile time?
a) By adding `implements MyInterface` to the struct declaration.
b) By assigning a blank identifier: `var _ MyInterface = (*MyStruct)(nil)`.
c) It is impossible to check at compile time.
d) By using the `CheckInterface()` standard library function.
**Answer: b) `var _ MyInterface = (*MyStruct)(nil)`.** This forces the compiler to verify that the pointer to `MyStruct` satisfies the interface.

**197** What happens if you panic inside a Goroutine, but put a `recover()` in the `main` Goroutine?
a) The `main` Goroutine catches the panic and continues safely.
b) The panic cannot cross Goroutine boundaries; the whole program crashes.
c) The panicked Goroutine dies quietly, the rest of the program continues.
d) It creates a deadlock.
**Answer: b) The program crashes.** A `recover()` must be placed in a `defer` block *within the same Goroutine* where the panic occurs.

**198** What is the correct way to perform simple, lock-free counter increments across multiple Goroutines?
a) Using a standard integer `count++`
b) Using `sync.Mutex` (though this introduces locking)
c) Using the `sync/atomic` package: `atomic.AddInt64(&count, 1)`
d) Using a `select` statement.
**Answer: c) Using `sync/atomic`.** The `atomic` package provides low-level, lock-free hardware-level synchronization, which is faster than Mutexes for simple counters.

**199** Which of the following is a key difference between `sync.Map` and a standard `map` wrapped in an `RWMutex`?
a) `sync.Map` is slower for all operations.
b) `sync.Map` is optimized for append-only data or disjoint sets of keys (e.g., caches where keys are written once and read many times).
c) `sync.Map` is strongly typed, standard maps are not.
d) `sync.Map` does not require `make()`.
**Answer: b) Optimized for specific concurrent workloads.** For heavily contended, write-heavy workloads, an RWMutex-wrapped map is often faster.

**200** If you have a buffered channel of size 5, and 3 elements are in it, what happens when you read from it?
a) It blocks until the channel is full.
b) It reads the first element immediately without blocking.
c) It reads all 3 elements at once.
d) It panics.
**Answer: b) It reads immediately.** Buffered channels only block on read when they are completely empty, and block on write when completely full.

**201** What is "variable shadowing" in Go?
a) Changing the type of a variable at runtime.
b) Declaring a new variable with the same name in an inner scope, effectively hiding the outer variable.
c) A security feature to hide memory addresses.
d) A compilation optimization technique.
**Answer: b) Declaring a new variable in an inner scope.** This often happens accidentally with the `:=` operator inside `if` statements or loops.

**202** Can a Go program have a memory leak?
a) No, the Garbage Collector prevents all memory leaks.
b) Yes, primarily through unclosed goroutines (Goroutine leaks) or keeping references to large objects (like sliced arrays).
c) Yes, but only if you use the `unsafe` package.
d) Yes, but the compiler catches them.
**Answer: b) Yes.** Goroutine leaks are one of the most common causes of memory leaks in Go, occurring when a Goroutine is blocked forever waiting on a channel that will never be written to or read from.

**203** What does the Context package (`context.Context`) provide in Go?
a) It replaces global variables for application state.
b) It provides deadlines, cancellation signals, and request-scoped values across API boundaries and goroutines.
c) It is primarily used for database migrations.
d) It manages memory allocation.
**Answer: b) Deadlines, cancellations, and request-scoped values.** It is the idiomatic way to handle request timeouts and cancellations in Go backend services.

**204** How are slices passed to functions in Go?
a) By value, meaning the entire backing array is copied.
b) By reference, meaning a pointer to the slice header is passed.
c) By value, meaning a copy of the slice header (pointer to array, length, and capacity) is passed.
d) Slices cannot be passed to functions.
**Answer: c) By value, copying the slice header.** This is why modifying the slice's contents affects the original array, but appending to a slice may not affect the original slice length or capacity if the backing array needs to grow.

**205** What is the purpose of `sync.Pool`?
a) To manage a pool of database connections.
b) To cache allocated but unused items for later reuse, relieving pressure on the garbage collector.
c) To limit the number of goroutines running concurrently.
d) To synchronize network requests across multiple servers.
**Answer: b) To cache allocated but unused items.** This is particularly useful for objects that are expensive to allocate and frequently created and destroyed, like buffers.

**206** In Go 1.18+, what is the difference between `any` and `interface{}`?
a) `any` is strictly evaluated at compile time, `interface{}` at runtime.
b) `any` can only be used with generics, `interface{}` cannot.
c) There is no difference; `any` is simply an alias for `interface{}`.
d) `any` prevents type assertions.
**Answer: c) There is no difference.** `any` is simply a type alias for `interface{}` introduced to make code, especially generics, more readable.

**207** What is the correct way to prevent a struct from being compared using the `==` operator?
a) Add a field of type `[]int` (or any non-comparable type, like a function or map) to the struct.
b) Declare the struct as `uncomparable type {}`.
c) It is not possible; all structs are inherently comparable.
d) Use the `sync.Mutex` inside the struct.
**Answer: a) Add a field of a non-comparable type.** Types like slices, maps, and functions are not comparable. A common idiom is adding `_ [0]func()` to explicitly make it uncomparable without consuming memory.

**208** How does `defer` handle panic within the same function?
a) The deferred function is not executed if a panic occurs.
b) The deferred function is executed, and it can optionally capture the panic using `recover()`.
c) The deferred function panics immediately as well.
d) The program crashes before the deferred function runs.
**Answer: b) The deferred function is executed.** This guarantees cleanup and allows the `recover()` function to capture and handle the panic state.


### Level 9: Deep Dive - Arrays, Slices & Memory Allocation

**209** What happens when you pass an array (not a slice) of 10,000 integers to a Go function?
a) A pointer to the array is passed, which is very efficient.
b) The entire array of 10,000 integers is copied to the function's stack frame.
c) Go automatically converts the array to a slice before passing.
d) The program will panic due to stack overflow.
**Answer: b) The entire array is copied.** Arrays in Go are value types. Passing a large array by value is inefficient; this is why slices are preferred for passing collections.

**210** If `s` is a slice, what is the capacity of `s2 := s[2:4]`?
a) It is exactly 2.
b) It is the capacity of `s` minus 2.
c) It is the length of `s` minus 2.
d) It becomes 0 because it's a new slice.
**Answer: b) It is the capacity of `s` minus 2.** Slicing a slice `s[i:j]` results in a slice with length `j-i` and capacity `cap(s)-i`.

**211** What is a memory leak associated with Go slices?
a) Slices don't get garbage collected until the main function ends.
b) Slicing a large array creates a new slice that keeps the entire backing array in memory, even if only a tiny portion is referenced.
c) Slices created with `make()` are never garbage collected.
d) Appending to a slice infinitely causes memory leaks.
**Answer: b) Slicing a large array keeps the backing array alive.** If you read a 1GB file into memory, slice out just 10 bytes, and keep that 10-byte slice, the entire 1GB backing array remains in memory until the 10-byte slice is garbage collected.

**212** How can you avoid the slice memory leak mentioned in the previous question?
a) By explicitly calling `runtime.GC()`.
b) By setting the slice to `nil` when done.
c) By allocating a new slice of the exact size needed and using `copy()` to copy only the required elements.
d) By using the `unsafe` package to free the unused memory.
**Answer: c) By allocating a new slice and using `copy()`.** This ensures the new slice has its own small backing array, allowing the large original array to be garbage collected.

**213** What is a "Full Slice Expression" in Go (e.g., `s[low:high:max]`)?
a) It allows slicing backwards.
b) It allows you to specify the length and explicitly limit the capacity of the new slice.
c) It automatically shrinks the backing array.
d) It creates a two-dimensional slice.
**Answer: b) It allows you to explicitly limit the capacity.** By setting `max`, the capacity of the new slice becomes `max - low`. This prevents accidental modification of the original backing array via `append()`.

**214** When `append()` causes a slice to exceed its capacity, how much does the capacity typically grow?
a) By exactly 1 element.
b) By a fixed chunk size of 4KB.
c) By doubling its previous capacity (for smaller slices).
d) It never grows automatically; you must use `make()`.
**Answer: c) By doubling its capacity (for smaller slices).** Historically, it doubled until 1024 elements, after which it grew by ~25%. In modern Go (1.18+), the growth factor transitions more smoothly, but for small slices, it doubles.

**215** Is it safe to concurrently read and write to the same map in Go?
a) Yes, Go maps are inherently thread-safe.
b) No, it will result in a fatal runtime panic.
c) Yes, but only if you are using Go 1.20 or later.
d) No, but the data will just be silently corrupted.
**Answer: b) No, it will result in a fatal runtime panic.** Go's map implementation includes a concurrency check that detects concurrent reads and writes, resulting in a fatal error that cannot be recovered from.

**216** Which data structure provides concurrent, thread-safe map access without external mutexes?
a) `map[string]interface{}`
b) `sync.Map`
c) `sync.MutexMap`
d) `container/map`
**Answer: b) `sync.Map`.** It is optimized for use cases where entries are written once but read many times (e.g., caches). For write-heavy workloads, a standard map protected by a `sync.RWMutex` is often faster.

**217** What is the underlying data structure of a Go slice?
a) A linked list.
b) A hash table.
c) A struct containing a pointer to an array, the length, and the capacity.
d) A continuous block of memory managed directly by the CPU.
**Answer: c) A struct containing a pointer to an array, length, and capacity.** This is known as a slice header (`reflect.SliceHeader`).

**218** If you define a constant array in Go, e.g. `const arr = [3]int{1, 2, 3}`, what happens?
a) The compiler optimizes it for fast access.
b) It creates an immutable array on the heap.
c) It is an invalid operation; Go does not support constant arrays.
d) It behaves like a slice.
**Answer: c) It is an invalid operation.** In Go, constants can only be booleans, runes, numbers, or strings. You cannot have constant arrays, slices, maps, or structs.

### Level 10: Deep Dive - Data Structures & Algorithms in Go

**219** Which built-in Go package provides implementations for Heaps, Linked Lists, and Ring Buffers?
a) `collections`
b) `container`
c) `structs`
d) `data`
**Answer: b) `container`.** Specifically, `container/heap`, `container/list`, and `container/ring`.

**220** To implement a Min-Heap using `container/heap`, your type must implement which interface?
a) `sort.Interface` (Len, Less, Swap) plus Push and Pop.
b) `heap.MinHeap`
c) `sort.Sortable`
d) `container.Queue`
**Answer: a) `sort.Interface` plus Push and Pop.** The `heap.Interface` embeds `sort.Interface`, meaning your type must define how to compare elements (`Less`) and how to add/remove them.

**221** How do you efficiently reverse a slice in Go?
a) Use `slices.Reverse()` from the standard library (Go 1.21+).
b) Use the `reverse` keyword.
c) Append it backwards to a new slice.
d) Use `sort.Reverse()`.
**Answer: a) Use `slices.Reverse()`.** Go 1.21 introduced the `slices` package which provides generic functions for slice manipulation, including `Reverse`. Before 1.21, a two-pointer loop swapping elements was required.

**222** What is the time complexity of appending to a Go slice?
a) O(1) in all cases.
b) O(N) in all cases.
c) Amortized O(1), but O(N) worst-case when reallocation is needed.
d) O(log N).
**Answer: c) Amortized O(1).** Most appends are O(1) as they just place an item in the existing backing array. When capacity is exceeded, Go allocates a new array and copies existing elements (O(N)), but because it doubles capacity, this happens infrequently enough to be amortized O(1).

**223** Which is the most memory-efficient way to represent a set of items (where only uniqueness matters) in Go?
a) `map[T]bool`
b) `[]T` with manual duplicate checking.
c) `map[T]struct{}`
d) `sync.Map`
**Answer: c) `map[T]struct{}`.** An empty struct (`struct{}`) consumes exactly 0 bytes of memory. Using `bool` consumes 1 byte per element.

**224** How does the Go `sort` package handle sorting algorithms internally?
a) It strictly uses QuickSort.
b) It strictly uses MergeSort.
c) It uses Pattern-Defeating Quicksort (pdqsort) since Go 1.19.
d) It uses BubbleSort for small arrays and HeapSort for large ones.
**Answer: c) It uses Pattern-Defeating Quicksort (pdqsort).** pdqsort provides O(N log N) worst-case performance while being significantly faster than standard quicksort on data with existing patterns or sorted segments.

**225** How would you implement a Queue (FIFO) using standard Go slices?
a) `queue = append(queue, item)` to enqueue, `item = queue[0]; queue = queue[1:]` to dequeue.
b) `queue = append([]T{item}, queue...)` to enqueue, `item = queue[len(queue)-1]` to dequeue.
c) Using the `queue` package.
d) By using a doubly linked list only.
**Answer: a) `append` for enqueue, slicing `[1:]` for dequeue.** However, note that continuously slicing from the front can leave the backing array continuously growing until reallocated. For high-performance queues, a Ring Buffer is preferred.

**226** What is a Ring Buffer, and when would you use it over a standard slice?
a) A buffer that encrypts data in a ring.
b) A fixed-size buffer that wraps around, useful for high-performance, lock-free queues where you don't want continuous memory allocations.
c) A buffer used exclusively for network sockets.
d) A slice that automatically sorts itself.
**Answer: b) A fixed-size buffer that wraps around.** It is heavily used in networking, audio processing, and high-performance concurrency patterns to avoid memory allocations and garbage collection overhead.

**227** If you need to search a large, unsorted slice for a specific element, what is the time complexity?
a) O(1)
b) O(log N)
c) O(N)
d) O(N log N)
**Answer: c) O(N).** For an unsorted slice, you must iterate through potentially every element (Linear Search) to find a match.

**228** How do you perform a Binary Search in Go?
a) Use `slices.BinarySearch()` (Go 1.21+) or `sort.Search()` on a sorted slice.
b) Use `search.Binary()`.
c) It happens automatically if you use `map`.
d) Write a custom recursive function; there is no standard library support.
**Answer: a) `slices.BinarySearch()` or `sort.Search()`.** The slice must be sorted first. This reduces the search time complexity to O(log N).

### Level 11: Deep Dive - Goroutines & The Go Scheduler

**229** What architectural model does the Go Scheduler use?
a) 1:1 threading (One OS thread per goroutine).
b) M:1 threading (All goroutines run on a single OS thread).
c) M:N scheduling (M goroutines are multiplexed onto N OS threads).
d) Event-loop driven architecture (like Node.js).
**Answer: c) M:N scheduling.** This allows Go to manage millions of lightweight goroutines while efficiently utilizing a small number of heavier OS threads across available CPU cores.

**230** What does `GOMAXPROCS` control?
a) The maximum number of goroutines that can be created.
b) The number of OS threads that can execute user-level Go code simultaneously.
c) The maximum amount of memory the GC is allowed to use.
d) The number of background network threads.
**Answer: b) The number of OS threads executing user-level code.** By default, this is set to the number of logical CPUs on the machine.

**231** In the context of the Go scheduler, what are G, M, and P?
a) Go, Mutex, Pointer.
b) Goroutine (G), Machine/OS Thread (M), Processor/Context (P).
c) Garbage, Memory, Pacing.
d) Group, Map, Package.
**Answer: b) Goroutine, Machine, Processor.** The P holds the local run queue of G's. The M executes G's by attaching to a P.

**232** What is "Work Stealing" in the Go Scheduler?
a) When a hacker steals compute resources.
b) When one P (Processor) runs out of work in its local queue, it looks at other Ps' queues and "steals" half of their pending goroutines to maintain load balancing.
c) When the OS interrupts a goroutine to run another process.
d) When a goroutine bypasses a Mutex lock.
**Answer: b) Load balancing by stealing goroutines.** This ensures all CPU cores remain busy even if one core's queue is emptied rapidly.

**233** How much memory does a newly spawned Goroutine typically consume in modern Go?
a) 2 MB (Same as an OS thread).
b) 8 KB.
c) 2 KB.
d) It depends strictly on the system RAM.
**Answer: c) 2 KB.** Goroutines start with an extremely small, resizable stack. This is why you can easily spawn hundreds of thousands of them without running out of memory.

**234** What causes a Goroutine's stack to grow?
a) Spawning more goroutines.
b) Deeply nested function calls or declaring large variables inside functions.
c) Network latency.
d) Using too many channels.
**Answer: b) Deeply nested calls or large local variables.** If the 2 KB stack is insufficient, Go automatically copies the stack to a larger memory block and updates all pointers.

**235** What happens when a Goroutine makes a blocking Syscall (e.g., reading a file from disk)?
a) The entire OS thread (M) and the Processor (P) block until the syscall finishes.
b) The Go scheduler detects the block, detaches the M from the P, creates/wakes up a new M to attach to the P, and continues executing other goroutines.
c) Go crashes if a syscall takes more than 10ms.
d) The syscall is automatically made asynchronous.
**Answer: b) The scheduler detaches the M and assigns a new M to the P.** This ensures that other goroutines on that Processor's queue are not starved while waiting for the blocking syscall.

**236** What is "Cooperative Preemption" in Go 1.14+?
a) Goroutines must manually yield control using `runtime.Gosched()`.
b) The scheduler uses OS signals (SIGURG) to asynchronously preempt long-running goroutines, preventing them from monopolizing a CPU core.
c) Goroutines take turns executing based on priority queues.
d) Goroutines negotiate lock ownership.
**Answer: b) Asynchronous preemption using signals.** Before 1.14, a tight loop without function calls could freeze a thread forever. Go 1.14 introduced signal-based preemption to interrupt even tight CPU loops.

**237** Which function explicitly yields the processor, allowing other goroutines to run?
a) `runtime.Yield()`
b) `runtime.Gosched()`
c) `sync.Yield()`
d) `go.Defer()`
**Answer: b) `runtime.Gosched()`.** It pauses the current goroutine, puts it back into the global run queue, and executes the next available goroutine.

**238** Can you prioritize one Goroutine over another?
a) Yes, using `runtime.SetPriority(g, high)`.
b) Yes, by assigning it to a specific core.
c) No, the Go scheduler does not expose goroutine priorities to the developer; all goroutines are treated equally.
d) Yes, by using unbuffered channels exclusively.
**Answer: c) No, the scheduler does not expose priorities.** Go's design philosophy encourages letting the scheduler manage execution dynamically.

### Level 12: Deep Dive - Advanced Concurrency Patterns & Worker Pools

**239** What is a Worker Pool in Go?
a) A slice of network connections.
b) A pattern where a fixed number of goroutines (workers) continuously pull tasks from a shared channel and process them concurrently.
c) A database connection manager.
d) A pool of memory used by the Garbage Collector.
**Answer: b) A pattern of fixed goroutines pulling tasks from a channel.** This prevents the system from being overwhelmed by spawning thousands of simultaneous goroutines if a burst of tasks arrives.

**240** Why would you use a Worker Pool instead of just spawning a new goroutine for every single task?
a) Spawning a goroutine for every task is an anti-pattern and illegal in Go.
b) To control concurrency limits, prevent CPU/Memory exhaustion, and protect downstream dependencies (like a database) from being flooded with connections.
c) Worker pools are actually slower, but they use less bandwidth.
d) To prevent panic errors in the main function.
**Answer: b) To control concurrency limits and prevent resource exhaustion.** While goroutines are cheap, the resources they access (DBs, network sockets, APIs) are not. A worker pool acts as a throttle.

**241** In a standard Worker Pool, how do workers know when to stop?
a) You call `worker.Stop()`.
b) When the task channel they are reading from is closed (`close(jobs)`).
c) When `runtime.NumGoroutine()` reaches 0.
d) They timeout automatically after 5 minutes.
**Answer: b) When the task channel is closed.** A `range` loop over a channel (`for task := range jobs`) automatically exits when the channel is closed and drained, gracefully shutting down the worker.

**242** What is the Fan-Out concurrency pattern?
a) Merging multiple channels into one.
b) Starting multiple goroutines to handle input from a single channel, effectively distributing the workload.
c) Sending a HTTP request to multiple servers.
d) Closing a channel gracefully.
**Answer: b) Distributing workload from a single channel to multiple goroutines.** This is exactly what a Worker Pool does.

**243** What is the Fan-In concurrency pattern?
a) Reading from multiple channels and multiplexing all their outputs into a single output channel.
b) Combining multiple struct fields into one.
c) Funneling network requests to a single load balancer.
d) Using a WaitGroup to stop workers.
**Answer: a) Multiplexing multiple channels into one.** This is typically achieved by launching a goroutine for each input channel that writes to the single output channel, or by using a `select` statement.

**244** When implementing Fan-In, how do you know when to close the multiplexed output channel?
a) Close it immediately after reading the first value.
b) You should never close the output channel.
c) Use a `sync.WaitGroup` to wait for all the input-reading goroutines to finish, then close the output channel.
d) The garbage collector closes it automatically.
**Answer: c) Use a `sync.WaitGroup`.** Wait for all publishers to finish their work, and then safely close the single output channel.

**245** What is the Pipeline pattern in Go?
a) A CI/CD deployment tool written in Go.
b) A series of stages (goroutines) connected by channels, where the output of one stage is the input of the next.
c) A way to pipe data directly to the OS shell.
d) An HTTP routing mechanism.
**Answer: b) A series of goroutine stages connected by channels.** Data streams through the pipeline, allowing different stages of processing to happen concurrently.

**246** In a concurrent environment, what does `sync.Once` guarantee?
a) A function will execute exactly once, even if called simultaneously by multiple goroutines.
b) A variable can only be mutated once.
c) A channel can only receive one message.
d) A goroutine runs for exactly one millisecond.
**Answer: a) A function executes exactly once.** It is thread-safe and highly optimized, perfect for initializing singleton resources (like DB connections or configuration loads).

**247** What is a "Race Condition" in Go?
a) When two goroutines finish at exactly the same time.
b) When two or more goroutines access the same memory location concurrently, and at least one access is a write, without explicit synchronization.
c) A race to see which server responds fastest.
d) When a channel buffer gets full.
**Answer: b) Unsynchronized concurrent memory access.** This leads to unpredictable behavior and silent data corruption.

**248** How can you detect Race Conditions in your Go code?
a) By writing unit tests with `testing.B`.
b) By compiling/running the application with the `-race` flag (e.g., `go run -race main.go`).
c) By running a linter like `golint`.
d) By using `runtime.NumCPU()`.
**Answer: b) By using the `-race` flag.** The Go Race Detector instrumentation catches memory access violations at runtime. It incurs significant overhead, so it should be used during testing, not in production.

**249** What does `sync.Cond` provide?
a) A condition variable, allowing goroutines to wait for or announce the occurrence of a specific state or event.
b) A conditional if/else statement for channels.
c) A way to pause the garbage collector conditionally.
d) A mutex that locks based on a boolean condition.
**Answer: a) A condition variable.** It allows multiple goroutines to be suspended (via `Wait()`) and then woken up (via `Signal()` or `Broadcast()`) when a shared condition changes.

**250** When building a Worker Pool, what happens if the jobs channel is unbuffered and there are no idle workers?
a) The sender will block until a worker becomes available to process the job.
b) The job will be dropped and permanently lost.
c) The Go runtime will automatically spawn a new worker goroutine.
d) The program will panic.
**Answer: a) The sender will block.** This provides natural "backpressure," preventing the system from taking on more work than it can handle.
