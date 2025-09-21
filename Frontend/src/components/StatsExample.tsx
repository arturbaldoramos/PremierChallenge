import React from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { 
  useStats, 
  useHospitalsByState, 
  useDoctorsByState, 
  useHospitalsBySpecialty, 
  useHospitalsByMunicipality, 
  useDoctorDistribution, 
  useSpecialties 
} from '@/hooks/useStats'

export function StatsExample() {
  // Hooks para consumir todos os endpoints
  const { data: totals, isLoading: totalsLoading, error: totalsError, refetch: refetchTotals } = useStats()
  const { data: hospitalsByState, isLoading: hospitalsByStateLoading, error: hospitalsByStateError, refetch: refetchHospitalsByState } = useHospitalsByState()
  const { data: doctorsByState, isLoading: doctorsByStateLoading, error: doctorsByStateError, refetch: refetchDoctorsByState } = useDoctorsByState()
  const { data: hospitalsBySpecialty, isLoading: hospitalsBySpecialtyLoading, error: hospitalsBySpecialtyError, refetch: refetchHospitalsBySpecialty } = useHospitalsBySpecialty()
  const { data: hospitalsByMunicipality, isLoading: hospitalsByMunicipalityLoading, error: hospitalsByMunicipalityError, refetch: refetchHospitalsByMunicipality } = useHospitalsByMunicipality()
  const { data: doctorDistribution, isLoading: doctorDistributionLoading, error: doctorDistributionError, refetch: refetchDoctorDistribution } = useDoctorDistribution()
  const { data: specialties, isLoading: specialtiesLoading, error: specialtiesError, refetch: refetchSpecialties } = useSpecialties()

  // Estados para filtros
  const [selectedSpecialty, setSelectedSpecialty] = React.useState('')
  const [selectedMunicipality, setSelectedMunicipality] = React.useState('')

  // Funções para aplicar filtros
  const handleSpecialtyFilter = () => {
    refetchHospitalsBySpecialty(selectedSpecialty || undefined)
  }

  const handleMunicipalityFilter = () => {
    refetchHospitalsByMunicipality(selectedMunicipality || undefined)
  }

  return (
    <div className="space-y-6 p-6">
      <h1 className="text-3xl font-bold">Exemplo de Consumo dos Endpoints de Estatísticas</h1>
      
      {/* Estatísticas Gerais */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center justify-between">
            Estatísticas Gerais
            <Button onClick={refetchTotals} variant="outline" size="sm">
              Atualizar
            </Button>
          </CardTitle>
        </CardHeader>
        <CardContent>
          {totalsLoading && <p>Carregando...</p>}
          {totalsError && <p className="text-red-500">Erro: {totalsError}</p>}
          {totals && (
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              <div className="text-center">
                <div className="text-2xl font-bold text-blue-600">{totals.total_hospitais.toLocaleString()}</div>
                <div className="text-sm text-gray-600">Hospitais</div>
              </div>
              <div className="text-center">
                <div className="text-2xl font-bold text-green-600">{totals.total_medicos.toLocaleString()}</div>
                <div className="text-sm text-gray-600">Médicos</div>
              </div>
              <div className="text-center">
                <div className="text-2xl font-bold text-purple-600">{totals.total_estados}</div>
                <div className="text-sm text-gray-600">Estados</div>
              </div>
              <div className="text-center">
                <div className="text-2xl font-bold text-orange-600">{totals.total_municipios.toLocaleString()}</div>
                <div className="text-sm text-gray-600">Municípios</div>
              </div>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Hospitais por Estado */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center justify-between">
            Hospitais por Estado
            <Button onClick={refetchHospitalsByState} variant="outline" size="sm">
              Atualizar
            </Button>
          </CardTitle>
        </CardHeader>
        <CardContent>
          {hospitalsByStateLoading && <p>Carregando...</p>}
          {hospitalsByStateError && <p className="text-red-500">Erro: {hospitalsByStateError}</p>}
          {hospitalsByState && (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-2">
              {hospitalsByState.slice(0, 9).map((item, index) => (
                <div key={index} className="flex items-center justify-between p-2 border rounded">
                  <span className="font-medium">{item.nome_estado}</span>
                  <Badge variant="secondary">{item.count}</Badge>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Médicos por Estado */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center justify-between">
            Médicos por Estado
            <Button onClick={refetchDoctorsByState} variant="outline" size="sm">
              Atualizar
            </Button>
          </CardTitle>
        </CardHeader>
        <CardContent>
          {doctorsByStateLoading && <p>Carregando...</p>}
          {doctorsByStateError && <p className="text-red-500">Erro: {doctorsByStateError}</p>}
          {doctorsByState && (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-2">
              {doctorsByState.slice(0, 9).map((item, index) => (
                <div key={index} className="flex items-center justify-between p-2 border rounded">
                  <span className="font-medium">{item.nome_estado}</span>
                  <Badge variant="secondary">{item.count.toLocaleString()}</Badge>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Hospitais por Especialidade */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center justify-between">
            Hospitais por Especialidade
            <Button onClick={refetchHospitalsBySpecialty} variant="outline" size="sm">
              Atualizar
            </Button>
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="mb-4 flex gap-2">
            <Input
              placeholder="Digite uma especialidade para filtrar"
              value={selectedSpecialty}
              onChange={(e) => setSelectedSpecialty(e.target.value)}
            />
            <Button onClick={handleSpecialtyFilter} variant="outline">
              Filtrar
            </Button>
          </div>
          
          {hospitalsBySpecialtyLoading && <p>Carregando...</p>}
          {hospitalsBySpecialtyError && <p className="text-red-500">Erro: {hospitalsBySpecialtyError}</p>}
          {hospitalsBySpecialty && (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-2">
              {Array.isArray(hospitalsBySpecialty) && hospitalsBySpecialty.slice(0, 9).map((item, index) => (
                <div key={index} className="flex items-center justify-between p-2 border rounded">
                  <span className="font-medium">
                    {'especialidade' in item ? item.especialidade : item.nome}
                  </span>
                  <Badge variant="secondary">
                    {'count' in item ? item.count : 'Hospital'}
                  </Badge>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Hospitais por Município */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center justify-between">
            Hospitais por Município
            <Button onClick={refetchHospitalsByMunicipality} variant="outline" size="sm">
              Atualizar
            </Button>
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="mb-4 flex gap-2">
            <Input
              placeholder="Digite um município para filtrar"
              value={selectedMunicipality}
              onChange={(e) => setSelectedMunicipality(e.target.value)}
            />
            <Button onClick={handleMunicipalityFilter} variant="outline">
              Filtrar
            </Button>
          </div>
          
          {hospitalsByMunicipalityLoading && <p>Carregando...</p>}
          {hospitalsByMunicipalityError && <p className="text-red-500">Erro: {hospitalsByMunicipalityError}</p>}
          {hospitalsByMunicipality && (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-2">
              {Array.isArray(hospitalsByMunicipality) && hospitalsByMunicipality.slice(0, 9).map((item, index) => (
                <div key={index} className="flex items-center justify-between p-2 border rounded">
                  <span className="font-medium">
                    {'municipio' in item ? item.municipio : item.nome}
                  </span>
                  <Badge variant="secondary">
                    {'count' in item ? item.count : 'Hospital'}
                  </Badge>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Distribuição de Médicos */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center justify-between">
            Distribuição de Médicos
            <Button onClick={refetchDoctorDistribution} variant="outline" size="sm">
              Atualizar
            </Button>
          </CardTitle>
        </CardHeader>
        <CardContent>
          {doctorDistributionLoading && <p>Carregando...</p>}
          {doctorDistributionError && <p className="text-red-500">Erro: {doctorDistributionError}</p>}
          {doctorDistribution && (
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              <div className="text-center">
                <div className="text-xl font-bold text-blue-600">{doctorDistribution.com_um_hospital.toLocaleString()}</div>
                <div className="text-sm text-gray-600">Com 1 Hospital</div>
              </div>
              <div className="text-center">
                <div className="text-xl font-bold text-green-600">{doctorDistribution.com_dois_hospitais.toLocaleString()}</div>
                <div className="text-sm text-gray-600">Com 2 Hospitais</div>
              </div>
              <div className="text-center">
                <div className="text-xl font-bold text-purple-600">{doctorDistribution.com_tres_hospitais.toLocaleString()}</div>
                <div className="text-sm text-gray-600">Com 3 Hospitais</div>
              </div>
              <div className="text-center">
                <div className="text-xl font-bold text-red-600">{doctorDistribution.sem_hospitais.toLocaleString()}</div>
                <div className="text-sm text-gray-600">Sem Hospitais</div>
              </div>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Lista de Especialidades */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center justify-between">
            Especialidades Disponíveis
            <Button onClick={refetchSpecialties} variant="outline" size="sm">
              Atualizar
            </Button>
          </CardTitle>
        </CardHeader>
        <CardContent>
          {specialtiesLoading && <p>Carregando...</p>}
          {specialtiesError && <p className="text-red-500">Erro: {specialtiesError}</p>}
          {specialties && (
            <div className="flex flex-wrap gap-2">
              {specialties.map((specialty, index) => (
                <Badge key={index} variant="outline">
                  {specialty}
                </Badge>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )
}

export default StatsExample