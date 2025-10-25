# 构建阶段
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的构建工具
RUN apk add --no-cache git

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 运行阶段
FROM alpine:latest

# 安装 ca-certificates 用于 HTTPS 请求
RUN apk --no-cache add ca-certificates tzdata

# 设置时区为上海
ENV TZ=Asia/Shanghai

WORKDIR /root/

# 从构建阶段复制可执行文件
COPY --from=builder /app/main .

# 复制 .env 文件（如果存在）
COPY --from=builder /app/.env* ./

# 暴露端口（如果需要）
# EXPOSE 8080

# 运行应用
CMD ["./main", "-daemon"]

