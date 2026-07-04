package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

var ips = []string{
	"192.168.1.10", "192.168.1.11", "10.0.0.5", "10.0.0.23", "172.16.0.2",
	"8.8.8.8", "1.1.1.1", "104.21.43.12", "142.250.190.46", "[::1]",
	"[2001:0db8:85a3:0000:0000:8a2e:0370:7334]",
}
var endpoints = []string{"/", "/api/users", "/api/login", "/api/checkout", "/api/products", "/favicon.ico"}
var methods = []string{"GET", "POST", "PUT", "DELETE", "PATCH"}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func main() {
	fileName := "production.log"
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	rand.Seed(time.Now().UnixNano())
	now := time.Now().Add(-24 * time.Hour) // start 24 hours ago

	totalLogs := 500000
	fmt.Printf("Generating %d production log entries...\n", totalLogs)

	for i := 1; i <= totalLogs; i++ {
		ip := ips[rand.Intn(len(ips))]
		port := rand.Intn(50000) + 10000
		
		method := methods[rand.Intn(len(methods))]
		endpoint := endpoints[rand.Intn(len(endpoints))]
		
		status := 200
		size := rand.Intn(50000) + 100
		duration := time.Duration(rand.Intn(100)) * time.Millisecond // 0-100ms
		
		// Traffic Pattern Skewing to make it realistic
		randVal := rand.Float32()
		
		if endpoint == "/api/login" && method == "POST" {
			if rand.Float32() < 0.1 {
				status = 401 // 10% Unauthorized on login
			} else {
				status = 200
			}
			duration += 150 * time.Millisecond // Logins take longer
		} else if endpoint == "/api/checkout" {
			// Checkout sometimes has high latency or fails (interviews love this!)
			if randVal < 0.05 {
				status = 500
				duration += 2 * time.Second
			} else if randVal < 0.15 {
				status = 502
				duration += 5 * time.Second
			} else {
				status = 201
				duration += 300 * time.Millisecond
			}
		} else if endpoint == "/favicon.ico" {
			status = 404
			size = 19
			duration = time.Duration(rand.Intn(50)) * time.Microsecond
		} else if randVal < 0.6 { // 60% standard successful requests
			method = "GET"
			status = 200
		} else if randVal < 0.7 {
			status = 400
		} else if randVal < 0.75 {
			status = 404
		} else if randVal < 0.78 {
			status = 429 // rate limiting
		}
		
		// Introduce some extreme long-tail latency (99th percentile spikes)
		if rand.Float32() < 0.01 { // 1% of requests are terribly slow
			duration += time.Duration(rand.Intn(5000)) * time.Millisecond
		}

		// Advance time slightly between 10ms to 50ms per request
		now = now.Add(time.Duration(rand.Intn(40)+10) * time.Millisecond)
		timestamp := now.Format("2006/01/02 15:04:05")
		
		reqID := fmt.Sprintf("fedora/%s-%06d", randomString(10), i)
		ipWithPort := fmt.Sprintf("%s:%d", ip, port)
		if ip[0] == '[' {
			ipWithPort = fmt.Sprintf("%s:%d", ip, port) // [::1]:port
		}

		// Format: 2026/07/04 10:41:31 [fedora/Hl8HBjSjWt-000001] "GET http://localhost:8080/ HTTP/1.1" from [::1]:35196 - 200 15624B in 3.84544ms
		logLine := fmt.Sprintf(`%s [%s] "%s http://localhost:8080%s HTTP/1.1" from %s - %d %dB in %v`+"\n",
			timestamp, reqID, method, endpoint, ipWithPort, status, size, duration)
		
		file.WriteString(logLine)
	}
	
	fmt.Println("Generated production.log successfully!")
}
