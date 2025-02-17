build:
	cd ./cmd/agent && go build -o agent main.go && cd ../server && go build -o server main.go
