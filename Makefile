# Ensure paths are correct
SERVER_PATH=server/cmd/server
CLIENT_PATH=client/cmd

build-linux-server:
	GOOS=linux GOARCH=amd64 go build -o $(SERVER_PATH)/myserver-linux $(SERVER_PATH)/main.go

build-macos-server:
	GOOS=darwin GOARCH=amd64 go build -o $(SERVER_PATH)/myserver-macos $(SERVER_PATH)/main.go

build-windows-server:
	GOOS=windows GOARCH=amd64 go build -o $(SERVER_PATH)/myserver.exe $(SERVER_PATH)/main.go

build-linux-client:
	GOOS=linux GOARCH=amd64 go build -o $(CLIENT_PATH)/mycli-linux $(CLIENT_PATH)/main.go

build-macos-client:
	GOOS=darwin GOARCH=amd64 go build -o $(CLIENT_PATH)/mycli-macos $(CLIENT_PATH)/main.go

build-windows-client:
	GOOS=windows GOARCH=amd64 go build -o $(CLIENT_PATH)/mycli.exe $(CLIENT_PATH)/main.go

build-all: build-linux-server build-macos-server build-windows-server build-linux-client build-macos-client build-windows-client