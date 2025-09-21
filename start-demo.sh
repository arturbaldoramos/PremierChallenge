#!/bin/bash

echo "🚀 Iniciando Premiere Challenge Demo..."
echo ""

# Parar qualquer instância anterior
echo "🛑 Parando instâncias anteriores..."
docker-compose down 2>/dev/null

echo ""
echo "📦 Iniciando containers..."
docker-compose up -d

echo ""
echo "⏳ Aguardando serviços iniciarem (30 segundos)..."
sleep 30

echo ""
echo "📊 Status dos containers:"
docker-compose ps

echo ""
echo "🌐 URLs disponíveis:"
echo "   Frontend:  http://localhost:3000"
echo "   Backend:   http://localhost:8080/api/v1/stats/totals"
echo ""

# Testar se está funcionando
echo "🧪 Testando API..."
if curl -s http://localhost:8080/api/v1/stats/totals > /dev/null; then
    echo "✅ API funcionando!"
else
    echo "⚠️  API ainda carregando... aguarde mais um pouco"
fi

echo ""
echo "🎉 Sistema pronto para apresentação!"
echo "   Para parar: docker-compose down"