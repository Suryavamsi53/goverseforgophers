# Implicit Implementation & Architecture

We saw that Go uses **Implicit Satisfaction** for interfaces. But how does this change the way we architect software?

## 1. The Decoupling Power

In Java or C#, the class defining the behavior must explicitly know about the interface it is implementing:

```java
// Java: The Dog class MUST import and know about the Animal interface
import com.example.Animal;

public class Dog implements Animal { ... }
```

In Go, the struct and the interface can be completely blind to each other. They can exist in entirely separate packages, written by different developers, and they will still work together.

```go
// Package A (Written by Google)
type Dog struct {}
func (d Dog) Speak() string { return "Woof" }

// Package B (Written by You)
type Speaker interface {
    Speak() string
}

// You can use Google's Dog as a Speaker without modifying Google's code!
```

## 2. Postel's Law in Go

A famous architectural guideline in Go is based on Postel's Law (The Robustness Principle):

> **"Accept Interfaces, Return Structs."**

### ❌ Bad Architecture:
```go
// Demands a specific concrete type. Hard to test!
func SaveUser(db *MySQLDatabase, u User) { ... }

// Returns an interface. Forces the caller to use type assertions to get underlying data.
func GetUser() DatabaseEntity { ... }
```

### ✅ Good Architecture:
```go
// Accepts an interface. You can pass a MySQLDatabase, or a MockDatabase!
func SaveUser(db DataSaver, u User) { ... }

// Returns a concrete struct. The caller knows exactly what they are getting.
func GetUser() *User { ... }
```

By returning structs, you don't force consumers of your package to deal with interface abstraction unless they want to. By accepting interfaces, you make your functions incredibly flexible and easy to mock during Unit Testing.
