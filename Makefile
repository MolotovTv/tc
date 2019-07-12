PROTOC := protoc \
			--proto_path=github.com/gogo/protobuf/protobuf/:. \
			--gogo_out=plugins=grpc:.

proto:
	cd $(GOPATH)/src && $(PROTOC) github.com/molotovtv/tc/tc/*.proto
	cd $(GOPATH)/src && $(PROTOC) github.com/molotovtv/tc/server/*.proto

GIT_COMMIT := $(shell git rev-list -1 HEAD)

build:
	go build -ldflags "-X github.com/molotovtv/tc/cmd.gitCommit=$(GIT_COMMIT)" .

install:
	go install -ldflags "-X github.com/molotovtv/tc/cmd.gitCommit=$(GIT_COMMIT)" .