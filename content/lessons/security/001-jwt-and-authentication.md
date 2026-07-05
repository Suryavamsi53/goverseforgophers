# JWT and Authentication

Authentication answers the question: *"Who are you?"*

In a traditional Monolithic architecture, when a user logs in, the Go server creates a Session object in RAM (or Redis) and sends a random `session_id` back to the browser via a Cookie. On every subsequent request, the Go server looks up the `session_id` in Redis to see who the user is.

In a massive Microservice architecture, checking a central Redis cluster for every single HTTP request creates a massive bottleneck.
We solve this using **Stateless Authentication: JSON Web Tokens (JWT)**.

## 1. What is a JWT?

A JWT is a Base64 encoded JSON string that contains information about the user. 
It consists of 3 parts separated by dots (`Header.Payload.Signature`).

* **Header**: Contains the hashing algorithm (e.g., `HS256` or `RS256`).
* **Payload**: The actual JSON data (e.g., `{"user_id": 42, "role": "admin", "exp": 1700000000}`).
* **Signature**: A cryptographic hash of the Header + Payload + a Secret Key.

Because the JWT is stateless, the `Order Service` does not need to talk to the `Auth Service` or Redis to verify the user! It simply runs the cryptographic hash function locally. If the resulting hash matches the Signature attached to the token, the token is 100% authentic and hasn't been tampered with!

## 2. JWT Implementation in Go

Never write your own cryptography code. Use `github.com/golang-jwt/jwt/v5`.

```go
// 1. Generating a Token (Usually done by the Auth Service)
func GenerateToken(userID int) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(time.Hour * 24).Unix(), // Expires in 24 hrs
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    // The Secret Key must NEVER leave the server!
    return token.SignedString([]byte("my-super-secret-key"))
}

// 2. Verifying a Token (Done by every microservice)
func VerifyToken(tokenString string) (*jwt.Token, error) {
    return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        // Validate the algorithm
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method")
        }
        return []byte("my-super-secret-key"), nil
    })
}
```

## 3. The Security Flaw of JWTs (Revocation)

If a hacker steals a user's JWT, the hacker has full access to the account. 
In a Stateful Session system, the Admin can just delete the `session_id` from Redis, and the hacker is instantly logged out.

**You cannot do this with JWTs!**
Because JWTs are stateless, they are mathematically valid until their `exp` (Expiration Time) runs out. You cannot simply "delete" a JWT.

### The Refresh Token Solution
To mitigate this, you must keep JWT lifespans incredibly short.
1. The Auth Service issues a **Short-Lived Access Token** (expires in 15 minutes) and a **Long-Lived Refresh Token** (expires in 30 days).
2. The Go API only accepts the 15-minute Access Token.
3. If the Access Token expires, the client silently sends the Refresh Token to the Auth Service.
4. The Auth Service checks the Database. If the Refresh Token is valid (not revoked), it issues a brand new 15-minute Access Token.

If a hacker steals an Access Token, they only have access for a maximum of 15 minutes. If an Admin bans the user, they delete the Refresh Token from the database. When the 15 minutes are up, the user is permanently locked out!

## 4. Symmetric vs Asymmetric Keys (RS256)

In the Go code above, we used `HS256` (Symmetric). This means the `Order Service` needs the exact same Secret Key as the `Auth Service` to verify the token. Sharing the Secret Key across 50 microservices is a massive security risk!

Enterprise architectures use **RS256 (Asymmetric RSA Keys)**.
* The `Auth Service` holds the **Private Key** (used to generate the token).
* The 50 other microservices hold the **Public Key** (used to verify the token).
This guarantees that if the `Order Service` is hacked, the hacker only gets the Public Key, which cannot be used to forge new fake JWTs!
