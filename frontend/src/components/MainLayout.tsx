import { useNavigate, useLocation } from 'react-router-dom'
import { useAuth } from '@/contexts/AuthContext'
import {
  LayoutDashboard,
  User,
  Settings,
  LogOut,
} from 'lucide-react'
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarProvider,
  SidebarTrigger,
} from '@/components/ui/sidebar'
import { TooltipProvider } from '@/components/ui/tooltip'

const menuItems = [
  { title: '仪表板', icon: LayoutDashboard, path: '/dashboard' },
  { title: '个人主页', icon: User, path: '/home' },
  { title: '设置', icon: Settings, path: '/settings' },
]

interface MainLayoutProps {
  children: React.ReactNode
}

export function MainLayout({ children }: MainLayoutProps) {
  const navigate = useNavigate()
  const location = useLocation()
  const { user, logout } = useAuth()

  const handleLogout = async () => {
    try {
      await logout()
      navigate('/login')
    } catch {
      console.error('退出登录失败')
    }
  }

  const isActive = (path: string) => location.pathname === path

  return (
    <TooltipProvider>
      <SidebarProvider defaultOpen={true}>
        <div className="flex h-screen w-full">
        <Sidebar collapsible="icon">
          <SidebarHeader className="border-b p-4">
            <div className="flex items-center justify-between gap-2 group-data-[collapsible=icon]:hidden">
              <div className="flex items-center gap-2">
                <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary text-primary-foreground">
                  <span className="font-bold">M</span>
                </div>
                <span className="font-display text-lg font-semibold">
                  Moon Hub
                </span>
              </div>
              <div className="group/trigger relative">
                <SidebarTrigger />
                <div className="absolute left-full ml-2 hidden rounded-md bg-popover px-2 py-1 text-xs text-popover-foreground shadow-md group-data-[collapsible=icon]/trigger:block">
                  折叠侧边栏
                </div>
              </div>
            </div>
            <div className="group-data-[collapsible=icon]:flex hidden items-center justify-center">
              <SidebarTrigger />
            </div>
          </SidebarHeader>
          <SidebarContent>
            <SidebarGroup>
              <SidebarGroupContent>
                <SidebarMenu>
                  {menuItems.map((item) => (
                    <SidebarMenuItem key={item.path}>
                      <SidebarMenuButton
                        isActive={isActive(item.path)}
                        tooltip={item.title}
                        onClick={() => navigate(item.path)}
                      >
                        <item.icon className="h-4 w-4" />
                        <span className="group-data-[collapsible=icon]:hidden">{item.title}</span>
                      </SidebarMenuButton>
                    </SidebarMenuItem>
                  ))}
                </SidebarMenu>
              </SidebarGroupContent>
            </SidebarGroup>
          </SidebarContent>
          <SidebarFooter className="border-t p-2">
            <SidebarMenu>
              <SidebarMenuItem>
                <SidebarMenuButton tooltip="退出登录" onClick={handleLogout}>
                  <LogOut className="h-4 w-4" />
                  <span className="group-data-[collapsible=icon]:hidden">退出登录</span>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarFooter>
        </Sidebar>
        <main className="flex flex-1 flex-col overflow-hidden">
          <header className="flex h-14 items-center justify-between px-4">
            <div className="flex items-center gap-2">
              <span className="text-sm font-medium text-muted-foreground">
                {menuItems.find((item) => isActive(item.path))?.title || 'Moon Hub'}
              </span>
            </div>
            {user && (
              <div className="flex items-center gap-2">
                <span className="text-sm text-muted-foreground">
                  {user.nickname || user.email}
                </span>
              </div>
            )}
          </header>
          <div className="flex-1 overflow-auto p-6">
            {children}
          </div>
        </main>
      </div>
    </SidebarProvider>
    </TooltipProvider>
  )
}
