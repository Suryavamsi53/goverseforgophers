# Struct Embedding (Promoted Fields)

While Go enforces Composition over Inheritance, writing `robot.Voice.Play()` can sometimes feel verbose. If `RobotDog` is fundamentally a wrapper around a `Speaker`, wouldn't it be nice to just call `robot.Play()` directly?

Go allows this through **Struct Embedding**.

## 1. Anonymous Fields

When defining a struct, you can include another struct without giving it a field name. This is called an embedded (or anonymous) field.

```go
type Logger struct {}

func (l *Logger) Log(msg string) {
    fmt.Println("[LOG]:", msg)
}

// Server embeds Logger
type Server struct {
    Host   string
    Logger // No field name!
}
```

## 2. Promoted Fields and Methods

Because `Logger` is embedded, all of its fields and methods are automatically "promoted" to the `Server` struct. 

You can call `.Log()` directly on the server, as if the server itself implemented it!

```go
func main() {
    s := Server{
        Host: "localhost",
        Logger: Logger{}, 
    }
    
    // We don't have to write s.Logger.Log()
    // The method was promoted directly to the Server!
    s.Log("Server started") 
}
```

## 3. Name Collisions

What happens if the `Server` has its own `.Log()` method, and the embedded `Logger` also has a `.Log()` method?

**The outermost method wins.**

The `Server`'s method will "shadow" the embedded method. If you still need to call the embedded method, you can access it explicitly via the type name: `s.Logger.Log()`.

## 4. Warning: Embedding is NOT Polymorphism

Developers coming from Object-Oriented backgrounds often look at Struct Embedding and think, *"Aha! This is just inheritance in disguise!"*

**It is not.**

If a function expects a `Logger`, you cannot pass a `Server` to it, even if `Server` embeds `Logger`. The Go compiler will strictly reject it because a `Server` *is not* a `Logger`; it simply *has* a `Logger` inside it.

To achieve true polymorphism in Go (where a function can accept multiple different types), you must use **Interfaces**.
