# Strategy Pattern

The Strategy Pattern is a behavioral design pattern that defines a family of algorithms, encapsulates each one, and makes them interchangeable. It allows the algorithm to vary independently from the clients that use it.

## 1. The Problem

Imagine an E-Commerce system calculating shipping costs. 
If you hardcode the logic using `if/else` or `switch` statements, the code becomes a massive, unmaintainable bottleneck.

```go
// BAD: Violates the Open/Closed Principle!
// Every time a new carrier is added, we have to modify this core function!
func CalculateShipping(carrier string, weight float64) float64 {
    if carrier == "fedex" {
        return weight * 2.5
    } else if carrier == "ups" {
        return weight * 3.0
    } else if carrier == "dhl" {
        return weight * 4.5
    }
    return 0
}
```

## 2. The Solution (The Strategy)

We define an Interface (the Strategy), and create separate structs for each algorithm.

```go
// 1. The Strategy Interface
type ShippingStrategy interface {
    Calculate(weight float64) float64
}

// 2. The Concrete Strategies
type FedEx struct{}
func (f *FedEx) Calculate(weight float64) float64 { return weight * 2.5 }

type UPS struct{}
func (u *UPS) Calculate(weight float64) float64 { return weight * 3.0 }

// 3. The Context (The object that uses the strategy)
type Order struct {
    weight   float64
    // The Order doesn't care HOW shipping is calculated, 
    // it just holds a reference to the Interface!
    strategy ShippingStrategy 
}

func (o *Order) SetStrategy(s ShippingStrategy) {
    o.strategy = s
}

func (o *Order) GetShippingCost() float64 {
    // Delegate the math to the injected strategy!
    return o.strategy.Calculate(o.weight)
}
```

## 3. The Usage (Runtime Swapping)

Because the Context (`Order`) only relies on the Interface, we can swap the algorithm dynamically at runtime without modifying a single line of the `Order` struct!

```go
func main() {
    order := &Order{weight: 10.0}

    // User selects FedEx in the UI
    order.SetStrategy(&FedEx{})
    fmt.Println(order.GetShippingCost()) // 25.0

    // User changes their mind, selects UPS
    order.SetStrategy(&UPS{})
    fmt.Println(order.GetShippingCost()) // 30.0
}
```

## 4. Go-Specific: First-Class Functions

In Java, Strategies require defining an Interface and multiple Classes.
In Go, functions are **First-Class Citizens**. You don't actually need an Interface if the Strategy is a single function! You can just pass a function signature directly.

```go
// The Strategy is just a function signature!
type ShippingStrategy func(weight float64) float64

var FedExStrategy = func(w float64) float64 { return w * 2.5 }
var UPSStrategy = func(w float64) float64 { return w * 3.0 }

func CalculateShipping(weight float64, strategy ShippingStrategy) float64 {
    return strategy(weight)
}

func main() {
    cost := CalculateShipping(10.0, FedExStrategy)
}
```
This is drastically simpler and more idiomatic in Go, provided the strategy only requires a single method. If the Strategy requires complex state (e.g., maintaining an API Key for FedEx), you should revert to the Interface/Struct approach.
