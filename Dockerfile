FROM golang:alpine AS build-base
WORKDIR /go/src/kafka2mongo

COPY go.mod .
COPY go.sum .

RUN go env -w GOPROXY=https://goproxy.cn,direct && \
    go mod download

FROM build-base AS pre-build
COPY . .

RUN go build -o kafka2mongo
ENV TZ=Asia/Shanghai

RUN echo "http://mirrors.aliyun.com/alpine/v3.10/main/" > /etc/apk/repositories && \
    apk update && \
    apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

FROM alpine
WORKDIR /app/

ENV TZ=Asia/Shanghai \
    MONGODB_DSN=mongodb://localhost:27017 \
    KAFKA_ADDR=localhost:9092


COPY --from=pre-build /go/src/kafka2mongo/kafka2mongo .
COPY --from=pre-build /etc/timezone /etc/timezone
COPY --from=pre-build /etc/localtime /etc/localtime
RUN mkdir -p /app/storage && \
    apk add tzdata

CMD ["/app/kafka2mongo"]