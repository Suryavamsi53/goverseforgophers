# Slices: The Dynamic Window

Because fixed-size Arrays are rigid and expensive to pass around, Go provides **Slices**. Slices are the primary data structure you will use to manage sequences of data. 

A Slice is a dynamically-sized, flexible view into the elements of an underlying Array.

## 1. Creating Slices

A slice type looks exactly like an array type, but without the size specified: `[]T`.

```go
// 1. Creating a slice literal (creates an invisible underlying array)
nums := []int{1, 2, 3, 4, 5}

// 2. Slicing an existing array
arr := [5]string{"a", "b", "c", "d", "e"}
slice := arr[1:4] // Includes indices 1, 2, and 3: ["b", "c", "d"]
```

## 2. The `make` Function

When building applications, you often don't know the exact data at compile time. The `make` function dynamically allocates a slice and its hidden backing array on the heap.

```go
// make(type, length, capacity)
users := make([]string, 0, 100) 
// Creates an empty slice, but pre-allocates memory for 100 strings!
```

## 3. Length vs. Capacity

Every slice has two critical properties that dictate how it behaves in memory:
1. **Length (`len`)**: The number of elements currently accessible in the slice.
2. **Capacity (`cap`)**: The total amount of memory allocated in the underlying array, starting from the first element of the slice.

```mermaid
block-beta
  columns 6
  space:6
  A["Backing Array (Capacity: 6)"]:6
  block:array:6
    B["[0]"]
    C["[1]"]
    D["[2]"]
    E["[3]"]
    F["[4]"]
    G["[5]"]
  end
  space:6
  H["Slice (Length: 3)"]:::sliceBlock
  style H fill:#4CAF50,stroke:#333,stroke-width:2px;
```

```go
arr := [6]int{10, 20, 30, 40, 50, 60}
sl := arr[1:4] // [20, 30, 40]

fmt.Println(len(sl)) // Length: 3
fmt.Println(cap(sl)) // Capacity: 5 (From index 1 to the end of the array)
```

## 4. Pass by Reference (Kind of)

Unlike arrays, passing a slice into a function is incredibly cheap. Why? Because you are not copying the underlying data; you are only copying the tiny slice descriptor struct.

```go
func modify(s []int) {
    s[0] = 99 // This modifies the original backing array!
}

func main() {
    nums := []int{1, 2, 3}
    modify(nums)
    fmt.Println(nums[0]) // Prints 99
}
```
