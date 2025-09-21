# Integração dos Endpoints de Estatísticas no Dashboard

Este documento mostra como integrar os dados dos endpoints de estatísticas no dashboard da aplicação.

## Exemplos Práticos

### 1. **StatsCards.tsx** - Cards de Estatísticas
Componente que mostra os dados principais em cards visuais:

```typescript
import { StatsCards } from '@/components/StatsCards'

// No seu dashboard
<StatsCards />
```

**Funcionalidades:**
- Cards com totais (hospitais, médicos, estados, municípios)
- Top estados com mais hospitais
- Top estados com mais médicos
- Distribuição de médicos por hospitais
- Botão de atualização de dados

### 2. **StatsCharts.tsx** - Gráficos de Dados
Componente com gráficos usando os dados da API:

```typescript
import { StatsCharts } from '@/components/StatsCharts'

// No seu dashboard
<StatsCharts />
```

**Funcionalidades:**
- Gráfico de barras: Hospitais por estado (Top 10)
- Gráfico de barras: Médicos por estado (Top 10)
- Gráfico de pizza: Hospitais por especialidade (Top 8)
- Gráfico de pizza: Distribuição de médicos

### 3. **StatsFilters.tsx** - Filtros Dinâmicos
Componente com filtros interativos:

```typescript
import { StatsFilters } from '@/components/StatsFilters'

// No seu dashboard
<StatsFilters />
```

**Funcionalidades:**
- Filtro por especialidade (dropdown + busca personalizada)
- Filtro por município (dropdown + busca personalizada)
- Resultados em tempo real
- Limpeza de filtros

### 4. **DashboardWithRealData.tsx** - Dashboard Completo
Dashboard completo integrando todos os endpoints:

```typescript
import DashboardWithRealData from '@/pages/DashboardWithRealData'

// Use este componente como página completa
<DashboardWithRealData />
```

## Como Usar no Dashboard Atual

### Opção 1: Substituir o Dashboard Atual
```typescript
// Em src/pages/Dashboard.tsx
import DashboardWithRealData from './DashboardWithRealData'

export default function Dashboard() {
  return <DashboardWithRealData />
}
```

### Opção 2: Integrar Componentes Específicos
```typescript
// Em src/pages/Dashboard.tsx
import { StatsCards, StatsCharts, StatsFilters } from '@/components'

export default function Dashboard() {
  return (
    <div className="flex-1 space-y-4 p-4 md:p-8 pt-6">
      <h2 className="text-3xl font-bold tracking-tight">Dashboard</h2>
      
      {/* Seus componentes existentes */}
      <StatsCards />
      <StatsCharts />
      <StatsFilters />
      
      {/* Outros componentes do dashboard */}
    </div>
  )
}
```

### Opção 3: Usar Hooks Individualmente
```typescript
import { 
  useStats, 
  useHospitalsByState, 
  useDoctorsByState 
} from '@/hooks/useStats'

export default function Dashboard() {
  const { data: totals, isLoading, error } = useStats()
  const { data: hospitalsByState } = useHospitalsByState()
  const { data: doctorsByState } = useDoctorsByState()

  return (
    <div>
      {/* Seus componentes personalizados usando os dados */}
      {totals && (
        <div>
          <p>Hospitais: {totals.total_hospitais}</p>
          <p>Médicos: {totals.total_medicos}</p>
        </div>
      )}
    </div>
  )
}
```

## Exemplos de Uso Específicos

### 1. Card Simples com Total de Hospitais
```typescript
import { useStats } from '@/hooks/useStats'

function HospitalCard() {
  const { data: totals, isLoading, error } = useStats()

  return (
    <Card>
      <CardHeader>
        <CardTitle>Total de Hospitais</CardTitle>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <p>Carregando...</p>
        ) : error ? (
          <p className="text-red-500">Erro: {error}</p>
        ) : (
          <div className="text-2xl font-bold">
            {totals?.total_hospitais.toLocaleString('pt-BR')}
          </div>
        )}
      </CardContent>
    </Card>
  )
}
```

### 2. Lista de Estados com Mais Hospitais
```typescript
import { useHospitalsByState } from '@/hooks/useStats'

function TopStatesList() {
  const { data: hospitalsByState, isLoading, error } = useHospitalsByState()

  return (
    <Card>
      <CardHeader>
        <CardTitle>Estados com Mais Hospitais</CardTitle>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <p>Carregando...</p>
        ) : error ? (
          <p className="text-red-500">Erro: {error}</p>
        ) : (
          <div className="space-y-2">
            {hospitalsByState?.slice(0, 5).map((item, index) => (
              <div key={index} className="flex justify-between">
                <span>{item.nome_estado}</span>
                <Badge>{item.count}</Badge>
              </div>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  )
}
```

### 3. Gráfico de Pizza com Distribuição de Médicos
```typescript
import { useDoctorDistribution } from '@/hooks/useStats'
import { PieChart, Pie, Cell } from 'recharts'

function DoctorDistributionChart() {
  const { data: distribution, isLoading, error } = useDoctorDistribution()

  const chartData = distribution ? [
    { name: "1 Hospital", value: distribution.com_um_hospital, color: "#3b82f6" },
    { name: "2 Hospitais", value: distribution.com_dois_hospitais, color: "#10b981" },
    { name: "3 Hospitais", value: distribution.com_tres_hospitais, color: "#f59e0b" },
    { name: "Sem Hospitais", value: distribution.sem_hospitais, color: "#ef4444" },
  ] : []

  return (
    <Card>
      <CardHeader>
        <CardTitle>Distribuição de Médicos</CardTitle>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <p>Carregando...</p>
        ) : error ? (
          <p className="text-red-500">Erro: {error}</p>
        ) : (
          <PieChart width={400} height={300}>
            <Pie
              data={chartData}
              dataKey="value"
              nameKey="name"
              cx="50%"
              cy="50%"
              outerRadius={80}
            >
              {chartData.map((entry, index) => (
                <Cell key={`cell-${index}`} fill={entry.color} />
              ))}
            </Pie>
          </PieChart>
        )}
      </CardContent>
    </Card>
  )
}
```

### 4. Filtro por Especialidade
```typescript
import { useHospitalsBySpecialty } from '@/hooks/useStats'

function SpecialtyFilter() {
  const [selectedSpecialty, setSelectedSpecialty] = useState("")
  const { data: hospitals, refetch } = useHospitalsBySpecialty()

  const handleSpecialtyChange = (specialty: string) => {
    setSelectedSpecialty(specialty)
    refetch(specialty || undefined)
  }

  return (
    <div>
      <Select value={selectedSpecialty} onValueChange={handleSpecialtyChange}>
        <SelectTrigger>
          <SelectValue placeholder="Selecione uma especialidade" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="">Todas</SelectItem>
          <SelectItem value="Cardiologia">Cardiologia</SelectItem>
          <SelectItem value="Neurologia">Neurologia</SelectItem>
          {/* Mais opções */}
        </SelectContent>
      </Select>
      
      {/* Mostrar resultados */}
      <div className="mt-4">
        {hospitals?.map((hospital, index) => (
          <div key={index}>
            {'especialidade' in hospital ? (
              <p>{hospital.especialidade}: {hospital.count} hospitais</p>
            ) : (
              <p>{hospital.nome} - {hospital.especialidades}</p>
            )}
          </div>
        ))}
      </div>
    </div>
  )
}
```

## Tratamento de Estados

Todos os hooks retornam os seguintes estados:

```typescript
const { 
  data,        // Dados retornados pela API
  isLoading,   // Estado de carregamento
  error,       // Mensagem de erro (se houver)
  refetch      // Função para recarregar dados
} = useStats()
```

### Padrão de Tratamento:
```typescript
if (isLoading) {
  return <div>Carregando...</div>
}

if (error) {
  return <div className="text-red-500">Erro: {error}</div>
}

if (!data) {
  return <div>Nenhum dado disponível</div>
}

// Renderizar dados
return <div>{/* Seus dados aqui */}</div>
```

## Dicas de Performance

1. **Use os hooks apenas onde necessário**
2. **Implemente debounce para filtros de busca**
3. **Use React.memo para componentes que não precisam re-renderizar**
4. **Implemente cache se necessário**

## Exemplo de Debounce para Busca
```typescript
import { useDebounce } from 'use-debounce'

function SearchComponent() {
  const [searchTerm, setSearchTerm] = useState("")
  const [debouncedSearchTerm] = useDebounce(searchTerm, 500)
  
  const { data, refetch } = useHospitalsBySpecialty()
  
  useEffect(() => {
    if (debouncedSearchTerm) {
      refetch(debouncedSearchTerm)
    }
  }, [debouncedSearchTerm, refetch])

  return (
    <Input
      value={searchTerm}
      onChange={(e) => setSearchTerm(e.target.value)}
      placeholder="Buscar especialidade..."
    />
  )
}
```

## Conclusão

Com esses exemplos, você pode integrar facilmente os dados dos endpoints de estatísticas no seu dashboard, criando uma experiência rica e interativa para os usuários.