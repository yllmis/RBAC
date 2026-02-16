# RBAC 权限管理系统

一个基于 Go + Gin + MySQL + Redis 的轻量级 RBAC（Role-Based Access Control）后端项目，提供登录认证、权限鉴权、用户查询与角色分配功能，并内置一个可直接使用的前端管理台页面。

## 功能概览

- 登录认证
  - 账号密码登录
  - JWT 签发与解析
- 鉴权中间件
  - Bearer Token 校验
  - 未登录返回 401、无权限返回 403
- RBAC 权限模型
  - 用户 -> 角色 -> 权限
  - 支持按 `perm_code` 鉴权
  - 支持按 `method + api_path` 路由维度鉴权
- 权限缓存
  - 登录后预热用户权限到 Redis
  - 鉴权时优先读缓存，未命中回源 MySQL 并回填
  - 正向/负向权限使用不同 TTL 策略
- 角色分配
  - 分配前参数校验
  - 用户存在性、角色存在性校验
  - 角色冲突校验（当前实现为单用户单角色）
- 前端管理台（单页）
  - 登录页面
  - 顶部导航栏 + 左侧边栏
  - 用户列表页
  - 分配角色页

## 技术栈

- 后端框架：Gin
- ORM：GORM
- 数据库：MySQL
- 缓存：Redis
- 鉴权：JWT
- 日志：Zap + Lumberjack
- 前端：Vue3（CDN）+ Axios（CDN）

## 项目结构

```text
cmd/                    程序入口
internal/
  handler/              HTTP 处理层
  middleware/           鉴权中间件
  model/                数据模型
  repository/           数据访问层（MySQL/Redis）
  router/               路由注册
  service/              业务层
pkg/
  ecode/                统一错误码
  log/                  日志初始化
  utils/                JWT 工具
web/
  login.html            前端管理台页面
```

## RBAC 数据模型

核心表关系：

- `user`：用户
- `role`：角色
- `permission`：权限点（`perm_code`, `method`, `api_path`）
- `user_role`：用户-角色关联
- `role_perm`：角色-权限关联

权限判定链路：

1. 根据 `user_id` 查询 `user_role` 得到角色列表
2. 通过 `role_perm` 关联 `permission`
3. 匹配条件：
   - `perm_code` 命中，或
   - `(method, api_path)` 命中

## 接口说明

### 1) 登录

- `POST /login`
- 请求体：

```json
{
  "account": "admin",
  "password": "123456"
}
```

- 成功响应：

```json
{
  "code": 0,
  "message": "",
  "data": {
    "token": "..."
  }
}
```

### 2) 查询用户列表（需权限）

- `GET /api/users`
- 所需权限：`user:list`
- Header：`Authorization: Bearer <token>`

### 3) 分配角色（需权限）

- `POST /api/user/role`
- 所需权限：`user:role:set`
- Header：`Authorization: Bearer <token>`
- 请求体：

```json
{
  "user_id": 1,
  "role_id": 2
}
```

## 角色分配规则（当前实现）

`SetUserRole` 业务规则如下：

- `user_id <= 0`：返回“请选择一个用户”
- `role_id <= 0`：返回“请选择一个角色”
- 用户不存在：返回“用户不存在”
- 角色不存在：返回“角色不存在”
- 用户已拥有目标角色：返回“用户已拥有该角色”
- 用户已分配其他角色：返回“用户已分配其他角色，请先解绑”
- 否则写入 `user_role`

> 说明：当前策略是“单用户单角色”。若要支持多角色，可在 service 层放开冲突限制。

## 本地运行

### 1) 环境准备

- Go 1.25+
- MySQL 8+
- Redis 6+

### 2) 安装依赖

```bash
go mod tidy
```

### 3) 配置数据库与 Redis

当前仓库版本使用了硬编码连接信息，请根据你的环境修改以下文件：

- `internal/repository/mysql.go`
- `internal/repository/redis.go`

### 4) 启动服务

```bash
go run ./cmd/main.go
```

默认监听：`http://localhost:8080`

- 前端管理台：`http://localhost:8080/`
- 登录接口：`http://localhost:8080/login`

## 前端管理台说明

`web/login.html` 已集成：

- 登录页
- 顶部导航栏
- 左侧边栏
- 用户列表视图
- 分配角色视图

登录成功后会自动携带 token 请求受保护接口。

## 测试

运行全部测试：

```bash
go test ./...
```

说明：

- `internal/middleware` 覆盖鉴权中间件关键路径
- `internal/service/role_test.go` 覆盖角色分配核心业务规则
- `internal/repository/repository_test.go` 为集成测试，需要环境变量控制：

```bash
RBAC_RUN_INTEGRATION=1 go test ./internal/repository -v
```

## 常见问题

### 1) 管理台提示“无权访问”

请检查：

- 当前用户所属角色是否绑定 `user:list` 与 `user:role:set`
- `permission` 表中 `perm_code` 是否准确
- 是否存在旧缓存（`user_perms_{userId}`），必要时删除后重新登录

### 2) 启动失败

通常是 MySQL/Redis 不可达或连接配置不匹配，请先确认：

- MySQL 地址、端口、账号、库名
- Redis 地址、端口

## 后续可扩展方向

- 角色管理与权限管理 CRUD
- 用户角色解绑/批量分配
- 配置外置（env/config）
- 审计日志与操作记录
- 更完整的前端管理页面

