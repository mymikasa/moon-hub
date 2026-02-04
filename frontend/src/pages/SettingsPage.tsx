import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Separator } from '@/components/ui/separator'
import { User, Bell, Shield, Palette } from 'lucide-react'

export function SettingsPage() {
  const settingsSections = [
    {
      title: '账户设置',
      description: '管理您的账户信息和偏好',
      icon: User,
      items: [
        { name: '个人资料', description: '更新您的个人信息' },
        { name: '账户安全', description: '密码和身份验证' },
        { name: '隐私设置', description: '控制您的数据隐私' },
      ],
    },
    {
      title: '通知设置',
      description: '选择您希望接收的通知',
      icon: Bell,
      items: [
        { name: '电子邮件通知', description: '通过电子邮件接收更新' },
        { name: '推送通知', description: '启用浏览器推送通知' },
        { name: '短信通知', description: '通过短信接收重要提醒' },
      ],
    },
    {
      title: '外观',
      description: '自定义应用程序的外观',
      icon: Palette,
      items: [
        { name: '主题', description: '在浅色和深色模式之间切换' },
        { name: '语言', description: '选择您的首选语言' },
        { name: '字体大小', description: '调整文本大小' },
      ],
    },
    {
      title: '安全',
      description: '保护您的账户安全',
      icon: Shield,
      items: [
        { name: '两步验证', description: '添加额外的安全层' },
        { name: '登录活动', description: '查看最近的登录历史' },
        { name: '已连接的设备', description: '管理已登录的设备' },
      ],
    },
  ]

  return (
    <div className="space-y-6">
      <div>
        <h1 className="font-display text-3xl">设置</h1>
        <p className="text-sm text-muted-foreground">管理您的应用设置和偏好</p>
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        {settingsSections.map((section) => (
          <Card key={section.title}>
            <CardHeader>
              <div className="flex items-center gap-2">
                <section.icon className="h-5 w-5 text-muted-foreground" />
                <div>
                  <CardTitle className="text-lg">{section.title}</CardTitle>
                  <CardDescription className="text-sm">{section.description}</CardDescription>
                </div>
              </div>
            </CardHeader>
            <CardContent className="space-y-1">
              {section.items.map((item) => (
                <div key={item.name}>
                  <button className="w-full flex items-center justify-between rounded-md p-3 text-left transition-colors hover:bg-accent hover:text-accent-foreground">
                    <div>
                      <div className="font-medium text-sm">{item.name}</div>
                      <div className="text-xs text-muted-foreground">{item.description}</div>
                    </div>
                  </button>
                  <Separator className="last:hidden" />
                </div>
              ))}
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  )
}
