# Teste do Dashboard API

## Endpoints Disponíveis

### 1. Estatísticas Gerais do Dashboard
**GET** `/api/v1/dashboard/stats`

Retorna totalizadores, gráficos e dados de regiões.

**Exemplo de Resposta:**
```json
{
  "totalizadores": {
    "total_hospitais": 150,
    "total_medicos": 2500,
    "total_pacientes": 50000,
    "total_leitos": 3000,
    "media_leitos_hospital": 20.0,
    "percentual_convenio": 65.5
  },
  "graficos": {
    "top_10_doencas": [
      {
        "codigo": "A15",
        "descricao": "Tuberculose respiratória",
        "quantidade": 1500
      }
    ],
    "doenca_por_estado": [
      {
        "estado": "São Paulo",
        "doenca": "Hipertensão",
        "quantidade": 5000
      }
    ],
    "percentual_populacao_doente": 2.5
  },
  "regioes": {
    "regiao_mais_doente": "Sudeste",
    "percentual_regiao": 3.2
  }
}
```

### 2. Hospitais com Filtros
**GET** `/api/v1/dashboard/hospitais?estado=São Paulo&cidade=São Paulo&especialidade=Cardiologia&page=1&limit=20`

**Parâmetros:**
- `estado` (opcional): Filtro por estado
- `cidade` (opcional): Filtro por cidade
- `especialidade` (opcional): Filtro por especialidade
- `page` (opcional): Página (padrão: 1)
- `limit` (opcional): Limite por página (padrão: 20, máximo: 100)

**Exemplo de Resposta:**
```json
{
  "hospitais": [
    {
      "uuid": "1b2c137e-75e1-4644-b1ca-c04e1055443a",
      "nome": "Hospital Municipal Santo Antônio",
      "bairro": "Jardim",
      "leitos_totais": 335,
      "especialidades": ["Neurologia", "Ortopedia"],
      "municipio": "São Paulo",
      "estado": "São Paulo",
      "regiao": "Sudeste"
    }
  ],
  "total": 25,
  "page": 1,
  "limit": 20,
  "pages": 2
}
```

### 3. Detalhes de Hospital
**GET** `/api/v1/dashboard/hospital?uuid=1b2c137e-75e1-4644-b1ca-c04e1055443a`

**Exemplo de Resposta:**
```json
{
  "hospital": {
    "uuid": "1b2c137e-75e1-4644-b1ca-c04e1055443a",
    "nome": "Hospital Municipal Santo Antônio",
    "bairro": "Jardim",
    "leitos_totais": 335,
    "especialidades": ["Neurologia", "Ortopedia"],
    "municipio": "São Paulo",
    "estado": "São Paulo",
    "regiao": "Sudeste",
    "latitude": "-23.5505",
    "longitude": "-46.6333"
  },
  "medicos": [
    {
      "uuid": "cfa9a86c-def0-41cb-bfe6-c2a11d962e4f",
      "nome": "Fernanda Pereira Araújo",
      "especialidade": "Neurologia",
      "municipio": "São Paulo",
      "latitude": "-23.5505",
      "longitude": "-46.6333",
      "distancia_km": 2.5
    }
  ]
}
```

## Como Testar

### 1. Usando curl
```bash
# Estatísticas gerais
curl -X GET "http://localhost:8080/api/v1/dashboard/stats"

# Hospitais filtrados
curl -X GET "http://localhost:8080/api/v1/dashboard/hospitais?estado=São Paulo&page=1&limit=10"

# Detalhes de hospital
curl -X GET "http://localhost:8080/api/v1/dashboard/hospital?uuid=1b2c137e-75e1-4644-b1ca-c04e1055443a"
```

### 2. Usando Postman
1. Importe as URLs acima
2. Configure headers: `Content-Type: application/json`
3. Execute as requisições

### 3. Usando JavaScript (Frontend)
```javascript
// Buscar estatísticas
const response = await fetch('http://localhost:8080/api/v1/dashboard/stats');
const data = await response.json();
console.log(data);

// Buscar hospitais com filtros
const hospitaisResponse = await fetch('http://localhost:8080/api/v1/dashboard/hospitais?estado=São Paulo');
const hospitais = await hospitaisResponse.json();
console.log(hospitais);
```

## Notas Importantes

1. **Distância Haversine**: Os médicos são filtrados por distância de 30km do hospital usando a fórmula de Haversine
2. **Especialidades**: Médicos de Clínica Geral podem trabalhar em qualquer especialidade
3. **Performance**: As queries são otimizadas com índices geográficos
4. **Paginação**: Use os parâmetros `page` e `limit` para controlar a paginação
5. **Filtros**: Todos os filtros são opcionais e podem ser combinados

## Estrutura do Banco

### Tabelas Principais
- `estados`: Estados brasileiros com coordenadas
- `municipios`: Municípios com população e coordenadas
- `hospitais`: Hospitais com especialidades e leitos
- `medicos`: Médicos com especialidades e localização
- `pacientes`: Pacientes com CID-10 e convênio
- `cid10s`: Códigos CID-10 com especialidades

### Índices Criados
- Índices geográficos para consultas de distância
- Índices em campos de filtro para performance
- Índices em campos de relacionamento
