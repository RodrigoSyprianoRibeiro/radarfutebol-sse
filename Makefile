.PHONY: build run dev clean test

# Nome do binario
BINARY_NAME=radarfutebol-sse

# Diretorio de build
BUILD_DIR=bin

# Build do projeto
build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/main.go

# Build para producao (otimizado)
build-prod:
	@echo "Building for production..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/main.go

# Executa o servidor
run: build
	@echo "Running..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Executa em modo desenvolvimento (com hot reload via air se instalado)
dev:
	@if command -v air > /dev/null; then \
		air; \
	else \
		go run ./cmd/main.go; \
	fi

# Baixa dependencias
deps:
	go mod download
	go mod tidy

# Limpa arquivos de build
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)

# Executa testes
test:
	go test -v ./...

# Verifica o codigo
lint:
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		go vet ./...; \
	fi

# Mostra ajuda
help:
	@echo "Comandos disponiveis:"
	@echo "  make build      - Compila o projeto"
	@echo "  make build-prod - Compila para producao (otimizado)"
	@echo "  make run        - Compila e executa"
	@echo "  make dev        - Executa em modo desenvolvimento"
	@echo "  make deps       - Baixa dependencias"
	@echo "  make clean      - Remove arquivos de build"
	@echo "  make test       - Executa testes"
	@echo "  make lint       - Verifica codigo"
