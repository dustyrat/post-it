for OS in "linux" "windows" "darwin" "freebsd"; do
    for ARCH in "386" "amd64"; do
        CGO_ENABLED=0 GOOS=$OS GOARCH=$ARCH go build -o ${OS}_${ARCH}/post-it ./main.go
    done
done