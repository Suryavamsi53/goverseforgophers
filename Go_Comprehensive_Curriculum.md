# Go Programming Curriculum — BEGINNER TIER
### Covers Level 1–6: Basics & Syntax, Control Flow, Functions, Arrays/Slices/Maps, Structs & Methods, Interfaces
### Format: Question → Answer shown immediately after

---

# LEVEL 1: Basics & Syntax

**1.1** What is the zero value of an `int` in Go?
a) null  b) 0  c) undefined  d) Compile error
**Answer: b) 0**

**1.2** What is the zero value of a `string` in Go?
a) `""`  b) `null`  c) `nil`  d) undefined
**Answer: a) `""` (empty string)**

**1.3** What is the zero value of a `bool` in Go?
a) `true`  b) `false`  c) `nil`  d) `0`
**Answer: b) `false`**

**1.4** Which is invalid inside a function body?
a) `var x int = 10`  b) `x := 10`  c) `x := 10` then `x := 20` in same scope  d) `var x = 10`
**Answer: c) — `:=` requires at least one new variable on the left; redeclaring the same single variable with `:=` in the same scope is a compile error**

**1.5** What does `const Pi = 3.14` create?
a) A typed float64 constant  b) An untyped constant that adapts to context  c) A reassignable variable  d) Compile error
**Answer: b) Untyped constant — takes the type required by context (float32, float64, etc.)**

**1.6** What happens if you declare `var x int` and never use it?
a) Nothing  b) Compile error  c) Runtime panic  d) Warning only
**Answer: b) Compile error — Go disallows unused local variables**

**1.7** What is `7 / 2` when both operands are `int`?
a) 3.5  b) 3  c) 4  d) Compile error
**Answer: b) 3 — integer division truncates**

**1.8** What does `iota` do inside a `const` block?
a) Generates random numbers  b) Auto-increments from 0 per line in the block  c) Marks immutability  d) Invalid syntax
**Answer: b) Auto-increments starting at 0 for each ConstSpec line**

**1.9** What is the type of `x := 5.0 / 2`?
a) int  b) float64  c) Compile error  d) float32
**Answer: b) float64 — untyped constant `5.0` forces floating-point division**

**1.10** What does `:=` do that `var` cannot?
a) Declare a constant  b) Infer type, only valid inside functions  c) Allow global declarations  d) Nothing, identical
**Answer: b) Short variable declaration — infers type, can only be used inside function bodies**

**1.11** Which statement about Go type conversion is true?
a) Implicit numeric conversion is allowed  b) Explicit conversion required, e.g. `float64(x)`  c) String+int auto-converts  d) Only pointers need conversion
**Answer: b) Explicit conversion required — Go has no implicit numeric widening**

**1.12** What does `x := 10; x := 20` (new scope, e.g. inside an if-block) do?
a) Compile error, x already declared  b) Shadows outer x within the block  c) Reassigns outer x  d) Undefined behavior
**Answer: b) Shadows — a new `x` is created scoped to the inner block**

**1.13** Which is a valid multiple-variable declaration with different types?
a) `var a, b int, string = 1, "x"`  b) `var ( a int = 1; b string = "x" )`  c) `var a, b = 1, "x"`  d) both b and c
**Answer: d) Both are valid — grouped var block, or `var a, b = 1, "x"` with inferred types**

**1.14** What is the output of `fmt.Println(10 % 3)`?
a) 3  b) 1  c) 0  d) 3.33
**Answer: b) 1 — modulo remainder**

**1.15** What does `const MaxRetries = 5` followed by `MaxRetries = 10` do?
a) Reassigns fine  b) Compile error, cannot assign to constant  c) Runtime panic  d) Creates shadow variable
**Answer: b) Compile error — constants are immutable**

**1.16** In `var b byte = 300`, what happens?
a) b becomes 300  b) Compile error — constant 300 overflows byte  c) b wraps to 44  d) Runtime panic
**Answer: b) Compile error — byte is uint8, max value 255, and this overflow is caught at compile time for constant expressions**

**1.17** What's the difference between `rune` and `byte` in Go?
a) No difference  b) `rune` is alias for int32 (Unicode code point), `byte` is alias for uint8  c) `byte` is signed, `rune` unsigned  d) `rune` is deprecated
**Answer: b) `rune` = int32 (Unicode code point), `byte` = uint8 (raw byte)**

**1.18** What does this print? `fmt.Println("go" + "lang")`
a) Compile error  b) "golang"  c) "go lang"  d) 5
**Answer: b) "golang" — string concatenation with `+`**

**1.19** Predict output:
```go
x := 5
y := &x
*y = 10
fmt.Println(x)
```
a) 5  b) 10  c) Compile error  d) nil
**Answer: b) 10 — `y` points to `x`'s address; dereferencing and assigning mutates `x`**

**1.20** What's wrong with this code?
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

**1.21** Write a program that declares `maxConnections` (int), `timeoutSeconds` (float64), and `serviceName` (string), then prints all three in a formatted single line using `fmt.Printf`.
**Answer:**
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

**1.22** Using `iota`, define byte-size constants: `KB`, `MB`, `GB` as powers of 1024 (bit-shift pattern common in infra code).
**Answer:**
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

**1.23** Write a program that takes memory in MB (int) and converts it to GB using proper float division, avoiding integer truncation.
**Answer:**
```go
package main

import "fmt"

func main() {
	memoryMB := 2560
	memoryGB := float64(memoryMB) / 1024
	fmt.Printf("%d MB = %.2f GB\n", memoryMB, memoryGB)
}
```

**1.24** Write a program simulating a rate limiter: given `maxRequestsPerMinute` = 600, compute and print `requestsPerSecond` as a float64.
**Answer:**
```go
package main

import "fmt"

func main() {
	maxRequestsPerMinute := 600
	requestsPerSecond := float64(maxRequestsPerMinute) / 60.0
	fmt.Printf("Requests/sec: %.2f\n", requestsPerSecond)
}
```

**1.25** Declare HTTP status-like constants `StatusOK = 200`, `StatusNotFound = 404`, `StatusServerError = 500` — explain why plain sequential `iota` can't produce these directly.
**Answer:**
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

**2.1** Does Go require parentheses around `if` conditions?
a) Yes, always  b) No, and braces are mandatory  c) No, but parentheses are also disallowed  d) Optional both ways
**Answer: b) No parens needed, but `{ }` braces are mandatory even for single statements**

**2.2** What does Go's `switch` do differently from C/Java by default?
a) Nothing, identical  b) No fallthrough by default — each case breaks automatically  c) Requires `break` explicitly  d) Only works on integers
**Answer: b) Cases don't fall through unless you use the `fallthrough` keyword**

**2.3** Which loop construct does Go NOT have?
a) `for`  b) `while`  c) `do-while`  d) both b and c
**Answer: d) Go only has `for` — no dedicated `while` or `do-while` keywords (achieved via `for` variants)**

**2.4** What does `for i := 0; i < 5; i++ {}` with empty body do?
a) Compile error  b) Infinite loop  c) Runs 5 times doing nothing, valid  d) Runs once
**Answer: c) Valid — loop runs 5 times with an empty body**

**2.5** What does `for { }` (no condition) do?
a) Compile error  b) Infinite loop  c) Runs zero times  d) Runs once
**Answer: b) Infinite loop, equivalent to `while(true)`**

**2.6** In a `switch` statement, what does a case with multiple values look like?
a) `case 1, 2, 3:`  b) `case 1 | 2 | 3:`  c) `case (1,2,3):`  d) Not supported
**Answer: a) `case 1, 2, 3:` — comma-separated values in one case**

**2.7** What is a "switch with no expression" used for?
a) Invalid syntax  b) Acts like an if-else chain, each case is a boolean condition  c) Only works on strings  d) Same as regular switch
**Answer: b) `switch { case x > 10: ... }` — clean alternative to long if-else chains**

**2.8** What does `continue` do inside a nested loop by default?
a) Breaks all loops  b) Skips to next iteration of innermost enclosing loop  c) Compile error in nested loops  d) Skips to next iteration of outer loop
**Answer: b) Affects only the innermost loop unless a label is used**

**2.9** How do you break out of an outer loop from within a nested inner loop?
a) `break 2`  b) Using a labeled break: `break OuterLoop`  c) Not possible in Go  d) `return`
**Answer: b) Labeled break, e.g. `OuterLoop: for {...}` then `break OuterLoop`**

**2.10** What does `range` return when iterating over a slice?
a) Only the value  b) Only the index  c) Index and value  d) A pointer to each element
**Answer: c) Index, value (in that order) — `for i, v := range slice`**

**2.11** What happens if you modify a slice element via `for _, v := range slice { v = 100 }`?
a) Modifies the original slice  b) `v` is a copy; original slice unchanged  c) Compile error  d) Panic
**Answer: b) `v` is a copy of the value — mutating it does not affect the underlying slice**

**2.12** Predict the output:
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

**2.13** What does `fallthrough` do in a switch case?
a) Nothing, invalid  b) Forces execution to continue into the next case block regardless of its condition  c) Skips remaining cases  d) Restarts the switch
**Answer: b) Forces the next case's body to execute unconditionally**

**2.14** Predict output:
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

**2.15** What's the idiomatic Go way to iterate over a map's keys and values?
a) `for k, v := range myMap`  b) `for k in myMap`  c) `myMap.each(...)`  d) Maps can't be iterated
**Answer: a) `for k, v := range myMap`**

**2.16** Is map iteration order guaranteed in Go?
a) Yes, insertion order  b) Yes, sorted key order  c) No, deliberately randomized  d) Yes, but only for small maps
**Answer: c) No — Go intentionally randomizes map iteration order to prevent reliance on it**

**2.17** What does this print?
```go
count := 0
for count < 5 {
	count++
}
fmt.Println(count)
```
a) 4  b) 5  c) Infinite loop  d) 0
**Answer: b) 5 — this is Go's "while-style" for loop (condition-only)**

**2.18** What is wrong (if anything) with:
```go
if x := getValue(); x > 0 {
	fmt.Println(x)
}
fmt.Println(x)
```
a) Nothing wrong  b) Compile error — `x` is scoped to the if statement and unavailable outside  c) Runtime panic  d) x is 0 outside
**Answer: b) Compile error — `x` declared in the if-statement's initializer is only in scope within the if/else block**

**2.19** In a health-check retry loop, what's idiomatic for "retry up to N times with early exit on success"?
a) Recursive function only  b) `for i := 0; i < maxRetries; i++ { if success { break } }`  c) `while` loop  d) `goto` only
**Answer: b) A bounded for-loop with a `break` on success is the idiomatic pattern**

**2.20** What does `goto` do in Go, and when is it discouraged?
a) Jumps to a labeled statement; discouraged for readability except rare cases like breaking nested loops/cleanup  b) Not supported in Go  c) Only works in switch statements  d) Same as break
**Answer: a) Go supports `goto` but idiomatic Go rarely uses it — labeled break/continue or restructuring is usually preferred**

---

### Level 2 — Coding Problems

**2.21** Write a function that classifies an HTTP status code into "Success", "Client Error", "Server Error", or "Unknown" using a switch statement with ranges.
**Answer:**
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
	fmt.Println(classify(200))
	fmt.Println(classify(503))
}
```

**2.22** Write a retry loop that attempts a (simulated) API call up to 3 times, breaking early on simulated success at attempt 2.
**Answer:**
```go
package main

import "fmt"

func callAPI(attempt int) bool {
	return attempt == 2 // simulate success on 2nd try
}

func main() {
	maxRetries := 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		fmt.Printf("Attempt %d...\n", attempt)
		if callAPI(attempt) {
			fmt.Println("Success!")
			break
		}
	}
}
```

**2.23** Write a nested loop that scans a 2D grid (e.g., simulating server rack positions) and uses a labeled break to stop entirely once a target server ID is found.
**Answer:**
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
				fmt.Printf("Found %d at [%d][%d]\n", target, row, col)
				break Search
			}
		}
	}
}
```

**2.24** Write a function using `switch` with `fallthrough` to build a cumulative permission string (e.g., "read" implies checking "write" too).
**Answer:**
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
	fmt.Println(permissions(3))
	fmt.Println(permissions(1))
}
```

---

# LEVEL 3: Functions

**3.1** How many values can a Go function return?
a) Only 1  b) Up to 2  c) Any number  d) Up to 4
**Answer: c) Any number of return values**

**3.2** What is a "named return"?
a) A function name that's a keyword  b) Return values declared in the function signature that act as local variables  c) A return value with a comment  d) Not valid Go
**Answer: b) e.g. `func div(a, b int) (result int, err error)` — `result` and `err` are pre-declared**

**3.3** When does `defer` execute?
a) Immediately  b) Just before the enclosing function returns, LIFO order  c) At program exit only  d) Before the function starts
**Answer: b) Deferred calls run in LIFO order right before the surrounding function returns**

**3.4** Predict output:
```go
func main() {
	defer fmt.Println("1")
	defer fmt.Println("2")
	defer fmt.Println("3")
}
```
a) 1 2 3  b) 3 2 1  c) Compile error  d) Only 3 prints
**Answer: b) 3 2 1 — LIFO (stack) order**

**3.5** What is a variadic function parameter?
a) A parameter with default value  b) `...T` — accepts zero or more arguments of type T as a slice  c) A pointer parameter  d) Not supported in Go
**Answer: b) `func sum(nums ...int)` accepts any number of int arguments, accessible as a slice inside**

**3.6** What is a closure in Go?
a) A function that closes files  b) A function value that references variables from outside its own body  c) Same as a method  d) A deprecated feature
**Answer: b) An anonymous function capturing variables from its enclosing scope**

**3.7** Predict output:
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

**3.8** Can you pass a slice `...int` variadic call using an existing slice?
a) No, must list elements individually  b) Yes, using `mySlice...` spread syntax  c) Yes, using `*mySlice`  d) Only with arrays
**Answer: b) `sum(mySlice...)` spreads the slice into variadic args**

**3.9** What happens if a deferred function's arguments reference a variable that changes later?
a) Deferred call sees the latest value  b) Arguments are evaluated at defer-time, not execution-time  c) Compile error  d) Runtime panic
**Answer: b) Arguments to deferred calls are evaluated immediately when `defer` executes, not when the call actually runs**

**3.10** Predict output:
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

**3.11** What's a common use of `defer` in real backend code?
a) Looping  b) Closing files/DB connections/unlocking mutexes reliably  c) Declaring constants  d) String formatting
**Answer: b) Resource cleanup — e.g. `defer file.Close()`, `defer mu.Unlock()`**

**3.12** Are Go functions first-class values?
a) No  b) Yes — can be assigned to variables, passed as arguments, returned from functions  c) Only named functions  d) Only methods
**Answer: b) Yes, functions are first-class citizens in Go**

**3.13** What does this function signature mean? `func process(data []byte) (result string, err error)`
a) Two required args  b) Takes bytes, returns a string and an error (named returns)  c) Invalid syntax  d) Takes no arguments
**Answer: b) Takes a byte slice, returns a named string result and named error**

**3.14** Can a Go function be recursive?
a) No  b) Yes, functions can call themselves  c) Only methods can  d) Only with special syntax
**Answer: b) Yes, standard recursion is supported**

**3.15** What's the zero value returned by named returns if you just call `return` with no values?
a) Compile error  b) The current values of the named return variables  c) Always nil  d) Always zero regardless of prior assignment
**Answer: b) Whatever the named return variables currently hold**

**3.16** Predict output:
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

**3.17** What is the difference between a function and a method in Go?
a) No difference  b) A method has a receiver argument, associating it with a type  c) Methods can't return values  d) Functions can't take structs
**Answer: b) Methods are functions with a receiver: `func (r ReceiverType) Name(...)`**

**3.18** Can you have a variadic parameter combined with regular parameters?
a) No  b) Yes, but variadic must be last: `func f(a int, b ...string)`  c) Yes, variadic must be first  d) Only in generics
**Answer: b) Variadic parameter must always be the last parameter**

**3.19** What does calling a variadic function with zero arguments do?
a) Compile error  b) The variadic parameter is an empty (nil) slice  c) Runtime panic  d) Uses default values
**Answer: b) It becomes a nil/empty slice, safely iterable with zero length**

**3.20** True or False: Deferred functions can modify named return values even after a panic is recovered.
a) True  b) False
**Answer: a) True — this is exactly how idiomatic panic-recovery-with-cleanup patterns work in Go**

---

### Level 3 — Coding Problems

**3.21** Write a variadic function `sum(nums ...int) int` that returns the total, and call it both with individual args and with a spread slice.
**Answer:**
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

**3.22** Write a function `safeDivide(a, b float64) (result float64, err error)` using named returns that returns an error instead of panicking on division by zero.
**Answer:**
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

**3.23** Write a closure-based rate limiter: a function `makeLimiter(max int)` that returns a function which returns true if calls are under the max, false otherwise (tracking count internally).
**Answer:**
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

**3.24** Write a function that uses `defer` to log entry/exit timing of a simulated DB query function.
**Answer:**
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

**3.25** Write a function using `recover()` inside a deferred closure to prevent a worker goroutine's panic from crashing the whole program, logging the recovered error instead.
**Answer:**
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

**4.1** What's the key difference between an array and a slice in Go?
a) No difference  b) Arrays have fixed size at compile time; slices are dynamic, backed by an array  c) Slices are always faster  d) Arrays can't hold structs
**Answer: b) Arrays are fixed-length value types; slices are flexible, reference-like views over an underlying array**

**4.2** What does `len()` vs `cap()` return for a slice?
a) Same thing  b) `len` = number of elements currently in use; `cap` = size of underlying array from the slice's start  c) `cap` is always double `len`  d) `len` is for arrays only
**Answer: b) len = current element count, cap = max elements before reallocation is needed**

**4.3** What happens when you append beyond a slice's capacity?
a) Panic  b) Go allocates a new, larger underlying array and copies data  c) Silently truncates  d) Compile error
**Answer: b) A new array is allocated (typically growth factor ~2x for small slices) and the old data is copied over**

**4.4** Predict output:
```go
a := []int{1, 2, 3}
b := a
b[0] = 100
fmt.Println(a[0])
```
a) 1  b) 100  c) Compile error  d) 0
**Answer: b) 100 — slices share the same underlying array; `b := a` copies the slice header, not the data**

**4.5** What's the zero value of a slice?
a) Empty slice `[]int{}`  b) `nil`  c) Compile error, must initialize  d) Array of zeros
**Answer: b) `nil` — a nil slice has len=0, cap=0, and is usable (e.g., append works on it)**

**4.6** How do you create a slice with initial length 5 and capacity 10?
a) `make([]int, 5, 10)`  b) `make([]int, 10, 5)`  c) `[]int{5, 10}`  d) `new([]int, 5, 10)`
**Answer: a) `make([]int, len, cap)`**

**4.7** What does map lookup return for a missing key?
a) Panic  b) Zero value of the value type, and `false` if using the "comma ok" form  c) nil always  d) Compile error
**Answer: b) e.g. `v, ok := myMap["missing"]` — v is zero value, ok is false**

**4.8** What happens if you write to a nil map?
a) Silently no-ops  b) Panic: "assignment to entry in nil map"  c) Auto-initializes  d) Compile error
**Answer: b) Panics at runtime — nil maps must be initialized with `make` or a literal before writing**

**4.9** Can you read from a nil map?
a) No, panics  b) Yes, returns zero value  c) Compile error  d) Only with comma-ok
**Answer: b) Reading from a nil map is safe and returns the zero value (unlike writing)**

**4.10** What does `a[1:3]` do for `a := []int{10, 20, 30, 40}`?
a) `[20, 30]`  b) `[10, 20, 30]`  c) `[20, 30, 40]`  d) `[30, 40]`
**Answer: a) `[20, 30]` — slicing is [low:high), high exclusive**

**4.11** What is the danger of slicing a large array/slice and keeping only a small sub-slice long-term?
a) None  b) The sub-slice keeps the entire original backing array alive, preventing GC of the rest  c) It copies data unnecessarily  d) Not possible in Go
**Answer: b) Memory leak risk — small slices can pin large backing arrays in memory unless explicitly copied**

**4.12** How do you safely copy a slice to avoid backing-array sharing issues?
a) `b := a`  b) `b := make([]int, len(a)); copy(b, a)`  c) `b := a[:]`  d) Not possible
**Answer: b) `make` + `copy` creates an independent backing array**

**4.13** What does `delete(myMap, "key")` do if "key" doesn't exist?
a) Panics  b) No-op, safe  c) Returns an error  d) Compile error
**Answer: b) Safe no-op — deleting a non-existent key does nothing**

**4.14** Are Go maps safe for concurrent read/write from multiple goroutines?
a) Yes, always  b) No — concurrent map writes (or write+read) cause a runtime panic/race  c) Only reads are safe, writes need locking too but reads are always fine  d) Yes, if buffered
**Answer: b) No — maps are not safe for concurrent use without external synchronization (sync.Mutex or sync.Map)**

**4.15** Predict output:
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

**4.16** What's the idiomatic way to check if a key exists in a map without caring about its value?
a) `if myMap["key"] != nil`  b) `if _, ok := myMap["key"]; ok`  c) `if myMap.has("key")`  d) `if len(myMap["key"]) > 0`
**Answer: b) The "comma ok" idiom**

**4.17** What does `append(a, b...)` do when both `a` and `b` are `[]int`?
a) Compile error  b) Appends all elements of b onto a  c) Appends b as a single nested element  d) Only works with `copy`
**Answer: b) Spreads b's elements and appends them individually**

**4.18** What is the output?
```go
m := map[string]int{"a": 1, "b": 2}
m["c"] = 3
fmt.Println(len(m))
```
a) 2  b) 3  c) Compile error  d) 0
**Answer: b) 3 — map now has three key-value pairs**

**4.19** What's a struct-keyed map used for in real systems?
a) Not allowed in Go  b) Using composite keys (e.g., struct{UserID, ResourceID}) for caching/dedup logic  c) Only string keys are allowed  d) Structs can't be map keys unless they implement an interface
**Answer: b) Structs with only comparable fields can be map keys — common for composite-key caches**

**4.20** What does `make([]int, 0, 100)` optimize for in code that appends heavily?
a) Nothing  b) Pre-allocates capacity to avoid repeated reallocation during appends  c) Creates a fixed array  d) Wastes memory always
**Answer: b) Pre-allocating capacity avoids multiple reallocations/copies as elements are appended — common perf optimization**

---

### Level 4 — Coding Problems

**4.21** Write a function that deduplicates a slice of strings using a map as a set.
**Answer:**
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

**4.22** Write a function that counts word frequency in a slice of strings (simulating log-line tokens) and returns a `map[string]int`.
**Answer:**
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

**4.23** Write a function `batch(items []int, size int) [][]int` that splits a slice into chunks of a given size (common for batch-processing API requests).
**Answer:**
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

**4.24** Write a function that safely merges two `map[string]int` config maps, where the second map's values override the first's on key conflicts.
**Answer:**
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
	fmt.Println(mergeConfigs(base, override))
}
```

**4.25** Write a function demonstrating the slice-capacity growth pitfall: show how appending to a sliced sub-slice can unexpectedly mutate a sibling slice sharing the same backing array.
**Answer:**
```go
package main

import "fmt"

func main() {
	original := make([]int, 3, 5)
	original[0], original[1], original[2] = 1, 2, 3

	a := original[:2] // len=2, cap=5 (shares backing array)
	a = append(a, 999) // overwrites original[2] since cap allows it in-place

	fmt.Println("original:", original) // [1 2 999] <- unexpectedly mutated
	fmt.Println("a:", a)
}
```

---

# LEVEL 5: Structs and Methods

**5.1** How do you define a struct in Go?
a) `struct MyStruct { ... }`  b) `type MyStruct struct { Field1 Type1; Field2 Type2 }`  c) `class MyStruct { ... }`  d) `interface MyStruct { ... }`
**Answer: b)**

**5.2** What's the difference between a value receiver and a pointer receiver on a method?
a) No difference  b) Value receiver operates on a copy; pointer receiver can mutate the original  c) Pointer receivers are always faster  d) Value receivers can't call other methods
**Answer: b) `func (s Struct) Method()` copies; `func (s *Struct) Method()` operates on original via pointer**

**5.3** Predict output:
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

**5.4** Can you embed one struct inside another in Go (composition)?
a) No, Go has no inheritance-like feature  b) Yes, via anonymous struct fields (embedding)  c) Only with interfaces  d) Only pointers can be embedded
**Answer: b) Yes — embedding promotes the embedded struct's fields/methods to the outer struct**

**5.5** What does struct embedding provide that's similar to inheritance?
a) True polymorphic override  b) Field and method promotion — outer struct can access embedded struct's members directly  c) Nothing, purely cosmetic  d) Multiple dispatch
**Answer: b) Promotion, not true inheritance — there's no dynamic override mechanism**

**5.6** What is a struct tag, e.g. `json:"name"`, used for?
a) Comments only  b) Metadata read via reflection, commonly for encoding/decoding (JSON, DB, etc.)  c) Compiler directives  d) Not valid Go syntax
**Answer: b) Struct tags provide metadata used by packages like encoding/json via reflection**

**5.7** What's the zero value of a struct?
a) nil  b) All fields set to their respective zero values  c) Compile error  d) Empty struct{}
**Answer: b) Every field takes its own zero value (0, "", false, nil, etc.)**

**5.8** Can structs be compared with `==`?
a) Never  b) Yes, if all fields are comparable types  c) Only pointer comparisons  d) Only with reflect.DeepEqual
**Answer: b) Structs are comparable if every field's type is comparable (no slices/maps/funcs as fields)**

**5.9** Predict output:
```go
type Point struct{ X, Y int }
p1 := Point{1, 2}
p2 := Point{1, 2}
fmt.Println(p1 == p2)
```
a) false  b) true  c) Compile error  d) panic
**Answer: b) true — struct equality compares all fields**

**5.10** What does `new(MyStruct)` return?
a) A value of type MyStruct  b) A pointer `*MyStruct` to a zero-valued struct  c) nil  d) Compile error
**Answer: b) `new()` allocates and returns a pointer to a zeroed value**

**5.11** How do you create a struct instance with named fields?
a) `Point{1, 2}` positional only  b) `Point{X: 1, Y: 2}`  c) Both a and b are valid  d) Neither is valid
**Answer: c) Both positional and named-field struct literals are valid Go**

**5.12** What happens when you pass a struct to a function by value (not pointer)?
a) Struct is passed by reference automatically  b) A full copy of the struct is made  c) Compile error for large structs  d) Only the first field is copied
**Answer: b) Go copies the entire struct when passed by value — can be costly for large structs**

**5.13** What does `func (s *Service) Start() error` typically indicate about design intent?
a) Nothing special  b) The method likely mutates the receiver's state (e.g., connection setup) so uses a pointer receiver  c) It's always faster than a value receiver  d) It must be called on a nil Service
**Answer: b) Pointer receivers are conventionally used when the method needs to mutate state or the struct is large**

**5.14** Can an embedded struct's method be "overridden" by defining a same-named method on the outer struct?
a) No, not possible  b) Yes — the outer struct's own method shadows the promoted one when called directly on the outer type  c) Compile error, ambiguous  d) Only for pointer receivers
**Answer: b) Yes, the outer type's own method takes precedence (shadowing), though the inner one can still be called explicitly via the field name**

**5.15** What's the convention for whether all methods on a type should use consistent receiver types (all value or all pointer)?
a) No convention, mix freely  b) Idiomatic Go generally keeps receiver type consistent across a type's methods to avoid confusion  c) Must always use pointer  d) Must always use value
**Answer: b) Consistency is the idiomatic Go convention, though not enforced by the compiler**

**5.16** What does `type ID int` (a defined type) let you do that a plain `int` doesn't?
a) Nothing different  b) Attach methods to it, gaining type safety distinct from raw int  c) Store it in maps only  d) Nothing, syntax error
**Answer: b) Named/defined types can have their own methods and provide type-safety (e.g., `UserID` vs `ProductID` won't mix accidentally)**

**5.17** Predict output:
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

**5.18** What's the difference between `Dog{Animal{"Rex"}, "Labrador"}` and `Dog{Animal: Animal{"Rex"}, Breed: "Labrador"}`?
a) Different behavior  b) Same result, positional vs named-field initialization  c) First is invalid  d) Second is invalid
**Answer: b) Functionally identical — named fields are just clearer/safer against field-order changes**

**5.19** Can methods be defined on non-struct named types, e.g. `type Celsius float64`?
a) No, only structs  b) Yes, any named type can have methods  c) Only if it embeds a struct  d) Only pointer types
**Answer: b) Yes — Go allows methods on any user-defined (named) type, not just structs**

**5.20** What happens if you call a pointer-receiver method on a value that is not addressable (e.g., a map value)?
a) Works fine always  b) Compile error — map values aren't addressable, so pointer methods can't be called directly on them  c) Panic  d) Auto-converts
**Answer: b) Compile error — you can't take the address of a map element, so pointer-receiver methods fail there directly (would need to extract, modify, reassign)**

---

### Level 5 — Coding Problems

**5.21** Define a `Server` struct with `Name string`, `Port int`, `Healthy bool`, and a pointer-receiver method `MarkUnhealthy()` that sets `Healthy = false`.
**Answer:**
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

**5.22** Use struct embedding to build a `BaseHandler` with a `Log(msg string)` method, embedded into an `AuthHandler`, and call `Log` via the outer struct.
**Answer:**
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
	h.Log("authenticating request")
}
```

**5.23** Define a `Config` struct with JSON tags for fields `APIKey string` (json: "api_key") and `MaxRetries int` (json: "max_retries"), then marshal it to JSON.
**Answer:**
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

**5.24** Write a `HealthCheck` struct with a method `Status() string` that returns "healthy" or "unhealthy" based on an internal `errCount int`, using a threshold constant.
**Answer:**
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

**5.25** Define a `Duration` type based on `int` (representing seconds) with a method `Minutes() float64` that converts it — demonstrating methods on non-struct named types.
**Answer:**
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

**6.1** How are interfaces implemented in Go?
a) Explicit `implements` keyword like Java  b) Implicitly — any type with matching methods satisfies the interface automatically  c) Interfaces must be declared inside the struct  d) Not supported in Go
**Answer: b) Implicit satisfaction — no explicit declaration needed**

**6.2** What is the empty interface `interface{}` (or `any` in modern Go)?
a) An interface with no methods, satisfied by every type  b) Invalid syntax  c) Only satisfied by structs  d) Equivalent to `nil`
**Answer: a) Zero-method interface — every type satisfies it, useful for generic-ish containers pre-generics**

**6.3** What is a type assertion, e.g. `v, ok := i.(string)`?
a) Type conversion  b) Extracts the concrete value from an interface if it holds that type, with `ok` indicating success  c) Compile-time only check  d) Not valid syntax
**Answer: b) Runtime check + extraction; the "comma ok" form avoids panicking on mismatch**

**6.4** What happens with `v := i.(string)` (no comma-ok) if `i` doesn't hold a string?
a) Returns zero value  b) Panics at runtime  c) Compile error  d) Returns nil
**Answer: b) Panics — single-value type assertion panics on failure**

**6.5** What is a type switch used for?
a) Same as regular switch  b) Branching logic based on the dynamic/concrete type stored in an interface  c) Not valid Go syntax  d) Only works with structs
**Answer: b) `switch v := i.(type) { case int: ... case string: ... }`**

**6.6** Can a nil concrete type stored in an interface make `interface == nil` false?
a) No, always true if nil  b) Yes — a non-nil interface can wrap a nil pointer, making the interface itself non-nil (classic Go gotcha)  c) Compile error  d) Always panics
**Answer: b) Yes — this is one of Go's most famous gotchas: an interface holding a typed nil is NOT equal to nil**

**6.7** What must a type do to satisfy `io.Writer`?
a) Nothing special  b) Implement `Write(p []byte) (n int, err error)`  c) Embed io.Writer  d) Be a struct only
**Answer: b) Implement exactly that method signature**

**6.8** What's the idiomatic Go interface design philosophy?
a) Large interfaces with many methods, defined upfront  b) Small, focused interfaces (often 1-2 methods), often defined at point of use by the consumer  c) One giant interface per package  d) Interfaces should mirror structs 1:1
**Answer: b) "Accept interfaces, return structs" — small interfaces like io.Reader, io.Writer are the Go idiom**

**6.9** Where should interfaces ideally be defined in Go — producer or consumer package?
a) Always in the producer/implementation package  b) Often in the consumer package that needs the behavior — Go encourages defining interfaces where they're used  c) Doesn't matter  d) Only in a shared "interfaces" package
**Answer: b) Consumer-side interface definition is idiomatic — avoids unnecessary coupling**

**6.10** Predict output:
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

**6.11** What does the `error` interface look like in Go's standard library?
a) `type error interface { Error() string }`  b) A struct  c) A built-in type with no methods  d) `type error interface { Message() string }`
**Answer: a) The error interface requires exactly one method: `Error() string`**

**6.12** Can you assign a `nil` value of a concrete pointer type to an interface variable and get a "non-nil interface"?
a) No  b) Yes, this is the classic "nil interface vs nil concrete type" trap  c) Compile error  d) Only for structs
**Answer: b) Yes — `var p *MyType = nil; var i MyInterface = p` makes `i != nil` even though `p == nil`**

**6.13** What is a "duck typing" language characteristic, and does Go have it (loosely)?
a) N/A  b) "If it walks like a duck and implements the methods, it satisfies the interface" — Go's implicit interfaces resemble this at compile time  c) Go has no such concept  d) Only applies to reflection
**Answer: b) Go's implicit, structural interface satisfaction is often compared to duck typing, but it's still statically type-checked**

**6.14** How many methods can an interface have?
a) Exactly 1  b) Any number, including zero  c) Max 5  d) Must be at least 2
**Answer: b) Zero or more — zero-method interfaces (empty interface) are valid and useful**

**6.15** What does interface embedding look like, e.g. `io.ReadWriter`?
a) Not supported  b) `type ReadWriter interface { Reader; Writer }` — composes multiple interfaces  c) Only structs can embed  d) Requires explicit method redeclaration
**Answer: b) Interfaces can embed other interfaces to compose larger interfaces from small ones**

**6.16** What's the output?
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

**6.17** Why is over-mocking with large interfaces considered an anti-pattern in Go?
a) It's not an anti-pattern  b) Large interfaces are harder to implement/mock and violate the "small interfaces" idiom, increasing coupling  c) Go doesn't support mocking  d) Interfaces can't be used in tests
**Answer: b) Smaller, focused interfaces are easier to implement, test, and mock — Go favors minimal interfaces**

**6.18** What does `var _ Shape = Circle{}` (blank identifier assignment) accomplish?
a) Nothing, invalid  b) A compile-time check that Circle satisfies the Shape interface, without creating a usable variable  c) Runtime assertion  d) Declares Circle as abstract
**Answer: b) A common idiom to statically verify interface satisfaction at compile time**

**6.19** Can a struct satisfy multiple interfaces simultaneously?
a) No, only one interface per type  b) Yes — implicit satisfaction means a type can satisfy any number of interfaces it happens to match  c) Only with embedding  d) Requires explicit declaration per interface
**Answer: b) Yes — as long as the method set matches, a type can satisfy many interfaces at once**

**6.20** What's the practical difference between `interface{}` / `any` and Go generics (1.18+) for writing reusable functions?
a) No difference  b) Empty interface loses type information (needs assertions) at runtime; generics preserve compile-time type safety  c) Generics are slower always  d) Generics replace interfaces entirely
**Answer: b) Generics give compile-time type checking and avoid the runtime type-assertion overhead/unsafety of `interface{}`**

---

### Level 6 — Coding Problems

**6.21** Define a `Notifier` interface with `Send(msg string) error`, and implement it for both `EmailNotifier` and `SlackNotifier` structs.
**Answer:**
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

type SlackNotifier struct{ Channel string }
func (s SlackNotifier) Send(msg string) error {
	fmt.Printf("Posting to #%s: %s\n", s.Channel, msg)
	return nil
}

func alertAll(notifiers []Notifier, msg string) {
	for _, n := range notifiers {
		n.Send(msg)
	}
}

func main() {
	notifiers := []Notifier{
		EmailNotifier{Address: "ops@company.com"},
		SlackNotifier{Channel: "alerts"},
	}
	alertAll(notifiers, "Server down!")
}
```

**6.22** Write a type switch function `describe(i interface{})` that prints different messages for int, string, bool, and a default case.
**Answer:**
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
	describe(3.14)
}
```

**6.23** Implement a custom error type `NotFoundError` satisfying the `error` interface, and demonstrate using it with `errors.As`.
**Answer:**
```go
package main

import (
	"errors"
	"fmt"
)

type NotFoundError struct {
	Resource string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", e.Resource)
}

func findUser(id int) error {
	if id != 1 {
		return &NotFoundError{Resource: fmt.Sprintf("user %d", id)}
	}
	return nil
}

func main() {
	err := findUser(99)
	var nfErr *NotFoundError
	if errors.As(err, &nfErr) {
		fmt.Println("Not found error:", nfErr.Resource)
	}
}
```

**6.24** Demonstrate the "typed nil in interface" gotcha: write code showing a function returning a `*MyError` (nil) assigned to an `error` interface, and check why `err != nil` is true.
**Answer:**
```go
package main

import "fmt"

type MyError struct{ msg string }

func (e *MyError) Error() string { return e.msg }

func doWork(fail bool) *MyError {
	if fail {
		return &MyError{msg: "failed"}
	}
	return nil // typed nil *MyError
}

func main() {
	var err error = doWork(false) // wraps a nil *MyError into a non-nil error interface
	if err != nil {
		fmt.Println("err is non-nil, even though the underlying pointer is nil!")
	} else {
		fmt.Println("err is nil")
	}
}
```

**6.25** Define a small `Cache` interface with `Get(key string) (string, bool)` and `Set(key, value string)`, and implement it with an in-memory `map`-backed struct.
**Answer:**
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

**7.1** What is the `error` type in Go fundamentally?
a) A concrete struct  b) A built-in interface with a single `Error() string` method  c) An alias for string  d) A special keyword
**Answer: b) An interface with `Error() string`**

**7.2** What is the idiomatic way to create a basic error in Go?
a) `new Error("msg")`  b) `errors.New("msg")`  c) `panic("msg")`  d) `err("msg")`
**Answer: b) `errors.New("msg")` or `fmt.Errorf("...")`**

**7.3** How do you check if an error `err` is specifically an `io.EOF` error?
a) `if err == io.EOF` or `if errors.Is(err, io.EOF)`  b) `if err.type == io.EOF`  c) `if err.contains("EOF")`  d) You cannot
**Answer: a) `if errors.Is(err, io.EOF)` is the modern, robust way**

**7.4** What does `fmt.Errorf("failed to open: %w", err)` do?
a) Formats a string only  b) Wraps the original `err` so it can be extracted later via `errors.Unwrap` or `errors.Is`  c) Panics if err is nil  d) Compile error
**Answer: b) The `%w` verb wraps the error, retaining the original error for inspection**

**7.5** When should you use `panic` in a standard Go application?
a) For all errors  b) For database connection failures  c) For truly unrecoverable programming errors (e.g., nil pointer dereference, out of bounds)  d) Instead of returning `error`
**Answer: c) Only for unrecoverable errors/bugs, not for expected runtime failures**

**7.6** What function is used to stop a panic and regain control of the program?
a) `catch()`  b) `recover()`  c) `rescue()`  d) `stopPanic()`
**Answer: b) `recover()`, which must be called inside a deferred function**

**7.7** Predict the output:
```go
err1 := errors.New("error")
err2 := errors.New("error")
fmt.Println(err1 == err2)
```
a) true  b) false  c) Compile error  d) Panic
**Answer: b) false — `errors.New` returns a pointer to a struct, so each call creates a distinct pointer**

**7.8** What is the result of `recover()` if the program is NOT panicking?
a) Panics  b) Returns `nil`  c) Compile error  d) Blocks forever
**Answer: b) Returns `nil`**

**7.9** Which package provides `Is` and `As` functions for error handling?
a) `fmt`  b) `log`  c) `errors`  d) `runtime`
**Answer: c) The `errors` package (introduced in Go 1.13)**

**7.10** If a function returns `(int, error)`, what is the recommended way to name the error return variable if it's named?
a) `e`  b) `error`  c) `err`  d) `Exception`
**Answer: c) `err` is the standard idiomatic name for an error variable**

---

### Level 7 — Coding Problems

**7.11** Write a custom error type `HTTPError` (a struct containing `Code int` and `Message string`) that implements the `error` interface.
**Answer:**
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

**7.12** Write a function that attempts to parse a config string and wraps any resulting error with additional context using `fmt.Errorf` and `%w`.
**Answer:**
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

**8.1** What does the `&` operator do in Go?
a) Logical AND  b) Bitwise AND  c) Takes the memory address of a variable  d) Both b and c
**Answer: d) Bitwise AND, and Address-of operator**

**8.2** What does the `*` operator do when placed before a pointer variable (e.g., `*ptr`)?
a) Multiplies it by zero  b) Dereferences the pointer to access or mutate the underlying value  c) Declares a pointer type  d) Nothing
**Answer: b) Dereferences the pointer**

**8.3** Does Go support pointer arithmetic (e.g., `ptr++`) like C/C++?
a) Yes  b) No, not by default (only via the unsafe package)  c) Only on arrays  d) Yes, but only for ints
**Answer: b) No, pointer arithmetic is not allowed in safe Go**

**8.4** What happens if you dereference a `nil` pointer?
a) Returns zero value  b) Compile error  c) Runtime panic  d) Silently ignored
**Answer: c) Runtime panic (invalid memory address or nil pointer dereference)**

**8.5** Why would you pass a pointer to a struct into a function instead of passing the struct by value?
a) To allow the function to mutate the original struct  b) To avoid copying a large struct (performance)  c) Both a and b  d) You shouldn't
**Answer: c) Both a and b — mutation and performance**

**8.6** What is the output of this code?
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

**8.7** What does `new(int)` return?
a) `int` initialized to 0  b) `*int` pointing to a newly allocated zeroed integer  c) Compile error  d) `nil`
**Answer: b) `*int` pointing to an integer with value 0**

**8.8** Can you take the address of a literal value directly (e.g., `&5`)?
a) Yes  b) No, you can only take the address of an addressable value like a variable  c) Yes, but only for strings  d) Only inside structs
**Answer: b) No, literals are not addressable**

**8.9** What is the output of this code?
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

**8.10** Which built-in function is often used interchangeably with struct initialization to get a pointer, e.g., `&MyStruct{}`?
a) `make`  b) `new`  c) `alloc`  d) `ptr`
**Answer: b) `new(MyStruct)` is equivalent to `&MyStruct{}`**

---

### Level 8 — Coding Problems

**8.11** Write a function `swap(a, b *int)` that swaps the values of two integers using their pointers.
**Answer:**
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

**8.12** Write a program that demonstrates returning a pointer to a local variable from a function (which is safe in Go due to escape analysis).
**Answer:**
```go
package main

import "fmt"

func createCounter() *int {
	count := 10
	return &count // Pe---

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

**INT.1** What is the default behavior of a channel when you send data to it without a receiver?
a) It buffers the data up to 10 elements.
b) It blocks the sending goroutine forever or until a receiver is ready.
c) It returns an error.
d) It drops the data.
**Answer: b) It blocks.** Unbuffered channels block until the other side is ready.

**INT.2** How does Go achieve inheritance?
a) Using the `extends` keyword.
b) Using struct embedding (composition).
c) Go does not support object-oriented programming.
d) Using abstract classes.
**Answer: b) Struct embedding.** Go uses composition over inheritance.

**INT.3** Which of the following statements about `defer` is true?
a) Deferred functions are executed in FIFO order.
b) Deferred functions are executed when the program exits.
c) Deferred functions are executed in LIFO order just before the surrounding function returns.
d) Arguments to deferred functions are evaluated at execution time.
**Answer: c) Executed in LIFO order.** Also note arguments are evaluated at defer-time.

**INT.4** What happens if you run a Go program with a data race?
a) It deterministically panics.
b) It prints a warning but continues.
c) It has undefined behavior unless compiled with `-race` which will panic.
d) Go prevents data races at compile time.
**Answer: c) Undefined behavior.** The race detector (`-race`) helps find them.

**INT.5** What does `sync.WaitGroup` do?
a) Limits the number of goroutines running concurrently.
b) Waits for a collection of goroutines to finish executing.
c) Blocks a channel until data is available.
d) Replaces the need for a Mutex.
**Answer: b) Waits for a collection of goroutines to finish.**

**INT.6** Can you return a pointer to a local variable safely in Go?
a) No, it will cause a segmentation fault.
b) Yes, Go's escape analysis moves the variable to the heap.
c) Yes, but only for primitive types.
d) No, it creates a memory leak.
**Answer: b) Yes, thanks to escape analysis.**

**INT.7** What is the correct way to initialize an empty slice with a specific capacity but zero length?
a) `s := make([]int, 0, 10)`
b) `s := make([]int, 10)`
c) `s := []int{10}`
d) `s := new([]int)`
**Answer: a) `make([]int, 0, 10)`**

**INT.8** What happens when you read from a closed channel?
a) Panic.
b) It blocks forever.
c) It yields the zero value of the channel's type immediately.
d) It returns a compiler error.
**Answer: c) Yields the zero value.** (You can also use the comma-ok idiom to check if it's open).

**INT.9** What is an empty interface (`interface{}`) used for in Go?
a) To represent a lack of data (null).
b) To define a struct with no fields.
c) To hold values of any type.
d) To define functions with no arguments.
**Answer: c) To hold values of any type.** Since every type implements zero methods, every type satisfies the empty interface.

**INT.10** How are map keys compared in Go?
a) Using a hash function provided by the developer.
b) They must implement the `Comparable` interface.
c) Map keys can be of any type, including slices.
d) Map keys must be of a comparable type (e.g., int, string, pointer, struct without slices/maps).
**Answer: d) Keys must be comparable.** You cannot use slices, maps, or functions as map keys.

**INT.11** What is the primary difference between a value receiver and a pointer receiver in a Go method?
a) Value receivers can mutate the original struct; pointer receivers cannot.
b) Pointer receivers avoid copying the struct and allow mutation of the original struct.
c) They are completely interchangeable with no performance or behavioral difference.
d) Value receivers are required for interfaces, pointer receivers are not.
**Answer: b) Pointer receivers allow mutation and avoid copying.** Use pointer receivers if the method needs to mutate the receiver or if the struct is large.

**INT.12** What does the `init()` function do in a Go package?
a) It acts as the main entry point of the application.
b) It is executed once per file when the package is initialized, before `main()`.
c) It must be called manually to initialize variables.
d) It runs concurrently in a separate goroutine.
**Answer: b) Executed automatically during package initialization.** You can have multiple `init()` functions in a single package or even a single file.

**INT.13** How does the `select` statement behave if multiple `case` channels are ready at the same time?
a) It executes the first one listed top-to-bottom.
b) It panics.
c) It chooses one pseudo-randomly.
d) It executes all of them concurrently.
**Answer: c) It chooses one pseudo-randomly.** This prevents starvation of cases further down the list.

**INT.14** When passing a `map` to a function, what is actually being passed?
a) A deep copy of the entire map.
b) A pointer to the map descriptor, meaning modifications inside the function affect the original map.
c) A read-only copy of the map.
d) Maps cannot be passed to functions.
**Answer: b) A pointer to the map descriptor.** Maps (like slices and channels) act as reference types, so mutating a map inside a function mutates the original.

**INT.15** What is the initial stack size of a new goroutine in modern Go (>= 1.4)?
a) 2 KB
b) 8 KB
c) 1 MB
d) 2 MB
**Answer: a) 2 KB.** This incredibly small footprint allows Go programs to spawn hundreds of thousands of goroutines easily. The stack grows and shrinks dynamically as needed.

**INT.16** What is a common way to cause a memory leak in Go using slices?
a) Appending to a slice in a loop.
b) Slicing a small portion of a massive array/slice and keeping it in memory.
c) Passing a slice to a function by value.
d) Using `make` instead of `new`.
**Answer: b) Slicing a small portion of a massive array.** The small slice retains a reference to the *entire* underlying array, preventing the garbage collector from freeing the massive array.

**INT.17** What is the difference between `new(T)` and `make(T)`?
a) `new` allocates memory and returns a pointer; `make` initializes slices, maps, and channels and returns the value itself.
b) `new` is for primitives; `make` is for structs.
c) `new` returns an initialized object; `make` returns a zeroed object.
d) There is no difference; they are aliases.
**Answer: a) `new` allocates zeroed memory and returns a pointer; `make` initializes internal data structures for slices/maps/channels and returns the value.**

**INT.18** How do you explicitly check if an interface value holds a specific underlying type?
a) By using a regular `if` statement like `if val == type`.
b) Using a Type Assertion: `v, ok := val.(SpecificType)`.
c) Using the `typeof()` function.
d) Interfaces cannot be checked at runtime.
**Answer: b) Using a Type Assertion.**

**INT.19** What does `errors.Is(err, targetErr)` do differently from `err == targetErr`?
a) It compares the string values of the errors.
b) It panics if the errors are not equal.
c) It unwraps the error chain to see if `targetErr` exists anywhere in the chain.
d) It casts the error to a struct.
**Answer: c) It unwraps the error chain.** This was introduced in Go 1.13 and is the standard way to check wrapped errors.

**INT.20** Which of the following is true about strings in Go?
a) They are mutable arrays of bytes.
b) They are immutable slices of bytes.
c) They are mutable slices of runes.
d) They are essentially linked lists of characters.
**Answer: b) They are immutable slices of bytes.** Once created, a string's contents cannot be changed.

**INT.21** In the Go scheduler's G-P-M model, what does the 'P' stand for?
a) Process
b) Pointer
c) Logical Processor
d) Program Counter
**Answer: c) Logical Processor.** 'P' represents a logical processor (context). 'M' is an OS thread, and 'G' is a Goroutine.

**INT.22** What is the "blank identifier" (`_`) used for in Go?
a) To define private variables.
b) To discard return values or avoid unused variable errors.
c) To act as a wildcard in regular expressions.
d) To define an untyped constant.
**Answer: b) To discard return values or avoid "unused variable/import" compile errors.**

**INT.23** How can you enforce that a struct `MyStruct` implements an interface `MyInterface` at compile time?
a) By adding `implements MyInterface` to the struct declaration.
b) By assigning a blank identifier: `var _ MyInterface = (*MyStruct)(nil)`.
c) It is impossible to check at compile time.
d) By using the `CheckInterface()` standard library function.
**Answer: b) `var _ MyInterface = (*MyStruct)(nil)`.** This forces the compiler to verify that the pointer to `MyStruct` satisfies the interface.

**INT.24** What happens if you panic inside a Goroutine, but put a `recover()` in the `main` Goroutine?
a) The `main` Goroutine catches the panic and continues safely.
b) The panic cannot cross Goroutine boundaries; the whole program crashes.
c) The panicked Goroutine dies quietly, the rest of the program continues.
d) It creates a deadlock.
**Answer: b) The program crashes.** A `recover()` must be placed in a `defer` block *within the same Goroutine* where the panic occurs.

**INT.25** What is the correct way to perform simple, lock-free counter increments across multiple Goroutines?
a) Using a standard integer `count++`
b) Using `sync.Mutex` (though this introduces locking)
c) Using the `sync/atomic` package: `atomic.AddInt64(&count, 1)`
d) Using a `select` statement.
**Answer: c) Using `sync/atomic`.** The `atomic` package provides low-level, lock-free hardware-level synchronization, which is faster than Mutexes for simple counters.

**INT.26** Which of the following is a key difference between `sync.Map` and a standard `map` wrapped in an `RWMutex`?
a) `sync.Map` is slower for all operations.
b) `sync.Map` is optimized for append-only data or disjoint sets of keys (e.g., caches where keys are written once and read many times).
c) `sync.Map` is strongly typed, standard maps are not.
d) `sync.Map` does not require `make()`.
**Answer: b) Optimized for specific concurrent workloads.** For heavily contended, write-heavy workloads, an RWMutex-wrapped map is often faster.

**INT.27** If you have a buffered channel of size 5, and 3 elements are in it, what happens when you read from it?
a) It blocks until the channel is full.
b) It reads the first element immediately without blocking.
c) It reads all 3 elements at once.
d) It panics.
**Answer: b) It reads immediately.** Buffered channels only block on read when they are completely empty, and block on write when completely full.

**INT.28** What is "variable shadowing" in Go?
a) Changing the type of a variable at runtime.
b) Declaring a new variable with the same name in an inner scope, effectively hiding the outer variable.
c) A security feature to hide memory addresses.
d) A compilation optimization technique.
**Answer: b) Declaring a new variable in an inner scope.** This often happens accidentally with the `:=` operator inside `if` statements or loops.

**INT.29** Can a Go program have a memory leak?
a) No, the Garbage Collector prevents all memory leaks.
b) Yes, primarily through unclosed goroutines (Goroutine leaks) or keeping references to large objects (like sliced arrays).
c) Yes, but only if you use the `unsafe` package.
d) Yes, but the compiler catches them.
**Answer: b) Yes.** Goroutine leaks are one of the most common causes of memory leaks in Go, occurring when a Goroutine is blocked forever waiting on a channel that will never be written to or read from.

**INT.30** What does the Context package (`context.Context`) provide in Go?
a) It replaces global variables for application state.
b) It provides deadlines, cancellation signals, and request-scoped values across API boundaries and goroutines.
c) It is primarily used for database migrations.
d) It manages memory allocation.
**Answer: b) Deadlines, cancellations, and request-scoped values.** It is the idiomatic way to handle request timeouts and cancellations in Go backend services.
