import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { 
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
  ChartLegend,
  ChartLegendContent,
} from "@/components/ui/chart"
import { 
  Users, 
  Activity,
  Download,
  TrendingUp,
  AlertTriangle,
  Filter,
  BarChart3,
  PieChart,
  MapPin,
  Building2,
  Stethoscope,
  RefreshCw
} from "lucide-react"
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, PieChart as RechartsPieChart, Pie, Cell, ResponsiveContainer } from "recharts"
import { useState } from "react"
import { 
  useStats, 
  useHospitalsByState, 
  useDoctorsByState, 
  useHospitalsBySpecialty, 
  useHospitalsByMunicipality, 
  useDoctorDistribution, 
  useSpecialties 
} from "@/hooks/useStats"

// Configuração dos gráficos
const chartConfig = {
  count: {
    label: "Quantidade",
    color: "hsl(var(--chart-1))",
  },
  hospitals: {
    label: "Hospitais",
    color: "hsl(var(--chart-2))",
  },
  doctors: {
    label: "Médicos",
    color: "hsl(var(--chart-3))",
  },
}

export default function DashboardWithRealData() {
  const [selectedState, setSelectedState] = useState<string>("")
  const [selectedSpecialty, setSelectedSpecialty] = useState<string>("")
  const [selectedMunicipality, setSelectedMunicipality] = useState<string>("")

  // Hooks para consumir todos os endpoints
  const { data: totals, isLoading: totalsLoading, error: totalsError, refetch: refetchTotals } = useStats()
  const { data: hospitalsByState, isLoading: hospitalsByStateLoading, error: hospitalsByStateError, refetch: refetchHospitalsByState } = useHospitalsByState()
  const { data: doctorsByState, isLoading: doctorsByStateLoading, error: doctorsByStateError, refetch: refetchDoctorsByState } = useDoctorsByState()
  const { data: hospitalsBySpecialty, isLoading: hospitalsBySpecialtyLoading, error: hospitalsBySpecialtyError, refetch: refetchHospitalsBySpecialty } = useHospitalsBySpecialty()
  const { data: hospitalsByMunicipality, isLoading: hospitalsByMunicipalityLoading, error: hospitalsByMunicipalityError, refetch: refetchHospitalsByMunicipality } = useHospitalsByMunicipality()
  const { data: doctorDistribution, isLoading: doctorDistributionLoading, error: doctorDistributionError, refetch: refetchDoctorDistribution } = useDoctorDistribution()
  const { data: specialties, isLoading: specialtiesLoading, error: specialtiesError, refetch: refetchSpecialties } = useSpecialties()

  // Função para atualizar todos os dados
  const refreshAllData = () => {
    refetchTotals()
    refetchHospitalsByState()
    refetchDoctorsByState()
    refetchHospitalsBySpecialty()
    refetchHospitalsByMunicipality()
    refetchDoctorDistribution()
    refetchSpecialties()
  }

  // Preparar dados para gráficos
  const hospitalsChartData = hospitalsByState?.slice(0, 10).map(item => ({
    estado: item.estado,
    hospitais: item.count,
    nome: item.nome_estado
  })) || []

  const doctorsChartData = doctorsByState?.slice(0, 10).map(item => ({
    estado: item.estado,
    medicos: item.count,
    nome: item.nome_estado
  })) || []

  const specialtyChartData = hospitalsBySpecialty?.slice(0, 8).map(item => {
    if ('especialidade' in item) {
      return {
        name: item.especialidade,
        value: item.count,
        color: `hsl(${Math.random() * 360}, 70%, 50%)`
      }
    }
    return null
  }).filter(Boolean) || []

  const distributionChartData = doctorDistribution ? [
    { name: "1 Hospital", value: doctorDistribution.com_um_hospital, color: "#3b82f6" },
    { name: "2 Hospitais", value: doctorDistribution.com_dois_hospitais, color: "#10b981" },
    { name: "3 Hospitais", value: doctorDistribution.com_tres_hospitais, color: "#f59e0b" },
    { name: "Sem Hospitais", value: doctorDistribution.sem_hospitais, color: "#ef4444" },
  ] : []

  return (
    <div className="flex-1 space-y-4 p-4 md:p-8 pt-6">
      <div className="flex items-center justify-between space-y-2">
        <h2 className="text-3xl font-bold tracking-tight">Dashboard de Gestão de Saúde</h2>
        <div className="flex items-center space-x-2">
          <Button onClick={refreshAllData} variant="outline" size="sm">
            <RefreshCw className="mr-2 h-4 w-4" />
            Atualizar Dados
          </Button>
          <Button variant="outline" size="sm">
            <Download className="mr-2 h-4 w-4" />
            Relatório
          </Button>
        </div>
      </div>
      
      {/* Filtros Dinâmicos */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Filter className="h-5 w-5" />
            Filtros de Análise
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="space-y-2">
              <label className="text-sm font-medium">Estado</label>
              <Select value={selectedState} onValueChange={setSelectedState}>
                <SelectTrigger>
                  <SelectValue placeholder="Selecione um estado" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="">Todos os Estados</SelectItem>
                  {hospitalsByState?.map((estado) => (
                    <SelectItem key={estado.estado} value={estado.estado}>
                      {estado.nome_estado} ({estado.count} hospitais)
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            
            <div className="space-y-2">
              <label className="text-sm font-medium">Especialidade</label>
              <Select value={selectedSpecialty} onValueChange={setSelectedSpecialty}>
                <SelectTrigger>
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
            </div>
            
            <div className="space-y-2">
              <label className="text-sm font-medium">Município</label>
              <Select value={selectedMunicipality} onValueChange={setSelectedMunicipality}>
                <SelectTrigger>
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
            </div>
          </div>
        </CardContent>
      </Card>
      
      {/* Cards de Estatísticas Principais */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Total de Hospitais
            </CardTitle>
            <Building2 className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            {totalsLoading ? (
              <div className="text-2xl font-bold">Carregando...</div>
            ) : totalsError ? (
              <>
                <div className="text-2xl font-bold">
                  {totals?.total_hospitais.toLocaleString('pt-BR') || '0'}
                </div>
                <p className="text-xs text-orange-500">
                  {totalsError}
                </p>
              </>
            ) : (
              <>
                <div className="text-2xl font-bold">
                  {totals?.total_hospitais.toLocaleString('pt-BR') || '0'}
                </div>
                <p className="text-xs text-muted-foreground">
                  Hospitais cadastrados no sistema
                </p>
              </>
            )}
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Total de Médicos
            </CardTitle>
            <Stethoscope className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            {totalsLoading ? (
              <div className="text-2xl font-bold">Carregando...</div>
            ) : totalsError ? (
              <>
                <div className="text-2xl font-bold">
                  {totals?.total_medicos.toLocaleString('pt-BR') || '0'}
                </div>
                <p className="text-xs text-orange-500">
                  {totalsError}
                </p>
              </>
            ) : (
              <>
                <div className="text-2xl font-bold">
                  {totals?.total_medicos.toLocaleString('pt-BR') || '0'}
                </div>
                <p className="text-xs text-muted-foreground">
                  Médicos cadastrados no sistema
                </p>
              </>
            )}
          </CardContent>
        </Card>
        
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
      
      {/* Gráficos Principais */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
        <Card className="col-span-4">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <BarChart3 className="h-5 w-5" />
              Hospitais por Estado (Top 10)
            </CardTitle>
          </CardHeader>
          <CardContent>
            {hospitalsByStateLoading ? (
              <div className="h-[300px] flex items-center justify-center">
                <p>Carregando dados...</p>
              </div>
            ) : hospitalsByStateError ? (
              <div className="h-[300px] flex items-center justify-center text-red-500">
                <p>Erro: {hospitalsByStateError}</p>
              </div>
            ) : (
              <ChartContainer config={chartConfig} className="h-[300px]">
                <BarChart data={hospitalsChartData}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis 
                    dataKey="estado" 
                    tick={{ fontSize: 12 }}
                    axisLine={false}
                    tickLine={false}
                  />
                  <YAxis 
                    tick={{ fontSize: 12 }}
                    axisLine={false}
                    tickLine={false}
                  />
                  <ChartTooltip content={<ChartTooltipContent />} />
                  <ChartLegend content={<ChartLegendContent />} />
                  <Bar 
                    dataKey="hospitais" 
                    fill="var(--color-count)"
                    radius={[4, 4, 0, 0]}
                  />
                </BarChart>
              </ChartContainer>
            )}
          </CardContent>
        </Card>
        
        <Card className="col-span-3">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <PieChart className="h-5 w-5" />
              Distribuição de Médicos
            </CardTitle>
          </CardHeader>
          <CardContent>
            {doctorDistributionLoading ? (
              <div className="h-[300px] flex items-center justify-center">
                <p>Carregando dados...</p>
              </div>
            ) : doctorDistributionError ? (
              <div className="h-[300px] flex items-center justify-center text-red-500">
                <p>Erro: {doctorDistributionError}</p>
              </div>
            ) : (
              <ChartContainer config={chartConfig} className="h-[300px]">
                <RechartsPieChart>
                  <ChartTooltip content={<ChartTooltipContent />} />
                  <ChartLegend content={<ChartLegendContent />} />
                  <Pie
                    data={distributionChartData}
                    dataKey="value"
                    nameKey="name"
                    cx="50%"
                    cy="50%"
                    outerRadius={80}
                    fill="#8884d8"
                  >
                    {distributionChartData.map((entry, index) => (
                      <Cell key={`cell-${index}`} fill={entry.color} />
                    ))}
                  </Pie>
                </RechartsPieChart>
              </ChartContainer>
            )}
          </CardContent>
        </Card>
      </div>
      
      {/* Gráficos Secundários */}
      <div className="grid gap-4 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <TrendingUp className="h-5 w-5" />
              Médicos por Estado (Top 10)
            </CardTitle>
          </CardHeader>
          <CardContent>
            {doctorsByStateLoading ? (
              <div className="h-[300px] flex items-center justify-center">
                <p>Carregando dados...</p>
              </div>
            ) : doctorsByStateError ? (
              <div className="h-[300px] flex items-center justify-center text-red-500">
                <p>Erro: {doctorsByStateError}</p>
              </div>
            ) : (
              <ChartContainer config={chartConfig} className="h-[300px]">
                <BarChart data={doctorsChartData}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis 
                    dataKey="estado" 
                    tick={{ fontSize: 12 }}
                    axisLine={false}
                    tickLine={false}
                  />
                  <YAxis 
                    tick={{ fontSize: 12 }}
                    axisLine={false}
                    tickLine={false}
                  />
                  <ChartTooltip content={<ChartTooltipContent />} />
                  <ChartLegend content={<ChartLegendContent />} />
                  <Bar 
                    dataKey="medicos" 
                    fill="var(--color-doctors)"
                    radius={[4, 4, 0, 0]}
                  />
                </BarChart>
              </ChartContainer>
            )}
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Activity className="h-5 w-5" />
              Hospitais por Especialidade (Top 8)
            </CardTitle>
          </CardHeader>
          <CardContent>
            {hospitalsBySpecialtyLoading ? (
              <div className="h-[300px] flex items-center justify-center">
                <p>Carregando dados...</p>
              </div>
            ) : hospitalsBySpecialtyError ? (
              <div className="h-[300px] flex items-center justify-center text-red-500">
                <p>Erro: {hospitalsBySpecialtyError}</p>
              </div>
            ) : (
              <ChartContainer config={chartConfig} className="h-[300px]">
                <RechartsPieChart>
                  <ChartTooltip content={<ChartTooltipContent />} />
                  <ChartLegend content={<ChartLegendContent />} />
                  <Pie
                    data={specialtyChartData}
                    dataKey="value"
                    nameKey="name"
                    cx="50%"
                    cy="50%"
                    outerRadius={80}
                    fill="#8884d8"
                  >
                    {specialtyChartData.map((entry, index) => (
                      <Cell key={`cell-${index}`} fill={entry.color} />
                    ))}
                  </Pie>
                </RechartsPieChart>
              </ChartContainer>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Lista de Especialidades */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Stethoscope className="h-5 w-5" />
            Especialidades Disponíveis
          </CardTitle>
        </CardHeader>
        <CardContent>
          {specialtiesLoading ? (
            <p>Carregando especialidades...</p>
          ) : specialtiesError ? (
            <p className="text-red-500">Erro: {specialtiesError}</p>
          ) : (
            <div className="flex flex-wrap gap-2">
              {specialties?.map((specialty, index) => (
                <Badge key={index} variant="outline" className="cursor-pointer hover:bg-primary hover:text-primary-foreground">
                  {specialty}
                </Badge>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Resumo dos Filtros Aplicados */}
      <Card>
        <CardHeader>
          <CardTitle>Resumo dos Filtros Aplicados</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
            <div className="space-y-1">
              <span className="text-muted-foreground">Estado:</span>
              <p className="font-medium">
                {selectedState ? 
                  hospitalsByState?.find(e => e.estado === selectedState)?.nome_estado : 
                  "Todos os Estados"
                }
              </p>
            </div>
            <div className="space-y-1">
              <span className="text-muted-foreground">Especialidade:</span>
              <p className="font-medium">
                {selectedSpecialty || "Todas as Especialidades"}
              </p>
            </div>
            <div className="space-y-1">
              <span className="text-muted-foreground">Município:</span>
              <p className="font-medium">
                {selectedMunicipality || "Todos os Municípios"}
              </p>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Métricas de Gestão */}
      <div className="grid gap-4 md:grid-cols-3">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <AlertTriangle className="h-5 w-5" />
              Indicadores Críticos
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">Leitos Disponíveis</span>
              <span className="font-semibold text-green-600">342</span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">Tempo Médio de Espera</span>
              <span className="font-semibold text-orange-600">2.3h</span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">Taxa de Alta</span>
              <span className="font-semibold text-blue-600">94.2%</span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">Reinternações</span>
              <span className="font-semibold text-red-600">3.1%</span>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Users className="h-5 w-5" />
              Recursos Humanos
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">Médicos Ativos</span>
              <span className="font-semibold">
                {doctorDistribution ? 
                  (doctorDistribution.com_um_hospital + doctorDistribution.com_dois_hospitais + doctorDistribution.com_tres_hospitais).toLocaleString('pt-BR') : 
                  '0'
                }
              </span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">Médicos Sem Hospitais</span>
              <span className="font-semibold text-orange-600">
                {doctorDistribution?.sem_hospitais.toLocaleString('pt-BR') || '0'}
              </span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">Taxa de Absenteísmo</span>
              <span className="font-semibold text-orange-600">4.2%</span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">Especialidades</span>
              <span className="font-semibold">
                {specialties?.length || '0'}
              </span>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <BarChart3 className="h-5 w-5" />
              Performance Operacional
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">Estados Cobertos</span>
              <span className="font-semibold text-green-600">
                {hospitalsByState?.length || '0'}
              </span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">Municípios Cobertos</span>
              <span className="font-semibold">
                {hospitalsByMunicipality?.length || '0'}
              </span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">Cobertura Nacional</span>
              <span className="font-semibold text-blue-600">100%</span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">ROI</span>
              <span className="font-semibold text-blue-600">8.4%</span>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}