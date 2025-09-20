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
  error?: string
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

     // Não há limite de tamanho - arquivos são enviados diretamente para o backend
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
      processFileSimple(file, fileId);
    });
  }


  // Função para processar arquivo
  const processFileSimple = async (file: File, fileId: string) => {
    const fileIndex = files.findIndex(f => f.id === fileId);
    
    // Atualiza status para processando
    setFiles(prev => prev.map((f, i) => 
      i === fileIndex ? { ...f, status: 'processing', progress: 0 } : f
    ));

    try {
      // Simula progresso de upload
      const progressInterval = setInterval(() => {
        setFiles(prev => prev.map((f, i) => 
          i === fileIndex ? { 
            ...f, 
            progress: Math.min(f.progress + 10, 90)
          } : f
        ));
      }, 200);

      // Detecta tipo de arquivo
      const fileType = FileDetector.detectFileType(file);
      
      // Simula envio para backend
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      clearInterval(progressInterval);

      // Atualiza com resultado
      setFiles(prev => prev.map((f, i) => 
        i === fileIndex ? { 
          ...f, 
          status: 'processed',
          progress: 100,
          fileType: fileType || { type: 'unknown', extension: '', mimeType: '', description: 'Unknown' },
          error: undefined
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



  // Função para obter estatísticas
  const getStats = () => {
    const totalFiles = files.length
    const processedFiles = files.filter(f => f.status === 'processed').length
    const totalSize = files.reduce((sum, f) => sum + f.size, 0)
    const categories = [...new Set(files.map(f => f.category))].length
    const totalErrors = files.reduce((sum, f) => sum + (f.error ? 1 : 0), 0)

    return { totalFiles, processedFiles, totalSize, categories, totalErrors }
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
                Arquivos Processados
              </CardTitle>
              <Database className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">
                {stats.processedFiles}
              </div>
              <p className="text-xs text-muted-foreground">
                Arquivos processados
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
                 Tipos suportados: XML, JSON, PDF, CSV, TXT • 
                 {selectedCategory === 'Pacientes' 
                   ? 'Sem limite de tamanho - processamento em chunks inteligentes'
                   : 'Detecção automática de tipo e envio para backend'
                 }
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
                            {file.status === 'processed' && (
                              <div className="flex gap-2">
                                <Button 
                                  size="sm" 
                                  variant="outline"
                                  onClick={() => {/* TODO: Implementar visualização de dados */}}
                                >
                                  <Eye className="h-4 w-4 mr-1" />
                                  Ver Dados
                                </Button>
                                <Button 
                                  size="sm" 
                                  variant="outline"
                                  onClick={() => {/* TODO: Implementar download */}}
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
                                  Processado e enviado para backend
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
