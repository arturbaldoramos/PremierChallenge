import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { useNavigate } from "react-router-dom"

export default function Home() {
  const navigate = useNavigate()

  return (
    <div className="min-h-screen bg-background flex items-center justify-center p-4">
      <Card className="w-full max-w-md">
        <CardHeader className="text-center">
          <CardTitle className="text-2xl font-bold">Bem-vindo!</CardTitle>
          <CardDescription>
            Este é um exemplo de aplicação com dashboard usando Shadcn UI
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <Button 
            onClick={() => navigate('/dashboard')} 
            className="w-full"
          >
            Acessar Dashboard
          </Button>
          <p className="text-sm text-muted-foreground text-center">
            Clique no botão acima para navegar para o dashboard com lazy loading
          </p>
        </CardContent>
      </Card>
    </div>
  )
}
