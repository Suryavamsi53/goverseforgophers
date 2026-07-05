# Accept Interfaces, Return Structs

In Java or C#, it is common practice for a library to define a massive `interface`, implement it, and then return the `interface` to the user from the constructor. 

Go flips this object-oriented paradigm on its head. The most famous Go proverb is:
**"Accept interfaces, return structs."**

## 1. Returning Structs

When you write a package, your constructor should almost always return a concrete struct pointer, not an interface.

```go
package storage

// BAD: Returning an interface
func NewDatabase() DBInterface { ... }

// GOOD: Returning a concrete struct
func NewDatabase() *Database {
    return &Database{}
}
```

### Why?
1. **You cannot predict the future.** If you define `DBInterface` inside your package, you are forcing the user to use *your* definition of what a database is. 
2. **Adding methods breaks code.** If you return an interface, and later add a `Ping()` method to that interface, you instantly break every single implementation of that interface in your user's codebase.

## 2. Accepting Interfaces (The Consumer Defines the Contract)

In Go, interfaces are implemented **implicitly**. The struct does not need to declare `implements MyInterface`.

Because of this, the interface belongs to the **Consumer** (the code calling the function), not the **Producer** (the code implementing the function).

Imagine a `BillingService` that needs to fetch a user from the database.

```go
package billing

// 1. The Consumer defines EXACTLY what it needs.
// It doesn't care about the 50 other methods the Database struct has.
// It only cares about GetUser.
type UserFetcher interface {
    GetUser(id int) (*User, error)
}

// 2. The Consumer accepts the interface
func ChargeUser(fetcher UserFetcher, userID int) error {
    user, _ := fetcher.GetUser(userID)
    // Charge the user...
    return nil
}
```

Now, in `main.go`, we wire them together:

```go
func main() {
    // 1. storage returns a concrete struct
    db := storage.NewDatabase() 
    
    // 2. We pass the concrete struct into Billing.
    // Because db has a GetUser() method, it implicitly satisfies UserFetcher!
    billing.ChargeUser(db, 42) 
}
```

## 3. The Power of Mocking

Because `BillingService` defined its own tiny `UserFetcher` interface, writing Unit Tests is incredibly trivial.

You don't need a massive mocking library. You just define a dummy struct in your test file that implements that one method.

```go
// In billing_test.go
type MockFetcher struct {}
func (m *MockFetcher) GetUser(id int) (*User, error) {
    return &User{Name: "TestUser"}, nil
}

func TestChargeUser(t *testing.T) {
    mock := &MockFetcher{}
    err := ChargeUser(mock, 1)
    // Assert...
}
```

By strictly adhering to "Accept Interfaces, Return Structs", your Go code becomes hyper-modular, perfectly decoupled, and incredibly easy to test.
