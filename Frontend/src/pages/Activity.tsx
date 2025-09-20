import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Badge } from "@/components/ui/badge"
import { Activity as ActivityIcon, Clock, User, ShoppingCart, Settings } from "lucide-react"

export default function Activity() {
  const activities = [
    {
      id: 1,
      user: "João Silva",
      action: "fez login",
      time: "2 minutos atrás",
      icon: User,
      type: "login"
    },
    {
      id: 2,
      user: "Maria Santos",
      action: "criou um novo pedido",
      time: "5 minutos atrás",
      icon: ShoppingCart,
      type: "order"
    },
    {
      id: 3,
      user: "Pedro Costa",
      action: "atualizou suas configurações",
      time: "10 minutos atrás",
      icon: Settings,
      type: "settings"
    },
    {
      id: 4,
      user: "Ana Oliveira",
      action: "fez logout",
      time: "15 minutos atrás",
      icon: User,
      type: "logout"
    },
    {
      id: 5,
      user: "Carlos Lima",
      action: "visualizou o dashboard",
      time: "20 minutos atrás",
      icon: ActivityIcon,
      type: "view"
    },
  ]

  const getTypeColor = (type: string) => {
    switch (type) {
      case "login":
        return "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300"
      case "order":
        return "bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300"
      case "settings":
        return "bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300"
      case "logout":
        return "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300"
      case "view":
        return "bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-300"
      default:
        return "bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-300"
    }
  }

  return (
    <div className="flex-1 space-y-4 p-4 md:p-8 pt-6">
      <div className="flex items-center justify-between space-y-2">
        <h2 className="text-3xl font-bold tracking-tight">Atividade</h2>
      </div>
      
      <div className="grid gap-4 md:grid-cols-3">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Usuários Online
            </CardTitle>
            <ActivityIcon className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">23</div>
            <p className="text-xs text-muted-foreground">
              +5 from last hour
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Ações Hoje
            </CardTitle>
            <Clock className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">1,456</div>
            <p className="text-xs text-muted-foreground">
              +12% from yesterday
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Picos de Atividade
            </CardTitle>
            <ActivityIcon className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">3</div>
            <p className="text-xs text-muted-foreground">
              Últimas 24h
            </p>
          </CardContent>
        </Card>
      </div>
      
      <Card>
        <CardHeader>
          <CardTitle>Atividade Recente</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {activities.map((activity) => (
              <div key={activity.id} className="flex items-center space-x-4 p-4 border rounded-lg">
                <Avatar className="h-10 w-10">
                  <AvatarImage src={`/avatars/${activity.id}.png`} alt={activity.user} />
                  <AvatarFallback>{activity.user.split(' ').map(n => n[0]).join('')}</AvatarFallback>
                </Avatar>
                <div className="flex-1">
                  <div className="flex items-center space-x-2">
                    <span className="text-sm font-medium">{activity.user}</span>
                    <span className="text-sm text-muted-foreground">{activity.action}</span>
                    <Badge className={getTypeColor(activity.type)}>
                      <activity.icon className="mr-1 h-3 w-3" />
                      {activity.type}
                    </Badge>
                  </div>
                  <p className="text-xs text-muted-foreground">{activity.time}</p>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
