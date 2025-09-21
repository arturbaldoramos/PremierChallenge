# Comandos da API
.PHONY: run build clean deps migrate migrate-reset migrate-status dev test

# Executar a API
run:
	go run cmd/api/main.go

# Instalar dependências
deps:
	go mod tidy
	go mod download

# Limpar arquivos de build
clean:
	rm -rf bin/

# Executar com live reload (requer air)
dev:
	air

# ====== COMANDOS DE MIGRAÇÃO ======

# Executar migrações do banco de dados
migrate:
	go run cmd/api/main.go -migrate

# Resetar banco de dados (CUIDADO: apaga todos os dados)
migrate-reset:
	@echo "ATENÇÃO: Este comando apagará todos os dados do banco!"
	@echo "Pressione Ctrl+C para cancelar ou Enter para continuar..."
	@read confirm
	go run cmd/api/main.go -reset

# Verificar status das migrações
migrate-status:
	go run cmd/api/main.go -status

# ====== COMANDOS DE DESENVOLVIMENTO ======

# Executar testes
test:
	go test ./...

# Executar testes com cobertura
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Build da aplicação
build:
	go build -o bin/api cmd/api/main.go

# ====== COMANDOS DOCKER ======

# Subir containers do banco
docker-up:
	docker-compose up -d

# Parar containers
docker-down:
	docker-compose down

# Resetar containers (remove dados)
docker-reset:
	docker-compose down -v
	docker-compose up -d

# ====== UTILITÁRIOS ======

# Verificar sintaxe Go
fmt:
	go fmt ./...

# Verificar problemas no código
vet:
	go vet ./...

# Executar linter (requer golangci-lint)
lint:
	golangci-lint run

# Setup completo do projeto
setup: deps docker-up migrate
	@echo "Setup completo! Use 'make run' para iniciar a API"