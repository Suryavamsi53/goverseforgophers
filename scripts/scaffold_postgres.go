package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

type Lesson struct {
	Title     string
	Subtopics []string
}

type Course struct {
	ID          string
	Slug        string
	Title       string
	Description string
	Difficulty  string
	Lessons     []Lesson
}

var postgresCourse = Course{
	ID:          "c0000000-0000-0000-0000-000000000002",
	Slug:        "postgresql",
	Title:       "PostgreSQL Mastery",
	Description: "From Basic to Advanced to Expert. Learn PostgreSQL internals, real-world backend patterns, and Go integration.",
	Difficulty:  "expert",
	Lessons: []Lesson{
		{"Level 1 - PostgreSQL Fundamentals", []string{"Introduction (Architecture, vs MySQL/MongoDB, ACID, CAP)", "Installation (Docker, pgAdmin, psql)", "Database Basics (Schema, Tables, Tuples, Relations)", "Data Types (Numeric, Character, Boolean, Date/Time, UUID, JSONB, Arrays)"}},
		{"Level 2 - CRUD Operations", []string{"Database Operations", "Table Operations", "Column Operations", "Insert Operations (RETURNING)", "Read Operations", "Update", "Delete"}},
		{"Level 3 - Constraints", []string{"PRIMARY KEY", "FOREIGN KEY", "UNIQUE", "NOT NULL", "CHECK", "DEFAULT", "EXCLUSION Constraint"}},
		{"Level 4 - Querying", []string{"Filtering (WHERE, AND, OR, BETWEEN, IN, LIKE, ILIKE, EXISTS)", "Sorting (ORDER BY, NULLS FIRST/LAST)", "Grouping (GROUP BY, HAVING)"}},
		{"Level 5 - Joins", []string{"INNER JOIN", "LEFT JOIN", "RIGHT JOIN", "FULL JOIN", "CROSS JOIN", "SELF JOIN", "NATURAL JOIN", "LATERAL JOIN"}},
		{"Level 6 - Aggregate Functions", []string{"COUNT()", "SUM()", "AVG()", "MAX()", "MIN()", "STRING_AGG()", "ARRAY_AGG()", "JSON_AGG()"}},
		{"Level 7 - Built-in Functions", []string{"String Functions (CONCAT, SUBSTRING, REPLACE, SPLIT_PART)", "Numeric Functions (ROUND, CEIL, RANDOM)", "Date Functions (NOW(), AGE(), DATE_TRUNC())", "Conditional (CASE, COALESCE, NULLIF)"}},
		{"Level 8 - Views", []string{"Views", "Materialized Views", "Refresh Materialized View"}},
		{"Level 9 - Indexing", []string{"What is an Index?", "Types (B-Tree, Hash, GIN, GiST, BRIN)", "Composite Index", "Partial Index", "Expression Index", "Covering Index", "Reindex"}},
		{"Level 10 - Transactions", []string{"BEGIN, COMMIT, ROLLBACK", "SAVEPOINT", "Isolation Levels (Read Uncommitted, Read Committed, Repeatable Read, Serializable)", "MVCC, Snapshot Isolation"}},
		{"Level 11 - Normalization", []string{"1NF", "2NF", "3NF", "BCNF", "4NF", "5NF", "Denormalization"}},
		{"Level 12 - Advanced SQL", []string{"Subqueries (Scalar, Correlated, Nested)", "Common Table Expressions (WITH, Recursive CTE)", "Window Functions (ROW_NUMBER, RANK, LAG, LEAD)", "Pivot & Unpivot"}},
		{"Level 13 - JSON & JSONB", []string{"Create JSON", "Query JSON", "JSON Operators", "JSON Functions", "JSON Indexing", "JSON Path Queries"}},
		{"Level 14 - Arrays", []string{"Create Arrays", "Update Arrays", "Search Arrays", "Array Functions"}},
		{"Level 15 - UUID", []string{"UUID Generation", "UUID vs SERIAL", "Performance"}},
		{"Level 16 - Sequences", []string{"CREATE SEQUENCE", "NEXTVAL", "CURRVAL", "SETVAL"}},
		{"Level 17 - Stored Procedures", []string{"Functions", "Procedures", "Parameters", "Return Types", "Exception Handling"}},
		{"Level 18 - Triggers", []string{"BEFORE Trigger", "AFTER Trigger", "INSTEAD OF Trigger", "Trigger Functions"}},
		{"Level 19 - PL/pgSQL", []string{"Variables", "Loops", "IF, CASE", "Exceptions", "Dynamic SQL"}},
		{"Level 20 - Security", []string{"Users, Roles, Privileges", "GRANT, REVOKE", "Row Level Security (RLS)", "Authentication & Authorization"}},
		{"Level 21 - Backup & Restore", []string{"pg_dump, pg_restore", "WAL", "Base Backup", "PITR (Point-In-Time Recovery)"}},
		{"Level 22 - Performance Tuning", []string{"EXPLAIN, EXPLAIN ANALYZE", "VACUUM, VACUUM FULL, ANALYZE, AUTOVACUUM", "Statistics", "Query Optimization & Execution Plans"}},
		{"Level 23 - Partitioning", []string{"Range Partition", "List Partition", "Hash Partition", "Declarative Partitioning", "Partition Pruning"}},
		{"Level 24 - Replication", []string{"Streaming Replication", "Logical Replication", "Physical Replication", "Hot Standby", "Failover"}},
		{"Level 25 - High Availability", []string{"WAL Shipping", "Replication Slots", "Patroni", "PgBouncer", "PgPool-II"}},
		{"Level 26 - Full-Text Search", []string{"tsvector", "tsquery", "Ranking", "Dictionaries", "GIN Indexes"}},
		{"Level 27 - Extensions", []string{"uuid-ossp", "pgcrypto", "citext", "hstore", "postgis", "pg_stat_statements", "pg_trgm"}},
		{"Level 28 - Monitoring", []string{"pg_stat_activity", "pg_locks", "pg_stat_database", "pg_stat_user_tables", "pg_stat_statements"}},
		{"Level 29 - Concurrency", []string{"Locks", "Deadlocks", "Advisory Locks", "Row/Table Locks", "Optimistic/Pessimistic Locking"}},
		{"Level 30 - Advanced Storage", []string{"TOAST", "Heap Storage", "WAL Internals", "Buffer Cache & Shared Buffers", "Checkpoints", "Visibility Map"}},
		{"Level 31 - PostgreSQL Internals", []string{"Parser", "Planner", "Optimizer", "Executor", "Buffer/Storage Manager", "Background/WAL Writer", "Checkpointer", "Autovacuum Worker"}},
		{"Level 32 - Distributed PostgreSQL", []string{"Citus", "CockroachDB concepts", "YugabyteDB", "Sharding", "Horizontal Scaling"}},
		{"Level 33 - PostgreSQL with Go", []string{"database/sql, pgx", "Connection Pooling", "Prepared Statements, Transactions", "Batch Inserts, Bulk Copy", "SQL Migrations, Repository Pattern"}},
		{"Level 34 - Interview Topics", []string{"Top 100 SQL queries", "MVCC, VACUUM, WAL", "Index internals", "Isolation levels", "Execution plans", "Partitioning", "JSONB vs JSON", "GIN vs B-Tree", "Deadlocks"}},
		{"Level 35 - Real-World Backend Patterns", []string{"Soft Deletes", "Audit Logs", "Optimistic Locking", "Pagination (Cursor)", "Multi-Tenant", "Event Sourcing", "CQRS", "Outbox Pattern", "Idempotency Keys", "Upserts", "Read Replicas"}},
	},
}

var templateMD = `# %s

## 1. Learning Objectives
* **What you'll learn**: Master the core mechanics of %s in PostgreSQL.
* **Why it matters**: Crucial for building scalable, high-performance, and robust backend systems.

---

## 2. Topics Covered
%s

---

## 3. Production Usage Scenarios (Real-world Examples)
For each concept, here is how we use it in a real production environment at scale:

%s

---

## 4. Code & Query Implementation
### 🔹 Basic Implementation
` + "```sql\n-- Standard query example\nSELECT * FROM table;\n```" + `

### 🔹 Advanced / Optimized Implementation
` + "```sql\n-- Optimized query with indexes or advanced features\n```" + `

---

## 5. Internals & Under the Hood
* **How PostgreSQL engine processes this**: (Parser -> Planner -> Executor)
* **Storage impact**: WAL logs, Heap, and TOAST considerations.

---

## 6. Performance & Benchmarking
* **EXPLAIN ANALYZE impacts**
* **Memory vs Disk Trade-offs**

---

## 7. Best Practices & Common Mistakes
* ✅ **Do**: Follow standard PostgreSQL conventions.
* ❌ **Don't**: Ignore index overhead or transaction locking.
* ⚠️ **Common Mistake**: N+1 queries, missing indexes on foreign keys.

---

## 8. Interview Questions
1. **Beginner**: Explain %s to a junior dev.
2. **Intermediate**: How would you optimize queries involving %s?
3. **Expert**: Describe the low-level locking and memory behavior of %s in high-concurrency environments.
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
	sqlFile, err := os.Create("scripts/seed_postgresql.sql")
	if err != nil {
		panic(err)
	}
	defer sqlFile.Close()

	sqlFile.WriteString("\n-- Seed PostgreSQL Course\n")
	sqlFile.WriteString("INSERT INTO courses (id, slug, title, description, difficulty) VALUES\n")

	c := postgresCourse
	safeTitle := strings.ReplaceAll(c.Title, "'", "''")
	safeDesc := strings.ReplaceAll(c.Description, "'", "''")
	sqlFile.WriteString(fmt.Sprintf("('%s', '%s', '%s', '%s', '%s')\n", c.ID, c.Slug, safeTitle, safeDesc, c.Difficulty))
	sqlFile.WriteString("ON CONFLICT (id) DO UPDATE SET slug = EXCLUDED.slug, title = EXCLUDED.title, description = EXCLUDED.description, difficulty = EXCLUDED.difficulty;\n\n")

	// Delete old lessons for this course to ensure clean state
	sqlFile.WriteString(fmt.Sprintf("DELETE FROM lessons WHERE course_id = '%s';\n\n", c.ID))

	sqlFile.WriteString("-- Seed PostgreSQL Lessons\n")
	sqlFile.WriteString("INSERT INTO lessons (id, course_id, slug, title, content, order_index) VALUES\n")

	courseDir := filepath.Join("content", "lessons", c.Slug)
	os.MkdirAll(courseDir, 0755)

	totalLessons := len(c.Lessons)
	for j, lesson := range c.Lessons {
		title := strings.TrimSpace(lesson.Title)
		
		slugRaw := slugify(title)
		fileName := fmt.Sprintf("%03d-%s.md", j+1, slugRaw)
		slug := fmt.Sprintf("%03d-%s", j+1, slugRaw) 
		filePath := filepath.Join(courseDir, fileName)

		// Format topics and scenarios
		var topicsBuilder strings.Builder
		var scenariosBuilder strings.Builder
		for _, topic := range lesson.Subtopics {
			topicsBuilder.WriteString(fmt.Sprintf("* %s\n", topic))
			scenariosBuilder.WriteString(fmt.Sprintf("### Scenario: %s\n", topic))
			scenariosBuilder.WriteString(fmt.Sprintf("**Context**: Imagine you are building a highly concurrent backend system.\n"))
			scenariosBuilder.WriteString(fmt.Sprintf("**The Problem**: (Define the engineering challenge)\n"))
			scenariosBuilder.WriteString(fmt.Sprintf("**The PostgreSQL Solution**: (How %s solves it effectively in production)\n\n", title))
		}

		// Create markdown file
		mdContent := fmt.Sprintf(templateMD, title, title, topicsBuilder.String(), scenariosBuilder.String(), title, title, title)
		os.WriteFile(filePath, []byte(mdContent), 0644)

		// Generate SQL
		id := fmt.Sprintf("30000000-0000-0000-0000-%012d", j+1)
		
		comma := ","
		if j == totalLessons-1 {
			comma = ""
		}

		safeLessonTitle := strings.ReplaceAll(title, "'", "''")
		sqlFile.WriteString(fmt.Sprintf("('%s', '%s', '%s', '%s', 'See markdown file', %d)%s\n", 
			id, c.ID, slug, safeLessonTitle, j+1, comma))
	}
	// No ON CONFLICT DO UPDATE needed because we DELETEd them above
	sqlFile.WriteString(";\n")

	fmt.Printf("Successfully generated %d PostgreSQL markdown lessons\n", totalLessons)
	fmt.Println("Successfully generated scripts/seed_postgresql.sql")
}
