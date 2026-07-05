# Testing Concurrency

Testing sequential code is straightforward: you provide an input and `assert.Equal` the output.

Testing concurrent code is a nightmare. Because the Go Scheduler is non-deterministic, a test might pass 99 times and fail once (Flaky Tests). Furthermore, testing time-based logic (like Timers and Tickers) often causes test suites to take minutes to run, ruining the CI/CD pipeline.

## 1. The Data Race Detector in Tests

The absolute most important rule of testing Go code: **Your CI pipeline must run tests with the `-race` flag.**

```bash
go test -race ./...
```
If you do not run tests with the `-race` flag, your tests might pass, but deploy a catastrophic memory corruption bug to production.

## 2. Dealing with Goroutine Leaks

If you write a function that starts a background Goroutine, how do you verify in your test that the Goroutine actually shut down properly and didn't leak?

The standard library provides `runtime.NumGoroutine()`.

```go
func TestWorkerPool(t *testing.T) {
    // 1. Record the baseline number of goroutines
    initialGoroutines := runtime.NumGoroutine()
    
    // 2. Run the concurrent logic
    RunWorkerPool()
    
    // 3. Verify no goroutines leaked!
    // We add a tiny sleep to let the scheduler clean up dead goroutines
    time.Sleep(10 * time.Millisecond)
    finalGoroutines := runtime.NumGoroutine()
    
    if finalGoroutines > initialGoroutines {
        t.Errorf("Goroutine leak detected! Started with %d, ended with %d", 
            initialGoroutines, finalGoroutines)
    }
}
```
*Note: For enterprise codebases, use the open-source library `go.uber.org/goleak` which automates this perfectly.*

## 3. Testing Timeouts (Don't use time.Sleep)

If you write a function that is supposed to timeout after 5 seconds, how do you test it? If you use `time.Sleep(5 * time.Second)` in your test, your entire test suite now takes 5 seconds to run. If you have 100 timeout tests, your test suite takes 8 minutes!

### The Solution: Mocking Time
You should abstract time injection into your structs.

```go
type JobRunner struct {
    // Instead of calling time.After directly, we inject a clock function!
    TimeoutFunc func(time.Duration) <-chan time.Time
}

// In Production:
runner := JobRunner{ TimeoutFunc: time.After }

// In Tests:
func TestJobTimeout(t *testing.T) {
    fakeTimeChan := make(chan time.Time)
    
    runner := JobRunner{
        TimeoutFunc: func(d time.Duration) <-chan time.Time {
            return fakeTimeChan
        },
    }

    // Instantly simulate a 5-second timeout by firing data into the channel!
    // The test completes in 0.001 seconds instead of 5 seconds!
    fakeTimeChan <- time.Now() 
}
```

## 4. The Eventually Pattern

If you are testing a complex distributed system (like a CQRS pipeline), you often need to assert that a database was updated *eventually*.

Do not use hardcoded `time.Sleep(2 * time.Second)` to wait for the data. Use a polling loop. (The open-source library `testify/require` provides this natively).

```go
// Polls the database every 10ms for up to 2 seconds.
// As soon as the condition is met, the test continues instantly!
require.Eventually(t, func() bool {
    user, _ := db.GetUser(1)
    return user.Status == "ACTIVE"
}, 2*time.Second, 10*time.Millisecond)
```
