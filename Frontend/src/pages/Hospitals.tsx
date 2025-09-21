import { useState, useEffect } from "react"
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
  Building2, 
  MapPin, 
  Users, 
  Stethoscope, 
  Search, 
  Filter,
  Download,
  FileText,
  FolderOpen
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

interface Hospital {
  codigo: string
  nome: string
  cidade: string
  bairro: string
  especialidades: string[]
  leitos_totais: number
  colaboradores: number
}

interface FilterState {
  cidade: string
  bairro: string
  specialty: string
  search: string
}

export default function Hospitals() {
  const navigate = useNavigate()
  const [hospitals, setHospitals] = useState<Hospital[]>([])
  const [filteredHospitals, setFilteredHospitals] = useState<Hospital[]>([])
  const [loading] = useState(false)
  const [filters, setFilters] = useState<FilterState>({
    cidade: '',
    bairro: '',
    specialty: '',
    search: ''
  })
  const [openSpecialty, setOpenSpecialty] = useState(false)

  // Estados para paginação
  const [currentPage, setCurrentPage] = useState(1)
  const [itemsPerPage] = useState(6) // 6 hospitais por página


  // Dados mockados para teste
  const mockHospitals: Hospital[] = [
    {
      codigo: "H001",
      nome: "Hospital São Paulo",
      cidade: "São Paulo",
      bairro: "Centro",
      especialidades: ["Cardiologia", "Neurologia", "Pediatria"],
      leitos_totais: 450,
      colaboradores: 320
    },
    {
      codigo: "H002", 
      nome: "Hospital das Clínicas",
      cidade: "São Paulo",
      bairro: "Cerqueira César",
      especialidades: ["Cardiologia", "Oncologia", "Transplantes"],
      leitos_totais: 1200,
      colaboradores: 850
    },
    {
      codigo: "H003",
      nome: "Hospital Sírio-Libanês",
      cidade: "São Paulo", 
      bairro: "Bela Vista",
      especialidades: ["Cardiologia", "Neurologia", "Ortopedia"],
      leitos_totais: 350,
      colaboradores: 280
    },
    {
      codigo: "H004",
      nome: "Hospital Copa D'Or",
      cidade: "Rio de Janeiro",
      bairro: "Copacabana",
      especialidades: ["Cardiologia", "Neurologia", "Pediatria"],
      leitos_totais: 280,
      colaboradores: 220
    },
    {
      codigo: "H005",
      nome: "Hospital Albert Einstein",
      cidade: "São Paulo",
      bairro: "Morumbi",
      especialidades: ["Cardiologia", "Oncologia", "Neurologia", "Pediatria"],
      leitos_totais: 600,
      colaboradores: 450
    },
    {
      codigo: "H006",
      nome: "Hospital Oswaldo Cruz",
      cidade: "São Paulo",
      bairro: "Paraíso",
      especialidades: ["Cardiologia", "Neurologia"],
      leitos_totais: 180,
      colaboradores: 150
    },
    {
      codigo: "H007",
      nome: "Hospital Samaritano",
      cidade: "Rio de Janeiro",
      bairro: "Botafogo",
      especialidades: ["Cardiologia", "Ortopedia", "Pediatria"],
      leitos_totais: 220,
      colaboradores: 180
    },
    {
      codigo: "H008",
      nome: "Hospital Beneficência Portuguesa",
      cidade: "São Paulo",
      bairro: "Vila Clementino",
      especialidades: ["Cardiologia", "Neurologia", "Oncologia"],
      leitos_totais: 400,
      colaboradores: 320
    },
    {
      codigo: "H009",
      nome: "Hospital Pró-Cardíaco",
      cidade: "Rio de Janeiro",
      bairro: "Botafogo",
      especialidades: ["Cardiologia"],
      leitos_totais: 150,
      colaboradores: 120
    },
    {
      codigo: "H010",
      nome: "Hospital Santa Catarina",
      cidade: "São Paulo",
      bairro: "Vila Mariana",
      especialidades: ["Cardiologia", "Neurologia", "Pediatria", "Ortopedia"],
      leitos_totais: 320,
      colaboradores: 250
    },
    {
      codigo: "H011",
      nome: "Hospital São Luiz",
      cidade: "São Paulo",
      bairro: "Itaim Bibi",
      especialidades: ["Cardiologia", "Neurologia"],
      leitos_totais: 200,
      colaboradores: 160
    },
    {
      codigo: "H012",
      nome: "Hospital Barra D'Or",
      cidade: "Rio de Janeiro",
      bairro: "Barra da Tijuca",
      especialidades: ["Cardiologia", "Neurologia", "Pediatria"],
      leitos_totais: 180,
      colaboradores: 140
    }
  ]

  // Especialidades médicas (será preenchida dinamicamente com os dados)
  const [specialties, setSpecialties] = useState<string[]>([])

  // Função para carregar dados mockados (para teste)
  const loadMockData = () => {
    setHospitals(mockHospitals)
    setFilteredHospitals(mockHospitals)
    
    // Extrair especialidades únicas dos dados mockados
    const uniqueSpecialties = [...new Set(mockHospitals.flatMap(h => h.especialidades))]
    setSpecialties(uniqueSpecialties)
  }


  // Aplicar filtros - apenas formatação de texto (lógica de filtro será feita pelo backend)
  useEffect(() => {
    // Se não há filtros ativos, mostra todos os hospitais
    const hasActiveFilters = filters.cidade || filters.bairro || filters.specialty || filters.search
    
    if (!hasActiveFilters) {
      setFilteredHospitals(hospitals)
    } else {
      // Quando há filtros, o backend será responsável pela lógica de filtro
      // Por enquanto, mostra todos os hospitais (será substituído pela resposta do backend)
      setFilteredHospitals(hospitals)
    }
    
    setCurrentPage(1) // Reset para primeira página quando filtros mudam
  }, [hospitals, filters])

  // Funções de paginação
  const totalPages = Math.ceil(filteredHospitals.length / itemsPerPage)
  const startIndex = (currentPage - 1) * itemsPerPage
  const endIndex = startIndex + itemsPerPage
  const currentHospitals = filteredHospitals.slice(startIndex, endIndex)

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
      bairro: '',
      specialty: '',
      search: ''
    })
    // Garantir que todos os hospitais sejam exibidos quando filtros são limpos
    setFilteredHospitals(hospitals)
    setCurrentPage(1)
  }


  return (
    <div className="flex-1 space-y-4 p-4 md:p-8 pt-6">
      <div className="flex items-center justify-between space-y-2">
        <h2 className="text-3xl font-bold tracking-tight">Hospitais</h2>
        <div className="flex items-center space-x-2">
          <Button variant="outline" size="sm" onClick={loadMockData}>
            <FileText className="mr-2 h-4 w-4" />
            Carregar Dados de Teste
          </Button>
          <Button variant="outline" size="sm" onClick={() => navigate('/file-manager')}>
            <FolderOpen className="mr-2 h-4 w-4" />
            Gerenciar Arquivos
          </Button>
          <Button variant="outline" size="sm" disabled={hospitals.length === 0}>
            <Download className="mr-2 h-4 w-4" />
            Exportar Dados
          </Button>
        </div>
      </div>



      {/* Botão para página de upload - só aparece quando há dados carregados */}
      {hospitals.length > 0 && (
        <Card>
          <CardContent className="pt-6">
            <div className="text-center">
              <p className="text-muted-foreground mb-4">
                Dados carregados com sucesso! ({hospitals.length} hospitais)
              </p>
              <Button variant="outline" size="sm" onClick={() => navigate('/file-manager')}>
                <FolderOpen className="mr-2 h-4 w-4" />
                Ir para Gerenciador de Arquivos
              </Button>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Estatísticas - só aparecem quando há dados carregados */}
      {hospitals.length > 0 && (
        <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Total de Hospitais
            </CardTitle>
            <Building2 className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{hospitals.length}</div>
            <p className="text-xs text-muted-foreground">
              {filteredHospitals.length} filtrados
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Especialidades Únicas
            </CardTitle>
            <Users className="h-4 w-4 text-muted-foreground" />
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
              Total de Leitos
            </CardTitle>
            <Stethoscope className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {hospitals.reduce((sum, h) => sum + h.leitos_totais, 0).toLocaleString()}
            </div>
            <p className="text-xs text-muted-foreground">
              Capacidade total
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Total de Colaboradores
            </CardTitle>
            <Users className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {hospitals.reduce((sum, h) => sum + h.colaboradores, 0).toLocaleString()}
            </div>
            <p className="text-xs text-muted-foreground">
              Médicos e funcionários
            </p>
          </CardContent>
        </Card>
        </div>
      )}

      {/* Filtros - só aparece quando há dados */}
      {hospitals.length > 0 && (
        <Card>
        <CardHeader>
                 <div className="flex items-center justify-between">
                   <CardTitle className="flex items-center gap-2">
                     <Filter className="h-5 w-5" />
                     Filtros
                   </CardTitle>
                   <div className="flex items-center space-x-2">
                     <Button variant="outline" size="sm" onClick={() => {
                       setFilteredHospitals(hospitals)
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

                   <div className="space-y-2">
                     <label className="text-sm font-medium">Bairro</label>
                     <div className="relative">
                       <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
                       <Input
                         placeholder="Buscar bairro..."
                         value={filters.bairro}
                         onChange={(e) => setFilters(prev => ({ ...prev, bairro: e.target.value }))}
                         className="pl-8"
                       />
                     </div>
                     <p className="text-xs text-muted-foreground">Filtro será processado pelo backend</p>
                   </div>

                   {/* Filtro por Nome - apenas formatação de texto */}
                   <div className="space-y-2">
                     <label className="text-sm font-medium">Nome do Hospital</label>
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
                    {filters.specialty || "Selecionar especialidade..."}
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
                            setFilters(prev => ({ ...prev, specialty: '' }))
                            setOpenSpecialty(false)
                          }}
                        >
                          <Check
                            className={cn(
                              "mr-2 h-4 w-4",
                              filters.specialty === '' ? "opacity-100" : "opacity-0"
                            )}
                          />
                          Todas as especialidades
                        </CommandItem>
                        {specialties.map((specialty) => (
                          <CommandItem
                            key={specialty}
                            value={specialty}
                            onSelect={() => {
                              setFilters(prev => ({ ...prev, specialty }))
                              setOpenSpecialty(false)
                            }}
                          >
                            <Check
                              className={cn(
                                "mr-2 h-4 w-4",
                                filters.specialty === specialty ? "opacity-100" : "opacity-0"
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
                {filteredHospitals.length} hospital{filteredHospitals.length !== 1 ? 'is' : ''}
              </div>
            </div>
          </div>
        </CardContent>
        </Card>
      )}

      {/* Lista de Hospitais */}
      <Card>
        <CardHeader>
          <CardTitle>Hospitais Encontrados</CardTitle>
        </CardHeader>
        <CardContent>
          
          {loading ? (
            <div className="flex items-center justify-center h-32">
              <div className="text-center">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-2"></div>
                <p className="text-muted-foreground">Carregando hospitais...</p>
              </div>
            </div>
          ) : hospitals.length === 0 ? (
            <div className="text-center py-8">
              <Building2 className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
              <h3 className="text-lg font-medium mb-2">Nenhum hospital carregado</h3>
              <p className="text-muted-foreground">
                Nenhum hospital encontrado. Use o Gerenciador de Arquivos para fazer upload dos dados.
              </p>
            </div>
          ) : filteredHospitals.length === 0 ? (
            <div className="text-center py-8">
              <Building2 className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
              <h3 className="text-lg font-medium mb-2">Nenhum hospital encontrado</h3>
              <p className="text-muted-foreground">
                Tente ajustar os filtros para encontrar hospitais.
              </p>
            </div>
          ) : (
            <div className="space-y-4">
              {/* Grid de hospitais paginados */}
              <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                {currentHospitals.map((hospital) => (
                  <Card key={hospital.codigo} className="hover:shadow-md transition-shadow">
                    <CardContent className="p-4">
                      <div className="flex items-start space-x-3">
                        <div className="flex-shrink-0">
                          <div className="w-12 h-12 bg-primary/10 rounded-lg flex items-center justify-center">
                            <Building2 className="h-6 w-6 text-primary" />
                          </div>
                        </div>
                        <div className="flex-1 min-w-0">
                          <h3 className="text-lg font-semibold truncate">{hospital.nome}</h3>
                          <div className="space-y-1 mt-2">
                            <div className="flex items-center space-x-1 text-sm text-muted-foreground">
                              <MapPin className="h-4 w-4" />
                              <span>{hospital.cidade}</span>
                            </div>
                            <div className="flex items-center space-x-1 text-sm text-muted-foreground">
                              <Stethoscope className="h-4 w-4" />
                              <span>{hospital.leitos_totais} leitos</span>
                            </div>
                            <div className="flex items-center space-x-1 text-sm text-muted-foreground">
                              <Users className="h-4 w-4" />
                              <span>{hospital.colaboradores} colaboradores</span>
                            </div>
                            <div className="flex items-center space-x-1 text-sm text-muted-foreground">
                              <span className="font-medium">Bairro:</span>
                              <span>{hospital.bairro}</span>
                            </div>
                          </div>
                          <div className="flex flex-wrap gap-1 mt-3">
                            {hospital.especialidades.map((specialty) => (
                              <Badge key={specialty} variant="secondary" className="text-xs">
                                {specialty}
                              </Badge>
                            ))}
                          </div>
                          <div className="mt-3 pt-3 border-t">
                            <div className="flex items-center justify-between">
                              <span className="text-xs text-muted-foreground">
                                Código: {hospital.codigo}
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
                    Mostrando {startIndex + 1} a {Math.min(endIndex, filteredHospitals.length)} de {filteredHospitals.length} hospitais
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
