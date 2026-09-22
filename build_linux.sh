CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -extldflags -static -extldflags -static"  -o nps_linux_amd64 ./cmd/nps/nps.go
chmod +x nps_linux_amd64
