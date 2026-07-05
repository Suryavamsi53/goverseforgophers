# Integration Testing (Testcontainers)

Unit Tests verify that your `UserService` logic works. But how do you verify that your `PostgresUserRepository` actually writes data to Postgres? 

If you Mock the database, you aren't testing the SQL! You might have a syntax error in your SQL (`INSERT ITNO users`), and your Unit Tests will pass 100%, but Production will crash instantly.

You must write **Integration Tests** to verify that your Go code correctly communicates with real infrastructure (Postgres, Redis, Kafka).

## 1. The Environment Problem

Historically, Integration Testing was a nightmare. 
To run the tests on your laptop, you had to manually install PostgreSQL on your Mac. When the CI/CD pipeline ran in GitHub Actions, the DevOps team had to write complex scripts to boot up a Postgres container before running `go test`. If a test crashed and didn't clean up the database, the next test would fail due to leftover data!

## 2. Testcontainers (The Industry Standard)

The modern solution is **Testcontainers** (`github.com/testcontainers/testcontainers-go`).

Testcontainers is a Go library that programmatically controls Docker directly from your `_test.go` file!

When you run `go test`, the Go code literally instructs Docker to download a PostgreSQL image, boot up an isolated container, map a random open port to it, and return the connection string to the test. When the test finishes, Go automatically destroys the container.

## 3. Implementation Example

```go
package repository_test

import (
    "context"
    "testing"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestPostgresUserRepository(t *testing.T) {
    ctx := context.Background()

    // 1. Programmatically boot a REAL Postgres database!
    // This blocks until the database is fully ready to accept connections.
    pgContainer, err := postgres.RunContainer(ctx,
        testcontainers.WithImage("postgres:15-alpine"),
        postgres.WithDatabase("testdb"),
        postgres.WithUsername("user"),
        postgres.WithPassword("pass"),
    )
    if err != nil {
        t.Fatal(err)
    }

    // 4. CRITICAL: Automatically destroy the container when the test ends!
    defer pgContainer.Terminate(ctx)

    // 2. Extract the random connection string
    connStr, _ := pgContainer.ConnectionString(ctx)

    // 3. Connect our real Repository to the isolated test database!
    db, _ := sql.Open("postgres", connStr)
    repo := repository.NewPostgresUserRepository(db)

    // Run migrations (e.g. CREATE TABLE users)
    RunMigrations(db)

    // Execute the real SQL!
    err = repo.Save(ctx, &User{Name: "Integration Test"})
    if err != nil {
        t.Errorf("Failed to save user: %v", err)
    }
}
```

## 4. The Benefits

1. **Zero Configuration**: Any developer can clone your repo and run `go test`. They don't need to install Postgres, Redis, or Kafka. They just need Docker installed on their laptop.
2. **Total Isolation**: Because Testcontainers boots a brand new, empty database for the test, there is absolutely zero risk of dirty data from previous tests failing your assertions.
3. **CI/CD Native**: GitHub Actions and GitLab CI natively support Docker-in-Docker, meaning these tests run flawlessly in the cloud without any complex YAML orchestration.

## 5. Separation of Test Suites

Because booting a Docker container takes ~2 seconds, Integration Tests are slow. If you have 500 Integration Tests, your test suite will take 20 minutes to run.

You must separate Unit Tests from Integration Tests using Go **Build Tags**.

At the top of `user_repo_test.go`:
```go
//go:build integration
```

Now, when you run `go test ./...`, Go will **ignore** the file! It only runs the blazing-fast unit tests.
When the CI/CD pipeline runs the nightly build, it explicitly requests them:
`go test -tags=integration ./...`
