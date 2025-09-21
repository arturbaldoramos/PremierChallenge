import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { useHospitaisMaisAcessados } from "@/hooks/use-api"
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, Legend } from 'recharts'
import { RefreshCw, AlertCircle, Building2 } from "lucide-react"

export function HospitaisChart() {
  const { data, loading, error } = useHospitaisMaisAcessados();

  if (loading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Building2 className="h-5 w-5" />
            Gráfico - Hospitais Mais Acessados
          </CardTitle>
          <CardDescription>Top 10 hospitais com mais pacientes</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-center py-8">
            <RefreshCw className="h-6 w-6 animate-spin text-muted-foreground" />
            <span className="ml-2 text-muted-foreground">Carregando dados...</span>
          </div>
        </CardContent>
      </Card>
    );
  }

  if (error) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Building2 className="h-5 w-5" />
            Gráfico - Hospitais Mais Acessados
          </CardTitle>
          <CardDescription>Top 10 hospitais com mais pacientes</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-center py-8">
            <AlertCircle className="h-6 w-6 text-destructive" />
            <span className="ml-2 text-destructive">Erro ao carregar dados</span>
          </div>
        </CardContent>
      </Card>
    );
  }

  if (!data || data.length === 0) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Building2 className="h-5 w-5" />
            Gráfico - Hospitais Mais Acessados
          </CardTitle>
          <CardDescription>Top 10 hospitais com mais pacientes</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="text-center py-8">
            <p className="text-muted-foreground">Nenhum dado disponível</p>
          </div>
        </CardContent>
      </Card>
    );
  }

  // Preparar dados para o gráfico - top 10
  const chartData = data.slice(0, 10).map((item, index) => ({
    nome: item.hospital_nome.length > 20 
      ? item.hospital_nome.substring(0, 20) + '...' 
      : item.hospital_nome,
    nomeCompleto: item.hospital_nome,
    pacientes: item.count,
    rank: index + 1,
    especialidades: item.especialidades?.split(',').slice(0, 2).join(', ') || 'N/A'
  }));

  const CustomTooltip = ({ active, payload }: any) => {
    if (active && payload && payload.length) {
      const data = payload[0].payload;
      return (
        <div className="bg-background border border-border rounded-lg p-3 shadow-lg">
          <p className="font-semibold text-sm mb-1">{data.nomeCompleto}</p>
          <p className="text-xs text-muted-foreground mb-1">
            Rank: #{data.rank}
          </p>
          <p className="text-xs text-muted-foreground mb-1">
            Especialidades: {data.especialidades}
          </p>
          <p className="font-medium text-primary">
            {data.pacientes.toLocaleString()} pacientes
          </p>
        </div>
      );
    }
    return null;
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Building2 className="h-5 w-5" />
          Gráfico - Hospitais Mais Acessados
        </CardTitle>
        <CardDescription>Top 10 hospitais com mais pacientes atendidos</CardDescription>
      </CardHeader>
      <CardContent>
        <div className="h-[400px] w-full">
          <ResponsiveContainer width="100%" height="100%">
            <BarChart
              data={chartData}
              margin={{
                top: 20,
                right: 30,
                left: 20,
                bottom: 60,
              }}
            >
              <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
              <XAxis 
                dataKey="nome" 
                angle={-45}
                textAnchor="end"
                height={80}
                fontSize={12}
                stroke="#6b7280"
              />
              <YAxis 
                tickFormatter={(value) => value.toLocaleString()}
                fontSize={12}
                stroke="#6b7280"
              />
              <Tooltip content={<CustomTooltip />} />
              <Legend />
              <Bar 
                dataKey="pacientes" 
                fill="#3b82f6" 
                name="Pacientes"
                radius={[4, 4, 0, 0]}
              />
            </BarChart>
          </ResponsiveContainer>
        </div>
        
        <div className="mt-4 text-center">
          <div className="text-sm text-muted-foreground">
            Total de {data.length} hospitais • Mostrando top 10
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
