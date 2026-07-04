# Value vs Pointer Receivers (The Golden Rule)

Deciding whether to use a Value Receiver `(u User)` or a Pointer Receiver `(u *User)` is one of the most common dilemmas for new Go developers. 

Here is the definitive guide to making that decision.

## 1. When to use a Pointer Receiver

You should use a pointer receiver if **ANY** of the following are true:

1. **Mutation**: The method needs to modify the receiver. (This is mandatory).
2. **Size**: The struct is very large (e.g., contains hundreds of fields). Copying it by value would impact performance.
3. **Synchronization**: The struct contains synchronization primitives like `sync.Mutex`. (Copying a Mutex by value causes catastrophic deadlocks and is caught by the `go vet` linter).

## 2. When to use a Value Receiver

You should use a value receiver if **ALL** of the following are true:

1. **Immutability**: The method does not need to modify the receiver.
2. **Small Size**: The struct is small (e.g., a simple coordinate `Point{X, Y}`). Small structs are allocated entirely on the ultra-fast Stack. Passing them by pointer forces them to escape to the Heap, causing Garbage Collection lag.
3. **Map Keys**: If you want to use the struct as a key in a `map`, it must be a value.

## 3. The Golden Rule: Consistency

If you are still unsure, default to a **Pointer Receiver**.

However, the most important rule in Go is **Consistency**. 
**Do not mix receiver types for the same struct.**

If a struct has 10 methods, and 1 of them requires a pointer receiver (because it mutates state), then *all 10 methods* should be written with pointer receivers. 

**❌ Bad (Mixed):**
```go
func (u User) GetName() string { ... }
func (u *User) SetName(n string) { ... }
```

**✅ Good (Consistent):**
```go
func (u *User) GetName() string { ... }
func (u *User) SetName(n string) { ... }
```

Why is mixing bad? Because it drastically complicates Interfaces. If you mix receivers, the concrete type `User` won't satisfy an interface, but `*User` will, leading to confusing compile errors when passing the struct around.
