# Types and Interfaces (The Go Foundation)

Unlike Java or C++, Go is **not** a traditional Object-Oriented language. It does not have Classes, it does not have Inheritance (`extends`), and it does not have Constructors.

Go uses a paradigm based entirely on **Composition** and **Implicit Interfaces**.

## 1. Structs (The Data)

In Go, you define the physical layout of memory using a `struct`.

```go
type User struct {
    ID    int
    Name  string
    Email string
}
```

You can attach behaviors (Methods) to these structs using a **Receiver**.

```go
// The (u *User) is the Receiver! It attaches the method to the Struct.
func (u *User) GetEmail() string {
    return u.Email
}
```

## 2. Interfaces (The Behavior)

An interface in Go is simply a list of method signatures. It defines *what* an object can do, but it doesn't care *how* the object does it.

```go
type Notifier interface {
    SendNotification(msg string) error
}
```

### Implicit Satisfaction (Go's Superpower)
In Java, a class must explicitly declare that it implements an interface: `class EmailService implements Notifier`.

In Go, satisfaction is **Implicit**. If a struct happens to have the exact methods defined in the interface, it automatically implements the interface!

```go
type EmailService struct {}

// Because EmailService has this method, it IS a Notifier!
func (e *EmailService) SendNotification(msg string) error {
    fmt.Println("Sending Email:", msg)
    return nil
}

type SMSService struct {}

// SMSService is ALSO a Notifier!
func (s *SMSService) SendNotification(msg string) error {
    fmt.Println("Sending SMS:", msg)
    return nil
}
```

## 3. Polymorphism 

Because both `EmailService` and `SMSService` implicitly satisfy the `Notifier` interface, you can write powerful, decoupled code.

```go
// This function accepts the Interface, not a specific struct!
func AlertAdmin(n Notifier) {
    n.SendNotification("The server is down!")
}

func main() {
    email := &EmailService{}
    sms := &SMSService{}

    // You can pass either one! The function doesn't care!
    AlertAdmin(email)
    AlertAdmin(sms)
}
```
This is the absolute foundation of Mocking and Unit Testing in Go. If you depend on Interfaces instead of concrete structs, you can easily swap out a real Postgres database with a fake Mock Database during a test!

## 4. The Empty Interface (`any`)

What if you want a function to accept *literally anything* (an int, a string, a User struct)?
You use the Empty Interface: `interface{}` (or in modern Go, simply `any`).

```go
func PrintAnything(val any) {
    fmt.Println(val)
}
```

Because an Empty Interface has zero methods required, every single type in the Go language implicitly satisfies it! 
*Warning: Overusing `any` destroys Go's strict type-safety. You should almost never use it unless you are writing generic library code like `fmt.Println` or JSON decoders.*
