import { useState, useRef } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Progress } from "@/components/ui/progress"
import { 
  Upload, 
  FileText, 
  FolderOpen,
  Plus,
  Trash2,
  Download,
  Eye,
  CheckCircle,
  AlertCircle,
  FileSpreadsheet,
  FileImage,
  FileCode,
  Database,
  BarChart3
} from "lucide-react"
import { FileDetector, formatFileSize } from '@/lib/file-detector'
import type { FileTypeInfo } from '@/lib/file-detector'
import { FileChunker } from '@/lib/file-chunker'
import type { ProcessingResult, Chunk } from '@/lib/file-chunker'

interface FileItem {
  id: string
  name: string
  size: number
  type: string
  category: string
  uploadDate: Date
  status: 'pending' | 'processing' | 'processed' | 'error'
  progress: number
  fileType?: FileTypeInfo
  result?: ProcessingResult
  error?: string
  chunks?: Chunk[]
}

const fileCategories = [
  'Hospitais',
  'Médicos', 
  'Pacientes',
  'Equipamentos',
  'Medicamentos',
  'Laboratórios',
  'Clínicas',
  'Outros'
]

export default function FileManager() {
  const [files, setFiles] = useState<FileItem[]>([])
  const [selectedCategory, setSelectedCategory] = useState('Hospitais')
  const [isDragOver, setIsDragOver] = useState(false)
  const [selectedChunks, setSelectedChunks] = useState<Chunk[]>([])
  const [showStats, setShowStats] = useState(false)
  const fileInputRef = useRef<HTMLInputElement>(null)

  const getFileIcon = (fileType: FileTypeInfo) => {
    switch (fileType.type) {
      case 'pdf': return <FileImage className="h-6 w-6 text-red-500" />;
      case 'csv': return <FileSpreadsheet className="h-6 w-6 text-green-500" />;
      case 'json': return <FileCode className="h-6 w-6 text-yellow-500" />;
      case 'xml': return <FileCode className="h-6 w-6 text-blue-500" />;
      default: return <FileText className="h-6 w-6 text-gray-500" />;
    }
  };

  const processFiles = (fileList: FileList) => {
    if (!fileList || fileList.length === 0) return

     // Não há limite de tamanho - o sistema de chunks foi projetado para arquivos grandes
     const validFiles = Array.from(fileList);

    const newFiles: FileItem[] = validFiles.map((file, index) => ({
      id: `file-${Date.now()}-${index}`,
      name: file.name,
      size: file.size,
      type: file.type || 'unknown',
      category: selectedCategory,
      uploadDate: new Date(),
      status: 'pending' as const,
      progress: 0,
      fileType: FileDetector.detectFileType(file) || undefined
    }))

    setFiles(prev => [...prev, ...newFiles])

    // Processar arquivos automaticamente
    validFiles.forEach((file, index) => {
      const fileId = `file-${Date.now()}-${index}`;
      processFileWithChunking(file, fileId);
    });
  }

  // Função para processar arquivo com chunking
  const processFileWithChunking = async (file: File, fileId: string) => {
    const fileIndex = files.findIndex(f => f.id === fileId);
    
    // Atualiza status para processando
    setFiles(prev => prev.map((f, i) => 
      i === fileIndex ? { ...f, status: 'processing', progress: 0 } : f
    ));

    try {
      // Progresso mais granular para arquivos grandes
      const progressInterval = setInterval(() => {
        setFiles(prev => prev.map((f, i) => 
          i === fileIndex ? { 
            ...f, 
            progress: Math.min(f.progress + (file.size > 1024 * 1024 * 1024 ? 2 : 5), 95) // Mais lento para arquivos > 1GB
          } : f
        ));
      }, file.size > 1024 * 1024 * 1024 ? 500 : 200); // Intervalo maior para arquivos grandes

      // Processa com chunk size otimizado baseado no tamanho do arquivo
      const chunkSize = file.size > 1024 * 1024 * 1024 ? 50000 : 10000; // Chunks maiores para arquivos grandes
      const result = await FileChunker.processFile(file, chunkSize);
      
      clearInterval(progressInterval);

      // Atualiza com resultado
      setFiles(prev => prev.map((f, i) => 
        i === fileIndex ? { 
          ...f, 
          status: result.errors.length > 0 ? 'error' : 'processed',
          progress: 100,
          result,
          chunks: result.chunks,
          error: result.errors.join(', ')
        } : f
      ));

    } catch (error) {
      setFiles(prev => prev.map((f, i) => 
        i === fileIndex ? { 
          ...f, 
          status: 'error',
          progress: 0,
          error: error instanceof Error ? error.message : 'Erro desconhecido'
        } : f
      ));
    }
  };

  // Função para lidar com upload de arquivo via input
  const handleFileUpload = async (event: React.ChangeEvent<HTMLInputElement>) => {
    processFiles(event.target.files!)
    
    // Limpar input
    if (fileInputRef.current) {
      fileInputRef.current.value = ''
    }
  }

  // Funções para drag and drop
  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault()
    setIsDragOver(true)
  }

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault()
    setIsDragOver(false)
  }

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault()
    setIsDragOver(false)
    
    const droppedFiles = e.dataTransfer.files
    processFiles(droppedFiles)
  }

  // Função para remover arquivo
  const removeFile = (fileId: string) => {
    setFiles(prev => prev.filter(f => f.id !== fileId))
  }

  // Função para filtrar arquivos por categoria
  const getFilesByCategory = (category: string) => {
    return files.filter(f => f.category === category)
  }

  // Função para download de chunk
  const downloadChunk = (chunk: Chunk) => {
    const content = JSON.stringify(chunk.content, null, 2);
    const blob = new Blob([content], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `chunk-${chunk.id}.json`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  };

  // Função para download de todos os chunks
  const downloadAllChunks = () => {
    const allChunks = files
      .filter(f => f.chunks)
      .flatMap(f => f.chunks!);
    
    const content = JSON.stringify(allChunks, null, 2);
    const blob = new Blob([content], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'all-chunks.json';
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  };

  // Função para visualizar chunks
  const viewChunks = (chunks: Chunk[]) => {
    setSelectedChunks(chunks);
  };

  // Função para obter estatísticas
  const getStats = () => {
    const totalFiles = files.length
    const processedFiles = files.filter(f => f.status === 'processed').length
    const totalSize = files.reduce((sum, f) => sum + f.size, 0)
    const categories = [...new Set(files.map(f => f.category))].length
    const totalChunks = files.reduce((sum, f) => sum + (f.chunks?.length || 0), 0)
    const totalErrors = files.reduce((sum, f) => sum + (f.error ? 1 : 0), 0)

    return { totalFiles, processedFiles, totalSize, categories, totalChunks, totalErrors }
  }

  const stats = getStats()

  return (
    <div className="flex-1 space-y-4 p-4 md:p-8 pt-6">
      <div className="flex items-center justify-between space-y-2">
        <h2 className="text-3xl font-bold tracking-tight">Gerenciador de Arquivos</h2>
        <div className="flex items-center space-x-2">
          <Button 
            variant="outline" 
            size="sm"
            onClick={() => setShowStats(!showStats)}
          >
            <BarChart3 className="mr-2 h-4 w-4" />
            {showStats ? 'Ocultar' : 'Mostrar'} Estatísticas
          </Button>
          <Button 
            variant="outline" 
            size="sm" 
            disabled={stats.totalChunks === 0}
            onClick={downloadAllChunks}
          >
            <Download className="mr-2 h-4 w-4" />
            Baixar Chunks
          </Button>
        </div>
      </div>

      {/* Estatísticas */}
      {showStats && (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">
                Total de Arquivos
              </CardTitle>
              <FileText className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stats.totalFiles}</div>
              <p className="text-xs text-muted-foreground">
                {stats.processedFiles} processados
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">
                Total de Chunks
              </CardTitle>
              <Database className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stats.totalChunks}</div>
              <p className="text-xs text-muted-foreground">
                Chunks gerados
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">
                Erros
              </CardTitle>
              <AlertCircle className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stats.totalErrors}</div>
              <p className="text-xs text-muted-foreground">
                Erros encontrados
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">
                Tamanho Total
              </CardTitle>
              <Upload className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">
                {(stats.totalSize / 1024 / 1024).toFixed(1)} MB
              </div>
              <p className="text-xs text-muted-foreground">
                Espaço utilizado
              </p>
            </CardContent>
          </Card>
        </div>
      )}

      {/* Área de Upload */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Upload className="h-5 w-5" />
            Upload de Arquivos
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {/* Seleção de Categoria */}
            <div className="space-y-2">
              <label className="text-sm font-medium">Categoria do Arquivo</label>
              <div className="flex flex-wrap gap-2">
                {fileCategories.map((category) => (
                  <Button
                    key={category}
                    variant={selectedCategory === category ? "default" : "outline"}
                    size="sm"
                    onClick={() => setSelectedCategory(category)}
                  >
                    {category}
                  </Button>
                ))}
              </div>
            </div>

            {/* Área de Upload com Drag & Drop */}
            <div 
              className={`border-2 border-dashed rounded-lg p-8 text-center transition-colors ${
                isDragOver 
                  ? 'border-primary bg-primary/5' 
                  : 'border-muted-foreground/25 hover:border-muted-foreground/50'
              }`}
              onDragOver={handleDragOver}
              onDragLeave={handleDragLeave}
              onDrop={handleDrop}
            >
              <Upload className={`h-12 w-12 mx-auto mb-4 ${isDragOver ? 'text-primary' : 'text-muted-foreground'}`} />
              <h3 className="text-lg font-medium mb-2">
                {isDragOver ? 'Solte os arquivos aqui' : 'Faça upload de arquivos'}
              </h3>
              <p className="text-muted-foreground mb-2">
                {isDragOver 
                  ? 'Arraste e solte os arquivos na categoria selecionada'
                  : `Arraste arquivos aqui ou selecione para a categoria: ${selectedCategory}`
                }
              </p>
               <p className="text-xs text-muted-foreground mb-4">
                 Tipos suportados: XML, JSON, PDF, CSV, TXT • Sem limite de tamanho - processamento em chunks inteligentes
               </p>
              <div className="flex items-center justify-center gap-4">
                <Button onClick={() => fileInputRef.current?.click()}>
                  <Plus className="mr-2 h-4 w-4" />
                  Selecionar Arquivos
                </Button>
              </div>
              <input
                ref={fileInputRef}
                type="file"
                accept="*/*"
                multiple
                onChange={handleFileUpload}
                className="hidden"
              />
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Lista de Arquivos por Categoria */}
      {files.length > 0 && (
        <div className="space-y-4">
          {fileCategories.map((category) => {
            const categoryFiles = getFilesByCategory(category)
            if (categoryFiles.length === 0) return null

            return (
              <Card key={category}>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <FolderOpen className="h-5 w-5" />
                    {category}
                    <Badge variant="outline">{categoryFiles.length}</Badge>
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-2">
                    {categoryFiles.map((file) => (
                      <div key={file.id} className="border rounded-lg p-4">
                        <div className="flex items-center justify-between mb-2">
                          <div className="flex items-center gap-3">
                            {file.fileType ? getFileIcon(file.fileType) : <FileText className="h-6 w-6 text-gray-500" />}
                            <div>
                              <p className="font-medium">{file.name}</p>
                              <p className="text-sm text-gray-500">
                                {formatFileSize(file.size)} • {file.fileType?.description || 'Tipo desconhecido'}
                              </p>
                            </div>
                          </div>
                          <div className="flex items-center gap-2">
                            {file.status === 'pending' && (
                              <Badge variant="secondary" className="text-xs">
                                Aguardando
                              </Badge>
                            )}
                            {file.status === 'processing' && (
                              <Badge variant="secondary" className="text-xs">
                                Processando...
                              </Badge>
                            )}
                            {file.status === 'processed' && file.chunks && (
                              <div className="flex gap-2">
                                <Button 
                                  size="sm" 
                                  variant="outline"
                                  onClick={() => viewChunks(file.chunks!)}
                                >
                                  <Eye className="h-4 w-4 mr-1" />
                                  Ver Chunks ({file.chunks.length})
                                </Button>
                                <Button 
                                  size="sm" 
                                  variant="outline"
                                  onClick={() => downloadChunk(file.chunks![0])}
                                >
                                  <Download className="h-4 w-4 mr-1" />
                                  Baixar
                                </Button>
                              </div>
                            )}
                            {file.status === 'error' && (
                              <Badge variant="destructive" className="text-xs">
                                Erro
                              </Badge>
                            )}
                            <Button 
                              size="sm" 
                              variant="outline"
                              onClick={() => removeFile(file.id)}
                            >
                              <Trash2 className="h-4 w-4" />
                            </Button>
                          </div>
                        </div>

                        {/* Status e Progresso */}
                        <div className="space-y-2">
                          <div className="flex items-center gap-2">
                            {file.status === 'pending' && (
                              <>
                                <div className="w-2 h-2 bg-gray-400 rounded-full"></div>
                                <span className="text-sm text-gray-600">Aguardando processamento</span>
                              </>
                            )}
                            {file.status === 'processing' && (
                              <>
                                <div className="w-2 h-2 bg-blue-500 rounded-full animate-pulse"></div>
                                <span className="text-sm text-blue-600">Processando...</span>
                              </>
                            )}
                            {file.status === 'processed' && (
                              <>
                                <CheckCircle className="h-4 w-4 text-green-500" />
                                <span className="text-sm text-green-600">
                                  Processado ({file.chunks?.length} chunks)
                                </span>
                              </>
                            )}
                            {file.status === 'error' && (
                              <>
                                <AlertCircle className="h-4 w-4 text-red-500" />
                                <span className="text-sm text-red-600">
                                  Erro: {file.error}
                                </span>
                              </>
                            )}
                          </div>
                          
                          {file.status === 'processing' && (
                            <Progress value={file.progress} className="h-2" />
                          )}
                        </div>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            )
          })}
        </div>
      )}

      {/* Visualizador de Chunks */}
      {selectedChunks.length > 0 && (
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle>Chunks Selecionados ({selectedChunks.length})</CardTitle>
              <Button 
                variant="outline" 
                onClick={() => setSelectedChunks([])}
              >
                Fechar
              </Button>
            </div>
          </CardHeader>
          <CardContent>
            <div className="space-y-4 max-h-96 overflow-y-auto">
              {selectedChunks.map((chunk) => (
                <div key={chunk.id} className="border rounded-lg p-4">
                  <div className="flex items-center justify-between mb-2">
                    <h4 className="font-medium">Chunk {chunk.metadata.chunkIndex + 1}</h4>
                    <div className="flex gap-2">
                      <Button 
                        size="sm" 
                        variant="outline"
                        onClick={() => downloadChunk(chunk)}
                      >
                        <Download className="h-4 w-4 mr-1" />
                        Baixar
                      </Button>
                    </div>
                  </div>
                  <div className="text-sm text-gray-600 mb-2">
                    <p>Tamanho: {formatFileSize(chunk.metadata.size)}</p>
                    <p>Processado em: {new Date(chunk.metadata.timestamp).toLocaleString()}</p>
                  </div>
                  <div className="bg-gray-50 rounded p-3 max-h-32 overflow-y-auto">
                    <pre className="text-xs">
                      {JSON.stringify(chunk.content, null, 2)}
                    </pre>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      )}

      {/* Estado vazio */}
      {files.length === 0 && (
        <Card>
          <CardContent className="pt-6">
            <div className="text-center py-8">
              <FolderOpen className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
              <h3 className="text-lg font-medium mb-2">Nenhum arquivo carregado</h3>
              <p className="text-muted-foreground mb-4">
                Faça upload de arquivos para começar a gerenciar seus dados.
              </p>
              <Button onClick={() => fileInputRef.current?.click()}>
                <Plus className="mr-2 h-4 w-4" />
                Fazer Primeiro Upload
              </Button>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  )
}
