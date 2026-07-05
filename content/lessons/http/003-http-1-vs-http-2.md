# HTTP/1.1 vs HTTP/2 vs HTTP/3

For two decades, the internet ran on HTTP/1.1. It is simple, text-based, and human-readable. 
But as modern websites grew from containing 3 images to containing 150 Javascript files, CSS files, and dynamic API requests, HTTP/1.1 became a massive bottleneck.

## 1. The HTTP/1.1 Bottleneck (Head-of-Line Blocking)

In HTTP/1.1, you can leave a TCP connection open (`Keep-Alive`). But you can only send **one request at a time** over that connection.

If a browser needs to download `script.js` (10MB) and `style.css` (10KB):
1. Browser requests `script.js`.
2. The TCP socket is now locked!
3. Browser has to wait 5 seconds for the massive 10MB script to download.
4. Only then can it request the tiny 10KB CSS file.

This is called **Head-of-Line Blocking**. The CSS file was blocked by the massive Javascript file ahead of it in the line!
*(Workaround: Browsers started opening 6 simultaneous TCP connections to the same server to download files in parallel, but this wastes huge amounts of server RAM).*

## 2. HTTP/2 (The Multiplexing Revolution)

In 2015, HTTP/2 solved this problem.

**It is fundamentally a Binary protocol, not Text.**
Instead of sending raw text, HTTP/2 splits the requests into tiny binary "Frames" and interleaves them!

Over a **single TCP connection**, the server can send a chunk of `script.js`, then a chunk of `style.css`, then another chunk of `script.js`! 

* **Multiplexing**: You can send 1,000 concurrent API requests over 1 TCP connection simultaneously without any blocking!
* **Header Compression (HPACK)**: In HTTP/1.1, the massive `User-Agent` and `Cookie` headers are sent in plain text on every single request. HTTP/2 compresses the headers, saving massive amounts of bandwidth.

**Go Support**: If you use `http.ListenAndServeTLS` in Go, your server automatically upgrades clients to HTTP/2 transparently! No code changes required! (This is also the underlying technology that makes gRPC possible).

## 3. HTTP/3 (Killing TCP)

HTTP/2 solved HTTP Head-of-Line blocking, but it exposed a deeper flaw: **TCP Head-of-Line blocking**.

Because HTTP/2 sends 1,000 requests over a single TCP connection, if the Wi-Fi drops a single TCP packet, the OS Kernel halts the *entire* TCP connection until that specific packet is re-transmitted. All 1,000 concurrent requests freeze because of one dropped packet!

**HTTP/3** abandons TCP entirely!
It runs on top of **UDP** using a new transport protocol built by Google called **QUIC**.

* UDP doesn't guarantee order. If a packet drops, only the specific file associated with that packet stalls. The other 999 requests continue streaming flawlessly!
* TCP requires a 3-Way Handshake. TLS requires another Handshake. QUIC combines them, allowing clients to establish a secure, encrypted connection in **0-RTT** (Zero Round Trips)!

HTTP/3 is actively being rolled out across the internet, heavily supported by CDNs like Cloudflare and Google.
