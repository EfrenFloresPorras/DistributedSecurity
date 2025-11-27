PROTO_FILES := $(wildcard *.proto)

gen:
	protoc --go_out=../.. --go-grpc_out=../.. $(PROTO_FILES)

clean:
	rm -rf ../*/proto/*pb
