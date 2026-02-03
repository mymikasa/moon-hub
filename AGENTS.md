# Agent 指南 - moon-hub

## 项目概述
采用 React + TypeScript 前端和 Go + Gin 后端的全栈应用，遵循 Clean Architecture 架构。

## API 接口文档
所有后端路由的详细信息请查看 [API_ROUTES.md](./API_ROUTES.md)，包含：
- 完整的路由列表（方法、路径、功能）
- 请求/响应格式示例
- 认证方式说明
- 错误码定义
- 前后端对应关系

**重要**: 添加或修改后端路由时，必须同步更新 API_ROUTES.md 文档。

## 构建与测试命令

### 前端 (React/TypeScript/Vite)
位置：`frontend/`
```bash
cd frontend
npm run dev          # 启动开发服务器
npm run build        # 生产构建（运行 tsc 检查）
npm run lint         # 运行 ESLint
npm run preview      # 预览生产构建
```
运行单个测试（需先添加测试框架）：`npm test -- --testNamePattern="testName"`

### 后端 (Go/Gin)
位置：`backend/moon/`
```bash
cd backend/moon
go build .           # 构建应用
go run .             # 直接运行
go test ./...        # 运行所有测试
go test ./internal/service -run TestUserService_Signup  # 运行指定测试
go test -v ./internal/service -run TestUserService_Signup  # 详细输出指定测试
```

生成 mock：`go generate ./...`

## 代码风格指南

### Go（后端）

**架构层次：**
- `domain/` - 业务实体和核心逻辑（无外部依赖）
- `service/` - 业务逻辑层，使用 repositories
- `repository/` - 数据访问抽象，支持缓存
- `repository/dao/` - 具体数据访问实现（GORM）
- `web/` - HTTP 处理器（Gin 框架）
- `ioc/` - 依赖初始化（数据库、Redis、日志）
- `pkg/` - 共享工具（ginx、logger）

**命名规范：**
- 接口：`UserService`, `UserRepository`（不使用 I 前缀）
- 实现：`userService`, `userRepository`（私有）
- 构造函数：`NewUserService(repo UserRepository) UserService`
- 错误：`ErrDuplicateEmail`, `ErrUserNotFound`（包级别变量）
- 文件：`user.go` 用于包代码，`user_test.go` 用于测试

**导入组织：**
1. 标准库
2. 内部包
3. 第三方包（按字母顺序）

**错误处理：**
- 将服务层错误定义为包级别变量
- 在处理器中使用 switch-case 匹配错误
- 将服务错误映射到 HTTP 错误码（`internal/errs/`）
- 使用 zap 日志记录错误：`L.Error("message", logger.Error(err))`

**Repository 模式：**
- 先定义接口，再用 GORM 实现
- 使用 `go:generate` 生成 mock
- 使用转换函数分离 domain 实体和 DAO 实体
- 支持缓存层，使用 `CachedUserRepository` 模式

**Web 层：**
- 使用 `ginx.WrapBody[Req]()` 处理请求
- 在 `*_vo.go` 文件中定义 VOs（值对象）
- 使用结构体标签：`json:"email" binding:"required"`
- 返回 `ginx.Result` 结构体，包含 Code、Msg、Data 字段

**测试：**
- 使用表驱动测试：`struct{ name string, ... }`
- 使用 testify/assert 进行断言
- 使用 gomock mock 接口
- 测试成功和错误两种路径

### TypeScript/React（前端）

**导入：**
- 使用路径别名：`@/components/ui/button`, `@/lib/utils`
- 分组：外部库在前，本地导入在后
- 组件优先使用命名导出

**组件风格：**
- 使用 shadcn/ui 组件（来自 `@/components/ui/`）
- 使用 `cn()` 工具合并 className（clsx + tailwind-merge）
- 使用 `cva()` 定义变体（class-variance-authority）
- 组件：使用 forwardRef 并设置 displayName

**TypeScript：**
- 启用严格模式
- 无未使用的局部变量/参数
- 定义 props 接口，继承 React.HTMLAttributes
- 使用泛型类型定义变体属性

**样式：**
- Tailwind CSS，自定义主题在 `tailwind.config.js`
- 设计令牌使用 HSL 格式的 CSS 变量
- 使用语义化颜色令牌：`primary`、`destructive`、`muted` 等

**模式：**
- 表单处理器返回 onSubmit 函数
- 使用 useState 进行状态管理
- 受控组件需正确设置 name/id 属性
- 可访问性：aria-live、autoComplete 属性

## 通用规则

- 除非明确要求，否则不添加注释
- 遵循每个文件中现有的代码模式
- Go 函数第一个参数使用 context.Context
- 保持 domain 实体纯净（无基础设施依赖）
- 前端：字符串使用单引号，JSX 属性使用双引号
- 后端：遵循 Go fmt 标准格式化
- 修改后始终运行 lint/typecheck：
  - 前端：`npm run lint` + `tsc -b`
  - 后端：`go fmt ./...` + `go vet ./...`

## 常见问题与解决方案

### CORS 与 Token 认证问题

#### 问题 1：前端无法读取 JWT Token

**现象：**
- 后端登录成功，响应头中包含 `x-jwt-token`
- 前端 `response.headers.get('x-jwt-token')` 返回 `null`
- 导致认证失败，无法访问需要认证的 API

**根本原因：**
根据 CORS 规范，浏览器默认只允许读取以下"简单响应头"：
- `Cache-Control`
- `Content-Language`
- `Content-Type`
- `Expires`
- `Last-Modified`
- `Pragma`

自定义响应头（如 `x-jwt-token`）必须在服务器端的 `Access-Control-Expose-Headers` 中显式声明，前端 JavaScript 才能读取。

**解决方案：**

1. **后端 CORS 中间件** (`backend/moon/internal/web/middleware/cors.go`)：
```go
ctx.Writer.Header().Set("Access-Control-Expose-Headers", "x-jwt-token, x-refresh-token")
```

2. **前端登录函数** (`frontend/src/lib/api.ts`)：
```typescript
export async function login(email: string, password: string): Promise<void> {
  const response = await fetch(`${API_BASE_URL}/users/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify({ email, password }),
  })

  const token = response.headers.get('x-jwt-token')
  if (token) {
    localStorage.setItem('access_token', token)
  }

  if (!response.ok) {
    throw new Error('登录失败')
  }
}
```

**验证步骤：**
1. 登录成功后，在浏览器控制台检查 Network 标签
2. 查看 `/users/login` 请求的 Headers
3. 确认 `Access-Control-Expose-Headers` 包含 `x-jwt-token`
4. 确认响应头中包含 `x-jwt-token`

---

### API 响应格式问题

#### 问题 2：前后端字段名大小写不匹配

**现象：**
- 后端返回：`{"code": 0, "msg": "success", "data": {...}}`
- 前端 TypeScript 接口：`interface ApiResponse { Code?: number; Msg: string; Data?: T }`
- 前端解析后 `data.Code` 为 `undefined`

**根本原因：**
Go 的 JSON 序列化默认使用结构体字段的小写形式，除非显式指定标签。Go 的命名习惯是首字母大写的公开字段，但 JSON 输出通常使用小写。

**解决方案：**

统一使用小写字段名：

1. **后端** (`backend/moon/pkg/ginx/result.go`)：
```go
type Result struct {
    Code int    `json:"code"`
    Msg  string `json:"msg"`
    Data any    `json:"data"`
}
```

2. **前端** (`frontend/src/lib/api.ts`)：
```typescript
interface ApiResponse<T = unknown> {
  code?: number
  msg: string
  data?: T
}
```

**规则：**
- 所有 API 响应使用统一的小写字段名
- 新建接口时参考现有格式
- 修改后端响应时同步更新前端 TypeScript 接口

---

### API 响应完整性问题

#### 问题 3：后端返回格式不一致

**现象：**
- 某些 Handler 返回 `ginx.Result{Data: resp}`，缺少 `Msg` 字段
- 导致前端解析时可能出现问题

**解决方案：**

所有 Handler 返回完整的 `ginx.Result`：

```go
ctx.JSON(http.StatusOK, ginx.Result{
    Code: 0,
    Msg:  "success",
    Data: resp,
})
```

**规则：**
- 成功响应：`Code: 0, Msg: "success" 或具体消息, Data: 数据`
- 错误响应：`Code: 错误码, Msg: 错误信息`
- 保持响应格式一致性，便于前端处理

---

### JWT 认证流程检查清单

**新增认证相关 API 时，必须检查：**

#### 后端检查清单：
- [ ] Handler 返回完整的 `ginx.Result` 格式
- [ ] 在 CORS 中间件中暴露 `x-jwt-token` 响应头
- [ ] JWT 中间件白名单包含新接口（如需要）
- [ ] API_ROUTES.md 文档已更新

#### 前端检查清单：
- [ ] API 函数正确处理响应头中的 token
- [ ] token 存储到 localStorage
- [ ] TypeScript 接口字段名与后端一致（小写）
- [ ] 请求头正确添加 `Authorization: Bearer <token>`
- [ ] 401 错误有统一处理逻辑

---

### 调试技巧

**后端调试：**
```go
fmt.Printf("LoginJWT: 收到登录请求，邮箱: %s\n", req.Email)
fmt.Printf("LoginJWT: token设置成功\n")
```

**前端调试：**
```typescript
console.log('Token received:', token)
console.log('getUserProfile called, token:', token ? 'exists' : 'missing')
console.log('request /users/profile response:', data)
```

**网络请求调试：**
1. 打开浏览器开发者工具（F12）
2. 切换到 Network 标签
3. 查看请求的完整信息：
   - 请求头（确认 Authorization）
   - 响应头（确认 x-jwt-token）
   - 响应内容（确认格式）
   - 状态码（确认 200 vs 401）

---

### 快速诊断步骤

遇到认证问题时，按以下步骤排查：

1. **检查 CORS 配置**
   - 访问 `http://localhost:8080/health`
   - 查看响应头是否有 `Access-Control-Expose-Headers`

2. **检查登录响应**
   - 登录后查看 Network 标签
   - 确认响应头中有 `x-jwt-token`
   - 确认 localStorage 中存储了 token

3. **检查 API 请求**
   - 访问需要认证的 API
   - 查看请求头是否包含 `Authorization: Bearer <token>`
   - 查看响应状态码和内容

4. **检查控制台错误**
   - 查看是否有 CORS 相关错误
   - 查看是否有 TypeScript 类型错误
   - 查看是否有网络请求失败

---

### 避免错误的最佳实践

1. **统一格式**：所有 API 响应使用统一的小写字段名
2. **显式声明**：CORS 响应头必须显式声明所有自定义头
3. **完整响应**：Handler 返回完整的 `ginx.Result` 结构
4. **同步文档**：修改 API 时同步更新 API_ROUTES.md
5. **类型安全**：前端使用 TypeScript 严格模式，确保类型正确
6. **调试日志**：开发阶段添加适当的日志，便于排查问题
