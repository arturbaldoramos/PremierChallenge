import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"
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
  PieChart
} from "lucide-react"
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, PieChart as RechartsPieChart, Pie, Cell } from "recharts"
import { useState } from "react"

// Dados mockados para demonstração
const estados = [
  { value: "sp", label: "São Paulo" },
  { value: "rj", label: "Rio de Janeiro" },
  { value: "mg", label: "Minas Gerais" },
  { value: "rs", label: "Rio Grande do Sul" },
  { value: "pr", label: "Paraná" },
]

const cidades = [
  { value: "sao-paulo", label: "São Paulo" },
  { value: "rio-de-janeiro", label: "Rio de Janeiro" },
  { value: "belo-horizonte", label: "Belo Horizonte" },
  { value: "porto-alegre", label: "Porto Alegre" },
  { value: "curitiba", label: "Curitiba" },
]

const hospitais = [
  { value: "hospital-central", label: "Hospital Central" },
  { value: "hospital-universitario", label: "Hospital Universitário" },
  { value: "hospital-municipal", label: "Hospital Municipal" },
  { value: "hospital-privado", label: "Hospital Privado" },
  { value: "hospital-regional", label: "Hospital Regional" },
]

const doencas = [
  { value: "hipertensao", label: "Hipertensão" },
  { value: "diabetes", label: "Diabetes" },
  { value: "resfriado", label: "Resfriado Comum" },
  { value: "dor-cabeca", label: "Dor de Cabeça" },
  { value: "ansiedade", label: "Ansiedade" },
  { value: "gripe", label: "Gripe" },
  { value: "asma", label: "Asma" },
]

// Dados para gráficos
const dadosAtendimentosPorRegiao = [
  { regiao: "Centro", atendimentos: 847, casos: 234 },
  { regiao: "Norte", atendimentos: 623, casos: 187 },
  { regiao: "Sul", atendimentos: 512, casos: 156 },
  { regiao: "Leste", atendimentos: 398, casos: 134 },
  { regiao: "Oeste", atendimentos: 334, casos: 98 },
]

const dadosDoencasRecorrentes = [
  { name: "Hipertensão", value: 234, color: "#ef4444" },
  { name: "Diabetes", value: 187, color: "#f97316" },
  { name: "Resfriado", value: 156, color: "#eab308" },
  { name: "Dor de Cabeça", value: 134, color: "#22c55e" },
  { name: "Ansiedade", value: 98, color: "#3b82f6" },
]

const dadosEficienciaHospital = [
  { hospital: "Hospital Central", eficiencia: 92, tempoMedio: 45 },
  { hospital: "Hospital Universitário", eficiencia: 88, tempoMedio: 52 },
  { hospital: "Hospital Municipal", eficiencia: 85, tempoMedio: 38 },
  { hospital: "Hospital Privado", eficiencia: 94, tempoMedio: 42 },
  { hospital: "Hospital Regional", eficiencia: 79, tempoMedio: 58 },
]

const chartConfig = {
  atendimentos: {
    label: "Atendimentos",
    color: "hsl(var(--chart-1))",
  },
  casos: {
    label: "Casos",
    color: "hsl(var(--chart-2))",
  },
  eficiencia: {
    label: "Eficiência (%)",
    color: "hsl(var(--chart-1))",
  },
  tempoMedio: {
    label: "Tempo Médio (min)",
    color: "hsl(var(--chart-2))",
  },
}

export default function Dashboard() {
  const [filtros, setFiltros] = useState({
    estado: "todos",
    cidade: "todos", 
    hospital: "todos",
    doenca: "todos"
  })

  const handleFiltroChange = (tipo: string, valor: string) => {
    setFiltros(prev => ({ ...prev, [tipo]: valor }))
  }

  return (
    <div className="flex-1 space-y-4 p-4 md:p-8 pt-6">
      <div className="flex items-center justify-between space-y-2">
        <h2 className="text-3xl font-bold tracking-tight">Dashboard de Gestão de Saúde</h2>
        <div className="flex items-center space-x-2">
          <Button variant="outline" size="sm">
            <Download className="mr-2 h-4 w-4" />
            Relatório
          </Button>
        </div>
      </div>
      
      {/* Filtros */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Filter className="h-5 w-5" />
            Filtros de Análise
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
            <div className="space-y-2">
              <label className="text-sm font-medium">Estado</label>
              <Select value={filtros.estado} onValueChange={(value) => handleFiltroChange("estado", value)}>
                <SelectTrigger>
                  <SelectValue placeholder="Selecione o estado" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="todos">Todos os Estados</SelectItem>
                  {estados.map((estado) => (
                    <SelectItem key={estado.value} value={estado.value}>
                      {estado.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            
            <div className="space-y-2">
              <label className="text-sm font-medium">Cidade</label>
              <Select value={filtros.cidade} onValueChange={(value) => handleFiltroChange("cidade", value)}>
                <SelectTrigger>
                  <SelectValue placeholder="Selecione a cidade" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="todos">Todas as Cidades</SelectItem>
                  {cidades.map((cidade) => (
                    <SelectItem key={cidade.value} value={cidade.value}>
                      {cidade.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            
            <div className="space-y-2">
              <label className="text-sm font-medium">Hospital</label>
              <Select value={filtros.hospital} onValueChange={(value) => handleFiltroChange("hospital", value)}>
                <SelectTrigger>
                  <SelectValue placeholder="Selecione o hospital" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="todos">Todos os Hospitais</SelectItem>
                  {hospitais.map((hospital) => (
                    <SelectItem key={hospital.value} value={hospital.value}>
                      {hospital.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            
            <div className="space-y-2">
              <label className="text-sm font-medium">Doença</label>
              <Select value={filtros.doenca} onValueChange={(value) => handleFiltroChange("doenca", value)}>
                <SelectTrigger>
                  <SelectValue placeholder="Selecione a doença" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="todos">Todas as Doenças</SelectItem>
                  {doencas.map((doenca) => (
                    <SelectItem key={doenca.value} value={doenca.value}>
                      {doenca.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>
        </CardContent>
      </Card>
      
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Taxa de Ocupação dos Leitos
            </CardTitle>
            <Activity className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">87.3%</div>
            <p className="text-xs text-muted-foreground">
              +2.1% em relação ao mês anterior
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Pacientes Ativos
            </CardTitle>
            <Users className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">1,234</div>
            <p className="text-xs text-muted-foreground">
              +8.5% em relação ao mês anterior
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total de Hospitais</CardTitle>
            <BarChart3 className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">127</div>
            <p className="text-xs text-muted-foreground">
              +3 novos hospitais este mês
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Casos Urgentes</CardTitle>
            <AlertTriangle className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">23</div>
            <p className="text-xs text-muted-foreground">
              -12% em relação à semana anterior
            </p>
          </CardContent>
        </Card>
      </div>
      
      {/* Gráficos Principais */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
        <Card className="col-span-4">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <BarChart3 className="h-5 w-5" />
              Atendimentos por Região
            </CardTitle>
          </CardHeader>
          <CardContent>
            <ChartContainer config={chartConfig} className="h-[300px]">
              <BarChart data={dadosAtendimentosPorRegiao}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis 
                  dataKey="regiao" 
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
                  dataKey="atendimentos" 
                  fill="var(--color-atendimentos)"
                  radius={[4, 4, 0, 0]}
                />
              </BarChart>
            </ChartContainer>
          </CardContent>
        </Card>
        
        <Card className="col-span-3">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <PieChart className="h-5 w-5" />
              Doenças Mais Recorrentes
            </CardTitle>
          </CardHeader>
          <CardContent>
            <ChartContainer config={chartConfig} className="h-[300px]">
              <RechartsPieChart>
                <ChartTooltip content={<ChartTooltipContent />} />
                <ChartLegend content={<ChartLegendContent />} />
                <Pie
                  data={dadosDoencasRecorrentes}
                  dataKey="value"
                  nameKey="name"
                  cx="50%"
                  cy="50%"
                  outerRadius={80}
                  fill="#8884d8"
                >
                  {dadosDoencasRecorrentes.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={entry.color} />
                  ))}
                </Pie>
              </RechartsPieChart>
            </ChartContainer>
          </CardContent>
        </Card>
      </div>
      
      {/* Gráficos Secundários */}
      <div className="grid gap-4 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <TrendingUp className="h-5 w-5" />
              Eficiência Operacional por Hospital
            </CardTitle>
          </CardHeader>
          <CardContent>
            <ChartContainer config={chartConfig} className="h-[300px]">
              <BarChart data={dadosEficienciaHospital}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis 
                  dataKey="hospital" 
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
                  dataKey="eficiencia" 
                  fill="var(--color-atendimentos)"
                  radius={[4, 4, 0, 0]}
                />
                <Bar 
                  dataKey="tempoMedio" 
                  fill="var(--color-casos)"
                  radius={[4, 4, 0, 0]}
                />
              </BarChart>
            </ChartContainer>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Activity className="h-5 w-5" />
              Médicos Mais Ativos
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-6">
              <div className="flex items-center">
                <Avatar className="h-9 w-9">
                  <AvatarFallback>DS</AvatarFallback>
                </Avatar>
                <div className="ml-4 space-y-1">
                  <p className="text-sm font-medium leading-none">
                    Dr. Silva
                  </p>
                  <p className="text-sm text-muted-foreground">
                    Cardiologia
                  </p>
                </div>
                <div className="ml-auto font-medium">127 atendimentos</div>
              </div>
              
              <div className="flex items-center">
                <Avatar className="h-9 w-9">
                  <AvatarFallback>MC</AvatarFallback>
                </Avatar>
                <div className="ml-4 space-y-1">
                  <p className="text-sm font-medium leading-none">
                    Dra. Costa
                  </p>
                  <p className="text-sm text-muted-foreground">
                    Pediatria
                  </p>
                </div>
                <div className="ml-auto font-medium">98 atendimentos</div>
              </div>
              
              <div className="flex items-center">
                <Avatar className="h-9 w-9">
                  <AvatarFallback>RS</AvatarFallback>
                </Avatar>
                <div className="ml-4 space-y-1">
                  <p className="text-sm font-medium leading-none">
                    Dr. Santos
                  </p>
                  <p className="text-sm text-muted-foreground">
                    Clínica Geral
                  </p>
                </div>
                <div className="ml-auto font-medium">89 atendimentos</div>
              </div>
              
              <div className="flex items-center">
                <Avatar className="h-9 w-9">
                  <AvatarFallback>AO</AvatarFallback>
                </Avatar>
                <div className="ml-4 space-y-1">
                  <p className="text-sm font-medium leading-none">
                    Dra. Oliveira
                  </p>
                  <p className="text-sm text-muted-foreground">
                    Ginecologia
                  </p>
                </div>
                <div className="ml-auto font-medium">76 atendimentos</div>
              </div>
              
              <div className="flex items-center">
                <Avatar className="h-9 w-9">
                  <AvatarFallback>FP</AvatarFallback>
                </Avatar>
                <div className="ml-4 space-y-1">
                  <p className="text-sm font-medium leading-none">
                    Dr. Pereira
                  </p>
                  <p className="text-sm text-muted-foreground">
                    Ortopedia
                  </p>
                </div>
                <div className="ml-auto font-medium">65 atendimentos</div>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
      
      {/* Resumo dos Filtros Aplicados */}
      <Card>
        <CardHeader>
          <CardTitle>Resumo dos Filtros Aplicados</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4 text-sm">
            <div className="space-y-1">
              <span className="text-muted-foreground">Estado:</span>
              <p className="font-medium">
                {filtros.estado === "todos" ? "Todos os Estados" : 
                 estados.find(e => e.value === filtros.estado)?.label}
              </p>
            </div>
            <div className="space-y-1">
              <span className="text-muted-foreground">Cidade:</span>
              <p className="font-medium">
                {filtros.cidade === "todos" ? "Todas as Cidades" : 
                 cidades.find(c => c.value === filtros.cidade)?.label}
              </p>
            </div>
            <div className="space-y-1">
              <span className="text-muted-foreground">Hospital:</span>
              <p className="font-medium">
                {filtros.hospital === "todos" ? "Todos os Hospitais" : 
                 hospitais.find(h => h.value === filtros.hospital)?.label}
              </p>
            </div>
            <div className="space-y-1">
              <span className="text-muted-foreground">Doença:</span>
              <p className="font-medium">
                {filtros.doenca === "todos" ? "Todas as Doenças" : 
                 doencas.find(d => d.value === filtros.doenca)?.label}
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
              <span className="font-semibold">247</span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">Enfermeiros</span>
              <span className="font-semibold">1,156</span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">Técnicos</span>
              <span className="font-semibold">892</span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">Taxa de Absenteísmo</span>
              <span className="font-semibold text-orange-600">4.2%</span>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <BarChart3 className="h-5 w-5" />
              Performance Financeira
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">Receita Mensal</span>
              <span className="font-semibold text-green-600">R$ 2.4M</span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">Custo por Paciente</span>
              <span className="font-semibold">R$ 1,247</span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">Margem de Lucro</span>
              <span className="font-semibold text-green-600">12.8%</span>
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
