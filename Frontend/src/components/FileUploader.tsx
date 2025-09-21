import React, { useCallback, useState } from 'react';
import { Upload, File, AlertCircle, CheckCircle2, Loader2, X, FileText, FileSpreadsheet, FileCode, FileImage, FileArchive } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Progress } from '@/components/ui/progress';
import { Badge } from '@/components/ui/badge';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { useWebSocket } from '@/hooks/useWebSocket';

export interface FileType {
  id: string;
  label: string;
  description: string;
  acceptedTypes: string[];
}

export const FILE_TYPES: FileType[] = [
  {
    id: 'estados',
    label: 'Estados',
    description: 'Dados dos estados brasileiros',
    acceptedTypes: ['.csv', '.xlsx', '.json', '.xml', '.txt']
  },
  {
    id: 'municipios',
    label: 'Municípios',
    description: 'Dados dos municípios brasileiros',
    acceptedTypes: ['.csv', '.xlsx', '.json', '.xml', '.txt']
  },
  {
    id: 'medicos',
    label: 'Médicos',
    description: 'Cadastro de médicos',
    acceptedTypes: ['.csv', '.xlsx', '.json', '.xml', '.txt']
  },
  {
    id: 'hospitais',
    label: 'Hospitais',
    description: 'Cadastro de hospitais',
    acceptedTypes: ['.csv', '.xlsx', '.json', '.xml', '.txt']
  },
  {
    id: 'pacientes',
    label: 'Pacientes',
    description: 'Cadastro de pacientes',
    acceptedTypes: ['.csv', '.xlsx', '.json', '.xml', '.txt']
  },
  {
    id: 'cid10',
    label: 'CID-10',
    description: 'Códigos de classificação de doenças',
    acceptedTypes: ['.csv', '.xlsx', '.json', '.xml', '.txt']
  }
];

interface FileUploaderProps {
  websocketUrl?: string;
  onUploadComplete?: (fileType: string, fileName: string) => void;
  onError?: (error: string) => void;
}

// Função para detectar a URL correta do WebSocket
const getWebSocketUrl = (): string => {
  // Se está rodando em desenvolvimento (npm run dev)
  if (window.location.hostname === 'localhost' && window.location.port === '5173') {
    return 'ws://localhost:8080/ws';
  }

  // Se está rodando via Docker ou produção, usar proxy do Nginx
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  const hostname = window.location.hostname;
  const port = window.location.port ? `:${window.location.port}` : '';
  return `${protocol}//${hostname}${port}/ws`;
};

export const FileUploader: React.FC<FileUploaderProps> = ({
  websocketUrl = getWebSocketUrl(),
  onUploadComplete,
  onError
}) => {
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [selectedFileType, setSelectedFileType] = useState<string>('estados');
  const [isUploading, setIsUploading] = useState(false);
  const [uploadProgress, setUploadProgress] = useState(0);
  const [isDragOver, setIsDragOver] = useState(false);

  console.log('FileUploader renderizado com URL:', websocketUrl);

  const {
    isConnected,
    isConnecting,
    sessionId,
    connect,
    disconnect,
    uploadFile,
    getProgress,
    getQueueStatus,
    progressData,
    queueStatus,
    error: wsError
  } = useWebSocket({
    url: websocketUrl,
    autoConnect: true
  });

  // Função para detectar tipo de arquivo e obter ícone
  const getFileIcon = useCallback((fileName: string) => {
    const extension = fileName.split('.').pop()?.toLowerCase();
    
    switch (extension) {
      case 'csv':
        return <FileSpreadsheet className="w-4 h-4 text-green-500" />;
      case 'xlsx':
      case 'xls':
        return <FileSpreadsheet className="w-4 h-4 text-green-600" />;
      case 'json':
        return <FileCode className="w-4 h-4 text-yellow-500" />;
      case 'xml':
        return <FileCode className="w-4 h-4 text-blue-500" />;
      case 'txt':
        return <FileText className="w-4 h-4 text-muted-foreground" />;
      case 'pdf':
        return <FileImage className="w-4 h-4 text-red-500" />;
      case 'zip':
      case 'rar':
      case '7z':
        return <FileArchive className="w-4 h-4 text-purple-500" />;
      default:
        return <File className="w-4 h-4 text-muted-foreground" />;
    }
  }, []);

  // Função para obter descrição do tipo de arquivo
  const getFileTypeDescription = useCallback((fileName: string) => {
    const extension = fileName.split('.').pop()?.toLowerCase();
    
    switch (extension) {
      case 'csv':
        return 'Arquivo CSV (Comma Separated Values)';
      case 'xlsx':
      case 'xls':
        return 'Planilha Excel';
      case 'json':
        return 'Arquivo JSON (JavaScript Object Notation)';
      case 'xml':
        return 'Arquivo XML (eXtensible Markup Language)';
      case 'txt':
        return 'Arquivo de texto';
      case 'pdf':
        return 'Documento PDF (Portable Document Format)';
      case 'zip':
      case 'rar':
      case '7z':
        return 'Arquivo compactado';
      default:
        return 'Arquivo desconhecido';
    }
  }, []);

  const handleFileSelect = useCallback((event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file) {
      setSelectedFile(file);
    }
  }, []);

  // Funções para drag and drop
  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragOver(true);
  }, []);

  const handleDragLeave = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragOver(false);
  }, []);

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragOver(false);
    
    const droppedFiles = e.dataTransfer.files;
    if (droppedFiles.length > 0) {
      setSelectedFile(droppedFiles[0]);
    }
  }, []);

  const handleUpload = useCallback(async () => {
    console.log('handleUpload chamado');
    console.log('selectedFile:', selectedFile);
    console.log('selectedFileType:', selectedFileType);
    console.log('isConnected:', isConnected);

    if (!selectedFile || !selectedFileType || !isConnected) {
      console.log('Upload bloqueado - arquivo, tipo ou conexão ausentes');
      return;
    }

    console.log('Iniciando upload...');
    setIsUploading(true);
    setUploadProgress(0);

    try {
      await uploadFile(selectedFile, selectedFileType);
      setUploadProgress(100);
      onUploadComplete?.(selectedFileType, selectedFile.name);

      // Limpar seleção após upload bem-sucedido
      setSelectedFile(null);
      setSelectedFileType('');

      // Reset do input file
      const fileInput = document.getElementById('file-input') as HTMLInputElement;
      if (fileInput) {
        fileInput.value = '';
      }
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Erro desconhecido no upload';
      onError?.(errorMessage);
    } finally {
      setIsUploading(false);
    }
  }, [selectedFile, selectedFileType, isConnected, uploadFile, onUploadComplete, onError]);

  const handleClearFile = useCallback(() => {
    setSelectedFile(null);
    setSelectedFileType('');
    const fileInput = document.getElementById('file-input') as HTMLInputElement;
    if (fileInput) {
      fileInput.value = '';
    }
  }, []);

  const formatFileSize = (bytes: number): string => {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  const getStatusBadge = () => {
    if (isConnecting) {
      return <Badge variant="outline"><Loader2 className="w-3 h-3 mr-1 animate-spin" />Conectando</Badge>;
    }
    if (isConnected) {
      return <Badge variant="default"><CheckCircle2 className="w-3 h-3 mr-1" />Conectado</Badge>;
    }
    return <Badge variant="destructive"><AlertCircle className="w-3 h-3 mr-1" />Desconectado</Badge>;
  };

  const selectedFileTypeData = FILE_TYPES.find(ft => ft.id === selectedFileType);

  return (
    <Card className="w-full max-w-2xl">
      <CardHeader>
        <div className="flex justify-between items-center">
          <CardTitle className="flex items-center gap-2">
            <Upload className="w-5 h-5" />
            Upload de Arquivos
          </CardTitle>
          {getStatusBadge()}
        </div>
      </CardHeader>
      <CardContent className="space-y-6">
        {/* Controles de Conexão */}
        <div className="flex gap-2">
          <Button
            onClick={connect}
            disabled={isConnected || isConnecting}
            variant="outline"
            size="sm"
          >
            {isConnecting ? <Loader2 className="w-4 h-4 mr-2 animate-spin" /> : null}
            Conectar
          </Button>
          <Button
            onClick={disconnect}
            disabled={!isConnected}
            variant="outline"
            size="sm"
          >
            Desconectar
          </Button>
          <Button
            onClick={getQueueStatus}
            disabled={!isConnected}
            variant="outline"
            size="sm"
          >
            Status das Filas
          </Button>
        </div>

        {/* Seleção do Tipo de Arquivo */}
        <div className="space-y-2">
          <label className="text-sm font-medium">Tipo de Arquivo</label>
          <div className="flex flex-wrap gap-2">
            {FILE_TYPES.map((fileType) => (
              <Button
                key={fileType.id}
                variant={selectedFileType === fileType.id ? "default" : "outline"}
                size="sm"
                onClick={() => setSelectedFileType(fileType.id)}
                className="text-xs"
              >
                {fileType.label}
              </Button>
            ))}
          </div>
          {selectedFileTypeData && (
            <p className="text-xs text-muted-foreground">
              {selectedFileTypeData.description}
            </p>
          )}
        </div>

        {/* Seleção de Arquivo */}
        <div className="space-y-2">
          <label className="text-sm font-medium">Arquivo</label>
          <div 
            className={`border-2 border-dashed rounded-lg p-6 text-center transition-colors ${
              isDragOver
                ? 'border-primary bg-primary/5'
                : 'border-border hover:border-border/80'
            }`}
            onDragOver={handleDragOver}
            onDragLeave={handleDragLeave}
            onDrop={handleDrop}
          >
            <input
              id="file-input"
              type="file"
              accept="*/*"
              onChange={handleFileSelect}
              className="hidden"
            />
            <label htmlFor="file-input" className="cursor-pointer">
              <Upload className={`w-8 h-8 mx-auto mb-2 ${isDragOver ? 'text-primary' : 'text-muted-foreground'}`} />
              <p className={`text-sm ${isDragOver ? 'text-primary' : 'text-muted-foreground'}`}>
                {isDragOver ? 'Solte o arquivo aqui' : 'Clique para selecionar ou arraste um arquivo'}
              </p>
              <p className="text-xs text-muted-foreground mt-1">
                Tipos aceitos: CSV, XLSX, JSON, XML, TXT, PDF e outros
              </p>
            </label>
          </div>
        </div>

        {/* Arquivo Selecionado */}
        {selectedFile && (
          <div className="flex items-center justify-between p-3 bg-muted/50 rounded-lg">
            <div className="flex items-center gap-3">
              {getFileIcon(selectedFile.name)}
              <div>
                <p className="text-sm font-medium">{selectedFile.name}</p>
                <p className="text-xs text-muted-foreground">
                  {formatFileSize(selectedFile.size)} • {getFileTypeDescription(selectedFile.name)}
                </p>
              </div>
            </div>
            <Button
              onClick={handleClearFile}
              variant="ghost"
              size="sm"
              className="h-8 w-8 p-0"
            >
              <X className="w-4 h-4" />
            </Button>
          </div>
        )}

        {/* Progress Bar */}
        {(isUploading || uploadProgress > 0) && (
          <div className="space-y-2">
            <div className="flex justify-between text-sm">
              <span>Progresso do Upload</span>
              <span>{uploadProgress}%</span>
            </div>
            <Progress value={uploadProgress} className="w-full" />
          </div>
        )}

        {/* Botão de Upload */}
        <Button
          onClick={handleUpload}
          disabled={!selectedFile || !selectedFileType || !isConnected || isUploading}
          className="w-full"
        >
          {isUploading ? (
            <>
              <Loader2 className="w-4 h-4 mr-2 animate-spin" />
              Enviando...
            </>
          ) : (
            <>
              <Upload className="w-4 h-4 mr-2" />
              Fazer Upload
            </>
          )}
        </Button>

        {/* Mensagens de Erro */}
        {wsError && (
          <Alert variant="destructive">
            <AlertCircle className="h-4 w-4" />
            <AlertDescription>{wsError}</AlertDescription>
          </Alert>
        )}

        {/* Status das Filas */}
        {queueStatus && Object.keys(queueStatus).length > 0 && (
          <div className="space-y-2">
            <h4 className="text-sm font-medium">Status das Filas</h4>
            <div className="grid grid-cols-2 gap-2">
              {Object.entries(queueStatus).map(([queue, count]) => (
                <div key={queue} className="flex justify-between p-2 bg-muted/50 rounded">
                  <span className="text-sm capitalize">{queue}</span>
                  <Badge variant="outline">{count} jobs</Badge>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Progresso dos Jobs */}
        {progressData.length > 0 && (
          <div className="space-y-2">
            <h4 className="text-sm font-medium">Progresso dos Jobs</h4>
            <div className="space-y-2">
              {progressData.map((job) => (
                <div key={job.job_id} className="p-3 border rounded-lg">
                  <div className="flex justify-between items-center mb-2">
                    <span className="text-sm font-medium capitalize">{job.type}</span>
                    <Badge variant={job.status === 'completed' ? 'default' : job.status === 'failed' ? 'destructive' : 'secondary'}>
                      {job.status}
                    </Badge>
                  </div>
                  <div className="space-y-1">
                    <div className="flex justify-between text-xs text-muted-foreground">
                      <span>{job.processed_items} de {job.total_items} processados</span>
                      <span>{job.progress}%</span>
                    </div>
                    <Progress value={job.progress} className="h-2" />
                    {job.message && (
                      <p className="text-xs text-muted-foreground">{job.message}</p>
                    )}
                  </div>
                </div>
              ))}
            </div>
            <Button
              onClick={getProgress}
              variant="outline"
              size="sm"
              className="w-full"
            >
              Atualizar Progresso
            </Button>
          </div>
        )}

        {/* Informações da Sessão */}
        {sessionId && (
          <div className="text-xs text-muted-foreground">
            Session ID: {sessionId}
          </div>
        )}
      </CardContent>
    </Card>
  );
};