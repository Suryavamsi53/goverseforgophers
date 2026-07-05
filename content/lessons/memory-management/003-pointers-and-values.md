# Pointers vs Values (Semantics)

One of the most debated topics in Go is when to use a Pointer (`*User`) and when to use a Value (`User`).

If you make the wrong choice, you will either create massive memory leaks, destroy CPU performance, or introduce terrifying data mutation bugs.

## 1. Value Semantics (Safety and Speed)

When you pass a variable by Value, Go creates a mathematically perfect, isolated **copy** of the data.

```go
func UpdateAge(u User) {
    u.Age = 99 // This only modifies the COPY!
}

func main() {
    u := User{Name: "Bob", Age: 30}
    UpdateAge(u)
    // u.Age is STILL 30!
}
```

* **Pros**: 100% thread-safe. You can pass a Value across 50 Goroutines concurrently without using Mutex Locks, because every Goroutine gets its own private copy! 
* **Pros**: Stays on the Stack! (No Garbage Collection overhead!).

## 2. Pointer Semantics (Mutation and Sharing)

When you pass a variable by Pointer, you are passing the memory address. Multiple functions now share the *exact same* piece of physical memory.

```go
func UpdateAge(u *User) {
    u.Age = 99 // This modifies the ORIGINAL struct!
}

func main() {
    u := User{Name: "Bob", Age: 30}
    UpdateAge(&u)
    // u.Age is now 99!
}
```

* **Pros**: Allows in-place mutation.
* **Cons**: If 2 Goroutines write to this Pointer at the same time, your app will instantly crash with a **Data Race** panic!
* **Cons**: Forces the variable onto the Heap (Garbage Collection overhead).

## 3. The 3 Enterprise Rules

How do you decide which to use? Follow these 3 rules strictly.

### Rule 1: Built-in Types (Value)
Never use pointers for `int`, `string`, `bool`, or `float64`. (Unless you specifically need a `nil` representation for a JSON/SQL field, in which case you should prefer `sql.NullString` or a wrapper type).

### Rule 2: Reference Types (Value)
Never use pointers for `slices`, `maps`, or `channels`.
These 3 types are already pointers under the hood!
If you pass a `[]int` into a function, the function receives a copy of the slice header, but the header points to the *exact same* underlying array in memory. You do not need `*[]int`.

### Rule 3: Structs (Context Dependent)
For your custom structs (like `User` or `Order`):
* **Use Values** if the struct is small (under 64 bytes) and you do not need to mutate it. Copying small structs on the Stack is exponentially faster than allocating them on the Heap.
* **Use Pointers** if the struct is massive (e.g., a 10KB configuration struct). Copying 10KB of memory on every function call will destroy your CPU cache.
* **Use Pointers** if the struct manages state (like a `sync.Mutex` or a `sql.DB`). You can never copy a Mutex, because the copy will have a different lock state than the original!
