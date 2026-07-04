package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

var ips = []string{"192.168.1.1", "10.0.0.5", "172.16.0.2", "8.8.8.8", "1.1.1.1"}
var endpoints = []string{"/", "/api/users", "/api/login", "/about", "/contact", "/products"}
var methods = []string{"GET", "POST", "PUT", "DELETE"}
var statuses = []int{200, 201, 400, 401, 403, 404, 500, 502, 503}

func main() {
	file, err := os.Create("server.log")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	rand.Seed(time.Now().UnixNano())
	now := time.Now()

	for i := 0; i < 50000; i++ {
		ip := ips[rand.Intn(len(ips))]
		method := methods[rand.Intn(len(methods))]
		endpoint := endpoints[rand.Intn(len(endpoints))]
		status := statuses[rand.Intn(len(statuses))]
		
		// Skew data to make it realistic
		if rand.Float32() < 0.7 {
			method = "GET"
			status = 200
		}

		timestamp := now.Add(time.Duration(i) * time.Second).Format("02/Jan/2006:15:04:05 -0700")
		
		logLine := fmt.Sprintf(`%s - - [%s] "%s %s HTTP/1.1" %d 1024 "-" "Mozilla/5.0"`+"\n", ip, timestamp, method, endpoint, status)
		file.WriteString(logLine)
	}
	fmt.Println("Generated server.log with 50,000 entries")
}
