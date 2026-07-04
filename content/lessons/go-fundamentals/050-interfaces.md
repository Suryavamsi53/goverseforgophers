# Interfaces

Interfaces are the crown jewel of the Go programming language. They provide the mechanism for **Polymorphism**—allowing different types to be treated exactly the same way, as long as they exhibit the same behavior.

## 1. What is an Interface?

An interface is simply a list of method signatures. It defines *what* an object should do, but leaves the *how* up to the object itself.

```go
// Any type that has a Speak() string method is legally an "Animal"
type Animal interface {
    Speak() string
}
```

## 2. Implicit Implementation (Duck Typing)

In Java, if a class implements an interface, it must explicitly declare it (`class Dog implements Animal`). 

**In Go, interfaces are satisfied implicitly.** 
There is no `implements` keyword. If a struct happens to have the exact methods defined in the interface, it automatically satisfies the interface! 

*("If it walks like a duck and quacks like a duck, it's a duck.")*

```go
type Dog struct{}

// Because Dog has Speak() string, it is automatically an Animal!
func (d Dog) Speak() string {
    return "Woof!"
}

type Cat struct{}

func (c Cat) Speak() string {
    return "Meow!"
}
```

## 3. Polymorphism in Action

Because both `Dog` and `Cat` satisfy the `Animal` interface, we can write a single function that accepts any `Animal`.

```go
func MakeSound(a Animal) {
    // The function doesn't know or care if 'a' is a Dog or a Cat.
    // It only cares that it can call Speak()
    fmt.Println(a.Speak())
}

func main() {
    MakeSound(Dog{}) // Outputs: Woof!
    MakeSound(Cat{}) // Outputs: Meow!
}
```

## 4. Under the Hood: The `iface` Struct

How does the runtime know what the underlying type is when you pass a `Dog` into an `Animal` variable?

When you assign a concrete value to an interface, Go wraps it in a hidden 16-byte struct called `iface`.

```mermaid
graph LR
    subgraph iface [Interface Struct (16 Bytes)]
        T[itab Pointer]
        D[Data Pointer]
    end
    
    subgraph itab [Type Information Table]
        Type[Type: Dog]
        Func[Method: Dog.Speak()]
    end
    
    subgraph memory [Heap Memory]
        DogStruct[{...Dog Data...}]
    end
    
    T --> itab
    D --> memory
```

1. **`Data Pointer`**: Points to the actual `Dog` struct in memory (causing it to escape to the Heap).
2. **`itab` (Interface Table)**: Points to a table containing the type information (`Dog`) and memory addresses for the methods that satisfy the interface (like `Dog.Speak()`).

Because of this dual-pointer architecture, executing an interface method involves dynamic dispatch (looking up the method address at runtime), which is slightly slower than calling a concrete method directly.
