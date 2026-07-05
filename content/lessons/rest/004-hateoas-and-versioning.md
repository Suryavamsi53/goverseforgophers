# HATEOAS and API Versioning

When you build a REST API, you must assume that multiple different clients (Web, iOS, Android) are consuming it. 

If you make a Breaking Change to the API, the iOS app (which takes 2 weeks to get approved by the Apple App Store) will instantly crash for all users. 
You must plan for evolution.

## 1. API Versioning

There are three common ways to version a REST API.

### Option A: URL Versioning (Most Common)
You embed the version directly into the URL path.
* `GET /api/v1/users`
* `GET /api/v2/users`

**Pros**: Extremely explicit. Easy to route in API Gateways (NGINX).
**Cons**: Purists argue that the URL should only represent the Resource, and the version is not part of the resource's identity.

### Option B: Header Versioning
The URL remains exactly the same (`GET /api/users`), but the client must pass a custom HTTP Header.
* `Accept-Version: v2`

**Pros**: Keeps URLs perfectly clean.
**Cons**: Harder to test in a browser without using Postman/cURL.

### Option C: Content Negotiation (Accept Header)
The most "RESTful" approach. The client tells the server exactly what JSON schema it expects.
* `Accept: application/vnd.mycompany.users.v2+json`

## 2. When to Create a New Version?

You **DO NOT** create a `v2` just because you added a new field to the JSON response. 
Adding a field is a **Non-Breaking Change**. The older `v1` clients will simply ignore the new JSON key they don't recognize.

You ONLY create a new version for a **Breaking Change**:
1. You delete a field that clients were relying on.
2. You rename a field (e.g., `name` to `first_name`).
3. You change the data type (e.g., `id` was an integer, now it's a string).

*Note: In gRPC, versioning is largely solved natively by the backward-compatible nature of Protobuf Field Tags.*

## 3. HATEOAS (Hypermedia as the Engine of Application State)

Roy Fielding (the creator of REST) has famously stated that 99% of APIs calling themselves "RESTful" are actually not RESTful, because they lack **HATEOAS**.

In a standard API, if a client fetches an Order, it gets data:
```json
{
    "id": 42,
    "status": "pending",
    "total": 100.00
}
```
If the client wants to cancel the order, the iOS developer has to read your PDF documentation to figure out that they need to call `POST /orders/42/cancel`. 

**HATEOAS** dictates that the API itself should tell the client exactly what actions are currently available based on the state of the resource!

```json
{
    "id": 42,
    "status": "pending",
    "total": 100.00,
    "links": [
        { "rel": "self", "href": "/orders/42", "method": "GET" },
        { "rel": "cancel", "href": "/orders/42/cancel", "method": "POST" },
        { "rel": "pay", "href": "/orders/42/pay", "method": "POST" }
    ]
}
```

If the order's status changes to `shipped`, the Go backend removes the `cancel` link from the JSON array. The iOS app doesn't need hardcoded logic; it simply loops through the `links` array and automatically hides the "Cancel Order" button on the screen!

*Reality Check: While HATEOAS is brilliant in theory, less than 5% of Enterprise APIs actually implement it, because it is incredibly complex to build and maintain.*
