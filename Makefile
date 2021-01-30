generate:
	go get \
		google.golang.org/grpc \
		google.golang.org/protobuf/cmd/protoc-gen-go \
        google.golang.org/grpc/cmd/protoc-gen-go-grpc
	mkdir -p internal/pb
	protoc --proto_path=api/ --go_out=internal/pb --go-grpc_out=internal/pb api/*.proto

run:
	go run cmd/sysmon/main.go