import { useState, useEffect, useRef } from "react"
import { useNavigate } from "react-router-dom"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { api } from "@/lib/api"
import { 
  ChevronLeft, 
  ChevronRight, 
  ChevronsLeft, 
  ChevronsRight 
} from "lucide-react"
import { 
  MapPin, 
  Users, 
  Building2, 
  Search, 
  Filter,
  Download,
  Upload,
  FileText,
  FolderOpen,
  X,
  Globe,
  TrendingUp,
  Activity
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

interface State {
  codigo: string
  nome: string
  sigla: string
  regiao: string
  populacao: number
  area_km2: number
  capital: string
  governador: string
  pib_per_capita: number
  municipios: number
}

interface FilterState {
  regiao: string
  search: string
  populacao_min: string
  populacao_max: string
}

export default function States() {
  const navigate = useNavigate()
  const [states, setStates] = useState<State[]>([])
  const [filteredStates, setFilteredStates] = useState<State[]>([])
  const [loading] = useState(false)
  const [uploading, setUploading] = useState(false)
  const [uploadedFile, setUploadedFile] = useState<File | null>(null)
  const [isDragOver, setIsDragOver] = useState(false)
  const [filters, setFilters] = useState<FilterState>({
    regiao: '',
    search: '',
    populacao_min: '',
    populacao_max: ''
  })
  const [openRegion, setOpenRegion] = useState(false)
  const fileInputRef = useRef<HTMLInputElement>(null)

  // Estados para paginação
  const [currentPage, setCurrentPage] = useState(1)
  const [itemsPerPage] = useState(6) // 6 estados por página


  // Regiões brasileiras (será preenchida dinamicamente com os dados)
  const [regions, setRegions] = useState<string[]>([])


  // Função para processar arquivo (usada tanto para input quanto drag&drop)
  const processFile = async (file: File) => {
    if (!file) return

    setUploading(true)
    setUploadedFile(file)

    try {
      const response = await api.states.uploadFile(file)
      
      // A resposta da API vem no formato { data: [...], success: true, message?: string }
      const data = response.data as State[]
      setStates(data)
      setFilteredStates(data)
      
      // Extrair regiões únicas dos dados
      const uniqueRegions = [...new Set(data.map((s: State) => s.regiao))] as string[]
      setRegions(uniqueRegions)
      
      console.log('Upload realizado com sucesso:', response.message)
      
    } catch (error) {
      console.error('Erro ao fazer upload do arquivo:', error)
      alert(`Erro ao processar arquivo: ${error instanceof Error ? error.message : 'Erro desconhecido'}`)
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
    setStates([])
    setFilteredStates([])
    setRegions([])
    if (fileInputRef.current) {
      fileInputRef.current.value = ''
    }
  }

  // Aplicar filtros - apenas formatação de texto (lógica de filtro será feita pelo backend)
  useEffect(() => {
    // Se não há filtros ativos, mostra todos os estados
    const hasActiveFilters = filters.regiao || filters.search || filters.populacao_min || filters.populacao_max
    
    if (!hasActiveFilters) {
      setFilteredStates(states)
    } else {
      // Quando há filtros, o backend será responsável pela lógica de filtro
      // Por enquanto, mostra todos os estados (será substituído pela resposta do backend)
      setFilteredStates(states)
    }
    
    setCurrentPage(1) // Reset para primeira página quando filtros mudam
  }, [states, filters])

  // Funções de paginação
  const totalPages = Math.ceil(filteredStates.length / itemsPerPage)
  const startIndex = (currentPage - 1) * itemsPerPage
  const endIndex = startIndex + itemsPerPage
  const currentStates = filteredStates.slice(startIndex, endIndex)

  const goToPage = (page: number) => {
    setCurrentPage(Math.max(1, Math.min(page, totalPages)))
  }

  const goToFirstPage = () => goToPage(1)
  const goToLastPage = () => goToPage(totalPages)
  const goToPreviousPage = () => goToPage(currentPage - 1)
  const goToNextPage = () => goToPage(currentPage + 1)

  const clearFilters = () => {
    setFilters({
      regiao: '',
      search: '',
      populacao_min: '',
      populacao_max: ''
    })
    // Garantir que todos os estados sejam exibidos quando filtros são limpos
    setFilteredStates(states)
    setCurrentPage(1)
  }

  // Função para formatar números
  const formatNumber = (num: number) => {
    return num.toLocaleString('pt-BR')
  }

  // Função para formatar área
  const formatArea = (area: number) => {
    return `${formatNumber(area)} km²`
  }

  // Função para formatar PIB per capita
  const formatPIB = (pib: number) => {
    return `R$ ${formatNumber(pib)}`
  }

  return (
    <div className="flex-1 space-y-4 p-4 md:p-8 pt-6">
      <div className="flex items-center justify-between space-y-2">
        <h2 className="text-3xl font-bold tracking-tight">Estados</h2>
        <div className="flex items-center space-x-2">
          <Button variant="outline" size="sm" onClick={() => navigate('/file-manager')}>
            <FolderOpen className="mr-2 h-4 w-4" />
            Gerenciar Arquivos
          </Button>
          <Button variant="outline" size="sm" disabled={states.length === 0}>
            <Download className="mr-2 h-4 w-4" />
            Exportar Dados
          </Button>
        </div>
      </div>

      {/* Área de Upload - só aparece quando não há dados carregados */}
      {states.length === 0 && (
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Upload className="h-5 w-5" />
              Upload de Dados dos Estados
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
                  {isDragOver ? 'Solte o arquivo aqui' : 'Faça upload do arquivo de estados'}
                </h3>
                <p className="text-muted-foreground mb-4">
                  {isDragOver 
                    ? 'Arraste e solte o arquivo aqui'
                    : 'Arraste um arquivo aqui ou selecione com os dados dos estados (CSV, Excel, etc.)'
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

      {/* Informação de sucesso - só aparece quando há dados carregados */}
      {states.length > 0 && (
        <Card>
          <CardContent className="pt-6">
            <div className="text-center">
              <p className="text-muted-foreground mb-4">
                Dados carregados com sucesso! ({states.length} estados)
              </p>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Estatísticas - só aparecem quando há dados carregados */}
      {states.length > 0 && (
        <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Total de Estados
            </CardTitle>
            <MapPin className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{states.length}</div>
            <p className="text-xs text-muted-foreground">
              {filteredStates.length} filtrados
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Regiões Únicas
            </CardTitle>
            <Globe className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {regions.length}
            </div>
            <p className="text-xs text-muted-foreground">
              Regiões do Brasil
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              População Total
            </CardTitle>
            <Users className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {formatNumber(states.reduce((sum, s) => sum + s.populacao, 0))}
            </div>
            <p className="text-xs text-muted-foreground">
              Habitantes
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Área Total
            </CardTitle>
            <Building2 className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {formatNumber(Math.round(states.reduce((sum, s) => sum + s.area_km2, 0)))}
            </div>
            <p className="text-xs text-muted-foreground">
              km²
            </p>
          </CardContent>
        </Card>
        </div>
      )}

      {/* Filtros - só aparece quando há dados */}
      {states.length > 0 && (
        <Card>
        <CardHeader>
                 <div className="flex items-center justify-between">
                   <CardTitle className="flex items-center gap-2">
                     <Filter className="h-5 w-5" />
                     Filtros
                   </CardTitle>
                   <div className="flex items-center space-x-2">
                     <Button variant="outline" size="sm" onClick={() => {
                       setFilteredStates(states)
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
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
                   {/* Filtro por Região */}
                   <div className="space-y-2">
                     <label className="text-sm font-medium">Região</label>
              <Popover open={openRegion} onOpenChange={setOpenRegion}>
                <PopoverTrigger asChild>
                  <Button
                    variant="outline"
                    role="combobox"
                    aria-expanded={openRegion}
                    className="w-full justify-between"
                  >
                    {filters.regiao || "Selecionar região..."}
                    <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                  </Button>
                </PopoverTrigger>
                <PopoverContent className="w-full p-0">
                  <Command>
                    <CommandInput placeholder="Buscar região..." />
                    <CommandList>
                      <CommandEmpty>Nenhuma região encontrada.</CommandEmpty>
                      <CommandGroup>
                        <CommandItem
                          value=""
                          onSelect={() => {
                            setFilters(prev => ({ ...prev, regiao: '' }))
                            setOpenRegion(false)
                          }}
                        >
                          <Check
                            className={cn(
                              "mr-2 h-4 w-4",
                              filters.regiao === '' ? "opacity-100" : "opacity-0"
                            )}
                          />
                          Todas as regiões
                        </CommandItem>
                        {regions.map((region) => (
                          <CommandItem
                            key={region}
                            value={region}
                            onSelect={() => {
                              setFilters(prev => ({ ...prev, regiao: region }))
                              setOpenRegion(false)
                            }}
                          >
                            <Check
                              className={cn(
                                "mr-2 h-4 w-4",
                                filters.regiao === region ? "opacity-100" : "opacity-0"
                              )}
                            />
                            {region}
                          </CommandItem>
                        ))}
                      </CommandGroup>
                    </CommandList>
                  </Command>
                </PopoverContent>
              </Popover>
              <p className="text-xs text-muted-foreground">Filtro será processado pelo backend</p>
            </div>

                   {/* Filtro por Nome */}
                   <div className="space-y-2">
                     <label className="text-sm font-medium">Nome do Estado</label>
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

                   {/* Filtro por População Mínima */}
                   <div className="space-y-2">
                     <label className="text-sm font-medium">População Mínima</label>
                     <Input
                       placeholder="Ex: 1000000"
                       value={filters.populacao_min}
                       onChange={(e) => setFilters(prev => ({ ...prev, populacao_min: e.target.value }))}
                       type="number"
                     />
                     <p className="text-xs text-muted-foreground">Filtro será processado pelo backend</p>
                   </div>

                   {/* Filtro por População Máxima */}
                   <div className="space-y-2">
                     <label className="text-sm font-medium">População Máxima</label>
                     <Input
                       placeholder="Ex: 50000000"
                       value={filters.populacao_max}
                       onChange={(e) => setFilters(prev => ({ ...prev, populacao_max: e.target.value }))}
                       type="number"
                     />
                     <p className="text-xs text-muted-foreground">Filtro será processado pelo backend</p>
                   </div>

            {/* Contador de resultados */}
            <div className="space-y-2">
              <label className="text-sm font-medium">Resultados</label>
              <div className="flex items-center justify-center h-10 px-3 py-2 text-sm border rounded-md bg-muted">
                {filteredStates.length} estado{filteredStates.length !== 1 ? 's' : ''}
              </div>
            </div>
          </div>
        </CardContent>
        </Card>
      )}

      {/* Lista de Estados */}
      <Card>
        <CardHeader>
          <CardTitle>Estados Encontrados</CardTitle>
        </CardHeader>
        <CardContent>
          
          {loading ? (
            <div className="flex items-center justify-center h-32">
              <div className="text-center">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-2"></div>
                <p className="text-muted-foreground">Carregando estados...</p>
              </div>
            </div>
          ) : states.length === 0 ? (
            <div className="text-center py-8">
              <MapPin className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
              <h3 className="text-lg font-medium mb-2">Nenhum estado carregado</h3>
              <p className="text-muted-foreground">
                Faça upload de um arquivo com os dados dos estados para começar.
              </p>
            </div>
          ) : filteredStates.length === 0 ? (
            <div className="text-center py-8">
              <MapPin className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
              <h3 className="text-lg font-medium mb-2">Nenhum estado encontrado</h3>
              <p className="text-muted-foreground">
                Tente ajustar os filtros para encontrar estados.
              </p>
            </div>
          ) : (
            <div className="space-y-4">
              {/* Grid de estados paginados */}
              <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                {currentStates.map((state) => (
                  <Card key={state.codigo} className="hover:shadow-md transition-shadow">
                    <CardContent className="p-4">
                      <div className="flex items-start space-x-3">
                        <div className="flex-shrink-0">
                          <div className="w-12 h-12 bg-primary/10 rounded-lg flex items-center justify-center">
                            <MapPin className="h-6 w-6 text-primary" />
                          </div>
                        </div>
                        <div className="flex-1 min-w-0">
                          <h3 className="text-lg font-semibold truncate">{state.nome}</h3>
                          <div className="space-y-1 mt-2">
                            <div className="flex items-center space-x-1 text-sm text-muted-foreground">
                              <Globe className="h-4 w-4" />
                              <span>{state.regiao}</span>
                            </div>
                            <div className="flex items-center space-x-1 text-sm text-muted-foreground">
                              <Users className="h-4 w-4" />
                              <span>{formatNumber(state.populacao)} hab.</span>
                            </div>
                            <div className="flex items-center space-x-1 text-sm text-muted-foreground">
                              <Building2 className="h-4 w-4" />
                              <span>{formatArea(state.area_km2)}</span>
                            </div>
                            <div className="flex items-center space-x-1 text-sm text-muted-foreground">
                              <span className="font-medium">Capital:</span>
                              <span>{state.capital}</span>
                            </div>
                            <div className="flex items-center space-x-1 text-sm text-muted-foreground">
                              <TrendingUp className="h-4 w-4" />
                              <span>PIB: {formatPIB(state.pib_per_capita)}</span>
                            </div>
                            <div className="flex items-center space-x-1 text-sm text-muted-foreground">
                              <Activity className="h-4 w-4" />
                              <span>{state.municipios} municípios</span>
                            </div>
                          </div>
                          <div className="mt-3 pt-3 border-t">
                            <div className="flex items-center justify-between">
                              <span className="text-xs text-muted-foreground">
                                {state.sigla} • Governador: {state.governador}
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
                    Mostrando {startIndex + 1} a {Math.min(endIndex, filteredStates.length)} de {filteredStates.length} estados
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
