# 使用官方的Go镜像作为基础镜像
FROM golang:latest

LABEL maintainer="Sunwish <isunwish@gmail.com>"

# 配置构建linux amd64
ENV GOARCH=amd64
ENV GOHOSTARCH=amd64
ENV GOOS=linux
ENV GOHOSTOS=linux

# 设置工作目录
WORKDIR /build

# 复制Go程序源代码到工作目录
COPY . /build

# 构建Go程序
RUN go build -o main

#######################

FROM ubuntu

# 设置工作目录
WORKDIR /app

# 使用上层构建好的程序
COPY --from=0 /build/main /app/main

# 存储本地证书到镜像中
ADD cacert.pem /etc/ssl/certs/

# 挂载点
VOLUME ["/app/archived"]

# 入口程序
CMD [ "./main" ]
