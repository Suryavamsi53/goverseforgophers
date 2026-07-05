# SQL Injection and XSS

The two most famous vulnerabilities in web development occur because a computer accidentally interprets **User Data** as **Executable Code**.

## 1. SQL Injection (SQLi)

If an attacker types `' OR 1=1; DROP TABLE users; --` into a login email field, and you concatenate that string directly into your SQL query, the database will interpret the attacker's string as a legitimate SQL command and execute it. 
Your entire database is instantly wiped out.

### The Go Solution: Parameterized Queries
You must NEVER use `fmt.Sprintf` or string concatenation to build SQL queries!

```go
// FATAL: The attacker controls the 'email' variable!
query := fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", email)
db.Exec(query)

// SECURE: Parameterized Query
query := "SELECT * FROM users WHERE email = $1"
db.Exec(query, email)
```

**How it works:**
When you use `$1`, the `database/sql` driver sends the query and the data to PostgreSQL in two completely separate network packets. 
PostgreSQL compiles the query first (`SELECT * FROM users WHERE email = ?`). Only after it is compiled does it inject the data. It is mathematically impossible for the data to alter the structure of the compiled query!

## 2. Cross-Site Scripting (XSS)

XSS occurs when an attacker saves malicious JavaScript into your database (e.g., as a comment on a blog post). 
When a victim views the blog post, the attacker's JavaScript is rendered by the victim's browser and executes. The JavaScript steals the victim's JWT token and sends it to the attacker's server!

### The Go Solution: Context-Aware Templating
If you are generating HTML on the server (Server-Side Rendering), you must use the `html/template` package, NEVER the `text/template` package.

```go
import "html/template"

// The template contains a variable injection
tmpl, _ := template.New("test").Parse(`<h1>Hello, {{.Name}}</h1>`)

// The attacker tries to inject a script!
maliciousData := struct{ Name string }{"<script>alert('Hacked!');</script>"}

tmpl.Execute(w, maliciousData)
```

**How it works:**
Because we used `html/template`, Go is context-aware. It sees that `{{.Name}}` is being injected into an HTML context. It automatically HTML-escapes the string before rendering it!
The output becomes:
`<h1>Hello, &lt;script&gt;alert(&#39;Hacked!&#39;);&lt;/script&gt;</h1>`
The browser renders it as harmless text, and the JavaScript does not execute!

## 3. Cross-Site Request Forgery (CSRF)

If you use Cookies to store session IDs instead of JWTs, the browser will automatically attach the Cookie to *every* request made to your domain.

An attacker builds a malicious website: `www.evil.com`. On that website, they place a hidden form:
`<form action="https://yourbank.com/transfer" method="POST">`
When a victim visits `evil.com`, the hidden form submits. Because the victim is logged into `yourbank.com` on another tab, the browser automatically attaches the victim's authentication Cookie to the malicious request! The transfer succeeds!

### The Solution: CSRF Tokens or SameSite Cookies
1. **CSRF Tokens**: The Go server generates a random string (CSRF Token) and places it in the HTML form. When the form submits, the Go server verifies the token. `evil.com` cannot read this token because of the browser's Same-Origin Policy! Use the `github.com/gorilla/csrf` middleware to automate this in Go.
2. **SameSite Cookies**: Modern browsers support the `SameSite=Strict` flag on Cookies. If you set this flag in Go, the browser will absolutely refuse to attach the Cookie if the request originated from a different domain (`evil.com`).
