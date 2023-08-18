# Assuming you have the protoc and the Golang protoc-gen-go plugins installed.
PROTO_FILE=balancer.proto
OUT_DIR=./gen

.PHONY: all clean

all: $(OUT_DIR)/balancer.pb.go

$(OUT_DIR)/balancer.pb.go: $(PROTO_FILE)
	@mkdir -p $(OUT_DIR)
	protoc --go_out=$(OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative $(PROTO_FILE)

clean:
	@rm -rf $(OUT_DIR)