PROTO_FILE=balancer.proto
OUT_DIR=./gen
BINARY_DIR=./bin

EP1=$(BINARY_DIR)/endpoint1
EP2=$(BINARY_DIR)/endpoint2
EP3=$(BINARY_DIR)/endpoint3
MAIN_EXEC=$(BINARY_DIR)/main_exec

.PHONY: all clean build_proto build_endpoints run_endpoints shutdown

# Default to build everything
all: clean build_proto run_endpoints

# Build protobuf files
build_proto: $(OUT_DIR)/balancer.pb.go

$(OUT_DIR)/balancer.pb.go: $(PROTO_FILE)
	@mkdir -p $(OUT_DIR)
	protoc --go_out=$(OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative $(PROTO_FILE)

$(EP1): 
	@echo "Building endpoint 1..."
	@mkdir -p $(BINARY_DIR)
	@go build -o $(EP1) /Users/rase/dev/LoadBalancer/endpoints/endpoint1/main.go

$(EP2): 
	@echo "Building endpoint 2..."
	@mkdir -p $(BINARY_DIR)
	@go build -o $(EP2) /Users/rase/dev/LoadBalancer/endpoints/endpoint2/main.go

$(EP3): 
	@echo "Building endpoint 3..."
	@mkdir -p $(BINARY_DIR)
	@go build -o $(EP3) /Users/rase/dev/LoadBalancer/endpoints/endpoint3/main.go

$(MAIN_EXEC):
	@echo "Building load balancer executable..."
	@mkdir -p $(BINARY_DIR)
	@go build -o $(MAIN_EXEC) /Users/rase/dev/LoadBalancer/main.go

# Run endpoints
run_endpoints: $(EP1) $(EP2) $(EP3) $(MAIN_EXEC) 
	@echo "Running endpoint 1..."
	@$(EP1) &
	@sleep 2
	@echo "Running endpoint 2..."
	@$(EP2) &
	@sleep 2
	@echo "Running endpoint 3..."
	@$(EP3) &
	@sleep 2
	@echo "Running main executable..."
	@$(MAIN_EXEC)

# Clean directories
clean:
	@rm -rf $(OUT_DIR)
	@rm -rf $(BINARY_DIR)
	@echo "Clearing port 8081..."
	-@kill -9 $$(lsof -t -i:8081) 2>/dev/null || true
	@echo "Clearing port 8082..."
	-@kill -9 $$(lsof -t -i:8082) 2>/dev/null || true
	@echo "Clearing port 8083..."
	-@kill -9 $$(lsof -t -i:8083) 2>/dev/null || true
	@echo "Clearing port 50051..."
	-@kill -9 $$(lsof -t -i:50051) 2>/dev/null || true