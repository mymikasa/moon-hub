import { useState } from 'react'
import { useAuth } from '@/contexts/AuthContext'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

export function HomePage() {
  const { user, refreshProfile } = useAuth()
  const [isEditing, setIsEditing] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const [status, setStatus] = useState('')
  const [statusType, setStatusType] = useState<'success' | 'error' | ''>('')

  const [formData, setFormData] = useState({
    nickname: user?.nickname || '',
    about_me: user?.about_me || '',
    phone: user?.phone || '',
  })

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value })
  }

  const handleSave = async (event: React.FormEvent) => {
    event.preventDefault()
    setIsLoading(true)
    setStatus('')

    try {
      await fetch(`${import.meta.env.VITE_API_BASE_URL}/users/profile`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
        },
        credentials: 'include',
        body: JSON.stringify(formData),
      })

      await refreshProfile()
      setStatusType('success')
      setStatus('更新成功')
      setIsEditing(false)
      window.setTimeout(() => setStatus(''), 2400)
    } catch {
      setStatusType('error')
      setStatus('更新失败')
      window.setTimeout(() => setStatus(''), 3000)
    } finally {
      setIsLoading(false)
    }
  }

  const handleCancel = () => {
    setFormData({
      nickname: user?.nickname || '',
      about_me: user?.about_me || '',
      phone: user?.phone || '',
    })
    setIsEditing(false)
    setStatus('')
  }

  if (!user) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <p>加载中...</p>
      </div>
    )
  }

  const birthday = user.birthday ? new Date(user.birthday).toLocaleDateString('zh-CN') : '未设置'

  return (
    <div className="space-y-6">
      <div>
        <h1 className="font-display text-3xl">个人主页</h1>
        <p className="text-sm text-muted-foreground">管理您的个人资料</p>
      </div>

      <Card>
          <CardHeader>
            <CardTitle>个人信息</CardTitle>
            <CardDescription>管理您的个人资料</CardDescription>
          </CardHeader>
          <CardContent>
            <form className="space-y-6" onSubmit={handleSave}>
              <div className="space-y-2">
                <Label htmlFor="email">邮箱</Label>
                <Input id="email" name="email" type="email" value={user.email} disabled />
              </div>

              <div className="space-y-2">
                <Label htmlFor="nickname">昵称</Label>
                <Input
                  id="nickname"
                  name="nickname"
                  type="text"
                  value={isEditing ? formData.nickname : user.nickname}
                  onChange={handleChange}
                  disabled={!isEditing}
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="birthday">生日</Label>
                <Input id="birthday" name="birthday" type="text" value={birthday} disabled />
              </div>

              <div className="space-y-2">
                <Label htmlFor="phone">手机号</Label>
                <Input
                  id="phone"
                  name="phone"
                  type="tel"
                  value={isEditing ? formData.phone : user.phone}
                  onChange={handleChange}
                  disabled={!isEditing}
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="about_me">个人简介</Label>
                <textarea
                  id="about_me"
                  name="about_me"
                  value={isEditing ? formData.about_me : user.about_me}
                  onChange={(e) => setFormData({ ...formData, about_me: e.target.value })}
                  disabled={!isEditing}
                  className="flex min-h-[120px] w-full rounded-md border border-input bg-background px-3 py-2 text-base ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 md:text-sm"
                  placeholder="介绍一下自己..."
                />
              </div>

              <div className="flex gap-2">
                {isEditing ? (
                  <>
                    <Button type="submit" disabled={isLoading}>
                      {isLoading ? '保存中...' : '保存'}
                    </Button>
                    <Button type="button" variant="outline" onClick={handleCancel} disabled={isLoading}>
                      取消
                    </Button>
                  </>
                ) : (
                  <Button type="button" onClick={() => setIsEditing(true)}>
                    编辑
                  </Button>
                )}
              </div>

              {status && (
                <div
                  className={`min-h-[20px] text-sm ${statusType === 'error' ? 'text-destructive' : statusType === 'success' ? 'text-primary' : ''}`}
                  aria-live="polite"
                >
                  {status}
                </div>
              )}
            </form>
          </CardContent>
        </Card>
    </div>
  )
}
