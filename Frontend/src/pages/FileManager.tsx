import { useState } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Progress } from "@/components/ui/progress"
import {
  FileText,
  FolderOpen,
  Trash2,
  Download,
  Eye,
  CheckCircle,
  AlertCircle,
  FileSpreadsheet,
  FileImage,
  FileCode
} from "lucide-react"
import { formatFileSize } from '@/lib/file-detector'
import type { FileTypeInfo } from '@/lib/file-detector'
import { FileUploader } from '@/components/FileUploader'

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
  const [uploadMessage, setUploadMessage] = useState<string | null>(null)

  const getFileIcon = (fileType: FileTypeInfo) => {
    switch (fileType.type) {
      case 'pdf': return <FileImage className="h-6 w-6 text-red-500" />;
      case 'csv': return <FileSpreadsheet className="h-6 w-6 text-green-500" />;
      case 'json': return <FileCode className="h-6 w-6 text-yellow-500" />;
      case 'xml': return <FileCode className="h-6 w-6 text-blue-500" />;
      default: return <FileText className="h-6 w-6 text-muted-foreground" />;
    }
  };


  // Função para remover arquivo
  const removeFile = (fileId: string) => {
    setFiles(prev => prev.filter(f => f.id !== fileId))
  }

  // Função para filtrar arquivos por categoria
  const getFilesByCategory = (category: string) => {
    return files.filter(f => f.category === category)
  }




  // Callbacks para o FileUploader
  const handleUploadComplete = (fileType: string, fileName: string) => {
    setUploadMessage(`✅ Arquivo ${fileName} (${fileType}) enviado com sucesso!`)
    setTimeout(() => setUploadMessage(null), 5000)

    // Adicionar arquivo à lista como processado
    const newFile: FileItem = {
      id: `ws-file-${Date.now()}`,
      name: fileName,
      size: 0, // Tamanho não conhecido via WebSocket
      type: 'application/csv',
      category: fileType.charAt(0).toUpperCase() + fileType.slice(1),
      uploadDate: new Date(),
      status: 'processed',
      progress: 100,
      fileType: { type: 'csv', extension: 'csv', mimeType: 'text/csv', description: 'CSV File' }
    }
    setFiles(prev => [...prev, newFile])
  }

  const handleUploadError = (error: string) => {
    setUploadMessage(`❌ Erro no upload: ${error}`)
    setTimeout(() => setUploadMessage(null), 5000)
  }

  return (
    <div className="flex-1 space-y-4 p-4 md:p-8 pt-6">
      <div className="flex items-center justify-between space-y-2">
        <h2 className="text-3xl font-bold tracking-tight">Gerenciador de Arquivos</h2>
      </div>


      {/* Mensagem de Upload */}
      {uploadMessage && (
        <Card className="border-green-200 bg-green-50">
          <CardContent className="pt-6">
            <p className="text-center text-green-800">{uploadMessage}</p>
          </CardContent>
        </Card>
      )}

      {/* Componente de Upload WebSocket */}
      <FileUploader
        onUploadComplete={handleUploadComplete}
        onError={handleUploadError}
      />

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
                            {file.fileType ? getFileIcon(file.fileType) : <FileText className="h-6 w-6 text-muted-foreground" />}
                            <div>
                              <p className="font-medium">{file.name}</p>
                              <p className="text-sm text-muted-foreground">
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
                        <div className="w-2 h-2 bg-muted-foreground rounded-full"></div>
                        <span className="text-sm text-muted-foreground">Aguardando processamento</span>
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
                Use o componente de upload acima para enviar arquivos CSV
              </p>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  )
}
