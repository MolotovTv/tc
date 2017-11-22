PROTOC := protoc \
			--proto_path=github.com/gogo/protobuf/protobuf/:. \
			--gogo_out=plugins=grpc:.

proto:
	cd $(GOPATH)/src && $(PROTOC) github.com/aestek/tc/tc/*.proto
	cd $(GOPATH)/src && $(PROTOC) github.com/aestek/tc/server/*.proto
