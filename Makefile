PROTOC := protoc \
			--proto_path=github.com/gogo/protobuf/protobuf/:. \
			--gogo_out=plugins=grpc:.
.PHONY: proto
proto:
	go get -u github.com/gogo/protobuf/proto
	go get -u github.com/gogo/protobuf/protoc-gen-gogo
	go get -u github.com/gogo/protobuf/gogoproto
	cd $(GOPATH)/src && $(PROTOC) github.com/molotovtv/tc/tc/*.proto

GIT_COMMIT := $(shell git rev-list -1 HEAD)

.PHONY: build
build:
	go build -ldflags "-X github.com/molotovtv/tc/cmd.gitCommit=$(GIT_COMMIT)" .

.PHONY: install
install:
	go install -ldflags "-X github.com/molotovtv/tc/cmd.gitCommit=$(GIT_COMMIT)"