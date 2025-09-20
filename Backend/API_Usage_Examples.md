# API de Upload de Arquivos - Guia de Uso

## Visão Geral

O sistema implementa rotas específicas para upload e processamento de arquivos para cada domínio. Atualmente suporta arquivos CSV com detecção automática de tipo e verificação de duplicatas.

## Endpoints Disponíveis

### 1. Upload de Estados
**POST** `/api/v1/upload/estados`

### 2. Upload de Municípios
**POST** `/api/v1/upload/municipios`

### 3. Upload de Hospitais
**POST** `/api/v1/upload/hospitais`

### 4. Upload de Pacientes
**POST** `/api/v1/upload/pacientes`

### 5. Upload de Médicos
**POST** `/api/v1/upload/medicos`

### 6. Upload de CID10
**POST** `/api/v1/upload/cid10`

## Formato da Requisição

Todas as rotas de upload esperam:
- **Content-Type**: `multipart/form-data`
- **Campo do arquivo**: `file`
- **Tamanho máximo**: 32MB

## Exemplos de Uso com cURL

### Estados
```bash
curl -X POST http://localhost:8080/api/v1/upload/estados \
  -F "file=@estados.csv"
```

### Municípios
```bash
curl -X POST http://localhost:8080/api/v1/upload/municipios \
  -F "file=@municipios.csv"
```

### Hospitais
```bash
curl -X POST http://localhost:8080/api/v1/upload/hospitais \
  -F "file=@hospitais.csv"
```

## Formato dos Arquivos CSV

### Estados (estados.csv)
```csv
id,codigo,unidade_federativa,nome,regiao,latitude,longitude
1,AC,AC,Acre,Norte,-9.0238,-70.8120
2,AL,AL,Alagoas,Nordeste,-9.5713,-36.7820
```

### Municípios (municipios.csv)
```csv
id,codigo,nome,latitude,longitude,capital,codigo_uf,siafi_id,ddd,fuso_hora,populacao
1,1200013,Assis Brasil,-10.9472,-69.5731,0,12,0139,68,America/Rio_Branco,6072
```

### Hospitais (hospitais.csv)
```csv
id,uuid,nome,cep,especialidades,leitos_totais,cod_municipio,bairro
1,550e8400-e29b-41d4-a716-446655440000,Hospital Municipal,12345678,"[""Cardiologia""]",100,1200013,Centro
```

### Pacientes (pacientes.csv)
```csv
id,nome,cpf,rg,data_nascimento,genero,tipo_sanguineo,endereco,municipio_id,cep,telefone,email,contato_emergencia,convenio,numero_carteira,status
550e8400-e29b-41d4-a716-446655440001,João Silva,12345678901,123456789,1990-01-15,M,O+,Rua A 123,1,12345678,11999999999,joao@email.com,11888888888,SUS,,ativo
```

### Médicos (medicos.csv)
```csv
id,nome,especialidade,cod_municipio
550e8400-e29b-41d4-a716-446655440002,Dr. Maria Santos,Cardiologia,1200013
```

### CID10 (cid10.csv)
```csv
id,codigo,descricao,categoria,grupo
1,A00,Cólera,Doenças infecciosas,A00-A09
```

## Resposta da API

### Sucesso (200 OK)
```json
{
  "success": true,
  "message": "Estados uploaded successfully",
  "file_type": "csv",
  "file_name": "estados.csv",
  "total_records": 27,
  "inserted_records": 20,
  "updated_records": 7
}
```

### Erro (400 Bad Request)
```json
{
  "success": false,
  "message": "Failed to parse CSV file",
  "errors": ["column 'invalid_field' not found"]
}
```

## Recursos Implementados

✅ **Detecção automática de tipo de arquivo** (CSV, XLSX, XML)
✅ **Verificação de duplicatas** - Evita inserir dados já existentes
✅ **Upsert automático** - Insere novos registros e atualiza existentes
✅ **Processamento em batch** - Otimizado para grandes volumes
✅ **Validação de dados** - Verifica formato e integridade
✅ **Resposta detalhada** - Informa quantos registros foram inseridos/atualizados

## Próximas Implementações

🔄 **Processamento em chunks** - Para arquivos muito grandes (>32MB)
🔄 **Suporte a XLSX** - Planilhas Excel
🔄 **Suporte a XML** - Arquivos XML estruturados
🔄 **Relatório de erros detalhado** - Linha por linha
🔄 **Progress tracking** - Acompanhamento do progresso para arquivos grandes

## Estratégia de Verificação de Duplicatas

- **Estados**: Verificação por `codigo`
- **Municípios**: Verificação por `codigo`
- **Hospitais**: Verificação por `uuid`
- **Pacientes**: Verificação por `cpf`
- **Médicos**: Verificação por `uuid`
- **CID10**: Verificação por `codigo`

## Executar o Servidor

```bash
go run cmd/api/main.go
```

O servidor estará disponível em `http://localhost:8080`