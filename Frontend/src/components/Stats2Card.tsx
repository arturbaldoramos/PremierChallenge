import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { useStats2 } from "@/hooks/use-api"
import { MapPin, RefreshCw, AlertCircle } from "lucide-react"

export function Stats2Card() {
  const { data, loading, error } = useStats2();

  if (loading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <MapPin className="h-5 w-5" />
            Localização
          </CardTitle>
          <CardDescription>Estados e municípios cadastrados</CardDescription>
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
            <MapPin className="h-5 w-5" />
            Localização
          </CardTitle>
          <CardDescription>Estados e municípios cadastrados</CardDescription>
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

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <MapPin className="h-5 w-5" />
          Localização
        </CardTitle>
        <CardDescription>Estados e municípios cadastrados</CardDescription>
      </CardHeader>
      <CardContent>
        <div className="grid grid-cols-2 gap-3">
          <div className="text-center p-4 border rounded-lg bg-gradient-to-br from-blue-50 to-blue-100/50 dark:from-blue-950/50 dark:to-blue-900/30">
            <div className="text-3xl font-bold text-blue-600 dark:text-blue-400">
              {data?.total_estados || 0}
            </div>
            <div className="text-sm text-muted-foreground font-medium">Estados</div>
            <div className="text-xs text-muted-foreground mt-1">Unidades Federativas</div>
          </div>
          <div className="text-center p-4 border rounded-lg bg-gradient-to-br from-green-50 to-green-100/50 dark:from-green-950/50 dark:to-green-900/30">
            <div className="text-3xl font-bold text-green-600 dark:text-green-400">
              {data?.total_municipios?.toLocaleString() || 0}
            </div>
            <div className="text-sm text-muted-foreground font-medium">Municípios</div>
            <div className="text-xs text-muted-foreground mt-1">Cidades cadastradas</div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
