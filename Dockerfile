# ใช้ Image ที่รวม Go ไว้เลย
FROM golang:1.21-alpine AS builder

# สร้างโฟลเดอร์ทำงาน
WORKDIR /app

# คัดลอกไฟล์ go.mod และ go.sum ก่อน (ถ้ามี) เพื่อ cache layers
COPY go.mod go.sum ./
# ถ้าไม่มีไฟล์เหล่านี้ ให้สร้างใหม่ด้วยคำสั่ง: go mod init c2-client && go mod tidy

# ติดตั้ง Dependencies ทั้งหมดก่อน
RUN go mod download

# คัดลอกโค้ดทั้งหมด
COPY . .

# Build โปรแกรม
# -ldflags="-s -w" คือลดขนาดไฟล์ binary
RUN go build -ldflags="-s -w" -o client client.go

# สร้าง Image สั้นๆ สำหรับรัน
FROM alpine:latest
WORKDIR /root
COPY --from=builder /app/client /client

# รันโปรแกรม
CMD ["/client"]
