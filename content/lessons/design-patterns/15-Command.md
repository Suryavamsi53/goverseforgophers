# Command Pattern

The Command Pattern is a behavioral design pattern that turns a request into a stand-alone object. This transformation lets you pass requests as arguments, delay or queue a request's execution, and support undoable operations.

## 1. The Core Idea

Instead of an object calling a method directly (`user.Delete()`), we create an object that *represents* the action (`DeleteUserCommand`). 

This is incredibly powerful because an object can be serialized to JSON, saved to a database, placed in a queue, or executed 3 days later!

## 2. Implementation

```go
// 1. The Command Interface
type Command interface {
    Execute() error
    Undo() error // Optional, but incredibly powerful
}

// 2. The Receiver (The actual business logic)
type BankAccount struct {
    Balance float64
}
func (b *BankAccount) Deposit(amount float64) { b.Balance += amount }
func (b *BankAccount) Withdraw(amount float64) { b.Balance -= amount }

// 3. A Concrete Command Struct
type DepositCommand struct {
    account *BankAccount
    amount  float64
}

// The Command encapsulates the Receiver and the Arguments!
func (c *DepositCommand) Execute() error {
    c.account.Deposit(c.amount)
    return nil
}
func (c *DepositCommand) Undo() error {
    c.account.Withdraw(c.amount) // The exact opposite action!
    return nil
}
```

## 3. The Invoker (Queueing and Undo)

Now we can build a `TransactionManager` that acts as the Invoker. It doesn't know *what* the commands do, it just executes them and keeps a history.

```go
type TransactionManager struct {
    history []Command
}

func (t *TransactionManager) ExecuteCommand(c Command) error {
    err := c.Execute()
    if err == nil {
        t.history = append(t.history, c) // Save to history!
    }
    return err
}

func (t *TransactionManager) UndoLast() {
    if len(t.history) == 0 { return }
    
    // Pop the last command off the stack
    lastIdx := len(t.history) - 1
    lastCommand := t.history[lastIdx]
    t.history = t.history[:lastIdx]
    
    // Call Undo!
    lastCommand.Undo()
}
```

## 4. Real-World Use Cases

1. **Job Queues (Asynchronous Execution)**: You wrap an action in a Command, serialize it to JSON, and push it to RabbitMQ. A background worker pops it off the queue and calls `Execute()`.
2. **Sagas (Distributed Transactions)**: If step 3 of a microservice flow fails, the Orchestrator iterates backward through the Command History and calls `Undo()` on every single step (Compensating Actions).
3. **CQRS (Command Query Responsibility Segregation)**: In Enterprise architecture, Write operations (Commands) are often physically separated from Read operations (Queries) into completely different microservices and databases.
