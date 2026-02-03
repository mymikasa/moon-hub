# Agent 指南 - moon-hub

## 项目概述
采用 React + TypeScript 前端和 Go + Gin 后端的全栈应用，遵循 Clean Architecture 架构。

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
