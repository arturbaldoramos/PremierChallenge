import React from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { 
  Building2, 
  Stethoscope, 
  MapPin, 
  Users,
  RefreshCw,
  TrendingUp,
  Activity
} from 'lucide-react'
import { 
  useStats, 
  useHospitalsByState, 
  useDoctorsByState, 
  useDoctorDistribution 
} from '@/hooks/useStats'

export function StatsCards() {
  const { data: totals, isLoading: totalsLoading, error: totalsError, refetch: refetchTotals } = useStats()
  const { data: hospitalsByState, isLoading: hospitalsByStateLoading, error: hospitalsByStateError, refetch: refetchHospitalsByState } = useHospitalsByState()
  const { data: doctorsByState, isLoading: doctorsByStateLoading, error: doctorsByStateError, refetch: refetchDoctorsByState } = useDoctorsByState()
  const { data: doctorDistribution, isLoading: doctorDistributionLoading, error: doctorDistributionError, refetch: refetchDoctorDistribution } = useDoctorDistribution()

  const refreshAllData = () => {
    refetchTotals()
    refetchHospitalsByState()
    refetchDoctorsByState()
    refetchDoctorDistribution()
  }

  return (
    <div className="space-y-6">
      {/* Botão de Atualização */}
      <div className="flex justify-end">
        <Button onClick={refreshAllData} variant="outline" size="sm">
          <RefreshCw className="mr-2 h-4 w-4" />
          Atualizar Dados
        </Button>
      </div>

      {/* Cards Principais */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {/* Total de Hospitais */}
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total de Hospitais</CardTitle>
            <Building2 className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            {totalsLoading ? (
              <div className="text-2xl font-bold">Carregando...</div>
            ) : totalsError ? (
              <>
                <div className="text-2xl font-bold text-red-500">Erro</div>
                <p className="text-xs text-red-500">{totalsError}</p>
              </>
            ) : (
              <>
                <div className="text-2xl font-bold">
                  {totals?.total_hospitais.toLocaleString('pt-BR') || '0'}
                </div>
                <p className="text-xs text-muted-foreground">
                  Hospitais cadastrados
                </p>
              </>
            )}
          </CardContent>
        </Card>

        {/* Total de Médicos */}
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total de Médicos</CardTitle>
            <Stethoscope className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            {totalsLoading ? (
              <div className="text-2xl font-bold">Carregando...</div>
            ) : totalsError ? (
              <>
                <div className="text-2xl font-bold text-red-500">Erro</div>
                <p className="text-xs text-red-500">{totalsError}</p>
              </>
            ) : (
              <>
                <div className="text-2xl font-bold">
                  {totals?.total_medicos.toLocaleString('pt-BR') || '0'}
                </div>
                <p className="text-xs text-muted-foreground">
                  Médicos cadastrados
                </p>
              </>
            )}
          </CardContent>
        </Card>

        {/* Estados */}
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Estados</CardTitle>
            <MapPin className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            {totalsLoading ? (
              <div className="text-2xl font-bold">Carregando...</div>
            ) : (
              <>
                <div className="text-2xl font-bold">
                  {totals?.total_estados || '0'}
                </div>
                <p className="text-xs text-muted-foreground">
                  Estados brasileiros
                </p>
              </>
            )}
          </CardContent>
        </Card>

        {/* Municípios */}
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Municípios</CardTitle>
            <MapPin className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            {totalsLoading ? (
              <div className="text-2xl font-bold">Carregando...</div>
            ) : (
              <>
                <div className="text-2xl font-bold">
                  {totals?.total_municipios.toLocaleString('pt-BR') || '0'}
                </div>
                <p className="text-xs text-muted-foreground">
                  Municípios brasileiros
                </p>
              </>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Top Estados com Hospitais */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Building2 className="h-5 w-5" />
            Estados com Mais Hospitais
          </CardTitle>
        </CardHeader>
        <CardContent>
          {hospitalsByStateLoading ? (
            <p>Carregando dados...</p>
          ) : hospitalsByStateError ? (
            <p className="text-red-500">Erro: {hospitalsByStateError}</p>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
              {hospitalsByState?.slice(0, 9).map((item, index) => (
                <div key={index} className="flex items-center justify-between p-3 border rounded-lg">
                  <div>
                    <p className="font-medium">{item.nome_estado}</p>
                    <p className="text-sm text-muted-foreground">{item.estado}</p>
                  </div>
                  <Badge variant="secondary" className="text-lg font-bold">
                    {item.count}
                  </Badge>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Top Estados com Médicos */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Stethoscope className="h-5 w-5" />
            Estados com Mais Médicos
          </CardTitle>
        </CardHeader>
        <CardContent>
          {doctorsByStateLoading ? (
            <p>Carregando dados...</p>
          ) : doctorsByStateError ? (
            <p className="text-red-500">Erro: {doctorsByStateError}</p>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
              {doctorsByState?.slice(0, 9).map((item, index) => (
                <div key={index} className="flex items-center justify-between p-3 border rounded-lg">
                  <div>
                    <p className="font-medium">{item.nome_estado}</p>
                    <p className="text-sm text-muted-foreground">{item.estado}</p>
                  </div>
                  <Badge variant="secondary" className="text-lg font-bold">
                    {item.count.toLocaleString('pt-BR')}
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
          <CardTitle className="flex items-center gap-2">
            <Users className="h-5 w-5" />
            Distribuição de Médicos por Hospitais
          </CardTitle>
        </CardHeader>
        <CardContent>
          {doctorDistributionLoading ? (
            <p>Carregando dados...</p>
          ) : doctorDistributionError ? (
            <p className="text-red-500">Erro: {doctorDistributionError}</p>
          ) : doctorDistribution ? (
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              <div className="text-center p-4 bg-blue-50 rounded-lg">
                <div className="text-2xl font-bold text-blue-600">
                  {doctorDistribution.com_um_hospital.toLocaleString('pt-BR')}
                </div>
                <div className="text-sm text-blue-800">Com 1 Hospital</div>
              </div>
              <div className="text-center p-4 bg-green-50 rounded-lg">
                <div className="text-2xl font-bold text-green-600">
                  {doctorDistribution.com_dois_hospitais.toLocaleString('pt-BR')}
                </div>
                <div className="text-sm text-green-800">Com 2 Hospitais</div>
              </div>
              <div className="text-center p-4 bg-yellow-50 rounded-lg">
                <div className="text-2xl font-bold text-yellow-600">
                  {doctorDistribution.com_tres_hospitais.toLocaleString('pt-BR')}
                </div>
                <div className="text-sm text-yellow-800">Com 3 Hospitais</div>
              </div>
              <div className="text-center p-4 bg-red-50 rounded-lg">
                <div className="text-2xl font-bold text-red-600">
                  {doctorDistribution.sem_hospitais.toLocaleString('pt-BR')}
                </div>
                <div className="text-sm text-red-800">Sem Hospitais</div>
              </div>
            </div>
          ) : null}
        </CardContent>
      </Card>
    </div>
  )
}

export default StatsCards