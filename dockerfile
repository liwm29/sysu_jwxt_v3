from --platform=linux/amd64 golang:1.16.0-alpine as baseimage

# run sed -i s@/deb.debian.org/@/mirrors.aliyun.com/@g /etc/apt/sources.list 
# run sed -i s@/archive.ubuntu.com/@/mirrors.aliyun.com/@g /etc/apt/sources.list
# run apt-get clean 
# run apt-get update 
# run apt-get install vim -y
# run apt-get upgrade -y
# run apt-get install net-tools

run echo "export GOPROXY=https://goproxy.io/" >> ~/.bashrc 
workdir /app
copy . .
run export GOPROXY=https://goproxy.io/ && go mod tidy
run go install ./cmd/jwxt

from --platform=linux/amd64 golang:1.16.0-alpine
workdir /go/bin
copy --from=baseimage /go/bin .

cmd ["/bin/bash"]