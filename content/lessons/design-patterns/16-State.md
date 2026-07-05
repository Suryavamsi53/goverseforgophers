# State Pattern

The State Pattern is a behavioral design pattern that allows an object to alter its behavior when its internal state changes. It appears as if the object changed its class.

It is closely related to the concept of a **Finite State Machine (FSM)**.

## 1. The Problem

Imagine an E-Commerce `Order`. It transitions through statuses: `PENDING` -> `PAID` -> `SHIPPED` -> `DELIVERED`.

If you use a massive `switch` statement in every function, your code becomes unreadable.

```go
func (o *Order) Cancel() error {
    switch o.Status {
    case "PENDING":
        o.Status = "CANCELLED"
        return nil
    case "PAID":
        // Have to issue a refund first!
        Refund(o)
        o.Status = "CANCELLED"
        return nil
    case "SHIPPED":
        return errors.New("cannot cancel, already shipped")
    }
    return nil
}
```
If you add a new "ON_HOLD" status, you have to find and modify every single `switch` statement across the entire codebase!

## 2. The Solution (State Interface)

Instead of a string, the State becomes an Interface. We create a concrete struct for *every* possible state.

```go
// 1. The Context (The Order)
type Order struct {
    State OrderState
}
func (o *Order) SetState(s OrderState) { o.State = s }

// 2. The State Interface
type OrderState interface {
    Pay(order *Order) error
    Cancel(order *Order) error
}
```

## 3. The Concrete States

Now, all the logic for what a "Pending" order can do is isolated into a single, perfectly encapsulated file.

```go
// --- PendingState ---
type PendingState struct{}

func (p *PendingState) Pay(o *Order) error {
    fmt.Println("Payment successful!")
    // Transition to the next state!
    o.SetState(&PaidState{})
    return nil
}

func (p *PendingState) Cancel(o *Order) error {
    fmt.Println("Order cancelled safely.")
    o.SetState(&CancelledState{})
    return nil
}

// --- ShippedState ---
type ShippedState struct{}

func (s *ShippedState) Pay(o *Order) error {
    return errors.New("already paid and shipped")
}

func (s *ShippedState) Cancel(o *Order) error {
    return errors.New("too late to cancel, already shipped")
}
```

## 4. The Usage

The Context (`Order`) delegates all logic to the State interface. 
When `order.Pay()` is called, the `PendingState` executes its logic, and magically updates the Order's internal pointer to `PaidState`. The next time `order.Pay()` is called, the `PaidState` object intercepts it and rejects it!

```go
func main() {
    order := &Order{State: &PendingState{}}
    
    // Delegates to PendingState.Pay()
    // Transitions internal pointer to PaidState
    order.State.Pay(order) 
    
    // Delegates to PaidState.Cancel()
    // PaidState handles the complex refund logic automatically!
    order.State.Cancel(order)
}
```

The State pattern eliminates `switch` statements entirely, guaranteeing that illegal state transitions are mathematically impossible.
