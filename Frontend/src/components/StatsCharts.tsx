import React from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
  ChartLegend,
  ChartLegendContent,
} from '@/components/ui/chart'
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, PieChart as RechartsPieChart, Pie, Cell } from 'recharts'
import { 
  useHospitalsByState, 
  useDoctorsByState, 
  useHospitalsBySpecialty, 
  useDoctorDistribution 
} from '@/hooks/useStats'

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

export function StatsCharts() {
  const { data: hospitalsByState, isLoading: hospitalsByStateLoading, error: hospitalsByStateError } = useHospitalsByState()
  const { data: doctorsByState, isLoading: doctorsByStateLoading, error: doctorsByStateError } = useDoctorsByState()
  const { data: hospitalsBySpecialty, isLoading: hospitalsBySpecialtyLoading, error: hospitalsBySpecialtyError } = useHospitalsBySpecialty()
  const { data: doctorDistribution, isLoading: doctorDistributionLoading, error: doctorDistributionError } = useDoctorDistribution()

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
    <div className="space-y-6">
      {/* Gráfico de Hospitais por Estado */}
      <Card>
        <CardHeader>
          <CardTitle>Hospitais por Estado (Top 10)</CardTitle>
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

      {/* Gráfico de Médicos por Estado */}
      <Card>
        <CardHeader>
          <CardTitle>Médicos por Estado (Top 10)</CardTitle>
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

      {/* Gráfico de Pizza - Hospitais por Especialidade */}
      <Card>
        <CardHeader>
          <CardTitle>Hospitais por Especialidade (Top 8)</CardTitle>
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

      {/* Gráfico de Pizza - Distribuição de Médicos */}
      <Card>
        <CardHeader>
          <CardTitle>Distribuição de Médicos por Hospitais</CardTitle>
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
  )
}

export default StatsCharts