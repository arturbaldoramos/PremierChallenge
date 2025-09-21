# Premiere Challenge - Sistema de Análise de Dados Médicos

## 📋 Visão Geral

Este projeto foi desenvolvido como parte do **Premier Challenge** e consiste em um sistema completo para análise e cruzamento de dados médicos, incluindo informações sobre médicos, hospitais, pacientes e especialidades médicas. O sistema permite o upload de arquivos em múltiplos formatos e fornece estatísticas detalhadas através de uma interface web moderna.

## 🏗️ Arquitetura do Sistema

### Backend (Go)
- **Framework**: Gorilla Mux para roteamento HTTP
- **Banco de Dados**: PostgreSQL com GORM como ORM
- **Cache**: Redis para otimização de performance
- **Processamento**: Sistema de workers para processamento assíncrono de arquivos
- **Monitoramento**: WebSocket para atualizações em tempo real

### Frontend (React + TypeScript)
- **Framework**: React 19 com TypeScript
- **UI Library**: Radix UI + Tailwind CSS
- **Roteamento**: React Router DOM
- **Gráficos**: Recharts para visualização de dados
- **Tema**: Sistema de tema claro/escuro

## 📊 Processo de Leitura e Cruzamento de Dados

### 1. Upload e Detecção de Arquivos
O sistema suporta múltiplos formatos de arquivo:
- **CSV**: Arquivos separados por vírgula
- **XML**: Documentos estruturados
- **Excel**: Planilhas (.xlsx)
- **Detecção Automática**: O sistema identifica automaticamente o formato do arquivo

### 2. Parsers Especializados
Cada formato possui um parser dedicado:

#### CSV Parser (`Backend/internal/parsers/csv_parser.go`)
- Detecta automaticamente o delimitador (vírgula, ponto-e-vírgula)
- Mapeia colunas para entidades do domínio
- Suporte a encoding UTF-8

#### XML Parser (`Backend/internal/parsers/xml_parser.go`)
- Processa documentos XML estruturados
- Extração de dados hierárquicos
- Validação de schema

#### Unified Parser (`Backend/internal/parsers/unified_parser.go`)
- Coordena o processamento entre diferentes parsers
- Normalização de dados
- Validação de integridade

### 3. Entidades do Domínio

#### Hospital (`Backend/internal/domain/Hospital.go`)
```go
type Hospital struct {
    ID             int       // Chave primária
    UUID           uuid.UUID // Identificador único
    Nome           string    // Nome do hospital
    CEP            string    // Código postal
    Especialidades string    // Lista de especialidades atendidas
    LeitosTotais   int       // Capacidade total de leitos
    CodMunicipio   string    // Código do município (IBGE)
    Bairro         string    // Bairro de localização
}
```

#### Médico (`Backend/internal/domain/Medico.go`)
```go
type Medico struct {
    UUID          uuid.UUID // Identificador único
    Nome          string    // Nome do médico
    Especialidade string    // Especialidade médica
    CodMunicipio  string    // Código do município (IBGE)
}
```

### 4. Processamento Assíncrono
- **Workers**: Sistema de workers para processamento em background
- **Chunks**: Divisão de arquivos grandes em partes menores
- **Redis**: Fila de tarefas e cache de resultados
- **Monitor**: Acompanhamento do progresso via WebSocket

### 5. Lógica de Negócio - Cruzamento de Dados

#### Atribuição Médicos-Hospitais
O sistema implementa uma lógica sofisticada para vincular médicos a hospitais baseada em:

1. **Proximidade Geográfica**: Raio máximo de 30km
2. **Especialidade Compatível**: Hospital deve atender a especialidade do médico
3. **Limite de Hospitais**: Médico pode atender até 3 hospitais simultaneamente

#### Cálculo de Distância
- Utiliza coordenadas geográficas baseadas no código de município (IBGE)
- Fórmula de distância euclidiana para cálculo aproximado
- Cache de resultados para otimização

## 🔄 Fluxo de Processamento de Dados

### 1. Upload
```
Usuário → Frontend → API Upload → Detecção de Formato → Parser Específico
```

### 2. Processamento
```
Parser → Validação → Normalização → Inserção no BD → Cache Redis → Notificação WebSocket
```

### 3. Análise
```
Requisição API → Consulta BD/Cache → Agregação → Estatísticas → Response JSON
```

## 📈 Estatísticas e Métricas Disponíveis

### Estatísticas Gerais
- Total de hospitais cadastrados (~1.500)
- Total de médicos cadastrados (~280.000)
- Total de estados (27)
- Total de municípios (~5.570)

### Distribuições por Estado
- Médicos por estado
- Hospitais por estado
- Ranking de estados por densidade médica

### Análise por Especialidade
- Hospitais que atendem cada especialidade
- Distribuição de médicos por especialidade
- Cobertura geográfica por especialidade

### Análise Geográfica
- Hospitais por município
- Densidade de atendimento médico
- Áreas com déficit de cobertura

### Distribuição de Médicos
- Médicos que atendem 1 hospital
- Médicos que atendem 2 hospitais
- Médicos que atendem 3 hospitais
- Médicos sem hospital vinculado

## 🚀 Como Executar o Projeto

### Pré-requisitos
- Docker e Docker Compose
- Go 1.25+ (para desenvolvimento)
- Node.js 18+ (para desenvolvimento)

### Execução com Docker
```bash
# Subir a infraestrutura (PostgreSQL + Redis)
cd Backend
docker-compose up -d

# Executar o backend
make run

# Executar o frontend (em outro terminal)
cd Frontend
npm install
npm run dev
```

### URLs de Acesso
- **Frontend**: http://localhost:5173
- **Backend API**: http://localhost:8080
- **Redis Commander**: http://localhost:8081 (admin/admin)

## 📋 Endpoints da API

### Estatísticas Gerais
- `GET /api/v1/stats/totals` - Contadores gerais
- `GET /api/v1/stats/medicos-por-estado` - Médicos agrupados por estado
- `GET /api/v1/stats/hospitais-por-estado` - Hospitais agrupados por estado

### Filtros e Buscas
- `GET /api/v1/stats/hospitais-por-especialidade` - Hospitais por especialidade
- `GET /api/v1/stats/hospitais-por-municipio` - Hospitais por município
- `GET /api/v1/stats/especialidades` - Lista de especialidades

### Lógica de Negócio
- `POST /api/v1/stats/assign-medicos-hospitais` - Executa atribuição médico-hospital

## 🛠️ Tecnologias Utilizadas

### Backend
- **Go 1.25**: Linguagem principal
- **Gorilla Mux**: Roteamento HTTP
- **GORM**: ORM para PostgreSQL
- **Redis**: Cache e filas
- **UUID**: Identificadores únicos
- **gocsv**: Processamento de CSV
- **Excelize**: Processamento de Excel

### Frontend
- **React 19**: Framework principal
- **TypeScript**: Tipagem estática
- **Vite**: Build tool
- **Tailwind CSS**: Estilização
- **Radix UI**: Componentes acessíveis
- **Recharts**: Gráficos e visualizações
- **React Router**: Roteamento

### Infraestrutura
- **PostgreSQL 17**: Banco de dados principal
- **Redis 7**: Cache e filas
- **Docker**: Containerização
- **Docker Compose**: Orquestração

## 📊 Estrutura do Banco de Dados

### Tabelas Principais
- **hospitais**: Dados dos hospitais
- **medicos**: Dados dos médicos
- **estados**: Informações dos estados
- **municipios**: Informações dos municípios
- **pacientes**: Dados dos pacientes (se aplicável)

### Índices Otimizados
- Índices por estado para consultas geográficas
- Índices por especialidade para filtros
- Índices UUID para joins eficientes

## 🔍 Monitoramento e Observabilidade

### Logs
- Logs estruturados no backend
- Rastreamento de operações de upload
- Métricas de performance

### WebSocket
- Atualizações em tempo real do progresso
- Notificações de conclusão de processamento
- Status de saúde do sistema

### Health Checks
- Verificação de conectividade com PostgreSQL
- Verificação de conectividade com Redis
- Monitoramento de workers

## 🎯 Resultados Alcançados

### Performance
- Processamento de arquivos grandes (>100MB) em chunks
- Cache inteligente para consultas frequentes
- Resposta rápida para estatísticas complexas

### Escalabilidade
- Arquitetura preparada para múltiplos workers
- Sistema de filas para processamento assíncrono
- Separação clara entre frontend e backend

### Usabilidade
- Interface moderna e responsiva
- Feedback visual do progresso de upload
- Visualizações interativas de dados

---

**Desenvolvido para o Premiere Challenge**
Demonstração de capacidades em processamento de dados, análise estatística e desenvolvimento full-stack.