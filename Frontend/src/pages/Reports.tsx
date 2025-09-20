import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { TrendingUp, Download, Calendar, FileText, BarChart3 } from "lucide-react"

export default function Reports() {
  const reports = [
    {
      id: 1,
      name: "Relatório de Vendas - Janeiro 2024",
      type: "Vendas",
      date: "2024-01-31",
      status: "Concluído",
      size: "2.3 MB"
    },
    {
      id: 2,
      name: "Análise de Usuários - Q4 2023",
      type: "Usuários",
      date: "2023-12-31",
      status: "Concluído",
      size: "1.8 MB"
    },
    {
      id: 3,
      name: "Relatório Financeiro - Dezembro 2023",
      type: "Financeiro",
      date: "2023-12-31",
      status: "Processando",
      size: "3.1 MB"
    },
    {
      id: 4,
      name: "Dashboard Analytics - Novembro 2023",
      type: "Analytics",
      date: "2023-11-30",
      status: "Concluído",
      size: "1.2 MB"
    },
    {
      id: 5,
      name: "Relatório de Atividade - Outubro 2023",
      type: "Atividade",
      date: "2023-10-31",
      status: "Concluído",
      size: "2.7 MB"
    },
  ]

  const getStatusVariant = (status: string) => {
    switch (status) {
      case "Concluído":
        return "default"
      case "Processando":
        return "secondary"
      case "Erro":
        return "destructive"
      default:
        return "outline"
    }
  }

  const getTypeIcon = (type: string) => {
    switch (type) {
      case "Vendas":
        return <TrendingUp className="h-4 w-4" />
      case "Usuários":
        return <BarChart3 className="h-4 w-4" />
      case "Financeiro":
        return <FileText className="h-4 w-4" />
      case "Analytics":
        return <BarChart3 className="h-4 w-4" />
      case "Atividade":
        return <Calendar className="h-4 w-4" />
      default:
        return <FileText className="h-4 w-4" />
    }
  }

  return (
    <div className="flex-1 space-y-4 p-4 md:p-8 pt-6">
      <div className="flex items-center justify-between space-y-2">
        <h2 className="text-3xl font-bold tracking-tight">Relatórios</h2>
        <Button>
          <FileText className="mr-2 h-4 w-4" />
          Gerar Relatório
        </Button>
      </div>
      
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Total de Relatórios
            </CardTitle>
            <FileText className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">156</div>
            <p className="text-xs text-muted-foreground">
              +8 novos este mês
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Downloads Hoje
            </CardTitle>
            <Download className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">23</div>
            <p className="text-xs text-muted-foreground">
              +15% from yesterday
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Em Processamento
            </CardTitle>
            <Calendar className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">3</div>
            <p className="text-xs text-muted-foreground">
              Aguardando conclusão
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Armazenamento
            </CardTitle>
            <FileText className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">2.3 GB</div>
            <p className="text-xs text-muted-foreground">
              de 10 GB utilizados
            </p>
          </CardContent>
        </Card>
      </div>
      
      <Card>
        <CardHeader>
          <CardTitle>Relatórios Disponíveis</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {reports.map((report) => (
              <div key={report.id} className="flex items-center justify-between p-4 border rounded-lg">
                <div className="flex items-center space-x-4">
                  {getTypeIcon(report.type)}
                  <div>
                    <p className="text-sm font-medium">{report.name}</p>
                    <div className="flex items-center space-x-2 mt-1">
                      <Badge variant="outline">{report.type}</Badge>
                      <span className="text-xs text-muted-foreground">{report.date}</span>
                      <span className="text-xs text-muted-foreground">•</span>
                      <span className="text-xs text-muted-foreground">{report.size}</span>
                    </div>
                  </div>
                </div>
                <div className="flex items-center space-x-2">
                  <Badge variant={getStatusVariant(report.status) as any}>
                    {report.status}
                  </Badge>
                  {report.status === "Concluído" && (
                    <Button variant="outline" size="sm">
                      <Download className="mr-2 h-4 w-4" />
                      Download
                    </Button>
                  )}
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
