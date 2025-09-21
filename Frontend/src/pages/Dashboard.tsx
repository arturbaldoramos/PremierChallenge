import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { 
  UserCheck,
  Bed,
  Heart,
  RefreshCw,
  AlertCircle,
  Download,
  Building2
} from "lucide-react"
import { useTotalStats, useStats2 } from "@/hooks/use-api"
import { Cid10Stats } from "@/components/Cid10Stats"
import { HospitaisMaisAcessados } from "@/components/HospitaisMaisAcessados"
import { HospitaisChart } from "@/components/HospitaisChart"
import { Cid10Chart } from "@/components/Cid10Chart"

export default function Dashboard() {
  const { data: totalStats, loading, error, refetch } = useTotalStats();
  const { data: stats2 } = useStats2();

  return (
    <div className="flex-1 space-y-4 p-4 md:p-8 pt-6">
      <div className="flex items-center justify-between space-y-2">
        <h2 className="text-3xl font-bold tracking-tight">Dashboard</h2>
        <div className="flex items-center space-x-2">
          <Button 
            variant="outline" 
            size="sm"
            onClick={refetch}
            disabled={loading}
          >
            <RefreshCw className={`mr-2 h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
            {loading ? 'Atualizando...' : 'Atualizar'}
          </Button>
          <Button variant="outline" size="sm">
            <Download className="mr-2 h-4 w-4" />
            Download
          </Button>
        </div>
      </div>

      {/* Texto indicativo de estados e municípios */}
      {stats2 && (
        <div className="text-center text-sm text-muted-foreground">
          Dados de <span className="font-semibold text-primary">{stats2.total_estados}</span> estados e <span className="font-semibold text-primary">{stats2.total_municipios?.toLocaleString()}</span> municípios
        </div>
      )}
      
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Total de Hospitais
            </CardTitle>
            <Building2 className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {loading ? (
                <div className="flex items-center gap-2">
                  <RefreshCw className="h-4 w-4 animate-spin" />
                  Carregando...
                </div>
              ) : error ? (
                <div className="flex items-center gap-2 text-red-500">
                  <AlertCircle className="h-4 w-4" />
                  Erro
                </div>
              ) : (
                totalStats?.total_hospitais?.toLocaleString() || '0'
              )}
            </div>
            <p className="text-xs text-muted-foreground">
              Hospitais cadastrados
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Total de Médicos
            </CardTitle>
            <UserCheck className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {loading ? (
                <div className="flex items-center gap-2">
                  <RefreshCw className="h-4 w-4 animate-spin" />
                  Carregando...
                </div>
              ) : error ? (
                <div className="flex items-center gap-2 text-red-500">
                  <AlertCircle className="h-4 w-4" />
                  Erro
                </div>
              ) : (
                totalStats?.total_medicos?.toLocaleString() || '0'
              )}
            </div>
            <p className="text-xs text-muted-foreground">
              Médicos cadastrados
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total de Leitos</CardTitle>
            <Bed className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {loading ? (
                <div className="flex items-center gap-2">
                  <RefreshCw className="h-4 w-4 animate-spin" />
                  Carregando...
                </div>
              ) : error ? (
                <div className="flex items-center gap-2 text-red-500">
                  <AlertCircle className="h-4 w-4" />
                  Erro
                </div>
              ) : (
                totalStats?.total_leitos?.toLocaleString() || '0'
              )}
            </div>
            <p className="text-xs text-muted-foreground">
              Leitos disponíveis
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total de Pacientes</CardTitle>
            <Heart className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {loading ? (
                <div className="flex items-center gap-2">
                  <RefreshCw className="h-4 w-4 animate-spin" />
                  Carregando...
                </div>
              ) : error ? (
                <div className="flex items-center gap-2 text-red-500">
                  <AlertCircle className="h-4 w-4" />
                  Erro
                </div>
              ) : (
                totalStats?.total_pacientes?.toLocaleString() || '0'
              )}
            </div>
            <p className="text-xs text-muted-foreground">
              Pacientes cadastrados
            </p>
          </CardContent>
        </Card>
      </div>


      {error && (
        <Card className="border-red-200 bg-red-50">
          <CardContent className="pt-6">
            <div className="flex items-center gap-2 text-red-800">
              <AlertCircle className="h-4 w-4" />
              <span className="font-medium">Erro ao carregar dados</span>
            </div>
            <p className="text-sm text-red-700 mt-1">{error}</p>
            <Button 
              variant="outline" 
              size="sm" 
              onClick={refetch}
              className="mt-2"
            >
              <RefreshCw className="mr-2 h-4 w-4" />
              Tentar Novamente
            </Button>
          </CardContent>
        </Card>
      )}

      {totalStats && !loading && (
        <Card className="border-green-200 bg-green-50">
          <CardContent className="pt-6">
            <div className="flex items-center gap-2 text-green-800">
              <div className="w-2 h-2 bg-green-500 rounded-full"></div>
              <span className="text-sm font-medium">Dados atualizados com sucesso</span>
            </div>
            <p className="text-xs text-green-700 mt-1">
              Última atualização: {new Date().toLocaleString('pt-BR')}
            </p>
          </CardContent>
        </Card>
      )}

      {/* Nova seção com estatísticas adicionais */}
      <div className="space-y-6">
        {/* Primeira linha: Hospitais mais acessados (horizontal) */}
        <div className="w-full">
          <HospitaisMaisAcessados />
        </div>
        
        {/* Segunda linha: Doenças mais comuns */}
        <div className="w-full">
          <Cid10Stats />
        </div>
      </div>

      {/* Seção de gráficos */}
      <div className="grid gap-6 grid-cols-1 xl:grid-cols-2">
        <HospitaisChart />
        <Cid10Chart />
      </div>
      
    </div>
  )
}
