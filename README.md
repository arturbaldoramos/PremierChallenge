# Premiere Challenge

Sistema de gerenciamento hospitalar com API em Go e banco PostgreSQL 17.

## 🚀 Como executar

### Pré-requisitos
- Docker e Docker Compose
- Go 1.25+ (para desenvolvimento)

### 1. Iniciar o banco de dados

```bash
# Subir o PostgreSQL com Docker Compose
docker-compose up -d postgres

# Verificar se está rodando
docker-compose ps
```

### 2. Executar a API

```bash
cd Backend

# Instalar dependências
go mod tidy

# Executar migrações
make migrate-up
# ou
go run cmd/migrate/main.go -up

# Iniciar a API
make run
# ou
go run cmd/api/main.go
```

### 3. Verificar funcionamento

```bash
# Health check da API
curl http://localhost:8080/health
```

## 🐘 Gerenciar PostgreSQL

### Comandos Docker Compose

```bash
# Iniciar apenas o PostgreSQL
docker-compose up -d postgres

# Iniciar PostgreSQL + PgAdmin
docker-compose up -d

# Parar serviços
docker-compose down

# Remover volumes (dados)
docker-compose down -v
```

### Acessar PgAdmin (Interface Web)
- URL: http://localhost:5050
- Email: admin@premiere.com
- Senha: admin123

**Configuração do servidor no PgAdmin:**
- Host: postgres
- Port: 5432
- Database: premieredb
- Username: postgres
- Password: postgres123

### Conectar diretamente ao PostgreSQL

```bash
# Via psql
docker exec -it premiere_postgres psql -U postgres -d premieredb

# Via Docker Compose
docker-compose exec postgres psql -U postgres -d premieredb
```

## 🛠️ Comandos úteis

```bash
cd Backend

# Executar migrações
make migrate-up

# Reverter migrações
make migrate-down

# Executar API
make run

# Build da aplicação
make build

# Executar testes
make test

# Instalar dependências
make deps
```

## 📦 Estrutura do projeto

```
├── Backend/
│   ├── cmd/
│   │   ├── api/main.go          # Servidor HTTP
│   │   └── migrate/main.go      # CLI de migrações
│   ├── internal/
│   │   ├── config/              # Configurações (Viper)
│   │   ├── database/            # Migrações GORM
│   │   └── domain/              # Models
│   ├── Makefile                 # Comandos simplificados
│   └── go.mod
├── docker-compose.yml           # PostgreSQL + PgAdmin
└── init.sql                     # Inicialização do banco
```

## 🔧 Configuração

### Via arquivo .env
```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres123
DB_NAME=premieredb
```

### Via config.yaml
```yaml
database:
  host: localhost
  port: "5432"
  user: postgres
  password: postgres123
  name: premieredb
```