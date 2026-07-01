package runner

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
)

type WSTerminalMessage struct {
	Type  string            `json:"type"` // "init", "input", "resize"
	Data  string            `json:"data,omitempty"`
	Files map[string]string `json:"files,omitempty"`
	Rows  uint16            `json:"rows,omitempty"`
	Cols  uint16            `json:"cols,omitempty"`
}

func HandleWSTerminalSession(ws *websocket.Conn) {
	defer ws.Close()

	// Wait for init message
	_, initMsg, err := ws.ReadMessage()
	if err != nil {
		log.Println("Error reading init message:", err)
		return
	}

	var req WSTerminalMessage
	if err := json.Unmarshal(initMsg, &req); err != nil || req.Type != "init" {
		log.Println("Invalid init message")
		return
	}

	// Create temp directory for the session
	tmpDir, err := os.MkdirTemp("", "goverse-pty-*")
	if err != nil {
		log.Println("Error creating temp dir:", err)
		return
	}
	defer os.RemoveAll(tmpDir)

	// Dump files
	for path, content := range req.Files {
		fullPath := filepath.Join(tmpDir, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			continue
		}
		os.WriteFile(fullPath, []byte(content), 0644)
	}

	// Ensure there's a go.mod
	if _, err := os.Stat(filepath.Join(tmpDir, "go.mod")); os.IsNotExist(err) {
		os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module example\n\ngo 1.21\n"), 0644)
	}

	// Start shell with PTY inside Docker
	cmd := exec.Command("docker", "run", "-it", "--rm",
		"--memory", "256m",
		"-cpus", "0.5",
		"-v", tmpDir+":/app:z",
		"-w", "/app",
		"-e", "TERM=xterm-256color",
		"golang:1.21-alpine",
		"/bin/sh")

	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Println("Error starting pty:", err)
		return
	}
	defer func() { _ = ptmx.Close() }()

	// Resize if provided in init
	if req.Rows > 0 && req.Cols > 0 {
		pty.Setsize(ptmx, &pty.Winsize{Rows: req.Rows, Cols: req.Cols})
	}

	done := make(chan struct{})

	// Goroutine to read from PTY and write to WS
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := ptmx.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Println("PTY read error:", err)
				}
				break
			}
			
			// Send output as JSON
			msg := map[string]interface{}{
				"type": "output",
				"data": string(buf[:n]),
			}
			if err := ws.WriteJSON(msg); err != nil {
				break
			}
		}
		close(done)
	}()

	// Goroutine to periodically sync files back to frontend
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				files := make(map[string]string)
				filepath.Walk(tmpDir, func(path string, info os.FileInfo, err error) error {
					if err != nil || info.IsDir() {
						return nil
					}
					rel, _ := filepath.Rel(tmpDir, path)
					if rel == "go.mod" {
						return nil // ignore go.mod
					}
					content, err := os.ReadFile(path)
					if err == nil {
						files[rel] = string(content)
					}
					return nil
				})
				
				if len(files) > 0 {
					msg := map[string]interface{}{
						"type": "files",
						"files": files,
					}
					ws.WriteJSON(msg)
				}
			}
		}
	}()

	// Goroutine to read from WS and write to PTY
	go func() {
		for {
			_, msg, err := ws.ReadMessage()
			if err != nil {
				break
			}
			
			var wsm WSTerminalMessage
			if err := json.Unmarshal(msg, &wsm); err != nil {
				// Fallback: treat raw message as input if not JSON
				ptmx.Write(msg)
				continue
			}

			switch wsm.Type {
			case "input":
				ptmx.Write([]byte(wsm.Data))
			case "resize":
				pty.Setsize(ptmx, &pty.Winsize{Rows: wsm.Rows, Cols: wsm.Cols})
			}
		}
		// If WS closes, kill the process
		cmd.Process.Kill()
	}()

	// Wait for process to exit or read error
	select {
	case <-done:
	case <-time.After(2 * time.Hour): // Timeout to prevent infinite sessions
		cmd.Process.Kill()
	}
	
	cmd.Wait()
}
