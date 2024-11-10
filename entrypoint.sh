export CGO_ENABLED=1
go mod download
go build -ldflags="-s -w" -o main .
./main