
.DEFAULT_GOAL := help

help:
	@echo "------------------------------------------------------------------"
	@echo " Makefile for curta"
	@echo "------------------------------------------------------------------"
	@echo " > make help   # show this help info"
	@echo " > make build  # build "
	@echo " > make proto  # regenerate artifacts from proto"
	@echo " > make test   # run all the go tests"
	@echo ""

proto-lint:
	@echo "Running linter on proto.  No output is good."
	protoc --lint_out=./ ./proto/mux.proto

.PHONY: proto
proto:
	protoc --go_out=plugins=grpc:. ./proto/mux.proto

build:
	GOBIN=$(PWD)/bin \
		  go install ./...

test:
	go test ./pkg/mux/

clean:
	rm -rf ./bin/ ./cmd/muxd/muxd

