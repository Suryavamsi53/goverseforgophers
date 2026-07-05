# Resource-Oriented Design

REST (Representational State Transfer) is an architectural style for designing networked applications. It was defined by Roy Fielding in 2000 and has become the de facto standard for building web APIs.

The core principle of REST is **Resource-Oriented Design**.

## 1. What is a Resource?

In RPC (Remote Procedure Call) architectures, your API revolves around Actions (Verbs).
* `POST /createUser`
* `POST /deleteUser`
* `POST /getAllUsers`

In REST, your API revolves around **Resources (Nouns)**. A resource is any entity your application manages (e.g., Users, Posts, Comments, Invoices).
Instead of changing the URL to define the action, you keep the URL (the Noun) exactly the same, and you change the HTTP Method (the Verb) to define the action!

## 2. URL Naming Conventions

There are strict industry standards for naming RESTful URLs. If you violate these, other engineers will struggle to use your API.

### Rule 1: Always use Plural Nouns
* ❌ `GET /user/42` (Bad)
* ❌ `GET /get_user/42` (Terrible)
* ✅ `GET /users/42` (Good)

### Rule 2: Keep URLs flat (Max 2 levels deep)
What if you want to get all comments on a specific post written by a specific user?
* ❌ `GET /users/42/posts/99/comments` (Too deep, hard to parse).

Instead, you flatten the relationship. If Post 99 is globally unique, you don't need the User ID in the URL!
* ✅ `GET /posts/99/comments`

### Rule 3: Use Kebab-Case
Never use camelCase or snake_case in URLs.
* ❌ `GET /api/userProfiles`
* ✅ `GET /api/user-profiles`

## 3. The 4 Core Principles of REST

To be truly RESTful, your Go API must adhere to these constraints:

1. **Client-Server Architecture**: The Go backend and the React frontend must be completely decoupled. They only communicate via HTTP JSON.
2. **Statelessness**: The Go server must NEVER store a "Session" for a client in RAM. Every single HTTP request from the client must contain all the information necessary to authenticate and authorize it (e.g., passing a JWT in the Header on every request).
3. **Cacheability**: The server must explicitly tell the client if a JSON response can be cached (using `Cache-Control` headers) to reduce network traffic.
4. **Uniform Interface**: All resources must be accessed using a standardized, predictable format (the URL conventions and standard HTTP methods).

## 4. Designing Sub-Resources

How do you design an API where a User "Likes" a Post? 
A "Like" isn't a physical resource you edit. It is an action. 

**Option A (The Sub-Resource creation):**
Treat the "Like" as a sub-resource.
`POST /posts/99/likes` (Creates a new Like entity).
`DELETE /posts/99/likes` (Removes the Like).

**Option B (The Custom Verb workaround):**
Sometimes, REST is too rigid. If you have an action like "Refund Payment", treating it as a resource (`POST /refunds`) might be too complex. 
The industry acceptable workaround is to use a custom verb at the end of the URL, preceded by a colon:
`POST /payments/123:refund`
