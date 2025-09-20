import { useState, useEffect, useRef } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Badge } from "@/components/ui/badge"
import { 
  Building2, 
  MapPin, 
  Users, 
  Stethoscope, 
  Search, 
  Filter,
  Download,
  Upload,
  FileText,
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

interface Hospital {
  codigo: string
  nome: string
  cidade: string
  jardim: string
  especialidades: string[]
  leitos_totais: number
}

interface FilterState {
  cidade: string
  jardim: string
  specialty: string
  search: string
}

export default function Hospitals() {
  const [hospitals, setHospitals] = useState<Hospital[]>([])
  const [filteredHospitals, setFilteredHospitals] = useState<Hospital[]>([])
  const [loading] = useState(false)
  const [uploading, setUploading] = useState(false)
  const [uploadedFile, setUploadedFile] = useState<File | null>(null)
  const [filters, setFilters] = useState<FilterState>({
    cidade: '',
    jardim: '',
    specialty: '',
    search: ''
  })
  const [openSpecialty, setOpenSpecialty] = useState(false)
  const [apiUrl, setApiUrl] = useState('')
  const fileInputRef = useRef<HTMLInputElement>(null)

  // Especialidades médicas (será preenchida dinamicamente com os dados)
  const [specialties, setSpecialties] = useState<string[]>([])

  // Função para lidar com upload de arquivo
  const handleFileUpload = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0]
    if (!file || !apiUrl) return

    setUploading(true)
    setUploadedFile(file)

    try {
      const formData = new FormData()
      formData.append('file', file)
      
      const response = await fetch(apiUrl, {
        method: 'POST',
        body: formData
      })
      
      if (!response.ok) {
        throw new Error(`Erro na API: ${response.status}`)
      }
      
      const data = await response.json()
      setHospitals(data)
      setFilteredHospitals(data)
      
      // Extrair especialidades únicas dos dados
      const uniqueSpecialties = [...new Set(data.flatMap((h: Hospital) => h.especialidades))] as string[]
      setSpecialties(uniqueSpecialties)
      
    } catch (error) {
      console.error('Erro ao fazer upload do arquivo:', error)
      alert('Erro ao processar arquivo. Verifique a URL da API.')
    } finally {
      setUploading(false)
    }
  }

  // Função para remover arquivo
  const removeFile = () => {
    setUploadedFile(null)
    setHospitals([])
    setFilteredHospitals([])
    setSpecialties([])
    if (fileInputRef.current) {
      fileInputRef.current.value = ''
    }
  }

  // Aplicar filtros
  useEffect(() => {
    let filtered = hospitals

    if (filters.cidade) {
      filtered = filtered.filter(h => 
        h.cidade.toLowerCase().includes(filters.cidade.toLowerCase())
      )
    }

    if (filters.jardim) {
      filtered = filtered.filter(h => 
        h.jardim.toLowerCase().includes(filters.jardim.toLowerCase())
      )
    }

    if (filters.specialty) {
      filtered = filtered.filter(h => 
        h.especialidades.includes(filters.specialty)
      )
    }

    if (filters.search) {
      filtered = filtered.filter(h => 
        h.nome.toLowerCase().includes(filters.search.toLowerCase()) ||
        h.cidade.toLowerCase().includes(filters.search.toLowerCase())
      )
    }

    setFilteredHospitals(filtered)
  }, [hospitals, filters])

  const clearFilters = () => {
    setFilters({
      cidade: '',
      jardim: '',
      specialty: '',
      search: ''
    })
  }

  return (
    <div className="flex-1 space-y-4 p-4 md:p-8 pt-6">
      <div className="flex items-center justify-between space-y-2">
        <h2 className="text-3xl font-bold tracking-tight">Hospitais</h2>
        <div className="flex items-center space-x-2">
          <Button variant="outline" size="sm" disabled={hospitals.length === 0}>
            <Download className="mr-2 h-4 w-4" />
            Exportar Dados
          </Button>
        </div>
      </div>

      {/* Configuração da API */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <FileText className="h-5 w-5" />
            Configuração da API
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-2">
            <label className="text-sm font-medium">URL da API</label>
            <Input
              placeholder="https://sua-api.com/upload"
              value={apiUrl}
              onChange={(e) => setApiUrl(e.target.value)}
            />
            <p className="text-xs text-muted-foreground">
              Insira a URL da sua API que receberá o arquivo via POST
            </p>
          </div>
        </CardContent>
      </Card>

      {/* Área de Upload */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Upload className="h-5 w-5" />
            Upload de Dados dos Hospitais
          </CardTitle>
        </CardHeader>
        <CardContent>
          {!uploadedFile ? (
            <div className="border-2 border-dashed border-muted-foreground/25 rounded-lg p-8 text-center hover:border-muted-foreground/50 transition-colors">
              <Upload className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
              <h3 className="text-lg font-medium mb-2">Faça upload do arquivo de hospitais</h3>
              <p className="text-muted-foreground mb-4">
                Selecione um arquivo com os dados dos hospitais (CSV, Excel, etc.)
              </p>
              <div className="flex items-center justify-center gap-4">
                <Button 
                  onClick={() => fileInputRef.current?.click()}
                  disabled={!apiUrl}
                >
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
              {!apiUrl && (
                <p className="text-sm text-destructive mt-2">
                  Configure a URL da API antes de fazer upload
                </p>
              )}
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

      {/* Estatísticas */}
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
              Cidades Cobertas
            </CardTitle>
            <MapPin className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {new Set(hospitals.map(h => h.cidade)).size}
            </div>
            <p className="text-xs text-muted-foreground">
              Cidades diferentes
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Filtros - só aparece quando há dados */}
      {hospitals.length > 0 && (
        <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle className="flex items-center gap-2">
              <Filter className="h-5 w-5" />
              Filtros
            </CardTitle>
            <Button variant="outline" size="sm" onClick={clearFilters}>
              Limpar Filtros
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
            {/* Filtro por Cidade */}
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
            </div>

            {/* Filtro por Jardim */}
            <div className="space-y-2">
              <label className="text-sm font-medium">Jardim</label>
              <div className="relative">
                <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
                <Input
                  placeholder="Buscar jardim..."
                  value={filters.jardim}
                  onChange={(e) => setFilters(prev => ({ ...prev, jardim: e.target.value }))}
                  className="pl-8"
                />
              </div>
            </div>

            {/* Filtro por Nome */}
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
            </div>

            {/* Filtro por Especialidade */}
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
              <p className="text-muted-foreground mb-4">
                Faça upload de um arquivo CSV com os dados dos hospitais para começar.
              </p>
              <Button onClick={() => fileInputRef.current?.click()}>
                <Upload className="mr-2 h-4 w-4" />
                Fazer Upload
              </Button>
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
               {filteredHospitals.map((hospital) => (
                 <div key={hospital.codigo} className="flex items-center justify-between p-4 border rounded-lg">
                  <div className="flex items-center space-x-4">
                    <div className="flex-shrink-0">
                      <Building2 className="h-8 w-8 text-primary" />
                    </div>
                     <div className="flex-1">
                       <h3 className="text-lg font-medium">{hospital.nome}</h3>
                       <div className="flex items-center space-x-4 mt-1">
                         <div className="flex items-center space-x-1 text-sm text-muted-foreground">
                           <MapPin className="h-4 w-4" />
                           <span>{hospital.cidade}</span>
                         </div>
                         <div className="flex items-center space-x-1 text-sm text-muted-foreground">
                           <Stethoscope className="h-4 w-4" />
                           <span>{hospital.leitos_totais} leitos</span>
                         </div>
                       </div>
                       <div className="flex flex-wrap gap-1 mt-2">
                         <Badge variant="outline" className="text-xs">
                           {hospital.jardim}
                         </Badge>
                         {hospital.especialidades.slice(0, 2).map((specialty) => (
                           <Badge key={specialty} variant="outline" className="text-xs">
                             {specialty}
                           </Badge>
                         ))}
                         {hospital.especialidades.length > 2 && (
                           <Badge variant="outline" className="text-xs">
                             +{hospital.especialidades.length - 2} mais
                           </Badge>
                         )}
                       </div>
                     </div>
                   </div>
                   <div className="flex items-center space-x-4">
                     <div className="text-right">
                       <p className="text-xs text-muted-foreground">
                         Código: {hospital.codigo}
                       </p>
                     </div>
                     <Button variant="outline" size="sm">
                       Ver Detalhes
                     </Button>
                   </div>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
