# Miniapp Frontend

小程序前端（Vue 3）

## 技术栈

- Vue 3 (Composition API)
- Vue Router 4
- Pinia (状态管理)
- Axios (HTTP请求)
- Vite (构建工具)

## 项目结构

```
src/
├── api/              # API接口
├── components/       # 组件
├── pages/           # 页面
│   ├── auth/        # 认证页面（登录、注册）
│   └── home/        # 首页
├── router/          # 路由配置
├── stores/          # 状态管理
├── utils/           # 工具函数
├── App.vue          # 根组件
└── main.js          # 入口文件
```

## 开发

### 安装依赖

```bash
npm install
```

### 配置环境变量

复制 `.env.example` 为 `.env` 并配置：

```bash
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

### 启动开发服务器

```bash
npm run dev
```

访问 http://localhost:3000

## 构建

```bash
npm run build
```

## 功能说明

### 用户认证模块

已实现的功能：

1. **用户注册**
   - 手机号验证
   - 验证码输入（UI已实现，接口待对接）
   - 密码验证（6-20位）
   - 昵称输入
   - 表单验证和错误提示

2. **用户登录**
   - 手机号登录
   - 密码登录
   - Token存储和管理
   - 登录状态管理

3. **Token管理**
   - 自动存储Token到localStorage
   - 请求拦截器自动添加Token
   - 401错误自动清除Token并跳转登录

4. **路由守卫**
   - 需要登录的页面自动验证
   - 未登录自动跳转到登录页
   - 登录后跳转到目标页面

## API接口

### 用户认证接口

- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录
- `GET /api/v1/user/info` - 获取用户信息

详细接口文档请参考：`docs/接口文档.md`

## 注意事项

1. 验证码发送功能需要对接后端接口
2. 小程序环境需要适配微信小程序API
3. 生产环境需要配置正确的API地址

