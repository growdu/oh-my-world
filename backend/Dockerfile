FROM golang:1.24.1

WORKDIR /app

COPY go.mod go.sum ./
ENV GOPROXY=https://goproxy.cn,direct

RUN go mod tidy

COPY . .

RUN go build -o main .

CMD ["./main"]