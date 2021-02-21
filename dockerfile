from golang


run sed -i s@/deb.debian.org/@/mirrors.aliyun.com/@g /etc/apt/sources.list 
run sed -i s@/archive.ubuntu.com/@/mirrors.aliyun.com/@g /etc/apt/sources.list

run apt-get clean 
run apt-get update 

run apt-get install vim -y

# run apt-get upgrade -y
# run apt-get install net-tools

# note: the image layer will be built and removed one by one
# note: so, the shell will be run one by one,repeatedly lose its config

run echo "export GOPROXY=https://goproxy.io/" >> ~/.bashrc 

workdir /go/src/github.com/liwm29/sysu_jwxt_v3
copy . .

run export GOPROXY=https://goproxy.io/ && go mod tidy

run go build -o ./jwxt ./cmd && cp ./jwxt /go/bin

cmd ["/bin/bash"]