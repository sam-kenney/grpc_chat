build-proto:
	protoc --go-grpc_out=. --go_out=. proto/service.proto
