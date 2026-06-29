# Command Pattern

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

The **Command Pattern** is a Behavioral Design Pattern that encapsulates a request as a standalone object containing all information about the request. 

This transformation lets you pass requests as method arguments, delay or queue a request's execution, and support undoable operations. In Go, it is heavily used for building Job Queues (Worker Pools) and CLI (Command Line Interface) applications.

---

# Learning Objectives

After completing this chapter you will be able to:

* Encapsulate method calls into standalone structs.
* Implement a robust Undo/Redo mechanism.
* Decouple the object that *invokes* the operation from the object that *knows how to perform* it.
* Build a simple Job Queue using the Command pattern.

---

# Prerequisites

Before reading this chapter you should know:

* Structs and Interfaces.
* Worker Pools (`34-Worker-Pools.md`).

---

# Why This Topic Exists

Imagine you are building a text editor. You have a UI Button for "Copy". 
If the UI Button directly calls `Editor.CopyText()`, the button is tightly coupled to the Editor. What if you want to trigger "Copy" using a keyboard shortcut (Ctrl+C)? Do you duplicate the code?

Furthermore, how do you implement "Undo" (Ctrl+Z)? 

The Command pattern solves this. The Button doesn't know about the Editor. The Button simply holds a `Command` interface and calls `cmd.Execute()`. The specific `CopyCommand` struct knows how to talk to the Editor. Because commands are objects, you can save them in a history slice. When the user hits Ctrl+Z, you pop the last command off the slice and call `cmd.Undo()`.

---

# Real-World Analogy

### The Restaurant Order

* **The Client**: The Customer. They want a steak.
* **The Command**: The Waiter writes "1 Steak" on a piece of paper (The Order Slip). This piece of paper encapsulates the entire request.
* **The Invoker**: The Waiter takes the slip and puts it on the order rack in the kitchen. The Waiter has no idea how to cook a steak; they just invoke the request.
* **The Receiver**: The Chef picks up the slip from the rack and reads it. The Chef contains the actual business logic to cook the steak.

Because the order is on a piece of paper (an object), it can be queued, delayed, or thrown away if the customer cancels!

---

# Core Concepts

* **Command Interface**: Usually defines a single method: `Execute()`. Sometimes includes `Undo()`.
* **Concrete Command**: A struct that implements the interface. It holds the parameters needed for the action and a reference to the Receiver.
* **Receiver**: The core business logic object that actually performs the work (e.g., The Editor, The Chef).
* **Invoker**: The object that holds the Command and calls `Execute()` when appropriate (e.g., The Button, The Job Queue).

---

# Architecture Diagram

```mermaid
flowchart LR
    Client[Client]
    Invoker[UI Button<br/>Invoker]
    Command((Command Interface<br/>Execute()))
    Concrete[CopyCommand]
    Receiver[TextEditor<br/>Receiver]

    Client -- "1. Creates" --> Concrete
    Client -- "2. Assigns to" --> Invoker
    
    Invoker -- "3. Calls Execute()" --> Command
    Concrete -. "Implements" .-> Command
    Concrete -- "4. Calls CopyText()" --> Receiver
```

---

# Step-by-Step Implementation

1. Declare the **Command Interface** with an `Execute()` method.
2. Create **Receiver** structs that contain the actual heavy business logic.
3. Extract requests into **Concrete Command** structs that implement the interface. Give them fields to store a pointer to the Receiver and any arguments needed.
4. Create the **Invoker** struct (like a Button or a Queue) that accepts a Command via a `SetCommand()` method.
5. The Invoker calls `Execute()` when triggered.

---

# Syntax

```go
type Command interface { Execute() }

type PrintCommand struct {
    receiver *Printer
    text     string
}

func (c *PrintCommand) Execute() {
    c.receiver.PrintText(c.text)
}
```

---

# Beginner Example

A simple Smart Home Remote Control.

```go
package main

import "fmt"

// 1. The Receiver (The actual device)
type Light struct{}
func (l *Light) TurnOn()  { fmt.Println("Light is ON") }
func (l *Light) TurnOff() { fmt.Println("Light is OFF") }

// 2. The Command Interface
type Command interface {
	Execute()
}

// 3. Concrete Commands
type LightOnCommand struct {
	light *Light
}
func (c *LightOnCommand) Execute() { c.light.TurnOn() }

type LightOffCommand struct {
	light *Light
}
func (c *LightOffCommand) Execute() { c.light.TurnOff() }

// 4. The Invoker (The Remote Control Button)
type Button struct {
	command Command
}
func (b *Button) Press() {
	b.command.Execute()
}

func main() {
	// Setup
	livingRoomLight := &Light{}
	
	onCommand := &LightOnCommand{light: livingRoomLight}
	offCommand := &LightOffCommand{light: livingRoomLight}

	// Client configures the invoker
	button := &Button{}

	// Program the button to turn light ON
	button.command = onCommand
	button.Press()

	// Reprogram the button to turn light OFF
	button.command = offCommand
	button.Press()
}
```

---

# Intermediate Example

Implementing an **Undo/Redo History**. This is where the Command pattern becomes irreplaceable.

```go
package main

import "fmt"

// Receiver
type BankAccount struct {
	Balance int
}

// Command Interface
type Command interface {
	Execute()
	Undo()
}

// Concrete Command
type DepositCommand struct {
	account *BankAccount
	amount  int
}
func (c *DepositCommand) Execute() {
	c.account.Balance += c.amount
	fmt.Printf("Deposited $%d. Balance: $%d\n", c.amount, c.account.Balance)
}
func (c *DepositCommand) Undo() {
	c.account.Balance -= c.amount
	fmt.Printf("Undid deposit of $%d. Balance: $%d\n", c.amount, c.account.Balance)
}

// Invoker (The History Manager)
type TransactionManager struct {
	history []Command
}

func (t *TransactionManager) ExecuteCommand(cmd Command) {
	cmd.Execute()
	t.history = append(t.history, cmd) // Save to history!
}

func (t *TransactionManager) UndoLast() {
	if len(t.history) == 0 {
		fmt.Println("Nothing to undo!")
		return
	}
	// Pop the last command
	lastIdx := len(t.history) - 1
	lastCmd := t.history[lastIdx]
	t.history = t.history[:lastIdx]
	
	// Execute the Undo logic
	lastCmd.Undo()
}

func main() {
	account := &BankAccount{Balance: 100}
	manager := &TransactionManager{}

	// Perform transactions
	cmd1 := &DepositCommand{account: account, amount: 50}
	cmd2 := &DepositCommand{account: account, amount: 25}

	manager.ExecuteCommand(cmd1)
	manager.ExecuteCommand(cmd2)

	// Oops, mistake! Undo the last transaction!
	manager.UndoLast()
}
```

---

# Advanced Example

Job Queues and Worker Pools. By turning tasks into Command interfaces, you can push them into a Go Channel. A pool of background Goroutines can pop them off the channel and execute them asynchronously.

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

// Command
type Job interface {
	Execute()
}

// Concrete Job 1
type EmailJob struct { email string }
func (e *EmailJob) Execute() {
	time.Sleep(50 * time.Millisecond)
	fmt.Println("Sent email to:", e.email)
}

// Concrete Job 2
type DBJob struct { query string }
func (d *DBJob) Execute() {
	time.Sleep(50 * time.Millisecond)
	fmt.Println("Executed query:", d.query)
}

// Invoker (The Worker Pool)
type JobQueue struct {
	queue chan Job
	wg    sync.WaitGroup
}

func NewJobQueue(workers int) *JobQueue {
	jq := &JobQueue{
		queue: make(chan Job, 100),
	}
	// Start the background workers
	for i := 0; i < workers; i++ {
		go jq.worker()
	}
	return jq
}

func (jq *JobQueue) worker() {
	for job := range jq.queue {
		job.Execute() // The worker doesn't know what the job is!
		jq.wg.Done()
	}
}

func (jq *JobQueue) AddJob(job Job) {
	jq.wg.Add(1)
	jq.queue <- job
}

func main() {
	// Create a queue with 3 concurrent workers
	jq := NewJobQueue(3)

	// Add various commands to the queue
	jq.AddJob(&EmailJob{email: "alice@example.com"})
	jq.AddJob(&DBJob{query: "UPDATE users SET status=1"})
	jq.AddJob(&EmailJob{email: "bob@example.com"})

	// Wait for all jobs to finish
	jq.wg.Wait()
	fmt.Println("All jobs processed.")
}
```

---

# Production Use Cases

### 1. CLI Applications (e.g., Cobra)
If you've ever used a Go CLI tool built with the popular `cobra` library (like `kubectl` or `hugo`), you've used the Command pattern. Every command (`kubectl get pods`) maps to a `cobra.Command` struct, which encapsulates the arguments and contains a `Run` function (the Execute method).

### 2. Transactional Systems (Saga Pattern)
In distributed microservices, you cannot use a simple Database Transaction. Instead, you create an array of Commands. If Command 3 fails, the system iterates backward through the array, calling `Undo()` on Command 2 and Command 1 to roll back the state. This is called the Saga Pattern.

---

# Performance Analysis

Wrapping method calls into Command structs allocates memory on the heap (especially if pushed into channels or history slices). However, for UI events, CLI execution, or Job Queues, this overhead is entirely negligible.

---

# Best Practices

* **Keep Commands Immutable**: Once a Command is created and configured by the Client, its internal state (arguments) should not be modified. If it is modified, the `Undo()` mechanism will become corrupted.
* **State Snapshotting for Undo**: Sometimes, writing an exact reverse calculation for `Undo()` is impossible (e.g., you can't easily undo a "Delete" operation). Instead, have the Command take a snapshot (backup) of the Receiver's state before calling `Execute()`. Then, `Undo()` simply restores the snapshot.

---

# Common Mistakes

### Bloated Commands
Commands should merely act as a transport layer. Do not put 500 lines of complex database querying inside the `Execute()` method. The Command should just be `c.receiver.ComplexQuery()`, delegating the heavy lifting to the Receiver object.

---

# Debugging Guide

* **"Undo corrupts data"**: The most common issue. Ensure that the arguments captured by the Command during creation were passed by *value*, not by pointer. If you pass a pointer to a struct, and the struct is mutated later, the Command's history is ruined.

---

# Exercises

## Beginner
Create an interface `Task` with an `Execute()` method. Create two structs: `PrintTask` and `BeepTask`. Create an array `[]Task`. Loop over the array and call `Execute()` on each one.

## Intermediate
Create a `Calculator` receiver with a value `Total`. Create an `AddCommand` with an `Execute()` and `Undo()` method. Execute 3 `AddCommand`s, print the total, then loop backward through your history array calling `Undo()`. Verify the total returns to 0.

---

# Quiz

## Multiple Choice Questions
**1. Why is the Command pattern essential for implementing "Undo" functionality?**
A) Because it stores the request as an object in memory, allowing you to keep a history stack and call reverse methods on those specific objects.
B) Because it uses Go channels.
C) Because it forces the Garbage Collector to ignore the objects.
*Answer*: A

## True or False
**The Invoker (e.g., the UI Button) must intimately understand the business logic of the Receiver (the Editor) to execute the command.**
*Answer*: False. The entire point of the pattern is that the Invoker only knows about the `Command Interface`. It just calls `Execute()`. It is completely decoupled from the Receiver.

---

# Interview Questions

## Beginner
**Q**: What is the Command Pattern?
*Answer*: It encapsulates a request into a standalone object. This allows you to parameterize clients with different requests, queue or log requests, and support undoable operations.

## Intermediate
**Q**: How would you implement a simple Job Queue in Go using the Command Pattern?
*Answer*: I would define a `Job` interface with an `Execute()` method. I would create a buffered Go channel of type `chan Job`. I would spin up several Worker goroutines that range over the channel, continuously popping `Job` objects off and calling `job.Execute()`. 

## Advanced
**Q**: When implementing the `Undo()` method for a complex operation where mathematical reversal is impossible, what technique should the Command use?
*Answer*: The Command should use the Memento pattern (State Snapshotting). Inside the Command's `Execute()` method, right before the mutation happens, it should take a deep copy of the Receiver's current state and store it inside the Command struct. The `Undo()` method simply overwrites the Receiver's state with the saved snapshot.

---

# Cheat Sheet

* **Command Interface**: `type Cmd interface { Execute(); Undo() }`
* **Concrete Command**: `type AddCmd struct { rec *Calc; val int }`
* **Invoker (Queue)**: `queue := make(chan Cmd, 10)`
* **Worker**: `for c := range queue { c.Execute() }`

---

# Summary

The Command pattern bridges the gap between intention and execution. By turning "actions" into "nouns" (objects), it unlocks powerful architectural capabilities like delayed execution, worker pools, and historical undos, forming the foundation of many robust CLI frameworks and distributed job processing systems.

---

# Key Takeaways

* ✔ Commands encapsulate requests as objects.
* ✔ Decouples the Invoker (Button) from the Receiver (Business Logic).
* ✔ Enables Undo/Redo by storing Commands in a history slice.
* ✔ Essential for building robust Job Queues via Go channels.

---

# Further Reading
* [Refactoring.guru: Command Pattern](https://refactoring.guru/design-patterns/command)

---

# Next Chapter
➡️ **Next:** `16-State.md`
