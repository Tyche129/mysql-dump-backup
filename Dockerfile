# 阶段 1: 构建
FROM golang:1.25.5-alpine AS builder

WORKDIR /app

# 安装必要的构建工具
RUN apk add --no-cache git build-base

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    GOTOOLCHAIN=auto \
    CGO_ENABLED=0  # 建议显式关闭 CGO，以确保生成静态链接的二进制文件

COPY go.mod go.sum ./
RUN go mod download

# 将源码复制到构建镜像中
COPY . .
# 编译
RUN go build -ldflags="-s -w" -o mysql-dump-backup .

# 阶段 2: 运行
FROM alpine:latest

WORKDIR /app

# 生产环境通常需要 ca-certificates 以便发起 HTTPS 请求
RUN apk --no-cache add ca-certificates

# 从构建阶段复制二进制文件
COPY --from=builder /app/mysql-dump-backup .

# 不需要 EXPOSE，因为你的服务不监听端口

# 启动命令
CMD ["./mysql-dump-backup"]