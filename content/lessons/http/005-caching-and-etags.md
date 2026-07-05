# Caching and ETags

In the Redis module, we learned how to cache data on the Server. 
But the fastest HTTP request is the one that never leaves the user's computer! 

HTTP has a brilliantly complex caching mechanism built directly into the protocol. You can instruct the Web Browser (or intermediary CDNs like Cloudflare) to cache your Go API responses locally!

## 1. Cache-Control (Time-Based Caching)

The simplest form of caching is telling the browser exactly how long it can keep the file.

If your Go server returns a profile image, you inject this header:
`Cache-Control: public, max-age=3600`

* `max-age=3600`: The browser will save this image to the hard drive for 3,600 seconds (1 hour).
* `public`: Any intermediary server (like Cloudflare or a corporate proxy) is also allowed to cache this file and serve it to other users! (Use `private` if the data contains sensitive user information).

For the next hour, if the user visits the page again, the browser will NOT make a network request to your Go server. It instantly loads the image from the local hard drive!

## 2. The Invalidation Problem

What if the user updates their profile image 5 minutes later?
The browser will still show the old image for another 55 minutes because of the `max-age`!

This is why you **never** cache dynamic API JSON responses (`/users/42`) using large `max-age` values. 

## 3. ETags (Validation-Based Caching)

Instead of relying on time, we can rely on **Hashes**.

1. The React app requests `GET /users/42`.
2. The Go server fetches the user from Postgres, generates the JSON, and hashes the JSON string into an **ETag** (Entity Tag).
3. The Go server returns the JSON and the Header: `ETag: "hash123"`.
4. The browser saves the JSON and the ETag.

Later, the user refreshes the page:
1. The browser makes the request, but includes a special header: `If-None-Match: "hash123"`.
2. The Go server fetches the user from Postgres, generates the JSON, and hashes it. 
3. The Go server sees the new hash is still `"hash123"`. The data hasn't changed!
4. **The Magic**: Instead of sending the massive 5MB JSON string back over the network, the Go server simply returns a `304 Not Modified` status code with an empty body!

The browser sees the `304` and instantly loads the JSON from its local hard drive. You just saved 5MB of network bandwidth!

## 4. Cache Busting (Static Assets)

If you have a React app, it compiles to `main.js`. If you cache `main.js` for 1 year, and you push a critical bug fix, no one will get the update!

Webpack/Vite solves this using **Cache Busting**. 
Every time you compile the React app, it injects a hash of the file contents into the filename itself!
* `main.a4b9c.js`

You can safely configure your Go server to cache this file for 10 years (`Cache-Control: max-age=31536000, immutable`). 
If you update the code, Webpack generates a brand new filename (`main.f82d1.js`). The HTML file points to the new filename, forcing the browser to download it, instantly bypassing the old cache!
