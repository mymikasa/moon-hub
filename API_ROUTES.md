# API 接口文档

本文档描述后端所有 API 路由的详细信息，包括请求/响应格式、错误码和认证方式。

---

## 认证方式

### JWT Token 认证
- **Token 传递方式**: 通过 `Authorization` header 传递
- **Header 格式**: `Authorization: Bearer <token>`
- **Token 存储**: Redis
- **刷新机制**: 使用 `/users/refresh_token` 接口刷新 access_token

### 前端 Token 管理
```typescript
// 获取 Token
localStorage.getItem('access_token')

// 设置 Token
localStorage.setItem('access_token', token)

// 清除 Token
localStorage.removeItem('access_token')
```

---

## 统一响应格式

所有接口返回统一的 JSON 格式：

```json
{
  "code": 0,
  "msg": "success",
  "data": {}
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| code | number | 0 表示成功，非 0 表示错误 |
| msg | string | 响应消息 |
| data | any | 响应数据（可选，成功时返回具体数据） |

---

## 路由列表

### 用户模块

#### 1. 用户注册
- **方法**: `POST`
- **路径**: `/users/signup`
- **认证**: 否

**请求体**:
```json
{
  "email": "user@example.com",
  "password": "Password123!",
  "confirm_password": "Password123!",
  "nickname": "JohnDoe"
}
```

| 字段 | 类型 | 必填 | 验证规则 |
|------|------|------|----------|
| email | string | 是 | 符合邮箱格式正则 |
| password | string | 是 | 至少8位，包含字母、数字、特殊字符 |
| confirm_password | string | 是 | 必须与 password 相同 |
| nickname | string | 是 | - |

**成功响应** (200 OK):
```json
{
  "code": 0,
  "msg": "注册成功"
}
```

**错误响应**:
- 邮箱冲突 (401003):
  ```json
  {
    "code": 401003,
    "msg": "邮箱冲突"
  }
  ```
- 非法输入 (401001):
  ```json
  {
    "code": 401001,
    "msg": "两次输入的密码不相等"
  }
  ```

---

#### 2. 用户登录 (JWT)
- **方法**: `POST`
- **路径**: `/users/login`
- **认证**: 否

**请求体**:
```json
{
  "email": "user@example.com",
  "password": "Password123!"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| email | string | 是 | 用户邮箱 |
| password | string | 是 | 用户密码 |

**成功响应** (200 OK):
```json
{
  "code": 0,
  "msg": "OK"
}
```
注：Token 会通过 Set-Cookie 设置到浏览器，同时前端可通过响应获取 access_token 存储到 localStorage。

**错误响应**:
- 用户名或密码错误:
  ```json
  {
    "code": 0,
    "msg": "用户名或者密码错误"
  }
  ```
- 系统错误:
  ```json
  {
    "code": 5,
    "msg": "系统错误"
  }
  ```

---

#### 3. 用户登出
- **方法**: `POST`
- **路径**: `/users/logout`
- **认证**: 是 (需要有效的 JWT Token)

**请求体**: 无

**请求头**:
```
Authorization: Bearer <access_token>
```

**成功响应** (200 OK):
```json
{
  "code": 0,
  "msg": "退出登录成功"
}
```

**错误响应** (501001):
```json
{
  "code": 501001,
  "msg": "系统错误"
}
```

---

#### 4. 获取用户信息
- **方法**: `GET`
- **路径**: `/users/profile`
- **认证**: 是 (需要有效的 JWT Token)

**请求头**:
```
Authorization: Bearer <access_token>
```

**成功响应** (200 OK):
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1,
    "email": "user@example.com",
    "nickname": "JohnDoe",
    "birthday": 946684800000,
    "about_me": "Hello, I'm John!",
    "phone": "+8613812345678"
  }
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int64 | 用户 ID |
| email | string | 用户邮箱 |
| nickname | string | 用户昵称 |
| birthday | int64 | 生日（Unix 毫秒时间戳） |
| about_me | string | 个人简介 |
| phone | string | 手机号 |

**错误响应** (401 Unauthorized):
- Token 无效或过期

---

#### 5. 更新用户信息
- **方法**: `PUT`
- **路径**: `/users/profile`
- **认证**: 是 (需要有效的 JWT Token)

**请求头**:
```
Authorization: Bearer <access_token>
```

**请求体**:
```json
{
  "nickname": "NewNick",
  "birthday": 946684800000,
  "about_me": "Updated profile",
  "phone": "+8613898765432"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| nickname | string | 否 | 用户昵称 |
| birthday | int64 | 否 | 生日（Unix 毫秒时间戳） |
| about_me | string | 否 | 个人简介 |
| phone | string | 否 | 手机号 |

**成功响应** (200 OK):
```json
{
  "code": 0,
  "msg": "更新成功"
}
```

**错误响应**:
- 系统错误:
  ```json
  {
    "code": 5,
    "msg": "系统错误"
  }
  ```

---

#### 6. 刷新 Token
- **方法**: `GET`
- **路径**: `/users/refresh_token`
- **认证**: 是 (需要 refresh_token)

**请求头**:
```
Authorization: Bearer <refresh_token>
```

**成功响应** (200 OK):
```json
{
  "code": 0,
  "msg": "OK"
}
```
注：新的 access_token 会通过 Set-Cookie 设置到浏览器。

**错误响应** (401 Unauthorized):
- Token 无效或过期

---

### 其他模块

#### 5. 健康检查
- **方法**: `GET`
- **路径**: `/health`
- **认证**: 否

**成功响应** (200 OK):
```json
{
  "status": "ok"
}
```

---

## 错误码定义

### 用户模块错误码

| 错误码 | 说明 | HTTP 状态码 |
|--------|------|-------------|
| 401001 | 用户输入错误 | 200 |
| 401002 | 用户名或密码错误 | 200 |
| 401003 | 邮箱冲突 | 200 |
| 501001 | 用户模块系统错误 | 200 |
| 5 | 系统错误（通用） | 200 |

### 文章模块错误码

| 错误码 | 说明 | HTTP 状态码 |
|--------|------|-------------|
| 402001 | 文章输入错误 | 200 |
| 502001 | 文章模块系统错误 | 200 |

---

## 前后端对应关系

| 前端函数 | 后端路由 | 前端文件位置 | 后端处理器 |
|----------|----------|--------------|------------|
| `register()` | POST /users/signup | frontend/src/lib/api.ts:50 | UserHandler.SignUp |
| `login()` | POST /users/login | frontend/src/lib/api.ts:43 | UserHandler.LoginJWT |
| `logout()` | POST /users/logout | frontend/src/lib/api.ts:62 | UserHandler.LogoutJWT |
| `refreshToken()` | GET /users/refresh_token | frontend/src/lib/api.ts:69 | UserHandler.RefreshToken |
| `getUserProfile()` | GET /users/profile | frontend/src/lib/api.ts | UserHandler.Profile |
| `updateUserProfile()` | PUT /users/profile | frontend/src/lib/api.ts | UserHandler.UpdateProfile |

---

## 维护说明

**重要**: 当添加、修改或删除后端路由时，请同步更新本文档。

### 更新检查清单
- [ ] 在"路由列表"中添加/修改/删除接口
- [ ] 更新"前后端对应关系"表格
- [ ] 如有新错误码，更新"错误码定义"
- [ ] 更新前端 API 函数对应关系

### 后端路由添加流程
1. 在 `backend/moon/internal/web/` 下创建或修改 Handler
2. 在 Handler 的 `RegisterRoutes` 方法中注册路由
3. 在 `backend/moon/internal/web/` 下创建或修改 VOs (`*_vo.go`)
4. 更新本文档
5. 在 `frontend/src/lib/api.ts` 中添加对应的 API 函数
