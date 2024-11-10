export CGO_ENABLED=1
go mod down
go build -ldflags="-s -w" -o main .
./main