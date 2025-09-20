import { 
  Home, 
  BarChart3, 
  Users, 
  Settings, 
  Bell,
  Search,
  CreditCard,
  Activity,
  TrendingUp
} from "lucide-react"
import { Button } from "@/components/ui/button"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Badge } from "@/components/ui/badge"

export function Sidebar() {
  return (
    <div className="flex h-full w-64 flex-col border-r bg-background">
      <div className="flex h-14 items-center border-b px-4">
        <h1 className="text-lg font-semibold">Dashboard</h1>
      </div>
      
      <div className="flex-1 overflow-auto py-2">
        <nav className="grid items-start px-2 text-sm font-medium lg:px-4">
          <Button variant="ghost" className="justify-start">
            <Home className="mr-2 h-4 w-4" />
            Dashboard
          </Button>
          
          <Button variant="ghost" className="justify-start">
            <BarChart3 className="mr-2 h-4 w-4" />
            Analytics
          </Button>
          
          <Button variant="ghost" className="justify-start">
            <Users className="mr-2 h-4 w-4" />
            Users
          </Button>
          
          <Button variant="ghost" className="justify-start">
            <CreditCard className="mr-2 h-4 w-4" />
            Orders
          </Button>
          
          <Button variant="ghost" className="justify-start">
            <Activity className="mr-2 h-4 w-4" />
            Activity
          </Button>
          
          <Button variant="ghost" className="justify-start">
            <TrendingUp className="mr-2 h-4 w-4" />
            Reports
          </Button>
        </nav>
      </div>
      
      <div className="border-t p-4">
        <div className="flex items-center space-x-2">
          <Avatar className="h-8 w-8">
            <AvatarImage src="/avatars/user.png" alt="User" />
            <AvatarFallback>U</AvatarFallback>
          </Avatar>
          <div className="flex-1 space-y-1">
            <p className="text-sm font-medium leading-none">Usuário</p>
            <p className="text-xs text-muted-foreground">usuario@email.com</p>
          </div>
          <Button variant="ghost" size="sm">
            <Settings className="h-4 w-4" />
          </Button>
        </div>
      </div>
    </div>
  )
}
