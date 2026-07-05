# RBAC and Authorization

Authentication answers *"Who are you?"* (JWT).
Authorization answers *"Are you allowed to do this?"* (RBAC).

If User 42 is authenticated, can they call `DELETE /api/users/99`? 
If you forget to implement Authorization, any authenticated user can delete the entire database. This is known as a **Broken Object Level Authorization (BOLA)** vulnerability, and it is the #1 most common API vulnerability in the world.

## 1. Role-Based Access Control (RBAC)

The industry standard for API authorization is RBAC. 
Instead of assigning permissions directly to a User, you assign permissions to a **Role**, and then you assign Roles to the User.

1. **Role**: `Admin`, `Manager`, `Viewer`.
2. **Permissions**: `users:delete`, `reports:read`.

When the user authenticates, the JWT payload contains their Role:
`{"user_id": 42, "role": "admin"}`

## 2. Implementing RBAC Middleware in Go

Authorization should never be hardcoded into your core Business Logic. It must be implemented as an HTTP Middleware that intercepts the request *before* it hits your HTTP Handler.

```go
// 1. A dummy permission checker
func hasPermission(role string, requiredPermission string) bool {
    // In a real app, this might query a fast in-memory map or Redis
    permissions := map[string][]string{
        "admin":  {"users:delete", "users:read", "users:write"},
        "viewer": {"users:read"},
    }
    
    for _, p := range permissions[role] {
        if p == requiredPermission { return true }
    }
    return false
}

// 2. The Authorization Middleware Factory
func RequirePermission(required string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            
            // Extract the role from the Context (injected earlier by the JWT Middleware!)
            role := r.Context().Value("user_role").(string)

            if !hasPermission(role, required) {
                http.Error(w, "Forbidden", http.StatusForbidden) // 403!
                return
            }

            // User is authorized! Proceed to the actual handler.
            next.ServeHTTP(w, r)
        })
    }
}
```

Now you can flawlessly protect your routes in `main.go`:
```go
router.Handle("DELETE /users/{id}", RequirePermission("users:delete")(deleteUserHandler))
```

## 3. ABAC (Attribute-Based Access Control)

RBAC has a major limitation. 
What if you want to allow a User to edit a Post, but ONLY if they are the original author of the Post?

A static "Role" cannot solve this, because the permission depends on the *Attribute* of the specific database row being accessed!

This requires **ABAC**.
ABAC logic is usually pushed down into the Service Layer (Business Logic) or the Database Query itself.

```go
func (s *PostService) EditPost(ctx context.Context, userID int, postID int) error {
    // ABAC Policy: You can only edit your own posts!
    post, _ := s.repo.GetPost(postID)
    
    if post.AuthorID != userID {
        // We reject the action because the Attributes don't match!
        return errors.New("unauthorized: you are not the author")
    }
    
    // Proceed with edit...
}
```

## 4. OPA (Open Policy Agent)

If you have 50 microservices written in Go, Java, and Python, implementing RBAC and ABAC rules in 3 different languages is a maintenance nightmare.

Enterprise companies use **OPA (Open Policy Agent)**.
OPA is a dedicated sidecar container that holds all your authorization rules written in a specialized language called `Rego`.

Instead of the Go application checking its own roles, the Go application makes a blazing-fast local gRPC call to the OPA sidecar: *"User 42 wants to delete Post 99. Is this allowed?"* OPA evaluates the complex enterprise policies and simply returns `true` or `false`!
