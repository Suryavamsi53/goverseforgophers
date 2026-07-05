# Factory Method Pattern

The Factory Method is a creational design pattern used to create objects without specifying the exact class of object that will be created. 

Because Go does not have Classes, Inheritance, or traditional Constructors, the Factory Method pattern relies heavily on **Interfaces**.

## 1. The Problem

Imagine a Payment processing system. You start with just Stripe.
Later, your manager asks you to add PayPal. Then Crypto. Then Apple Pay.

If you hardcode the instantiation of these structs throughout your codebase, every time you add a new payment provider, you have to modify 50 different files.

## 2. The Solution (Interface Factory)

We define a single interface, and a Factory function that returns different concrete structs depending on the input.

```go
package payment

import "fmt"

// 1. Define the Interface
type PaymentProcessor interface {
    Pay(amount float64) string
}

// 2. Define the Concrete Structs (Unexported)
type stripeProcessor struct{}
func (s *stripeProcessor) Pay(amount float64) string { return "Paid via Stripe" }

type paypalProcessor struct{}
func (p *paypalProcessor) Pay(amount float64) string { return "Paid via PayPal" }

// 3. The Factory Method
func GetPaymentProcessor(method string) (PaymentProcessor, error) {
    switch method {
    case "stripe":
        return &stripeProcessor{}, nil
    case "paypal":
        return &paypalProcessor{}, nil
    default:
        return nil, fmt.Errorf("unknown payment method: %s", method)
    }
}
```

## 3. The Usage

Now, the rest of your application does not need to know that `stripeProcessor` or `paypalProcessor` even exist! They are completely decoupled.

```go
func Checkout(method string, amount float64) {
    // The Factory gives us the correct implementation dynamically
    processor, err := payment.GetPaymentProcessor(method)
    if err != nil {
        log.Fatal(err)
    }

    // We execute the interface method perfectly
    result := processor.Pay(amount)
    fmt.Println(result)
}
```

## 4. Go-Specific Nuances

Wait! In Lesson 2, we said **"Accept Interfaces, Return Structs"**. 
Why is `GetPaymentProcessor` returning an Interface?!

That proverb applies to standard library design and general business logic. The Factory Method is one of the rare exceptions to the rule. 
If the entire purpose of the function is to encapsulate multiple different concrete types behind a single abstraction layer, you *must* return the interface. 

However, in Enterprise Go, true Factory patterns are relatively rare compared to Java. Instead of a runtime string switch (`case "stripe"`), Go developers usually prefer Dependency Injection, passing the concrete struct into the handler at startup.
