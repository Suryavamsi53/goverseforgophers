package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

type Course struct {
	ID          string
	Slug        string
	Title       string
	Description string
	Difficulty  string
	Lessons     []string
}

var advancedCourses = []Course{
	{
		ID:          "44444444-4444-4444-4444-444444444444",
		Slug:        "clean-architecture",
		Title:       "Clean Architecture",
		Description: "Learn to build maintainable, scalable, and testable Go applications.",
		Difficulty:  "advanced",
		Lessons: []string{
			"Layered Architecture",
			"Repository Pattern",
			"Service Layer",
			"Dependency Injection",
			"DDD Basics",
		},
	},
	{
		ID:          "55555555-5555-5555-5555-555555555555",
		Slug:        "backend-development",
		Title:       "Backend Development",
		Description: "Master modern backend APIs and real-time communication.",
		Difficulty:  "intermediate",
		Lessons: []string{
			"REST APIs",
			"JWT Authentication",
			"OAuth2",
			"gRPC",
			"GraphQL",
			"WebSockets",
			"Server-Sent Events",
		},
	},
	{
		ID:          "66666666-6666-6666-6666-666666666666",
		Slug:        "database",
		Title:       "Database Engineering",
		Description: "Deep dive into PostgreSQL, transactions, and optimizations.",
		Difficulty:  "advanced",
		Lessons: []string{
			"PostgreSQL",
			"Transactions",
			"Indexes",
			"Connection Pooling",
			"Migrations",
			"Query Optimization",
		},
	},
	{
		ID:          "77777777-7777-7777-7777-777777777777",
		Slug:        "microservices",
		Title:       "Microservices",
		Description: "Design, build, and scale distributed service architectures.",
		Difficulty:  "expert",
		Lessons: []string{
			"Service Discovery",
			"API Gateway",
			"Kafka",
			"RabbitMQ",
			"Redis",
			"Event-Driven Architecture",
			"Saga Pattern",
			"CQRS",
			"Outbox Pattern",
		},
	},
	{
		ID:          "88888888-8888-8888-8888-888888888888",
		Slug:        "cloud-devops",
		Title:       "Cloud & DevOps",
		Description: "Deploy and manage Go applications in the cloud.",
		Difficulty:  "advanced",
		Lessons: []string{
			"Docker",
			"Docker Compose",
			"Kubernetes",
			"Helm",
			"GitHub Actions",
			"Terraform",
			"Cloud Run",
			"GKE",
		},
	},
	{
		ID:          "99999999-9999-9999-9999-999999999999",
		Slug:        "observability",
		Title:       "Observability",
		Description: "Monitor, trace, and debug production Go systems.",
		Difficulty:  "advanced",
		Lessons: []string{
			"Prometheus",
			"Grafana",
			"OpenTelemetry",
			"Structured Logging",
			"Distributed Tracing",
			"Metrics",
		},
	},
	{
		ID:          "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		Slug:        "performance-engineering",
		Title:       "Performance Engineering",
		Description: "Optimize Go applications for high throughput and low latency.",
		Difficulty:  "expert",
		Lessons: []string{
			"pprof",
			"Memory Optimization",
			"Escape Analysis",
			"Scheduler",
			"Garbage Collector",
			"Benchmarking",
		},
	},
	{
		ID:          "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
		Slug:        "system-design",
		Title:       "System Design",
		Description: "Architect massive scale systems for Google-level interviews.",
		Difficulty:  "expert",
		Lessons: []string{
			"URL Shortener",
			"Chat System",
			"Notification Service",
			"Search Engine",
			"Rate Limiter",
			"Distributed Cache",
			"Load Balancer",
			"API Gateway",
		},
	},
}

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
	sqlFile, err := os.Create("scripts/seed_advanced.sql")
	if err != nil {
		panic(err)
	}
	defer sqlFile.Close()

	sqlFile.WriteString("\n-- Seed Advanced Courses\n")
	sqlFile.WriteString("INSERT INTO courses (id, slug, title, description, difficulty) VALUES\n")

	for i, c := range advancedCourses {
		comma := ","
		if i == len(advancedCourses)-1 {
			comma = ";"
		}
		safeTitle := strings.ReplaceAll(c.Title, "'", "''")
		safeDesc := strings.ReplaceAll(c.Description, "'", "''")
		sqlFile.WriteString(fmt.Sprintf("('%s', '%s', '%s', '%s', '%s')%s\n", c.ID, c.Slug, safeTitle, safeDesc, c.Difficulty, comma))
	}
	sqlFile.WriteString("ON CONFLICT DO NOTHING;\n\n")

	sqlFile.WriteString("-- Seed Advanced Lessons\n")
	sqlFile.WriteString("INSERT INTO lessons (id, course_id, slug, title, content, order_index) VALUES\n")

	var totalLessons int
	for _, c := range advancedCourses {
		totalLessons += len(c.Lessons)
	}

	lessonCounter := 0
	for _, c := range advancedCourses {
		courseDir := filepath.Join("content", "lessons", c.Slug)
		os.MkdirAll(courseDir, 0755)

		for j, title := range c.Lessons {
			title = strings.TrimSpace(title)
			if title == "" {
				continue
			}
			
			slugRaw := slugify(title)
			fileName := fmt.Sprintf("%03d-%s.md", j+1, slugRaw)
			slug := fmt.Sprintf("%03d-%s", j+1, slugRaw) 
			filePath := filepath.Join(courseDir, fileName)

			// Create markdown file
			mdContent := fmt.Sprintf(templateMD, title, title, title, title, title, title, title, title, title, title, title, title, title, title)
			os.WriteFile(filePath, []byte(mdContent), 0644)

			// Generate SQL
			id := fmt.Sprintf("20000000-0000-0000-0000-%012d", lessonCounter+1)
			
			comma := ","
			if lessonCounter == totalLessons-1 {
				comma = ";"
			}

			safeTitle := strings.ReplaceAll(title, "'", "''")
			sqlFile.WriteString(fmt.Sprintf("('%s', '%s', '%s', '%s', 'See markdown file', %d)%s\n", 
				id, c.ID, slug, safeTitle, j+1, comma))
			
			lessonCounter++
		}
	}
	sqlFile.WriteString("ON CONFLICT DO NOTHING;\n")

	fmt.Printf("Successfully generated %d advanced markdown lessons\n", totalLessons)
	fmt.Println("Successfully generated scripts/seed_advanced.sql")
}
