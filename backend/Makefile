.PHONY: build run restart stop clean help

BINARY := bin/api
MAIN   := ./cmd/api
PORT   := 8080

# Carrega variáveis do .env se existir
ifneq (,$(wildcard .env))
    include .env
    export
endif

help:
	@echo "Comandos disponíveis:"
	@echo "  make build   - Compila a API (binário em bin/api)"
	@echo "  make run     - Sobe a API (faz build se necessário)"
	@echo "  make restart - Para a API na porta $(PORT) e sobe de novo"
	@echo "  make stop    - Para o processo da API na porta $(PORT)"
	@echo "  make clean   - Remove o binário gerado"

build:
	@echo "Building..."
	@mkdir -p bin
	@go build -o $(BINARY) $(MAIN)
	@echo "Build concluído: $(BINARY)"

run: build
	@echo "Iniciando API na porta $(PORT)..."
	@./$(BINARY)

restart: stop build
	@echo "Iniciando API na porta $(PORT)..."
	@./$(BINARY)

stop:
	@-lsof -ti:$(PORT) | xargs kill -9 2>/dev/null || true
	@echo "Processo na porta $(PORT) encerrado (se existia)."

clean:
	@rm -rf bin
	@echo "Limpeza feita."
