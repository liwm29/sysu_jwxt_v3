FROM golang

workdir /go/app
copy . .
run export GOPROXY=https://mirrors.aliyun.com/goproxy/
run GO111MODULE=on
run go mod tidy

run go build -o jwxt ./cmd


cmd ["./jwxt"]