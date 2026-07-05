# Database Migrations (Schema as Code)

## 1. Learning Objectives
* **What you'll learn**: The principles of treating your Database Schema as Code, and how to execute safe, version-controlled migrations using `golang-migrate`.
* **Why it matters**: Mutating production databases manually via a GUI tool (like pgAdmin or DBeaver) guarantees catastrophic failure in a team environment. Migrations ensure every developer, CI/CD pipeline, and server has the exact same schema.
* **Where it's used**: Every single professional software engineering team in the world.

---

## 2. Real-world Story
Imagine building a house (The App). You realize the kitchen (The Database) needs an extra sink. If Bob just walks in with a sledgehammer and installs a sink, Alice (who is building the plumbing) has no idea what happened, and everything breaks.
Database Migrations are the Architect's Blueprints. If Bob wants a sink, he writes a numbered blueprint (`001_add_sink.up.sql`). The CI/CD pipeline applies the blueprint. Now, Alice, the Staging Server, and the Production Server all read the blueprint and get the exact same sink, deterministically.

---

## 3. Visual Learning (Execution Flow & Architecture)
```mermaid
graph TD
    A[git push to main] --> B(GitHub Actions)
    B -->|1. Run Tests| C{Tests Pass?}
    C -->|Yes| D[Download migrate CLI]
    D -->|2. Connect to Prod DB| E[(PostgreSQL)]
    E --> F{Check schema_migrations table}
    F -->|Version 5 installed| G[Run 006_add_users.up.sql]
    G --> H[Update schema_migrations to 6]
    H --> I[Deploy new Go App!]
```

---

## 4. Internal Working (Under the Hood)
How does the Migration engine know which SQL files to run?
1. The engine maintains a hidden table in your database called `schema_migrations` (or similar).
2. This table stores a single integer (e.g., `Version: 5`).
3. You have files like `005_init.up.sql` and `006_add_index.up.sql`.
4. When you run `migrate up`, the engine sees the DB is on Version 5, and automatically executes `006_add_index.up.sql`, updating the table to Version 6.

---

## 5. Compiler Behavior
* **Embedding Migrations**: In Go 1.16+, you can use the `//go:embed` directive to physically bundle all your `.sql` migration files directly into the compiled Go binary. This means your Docker image is just one file—no need to copy a separate `migrations/` folder into the container!

---

## 6. Memory Management
* **Zero Overhead**: Migrations execute only once during startup or deployment. They consume zero memory during the actual execution of your Go business logic.

---

## 7. Code Examples

### 🔹 Example 1: Simple
```sql
-- 001_create_users_table.up.sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 001_create_users_table.down.sql
DROP TABLE IF EXISTS users;
```

### 🔹 Example 2: Intermediate
```bash
# Running migrations via the golang-migrate CLI
migrate -path db/migrations -database "postgres://user:pass@localhost:5432/mydb?sslmode=disable" up
```

### 🔹 Example 3: Advanced
```go
// Running migrations inside Go programmatically using go:embed
import (
    "embed"
    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed db/migrations/*.sql
var fs embed.FS

func RunMigrations(dbURL string) {
    d, err := iofs.New(fs, "db/migrations")
    m, err := migrate.NewWithSourceInstance("iofs", d, dbURL)
    
    // Automatically apply all pending migrations!
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatal(err)
    }
}
```

### 🔹 Example 4: Production
```sql
-- Safe Migrations for massive tables
-- Adding a column to a 1-billion row table without a default value is instantaneous.
-- Adding it WITH a default value will lock the table for hours in older Postgres!
ALTER TABLE orders ADD COLUMN status VARCHAR(50); -- FAST
```

### 🔹 Example 5: Interview
```go
// Q: Should the Go App run migrations on startup (e.g., in main.go)?
// A: In a simple app, yes. In a Kubernetes cluster with 50 replicas, NO! 
// 50 pods starting at the same time will all try to run migrations, causing race conditions and deadlocks. 
// Run migrations in a dedicated CI/CD Job before the new app pods deploy.
```

---

## 8. Production Examples
1. **Rollbacks**: If you deploy `002_add_billing.up.sql` and it breaks production, you run `migrate down 1` via your CI/CD pipeline, and it automatically executes `002_add_billing.down.sql`, saving the company.
2. **Local Development**: When a new developer joins the team, they just spin up an empty Postgres docker container and run `migrate up`. They instantly have the exact same database schema as production.

---

## 9. Performance & Benchmarking
* **Locking**: Migrations alter the physical structure of tables (DDL). They require heavy `ACCESS EXCLUSIVE` locks on the database. If you drop a column during peak traffic hours, it will block all incoming Go `SELECT` requests until the migration finishes, causing a massive latency spike.

---

## 10. Best Practices
* ✅ **Do**: Make your migrations strictly additive (Adding tables, adding columns) as much as possible.
* ❌ **Don't**: Drop columns or rename columns unless you follow the "Expand and Contract" deployment pattern. (If you rename a column, the currently running Go app will instantly crash because it's looking for the old name!).
* 🏢 **Google / Uber / Netflix Style**: Use tools like `gh-ost` or `pg-osc` for Zero-Downtime schema migrations on massive tables without holding heavy locks.

---

## 11. Common Mistakes
1. **Modifying Existing Migrations**: Never modify `001_init.up.sql` after it has been pushed to the main branch. Your local DB will have the new change, but production DB already ran `001` weeks ago and will ignore your modification! Always create a new file `002_update.up.sql`.
2. **Missing `down.sql` files**: Every `up.sql` MUST have a perfectly corresponding `down.sql`. If you add a table, the down file drops it.

---

## 12. Debugging
How to troubleshoot Migrations in production:
* **The "Dirty" State**: If a migration fails halfway through (e.g., a syntax error in your SQL), the `schema_migrations` table marks the version as `dirty = true`. You cannot run any more migrations until you manually fix the database state and run `migrate force <version>`.

---

## 13. Exercises
1. **Easy**: Install `golang-migrate` CLI and generate an up/down file pair.
2. **Medium**: Write the SQL to create a `products` table and successfully `migrate up`.
3. **Hard**: Modify your Go server `main.go` to use `go:embed` to run migrations automatically on boot.
4. **Expert**: Recreate a "Dirty State" failure by adding a syntax error to a migration, and use the CLI to force-fix the version back to a healthy state.

---

## 14. Quiz
1. **MCQ**: Why is renaming a column in a single migration highly dangerous in a zero-downtime microservice architecture?
   * (A) Postgres doesn't support renaming (B) The old Go pods (V1) are still running and will crash when the DB schema suddenly changes to V2 (C) It takes up too much RAM. *(Answer: B)*
2. **Code Review**: `DROP TABLE users;` inside an `up.sql` file. Why is this a career-ending move? *(You just deleted production data with no backup).*

---

## 15. FAANG Interview Questions
* **Beginner**: What is the purpose of the `schema_migrations` table?
* **Intermediate**: Explain the "Expand and Contract" pattern for zero-downtime database migrations.
* **Senior (Google/Meta)**: How do you handle schema migrations across a geo-distributed multi-master database cluster without causing replication conflicts?

---

## 16. Mini Project
**Zero-Downtime Column Rename**
* Goal: Rename `first_name` to `given_name`.
* Step 1: Create a migration to add `given_name`. Deploy App V2 that writes to both columns but reads from `first_name`.
* Step 2: Create a migration to copy existing data. Deploy App V3 that reads from `given_name`.
* Step 3: Create a migration to drop `first_name`.

---

## 17. Enterprise Features & Observability
* **Atlas & Prisma**: Modern declarative migration tools (like `ariga/atlas`) allow you to write your desired schema in HCL (Terraform style) or Go code, and the tool automatically calculates the diff and generates the SQL migration files for you!

---

## 18. Source Code Reading
Walkthrough of `github.com/golang-migrate/migrate`.
* **The Driver Interface**: Look at how the engine abstracts the backend. The core engine doesn't care if it's migrating Postgres, MySQL, or Cassandra. It simply uses standard Go interfaces to read bytes from the file source and send them to the database driver.

---

## 19. Architecture
* **The Migration Container**: In Kubernetes, migrations are best run as a `Job` via an `initContainer` or a Helm hook, guaranteeing the schema is fully updated *before* the new Go application Pods are allowed to boot up and serve traffic.

---

## 20. Summary & Cheat Sheet
* **Tool**: `golang-migrate` CLI.
* **File Naming**: `<version>_<title>.up.sql`.
* **Immutability**: Never change an applied migration. Add a new one.
* **Automation**: Run via CI/CD or Go `init()`.
