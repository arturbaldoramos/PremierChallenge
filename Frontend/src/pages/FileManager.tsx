import { useState, useRef } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { 
  Upload, 
  FileText, 
  FolderOpen,
  Plus,
  Trash2,
  Download,
  Eye
} from "lucide-react"

interface FileItem {
  id: string
  name: string
  size: number
  type: string
  category: string
  uploadDate: Date
  status: 'uploading' | 'processed' | 'error'
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
  const fileInputRef = useRef<HTMLInputElement>(null)

  // Função para processar arquivos (usada tanto para input quanto drag&drop)
  const processFiles = (fileList: FileList) => {
    if (!fileList || fileList.length === 0) return

    // Simular upload de múltiplos arquivos
    const newFiles: FileItem[] = Array.from(fileList).map((file, index) => ({
      id: `file-${Date.now()}-${index}`,
      name: file.name,
      size: file.size,
      type: file.type || 'unknown',
      category: selectedCategory,
      uploadDate: new Date(),
      status: 'uploading' as const
    }))

    setFiles(prev => [...prev, ...newFiles])

    // Simular processamento
    setTimeout(() => {
      setFiles(prev => prev.map(f => 
        f.status === 'uploading' ? { ...f, status: 'processed' as const } : f
      ))
    }, 2000)
  }

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

    return { totalFiles, processedFiles, totalSize, categories }
  }

  const stats = getStats()

  return (
    <div className="flex-1 space-y-4 p-4 md:p-8 pt-6">
      <div className="flex items-center justify-between space-y-2">
        <h2 className="text-3xl font-bold tracking-tight">Gerenciador de Arquivos</h2>
        <div className="flex items-center space-x-2">
          <Button variant="outline" size="sm" disabled={files.length === 0}>
            <Download className="mr-2 h-4 w-4" />
            Exportar Lista
          </Button>
        </div>
      </div>

      {/* Estatísticas */}
      <div className="grid gap-4 md:grid-cols-4">
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
              Categorias
            </CardTitle>
            <FolderOpen className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.categories}</div>
            <p className="text-xs text-muted-foreground">
              Tipos diferentes
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

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Status
            </CardTitle>
            <Eye className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {stats.processedFiles}/{stats.totalFiles}
            </div>
            <p className="text-xs text-muted-foreground">
              Processados
            </p>
          </CardContent>
        </Card>
      </div>

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
              <p className="text-muted-foreground mb-4">
                {isDragOver 
                  ? 'Arraste e solte os arquivos na categoria selecionada'
                  : `Arraste arquivos aqui ou selecione para a categoria: ${selectedCategory}`
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
                      <div key={file.id} className="flex items-center justify-between p-3 border rounded-lg">
                        <div className="flex items-center space-x-3">
                          <FileText className="h-6 w-6 text-primary" />
                          <div>
                            <p className="font-medium">{file.name}</p>
                            <div className="flex items-center space-x-4 text-sm text-muted-foreground">
                              <span>{(file.size / 1024).toFixed(2)} KB</span>
                              <span>{file.uploadDate.toLocaleDateString()}</span>
                              <Badge 
                                variant={file.status === 'processed' ? 'default' : 'secondary'}
                                className="text-xs"
                              >
                                {file.status === 'uploading' ? 'Processando...' : 
                                 file.status === 'processed' ? 'Processado' : 'Erro'}
                              </Badge>
                            </div>
                          </div>
                        </div>
                        <div className="flex items-center space-x-2">
                          {file.status === 'uploading' && (
                            <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-primary"></div>
                          )}
                          <Button variant="outline" size="sm" onClick={() => removeFile(file.id)}>
                            <Trash2 className="h-4 w-4" />
                          </Button>
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
