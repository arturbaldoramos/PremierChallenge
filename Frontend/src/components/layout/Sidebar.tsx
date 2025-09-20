import {
  Home,
  Settings,
  Building2,
  FolderOpen,
  Stethoscope,
  MapPin,
  User,
  BookOpen,
  Building
} from "lucide-react"
import { Button } from "@/components/ui/button"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Link, useLocation } from "react-router-dom"

export function Sidebar() {
  const location = useLocation()
  
  const navigationItems = [
    { path: '/dashboard', label: 'Dashboard', icon: Home },
    { path: '/hospitals', label: 'Hospitais', icon: Building2 },
    { path: '/doctors', label: 'Médicos', icon: Stethoscope },
    { path: '/patients', label: 'Pacientes', icon: User },
    { path: '/states', label: 'Estados', icon: MapPin },
    { path: '/municipalities', label: 'Municípios', icon: Building },
    { path: '/cid', label: 'CID', icon: BookOpen },
    { path: '/file-manager', label: 'Gerenciador de Arquivos', icon: FolderOpen },
  ]

  return (
    <div className="flex h-full w-64 flex-col border-r bg-background">
      <div className="flex h-14 items-center border-b px-4">
        <Link to="/dashboard" className="text-lg font-semibold hover:text-primary transition-colors">
          Dashboard
        </Link>
      </div>
      
      <div className="flex-1 overflow-auto py-2">
        <nav className="grid items-start px-2 text-sm font-medium lg:px-4">
          {navigationItems.map((item) => (
            <Button 
              key={item.path}
              variant={location.pathname === item.path ? "secondary" : "ghost"} 
              className="justify-start"
              asChild
            >
              <Link to={item.path}>
                <item.icon className="mr-2 h-4 w-4" />
                {item.label}
              </Link>
            </Button>
          ))}
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
