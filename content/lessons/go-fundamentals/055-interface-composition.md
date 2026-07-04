# Interface Composition & Segregation

The **SOLID** principles are the cornerstone of good software design. The 'I' in SOLID stands for **Interface Segregation**: *"No client should be forced to depend on methods it does not use."*

Go's interface system is specifically designed to enforce this principle flawlessly through Interface Composition.

## 1. The Danger of Fat Interfaces

In older languages, it's common to see massive interfaces with dozens of methods.

```go
// ❌ BAD: Fat Interface
type File interface {
    Read(b []byte) (int, error)
    Write(b []byte) (int, error)
    Close() error
    Seek(offset int64) (int64, error)
    Stat() (FileInfo, error)
}
```
If a function only needs to *read* data, forcing it to accept this `File` interface is dangerous. What if someone passes in a mock object that implements `Read` but panics on `Write`? 

## 2. Small Interfaces

In Go, the standard library is built on micro-interfaces, often containing only a **single method**. 

```go
// ✅ GOOD: Segregated Interfaces
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

type Closer interface {
    Close() error
}
```

If your function only needs to read, it accepts an `io.Reader`. This allows it to accept a File, a Network Socket, or an HTTP Request Body, because they all implicitly implement that one tiny `Read` method!

## 3. Interface Composition

What if you *do* need an object that can both Read and Write? 
Instead of building a new fat interface, you **compose** (embed) the small interfaces together.

```go
// Composing Reader and Writer into a new Interface!
type ReadWriter interface {
    Reader
    Writer
}
```
This is identical to struct embedding. The `ReadWriter` interface now legally requires both the `Read` and `Write` methods.

### 🧠 Architecture Insight
By combining small, segregated interfaces using composition, Go codebases remain highly modular. You build exactly the abstraction you need, exactly where you need it, without forcing other packages to implement bloated, unnecessary methods.
