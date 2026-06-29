# State Pattern

---

# Table of Contents

* Introduction
* Learning Objectives
* Prerequisites
* Why This Topic Exists
* Real-World Analogy
* Core Concepts
* Architecture Diagram
* Step-by-Step Implementation
* Syntax
* Beginner Example
* Intermediate Example
* Advanced Example
* Production Use Cases
* Performance Analysis
* Best Practices
* Common Mistakes
* Debugging Guide
* Exercises
* Quiz
* Interview Questions
* Cheat Sheet
* Summary
* Key Takeaways
* Further Reading
* Next Chapter

---

# Introduction

The **State Pattern** is a Behavioral Design Pattern that allows an object to alter its behavior when its internal state changes. It appears as if the object changed its class.

It is closely related to the concept of a Finite State Machine (FSM). Instead of writing massive `switch(currentState)` statements inside every method of your object, you extract the state-specific behaviors into separate structs (States). The main object (Context) simply holds a pointer to its current State struct and delegates all work to it.

---

# Learning Objectives

After completing this chapter you will be able to:

* Identify when to replace large conditional state-checks with the State pattern.
* Implement a Finite State Machine using Go interfaces.
* Manage safe transitions between different states.

---

# Prerequisites

Before reading this chapter you should know:

* Structs and Interfaces.
* The Strategy Pattern (`13-Strategy.md`) - The State pattern is structurally identical, but with a different intent.

---

# Why This Topic Exists

Imagine you are programming a Vending Machine. It has methods like `InsertCoin()`, `PressButton()`, and `DispenseItem()`.
Depending on the machine's state (HasCoin, NoCoin, OutOfStock), these methods behave differently.

If you write it normally, your `InsertCoin()` method looks like this:
```go
func (v *VendingMachine) InsertCoin() {
    if v.state == "HasCoin" { fmt.Println("Already has coin") }
    else if v.state == "OutOfStock" { fmt.Println("Cannot accept, empty") }
    else if v.state == "NoCoin" { v.state = "HasCoin" }
}
```
You have to write this massive `if/else` block inside *every single method*. When you add a new state ("MaintenanceMode"), you have to edit 10 different methods. This is a nightmare to maintain.

With the State pattern, you extract `HasCoinState`, `NoCoinState`, and `OutOfStockState` into distinct structs. The Vending Machine delegates `InsertCoin()` to the active state struct, eliminating all `if/else` statements.

---

# Real-World Analogy

### The Smartphone

* **The Object**: Your Smartphone.
* **The Buttons**: Volume Up, Volume Down, Power Button.
* **The States**: 
  1. *Screen Locked State*: Pressing Volume does nothing. Pressing Power turns on the screen.
  2. *Unlocked State*: Pressing Volume changes ringer volume. Pressing Power locks the screen.
  3. *In-Call State*: Pressing Volume changes the earpiece volume. Pressing Power hangs up the call.

The buttons (methods) remain exactly the same, but the smartphone's behavior changes completely depending on its current State.

---

# Core Concepts

* **State Interface**: Defines the methods that represent the actions the Context can perform.
* **Concrete States**: Structs that implement the interface, providing state-specific behavior.
* **Context**: The main object the client interacts with. It maintains a reference to a Concrete State object representing its current state.
* **State Transition**: The act of swapping the Context's current state pointer to a different Concrete State struct.

---

# Architecture Diagram

```mermaid
flowchart TD
    Client[Client Code]
    Context[VendingMachine<br/>Holds 'currentState']
    Interface((State Interface<br/>InsertCoin()<br/>Dispense()))
    
    NoCoin[NoCoinState]
    HasCoin[HasCoinState]

    Client -- "Calls InsertCoin()" --> Context
    Context -- "Delegates to" --> Interface
    
    NoCoin -. implements .-> Interface
    HasCoin -. implements .-> Interface
    
    NoCoin -- "Transitions Context to" --> HasCoin
```

---

# Step-by-Step Implementation

1. Define the **State Interface** with all the context's action methods.
2. Create the **Context** struct. It must hold a field `currentState StateInterface`.
3. Create methods on the Context that map exactly to the Interface. Inside them, delegate the work: `c.currentState.Action()`.
4. Create **Concrete State** structs. They should usually hold a pointer back to the Context so they can trigger state transitions.
5. Implement the interface methods on the Concrete States. When a state finishes its logic, it explicitly changes the Context's state: `c.vendingMachine.SetState(c.vendingMachine.hasCoinState)`.

---

# Syntax

```go
type State interface { Action() }

type Context struct {
    currentState State
    stateA       State
    stateB       State
}
func (c *Context) SetState(s State) { c.currentState = s }
func (c *Context) Action() { c.currentState.Action() }

// Concrete State
type StateA struct { ctx *Context }
func (s *StateA) Action() { 
    fmt.Println("Doing A")
    s.ctx.SetState(s.ctx.stateB) // Transition!
}
```

---

# Beginner Example

A simple Document approval workflow (Draft -> Review -> Published).

```go
package main

import "fmt"

// 1. State Interface
type State interface {
	Publish()
}

// 2. Context
type Document struct {
	draftState   State
	reviewState  State
	publishState State
	currentState State
}

func NewDocument() *Document {
	doc := &Document{}
	// Initialize states with back-references to the context
	doc.draftState = &DraftState{doc: doc}
	doc.reviewState = &ReviewState{doc: doc}
	doc.publishState = &PublishedState{doc: doc}
	
	// Initial state
	doc.currentState = doc.draftState
	return doc
}

func (d *Document) SetState(s State) { d.currentState = s }
func (d *Document) Publish()         { d.currentState.Publish() } // Delegate


// 3. Concrete States
type DraftState struct{ doc *Document }
func (s *DraftState) Publish() {
	fmt.Println("Draft sent for Review.")
	s.doc.SetState(s.doc.reviewState) // Transition
}

type ReviewState struct{ doc *Document }
func (s *ReviewState) Publish() {
	fmt.Println("Review approved. Document Published!")
	s.doc.SetState(s.doc.publishState) // Transition
}

type PublishedState struct{ doc *Document }
func (s *PublishedState) Publish() {
	fmt.Println("Error: Document is already published.")
	// No transition
}


func main() {
	doc := NewDocument()

	doc.Publish() // Draft -> Review
	doc.Publish() // Review -> Published
	doc.Publish() // Already published (Error)
}
```

---

# Intermediate Example

The Vending Machine (Handling multiple actions).

```go
package main

import "fmt"

type State interface {
	InsertCoin()
	Dispense()
}

type VendingMachine struct {
	hasCoin State
	noCoin  State
	current State
}

func NewMachine() *VendingMachine {
	m := &VendingMachine{}
	m.hasCoin = &HasCoinState{machine: m}
	m.noCoin = &NoCoinState{machine: m}
	m.current = m.noCoin
	return m
}

func (m *VendingMachine) InsertCoin() { m.current.InsertCoin() }
func (m *VendingMachine) Dispense()   { m.current.Dispense() }


// --- NO COIN STATE ---
type NoCoinState struct{ machine *VendingMachine }

func (s *NoCoinState) InsertCoin() {
	fmt.Println("Coin accepted.")
	s.machine.current = s.machine.hasCoin // Transition
}
func (s *NoCoinState) Dispense() {
	fmt.Println("Error: Insert coin first.")
}

// --- HAS COIN STATE ---
type HasCoinState struct{ machine *VendingMachine }

func (s *HasCoinState) InsertCoin() {
	fmt.Println("Error: Already has coin.")
}
func (s *HasCoinState) Dispense() {
	fmt.Println("Dispensing item...")
	s.machine.current = s.machine.noCoin // Transition back
}


func main() {
	machine := NewMachine()

	machine.Dispense()   // Error
	machine.InsertCoin() // Accepted
	machine.InsertCoin() // Error
	machine.Dispense()   // Dispensed!
	machine.Dispense()   // Error
}
```

---

# Strategy vs State Pattern

Structurally, these two patterns look identical (A Context delegates to an Interface). 
The difference is **Intent and Transitions**:
* **Strategy Pattern**: The Client (the code in `main`) creates the strategy and manually injects it into the context. The strategies know nothing about each other and rarely transition between each other.
* **State Pattern**: The Client interacts with the Context, entirely blind to the states. The States contain the transition logic (`c.SetState(c.stateB)`) and actively swap *themselves* out based on the rules of the Finite State Machine.

---

# Production Use Cases

### 1. TCP Connection Management
A TCP socket is a classic Finite State Machine (LISTEN, SYN-SENT, SYN-RECEIVED, ESTABLISHED, FIN-WAIT). Network programming libraries often use the State pattern to handle packets differently depending on whether the connection is currently handshaking or fully established.

### 2. Lexers and Parsers
When writing a compiler or a JSON parser, the parser reads characters one by one. If it reads a `"`, it enters the `StringState`. In this state, spaces are treated as literal characters. If it reads another `"`, it transitions back to `CodeState`, where spaces are ignored.

---

# Performance Analysis

The State pattern removes massive `switch` statements, replacing them with interface method calls. Interface method calls in Go have a tiny dynamic dispatch overhead, but this is negligible. The memory overhead is also tiny; you usually instantiate all Concrete State structs once (when the Context is created) and just swap the pointers around.

---

# Best Practices

* **Store States in Context**: To avoid allocating new Memory every time a transition happens, initialize all possible Concrete State structs inside the Context's constructor, and just pass pointers to them during transitions.
* **Define Default Behaviors**: If you have 10 states and 10 methods, writing 100 empty or error-returning methods is tedious. You can create an embedded "BaseState" struct that returns errors for all methods, and have your Concrete States embed it, overriding only the methods they care about.

---

# Common Mistakes

### Transition Spaghetti
If every State can transition to every other State, the logic becomes an unreadable spiderweb. Ensure you map out your Finite State Machine on a piece of paper before coding it. States should only transition through logical, predefined paths.

---

# Debugging Guide

* **"Infinite Transition Loops"**: State A transitions to State B. State B's initialization code immediately calls a method that transitions back to State A. This causes a stack overflow. Ensure transitions only happen as a direct result of an explicit action.

---

# Exercises

## Beginner
Create an `AudioPlayer` context. Create a `State` interface with `ClickPlay()` and `ClickLock()`. Implement three states: `LockedState`, `PlayingState`, and `PausedState`. Ensure clicking play while in `LockedState` does nothing.

## Intermediate
Implement the "BaseState" embedded struct approach. Create a BaseState that prints "Action not allowed in this state" for all interface methods. Embed it into a Concrete State to save yourself from writing boilerplate methods.

---

# Quiz

## Multiple Choice Questions
**1. Who is primarily responsible for triggering the transition to a new State in the State pattern?**
A) The Client (main function).
B) The Context.
C) The currently active Concrete State itself.
*Answer*: C. While the Context *can* do it, the true power of the State pattern is that the active State evaluates the input and decides when it is time to transition the Context to the next State.

## True or False
**The State pattern is identical to the Strategy pattern in both structure and architectural intent.**
*Answer*: False. The structure is identical, but the intent is opposite. Strategy encapsulates independent algorithms that are swapped by the Client. State encapsulates dependent states that swap themselves based on FSM rules.

---

# Interview Questions

## Beginner
**Q**: What code smell usually indicates that you need the State Pattern?
*Answer*: Massive `switch(currentState)` or `if/else` blocks duplicated across multiple methods inside a single struct.

## Intermediate
**Q**: Why is it recommended for the Context to instantiate and hold pointers to all its possible States?
*Answer*: If States are instantiated dynamically inside the transition methods (e.g., `SetState(&NewState{})`), you will generate garbage collection overhead every time the object changes state. By instantiating them once and just swapping pointers, the transitions become extremely fast and zero-allocation.

## Advanced
**Q**: Explain how a Lexer (like the one used in the Go `text/template` package) utilizes a functional variation of the State Pattern.
*Answer*: Instead of using Interfaces and Structs, Rob Pike famously implemented a Lexer using First-Class Functions. A `stateFn` is defined as a function that returns another `stateFn`. The context just runs a loop: `for state != nil { state = state(context) }`. Each state function reads characters, performs its logic, and returns the next state function to execute, creating a brilliant, allocation-free FSM.

---

# Cheat Sheet

* **Interface**: `type State interface { Do() }`
* **Context**:
```go
type Context struct {
    StateA State
    StateB State
    Active State
}
func (c *Context) Do() { c.Active.Do() }
```
* **Concrete State (Transitions)**:
```go
type StateA struct { c *Context }
func (s *StateA) Do() { 
    s.c.Active = s.c.StateB // Transition!
}
```

---

# Summary

The State pattern transforms tangled, conditional FSM spaghetti into clean, isolated, and highly readable components. By giving each state its own struct, the rules for transitioning become explicit and bugs regarding invalid state actions are practically eliminated at compile-time.

---

# Key Takeaways

* ✔ Replaces large `switch` statements with polymorphism.
* ✔ Structurally identical to Strategy, but used for Finite State Machines.
* ✔ Concrete states usually hold a pointer to the Context to trigger transitions.
* ✔ Pre-allocate state structs in the Context to avoid GC pressure.

---

# Further Reading
* [Refactoring.guru: State Pattern](https://refactoring.guru/design-patterns/state)
* [Lexical Scanning in Go (Rob Pike's Functional State Pattern)](https://www.youtube.com/watch?v=HxaD_trXwRE)

---

# Next Chapter
➡️ **Next:** `17-Chain-of-Responsibility.md`
