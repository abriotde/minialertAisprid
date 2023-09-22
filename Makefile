
PROTOC_CMD = protoc
PROTOC_GRPC_CMD = protoc --plugin=$(shell go env GOPATH)/bin/protoc-gen-go-grpc

all: minialertAisprid

minialertAisprid: messages/setIntVar_grpc.pb.go messages/setIntVar.pb.go
	go build

messages/setIntVar_grpc.pb.go: messages/setIntVar.proto
	$(PROTOC_GRPC_CMD) --go-grpc_out=. --go-grpc_opt=paths=source_relative  messages/setIntVar.proto

messages/setIntVar.pb.go: messages/setIntVar.proto
	$(PROTOC_CMD) --go_out=. --go_opt=paths=source_relative messages/setIntVar.proto

dependencies:
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

clean:
	rm -f minialertAisprid
	rm -f messages/*.go

