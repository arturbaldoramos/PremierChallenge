import { useState, useEffect, useRef } from "react"
import { useNavigate } from "react-router-dom"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Badge } from "@/components/ui/badge"
import { 
  ChevronLeft, 
  ChevronRight, 
  ChevronsLeft, 
  ChevronsRight 
} from "lucide-react"
import { 
  User, 
  MapPin, 
  Stethoscope, 
  Search, 
  Filter,
  Download,
  Upload,
  FileText,
  FolderOpen,
  X
} from "lucide-react"
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover"
import { Check, ChevronsUpDown } from "lucide-react"
import { cn } from "@/lib/utils"

interface Doctor {
  codigo: string
  nome_completo: string
  especialidade: string
  cidade: number
}

interface FilterState {
  cidade: string
  especialidade: string
  search: string
}

export default function Doctors() {
  const navigate = useNavigate()
  const [doctors, setDoctors] = useState<Doctor[]>([])
  const [filteredDoctors, setFilteredDoctors] = useState<Doctor[]>([])
  const [loading] = useState(false)
  const [uploading, setUploading] = useState(false)
  const [uploadedFile, setUploadedFile] = useState<File | null>(null)
  const [isDragOver, setIsDragOver] = useState(false)
  const [filters, setFilters] = useState<FilterState>({
    cidade: '',
    especialidade: '',
    search: ''
  })
  const [openSpecialty, setOpenSpecialty] = useState(false)
  const fileInputRef = useRef<HTMLInputElement>(null)

  // Estados para paginação
  const [currentPage, setCurrentPage] = useState(1)
  const [itemsPerPage] = useState(6) // 6 médicos por página

  // Configuração da API - substitua pela sua URL
  const API_URL = 'https://sua-api.com/upload'

  // Dados mockados para teste
  const mockDoctors: Doctor[] = [
    {
      codigo: "M001",
      nome_completo: "Dr. João Silva Santos",
      especialidade: "Cardiologia",
      cidade: 1
    },
    {
      codigo: "M002", 
      nome_completo: "Dra. Maria Oliveira Costa",
      especialidade: "Neurologia",
      cidade: 1
    },
    {
      codigo: "M003",
      nome_completo: "Dr. Carlos Eduardo Ferreira",
      especialidade: "Pediatria",
      cidade: 2
    },
    {
      codigo: "M004",
      nome_completo: "Dra. Ana Paula Rodrigues",
      especialidade: "Cardiologia",
      cidade: 2
    },
    {
      codigo: "M005",
      nome_completo: "Dr. Roberto Almeida Lima",
      especialidade: "Ortopedia",
      cidade: 1
    },
    {
      codigo: "M006",
      nome_completo: "Dra. Fernanda Souza Martins",
      especialidade: "Neurologia",
      cidade: 3
    },
    {
      codigo: "M007",
      nome_completo: "Dr. Pedro Henrique Gomes",
      especialidade: "Pediatria",
      cidade: 1
    },
    {
      codigo: "M008",
      nome_completo: "Dra. Juliana Mendes Pereira",
      especialidade: "Cardiologia",
      cidade: 2
    },
    {
      codigo: "M009",
      nome_completo: "Dr. Rafael Barbosa Silva",
      especialidade: "Ortopedia",
      cidade: 3
    },
    {
      codigo: "M010",
      nome_completo: "Dra. Camila Santos Oliveira",
      especialidade: "Neurologia",
      cidade: 1
    },
    {
      codigo: "M011",
      nome_completo: "Dr. Lucas Ferreira Costa",
      especialidade: "Pediatria",
      cidade: 2
    },
    {
      codigo: "M012",
      nome_completo: "Dra. Beatriz Alves Rodrigues",
      especialidade: "Cardiologia",
      cidade: 3
    }
  ]

  // Especialidades médicas (será preenchida dinamicamente com os dados)
  const [specialties, setSpecialties] = useState<string[]>([])

  // Função para carregar dados mockados (para teste)
  const loadMockData = () => {
    setDoctors(mockDoctors)
    setFilteredDoctors(mockDoctors)
    
    // Extrair especialidades únicas dos dados mockados
    const uniqueSpecialties = [...new Set(mockDoctors.map(d => d.especialidade))]
    setSpecialties(uniqueSpecialties)
  }

  // Função para processar arquivo (usada tanto para input quanto drag&drop)
  const processFile = async (file: File) => {
    if (!file) return

    setUploading(true)
    setUploadedFile(file)

    try {
      const formData = new FormData()
      formData.append('file', file)
      
      const response = await fetch(API_URL, {
        method: 'POST',
        body: formData
      })
      
      if (!response.ok) {
        throw new Error(`Erro na API: ${response.status}`)
      }
      
      const data = await response.json()
      setDoctors(data)
      setFilteredDoctors(data)
      
      // Extrair especialidades únicas dos dados
      const uniqueSpecialties = [...new Set(data.map((d: Doctor) => d.especialidade))] as string[]
      setSpecialties(uniqueSpecialties)
      
    } catch (error) {
      console.error('Erro ao fazer upload do arquivo:', error)
      alert('Erro ao processar arquivo. Verifique a configuração da API.')
    } finally {
      setUploading(false)
    }
  }

  // Função para lidar com upload de arquivo via input
  const handleFileUpload = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0]
    if (!file) return

    await processFile(file)
    
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

  const handleDrop = async (e: React.DragEvent) => {
    e.preventDefault()
    setIsDragOver(false)
    
    const droppedFiles = e.dataTransfer.files
    if (droppedFiles.length > 0) {
      await processFile(droppedFiles[0]) // Pega apenas o primeiro arquivo
    }
  }

  // Função para remover arquivo
  const removeFile = () => {
    setUploadedFile(null)
    setDoctors([])
    setFilteredDoctors([])
    setSpecialties([])
    if (fileInputRef.current) {
      fileInputRef.current.value = ''
    }
  }

  // Aplicar filtros - apenas formatação de texto (lógica de filtro será feita pelo backend)
  useEffect(() => {
    // Se não há filtros ativos, mostra todos os médicos
    const hasActiveFilters = filters.cidade || filters.especialidade || filters.search
    
    if (!hasActiveFilters) {
      setFilteredDoctors(doctors)
    } else {
      // Quando há filtros, o backend será responsável pela lógica de filtro
      // Por enquanto, mostra todos os médicos (será substituído pela resposta do backend)
      setFilteredDoctors(doctors)
    }
    
    setCurrentPage(1) // Reset para primeira página quando filtros mudam
  }, [doctors, filters])

  // Funções de paginação
  const totalPages = Math.ceil(filteredDoctors.length / itemsPerPage)
  const startIndex = (currentPage - 1) * itemsPerPage
  const endIndex = startIndex + itemsPerPage
  const currentDoctors = filteredDoctors.slice(startIndex, endIndex)

  const goToPage = (page: number) => {
    setCurrentPage(Math.max(1, Math.min(page, totalPages)))
  }

  const goToFirstPage = () => goToPage(1)
  const goToLastPage = () => goToPage(totalPages)
  const goToPreviousPage = () => goToPage(currentPage - 1)
  const goToNextPage = () => goToPage(currentPage + 1)

  const clearFilters = () => {
    setFilters({
      cidade: '',
      especialidade: '',
      search: ''
    })
    // Garantir que todos os médicos sejam exibidos quando filtros são limpos
    setFilteredDoctors(doctors)
    setCurrentPage(1)
  }

  // Função para obter nome da cidade baseado no código
  const getCityName = (cityCode: number) => {
    const cityNames: { [key: number]: string } = {
      1: "São Paulo",
      2: "Rio de Janeiro", 
      3: "Belo Horizonte",
      4: "Salvador",
      5: "Brasília"
    }
    return cityNames[cityCode] || `Cidade ${cityCode}`
  }

  return (
    <div className="flex-1 space-y-4 p-4 md:p-8 pt-6">
      <div className="flex items-center justify-between space-y-2">
        <h2 className="text-3xl font-bold tracking-tight">Médicos</h2>
        <div className="flex items-center space-x-2">
          <Button variant="outline" size="sm" onClick={loadMockData}>
            <FileText className="mr-2 h-4 w-4" />
            Carregar Dados de Teste
          </Button>
          <Button variant="outline" size="sm" onClick={() => navigate('/file-manager')}>
            <FolderOpen className="mr-2 h-4 w-4" />
            Gerenciar Arquivos
          </Button>
          <Button variant="outline" size="sm" disabled={doctors.length === 0}>
            <Download className="mr-2 h-4 w-4" />
            Exportar Dados
          </Button>
        </div>
      </div>

      {/* Área de Upload - só aparece quando não há dados carregados */}
      {doctors.length === 0 && (
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Upload className="h-5 w-5" />
              Upload de Dados dos Médicos
            </CardTitle>
          </CardHeader>
          <CardContent>
            {!uploadedFile ? (
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
                  {isDragOver ? 'Solte o arquivo aqui' : 'Faça upload do arquivo de médicos'}
                </h3>
                <p className="text-muted-foreground mb-4">
                  {isDragOver 
                    ? 'Arraste e solte o arquivo aqui'
                    : 'Arraste um arquivo aqui ou selecione com os dados dos médicos (CSV, Excel, etc.)'
                  }
                </p>
                <div className="flex items-center justify-center gap-4">
                  <Button onClick={() => fileInputRef.current?.click()}>
                    <Upload className="mr-2 h-4 w-4" />
                    Selecionar Arquivo
                  </Button>
                </div>
                <input
                  ref={fileInputRef}
                  type="file"
                  accept="*/*"
                  onChange={handleFileUpload}
                  className="hidden"
                />
              </div>
            ) : (
              <div className="flex items-center justify-between p-4 border rounded-lg bg-muted/50">
                <div className="flex items-center space-x-3">
                  <FileText className="h-8 w-8 text-primary" />
                  <div>
                    <p className="font-medium">{uploadedFile.name}</p>
                    <p className="text-sm text-muted-foreground">
                      {(uploadedFile.size / 1024).toFixed(2)} KB
                    </p>
                  </div>
                </div>
                <div className="flex items-center space-x-2">
                  {uploading && (
                    <div className="flex items-center space-x-2 text-sm text-muted-foreground">
                      <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-primary"></div>
                      <span>Processando...</span>
                    </div>
                  )}
                  <Button variant="outline" size="sm" onClick={removeFile} disabled={uploading}>
                    <X className="h-4 w-4" />
                  </Button>
                </div>
              </div>
            )}
          </CardContent>
        </Card>
      )}

      {/* Botão para página de upload - só aparece quando há dados carregados */}
      {doctors.length > 0 && (
        <Card>
          <CardContent className="pt-6">
            <div className="text-center">
              <p className="text-muted-foreground mb-4">
                Dados carregados com sucesso! ({doctors.length} médicos)
              </p>
              <Button variant="outline" size="sm">
                <Upload className="mr-2 h-4 w-4" />
                Ir para Upload de Arquivos
              </Button>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Estatísticas - só aparecem quando há dados carregados */}
      {doctors.length > 0 && (
        <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Total de Médicos
            </CardTitle>
            <User className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{doctors.length}</div>
            <p className="text-xs text-muted-foreground">
              {filteredDoctors.length} filtrados
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Especialidades Únicas
            </CardTitle>
            <Stethoscope className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {specialties.length}
            </div>
            <p className="text-xs text-muted-foreground">
              Tipos de especialidades
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Cidades Atendidas
            </CardTitle>
            <MapPin className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {new Set(doctors.map(d => d.cidade)).size}
            </div>
            <p className="text-xs text-muted-foreground">
              Cidades diferentes
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Especialidade Mais Comum
            </CardTitle>
            <Stethoscope className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {specialties.length > 0 ? 
                specialties.reduce((a, b) => 
                  doctors.filter(d => d.especialidade === a).length > 
                  doctors.filter(d => d.especialidade === b).length ? a : b
                ) : 'N/A'
              }
            </div>
            <p className="text-xs text-muted-foreground">
              Mais frequente
            </p>
          </CardContent>
        </Card>
        </div>
      )}

      {/* Filtros - só aparece quando há dados */}
      {doctors.length > 0 && (
        <Card>
        <CardHeader>
                 <div className="flex items-center justify-between">
                   <CardTitle className="flex items-center gap-2">
                     <Filter className="h-5 w-5" />
                     Filtros
                   </CardTitle>
                   <div className="flex items-center space-x-2">
                     <Button variant="outline" size="sm" onClick={() => {
                       setFilteredDoctors(doctors)
                       setCurrentPage(1)
                     }}>
                       Ver Todos
                     </Button>
                     <Button variant="outline" size="sm" onClick={clearFilters}>
                       Limpar Filtros
                     </Button>
                   </div>
                 </div>
        </CardHeader>
        <CardContent>
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                   {/* Filtro por Cidade - apenas formatação de texto */}
                   <div className="space-y-2">
                     <label className="text-sm font-medium">Cidade</label>
                     <div className="relative">
                       <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
                       <Input
                         placeholder="Buscar cidade..."
                         value={filters.cidade}
                         onChange={(e) => setFilters(prev => ({ ...prev, cidade: e.target.value }))}
                         className="pl-8"
                       />
                     </div>
                     <p className="text-xs text-muted-foreground">Filtro será processado pelo backend</p>
                   </div>

                   {/* Filtro por Nome - apenas formatação de texto */}
                   <div className="space-y-2">
                     <label className="text-sm font-medium">Nome do Médico</label>
                     <div className="relative">
                       <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
                       <Input
                         placeholder="Buscar nome..."
                         value={filters.search}
                         onChange={(e) => setFilters(prev => ({ ...prev, search: e.target.value }))}
                         className="pl-8"
                       />
                     </div>
                     <p className="text-xs text-muted-foreground">Filtro será processado pelo backend</p>
                   </div>

                   {/* Filtro por Especialidade - apenas formatação de texto */}
                   <div className="space-y-2">
                     <label className="text-sm font-medium">Especialidade</label>
              <Popover open={openSpecialty} onOpenChange={setOpenSpecialty}>
                <PopoverTrigger asChild>
                  <Button
                    variant="outline"
                    role="combobox"
                    aria-expanded={openSpecialty}
                    className="w-full justify-between"
                  >
                    {filters.especialidade || "Selecionar especialidade..."}
                    <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                  </Button>
                </PopoverTrigger>
                <PopoverContent className="w-full p-0">
                  <Command>
                    <CommandInput placeholder="Buscar especialidade..." />
                    <CommandList>
                      <CommandEmpty>Nenhuma especialidade encontrada.</CommandEmpty>
                      <CommandGroup>
                        <CommandItem
                          value=""
                          onSelect={() => {
                            setFilters(prev => ({ ...prev, especialidade: '' }))
                            setOpenSpecialty(false)
                          }}
                        >
                          <Check
                            className={cn(
                              "mr-2 h-4 w-4",
                              filters.especialidade === '' ? "opacity-100" : "opacity-0"
                            )}
                          />
                          Todas as especialidades
                        </CommandItem>
                        {specialties.map((specialty) => (
                          <CommandItem
                            key={specialty}
                            value={specialty}
                            onSelect={() => {
                              setFilters(prev => ({ ...prev, especialidade: specialty }))
                              setOpenSpecialty(false)
                            }}
                          >
                            <Check
                              className={cn(
                                "mr-2 h-4 w-4",
                                filters.especialidade === specialty ? "opacity-100" : "opacity-0"
                              )}
                            />
                            {specialty}
                          </CommandItem>
                        ))}
                      </CommandGroup>
                    </CommandList>
                  </Command>
                </PopoverContent>
              </Popover>
              <p className="text-xs text-muted-foreground">Filtro será processado pelo backend</p>
            </div>

            {/* Contador de resultados */}
            <div className="space-y-2">
              <label className="text-sm font-medium">Resultados</label>
              <div className="flex items-center justify-center h-10 px-3 py-2 text-sm border rounded-md bg-muted">
                {filteredDoctors.length} médico{filteredDoctors.length !== 1 ? 's' : ''}
              </div>
            </div>
          </div>
        </CardContent>
        </Card>
      )}

      {/* Lista de Médicos */}
      <Card>
        <CardHeader>
          <CardTitle>Médicos Encontrados</CardTitle>
        </CardHeader>
        <CardContent>
          
          {loading ? (
            <div className="flex items-center justify-center h-32">
              <div className="text-center">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-2"></div>
                <p className="text-muted-foreground">Carregando médicos...</p>
              </div>
            </div>
          ) : doctors.length === 0 ? (
            <div className="text-center py-8">
              <User className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
              <h3 className="text-lg font-medium mb-2">Nenhum médico carregado</h3>
              <p className="text-muted-foreground">
                Faça upload de um arquivo com os dados dos médicos para começar.
              </p>
            </div>
          ) : filteredDoctors.length === 0 ? (
            <div className="text-center py-8">
              <User className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
              <h3 className="text-lg font-medium mb-2">Nenhum médico encontrado</h3>
              <p className="text-muted-foreground">
                Tente ajustar os filtros para encontrar médicos.
              </p>
            </div>
          ) : (
            <div className="space-y-4">
              {/* Grid de médicos paginados */}
              <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                {currentDoctors.map((doctor) => (
                  <Card key={doctor.codigo} className="hover:shadow-md transition-shadow">
                    <CardContent className="p-4">
                      <div className="flex items-start space-x-3">
                        <div className="flex-shrink-0">
                          <div className="w-12 h-12 bg-primary/10 rounded-lg flex items-center justify-center">
                            <User className="h-6 w-6 text-primary" />
                          </div>
                        </div>
                        <div className="flex-1 min-w-0">
                          <h3 className="text-lg font-semibold truncate">{doctor.nome_completo}</h3>
                          <div className="space-y-1 mt-2">
                            <div className="flex items-center space-x-1 text-sm text-muted-foreground">
                              <MapPin className="h-4 w-4" />
                              <span>{getCityName(doctor.cidade)}</span>
                            </div>
                            <div className="flex items-center space-x-1 text-sm text-muted-foreground">
                              <Stethoscope className="h-4 w-4" />
                              <span>{doctor.especialidade}</span>
                            </div>
                          </div>
                          <div className="mt-3 pt-3 border-t">
                            <div className="flex items-center justify-between">
                              <span className="text-xs text-muted-foreground">
                                Código: {doctor.codigo}
                              </span>
                              <Button variant="outline" size="sm">
                                Ver Detalhes
                              </Button>
                            </div>
                          </div>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>

              {/* Controles de paginação */}
              {totalPages > 1 && (
                <div className="flex items-center justify-between pt-4">
                  <div className="text-sm text-muted-foreground">
                    Mostrando {startIndex + 1} a {Math.min(endIndex, filteredDoctors.length)} de {filteredDoctors.length} médicos
                  </div>
                  <div className="flex items-center space-x-2">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={goToFirstPage}
                      disabled={currentPage === 1}
                    >
                      <ChevronsLeft className="h-4 w-4" />
                    </Button>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={goToPreviousPage}
                      disabled={currentPage === 1}
                    >
                      <ChevronLeft className="h-4 w-4" />
                    </Button>
                    
                    <div className="flex items-center space-x-1">
                      {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
                        let pageNum;
                        if (totalPages <= 5) {
                          pageNum = i + 1;
                        } else if (currentPage <= 3) {
                          pageNum = i + 1;
                        } else if (currentPage >= totalPages - 2) {
                          pageNum = totalPages - 4 + i;
                        } else {
                          pageNum = currentPage - 2 + i;
                        }
                        
                        return (
                          <Button
                            key={pageNum}
                            variant={currentPage === pageNum ? "default" : "outline"}
                            size="sm"
                            onClick={() => goToPage(pageNum)}
                            className="w-8 h-8 p-0"
                          >
                            {pageNum}
                          </Button>
                        );
                      })}
                    </div>

                    <Button
                      variant="outline"
                      size="sm"
                      onClick={goToNextPage}
                      disabled={currentPage === totalPages}
                    >
                      <ChevronRight className="h-4 w-4" />
                    </Button>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={goToLastPage}
                      disabled={currentPage === totalPages}
                    >
                      <ChevronsRight className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              )}
            </div>
          )}
        </CardContent>
        </Card>
    </div>
  )
}
