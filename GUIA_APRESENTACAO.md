# 🚀 Guia de Apresentação - Premiere Challenge

## 📋 **Para rodar em qualquer PC:**

### **Pré-requisitos:**
- Docker Desktop instalado
- Conexão com internet (para baixar as imagens)

### **Passo 1: Baixar o projeto**
```bash
# Clone ou baixe apenas estes arquivos:
# - docker-compose.yml
# - README.md
# - GUIA_APRESENTACAO.md
```

### **Passo 2: Executar**
```bash
# No diretório do projeto:
docker-compose up -d
```

### **Passo 3: Aguardar (1-2 minutos)**
O sistema vai:
1. Baixar as imagens do DockerHub
2. Iniciar PostgreSQL e Redis
3. Fazer migrations automaticamente
4. Subir Backend e Frontend

### **Passo 4: Acessar**
- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080/api/v1/stats/totals

---

## 🎯 **Para a Apresentação:**

### **1. Mostrar o Sistema Funcionando:**
```bash
# Verificar se tudo está rodando:
docker-compose ps

# Deve mostrar:
# premiere_frontend    Up    0.0.0.0:3000->80/tcp
# premiere_backend     Up    0.0.0.0:8080->8080/tcp
# premiere_postgres    Up    0.0.0.0:5432->5432/tcp
# premiere_redis       Up    0.0.0.0:6379->6379/tcp
```

### **2. Demonstrar API:**
```bash
# Estatísticas totais:
curl http://localhost:8080/api/v1/stats/totals

# Especialidades disponíveis:
curl http://localhost:8080/api/v1/stats/especialidades

# Médicos por estado:
curl http://localhost:8080/api/v1/stats/medicos-por-estado
```

### **3. Mostrar Frontend:**
- Acesse: http://localhost:3000
- Interface moderna com React
- Upload de arquivos
- Visualização de estatísticas

### **4. Mostrar Escalabilidade:**
```bash
# Ver logs em tempo real:
docker logs premiere_backend -f

# Mostrar workers rodando
# Mostrar processamento assíncrono
```

---

## 🛠️ **Comandos Úteis:**

### **Parar tudo:**
```bash
docker-compose down
```

### **Ver logs:**
```bash
# Backend:
docker logs premiere_backend

# Frontend:
docker logs premiere_frontend

# Banco:
docker logs premiere_postgres
```

### **Reiniciar se necessário:**
```bash
docker-compose restart
```

### **Ver uso de recursos:**
```bash
docker stats
```

---

## 📊 **Pontos para Destacar:**

### **1. Arquitetura Moderna:**
- ✅ Microserviços containerizados
- ✅ Frontend React + Backend Go
- ✅ Banco PostgreSQL + Cache Redis
- ✅ Comunicação via REST API + WebSocket

### **2. DevOps/Produção:**
- ✅ Imagens no DockerHub públicas
- ✅ Multi-stage builds otimizados
- ✅ Health checks configurados
- ✅ Volumes persistentes
- ✅ Networking isolado

### **3. Escalabilidade:**
- ✅ Workers assíncronos
- ✅ Processamento em chunks
- ✅ Cache inteligente
- ✅ Fila de tarefas

### **4. Facilidade de Deploy:**
- ✅ Um comando: `docker-compose up -d`
- ✅ Funciona em qualquer PC
- ✅ Auto-configuração
- ✅ Zero configuração manual

---

## 🔗 **URLs Importantes:**

| Serviço | URL | Descrição |
|---------|-----|-----------|
| **Frontend** | http://localhost:3000 | Interface principal |
| **API Stats** | http://localhost:8080/api/v1/stats/totals | Estatísticas |
| **Health** | http://localhost:8080/health | Status do sistema |
| **DockerHub** | https://hub.docker.com/u/blamet | Imagens publicadas |

---

## 🎉 **Script de Demo Automático:**

```bash
#!/bin/bash
echo "🚀 Iniciando Premiere Challenge..."
docker-compose up -d

echo "⏳ Aguardando serviços (30s)..."
sleep 30

echo "🌐 Frontend: http://localhost:3000"
echo "📊 API: http://localhost:8080/api/v1/stats/totals"

echo "📋 Status dos containers:"
docker-compose ps

echo "✅ Sistema pronto para apresentação!"
```

---

## 🆘 **Troubleshooting:**

### **Se der erro de porta:**
```bash
# Verificar se portas estão livres:
netstat -an | grep :3000
netstat -an | grep :8080

# Se ocupadas, parar outros serviços ou usar outras portas
```

### **Se não baixar as imagens:**
```bash
# Forçar download:
docker pull blamet/premiere-challenge-frontend:latest
docker pull blamet/premiere-challenge-backend:latest
```

### **Se container não subir:**
```bash
# Ver detalhes do erro:
docker-compose logs [nome-do-servico]
```

---

**🎯 Tudo pronto para uma apresentação perfeita!**