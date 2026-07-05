# Facade Pattern

The Facade Pattern is a structural design pattern that provides a simplified, higher-level interface to a complex subsystem. 

Imagine driving a car. You press the accelerator pedal. You do not manually inject fuel into the cylinders, adjust the timing belt, or monitor the oxygen sensors. The accelerator pedal is a **Facade** that hides the terrifying complexity of the internal combustion engine.

## 1. The Problem

In enterprise Go applications, performing a seemingly simple business action might require orchestrating 5 different micro-packages.

```go
// A complex, multi-step process for checking out an E-Commerce Cart
func checkout() {
    // 1. Check inventory
    inv := inventory.NewSubsystem()
    if !inv.CheckStock("item_1") { return }

    // 2. Calculate tax
    tax := taxcalc.NewSubsystem()
    total := tax.Calculate(100.0, "US-CA")

    // 3. Process payment
    pay := payment.NewSubsystem()
    pay.Charge("user_42", total)

    // 4. Send email
    mail := email.NewSubsystem()
    mail.SendReceipt("user_42", total)
}
```

If you force the HTTP Handler (the Controller) to write this orchestration logic, your handlers become massive, bloated, and impossible to test. Furthermore, if you need to perform a checkout from a gRPC endpoint instead of an HTTP endpoint, you have to copy-paste all 4 steps!

## 2. The Solution (The Facade)

We create a single struct (the Facade) that encapsulates the subsystems.

```go
// 1. The Facade Struct
type OrderFacade struct {
    inventory *inventory.Subsystem
    tax       *taxcalc.Subsystem
    payment   *payment.Subsystem
    email     *email.Subsystem
}

// 2. The Constructor (Dependency Injection!)
func NewOrderFacade(i *inventory.Subsystem, t *taxcalc.Subsystem, p *payment.Subsystem, e *email.Subsystem) *OrderFacade {
    return &OrderFacade{
        inventory: i,
        tax:       t,
        payment:   p,
        email:     e,
    }
}

// 3. The Simplified Method
func (f *OrderFacade) CheckoutCart(userID string, itemID string, amount float64) error {
    if !f.inventory.CheckStock(itemID) {
        return errors.New("out of stock")
    }
    total := f.tax.Calculate(amount, "US-CA")
    f.payment.Charge(userID, total)
    f.email.SendReceipt(userID, total)
    return nil
}
```

## 3. The Usage

Now, the HTTP Handler is incredibly clean and completely decoupled from the underlying sub-systems.

```go
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // The handler just pushes the "Accelerator Pedal"!
    err := h.orderFacade.CheckoutCart("user_42", "item_1", 100.0)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    w.Write([]byte("Success!"))
}
```

## 4. Facade vs Adapter vs Decorator

These three structural patterns look very similar, but have entirely different intents:
* **Adapter**: Changes the *interface* of an existing object so it matches what a client expects.
* **Decorator**: Adds *new behavior* to an existing object without changing its interface.
* **Facade**: Simplifies a *complex network* of objects into a single, easy-to-use method.
