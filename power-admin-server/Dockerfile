# 构建阶段
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安装依赖
RUN apk add --no-cache git make

# 复制go.mod和go.sum
COPY power-admin-server/go.mod .
COPY power-admin-server/go.sum .

# 下载依赖
RUN go mod download

# 复制源代码
COPY power-admin-server/ .

# 编译
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o power-admin .

# 运行阶段
FROM alpine:latest

WORKDIR /app

# 安装ca证书和时区数据
RUN apk --no-cache add ca-certificates tzdata

# 从构建阶段复制二进制文件
COPY --from=builder /app/power-admin .
COPY --from=builder /app/etc/ ./etc/
COPY --from=builder /app/db/ ./db/

# 暴露端口
EXPOSE 8888

# 启动应用
CMD ["./power-admin", "-f", "etc/power-api.yaml"]
