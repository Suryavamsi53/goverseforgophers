# Structs

Go is not an Object-Oriented language in the traditional sense. It does not have classes, inheritance, or objects. Instead, it relies on **Structs**.

A struct is a typed collection of fields, useful for grouping data together to form records.

## 1. Defining a Struct

You define a struct using the `type` and `struct` keywords. Field names must be unique.

```go
type User struct {
    ID       int
    Name     string
    Email    string
    IsActive bool
}
```
*Visibility Rule: Just like variables and functions, if a field name starts with a Capital letter (e.g., `Name`), it is exported and public to other packages. If it starts with a lowercase letter (e.g., `email`), it is strictly private to the package.*

## 2. Instantiating a Struct

There are three common ways to create a struct in memory:

```go
// 1. Zero Value Initialization
// Every field is automatically set to its zero value (0, "", false)
var u1 User

// 2. Struct Literal (Most Common)
// You explicitly define the fields. Unspecified fields default to zero.
u2 := User{
    Name:  "Alice",
    Email: "alice@example.com",
}

// 3. The `new` Keyword (Returns a Pointer)
// Allocates memory and returns *User
u3 := new(User)
u3.Name = "Bob"
```

## 3. Pointers to Structs

Unlike C/C++, where you must use the `->` operator to access fields of a struct pointer, Go automatically dereferences struct pointers for you using the standard `.` operator.

```go
u := &User{Name: "Charlie"}

// In C, you would write: (*u).Name or u->Name
// In Go, it's just:
fmt.Println(u.Name) 
```

## 4. Struct Equality

Structs are comparable by default. If you use the `==` operator on two structs, Go will compare them field-by-field.

```go
type Point struct {
    X, Y int
}

p1 := Point{X: 1, Y: 2}
p2 := Point{X: 1, Y: 2}

fmt.Println(p1 == p2) // true
```
*Warning: If a struct contains a field that is NOT comparable (like a Slice, Map, or Function), the compiler will throw an error if you try to use `==`. You must use `reflect.DeepEqual` in those cases.*
