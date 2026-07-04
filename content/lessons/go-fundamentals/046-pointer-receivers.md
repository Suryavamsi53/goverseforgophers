# Pointer Receivers

In the previous lesson, we used a **Value Receiver** (`func (c Circle) Area()`). Because Go is pass-by-value, when `myCircle.Area()` was called, the entire `Circle` struct was copied into the method.

What if we want a method to actually **modify** the struct?

## 1. Mutating State

If a method needs to change the data inside the struct, it **must** use a pointer receiver `*T`.

```go
type User struct {
    Name  string
    Score int
}

// ❌ Value Receiver (Modifies a COPY)
func (u User) AddScoreBad() {
    u.Score += 10
}

// ✅ Pointer Receiver (Modifies the ORIGINAL)
func (u *User) AddScoreGood() {
    u.Score += 10
}

func main() {
    player := User{Name: "Alice", Score: 0}
    
    player.AddScoreBad()
    fmt.Println(player.Score) // Still 0!
    
    player.AddScoreGood()
    fmt.Println(player.Score) // 10!
}
```

## 2. Automatic Dereferencing

Did you notice something weird in the code above? 
`player` is a standard value variable, not a pointer. Yet we called `player.AddScoreGood()`, which requires a pointer receiver.

Why didn't we have to write `(&player).AddScoreGood()`?

Because the Go compiler is smart. As a convenience, if you call a pointer method on a value, Go automatically injects the `&` address-of operator for you under the hood.

Likewise, if you have a pointer to a struct, you can call value methods on it without manually dereferencing it `(*player).Area()`.

## 3. Performance Benefits

Just like function parameters, passing a massive struct into a Value Receiver method forces the CPU to copy the entire struct memory on every single method call.

If your struct is large (e.g., it contains large arrays, caches, or database connection pools), you should use a **Pointer Receiver** for *all* methods on that struct, even if the method doesn't mutate the data, simply to prevent massive memory allocation overhead.
