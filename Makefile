.PHONY: help protoc

all:

protoc:
	rm -rf go/lib/helloworld
	protoc --go_out=./go --go-grpc_out=./go proto/helloworld.proto	

clean:
	rm -rf go/helloworld