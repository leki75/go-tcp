GOPATH := $(shell go env GOPATH)

.PHONY=generate
generate: $(GOPATH)/bin/protoc-gen-go-vtproto
	protoc \
		--plugin protoc-gen-go-vtproto="$(GOPATH)/bin/protoc-gen-go-vtproto" \
		--go_out=. \
		--go-grpc_out=. \
		--go-vtproto_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		--go-vtproto_opt=paths=source_relative \
		--go-vtproto_opt=features=marshal+unmarshal+size \
		schema/proto/stocks.proto

$(GOPATH)/bin/protoc-gen-go-vtproto:
	go install github.com/planetscale/vtprotobuf/cmd/protoc-gen-go-vtproto@v0.2.0
