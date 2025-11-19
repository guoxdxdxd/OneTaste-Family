# Backend 目录结构

本文档基于当前仓库 (`src/backend`) 的真实目录生成，列出了所有在仓库中的目录与文件，并为每个条目提供了用途说明，方便快速了解后端 Go 服务的组织方式。

## 目录结构（含注释）

```
src/backend/
├── .DS_Store                          # macOS 访达生成的缓存文件，可忽略
├── .idea/                             # GoLand/IDEA 的工程配置
├── README.md                          # 后端服务的使用说明与启动指引
├── cmd/                               # 应用入口所在目录
│   └── main.go                        # Go 主程序入口，初始化依赖并启动 HTTP 服务
├── config/                            # 配置文件目录
│   ├── config.example.yaml            # 配置示例，列出所有关键配置项
│   └── config.yaml                    # 本地默认配置，供开发环境使用
├── docs/                              # 后端使用的生成型文档
│   └── swagger/                       # Swagger 产物目录
│       ├── README.md                  # Swagger 文档的使用说明
│       ├── docs.go                    # swag 工具生成的 Go 注解文件
│       ├── swagger.json               # Swagger JSON 规格文件
│       └── swagger.yaml               # Swagger YAML 规格文件
├── go.mod                             # Go 模块定义，声明依赖与模块名
├── go.sum                             # Go 依赖的校验和记录
├── internal/                          # 仅供内部使用的业务代码
│   ├── .DS_Store                      # Finder 缓存文件，可忽略
│   ├── config/                        # 配置读取逻辑
│   │   ├── config.go                  # 配置结构体定义
│   │   └── loader.go                  # 读取 YAML/环境变量的加载器
│   ├── handlers/                      # HTTP 控制器层
│   │   ├── README.md                  # Handler 层开发约定
│   │   ├── auth_handler.go            # 认证相关接口（登录、注册等）
│   │   ├── family_handler.go          # 家庭数据的 HTTP 接口
│   │   ├── router.go                  # 路由初始化及依赖注入
│   │   ├── routes.go                  # 路由表与分组定义
│   │   └── user_handler.go            # 用户信息相关接口
│   ├── middleware/                    # HTTP 中间件集合
│   │   └── auth.go                    # JWT 鉴权中间件
│   ├── models/                        # 数据模型定义
│   │   ├── family.go                  # 家庭实体及数据库映射
│   │   └── user.go                    # 用户实体及数据库映射
│   ├── repositories/                  # 数据访问层
│   │   ├── family_repository.go       # 家庭表 CRUD 封装
│   │   └── user_repository.go         # 用户表 CRUD 封装
│   ├── services/                      # 业务逻辑层
│   │   ├── family_service.go          # 家庭相关业务逻辑
│   │   └── user_service.go            # 用户相关业务逻辑
│   └── utils/                         # 通用工具集合
│       ├── BINDING_USAGE.md           # binding 工具的使用说明
│       ├── binding.go                 # 请求参数绑定封装
│       ├── jwt.go                     # JWT 生成与校验工具
│       ├── password.go                # 密码哈希与验证工具
│       ├── response.go                # 统一响应格式输出
│       └── validator.go               # 自定义参数校验逻辑
├── main                               # go build 生成的本地可执行文件
├── migrations/                        # 数据库迁移脚本（golang-migrate）
│   ├── 001_create_update_function.down.sql  # 回滚更新触发函数
│   ├── 001_create_update_function.up.sql    # 创建自动更新时间函数
│   ├── 002_create_users_table.down.sql      # 回滚用户表结构
│   ├── 002_create_users_table.up.sql        # 创建用户表
│   ├── 003_create_families_tables.down.sql  # 回滚家庭相关表
│   ├── 003_create_families_tables.up.sql    # 创建家庭相关表
│   ├── 004_create_dishes_tables.down.sql    # 回滚菜品相关表
│   ├── 004_create_dishes_tables.up.sql      # 创建菜品相关表
│   ├── 005_create_menus_tables.down.sql     # 回滚菜单相关表
│   ├── 005_create_menus_tables.up.sql       # 创建菜单相关表
│   ├── 006_create_shopping_lists_tables.down.sql  # 回滚购物清单表
│   ├── 006_create_shopping_lists_tables.up.sql    # 创建购物清单表
│   ├── 007_create_health_records_table.down.sql   # 回滚健康记录表
│   ├── 007_create_health_records_table.up.sql     # 创建健康记录表
│   ├── 008_create_memberships_table.down.sql      # 回滚成员关系表
│   ├── 008_create_memberships_table.up.sql        # 创建成员关系表
│   ├── 009_create_ai_usage_logs_table.down.sql    # 回滚 AI 调用日志表
│   ├── 009_create_ai_usage_logs_table.up.sql      # 创建 AI 调用日志表
│   ├── 010_create_payment_orders_table.down.sql   # 回滚支付订单表
│   ├── 010_create_payment_orders_table.up.sql     # 创建支付订单表
│   ├── 011_create_system_configs_table.down.sql   # 回滚系统配置表
│   └── 011_create_system_configs_table.up.sql     # 创建系统配置表
├── pkg/                               # 可复用公共库
│   └── database/                      # 数据库连接封装
│       └── postgres.go                # PostgreSQL 实例初始化
└── scripts/                           # 项目运维脚本
    ├── fix_migration.sh               # 批量修正迁移文件序号的小工具
    ├── migrate.sh                     # 运行数据库迁移的脚本
    └── swagger.sh                     # 调用 swag 生成最新 Swagger 文档
```

## 目录说明

### 顶层文件与目录
- `README.md`：介绍后端服务目标、运行方式与配置方法。
- `.idea/`：JetBrains 系列 IDE 的项目配置文件，不参与编译。
- `.DS_Store`：macOS 自动生成的目录缓存，可忽略。
- `go.mod` / `go.sum`：Go 模块依赖声明与校验，需与代码同步维护。
- `main`：执行 `go build` 后产生的可执行文件，建议只在本地调试阶段存在。

### cmd/
存放应用入口。`main.go` 负责：加载配置、初始化数据库与依赖、注册路由、启动 HTTP 服务。

### config/
集中存放 YAML 配置文件：
- `config.example.yaml` 提供完整字段的示例以便团队成员参考。
- `config.yaml` 为当前默认开发环境配置，可根据需要调整环境变量覆盖。

### docs/swagger/
Swag 命令生成的描述文件，提供 API 文档产物：
- `docs.go` 用于在 Go 应用内注册 Swagger 元数据。
- `swagger.json`/`swagger.yaml` 发布到外部或文档平台。
- `README.md` 记录如何生成与查看 Swagger。

### internal/
Go 模块的核心业务逻辑所在，也是项目最多文件的目录：
- `config/`：`config.go` 与 `loader.go` 负责定义配置结构并注入默认值。
- `handlers/`：REST 接口层，定义 gin 路由及请求处理；`router.go` 构建服务器，`routes.go` 列出所有路径，`auth_handler.go`/`user_handler.go`/`family_handler.go` 等承担具体模块逻辑。
- `middleware/`：目前仅有 `auth.go` JWT 鉴权中间件，在 `router.go` 中注册。
- `models/`：使用 struct 定义数据库表字段以及 JSON 标签。
- `repositories/`：封装数据库访问，便于在 service 层通过接口调用。
- `services/`：承载业务逻辑与事务控制，按实体拆分为 user/family。
- `utils/`：自定义工具，包括参数绑定、JWT、密码、响应包装和验证。`BINDING_USAGE.md` 解释 binding 设计思路。

### migrations/
使用 `golang-migrate` 的 SQL 脚本，严格按数字前缀排序，每个变更都有 up/down 成对文件。执行脚本请使用 `scripts/migrate.sh` 以保证顺序正确。

### pkg/
放置可以被其他模块重用的包，目前包含 `database/postgres.go`，负责连接池创建与复用，供 repository 层注入。

### scripts/
自动化运维脚本：
- `migrate.sh`：根据 `.env`/配置运行数据库迁移。
- `swagger.sh`：执行 `swag init` 并刷新 `docs/swagger` 内容。
- `fix_migration.sh`：辅助处理错误编号或批量改名，避免手动操作失误。

如需新增目录或文件，请同步更新本文件中的树形结构与说明，以保持文档与源码一致。
