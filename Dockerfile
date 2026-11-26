# ใช้ Image ของ Python แทน Go
FROM python:3.9-slim

# ตั้งค่า Working Directory
WORKDIR /app

# คัดลอกไฟล์ server.py เข้าไปใน Container
COPY server.py .

# เปิด Port 8080
EXPOSE 8080

# คำสั่งรันโปรแกรม
CMD ["python", "server.py"]
