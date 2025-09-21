import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { useHospitaisMaisAcessados } from "@/hooks/use-api"
import { Building2, RefreshCw, AlertCircle, Users } from "lucide-react"

export function HospitaisMaisAcessados() {
  const { data, loading, error } = useHospitaisMaisAcessados();

  if (loading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Building2 className="h-5 w-5" />
            Hospitais Mais Acessados
          </CardTitle>
          <CardDescription>Top 20 hospitais mais acessados pelos pacientes</CardDescription>
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
            Hospitais Mais Acessados
          </CardTitle>
          <CardDescription>Top 20 hospitais mais acessados pelos pacientes</CardDescription>
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
            Hospitais Mais Acessados
          </CardTitle>
          <CardDescription>Top 20 hospitais mais acessados pelos pacientes</CardDescription>
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
          <Building2 className="h-5 w-5" />
          Hospitais Mais Acessados
        </CardTitle>
        <CardDescription>Top 20 hospitais mais acessados pelos pacientes</CardDescription>
      </CardHeader>
      <CardContent>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
          {data.slice(0, 12).map((item, index) => (
            <div key={item.hospital_uuid} className="flex flex-col p-4 border rounded-lg hover:bg-muted/50 transition-colors">
              <div className="flex items-center justify-between mb-2">
                <div className="flex items-center justify-center w-6 h-6 rounded-full bg-primary/10 text-primary font-semibold text-xs flex-shrink-0">
                  {index + 1}
                </div>
                <div className="flex items-center gap-1 font-semibold text-sm">
                  <Users className="h-3 w-3 text-primary" />
                  {item.count.toLocaleString()}
                </div>
              </div>
              <div className="flex-1 min-w-0">
                <div className="font-medium text-sm mb-2 line-clamp-2">
                  {item.hospital_nome}
                </div>
                <div className="flex flex-wrap gap-1">
                  {item.especialidades?.split(',').slice(0, 1).map((esp, idx) => (
                    <Badge key={idx} variant="secondary" className="text-xs px-2 py-0.5">
                      {esp.trim()}
                    </Badge>
                  ))}
                  {item.especialidades?.split(',').length > 1 && (
                    <Badge variant="outline" className="text-xs px-2 py-0.5">
                      +{item.especialidades.split(',').length - 1}
                    </Badge>
                  )}
                </div>
              </div>
            </div>
          ))}
        </div>
        {data.length > 12 && (
          <div className="mt-3 text-center text-sm text-muted-foreground">
            Mostrando top 12 de {data.length} hospitais
          </div>
        )}
      </CardContent>
    </Card>
  );
}
