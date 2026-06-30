package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

var rawCurriculum = `
Introduction to Go
Why Go? History & Philosophy
Installing Go
Go Workspace & Modules
Your First Go Program
Program Structure
Packages & Imports
Variables
Constants
Data Types
Operators
Type Conversion
User Input
Output Formatting
Comments & Documentation
if Statement
switch Statement
for Loop
break
continue
goto (when not to use)
Labels
Functions
Parameters
Multiple Return Values
Named Return Values
Variadic Functions
Anonymous Functions
Closures
Recursion
Arrays
Slices
Slice Internals
append()
copy()
Maps
Map Internals
Strings & Runes
Structs
Anonymous Structs
Nested Structs
Methods
Value vs Pointer Receivers
Composition
Embedding
Pointers
Pointer Receivers
Memory Layout
Escape Analysis Basics
Interfaces
Implicit Implementation
Empty Interface
Type Assertions
Type Switch
Interface Composition
Practical Interfaces
Dependency Injection
Best Practices
Errors
Custom Errors
Wrapping Errors
panic
defer
recover
Introduction to Concurrency
Goroutines
WaitGroup
Channels
Buffered Channels
Unbuffered Channels
Channel Directions
Closing Channels
Range over Channels
Select Statement
Context
Mutex
RWMutex
Atomic Operations
Worker Pool Pattern
Pipeline Pattern
fmt
strings
strconv
bytes
time
os
io
filepath
Reading Files
Writing Files
JSON
CSV
HTTP Client
HTTP Server
REST API Basics
Middleware
Unit Testing
Benchmarking
Profiling
Production Best Practices
Project Structure
Logging
Configuration
Dependency Injection
Graceful Shutdown
Docker
CI/CD
Observability
Performance Optimization
Deployment
`

var templateMD = `# %s

## 1️⃣ Learning Objectives
* **What you'll learn**: Master the core mechanics of %s.
* **Why it matters**: Crucial for building scalable, concurrent, and robust backend systems.
* **Where it's used**: Heavily utilized in API Gateways, Microservices, and High-throughput pipelines.

---

## 2️⃣ Real-world Story
Instead of a dry technical definition, imagine you're managing seats in a cinema... *(To be expanded: A real-world analogy explaining %s)*.

---

## 3️⃣ Visual Learning (Execution Flow & Architecture)
` + "```mermaid\ngraph TD\n    A[Heap Allocation] -->|Garbage Collector| B(Trace Pointers)\n    B --> C{Escape Analysis}\n    C -->|Stack| D[Fast Allocation]\n    C -->|Heap| E[Slower Allocation]\n```" + `

---

## 4️⃣ Internal Working (Under the Hood)
Deep dive into the Go runtime source code.
* **Struct definition**: Exploring ` + "`runtime`" + ` internals.
* **Field by field breakdown**: What does the runtime actually store?

---

## 5️⃣ Compiler Behavior
* **Escape Analysis**: Does this variable escape to the heap?
* **Inlining**: How the compiler optimizes the function call overhead.
* **SSA (Static Single Assignment)**: Optimization passes.

---

## 6️⃣ Memory Management
* **Heap vs Stack**: Memory locality.
* **Garbage Collection**: Impact on GC latency.
* **Pointer Analysis**: Safepoints and write barriers.

---

## 7️⃣ Code Examples

### 🔹 Example 1: Simple
` + "```go\n// Basic implementation\npackage main\n\nfunc main() {\n\t// TODO\n}\n```" + `

### 🔹 Example 2: Intermediate
` + "```go\n// Adding edge cases and error handling\n```" + `

### 🔹 Example 3: Advanced
` + "```go\n// Optimized for zero-allocation\n```" + `

### 🔹 Example 4: Production
` + "```go\n// Production-grade implementation with metrics and context\n```" + `

### 🔹 Example 5: Interview
` + "```go\n// Tricky edge-case testing understanding of pointers/state\n```" + `

---

## 8️⃣ Production Examples
How is %s used in real systems?
1. **Worker Pools**: Distributing tasks.
2. **API Gateways**: Managing request lifecycle.
3. **Kafka Streams**: Batching and dispatching events.

---

## 9️⃣ Performance & Benchmarking
* **CPU vs Memory Trade-offs**
* **Latency impacts**
* **Cache Locality & Branch Prediction**
` + "```bash\ngo test -bench=.\n```" + `

---

## 🔟 Best Practices
* ✅ **Do**: Follow Idiomatic Go patterns.
* ❌ **Don't**: Ignore context cancellation or leak goroutines.
* 🏢 **Google / Uber / Netflix Style**: Explicit error handling, minimal package surface area.

---

## 11️⃣ Common Mistakes
1. **Memory Leaks**: Forgetting to clean up pointers in slices.
2. **Deadlocks**: Improper channel synchronization.
3. **Race Conditions**: Shared state without Mutex.
4. **Shadow Variables**: Accidental re-declaration using ` + "`:=`" + `.

---

## 12️⃣ Debugging
How to troubleshoot %s in production:
* **pprof**: Analyzing heap and CPU profiles.
* **Trace**: Visualizing goroutine execution.
* **Race Detector**: ` + "`go run -race`" + `
* **Delve**: Stepping through memory.

---

## 13️⃣ Exercises
1. **Easy**: Write a basic %s.
2. **Medium**: Refactor to handle concurrent access.
3. **Hard**: Eliminate all heap allocations in the hot path.
4. **Expert**: Implement a custom scheduler utilizing %s.

---

## 14️⃣ Quiz
1. **MCQ**: What happens when you read from a closed %s?
2. **Output Prediction**: What does this program print?
3. **Debugging**: Find the hidden memory leak in this snippet.
4. **Code Review**: Critique this pull request.

---

## 15️⃣ FAANG Interview Questions
* **Beginner**: Explain %s to a junior dev.
* **Intermediate**: How would you optimize %s?
* **Senior (Google/Meta)**: Design a distributed lock manager using %s.
* **System Design Follow-up**: How does this impact your database connection pool?

---

## 16️⃣ Mini Project
**Real-Time %s Implementation**
Build a production-ready feature utilizing %s.
* **Examples**: A concurrent web crawler, an email queue worker, or a reverse proxy.

---

## 17️⃣ Enterprise Features & Observability
* **Logging**: Structured JSON logging.
* **Metrics**: Prometheus instrumentation.
* **Tracing**: OpenTelemetry spans.
* **Security**: Input sanitization.
* **CI/CD & Kubernetes**: Graceful shutdown and liveness probes.

---

## 18️⃣ Source Code Reading
Walkthrough of the Go source code for %s.
* **Why it was implemented this way**.
* **Trade-offs made by the Go core team**.

---

## 19️⃣ Architecture
For production projects integrating %s:
* **Folder Structure**
* **Clean Architecture & DDD**
* **Repository & Service Layers**
* **Testing & Deployment via GitHub Actions**

---

## 20️⃣ Summary & Cheat Sheet
* Key takeaways.
* 1-page quick reference code snippets.
`

func slugify(s string) string {
	s = strings.ToLower(s)
	var builder strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			builder.WriteRune(r)
		} else if r == ' ' || r == '-' {
			builder.WriteRune('-')
		}
	}
	res := builder.String()
	for strings.Contains(res, "--") {
		res = strings.ReplaceAll(res, "--", "-")
	}
	return strings.Trim(res, "-")
}

func main() {
	courseID := "22222222-2222-2222-2222-222222222222"
	lines := strings.Split(strings.TrimSpace(rawCurriculum), "\n")
	
	courseDir := "content/lessons/go-fundamentals"
	os.MkdirAll(courseDir, 0755)

	sqlFile, err := os.Create("scripts/seed_curriculum.sql")
	if err != nil {
		panic(err)
	}
	defer sqlFile.Close()

	sqlFile.WriteString("-- Generated 100-Lesson Curriculum Seed\n")
	sqlFile.WriteString("INSERT INTO lessons (id, course_id, slug, title, content, order_index) VALUES\n")

	for i, title := range lines {
		title = strings.TrimSpace(title)
		if title == "" {
			continue
		}
		
		slugRaw := slugify(title)
		fileName := fmt.Sprintf("%03d-%s.md", i+1, slugRaw)
		slug := fmt.Sprintf("%03d-%s", i+1, slugRaw) 
		filePath := filepath.Join(courseDir, fileName)

		// Create markdown file
		mdContent := fmt.Sprintf(templateMD, title, title, title, title, title, title, title, title, title, title, title, title, title, title)
		os.WriteFile(filePath, []byte(mdContent), 0644)

		// Generate SQL
		id := fmt.Sprintf("10000000-0000-0000-0000-%012d", i+1)
		
		comma := ","
		if i == len(lines)-1 {
			comma = ";"
		}

		safeTitle := strings.ReplaceAll(title, "'", "''")
		sqlFile.WriteString(fmt.Sprintf("('%s', '%s', '%s', '%s', 'See markdown file', %d)%s\n", 
			id, courseID, slug, safeTitle, i+1, comma))
	}

	fmt.Printf("Successfully regenerated %d markdown lessons with Ultimate Template in %s\n", len(lines), courseDir)
}
