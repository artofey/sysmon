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

run:
	go run cmd/sysmon/main.go

run_client:
	go run cmd/sysmon/client/main.go
