import { useAuth } from '@/contexts/AuthContext'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { LayoutDashboard, User, Clock } from 'lucide-react'

export function DashboardPage() {
  const { user } = useAuth()

  const stats = [
    {
      title: '欢迎回来',
      value: user?.nickname || user?.email || '用户',
      description: '今天是美好的一天',
      icon: User,
    },
    {
      title: '账户状态',
      value: '正常',
      description: '您的账户运行正常',
      icon: LayoutDashboard,
    },
    {
      title: '上次登录',
      value: '刚刚',
      description: '欢迎回到 Moon Hub',
      icon: Clock,
    },
  ]

  return (
    <div className="space-y-6">
      <div>
        <h1 className="font-display text-3xl">仪表板</h1>
        <p className="text-sm text-muted-foreground">欢迎使用 Moon Hub</p>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {stats.map((stat) => (
          <Card key={stat.title}>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">{stat.title}</CardTitle>
              <stat.icon className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stat.value}</div>
              <p className="text-xs text-muted-foreground">{stat.description}</p>
            </CardContent>
          </Card>
        ))}
      </div>

      <Card>
        <CardHeader>
          <CardTitle>快速开始</CardTitle>
          <CardDescription>探索 Moon Hub 的功能</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <p className="text-sm text-muted-foreground">
            Moon Hub 是一个现代化的 Web 应用程序，采用 React + TypeScript 前端和 Go + Gin 后端。
          </p>
          <div className="grid gap-4 md:grid-cols-2">
            <div className="rounded-lg border p-4">
              <h3 className="font-semibold mb-2">个人主页</h3>
              <p className="text-sm text-muted-foreground">管理您的个人资料和设置</p>
            </div>
            <div className="rounded-lg border p-4">
              <h3 className="font-semibold mb-2">设置</h3>
              <p className="text-sm text-muted-foreground">自定义您的应用体验</p>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
