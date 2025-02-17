build:
	cd ./cmd/agent && go build -o agent main.go && cd ../server && go build -o server main.go

run-a:
	go run ./cmd/agent/main.go

run-s:
	go run ./cmd/server/main.go
