# Backend Service

Go 后端服务

## 开发环境要求

- Go 1.25+
- PostgreSQL 14+
- Redis (可选)

## 安装依赖

```bash
go mod download
```

## 配置

1. 复制配置文件模板：
```bash
cp config/config.example.yaml config/config.yaml
```

2. 修改 `config/config.yaml` 中的数据库配置

3. 或使用环境变量覆盖配置：
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=onetaste_family
```

## 数据库迁移

### 安装 golang-migrate

**macOS:**
```bash
brew install golang-migrate
```

**Linux:**
```bash
# 下载二进制文件
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/migrate
```

**Windows:**
```powershell
# 使用 Chocolatey
choco install golang-migrate
```

### 执行迁移

使用提供的迁移脚本：

```bash
# 执行所有待执行的迁移
./scripts/migrate.sh up

# 回滚最后一次迁移
./scripts/migrate.sh down

# 查看当前迁移版本
./scripts/migrate.sh version

# 创建新的迁移文件
./scripts/migrate.sh create add_new_table
```

或直接使用 migrate 命令：

```bash
# 执行迁移
migrate -path ./migrations -database "postgres://user:password@localhost:5432/onetaste_family?sslmode=disable" up

# 回滚
migrate -path ./migrations -database "postgres://user:password@localhost:5432/onetaste_family?sslmode=disable" down
```

## 运行

```bash
go run cmd/main.go
```

## 媒体文件上传

- 接口：`POST /api/v1/media/upload`（需 Bearer Token）
- 请求：`multipart/form-data`，字段包括  
  - `file`：要上传的文件，仅支持常见图片（JPG/PNG/WebP/GIF/SVG）  
  - `path`：文件目录前缀，由前端控制，例如：
    - `/user/head/{userId}`：该用户历史头像
    - `/family/head/{familyId}`：家庭头像
    - `/family/{familyId}/dishes/{dishId}`：指定菜式图片
- 后端会在该目录下生成 `ULID + 扩展名` 的文件名并上传至 MinIO，接口返回可直接访问的 URL 与对象路径，便于后续覆盖或清理。

## 构建

```bash
go build -o app cmd/main.go
```

## 项目结构

```
src/backend/
├── cmd/                    # 应用程序入口
│   └── main.go
├── internal/               # 内部代码
│   ├── config/            # 配置管理
│   ├── models/            # 数据模型
│   ├── handlers/          # HTTP处理器
│   ├── services/          # 业务逻辑层
│   ├── repositories/      # 数据访问层
│   ├── middleware/        # 中间件
│   └── utils/             # 工具函数
├── pkg/                   # 可复用的公共包
│   └── database/         # 数据库连接
├── migrations/            # 数据库迁移脚本
├── config/                # 配置文件目录
└── scripts/               # 脚本文件
```

## 数据库表结构

数据库包含以下核心表：

1. **users** - 用户表
2. **families** - 家庭表
3. **family_members** - 家庭成员表
4. **dishes** - 菜式表
5. **ingredients** - 食材表
6. **cooking_steps** - 烹饪步骤表
7. **menus** - 菜单表
8. **menu_dishes** - 菜单菜式关联表
9. **shopping_lists** - 购物清单表
10. **shopping_list_items** - 购物清单项表
11. **health_records** - 身体状况记录表
12. **memberships** - 会员表
13. **ai_usage_logs** - AI调用记录表
14. **payment_orders** - 支付订单表
15. **system_configs** - 系统配置表

详细设计请参考 [数据库设计文档](../../docs/数据库设计.md)

## 开发规范

- 遵循 Go 代码规范
- 使用分层架构：Handler -> Service -> Repository
- 错误处理要明确
- 添加必要的注释

## 相关文档

- [API接口文档](../../docs/接口文档.md)
- [数据库设计文档](../../docs/数据库设计.md)
- [后端目录结构](../../docs/BACKEND_DIRECTORY.md)
