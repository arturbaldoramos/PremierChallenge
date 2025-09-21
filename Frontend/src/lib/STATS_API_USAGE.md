# Guia de Uso dos Endpoints de Estatísticas

Este documento explica como usar os hooks e funções da API para consumir os endpoints de estatísticas do backend.

## Endpoints Disponíveis

Todos os endpoints seguem o padrão: `GET http://localhost:8080/api/v1/stats/{endpoint}`

### 1. Estatísticas Gerais
- **Endpoint**: `/stats/totals`
- **Hook**: `useStats()`
- **Retorna**: Contadores gerais (hospitais, médicos, estados, municípios)

```typescript
import { useStats } from '@/hooks/useStats'

function MyComponent() {
  const { data, isLoading, error, refetch } = useStats()
  
  if (isLoading) return <div>Carregando...</div>
  if (error) return <div>Erro: {error}</div>
  
  return (
    <div>
      <p>Hospitais: {data?.total_hospitais}</p>
      <p>Médicos: {data?.total_medicos}</p>
      <p>Estados: {data?.total_estados}</p>
      <p>Municípios: {data?.total_municipios}</p>
    </div>
  )
}
```

### 2. Hospitais por Estado
- **Endpoint**: `/stats/hospitais-por-estado`
- **Hook**: `useHospitalsByState()`
- **Retorna**: Array com contagem de hospitais por estado

```typescript
import { useHospitalsByState } from '@/hooks/useStats'

function HospitalsByState() {
  const { data, isLoading, error, refetch } = useHospitalsByState()
  
  return (
    <div>
      {data?.map((item, index) => (
        <div key={index}>
          {item.nome_estado}: {item.count} hospitais
        </div>
      ))}
    </div>
  )
}
```

### 3. Médicos por Estado
- **Endpoint**: `/stats/medicos-por-estado`
- **Hook**: `useDoctorsByState()`
- **Retorna**: Array com contagem de médicos por estado

```typescript
import { useDoctorsByState } from '@/hooks/useStats'

function DoctorsByState() {
  const { data, isLoading, error, refetch } = useDoctorsByState()
  
  return (
    <div>
      {data?.map((item, index) => (
        <div key={index}>
          {item.nome_estado}: {item.count} médicos
        </div>
      ))}
    </div>
  )
}
```

### 4. Hospitais por Especialidade
- **Endpoint**: `/stats/hospitais-por-especialidade`
- **Hook**: `useHospitalsBySpecialty()`
- **Parâmetros**: `specialty` (opcional) - para filtrar por especialidade específica
- **Retorna**: Array com contagem por especialidade OU lista de hospitais filtrados

```typescript
import { useHospitalsBySpecialty } from '@/hooks/useStats'

function HospitalsBySpecialty() {
  const { data, isLoading, error, refetch } = useHospitalsBySpecialty()
  
  // Para buscar todas as especialidades
  const handleLoadAll = () => refetch()
  
  // Para filtrar por especialidade específica
  const handleFilter = (specialty: string) => refetch(specialty)
  
  return (
    <div>
      <button onClick={() => handleFilter('Cardiologia')}>
        Filtrar Cardiologia
      </button>
      {data?.map((item, index) => (
        <div key={index}>
          {/* Verifica se é contagem ou hospital específico */}
          {'especialidade' in item ? (
            <div>{item.especialidade}: {item.count}</div>
          ) : (
            <div>{item.nome} - {item.especialidades}</div>
          )}
        </div>
      ))}
    </div>
  )
}
```

### 5. Hospitais por Município
- **Endpoint**: `/stats/hospitais-por-municipio`
- **Hook**: `useHospitalsByMunicipality()`
- **Parâmetros**: `municipality` (opcional) - para filtrar por município específico
- **Retorna**: Array com contagem por município OU lista de hospitais filtrados

```typescript
import { useHospitalsByMunicipality } from '@/hooks/useStats'

function HospitalsByMunicipality() {
  const { data, isLoading, error, refetch } = useHospitalsByMunicipality()
  
  // Para buscar todos os municípios
  const handleLoadAll = () => refetch()
  
  // Para filtrar por município específico
  const handleFilter = (municipality: string) => refetch(municipality)
  
  return (
    <div>
      <button onClick={() => handleFilter('São Paulo')}>
        Filtrar São Paulo
      </button>
      {data?.map((item, index) => (
        <div key={index}>
          {/* Verifica se é contagem ou hospital específico */}
          {'municipio' in item ? (
            <div>{item.municipio}: {item.count}</div>
          ) : (
            <div>{item.nome} - {item.bairro}</div>
          )}
        </div>
      ))}
    </div>
  )
}
```

### 6. Distribuição de Médicos
- **Endpoint**: `/stats/medicos-distribuicao`
- **Hook**: `useDoctorDistribution()`
- **Retorna**: Objeto com distribuição de médicos por quantidade de hospitais

```typescript
import { useDoctorDistribution } from '@/hooks/useStats'

function DoctorDistribution() {
  const { data, isLoading, error, refetch } = useDoctorDistribution()
  
  return (
    <div>
      {data && (
        <>
          <p>Com 1 hospital: {data.com_um_hospital}</p>
          <p>Com 2 hospitais: {data.com_dois_hospitais}</p>
          <p>Com 3 hospitais: {data.com_tres_hospitais}</p>
          <p>Sem hospitais: {data.sem_hospitais}</p>
        </>
      )}
    </div>
  )
}
```

### 7. Lista de Especialidades
- **Endpoint**: `/stats/especialidades`
- **Hook**: `useSpecialties()`
- **Retorna**: Array com todas as especialidades disponíveis

```typescript
import { useSpecialties } from '@/hooks/useStats'

function SpecialtiesList() {
  const { data, isLoading, error, refetch } = useSpecialties()
  
  return (
    <div>
      {data?.map((specialty, index) => (
        <span key={index} className="badge">
          {specialty}
        </span>
      ))}
    </div>
  )
}
```

## Uso Direto da API (sem hooks)

Se preferir usar a API diretamente sem hooks:

```typescript
import { api } from '@/lib/api'

// Estatísticas gerais
const totals = await api.stats.getTotals()

// Hospitais por estado
const hospitalsByState = await api.stats.getHospitalsByState()

// Médicos por estado
const doctorsByState = await api.stats.getDoctorsByState()

// Hospitais por especialidade (todas)
const hospitalsBySpecialty = await api.stats.getHospitalsBySpecialty()

// Hospitais por especialidade (filtrado)
const cardiologyHospitals = await api.stats.getHospitalsBySpecialty('Cardiologia')

// Hospitais por município (todos)
const hospitalsByMunicipality = await api.stats.getHospitalsByMunicipality()

// Hospitais por município (filtrado)
const spHospitals = await api.stats.getHospitalsByMunicipality('São Paulo')

// Distribuição de médicos
const doctorDistribution = await api.stats.getDoctorDistribution()

// Lista de especialidades
const specialties = await api.stats.getSpecialties()
```

## Tratamento de Erros

Todos os hooks retornam um objeto com:
- `data`: Os dados retornados pela API
- `isLoading`: Estado de carregamento
- `error`: Mensagem de erro (se houver)
- `refetch`: Função para recarregar os dados

```typescript
const { data, isLoading, error, refetch } = useStats()

if (isLoading) {
  return <div>Carregando dados...</div>
}

if (error) {
  return (
    <div>
      <p>Erro: {error}</p>
      <button onClick={refetch}>Tentar novamente</button>
    </div>
  )
}

// Renderizar dados
return <div>{/* Seus dados aqui */}</div>
```

## Exemplo Completo

Veja o arquivo `StatsExample.tsx` para um exemplo completo de como usar todos os endpoints em um componente React.