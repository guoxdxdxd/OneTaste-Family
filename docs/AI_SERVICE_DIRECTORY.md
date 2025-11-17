# AI Service 目录结构

本文档描述了 Python AI 服务的详细目录结构。

## 目录结构

```
src/ai-service/
├── app/                    # 应用代码目录
│   ├── __init__.py        # 包初始化文件，使app成为Python包
│   │
│   ├── main.py            # FastAPI应用入口，创建应用实例并启动服务
│   │
│   ├── config/            # 配置管理
│   │   ├── __init__.py    # 包初始化文件
│   │   ├── settings.py    # 配置类，定义所有配置项的数据结构
│   │   └── env.py         # 环境变量加载，从.env文件加载环境变量
│   │
│   ├── models/            # 数据模型
│   │   ├── __init__.py    # 包初始化文件
│   │   ├── request.py     # 请求模型，定义API请求参数的数据模型
│   │   ├── response.py    # 响应模型，定义API响应数据的数据模型
│   │   └── recipe.py      # 食谱相关模型，食谱相关的数据模型定义
│   │
│   ├── routers/           # 路由（API端点）
│   │   ├── __init__.py    # 包初始化文件
│   │   ├── health.py      # 健康检查路由，提供健康检查API端点
│   │   ├── recipe.py      # 食谱推荐路由，食谱推荐相关的API端点
│   │   ├── chat.py        # 对话路由，AI对话相关的API端点
│   │   └── asr.py         # 语音识别路由，语音识别相关的API端点
│   │
│   ├── services/          # 业务逻辑层
│   │   ├── __init__.py    # 包初始化文件
│   │   ├── openai_service.py         # OpenAI服务，调用OpenAI API的封装
│   │   ├── baidu_asr_service.py      # 百度ASR服务，调用百度语音识别API
│   │   ├── recipe_service.py         # 食谱服务，食谱相关的业务逻辑
│   │   └── recommendation_service.py  # 推荐服务，智能推荐算法实现
│   │
│   ├── utils/             # 工具函数
│   │   ├── __init__.py    # 包初始化文件
│   │   ├── logger.py      # 日志工具，日志记录封装和配置
│   │   ├── cache.py       # 缓存工具，Redis缓存操作封装
│   │   ├── validator.py   # 验证工具，数据验证函数
│   │   └── formatter.py   # 格式化工具，数据格式化函数
│   │
│   ├── middleware/        # 中间件
│   │   ├── __init__.py    # 包初始化文件
│   │   ├── cors.py        # CORS中间件，处理跨域请求
│   │   ├── auth.py        # 认证中间件，API认证验证
│   │   └── logger.py      # 日志中间件，记录请求日志
│   │
│   └── exceptions/        # 异常处理
│       ├── __init__.py    # 包初始化文件
│       ├── base.py        # 基础异常类，定义自定义异常基类
│       └── handlers.py    # 异常处理器，全局异常处理和错误响应
│
├── tests/                  # 测试文件
│   ├── __init__.py        # 包初始化文件
│   ├── test_services/     # 服务层测试
│   │   ├── test_openai_service.py    # OpenAI服务测试
│   │   └── test_recipe_service.py    # 食谱服务测试
│   │
│   ├── test_routers/      # 路由层测试
│   │   ├── test_recipe.py # 食谱路由测试
│   │   └── test_chat.py   # 对话路由测试
│   │
│   └── conftest.py        # pytest配置文件，测试fixtures和配置
│
├── scripts/               # 脚本文件
│   ├── init_db.py        # 数据库初始化脚本，初始化数据库结构
│   └── migrate.py         # 数据迁移脚本，执行数据迁移
│
├── logs/                  # 日志目录（不提交到版本控制）
│   └── .gitkeep           # Git保持文件，确保logs目录被版本控制
│
├── .env                   # 环境变量文件（不提交到版本控制），包含敏感配置
├── .env.example           # 环境变量示例文件，环境变量配置模板
├── .gitignore             # Git忽略文件，定义不提交到版本控制的文件
│
├── requirements.txt       # Python依赖列表，生产环境依赖
├── requirements-dev.txt  # 开发依赖列表，开发环境额外依赖
│
├── Dockerfile             # Docker构建文件，用于构建Docker镜像
├── docker-compose.yml     # Docker Compose配置文件（可选），服务编排配置
│
├── README.md              # 项目说明文档，项目说明和使用指南
└── pytest.ini            # pytest配置文件，pytest测试框架配置
```

## 目录说明

### app/
应用代码主目录，包含所有业务逻辑代码。

#### app/main.py
FastAPI应用入口：
- 创建FastAPI应用实例
- 注册路由
- 注册中间件
- 启动配置

#### app/config/
配置管理模块：
- 从环境变量加载配置
- 配置类定义
- 配置验证

#### app/models/
数据模型（Pydantic模型）：
- 请求模型：定义API请求参数
- 响应模型：定义API响应格式
- 数据验证和序列化

#### app/routers/
API路由定义：
- 按功能模块组织路由
- 每个路由文件对应一个功能模块
- 使用FastAPI的APIRouter

#### app/services/
业务逻辑层：
- 封装外部API调用（OpenAI、百度ASR等）
- 实现核心业务逻辑
- 数据处理和转换

#### app/utils/
工具函数：
- 日志工具：统一日志格式
- 缓存工具：Redis缓存封装
- 验证工具：数据验证
- 格式化工具：数据格式化

#### app/middleware/
中间件：
- CORS处理
- 认证验证
- 请求日志记录
- 错误处理

#### app/exceptions/
异常处理：
- 自定义异常类
- 全局异常处理器
- 错误响应格式化

### tests/
测试文件目录：
- 单元测试
- 集成测试
- 测试工具和fixtures

### scripts/
脚本文件：
- 数据库初始化脚本
- 数据迁移脚本
- 工具脚本

### logs/
日志文件目录：
- 应用运行日志
- 错误日志
- 访问日志

## 文件命名规范

### Python文件
- 使用小写字母和下划线：`user_service.py`
- 测试文件：`test_user_service.py`
- 主程序：`main.py`

### 类命名
- 使用PascalCase：`UserService`, `RecipeModel`
- 异常类：`UserNotFoundError`

### 函数和变量命名
- 使用snake_case：`get_user()`, `user_id`
- 常量：全大写，下划线分隔：`MAX_RETRY_COUNT`

### 模块命名
- 使用小写字母和下划线：`recipe_service.py`
- 避免使用连字符

## 代码组织原则

1. **分层架构**：
   - Router层：处理HTTP请求
   - Service层：业务逻辑
   - Model层：数据模型

2. **依赖注入**：
   - 使用FastAPI的依赖注入系统
   - 便于测试和扩展

3. **错误处理**：
   - 使用自定义异常类
   - 统一的错误响应格式
   - 全局异常处理

4. **配置管理**：
   - 使用环境变量
   - Pydantic Settings管理配置
   - 支持多环境配置

5. **日志记录**：
   - 结构化日志
   - 不同级别的日志输出
   - 日志轮转

## 开发规范

### 导入顺序
```python
# 标准库
import os
from typing import List, Optional

# 第三方库
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import openai

# 本地模块
from app.models import RecipeRequest
from app.services import RecipeService
```

### 类型提示
```python
def get_recipe(recipe_id: int) -> Optional[Recipe]:
    """获取食谱信息"""
    pass
```

### 文档字符串
```python
def recommend_recipe(user_preferences: dict) -> List[Recipe]:
    """
    根据用户偏好推荐食谱
    
    Args:
        user_preferences: 用户偏好字典，包含口味、食材等信息
        
    Returns:
        推荐的食谱列表
        
    Raises:
        ValueError: 当用户偏好无效时
    """
    pass
```

### 异常处理
```python
try:
    result = await service.process()
except ServiceError as e:
    logger.error(f"Service error: {e}")
    raise HTTPException(status_code=500, detail=str(e))
```

## FastAPI 最佳实践

### 路由定义
```python
from fastapi import APIRouter, Depends
from app.models import RecipeRequest, RecipeResponse
from app.services import RecipeService

router = APIRouter(prefix="/api/v1/recipes", tags=["recipes"])

@router.post("/recommend", response_model=RecipeResponse)
async def recommend_recipe(
    request: RecipeRequest,
    service: RecipeService = Depends()
):
    """推荐食谱"""
    return await service.recommend(request)
```

### 依赖注入
```python
from fastapi import Depends
from app.services import RecipeService

def get_recipe_service() -> RecipeService:
    return RecipeService()

@router.get("/{recipe_id}")
async def get_recipe(
    recipe_id: int,
    service: RecipeService = Depends(get_recipe_service)
):
    return await service.get_by_id(recipe_id)
```

### 中间件使用
```python
from fastapi.middleware.cors import CORSMiddleware

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)
```

## 环境变量配置

### .env.example
```env
# OpenAI配置
OPENAI_API_KEY=your_openai_api_key
OPENAI_MODEL=gpt-3.5-turbo

# 百度ASR配置
BAIDU_ASR_API_KEY=your_baidu_key
BAIDU_ASR_SECRET_KEY=your_baidu_secret

# 服务配置
HOST=0.0.0.0
PORT=8000
DEBUG=False

# Redis配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_DB=0
```

## 测试规范

### 单元测试示例
```python
import pytest
from app.services import RecipeService

@pytest.mark.asyncio
async def test_recommend_recipe():
    service = RecipeService()
    result = await service.recommend({"taste": "sweet"})
    assert result is not None
    assert len(result) > 0
```

## 相关文档

- [API接口文档](./API.md)
- [部署文档](./DEPLOYMENT.md)
- [技术栈文档](./TECH_STACK.md)

