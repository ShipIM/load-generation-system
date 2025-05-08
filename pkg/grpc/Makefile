.PHONY: all go-generate

all:

go-generate:
	rm -rf go/
	mkdir go
	
	protoc --go_out=go --go_opt=Mload_generation_system_v1.proto=./pb --go-grpc_out=go --go-grpc_opt=Mload_generation_system_v1.proto=./pb load_generation_system_v1.proto