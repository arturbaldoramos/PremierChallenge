#!/bin/bash

# Script para testar os endpoints do dashboard
# Execute após iniciar o servidor com: go run cmd/api/main.go

BASE_URL="http://localhost:8080/api/v1"

echo "🏥 Testando Dashboard API - Premiere Challenge"
echo "=============================================="

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Função para testar endpoint
test_endpoint() {
    local name="$1"
    local url="$2"
    local expected_status="$3"
    
    echo -e "\n${YELLOW}🧪 Testando: $name${NC}"
    echo "URL: $url"
    
    response=$(curl -s -w "\n%{http_code}" "$url")
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n -1)
    
    if [ "$http_code" = "$expected_status" ]; then
        echo -e "${GREEN}✅ Sucesso! Status: $http_code${NC}"
        echo "Resposta:"
        echo "$body" | jq . 2>/dev/null || echo "$body"
    else
        echo -e "${RED}❌ Erro! Status esperado: $expected_status, recebido: $http_code${NC}"
        echo "Resposta:"
        echo "$body"
    fi
}

# Verificar se o servidor está rodando
echo -e "\n${YELLOW}🔍 Verificando se o servidor está rodando...${NC}"
if ! curl -s "$BASE_URL/status" > /dev/null; then
    echo -e "${RED}❌ Servidor não está rodando! Inicie com: go run cmd/api/main.go${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Servidor está rodando!${NC}"

# Testar endpoints
test_endpoint "Status da API" "$BASE_URL/status" "200"
test_endpoint "Estatísticas do Dashboard" "$BASE_URL/dashboard/stats" "200"
test_endpoint "Hospitais (sem filtros)" "$BASE_URL/dashboard/hospitais" "200"
test_endpoint "Hospitais (com filtro de estado)" "$BASE_URL/dashboard/hospitais?estado=São Paulo" "200"
test_endpoint "Hospitais (com paginação)" "$BASE_URL/dashboard/hospitais?page=1&limit=5" "200"
test_endpoint "Hospitais (com filtro de especialidade)" "$BASE_URL/dashboard/hospitais?especialidade=Cardiologia" "200"

# Testar endpoint de hospital específico (se existir)
echo -e "\n${YELLOW}🧪 Testando: Detalhes de Hospital${NC}"
echo "URL: $BASE_URL/dashboard/hospital?uuid=test-uuid"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/dashboard/hospital?uuid=test-uuid")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n -1)

if [ "$http_code" = "200" ] || [ "$http_code" = "400" ]; then
    echo -e "${GREEN}✅ Endpoint funcionando! Status: $http_code${NC}"
    echo "Resposta:"
    echo "$body" | jq . 2>/dev/null || echo "$body"
else
    echo -e "${RED}❌ Erro inesperado! Status: $http_code${NC}"
    echo "Resposta:"
    echo "$body"
fi

echo -e "\n${GREEN}🎉 Testes concluídos!${NC}"
echo -e "\n${YELLOW}📝 Próximos passos:${NC}"
echo "1. Verifique se os dados estão sendo retornados corretamente"
echo "2. Teste com dados reais no banco"
echo "3. Integre com o frontend React"
echo "4. Configure os filtros de busca no frontend"
