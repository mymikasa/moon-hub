import { useState } from "react"

import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { login, register } from "@/lib/api"

interface LoginFormData {
  email: string
  password: string
}

interface RegisterFormData {
  email: string
  password: string
  confirmPassword: string
  nickname: string
}

function App() {
  const [activeTab, setActiveTab] = useState("login")
  const [status, setStatus] = useState("")
  const [statusType, setStatusType] = useState<"success" | "error" | "">("")
  const [isLoading, setIsLoading] = useState(false)
  const [loginForm, setLoginForm] = useState<LoginFormData>({ email: "", password: "" })
  const [registerForm, setRegisterForm] = useState<RegisterFormData>({
    email: "",
    password: "",
    confirmPassword: "",
    nickname: "",
  })

  const handleLoginChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setLoginForm({ ...loginForm, [e.target.name]: e.target.value })
  }

  const handleRegisterChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setRegisterForm({ ...registerForm, [e.target.name]: e.target.value })
  }

  const handleLogin = async (event: React.FormEvent) => {
    event.preventDefault()
    setIsLoading(true)
    setStatus("")

    try {
      await login(loginForm.email, loginForm.password)
      setStatusType("success")
      setStatus("登录成功！")
      window.setTimeout(() => setStatus(""), 2400)
    } catch (error) {
      setStatusType("error")
      setStatus(error instanceof Error ? error.message : "登录失败")
      window.setTimeout(() => setStatus(""), 5000)
    } finally {
      setIsLoading(false)
    }
  }

  const handleRegister = async (event: React.FormEvent) => {
    event.preventDefault()
    setIsLoading(true)
    setStatus("")

    if (registerForm.password !== registerForm.confirmPassword) {
      setStatusType("error")
      setStatus("两次输入的密码不一致")
      setIsLoading(false)
      window.setTimeout(() => setStatus(""), 5000)
      return
    }

    try {
      await register(registerForm.email, registerForm.password, registerForm.confirmPassword, registerForm.nickname)
      setStatusType("success")
      setStatus("注册成功，欢迎加入 Moon Hub")
      setActiveTab("login")
      window.setTimeout(() => setStatus(""), 2400)
    } catch (error) {
      setStatusType("error")
      setStatus(error instanceof Error ? error.message : "注册失败")
      window.setTimeout(() => setStatus(""), 5000)
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-background px-6 py-16 text-foreground">
      <main className="mx-auto w-full max-w-md">
        <div className="mb-8 space-y-2 text-center">
          <h1 className="font-display text-3xl">Moon Hub</h1>
          <p className="text-sm text-muted-foreground">登录或注册以继续使用。</p>
        </div>

        <Card>
          <CardHeader className="space-y-1">
            <CardTitle className="text-xl">欢迎回来</CardTitle>
            <CardDescription>使用邮箱完成登录或注册。</CardDescription>
          </CardHeader>
          <CardContent>
            <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
              <TabsList>
                <TabsTrigger value="login">登录</TabsTrigger>
                <TabsTrigger value="register">注册</TabsTrigger>
              </TabsList>

              <TabsContent value="login">
                <form className="space-y-4" onSubmit={handleLogin}>
                  <div className="space-y-2">
                    <Label htmlFor="login-email">邮箱</Label>
                    <Input
                      id="login-email"
                      name="email"
                      type="email"
                      autoComplete="email"
                      placeholder="you@company.com"
                      required
                      value={loginForm.email}
                      onChange={handleLoginChange}
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="login-password">密码</Label>
                    <Input
                      id="login-password"
                      name="password"
                      type="password"
                      autoComplete="current-password"
                      placeholder="请输入密码"
                      required
                      value={loginForm.password}
                      onChange={handleLoginChange}
                    />
                  </div>
                  <Button className="w-full" disabled={isLoading}>
                    {isLoading ? "登录中..." : "登录"}
                  </Button>
                </form>
              </TabsContent>

              <TabsContent value="register">
                <form className="space-y-4" onSubmit={handleRegister}>
                  <div className="space-y-2">
                    <Label htmlFor="register-nickname">姓名</Label>
                    <Input
                      id="register-nickname"
                      name="nickname"
                      type="text"
                      autoComplete="name"
                      placeholder="你的名字"
                      required
                      value={registerForm.nickname}
                      onChange={handleRegisterChange}
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="register-email">邮箱</Label>
                    <Input
                      id="register-email"
                      name="email"
                      type="email"
                      autoComplete="email"
                      placeholder="you@company.com"
                      required
                      value={registerForm.email}
                      onChange={handleRegisterChange}
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="register-password">设置密码</Label>
                    <Input
                      id="register-password"
                      name="password"
                      type="password"
                      autoComplete="new-password"
                      placeholder="至少 8 位，包含字母、数字、特殊字符"
                      required
                      value={registerForm.password}
                      onChange={handleRegisterChange}
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="register-confirm-password">确认密码</Label>
                    <Input
                      id="register-confirm-password"
                      name="confirmPassword"
                      type="password"
                      autoComplete="new-password"
                      placeholder="再次输入密码"
                      required
                      value={registerForm.confirmPassword}
                      onChange={handleRegisterChange}
                    />
                  </div>
                  <Button className="w-full" disabled={isLoading}>
                    {isLoading ? "注册中..." : "注册"}
                  </Button>
                </form>
              </TabsContent>
            </Tabs>

            <div
              className={`mt-4 min-h-[20px] text-xs ${statusType === "error" ? "text-destructive" : statusType === "success" ? "text-primary" : ""}`}
              aria-live="polite"
            >
              {status}
            </div>
          </CardContent>
        </Card>
      </main>
    </div>
  )
}

export default App
