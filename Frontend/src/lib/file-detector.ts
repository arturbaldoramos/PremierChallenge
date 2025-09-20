// Utilitário para detecção de tipos de arquivo
export interface FileTypeInfo {
  type: string;
  extension: string;
  mimeType: string;
  description: string;
}

export const SUPPORTED_FILE_TYPES: FileTypeInfo[] = [
  {
    type: 'xml',
    extension: '.xml',
    mimeType: 'application/xml',
    description: 'Extensible Markup Language'
  },
  {
    type: 'json',
    extension: '.json',
    mimeType: 'application/json',
    description: 'JavaScript Object Notation'
  },
  {
    type: 'pdf',
    extension: '.pdf',
    mimeType: 'application/pdf',
    description: 'Portable Document Format'
  },
  {
    type: 'csv',
    extension: '.csv',
    mimeType: 'text/csv',
    description: 'Comma-Separated Values'
  },
  {
    type: 'txt',
    extension: '.txt',
    mimeType: 'text/plain',
    description: 'Plain Text'
  },
  {
    type: 'xlsx',
    extension: '.xlsx',
    mimeType: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
    description: 'Excel Spreadsheet'
  },
  {
    type: 'docx',
    extension: '.docx',
    mimeType: 'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
    description: 'Word Document'
  }
];

export class FileDetector {
  /**
   * Detecta o tipo de arquivo baseado no nome e MIME type
   */
  static detectFileType(file: File): FileTypeInfo | null {
    const fileName = file.name.toLowerCase();
    const mimeType = file.type.toLowerCase();
    
    // Primeiro, tenta detectar pela extensão
    for (const fileType of SUPPORTED_FILE_TYPES) {
      if (fileName.endsWith(fileType.extension)) {
        return fileType;
      }
    }
    
    // Se não encontrou pela extensão, tenta pelo MIME type
    for (const fileType of SUPPORTED_FILE_TYPES) {
      if (mimeType.includes(fileType.type) || mimeType === fileType.mimeType) {
        return fileType;
      }
    }
    
    return null;
  }

  /**
   * Detecta o tipo de arquivo de forma assíncrona (incluindo análise de conteúdo)
   */
  static async detectFileTypeAsync(file: File): Promise<FileTypeInfo | null> {
    const fileName = file.name.toLowerCase();
    const mimeType = file.type.toLowerCase();
    
    // Primeiro, tenta detectar pela extensão
    for (const fileType of SUPPORTED_FILE_TYPES) {
      if (fileName.endsWith(fileType.extension)) {
        return fileType;
      }
    }
    
    // Se não encontrou pela extensão, tenta pelo MIME type
    for (const fileType of SUPPORTED_FILE_TYPES) {
      if (mimeType.includes(fileType.type) || mimeType === fileType.mimeType) {
        return fileType;
      }
    }
    
    // Detecção adicional por conteúdo para alguns tipos
    return await this.detectByContent(file);
  }

  /**
   * Detecta o tipo de arquivo analisando o conteúdo (primeiros bytes)
   */
  private static async detectByContent(file: File): Promise<FileTypeInfo | null> {
    return new Promise((resolve) => {
      const reader = new FileReader();
      reader.onload = (e) => {
        const content = e.target?.result as string;
        
        // Detecta XML
        if (content.trim().startsWith('<?xml') || content.trim().startsWith('<')) {
          resolve(SUPPORTED_FILE_TYPES.find(ft => ft.type === 'xml') || null);
          return;
        }
        
        // Detecta JSON
        if (content.trim().startsWith('{') || content.trim().startsWith('[')) {
          try {
            JSON.parse(content);
            resolve(SUPPORTED_FILE_TYPES.find(ft => ft.type === 'json') || null);
            return;
          } catch {
            // Não é JSON válido
          }
        }
        
        // Detecta CSV (linhas separadas por vírgula)
        if (content.includes(',') && content.includes('\n')) {
          const lines = content.split('\n');
          if (lines.length > 1 && lines[0].includes(',')) {
            resolve(SUPPORTED_FILE_TYPES.find(ft => ft.type === 'csv') || null);
            return;
          }
        }
        
        resolve(null);
      };
      
      // Lê apenas os primeiros 1024 bytes para detecção
      reader.readAsText(file.slice(0, 1024));
    });
  }

  /**
   * Valida se o arquivo é suportado
   */
  static isSupportedFile(file: File): boolean {
    const fileType = this.detectFileType(file);
    return fileType !== null;
  }

  /**
   * Obtém informações sobre um tipo de arquivo específico
   */
  static getFileTypeInfo(type: string): FileTypeInfo | null {
    return SUPPORTED_FILE_TYPES.find(ft => ft.type === type) || null;
  }

  /**
   * Lista todos os tipos de arquivo suportados
   */
  static getSupportedTypes(): FileTypeInfo[] {
    return [...SUPPORTED_FILE_TYPES];
  }
}

// Função utilitária para formatar tamanho de arquivo
export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 Bytes';
  
  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}
