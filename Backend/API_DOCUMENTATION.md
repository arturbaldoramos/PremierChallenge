# API de Estatísticas e Filtros - Premiere Challenge

## Base URL
```
http://localhost:8080/api/v1
```

## Rotas Disponíveis

### 1. Estatísticas Gerais
**GET** `/stats/totals`

Retorna contadores gerais do sistema.

**Resposta:**
```json
{
  "success": true,
  "data": {
    "total_hospitais": 1500,
    "total_medicos": 280000,
    "total_leitos": 125000,
    "total_pacientes": 1500000
  }
}
```

---

### 2. Hospitais por Estado
**GET** `/stats/hospitais-por-estado`

Retorna quantidade de hospitais agrupados por estado.

**Resposta:**
```json
{
  "success": true,
  "data": [
    {
      "estado": "SP",
      "nome_estado": "São Paulo",
      "count": 350
    },
    {
      "estado": "RJ",
      "nome_estado": "Rio de Janeiro",
      "count": 180
    }
  ]
}
```

---

### 3. Médicos por Estado
**GET** `/stats/medicos-por-estado`

Retorna quantidade de médicos agrupados por estado.

**Resposta:**
```json
{
  "success": true,
  "data": [
    {
      "estado": "SP",
      "nome_estado": "São Paulo",
      "count": 85000
    },
    {
      "estado": "RJ",
      "nome_estado": "Rio de Janeiro",
      "count": 45000
    }
  ]
}
```

---

### 4. Hospitais por Especialidade
**GET** `/stats/hospitais-por-especialidade`

Lista todas as especialidades disponíveis com contadores.

**GET** `/stats/hospitais-por-especialidade?especialidade=Cardiologia`

Filtra hospitais que atendem uma especialidade específica.

**Resposta (sem parâmetro):**
```json
{
  "success": true,
  "data": [
    {
      "especialidade": "Cardiologia",
      "count": 450
    },
    {
      "especialidade": "Neurologia",
      "count": 320
    }
  ]
}
```

**Resposta (com parâmetro):**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "uuid": "550e8400-e29b-41d4-a716-446655440000",
      "nome": "Hospital do Coração",
      "cep": "01234567",
      "especialidades": "Cardiologia,Cirurgia Cardíaca",
      "leitos_totais": 200,
      "cod_municipio": "3550308",
      "bairro": "Vila Mariana"
    }
  ],
  "message": "Hospitais que atendem a especialidade: Cardiologia"
}
```

---

### 5. Hospitais por Município
**GET** `/stats/hospitais-por-municipio`

Lista todos os municípios com contadores de hospitais.

**GET** `/stats/hospitais-por-municipio?municipio=São Paulo`

Filtra hospitais de um município específico.

**Resposta (sem parâmetro):**
```json
{
  "success": true,
  "data": [
    {
      "municipio": "São Paulo",
      "count": 120
    },
    {
      "municipio": "Rio de Janeiro",
      "count": 85
    }
  ]
}
```

---

### 6. Distribuição de Médicos
**GET** `/stats/medicos-distribuicao`

Mostra quantos médicos estão alocados em 1, 2, 3 ou nenhum hospital.

**Resposta:**
```json
{
  "success": true,
  "data": {
    "com_um_hospital": 112000,
    "com_dois_hospitais": 98000,
    "com_tres_hospitais": 42000,
    "sem_hospitais": 28000
  },
  "message": "Distribuição simulada - implementar lógica de negócio real"
}
```

---

### 7. Lista de Especialidades
**GET** `/stats/especialidades`

Retorna todas as especialidades médicas disponíveis no sistema.

**Resposta:**
```json
{
  "success": true,
  "data": [
    "Cardiologia",
    "Dermatologia",
    "Endocrinologia",
    "Neurologia",
    "Ortopedia",
    "Pediatria"
  ]
}
```

---

### 8. Atribuição Médicos-Hospitais (Lógica de Negócio)
**POST** `/stats/assign-medicos-hospitais`

Implementa a lógica de negócio para atribuir médicos a hospitais baseado nas regras:
- Médicos podem atender até 3 hospitais
- Hospitais devem estar num raio de 30km do médico
- Hospital deve ter a mesma especialidade que o médico pratica

**Resposta:**
```json
{
  "success": true,
  "data": [
    {
      "medico_uuid": "550e8400-e29b-41d4-a716-446655440001",
      "medico_nome": "Dr. João Silva",
      "hospital_uuid": "550e8400-e29b-41d4-a716-446655440002",
      "hospital_nome": "Hospital do Coração",
      "especialidade": "Cardiologia",
      "distancia_km": 15.5
    }
  ],
  "message": "Atribuições médico-hospital baseadas na lógica de negócio (raio 30km, mesma especialidade, máx 3 hospitais)"
}
```

---

## Exemplos de Uso com cURL

### Testar total de médicos (deve retornar 280 mil):
```bash
curl -X GET "http://localhost:8080/api/v1/stats/totals"
```

### Buscar hospitais de cardiologia:
```bash
curl -X GET "http://localhost:8080/api/v1/stats/hospitais-por-especialidade?especialidade=Cardiologia"
```

### Buscar hospitais de São Paulo:
```bash
curl -X GET "http://localhost:8080/api/v1/stats/hospitais-por-municipio?municipio=São%20Paulo"
```

### Executar lógica de atribuição médico-hospital:
```bash
curl -X POST "http://localhost:8080/api/v1/stats/assign-medicos-hospitais"
```

---

## Códigos de Status HTTP

- `200 OK`: Requisição bem-sucedida
- `400 Bad Request`: Parâmetros inválidos
- `500 Internal Server Error`: Erro interno do servidor

## Formato de Resposta Padrão

Todas as respostas seguem o formato:
```json
{
  "success": boolean,
  "data": object|array,
  "message": "string (opcional)",
  "error": "string (apenas em caso de erro)"
}
```


