import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { useCid10MaisComuns } from "@/hooks/use-api"
import { Activity, RefreshCw, AlertCircle } from "lucide-react"

export function Cid10Stats() {
  const { data, loading, error } = useCid10MaisComuns();

  if (loading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Activity className="h-5 w-5" />
            Doenças Mais Comuns (CID10)
          </CardTitle>
          <CardDescription>Top 5 doenças mais comuns dos pacientes</CardDescription>
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
            Doenças Mais Comuns (CID10)
          </CardTitle>
          <CardDescription>Top 5 doenças mais comuns dos pacientes</CardDescription>
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
            Doenças Mais Comuns (CID10)
          </CardTitle>
          <CardDescription>Top 5 doenças mais comuns dos pacientes</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="text-center py-8">
            <p className="text-muted-foreground">Nenhum dado disponível</p>
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Activity className="h-5 w-5" />
          Doenças Mais Comuns (CID10)
        </CardTitle>
        <CardDescription>Top 5 doenças mais comuns dos pacientes</CardDescription>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          {data.map((item, index) => (
            <div key={item.cid10} className="flex items-center justify-between p-3 border rounded-lg">
              <div className="flex items-center gap-3">
                <div className="flex items-center justify-center w-8 h-8 rounded-full bg-primary/10 text-primary font-semibold text-sm">
                  {index + 1}
                </div>
                <div>
                  <div className="font-medium">{item.doenca}</div>
                  <Badge variant="secondary" className="text-xs">
                    CID10: {item.cid10}
                  </Badge>
                </div>
              </div>
              <div className="text-right">
                <div className="font-semibold text-lg">{item.count.toLocaleString()}</div>
                <div className="text-sm text-muted-foreground">pacientes</div>
              </div>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}
