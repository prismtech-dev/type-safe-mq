OUT_DIR=tests
PROTO=mock.proto

.PHONY: python go all clean

python:
	python -m grpc_tools.protoc \
		-I. \
		--python_out=python/$(OUT_DIR) \
		--plugin=protoc-gen-mypy=$(VIRTUAL_ENV)/bin/protoc-gen-mypy \
		--mypy_out=python/$(OUT_DIR) \
		$(PROTO)

go:
	protoc -I. \
		--go_out=go/$(OUT_DIR)/pb \
		--go_opt=paths=source_relative \
		${PROTO}

all: python go

clean:
	rm -f python/$(OUT_DIR)/mock_pb2.py
	rm -f go/$(OUT_DIR)/mock.pb.go