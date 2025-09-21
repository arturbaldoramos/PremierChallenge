# Sistema de Detecção de Arquivos

Este sistema permite detectar tipos de arquivo para upload e processamento.

## Funcionalidades

### 🔍 Detecção de Tipos de Arquivo
- **XML**: Detecção por extensão e análise de conteúdo
- **JSON**: Validação de sintaxe e estrutura
- **CSV**: Análise de separadores e estrutura tabular
- **PDF**: Processamento de texto extraído
- **TXT**: Divisão por linhas

## Como Usar

### 1. Detecção de Arquivo

```typescript
import { FileDetector } from '@/lib/file-detector';

const file = event.target.files[0];
const fileType = FileDetector.detectFileType(file);

if (fileType) {
  console.log(`Tipo detectado: ${fileType.description}`);
}
```

### 2. Upload de Arquivo

```typescript
import { api } from '@/lib/api';

const file = event.target.files[0];
const result = await api.states.uploadFile(file);

console.log('Arquivo enviado com sucesso:', result.message);
```

### 3. Componente de Upload

```tsx
import FileUpload from '@/components/FileUpload';

<FileUpload 
  onFilesProcessed={(results) => {
    console.log('Arquivos processados:', results);
  }}
  maxFiles={10}
  maxFileSize={50 * 1024 * 1024} // 50MB
/>
```

## Estrutura de Chunks

Cada chunk contém:

```typescript
interface Chunk {
  id: string;                    // ID único do chunk
  type: string;                  // Tipo do arquivo original
  content: any;                  // Conteúdo processado
  metadata: {
    originalFile: string;         // Nome do arquivo original
    chunkIndex: number;          // Índice do chunk
    totalChunks: number;          // Total de chunks
    size: number;                // Tamanho em bytes
    timestamp: string;           // Data/hora de processamento
  };
}
```

## Tipos de Conteúdo por Arquivo

### XML
```typescript
content: {
  elements: Array<{
    tag: string;
    value?: any;
    children?: any[];
    path: string;
  }>;
  chunkInfo: {
    startIndex: number;
    endIndex: number;
    totalElements: number;
  };
}
```

### JSON
```typescript
content: {
  items: any[];                  // Array de itens do JSON
  chunkInfo: {
    startIndex: number;
    endIndex: number;
    totalItems: number;
  };
}
```

### CSV
```typescript
content: {
  headers: string[];             // Cabeçalhos das colunas
  rows: any[][];                // Dados das linhas
  chunkInfo: {
    startIndex: number;
    endIndex: number;
    totalRows: number;
  };
}
```

### PDF/TXT
```typescript
content: {
  text: string;                  // Texto extraído
  chunkInfo: {
    startIndex: number;
    endIndex: number;
    totalWords: number;
  };
}
```

## Recursos Avançados

### Download de Chunks
```typescript
// Download individual
downloadChunk(chunk);

// Download em lote
downloadAllChunks(chunks);
```

### Reconstrução de Arquivo
```typescript
const originalContent = FileChunker.combineChunks(chunks);
```

### Estatísticas
```typescript
const stats = {
  totalFiles: results.length,
  totalChunks: results.reduce((sum, r) => sum + r.totalChunks, 0),
  totalErrors: results.reduce((sum, r) => sum + r.errors.length, 0),
  avgProcessingTime: results.reduce((sum, r) => sum + r.processingTime, 0) / results.length
};
```

## Limitações

- **PDF**: Atualmente usa simulação. Para produção, integrar com `pdf-parse`
- **Tamanho**: Arquivos muito grandes podem causar problemas de memória
- **Tipos**: Suporte limitado aos tipos implementados

## Próximos Passos

1. Integração com `pdf-parse` para PDFs reais
2. Suporte a mais tipos de arquivo (Excel, Word, etc.)
3. Processamento assíncrono para arquivos grandes
4. Compressão de chunks para economizar espaço
5. Cache de chunks processados
