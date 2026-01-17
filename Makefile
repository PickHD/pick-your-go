# Pick Your Go - Makefile
# A CLI tool to generate Go projects with various architecture patterns

# Application variables
APP_NAME=pick-your-go
CMD_DIR=./cmd/$(APP_NAME)
BINARY_NAME=$(APP_NAME)
BINARY_PATH=./bin/$(BINARY_NAME)

# Go variables
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOINSTALL=$(GOCMD) install

# Build variables
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildDate=$(DATE)"
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Colors for output
COLOR_RESET=\033[0m
COLOR_BOLD=\033[1m
COLOR_GREEN=\033[32m
COLOR_YELLOW=\033[33m
COLOR_BLUE=\033[34m

.PHONY: help
help: ## Display this help screen
	@echo "$(COLOR_BOLD)Pick Your Go - Available Commands:$(COLOR_RESET)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(COLOR_BLUE)%-20s$(COLOR_RESET) %s\n", $$1, $$2}'

.PHONY: all
all: clean deps lint test build ## Run all development tasks (clean, deps, lint, test, build)

.PHONY: build
build: ## Build the application binary
	@echo "$(COLOR_GREEN)Building $(APP_NAME)...$(COLOR_RESET)"
	@mkdir -p bin
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_PATH) $(CMD_DIR)
	@echo "$(COLOR_GREEN)Build complete: $(BINARY_PATH)$(COLOR_RESET)"

.PHONY: run
run: ## Run the application directly
	@echo "$(COLOR_GREEN)Running $(APP_NAME)...$(COLOR_RESET)"
	$(GORUN) $(CMD_DIR)

.PHONY: install
install: ## Install the application to $GOPATH/bin
	@echo "$(COLOR_GREEN)Installing $(APP_NAME) to \$$(go env GOPATH)/bin...$(COLOR_RESET)"
	$(GOINSTALL) $(LDFLAGS) $(CMD_DIR)
	@echo "$(COLOR_GREEN)Installation complete!$(COLOR_RESET)"
	@echo "$(COLOR_YELLOW)Make sure \$$(go env GOPATH)/bin is in your PATH$(COLOR_RESET)"

.PHONY: install-system
install-system: build ## Install the application to /usr/local/bin (requires sudo)
	@echo "$(COLOR_GREEN)Installing $(APP_NAME) to /usr/local/bin...$(COLOR_RESET)"
	@sudo cp $(BINARY_PATH) /usr/local/bin/$(BINARY_NAME)
	@sudo chmod +x /usr/local/bin/$(BINARY_NAME)
	@echo "$(COLOR_GREEN)Installation complete!$(COLOR_RESET)"
	@echo "$(COLOR_BLUE)Run '$(BINARY_NAME) --help' to get started$(COLOR_RESET)"

.PHONY: uninstall
uninstall: ## Uninstall the application
	@echo "$(COLOR_YELLOW)Uninstalling $(APP_NAME)...$(COLOR_RESET)"
	@rm -f $$(go env GOPATH)/bin/$(BINARY_NAME) 2>/dev/null || true
	@sudo rm -f /usr/local/bin/$(BINARY_NAME) 2>/dev/null || true
	@echo "$(COLOR_GREEN)Uninstallation complete!$(COLOR_RESET)"

.PHONY: clean
clean: ## Clean build artifacts and cache
	@echo "$(COLOR_YELLOW)Cleaning...$(COLOR_RESET)"
	@rm -rf bin/
	@rm -rf dist/
	@echo "$(COLOR_GREEN)Clean complete!$(COLOR_RESET)"

.PHONY: deps
deps: ## Download dependencies
	@echo "$(COLOR_GREEN)Downloading dependencies...$(COLOR_RESET)"
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "$(COLOR_GREEN)Dependencies ready!$(COLOR_RESET)"

.PHONY: test
test: ## Run tests
	@echo "$(COLOR_GREEN)Running tests...$(COLOR_RESET)"
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	@echo "$(COLOR_GREEN)Tests complete!$(COLOR_RESET)"

.PHONY: test-coverage
test-coverage: test ## Run tests and display coverage
	@echo "$(COLOR_GREEN)Coverage report:$(COLOR_RESET)"
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "$(COLOR_GREEN)Coverage report generated: coverage.html$(COLOR_RESET)"

.PHONY: lint
lint: ## Run linters
	@echo "$(COLOR_GREEN)Running linters...$(COLOR_RESET)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "$(COLOR_YELLOW)golangci-lint not found. Skipping...$(COLOR_RESET)"; \
		echo "$(COLOR_YELLOW)Install it from: https://golangci-lint.run/usage/install/$(COLOR_RESET)"; \
	fi
	@if command -v gofmt >/dev/null 2>&1; then \
		echo "$(COLOR_GREEN)Checking formatting...$(COLOR_RESET)"; \
		OUTPUT=$$(gofmt -l .); \
		if [ -n "$$OUTPUT" ]; then \
			echo "$(COLOR_YELLOW)The following files are not formatted:$(COLOR_RESET)"; \
			echo "$$OUTPUT"; \
			exit 1; \
		fi; \
	fi
	@echo "$(COLOR_GREEN)Linting complete!$(COLOR_RESET)"

.PHONY: fmt
fmt: ## Format code
	@echo "$(COLOR_GREEN)Formatting code...$(COLOR_RESET)"
	gofmt -w -s .
	@echo "$(COLOR_GREEN)Formatting complete!$(COLOR_RESET)"

.PHONY: vet
vet: ## Run go vet
	@echo "$(COLOR_GREEN)Running go vet...$(COLOR_RESET)"
	$(GOCMD) vet ./...
	@echo "$(COLOR_GREEN)Vet complete!$(COLOR_RESET)"

.PHONY: init-example
init-example: ## Run init command with example values
	@echo "$(COLOR_GREEN)Running init command...$(COLOR_RESET)"
	$(GORUN) $(CMD_DIR) init --architecture layered --name example-app --module github.com/example/example-app

.PHONY: templates-update
templates-update: ## Update template cache
	@echo "$(COLOR_GREEN)Updating templates...$(COLOR_RESET)"
	$(GORUN) $(CMD_DIR) templates update

.PHONY: templates-list
templates-list: ## List available templates
	@echo "$(COLOR_GREEN)Listing templates...$(COLOR_RESET)"
	$(GORUN) $(CMD_DIR) templates list

.PHONY: dev
dev: ## Run in development mode with hot reload (requires air)
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "$(COLOR_YELLOW)air not found. Installing...$(COLOR_RESET)"; \
		$(GOGET) -u github.com/cosmtrek/air@latest; \
		air; \
	fi

.PHONY: build-all
build-all: ## Build binaries for multiple platforms
	@echo "$(COLOR_GREEN)Building for multiple platforms...$(COLOR_RESET)"
	@mkdir -p dist
	@GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o dist/$(APP_NAME)-linux-amd64 $(CMD_DIR)
	@GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o dist/$(APP_NAME)-darwin-amd64 $(CMD_DIR)
	@GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o dist/$(APP_NAME)-darwin-arm64 $(CMD_DIR)
	@GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o dist/$(APP_NAME)-windows-amd64.exe $(CMD_DIR)
	@echo "$(COLOR_GREEN)Builds complete in dist/$(COLOR_RESET)"

.PHONY: release
release: clean lint test build-all ## Create a release build
	@echo "$(COLOR_GREEN)Release builds complete!$(COLOR_RESET)"

.PHONY: check-deps
check-deps: ## Check if required tools are installed
	@echo "$(COLOR_GREEN)Checking dependencies...$(COLOR_RESET)"
	@command -v git >/dev/null 2>&1 || { echo "$(COLOR_YELLOW)git not found$(COLOR_RESET)"; exit 1; }
	@command -v go >/dev/null 2>&1 || { echo "$(COLOR_YELLOW)go not found$(COLOR_RESET)"; exit 1; }
	@echo "$(COLOR_GREEN)All dependencies installed!$(COLOR_RESET)"

.PHONY: setup
setup: check-deps deps ## Setup development environment
	@echo "$(COLOR_GREEN)Setting up development environment...$(COLOR_RESET)"
	@echo "$(COLOR_GREEN)Setup complete! Run 'make help' for available commands.$(COLOR_RESET)"
