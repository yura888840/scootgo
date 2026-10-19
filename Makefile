APP_NAME=scootgo
BUILD_DIR=build
SRC=./pkg
GO=go
BATCH_SIZE?=1
TEST_PACKAGES := $(shell $(GO) list -tags=mysql ./... | grep -vE '/pkg/(ridetest)(/|$$)')

-include .env.secrets

# Build the binary
.PHONY: build
build:
	@echo "Building $(APP_NAME)..."
	$(GO) build -tags mysql -o $(BUILD_DIR)/$(APP_NAME) .

# Run the application
.PHONY: run-http
run-http:
	@echo "Running http server in $(APP_NAME)..."
	$(BUILD_DIR)/$(APP_NAME) serve

# Run rabbitMQ consumer
.PHONY: run-rmq-c
run-rmq-c:
ifeq ($(strip $(QUEUE_NAME)),)
	$(error QUEUE_NAME is mandatory. Use 'make run-rmq-c QUEUE_NAME=value ROUTING_KEY=value')
endif
ifeq ($(strip $(ROUTING_KEY)),)
	$(error ROUTING_KEY is mandatory. Use 'make run-rmq-c QUEUE_NAME=value ROUTING_KEY=value')
endif
	@echo "Running rmq consumer in $(APP_NAME) for queue $(QUEUE_NAME), routing key $(ROUTING_KEY) with batch size $(BATCH_SIZE)"
	$(BUILD_DIR)/$(APP_NAME) consume $(QUEUE_NAME) $(ROUTING_KEY) $(BATCH_SIZE)

.PHONY: run-command
run-command:
	@echo "Running sample command in $(APP_NAME)..."
	$(BUILD_DIR)/$(APP_NAME) exec hello

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	$(GO) test -tags=mysql $(TEST_PACKAGES) -v -skip "^TestFunctional"

# Format the code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

# Clean up build artifacts
.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -rf $(BUILD_DIR)

# Lint the code
.PHONY: lint
lint:
	@echo "Linting code..."
	go tool golangci-lint run --new-from-rev=origin/main --build-tags=mysql
