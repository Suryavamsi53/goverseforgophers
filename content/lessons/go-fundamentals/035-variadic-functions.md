# Variadic Functions

A variadic function is a function that can accept a variable (infinite) number of arguments. 
The most famous variadic function in Go is `fmt.Println()`, which is why you can pass it as many strings or numbers as you want.

## 1. Syntax

To make a function variadic, prefix the type of its **final** parameter with an ellipsis `...`.

```go
func sum(nums ...int) int {
    total := 0
    // Inside the function, 'nums' behaves exactly like a slice of ints: []int
    for _, num := range nums {
        total += num
    }
    return total
}

func main() {
    fmt.Println(sum(1, 2))       // Prints 3
    fmt.Println(sum(1, 2, 3, 4)) // Prints 10
    fmt.Println(sum())           // Prints 0
}
```
*Note: A function can only have one variadic parameter, and it MUST be the last parameter in the list.*

## 2. Under the Hood

When you call `sum(1, 2, 3)`, the Go compiler actually creates a hidden backing array, allocates a slice pointing to it, populates it with `[1, 2, 3]`, and passes that slice into the function.

```mermaid
graph LR
    A[Caller: sum 1, 2, 3 ] -->|Compiler Magic| B(Allocates slice: []int{1, 2, 3})
    B --> C[Function Receives Slice]
```

## 3. The Spread Operator (Passing Slices)

What if you already have a slice of data, and you want to pass it into a variadic function?
If you try to pass `[]int{1, 2, 3}` into `sum()`, the compiler will complain because `sum` expects individual integers, not a slice.

To solve this, use the spread operator `...` **after** the slice. This "unpacks" the slice into individual arguments.

```go
func main() {
    mySlice := []int{10, 20, 30}
    
    // ERROR: cannot use mySlice (type []int) as type int in argument to sum
    // total := sum(mySlice) 
    
    // CORRECT: Unpack the slice
    total := sum(mySlice...) 
    fmt.Println(total) // Prints 60
}
```
**Performance Insight**: When you use the `slice...` spread operator, the compiler is smart. It does *not* allocate a new hidden slice; it simply passes your existing slice directly into the variadic function, resulting in zero allocation overhead!
