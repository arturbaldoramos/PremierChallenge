import React, { useState } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { 
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Filter, Search, RefreshCw } from 'lucide-react'
import { 
  useHospitalsBySpecialty, 
  useHospitalsByMunicipality, 
  useSpecialties 
} from '@/hooks/useStats'

export function StatsFilters() {
  const [selectedSpecialty, setSelectedSpecialty] = useState<string>("")
  const [selectedMunicipality, setSelectedMunicipality] = useState<string>("")
  const [customSpecialty, setCustomSpecialty] = useState<string>("")
  const [customMunicipality, setCustomMunicipality] = useState<string>("")

  // Hooks para dados
  const { 
    data: hospitalsBySpecialty, 
    isLoading: hospitalsBySpecialtyLoading, 
    error: hospitalsBySpecialtyError, 
    refetch: refetchHospitalsBySpecialty 
  } = useHospitalsBySpecialty()

  const { 
    data: hospitalsByMunicipality, 
    isLoading: hospitalsByMunicipalityLoading, 
    error: hospitalsByMunicipalityError, 
    refetch: refetchHospitalsByMunicipality 
  } = useHospitalsByMunicipality()

  const { 
    data: specialties, 
    isLoading: specialtiesLoading, 
    error: specialtiesError 
  } = useSpecialties()

  // Funções para aplicar filtros
  const handleSpecialtyFilter = (specialty: string) => {
    setSelectedSpecialty(specialty)
    refetchHospitalsBySpecialty(specialty || undefined)
  }

  const handleMunicipalityFilter = (municipality: string) => {
    setSelectedMunicipality(municipality)
    refetchHospitalsByMunicipality(municipality || undefined)
  }

  const handleCustomSpecialtySearch = () => {
    if (customSpecialty.trim()) {
      handleSpecialtyFilter(customSpecialty.trim())
    }
  }

  const handleCustomMunicipalitySearch = () => {
    if (customMunicipality.trim()) {
      handleMunicipalityFilter(customMunicipality.trim())
    }
  }

  const clearFilters = () => {
    setSelectedSpecialty("")
    setSelectedMunicipality("")
    setCustomSpecialty("")
    setCustomMunicipality("")
    refetchHospitalsBySpecialty()
    refetchHospitalsByMunicipality()
  }

  return (
    <div className="space-y-6">
      {/* Filtros */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Filter className="h-5 w-5" />
            Filtros de Análise
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          {/* Filtro por Especialidade */}
          <div className="space-y-2">
            <label className="text-sm font-medium">Filtrar por Especialidade</label>
            <div className="flex gap-2">
              <Select value={selectedSpecialty} onValueChange={handleSpecialtyFilter}>
                <SelectTrigger className="flex-1">
                  <SelectValue placeholder="Selecione uma especialidade" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="">Todas as Especialidades</SelectItem>
                  {specialties?.map((specialty) => (
                    <SelectItem key={specialty} value={specialty}>
                      {specialty}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              <Button 
                onClick={() => refetchHospitalsBySpecialty()} 
                variant="outline" 
                size="sm"
              >
                <RefreshCw className="h-4 w-4" />
              </Button>
            </div>
          </div>

          {/* Busca Personalizada por Especialidade */}
          <div className="space-y-2">
            <label className="text-sm font-medium">Buscar Especialidade Personalizada</label>
            <div className="flex gap-2">
              <Input
                placeholder="Digite uma especialidade específica"
                value={customSpecialty}
                onChange={(e) => setCustomSpecialty(e.target.value)}
                onKeyPress={(e) => e.key === 'Enter' && handleCustomSpecialtySearch()}
              />
              <Button onClick={handleCustomSpecialtySearch} variant="outline">
                <Search className="h-4 w-4" />
              </Button>
            </div>
          </div>

          {/* Filtro por Município */}
          <div className="space-y-2">
            <label className="text-sm font-medium">Filtrar por Município</label>
            <div className="flex gap-2">
              <Select value={selectedMunicipality} onValueChange={handleMunicipalityFilter}>
                <SelectTrigger className="flex-1">
                  <SelectValue placeholder="Selecione um município" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="">Todos os Municípios</SelectItem>
                  {hospitalsByMunicipality?.slice(0, 20).map((item, index) => {
                    const municipio = 'municipio' in item ? item.municipio : item.nome
                    const count = 'count' in item ? item.count : 1
                    return (
                      <SelectItem key={index} value={municipio}>
                        {municipio} ({count} hospitais)
                      </SelectItem>
                    )
                  })}
                </SelectContent>
              </Select>
              <Button 
                onClick={() => refetchHospitalsByMunicipality()} 
                variant="outline" 
                size="sm"
              >
                <RefreshCw className="h-4 w-4" />
              </Button>
            </div>
          </div>

          {/* Busca Personalizada por Município */}
          <div className="space-y-2">
            <label className="text-sm font-medium">Buscar Município Personalizado</label>
            <div className="flex gap-2">
              <Input
                placeholder="Digite um município específico"
                value={customMunicipality}
                onChange={(e) => setCustomMunicipality(e.target.value)}
                onKeyPress={(e) => e.key === 'Enter' && handleCustomMunicipalitySearch()}
              />
              <Button onClick={handleCustomMunicipalitySearch} variant="outline">
                <Search className="h-4 w-4" />
              </Button>
            </div>
          </div>

          {/* Botão para limpar filtros */}
          <div className="flex justify-end">
            <Button onClick={clearFilters} variant="outline">
              Limpar Filtros
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Resultados dos Filtros */}
      <div className="grid gap-4 md:grid-cols-2">
        {/* Resultados por Especialidade */}
        <Card>
          <CardHeader>
            <CardTitle>
              Hospitais por Especialidade
              {selectedSpecialty && (
                <Badge variant="secondary" className="ml-2">
                  {selectedSpecialty}
                </Badge>
              )}
            </CardTitle>
          </CardHeader>
          <CardContent>
            {hospitalsBySpecialtyLoading ? (
              <p>Carregando dados...</p>
            ) : hospitalsBySpecialtyError ? (
              <p className="text-red-500">Erro: {hospitalsBySpecialtyError}</p>
            ) : (
              <div className="space-y-2">
                {Array.isArray(hospitalsBySpecialty) && hospitalsBySpecialty.length > 0 ? (
                  hospitalsBySpecialty.slice(0, 10).map((item, index) => (
                    <div key={index} className="flex items-center justify-between p-2 border rounded">
                      <div>
                        {'especialidade' in item ? (
                          <>
                            <p className="font-medium">{item.especialidade}</p>
                            <p className="text-sm text-muted-foreground">
                              {item.count} hospitais
                            </p>
                          </>
                        ) : (
                          <>
                            <p className="font-medium">{item.nome}</p>
                            <p className="text-sm text-muted-foreground">
                              {item.especialidades} • {item.bairro}
                            </p>
                          </>
                        )}
                      </div>
                      <Badge variant="outline">
                        {'count' in item ? item.count : 'Hospital'}
                      </Badge>
                    </div>
                  ))
                ) : (
                  <p className="text-muted-foreground">Nenhum resultado encontrado</p>
                )}
              </div>
            )}
          </CardContent>
        </Card>

        {/* Resultados por Município */}
        <Card>
          <CardHeader>
            <CardTitle>
              Hospitais por Município
              {selectedMunicipality && (
                <Badge variant="secondary" className="ml-2">
                  {selectedMunicipality}
                </Badge>
              )}
            </CardTitle>
          </CardHeader>
          <CardContent>
            {hospitalsByMunicipalityLoading ? (
              <p>Carregando dados...</p>
            ) : hospitalsByMunicipalityError ? (
              <p className="text-red-500">Erro: {hospitalsByMunicipalityError}</p>
            ) : (
              <div className="space-y-2">
                {Array.isArray(hospitalsByMunicipality) && hospitalsByMunicipality.length > 0 ? (
                  hospitalsByMunicipality.slice(0, 10).map((item, index) => (
                    <div key={index} className="flex items-center justify-between p-2 border rounded">
                      <div>
                        {'municipio' in item ? (
                          <>
                            <p className="font-medium">{item.municipio}</p>
                            <p className="text-sm text-muted-foreground">
                              {item.count} hospitais
                            </p>
                          </>
                        ) : (
                          <>
                            <p className="font-medium">{item.nome}</p>
                            <p className="text-sm text-muted-foreground">
                              {item.bairro} • {item.cod_municipio}
                            </p>
                          </>
                        )}
                      </div>
                      <Badge variant="outline">
                        {'count' in item ? item.count : 'Hospital'}
                      </Badge>
                    </div>
                  ))
                ) : (
                  <p className="text-muted-foreground">Nenhum resultado encontrado</p>
                )}
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Resumo dos Filtros Ativos */}
      <Card>
        <CardHeader>
          <CardTitle>Filtros Ativos</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex flex-wrap gap-2">
            {selectedSpecialty && (
              <Badge variant="default">
                Especialidade: {selectedSpecialty}
              </Badge>
            )}
            {selectedMunicipality && (
              <Badge variant="default">
                Município: {selectedMunicipality}
              </Badge>
            )}
            {!selectedSpecialty && !selectedMunicipality && (
              <p className="text-muted-foreground">Nenhum filtro ativo</p>
            )}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

export default StatsFilters