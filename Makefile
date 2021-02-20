
CLIENT_DOCKER_IMG="sysmon-client:develop"
SERVER_DOCKER_IMG="sysmon-server:develop"

generate:
	go get \
		google.golang.org/grpc \
		google.golang.org/protobuf/cmd/protoc-gen-go \
        google.golang.org/grpc/cmd/protoc-gen-go-grpc
	mkdir -p pkg/server/pb
	protoc --proto_path=api/ --go_out=pkg/server/pb --go-grpc_out=pkg/server/pb api/*.proto

lint:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run --fix
	golangci-lint run ./...

build:
	mkdir .build
	go build -o .build/ ./cmd/...

build-img-client:
	docker -v build -t $(CLIENT_DOCKER_IMG) -f cmd/client/Dockerfile .

build-img-server:
	docker -v build -t $(SERVER_DOCKER_IMG) -f cmd/server/Dockerfile .

run-img-client:
	docker run -it --rm $(CLIENT_DOCKER_IMG)

run-img-server:
	docker run -it --rm $(SERVER_DOCKER_IMG)

run:
	go run cmd/server/main.go

run_client:
	go run cmd/client/main.go
