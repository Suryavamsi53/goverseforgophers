# Headers, Cookies, and State

HTTP is fundamentally a **Stateless** protocol. The Go server does not inherently remember that Request B came from the exact same user as Request A.

To maintain state (like being logged into an account), we must use HTTP Headers and Cookies.

## 1. HTTP Headers

Headers are key-value pairs that provide metadata about the request or the response. 

* **`Content-Type`**: Tells the receiver how to parse the Body (e.g., `application/json` or `text/html`). If you forget to set this in Go, the `http.ResponseWriter` will automatically try to guess (Sniff) the content type by reading the first 512 bytes of the response!
* **`Authorization`**: Used to send JWTs or Bearer tokens.
* **`User-Agent`**: Tells the server what browser or device is making the request.

```go
// Setting a header in Go
w.Header().Set("Content-Type", "application/json")
w.Header().Set("X-Custom-Header", "Goverse")
```

## 2. What is a Cookie?

A Cookie is just a specific HTTP Header, but it has a magical property: **Web browsers automatically manage them.**

If your Go server sends this header in the Response:
`Set-Cookie: session_id=12345;`

The Web Browser intercepts this header. It saves `session_id=12345` to its local hard drive. 
For every single future request the browser makes to your domain, it will automatically attach this header:
`Cookie: session_id=12345`

You do not have to write any JavaScript to make this happen! It is built into the HTTP specification of every browser.

## 3. Cookie Security Attributes (Critical!)

If a hacker writes malicious JavaScript to steal the `document.cookie` variable, your users' accounts will be compromised (XSS).

To prevent this, you must configure strict security attributes when setting Cookies in Go.

```go
http.SetCookie(w, &http.Cookie{
    Name:     "session_id",
    Value:    "12345",
    Path:     "/",
    MaxAge:   3600, // Expires in 1 hour
    
    // SECURITY ATTRIBUTES:
    HttpOnly: true,  // JavaScript CANNOT access this cookie! Defeats XSS!
    Secure:   true,  // Only send over HTTPS!
    SameSite: http.SameSiteStrictMode, // Defeats CSRF attacks!
})
```

## 4. LocalStorage vs Cookies

Modern React/Angular developers often prefer storing JWTs in `LocalStorage` instead of Cookies.

* **LocalStorage Pros**: Easy to access via JavaScript. Prevents CSRF attacks natively (because the browser doesn't automatically attach LocalStorage data to requests).
* **LocalStorage Cons**: Massively vulnerable to XSS! Any rogue NPM package on your website can instantly steal the JWT from LocalStorage.

* **Cookie Pros**: Using `HttpOnly` completely neutralizes XSS token theft!
* **Cookie Cons**: Requires strict `SameSite` configurations to prevent CSRF attacks.

**Enterprise Standard**: The most secure way to store a session token in a web browser is inside an `HttpOnly`, `Secure`, `SameSite=Strict` Cookie.
