# CLAUDE.md

本文件为 Claude Code (claude.ai/code) 提供在此代码库中工作的指导。

## 项目概述

**moon-hub** 是一个全栈 Web 应用，前端采用 React + TypeScript，后端采用 Go + Gin。遵循 Clean Architecture 架构原则，使用 JWT 认证。

## 命令

### 前端 (React/TypeScript/Vite)
位置: `frontend/`
```bash
cd frontend
npm run dev          # 启动开发服务器
npm run build        # 生产构建（运行 tsc 检查）
npm run lint         # 运行 ESLint
npm run preview      # 预览生产构建
```

### 后端 (Go/Gin)
位置: `backend/moon/`
```bash
cd backend/moon
go build .           # 构建应用
go run .             # 直接运行
go test ./...        # 运行所有测试
go test ./internal/service -run TestUserService_Signup  # 运行指定测试
go test -v ./internal/service -run TestUserService_Signup  # 详细输出
go fmt ./...         # 格式化代码
go vet ./...         # 运行 vet
```

## 架构

### 后端分层结构 (`backend/moon/internal/`)
- `domain/` - 业务实体，无外部依赖
- `service/` - 业务逻辑层，使用 repositories
- `repository/` - 数据访问抽象，支持缓存
- `repository/dao/` - GORM 数据访问实现
- `web/` - HTTP 处理器（Gin 框架）
- `web/middleware/` - CORS、日志、JWT 认证中间件
- `errs/` - 错误码定义
- `ioc/` - 依赖注入（数据库、Redis、日志）
- `pkg/` - 共享工具（ginx、logger）

### 前端结构 (`frontend/src/`)
- `components/ui/` - shadcn/ui 组件（button、card、input、label、tabs）
- `components/` - 应用特定组件（ProtectedRoute、App）
- `contexts/` - React Context（AuthContext）
- `lib/` - API 客户端和工具函数
- `pages/` - 页面组件（HomePage、LoginPage）
- 路径别名: `@/` 映射到 `src/`

### API 文档
所有后端路由的详细文档请参考 [API_ROUTES.md](./API_ROUTES.md)。添加或修改路由时，请同步更新该文件。

### 配置
- 后端配置: `backend/moon/config.yaml`（数据库、Redis、服务器设置）
- 前端使用 Vite + React 插件 + Tailwind CSS

## 代码规范

### Go（后端）
- 接口命名: `UserService`、`UserRepository`（不使用 I 前缀）
- 实现命名: `userService`（私有）
- 构造函数: `NewUserService(repo UserRepository) UserService`
- 错误: 包级别变量，如 `ErrDuplicateEmail`、`ErrUserNotFound`
- Web 处理器: 使用 `ginx.WrapBody[Req]()`，返回 `ginx.Result{Code, Msg, Data}`
- 响应格式: `{"code": 0, "msg": "success", "data": {...}}`（小写字段名）
- Repository 模式: 先定义接口，使用 `go:generate` 生成 mock

### TypeScript/React（前端）
- 使用 `@/components/ui/` 下的 shadcn/ui 组件
- 使用 `cn()` 工具函数合并 className（clsx + tailwind-merge）
- 路径别名: `@/components/ui/button`、`@/lib/utils`
- 字符串使用单引号，JSX 属性使用双引号
- 组件使用 forwardRef 并设置 displayName

## 重要模式

### JWT 认证
- Token 通过 `Authorization: Bearer <token>` 请求头传递
- 登录时在 `x-jwt-token` 响应头返回 token
- CORS 中间件必须在 `Access-Control-Expose-Headers` 中暴露 `x-jwt-token, x-refresh-token`
- 前端将 token 存储在 `localStorage.access_token`

### 测试
- 后端: 表驱动测试，使用 testify/assert 断言，gomock 进行 mock
- 测试文件: `user_test.go` 与 `user.go` 放在一起

### 提交前务必运行
- 前端: `npm run lint` 和 `tsc -b`
- 后端: `go fmt ./...` 和 `go vet ./...`
