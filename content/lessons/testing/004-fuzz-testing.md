# Fuzz Testing

In Lesson 1, we learned how to use Table-Driven tests to verify 5 or 10 specific edge cases that we, as humans, could think of. 

But what about the edge cases we didn't think of? What happens if an API user sends a JSON string containing an ancient Egyptian hieroglyph, a 4-gigabyte string of null bytes, or a Unicode emoji that reverses text direction?

Humans cannot write tests for these. We need **Fuzzing**.

## 1. What is Fuzzing?

Fuzz testing is an automated testing technique where a machine rapidly fires random, invalid, or unexpected data inputs into your function, looking to cause a `panic`, a memory leak, or a freeze (infinite loop).

Historically, setting up a fuzzer required complex 3rd-party tools. 
**As of Go 1.18, Fuzzing is built directly into the standard library!**

## 2. Writing a Fuzz Test

Imagine a simple function that reverses a string:
```go
func Reverse(s string) string {
    b := []byte(s)
    for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
        b[i], b[j] = b[j], b[i]
    }
    return string(b)
}
```

To fuzz this, we create a function starting with `Fuzz` (instead of `Test`), and take `*testing.F` (instead of `*testing.T`).

```go
func FuzzReverse(f *testing.F) {
    // 1. Add "Seed Corpus" (Examples to give the fuzzer a starting point)
    f.Add("hello")
    f.Add("world")

    // 2. The Fuzzing Target!
    f.Fuzz(func(t *testing.T, randomString string) {
        
        // We run the function with the randomly generated string
        rev := Reverse(randomString)
        
        // We run it again, it should equal the original string!
        doubleRev := Reverse(rev)
        
        if randomString != doubleRev {
            t.Errorf("Reverse failed! Orig: %q, Rev: %q", randomString, doubleRev)
        }
    })
}
```

## 3. Running the Fuzzer

Unlike standard tests that finish in 0.01 seconds, Fuzz tests run infinitely until they find a crash, or until you manually stop them!

```bash
# Run the fuzzer for exactly 30 seconds
go test -fuzz=FuzzReverse -fuzztime=30s
```

## 4. The Crash (Finding the Bug)

If we run the fuzzer on our `Reverse` function, it will instantly crash!

```text
fuzz: minimizing 53-byte failing input file
--- FAIL: FuzzReverse (0.02s)
    fuzz_test.go:15: Reverse failed! Orig: "Hello, 世界", Rev: "Hello, "
```

**The Bug Detected:** Our `Reverse` function converted the string to a `[]byte`. In Go, strings are UTF-8. The Chinese characters `世界` take 3 bytes each. When we reversed the raw bytes, we shattered the UTF-8 encoding, destroying the string!

The Fuzzer found an edge case that a standard English Table-Driven test never would have caught!

## 5. The Corpus Directory

When the Fuzzer finds a crash, it doesn't just print it to the terminal. It creates a special file inside a new folder: `testdata/fuzz/FuzzReverse/`. 

This file permanently stores the exact string (`"Hello, 世界"`) that caused the crash. 
Now, every time you run a normal `go test` in the future, Go will automatically load that file and use it as a standard Unit Test to ensure you never regress and introduce that specific bug again!
