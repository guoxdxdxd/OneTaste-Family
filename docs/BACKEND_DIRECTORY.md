# Backend 目录结构

本文档描述了 Go 后端服务的详细目录结构。

## 目录结构

```
src/backend/
├── cmd/                    # 应用程序入口
│   └── main.go            # 主程序入口文件，初始化应用并启动HTTP服务器
│
├── internal/               # 内部代码（不对外暴露）
│   ├── config/            # 配置管理
│   │   ├── config.go      # 配置结构定义，包含所有配置项的结构体
│   │   └── loader.go      # 配置加载器，从文件或环境变量加载配置
│   │
│   ├── models/            # 数据模型
│   │   ├── user.go        # 用户模型，定义用户数据结构
│   │   ├── family.go      # 家庭模型，定义家庭数据结构
│   │   ├── recipe.go      # 食谱模型，定义食谱数据结构
│   │   ├── ingredient.go  # 食材模型，定义食材数据结构
│   │   └── shopping.go    # 购物清单模型，定义购物清单数据结构
│   │
│   ├── handlers/          # HTTP处理器（Controller层）
│   │   ├── user_handler.go        # 用户相关HTTP请求处理器
│   │   ├── family_handler.go      # 家庭相关HTTP请求处理器
│   │   ├── recipe_handler.go      # 食谱相关HTTP请求处理器
│   │   ├── ingredient_handler.go  # 食材相关HTTP请求处理器
│   │   └── shopping_handler.go    # 购物清单相关HTTP请求处理器
│   │
│   ├── services/          # 业务逻辑层
│   │   ├── user_service.go        # 用户业务逻辑服务
│   │   ├── family_service.go      # 家庭业务逻辑服务
│   │   ├── recipe_service.go      # 食谱业务逻辑服务
│   │   ├── ingredient_service.go  # 食材业务逻辑服务
│   │   ├── shopping_service.go    # 购物清单业务逻辑服务
│   │   └── ai_service.go         # AI服务调用封装，与AI服务交互
│   │
│   ├── repositories/      # 数据访问层（Repository层）
│   │   ├── user_repository.go        # 用户数据访问层，数据库CRUD操作
│   │   ├── family_repository.go      # 家庭数据访问层，数据库CRUD操作
│   │   ├── recipe_repository.go      # 食谱数据访问层，数据库CRUD操作
│   │   ├── ingredient_repository.go  # 食材数据访问层，数据库CRUD操作
│   │   └── shopping_repository.go    # 购物清单数据访问层，数据库CRUD操作
│   │
│   ├── middleware/        # 中间件
│   │   ├── auth.go        # 认证中间件，JWT token验证
│   │   ├── cors.go        # CORS中间件，跨域请求处理
│   │   ├── logger.go      # 日志中间件，记录请求日志
│   │   ├── recovery.go    # 错误恢复中间件，捕获panic并返回友好错误
│   │   └── validator.go   # 参数验证中间件，验证请求参数
│   │
│   └── utils/            # 工具函数
│       ├── jwt.go        # JWT工具，生成和验证JWT token
│       ├── password.go   # 密码加密工具，密码哈希和验证
│       ├── validator.go  # 验证工具，数据验证函数
│       ├── response.go   # 响应格式化，统一API响应格式
│       └── logger.go     # 日志工具，日志记录封装
│
├── pkg/                   # 可复用的公共包
│   ├── database/         # 数据库连接
│   │   ├── postgres.go   # PostgreSQL数据库连接封装
│   │   └── redis.go      # Redis连接封装
│   │
│   └── errors/           # 错误定义
│       └── errors.go     # 自定义错误类型定义
│
├── migrations/            # 数据库迁移脚本
│   ├── 001_create_users.up.sql        # 创建用户表的迁移脚本（向上）
│   ├── 001_create_users.down.sql      # 删除用户表的迁移脚本（向下）
│   ├── 002_create_families.up.sql     # 创建家庭表的迁移脚本（向上）
│   └── ...                            # 其他迁移脚本
│
├── config/                # 配置文件目录
│   ├── config.yaml        # 配置文件（示例），包含应用配置
│   └── config.example.yaml # 配置文件模板，供参考使用
│
├── tests/                 # 测试文件
│   ├── handlers/          # HTTP处理器测试文件
│   ├── services/          # 业务逻辑层测试文件
│   └── repositories/      # 数据访问层测试文件
│
├── scripts/              # 脚本文件
│   └── migrate.sh        # 数据库迁移脚本，执行迁移命令
│
├── go.mod                # Go模块定义文件，定义模块名称和依赖
├── go.sum                # Go依赖校验文件，依赖包的校验和
├── .gitignore            # Git忽略文件，定义不提交到版本控制的文件
├── Dockerfile            # Docker构建文件，用于构建Docker镜像
└── README.md             # 后端服务说明文档，项目说明和使用指南
```

## 目录说明

### cmd/
应用程序的入口点，包含 `main.go` 文件，负责：
- 初始化配置
- 初始化数据库连接
- 初始化路由
- 启动HTTP服务器

### internal/
内部代码目录，遵循 Go 的包可见性规则，这些代码不会被外部包导入。

#### internal/config/
配置管理模块：
- 读取配置文件（YAML/JSON/环境变量）
- 配置结构体定义
- 配置验证

#### internal/models/
数据模型定义：
- 数据库表对应的结构体
- JSON序列化标签
- 验证标签

#### internal/handlers/
HTTP处理器（Controller层）：
- 接收HTTP请求
- 参数验证
- 调用Service层
- 返回HTTP响应

#### internal/services/
业务逻辑层：
- 实现核心业务逻辑
- 调用Repository层
- 调用外部服务（如AI服务）
- 事务管理

#### internal/repositories/
数据访问层：
- 数据库CRUD操作
- SQL查询封装
- 数据库连接管理

#### internal/middleware/
中间件：
- 认证授权
- 请求日志
- 错误处理
- 跨域处理
- 参数验证

#### internal/utils/
工具函数：
- JWT生成和验证
- 密码加密和验证
- 响应格式化
- 日志工具

### pkg/
可复用的公共包，可以被外部项目导入：
- 数据库连接封装
- 错误类型定义
- 通用工具函数

### migrations/
数据库迁移脚本：
- 使用 `golang-migrate` 或类似工具
- 每个迁移包含 up 和 down 脚本
- 版本化管理数据库结构变更

### config/
配置文件目录：
- 开发环境配置
- 生产环境配置示例
- 配置文件模板

### tests/
测试文件：
- 单元测试
- 集成测试
- 测试工具和辅助函数

## 文件命名规范

### Go文件
- 使用小写字母和下划线：`user_service.go`
- 测试文件：`user_service_test.go`
- 主程序：`main.go`

### 包命名
- 使用小写字母，简短且有意义
- 避免使用下划线或混合大小写
- 单数形式：`user` 而不是 `users`

### 结构体命名
- 使用PascalCase：`User`, `FamilyRecipe`
- 导出结构体首字母大写

### 函数命名
- 导出函数使用PascalCase：`GetUser()`
- 私有函数使用camelCase：`validateUser()`

## 代码组织原则

1. **分层架构**：
   - Handler层：处理HTTP请求
   - Service层：业务逻辑
   - Repository层：数据访问

2. **依赖注入**：
   - 通过构造函数注入依赖
   - 便于测试和扩展

3. **错误处理**：
   - 使用自定义错误类型
   - 统一的错误响应格式

4. **配置管理**：
   - 使用配置文件和环境变量
   - 支持多环境配置

5. **日志记录**：
   - 结构化日志
   - 不同级别的日志输出

## 开发规范

### 导入顺序
```go
import (
    // 标准库
    "fmt"
    "net/http"
    
    // 第三方库
    "github.com/gin-gonic/gin"
    "github.com/lib/pq"
    
    // 内部包
    "onetaste-family/backend/internal/models"
    "onetaste-family/backend/internal/services"
)
```

### 错误处理
```go
if err != nil {
    log.Errorf("failed to get user: %v", err)
    return nil, fmt.Errorf("get user failed: %w", err)
}
```

### 注释规范
- 导出函数和类型必须有注释
- 注释以被注释对象名称开头
- 使用完整的句子

## 相关文档

- [API接口文档](./接口文档.md)
- [数据库设计文档](./数据库设计.md)
- [部署文档](./部署文档.md)

