# Testing Concurrency

---

# Table of Contents

* Introduction
* Learning Objectives
* Prerequisites
* Why This Topic Exists
* Real-World Analogy
* Core Concepts
* Step-by-Step Implementation
* Testing with the Race Detector
* Using `t.Parallel()`
* Testing Timeouts
* Mocking Channels
* Production Use Cases
* Best Practices
* Common Mistakes
* Debugging Guide
* Exercises
* Quiz
* Interview Questions
* Cheat Sheet
* Summary
* Key Takeaways
* Further Reading
* Next Chapter

---

# Introduction

Testing synchronous code is easy: you put data in, and you verify what comes out. Testing concurrent code is notoriously difficult. Goroutines run in an unpredictable order. If you run a test 10 times, it might pass 9 times and fail once because of a microscopic timing difference in the OS scheduler. This is called a "Flaky Test".

In this chapter, we will learn how to write deterministic tests for concurrent code, how to use Go's race detector in our CI/CD pipelines, and how to verify timeouts and channel behavior reliably.

---

# Learning Objectives

After completing this chapter you will be able to:

* Write tests that wait for Goroutines to finish before asserting results.
* Use `t.Parallel()` to run unit tests concurrently and save time.
* Test code that relies on `time.Sleep` without actually sleeping the test suite.
* Run the race detector `go test -race` to prove thread-safety.

---

# Prerequisites

Before reading this chapter you should know:

* Standard Go Testing (`go test`)
* WaitGroups (`09-WaitGroup.md`)
* Race Conditions (`28-Race-Conditions.md`)

---

# Why This Topic Exists

If you write a function that spawns a Goroutine to update a database, and you write a standard unit test for it, the test will reach the `assert` statement *before* the Goroutine has actually written to the database. The test will fail.

Or worse, the test might pass on your fast MacBook, but fail on the slow GitHub Actions CI runner. You must learn specific patterns to synchronize your tests with your background Goroutines to eliminate these flakes.

---

# Real-World Analogy

### The Food Critic

A food critic (the Test) orders a complex 5-course meal (Concurrent Function) for his table.
* **Bad Test**: The critic orders the food, instantly looks at the empty table, declares the restaurant a failure, and leaves before the chefs have even finished cooking.
* **Good Test**: The critic orders the food, is given a pager (WaitGroup/Channel) by the waiter, and waits. When the pager buzzes, the critic knows all 5 courses have arrived at the table concurrently, and *then* he evaluates the meal.

---

# Core Concepts

* **Synchronization**: Using Channels or WaitGroups inside your tests to pause the `t.Run` block until the background Goroutines are finished.
* **Flaky Tests**: A test that exhibits non-deterministic behavior (sometimes passes, sometimes fails) due to race conditions or timing assumptions.
* **Race Detector**: A built-in Go tool (`-race`) that instruments your code during testing to find unsafe memory accesses.
* **t.Parallel()**: A testing signal that tells `go test` to run this specific test concurrently alongside other parallel tests, drastically speeding up large test suites.

---

# Step-by-Step Implementation (Testing a Background Job)

If a function spawns a Goroutine and doesn't return anything to wait on, it's very hard to test. The best practice is to refactor the function to accept a `done` channel or return a `<-chan struct{}` that signals completion.

1. Refactor the concurrent function to signal when it is finished.
2. In the `_test.go` file, call the function.
3. Use a `select` block in the test to wait for the `done` signal, with a fallback `time.After` timeout to prevent the test from hanging forever if the function fails.
4. After the signal is received, `assert` the results.

---

# Testing with the Race Detector

Always, always, *always* run your tests with the race detector enabled locally and in CI:
```bash
go test -race ./...
```
If your tests pass, but the race detector output shows `WARNING: DATA RACE`, your test is theoretically a failure. The race detector adds CPU and Memory overhead, so it makes tests slower, but it is mandatory for catching hidden bugs in concurrent code.

---

# Using `t.Parallel()`

By default, `go test` runs every test function synchronously, one after the other. If you have 100 tests that each take 1 second, the suite takes 100 seconds. 

By adding `t.Parallel()` to the top of your test functions, Go will run them simultaneously.

```go
func TestMath1(t *testing.T) {
    t.Parallel() // Signal to run concurrently
    // ... test logic
}

func TestMath2(t *testing.T) {
    t.Parallel() 
    // ... test logic
}
```
*Note: If you use `t.Parallel()`, your tests MUST NOT share global variables (like a global database connection or global config), or they will race and fail!*

---

# Beginner Example

Testing a function that spawns a Goroutine using a WaitGroup.

```go
package main

import (
	"sync"
	"testing"
)

// The function we want to test
func ProcessItems(items []int, wg *sync.WaitGroup) []int {
	results := make([]int, len(items))
	for i, item := range items {
		wg.Add(1)
		go func(idx, val int) {
			defer wg.Done()
			results[idx] = val * 2 // Square it
		}(i, item)
	}
	return results
}

// The Test
func TestProcessItems(t *testing.T) {
	var wg sync.WaitGroup
	items := []int{1, 2, 3}
	
	// Call the function
	results := ProcessItems(items, &wg)
	
	// Wait for the background Goroutines to finish!
	wg.Wait()
	
	// Now it is safe to assert
	if results[0] != 2 {
		t.Errorf("Expected 2, got %d", results[0])
	}
	if results[1] != 4 {
		t.Errorf("Expected 4, got %d", results[1])
	}
}
```

---

# Intermediate Example

Testing a channel with a timeout. If the concurrent function deadlocks and never returns data, the test will hang forever. We must protect our tests with a timeout using `select`.

```go
package main

import (
	"testing"
	"time"
)

// Function that does work and sends to a channel
func AsyncWorker() <-chan string {
	out := make(chan string)
	go func() {
		time.Sleep(100 * time.Millisecond) // Simulate work
		out <- "SUCCESS"
	}()
	return out
}

func TestAsyncWorker(t *testing.T) {
	ch := AsyncWorker()

	// Wait for the result, OR a timeout
	select {
	case res := <-ch:
		if res != "SUCCESS" {
			t.Errorf("Expected SUCCESS, got %s", res)
		}
	case <-time.After(2 * time.Second):
		// If 2 seconds pass, the worker is definitely broken.
		// We fail the test explicitly instead of hanging the CI pipeline forever.
		t.Fatal("Test timed out waiting for worker")
	}
}
```

---

# Advanced Example

Mocking the Clock. If a concurrent function has a 10-minute retry loop, you can't have your unit test literally sleep for 10 minutes. You must refactor your code to accept an interface for time, or mock the sleep.

*(A simplified approach using a mocked sleep duration)*

```go
package main

import (
	"testing"
	"time"
)

// The Production Function
// It accepts a 'sleepFunc' so the test can inject a fake one!
func RetryWorker(sleepFunc func(time.Duration)) bool {
	for i := 0; i < 3; i++ {
		// In production, this might be time.Sleep
		sleepFunc(10 * time.Minute) 
	}
	return true
}

func TestRetryWorker(t *testing.T) {
	// Our Mock Sleep function does absolutely nothing!
	// This makes the 30-minute retry loop execute in 0.001 seconds.
	mockSleep := func(d time.Duration) {
		// Do nothing, just return instantly
	}

	result := RetryWorker(mockSleep)
	if !result {
		t.Error("Expected true")
	}
}
```

---

# Production Use Cases

### 1. CI/CD Pipelines
Every professional Go repository has a GitHub Action (or similar) that runs `go test -race ./...`. This is the ultimate gatekeeper that prevents developers from accidentally merging code that contains data races or deadlocks.

### 2. Testing HTTP Handlers
When testing Chi or Gin HTTP handlers that spawn background Goroutines (e.g., to send a welcome email after signup), developers often pass a mock email client with a buffered channel. The test hits the HTTP endpoint, gets a 200 OK, and then reads from the mock email channel to verify the background Goroutine actually executed.

---

# Best Practices

* **Never use `time.Sleep` in tests**: Do not write tests like: `go doWork(); time.Sleep(1 * time.Second); assert()`. It will flake. Always use deterministic synchronization (WaitGroups or Channels).
* **Avoid Globals**: If you use `t.Parallel()`, global variables will cause massive race conditions between your tests.
* **Test Timeouts**: Always wrap channel reads in your tests with a `select` + `time.After` to prevent infinite hangs during deadlocks.

---

# Common Mistakes

### The For-Loop Variable Trap in Tests
Prior to Go 1.22, using `t.Parallel()` inside a table-driven test (`for _, tc := range testCases`) caused all parallel tests to run using the exact same (last) test case because of how closures captured the loop variable.
*Fix*: Always redefine the loop variable inside the loop: `tc := tc`, or upgrade to Go 1.22+.

---

# Debugging Guide

* **"Test hangs forever"**: The test is blocked reading a channel that was never written to or closed. Add a timeout to your tests.
* **"DATA RACE"**: Look closely at the stack trace provided by the race detector. It will tell you exactly which two Goroutines (and which lines of code) accessed the same memory variable simultaneously.

---

# Exercises

## Beginner
Write a function `Sum(a, b int, out chan int)` that calculates the sum in a Goroutine and sends it to the channel. Write a test for this function that reads the channel and asserts the value.

## Intermediate
Write a test for a function that is supposed to close a channel when it finishes. How do you test if a channel is closed? (Hint: use the comma-ok idiom `val, ok := <-ch`).

---

# Quiz

## Multiple Choice Questions
**1. Why should you avoid using `time.Sleep(1 * time.Second)` in a unit test to wait for a Goroutine?**
A) It makes the test too fast.
B) It is non-deterministic (flaky). If the CI server is under heavy load, the Goroutine might take 1.1 seconds, causing the test to fail randomly.
C) The `time` package is not allowed in `_test.go` files.
*Answer*: B

## True or False
**You can safely use `t.Parallel()` on tests that both read and write to a shared global database table without any locks.**
*Answer*: False. Running them in parallel will cause race conditions on the global state, leading to flaky test failures.

---

# Interview Questions

## Beginner
**Q**: How do you prevent a test from hanging forever if a concurrent function deadlocks?
*Answer*: Use a `select` statement that waits on the function's output channel and a `time.After(timeout)` channel. If the timeout triggers first, call `t.Fatal()` to explicitly fail the test.

## Intermediate
**Q**: What is the `-race` flag and how does it work?
*Answer*: The `-race` flag enables the Go Race Detector during testing or building. It instruments memory accesses at compile time and tracks them at runtime to detect if two Goroutines access the same memory address simultaneously without synchronization, where at least one is a write.

## Advanced
**Q**: How do you test a background process that is designed to run in an infinite `for{}` loop?
*Answer*: You must design the infinite loop to be cancellable by accepting a `context.Context`. In your test, you launch the loop, let it run long enough to process a test item, verify the output, and then call `cancel()` on the context. Finally, use a WaitGroup to ensure the Goroutine cleanly exited the infinite loop before the test finishes.

---

# Cheat Sheet

* **Test Timeout Pattern**:
```go
select {
case res := <-ch:
    // assert res
case <-time.After(3 * time.Second):
    t.Fatal("test timed out")
}
```
* **Parallel Testing**:
```go
func TestSomething(t *testing.T) {
    t.Parallel()
    // ...
}
```

---

# Summary

Testing concurrent code forces you to write *better* concurrent code. If a Goroutine is impossible to test, it means its lifecycle is unmanaged and it is likely a source of bugs in production. By designing concurrent functions that explicitly signal when they are done (via Context, Channels, or WaitGroups), you can write deterministic, lightning-fast, parallel test suites.

---

# Key Takeaways

* ✔ Never use `time.Sleep` to synchronize tests.
* ✔ Always run tests with the `-race` flag.
* ✔ Use `select` and `time.After` to prevent hanging tests.
* ✔ Refactor code to be testable by injecting channels/waitgroups.

---

# Further Reading
* [Go Blog: Introducing the Go Race Detector](https://go.dev/blog/race-detector)
* [Testing Context and Goroutines](https://pkg.go.dev/testing)

---

# Next Chapter
➡️ **Next:** `41-Conclusion.md`
