# CORS and Security

If you deploy your Go API to `https://api.mycompany.com`, and you deploy your React frontend to `https://www.mycompany.com`, you will instantly encounter the most frustrating error in web development: **The CORS Error**.

## 1. The Same-Origin Policy (SOP)

Web browsers enforce a strict security rule called the Same-Origin Policy.
If a Javascript file loaded from `Origin A` attempts to make an HTTP request to `Origin B`, the browser will physically block the Javascript from reading the response!

*An Origin is the combination of Protocol, Domain, and Port (e.g., `https://www.mycompany.com:443`).*

**Why does this exist?**
Imagine you are logged into `facebook.com`. You visit an evil website: `hackers.com`. The Javascript on `hackers.com` makes a silent `GET` request to `facebook.com/messages`. Because of your cookies, Facebook returns your private messages! 
The SOP prevents the evil Javascript from reading that response!

## 2. Cross-Origin Resource Sharing (CORS)

But what if you *want* `www.mycompany.com` to be able to read data from `api.mycompany.com`? They are different origins!

You must configure **CORS** on your Go server. CORS is a set of HTTP Headers that tell the browser to relax the Same-Origin Policy for specific domains.

In your Go API, you must return this header:
`Access-Control-Allow-Origin: https://www.mycompany.com`

When the browser sees this header, it says, "Ah! The Go API explicitly trusts this React app. I will allow the Javascript to read the response!"

## 3. The Preflight Request (OPTIONS)

If the React app tries to make a `POST` request, or send a custom Header (like `Authorization`), the browser gets scared. It doesn't want to send a dangerous action to a server without permission.

The browser will automatically (and silently) pause the `POST` request, and send a preliminary `OPTIONS` request to the Go server. This is the **Preflight**.

The browser asks: *"Hey Go server, are you okay with me sending a POST request with an Authorization header from mycompany.com?"*

The Go server MUST intercept the `OPTIONS` request and respond with:
```text
Access-Control-Allow-Origin: https://www.mycompany.com
Access-Control-Allow-Methods: POST, GET, OPTIONS
Access-Control-Allow-Headers: Authorization, Content-Type
```

Once the browser sees this approval, it finally sends the actual `POST` request!

## 4. The Wildcard Danger (`*`)

Many junior developers get frustrated with CORS errors and blindly configure their Go server to return:
`Access-Control-Allow-Origin: *`

**This is a critical security vulnerability.**
It disables the Same-Origin Policy entirely. Any evil website on the internet can now make Javascript requests to your Go API and read the responses. 
Never use `*` unless you are building a completely public, unauthenticated API (like a public weather API).

## 5. Implementing CORS in Go

Writing the logic to intercept every `OPTIONS` request and inject these headers manually is tedious. 
The industry standard in Go is to use a Middleware package like `github.com/rs/cors`.

```go
c := cors.New(cors.Options{
    AllowedOrigins: []string{"https://www.mycompany.com"},
    AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
    AllowedHeaders: []string{"Authorization", "Content-Type"},
    // Required if you want the browser to send Cookies across origins!
    AllowCredentials: true, 
})

// Wrap your entire router!
handler := c.Handler(mux)
http.ListenAndServe(":8080", handler)
```
