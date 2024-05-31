to build run:
- `go get -C . -d ./...`
- `go mod tidy `
- `go build -C . -o ./bin -ldflags "-s -w"`

or use `task build-r` if you have task installed
