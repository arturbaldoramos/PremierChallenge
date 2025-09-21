import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { useCid10MaisComuns } from "@/hooks/use-api"
import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip, Legend } from 'recharts'
import { RefreshCw, AlertCircle, Activity } from "lucide-react"

export function Cid10Chart() {
  const { data, loading, error } = useCid10MaisComuns();

  // Cores para o gráfico de pizza
  const COLORS = ['#3b82f6', '#ef4444', '#10b981', '#f59e0b', '#8b5cf6'];

  if (loading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Activity className="h-5 w-5" />
            Gráfico - Doenças Mais Comuns
          </CardTitle>
          <CardDescription>Distribuição das top 5 doenças (CID10)</CardDescription>
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
            <Activity className="h-5 w-5" />
            Gráfico - Doenças Mais Comuns
          </CardTitle>
          <CardDescription>Distribuição das top 5 doenças (CID10)</CardDescription>
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
            <Activity className="h-5 w-5" />
            Gráfico - Doenças Mais Comuns
          </CardTitle>
          <CardDescription>Distribuição das top 5 doenças (CID10)</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="text-center py-8">
            <p className="text-muted-foreground">Nenhum dado disponível</p>
          </div>
        </CardContent>
      </Card>
    );
  }

  // Preparar dados para o gráfico de pizza
  const chartData = data.map((item, index) => ({
    nome: item.doenca.length > 25 
      ? item.doenca.substring(0, 25) + '...' 
      : item.doenca,
    nomeCompleto: item.doenca,
    cid10: item.cid10,
    count: item.count,
    percentual: 0, // Será calculado abaixo
    fill: COLORS[index % COLORS.length]
  }));

  // Calcular percentuais
  const total = chartData.reduce((sum, item) => sum + item.count, 0);
  chartData.forEach(item => {
    item.percentual = Math.round((item.count / total) * 100);
  });

  const CustomTooltip = ({ active, payload }: any) => {
    if (active && payload && payload.length) {
      const data = payload[0].payload;
      return (
        <div className="bg-background border border-border rounded-lg p-3 shadow-lg">
          <p className="font-semibold text-sm mb-1">{data.nomeCompleto}</p>
          <p className="text-xs text-muted-foreground mb-1">
            CID10: {data.cid10}
          </p>
          <p className="font-medium text-primary mb-1">
            {data.count.toLocaleString()} pacientes
          </p>
          <p className="text-xs text-muted-foreground">
            {data.percentual}% do total
          </p>
        </div>
      );
    }
    return null;
  };

  const renderLabel = (entry: any) => {
    return `${entry.cid10}`;
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Activity className="h-5 w-5" />
          Gráfico - Doenças Mais Comuns
        </CardTitle>
        <CardDescription>Distribuição das top 5 doenças (CID10)</CardDescription>
      </CardHeader>
      <CardContent>
        <div className="h-[400px] w-full">
          <ResponsiveContainer width="100%" height="100%">
            <PieChart>
              <Pie
                data={chartData}
                cx="50%"
                cy="50%"
                labelLine={false}
                label={renderLabel}
                outerRadius={120}
                fill="#8884d8"
                dataKey="count"
              >
                {chartData.map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={entry.fill} />
                ))}
              </Pie>
              <Tooltip content={<CustomTooltip />} />
              <Legend 
                verticalAlign="bottom" 
                height={36}
                formatter={(_value, entry: any) => (
                  <span style={{ color: entry.color, fontSize: '12px' }}>
                    {entry.payload.cid10} - {entry.payload.nome}
                  </span>
                )}
              />
            </PieChart>
          </ResponsiveContainer>
        </div>
        
        <div className="mt-4 grid grid-cols-2 gap-2 text-sm">
          {chartData.map((item, index) => (
            <div key={index} className="flex items-center gap-2">
              <div 
                className="w-3 h-3 rounded-full" 
                style={{ backgroundColor: item.fill }}
              ></div>
              <span className="text-muted-foreground">
                {item.cid10}: {item.percentual}%
              </span>
            </div>
          ))}
        </div>

        <div className="mt-4 text-center">
          <div className="text-sm text-muted-foreground">
            Total: {total.toLocaleString()} pacientes
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
