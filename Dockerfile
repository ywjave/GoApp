FROM golang:1.21-alpine AS builder

#为镜像设置环境变量
ENV GO111MODULE=on \
    CGO_ENABLE=0   \
    GOPROXY=goproxy.cn,direct \
    GOOS=linux   \
    GOARCH=amd64

#移动到工作目录
WORKDIR /build

#复制项目中的 go.mod go.sum并下载依赖信
COPY go.mod .
COPY go.sum .

RUN go mod download

#   将代码复制到容器中

COPY . .

#编译成二进制可执行文件 bullbell
RUN go build -o bluebell .

#创建小镜像
FROM scratch

COPY ./conf /conf

COPY --from=builder /build/bluebell  ./


#RUN ls

EXPOSE 8086

ENTRYPOINT ["/bluebell"]

