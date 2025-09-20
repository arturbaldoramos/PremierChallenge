import Papa from 'papaparse';
import * as xml2js from 'xml2js';
import { FileDetector } from './file-detector';
import type { FileTypeInfo } from './file-detector';

export interface Chunk {
  id: string;
  type: string;
  content: any;
  metadata: {
    originalFile: string;
    chunkIndex: number;
    totalChunks: number;
    size: number;
    timestamp: string;
  };
}

export interface ProcessingResult {
  fileType: FileTypeInfo;
  chunks: Chunk[];
  totalChunks: number;
  processingTime: number;
  errors: string[];
}

export class FileChunker {
  private static readonly DEFAULT_CHUNK_SIZE = 10000; // Tamanho padrão do chunk para arquivos grandes

  /**
   * Processa um arquivo e o divide em chunks baseado no tipo
   */
  static async processFile(file: File, chunkSize?: number): Promise<ProcessingResult> {
    const startTime = Date.now();
    const errors: string[] = [];
    const chunks: Chunk[] = [];
    
    try {
      const fileType = await FileDetector.detectFileTypeAsync(file);
      if (!fileType) {
        throw new Error('Tipo de arquivo não suportado');
      }

      const effectiveChunkSize = chunkSize || this.DEFAULT_CHUNK_SIZE;
      
      switch (fileType.type) {
        case 'xml':
          await this.processXML(file, effectiveChunkSize, chunks, errors);
          break;
        case 'json':
          await this.processJSON(file, effectiveChunkSize, chunks, errors);
          break;
        case 'csv':
          await this.processCSV(file, effectiveChunkSize, chunks, errors);
          break;
        case 'pdf':
          await this.processPDF(file, effectiveChunkSize, chunks, errors);
          break;
        case 'txt':
          await this.processTXT(file, effectiveChunkSize, chunks, errors);
          break;
        default:
          throw new Error(`Processamento não implementado para tipo: ${fileType.type}`);
      }

      const processingTime = Date.now() - startTime;

      return {
        fileType,
        chunks,
        totalChunks: chunks.length,
        processingTime,
        errors
      };

    } catch (error) {
      errors.push(`Erro geral: ${error instanceof Error ? error.message : 'Erro desconhecido'}`);
      return {
        fileType: { type: 'unknown', extension: '', mimeType: '', description: 'Unknown' },
        chunks: [],
        totalChunks: 0,
        processingTime: Date.now() - startTime,
        errors
      };
    }
  }


  /**
   * Processa arquivos XML
   */
  private static async processXML(file: File, chunkSize: number, chunks: Chunk[], errors: string[]): Promise<void> {
    try {
      // Para arquivos muito grandes, processa em streaming
      if (file.size > 100 * 1024 * 1024) { // > 100MB
        await this.processXMLStreaming(file, chunkSize, chunks, errors);
        return;
      }

      const content = await file.text();
      const parser = new xml2js.Parser({ 
        explicitArray: false,
        mergeAttrs: true,
        trim: true
      });
      const result = await parser.parseStringPromise(content);
      
      // Divide o XML em elementos menores preservando entidades
      const elements = this.extractXMLElements(result);
      const totalElements = elements.length;
      
      for (let i = 0; i < totalElements; i += chunkSize) {
        const chunkElements = elements.slice(i, i + chunkSize);
        const chunk: Chunk = {
          id: `xml-chunk-${i}-${Date.now()}`,
          type: 'xml',
          content: {
            elements: chunkElements,
            chunkInfo: {
              startIndex: i,
              endIndex: Math.min(i + chunkSize - 1, totalElements - 1),
              totalElements
            }
          },
          metadata: {
            originalFile: file.name,
            chunkIndex: Math.floor(i / chunkSize),
            totalChunks: Math.ceil(totalElements / chunkSize),
            size: JSON.stringify(chunkElements).length,
            timestamp: new Date().toISOString()
          }
        };
        chunks.push(chunk);
      }
    } catch (error) {
      errors.push(`Erro ao processar XML: ${error instanceof Error ? error.message : 'Erro desconhecido'}`);
    }
  }

  /**
   * Processa XML grandes em streaming para preservar entidades
   */
  private static async processXMLStreaming(file: File, chunkSize: number, chunks: Chunk[], errors: string[]): Promise<void> {
    try {
      const content = await file.text();
      
      // Processa XML grande dividindo por tags principais
      const xmlChunks = this.splitXMLByEntities(content, chunkSize);
      
      xmlChunks.forEach((xmlChunk, index) => {
        const chunk: Chunk = {
          id: `xml-streaming-chunk-${index}-${Date.now()}`,
          type: 'xml',
          content: {
            xmlContent: xmlChunk,
            chunkInfo: {
              chunkIndex: index,
              totalChunks: xmlChunks.length,
              isStreaming: true
            }
          },
          metadata: {
            originalFile: file.name,
            chunkIndex: index,
            totalChunks: xmlChunks.length,
            size: xmlChunk.length,
            timestamp: new Date().toISOString()
          }
        };
        chunks.push(chunk);
      });
    } catch (error) {
      errors.push(`Erro ao processar XML grande: ${error instanceof Error ? error.message : 'Erro desconhecido'}`);
    }
  }

  /**
   * Divide XML por entidades preservando integridade - NUNCA corta entidades no meio
   */
  private static splitXMLByEntities(xmlContent: string, chunkSize: number): string[] {
    const chunks: string[] = [];
    let currentChunk = '';
    let currentEntity = '';
    let depth = 0;
    let inEntity = false;
    
    const lines = xmlContent.split('\n');
    
    for (const line of lines) {
      currentEntity += line + '\n';
      
      // Conta tags de abertura e fechamento para detectar entidades completas
      const openTags = (line.match(/<[^\/][^>]*>/g) || []).length;
      const closeTags = (line.match(/<\/[^>]*>/g) || []).length;
      
      depth += openTags - closeTags;
      
      // Detecta início de entidade (primeira tag de abertura)
      if (openTags > 0 && !inEntity) {
        inEntity = true;
      }
      
      // Detecta fim de entidade (depth volta a 0)
      if (inEntity && depth === 0) {
        // Entidade completa encontrada
        if (currentChunk.length + currentEntity.length > chunkSize && currentChunk.length > 0) {
          // Chunk atual está cheio, salva ele
          chunks.push(currentChunk.trim());
          currentChunk = '';
        }
        
        // Adiciona a entidade completa ao chunk atual
        currentChunk += currentEntity;
        currentEntity = '';
        inEntity = false;
      }
    }
    
    // Adiciona qualquer conteúdo restante
    if (currentEntity.trim()) {
      currentChunk += currentEntity;
    }
    
    if (currentChunk.trim()) {
      chunks.push(currentChunk.trim());
    }
    
    return chunks.length > 0 ? chunks : [xmlContent];
  }

  /**
   * Processa arquivos JSON
   */
  private static async processJSON(file: File, chunkSize: number, chunks: Chunk[], errors: string[]): Promise<void> {
    try {
      // Para arquivos muito grandes, processa em streaming
      if (file.size > 100 * 1024 * 1024) { // > 100MB
        await this.processJSONStreaming(file, chunkSize, chunks, errors);
        return;
      }

      const content = await file.text();
      const jsonData = JSON.parse(content);
      
      let items: any[] = [];
      
      if (Array.isArray(jsonData)) {
        items = jsonData;
      } else if (typeof jsonData === 'object') {
        // Se for um objeto, converte em array de propriedades
        items = Object.entries(jsonData).map(([key, value]) => ({ key, value }));
      } else {
        items = [jsonData];
      }
      
      const totalItems = items.length;
      
      for (let i = 0; i < totalItems; i += chunkSize) {
        const chunkItems = items.slice(i, i + chunkSize);
        const chunk: Chunk = {
          id: `json-chunk-${i}-${Date.now()}`,
          type: 'json',
          content: {
            items: chunkItems,
            chunkInfo: {
              startIndex: i,
              endIndex: Math.min(i + chunkSize - 1, totalItems - 1),
              totalItems
            }
          },
          metadata: {
            originalFile: file.name,
            chunkIndex: Math.floor(i / chunkSize),
            totalChunks: Math.ceil(totalItems / chunkSize),
            size: JSON.stringify(chunkItems).length,
            timestamp: new Date().toISOString()
          }
        };
        chunks.push(chunk);
      }
    } catch (error) {
      errors.push(`Erro ao processar JSON: ${error instanceof Error ? error.message : 'Erro desconhecido'}`);
    }
  }

  /**
   * Processa JSON grandes em streaming para preservar entidades
   */
  private static async processJSONStreaming(file: File, chunkSize: number, chunks: Chunk[], errors: string[]): Promise<void> {
    try {
      const content = await file.text();
      
      // Processa JSON grande dividindo por objetos/arrays
      const jsonChunks = this.splitJSONByEntities(content, chunkSize);
      
      jsonChunks.forEach((jsonChunk, index) => {
        const chunk: Chunk = {
          id: `json-streaming-chunk-${index}-${Date.now()}`,
          type: 'json',
          content: {
            jsonContent: jsonChunk,
            chunkInfo: {
              chunkIndex: index,
              totalChunks: jsonChunks.length,
              isStreaming: true
            }
          },
          metadata: {
            originalFile: file.name,
            chunkIndex: index,
            totalChunks: jsonChunks.length,
            size: jsonChunk.length,
            timestamp: new Date().toISOString()
          }
        };
        chunks.push(chunk);
      });
    } catch (error) {
      errors.push(`Erro ao processar JSON grande: ${error instanceof Error ? error.message : 'Erro desconhecido'}`);
    }
  }

  /**
   * Divide JSON por entidades preservando integridade - NUNCA corta entidades no meio
   */
  private static splitJSONByEntities(jsonContent: string, chunkSize: number): string[] {
    const chunks: string[] = [];
    let currentChunk = '';
    let currentEntity = '';
    let braceCount = 0;
    let bracketCount = 0;
    let inString = false;
    let escapeNext = false;
    let inEntity = false;
    
    for (let i = 0; i < jsonContent.length; i++) {
      const char = jsonContent[i];
      currentEntity += char;
      
      if (escapeNext) {
        escapeNext = false;
        continue;
      }
      
      if (char === '\\') {
        escapeNext = true;
        continue;
      }
      
      if (char === '"') {
        inString = !inString;
        continue;
      }
      
      if (!inString) {
        if (char === '{') {
          braceCount++;
          if (!inEntity) inEntity = true; // Início de entidade
        }
        if (char === '}') braceCount--;
        if (char === '[') {
          bracketCount++;
          if (!inEntity) inEntity = true; // Início de entidade
        }
        if (char === ']') bracketCount--;
        
        // Detecta fim de entidade (objeto ou array completo)
        if (inEntity && braceCount === 0 && bracketCount === 0) {
          // Entidade completa encontrada
          if (currentChunk.length + currentEntity.length > chunkSize && currentChunk.length > 0) {
            // Chunk atual está cheio, salva ele
            chunks.push(currentChunk.trim());
            currentChunk = '';
          }
          
          // Adiciona a entidade completa ao chunk atual
          currentChunk += currentEntity;
          currentEntity = '';
          inEntity = false;
        }
      }
    }
    
    // Adiciona qualquer conteúdo restante
    if (currentEntity.trim()) {
      currentChunk += currentEntity;
    }
    
    if (currentChunk.trim()) {
      chunks.push(currentChunk.trim());
    }
    
    return chunks.length > 0 ? chunks : [jsonContent];
  }

  /**
   * Processa arquivos CSV
   */
  private static async processCSV(file: File, chunkSize: number, chunks: Chunk[], errors: string[]): Promise<void> {
    try {
      const content = await file.text();
      
      Papa.parse(content, {
        complete: (results) => {
          const data = results.data as any[][];
          const headers = data[0] || [];
          const rows = data.slice(1);
          
          const totalRows = rows.length;
          
          for (let i = 0; i < totalRows; i += chunkSize) {
            const chunkRows = rows.slice(i, i + chunkSize);
            const chunk: Chunk = {
              id: `csv-chunk-${i}-${Date.now()}`,
              type: 'csv',
              content: {
                headers,
                rows: chunkRows,
                chunkInfo: {
                  startIndex: i,
                  endIndex: Math.min(i + chunkSize - 1, totalRows - 1),
                  totalRows
                }
              },
              metadata: {
                originalFile: file.name,
                chunkIndex: Math.floor(i / chunkSize),
                totalChunks: Math.ceil(totalRows / chunkSize),
                size: JSON.stringify({ headers, rows: chunkRows }).length,
                timestamp: new Date().toISOString()
              }
            };
            chunks.push(chunk);
          }
        },
        error: (error: any) => {
          errors.push(`Erro ao processar CSV: ${error.message}`);
        }
      });
    } catch (error) {
      errors.push(`Erro ao processar CSV: ${error instanceof Error ? error.message : 'Erro desconhecido'}`);
    }
  }

  /**
   * Processa arquivos PDF
   */
  private static async processPDF(file: File, chunkSize: number, chunks: Chunk[], errors: string[]): Promise<void> {
    try {
      // Para PDF, vamos dividir por páginas ou por tamanho de texto
      await file.arrayBuffer();
      
      // Simulação de processamento de PDF (em produção, usar pdf-parse)
      const textContent = `Conteúdo extraído do PDF: ${file.name}`;
      const words = textContent.split(' ');
      const totalWords = words.length;
      
      for (let i = 0; i < totalWords; i += chunkSize) {
        const chunkWords = words.slice(i, i + chunkSize);
        const chunk: Chunk = {
          id: `pdf-chunk-${i}-${Date.now()}`,
          type: 'pdf',
          content: {
            text: chunkWords.join(' '),
            chunkInfo: {
              startIndex: i,
              endIndex: Math.min(i + chunkSize - 1, totalWords - 1),
              totalWords
            }
          },
          metadata: {
            originalFile: file.name,
            chunkIndex: Math.floor(i / chunkSize),
            totalChunks: Math.ceil(totalWords / chunkSize),
            size: chunkWords.join(' ').length,
            timestamp: new Date().toISOString()
          }
        };
        chunks.push(chunk);
      }
    } catch (error) {
      errors.push(`Erro ao processar PDF: ${error instanceof Error ? error.message : 'Erro desconhecido'}`);
    }
  }

  /**
   * Processa arquivos TXT
   */
  private static async processTXT(file: File, chunkSize: number, chunks: Chunk[], errors: string[]): Promise<void> {
    try {
      const content = await file.text();
      const lines = content.split('\n');
      const totalLines = lines.length;
      
      for (let i = 0; i < totalLines; i += chunkSize) {
        const chunkLines = lines.slice(i, i + chunkSize);
        const chunk: Chunk = {
          id: `txt-chunk-${i}-${Date.now()}`,
          type: 'txt',
          content: {
            lines: chunkLines,
            chunkInfo: {
              startIndex: i,
              endIndex: Math.min(i + chunkSize - 1, totalLines - 1),
              totalLines
            }
          },
          metadata: {
            originalFile: file.name,
            chunkIndex: Math.floor(i / chunkSize),
            totalChunks: Math.ceil(totalLines / chunkSize),
            size: chunkLines.join('\n').length,
            timestamp: new Date().toISOString()
          }
        };
        chunks.push(chunk);
      }
    } catch (error) {
      errors.push(`Erro ao processar TXT: ${error instanceof Error ? error.message : 'Erro desconhecido'}`);
    }
  }

  /**
   * Extrai elementos de um objeto XML parseado
   */
  private static extractXMLElements(xmlObject: any): any[] {
    const elements: any[] = [];
    
    function traverse(obj: any, path: string = '') {
      if (typeof obj === 'object' && obj !== null) {
        if (Array.isArray(obj)) {
          obj.forEach((item, index) => traverse(item, `${path}[${index}]`));
        } else {
          Object.entries(obj).forEach(([key, value]) => {
            if (typeof value === 'object' && value !== null) {
              if (Array.isArray(value)) {
                elements.push({ tag: key, children: value, path: `${path}.${key}` });
                value.forEach((item, index) => traverse(item, `${path}.${key}[${index}]`));
              } else {
                elements.push({ tag: key, value, path: `${path}.${key}` });
                traverse(value, `${path}.${key}`);
              }
            } else {
              elements.push({ tag: key, value, path: `${path}.${key}` });
            }
          });
        }
      }
    }
    
    traverse(xmlObject);
    return elements;
  }

  /**
   * Combina chunks de volta em um arquivo
   */
  static combineChunks(chunks: Chunk[]): string {
    if (chunks.length === 0) return '';
    
    const sortedChunks = chunks.sort((a, b) => a.metadata.chunkIndex - b.metadata.chunkIndex);
    
    switch (chunks[0].type) {
      case 'json':
        const jsonItems = sortedChunks.flatMap(chunk => chunk.content.items || []);
        return JSON.stringify(jsonItems, null, 2);
      
      case 'csv':
        const csvHeaders = sortedChunks[0]?.content.headers || [];
        const csvRows = sortedChunks.flatMap(chunk => chunk.content.rows || []);
        return Papa.unparse([csvHeaders, ...csvRows]);
      
      case 'txt':
        return sortedChunks.map(chunk => chunk.content.lines?.join('\n') || '').join('\n');
      
      case 'xml':
        // Para XML, seria necessário reconstruir a estrutura
        return '<!-- XML reconstruído -->';
      
      default:
        return sortedChunks.map(chunk => JSON.stringify(chunk.content)).join('\n');
    }
  }
}
