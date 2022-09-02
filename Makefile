# brew install FiloSottile/musl-cross/musl-cross
# brew install mingw-w64
CGO_ENABLED=1 GOOS=linux  GOARCH=amd64  CC=x86_64-linux-musl-gcc  CXX=x86_64-linux-musl-g++ go build -o main
# 服务端需要安装 musl-devmusl-tools ，否则无法识别二进制，提示文件不存在