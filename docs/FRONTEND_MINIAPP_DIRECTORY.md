# Frontend Miniapp 目录结构

本文档依据当前仓库 `src/frontend/miniapp` 的真实文件生成，完整列出项目包含的目录与文件，并给出简要功能说明，便于定位代码位置和职责。

## 目录结构（含注释）

```
src/frontend/miniapp/
├── .env                              # 本地开发环境变量（不会提交）
├── .env.example                      # 环境变量示例，说明需要配置的键
├── .gitignore                        # Git 忽略规则，排除 node_modules 等目录
├── README.md                         # 小程序端项目说明与运行指南
├── index.html                        # Vite 入口 HTML 模板
├── node_modules/                     # npm 安装的依赖目录
├── package-lock.json                 # 锁定依赖的版本与完整性
├── package.json                      # 项目依赖、脚本命令与元数据
├── src/                              # 小程序端源代码
│   ├── App.vue                       # 根组件，负责全局布局与路由出口
│   ├── api/                          # HTTP API 封装
│   │   ├── index.js                  # axios 实例与拦截器配置
│   │   └── user.js                   # 用户注册/登录/信息接口
│   ├── main.js                       # Vue/Vite 入口文件，挂载应用
│   ├── pages/                        # 业务页面
│   │   ├── auth/                     # 认证流程页面
│   │   │   ├── Login.vue             # 登录页面
│   │   │   └── Register.vue          # 注册页面
│   │   └── home/                     # 首页模块
│   │       └── Home.vue              # 首页展示页面
│   ├── router/                       # Vue Router 配置
│   │   ├── guards.js                 # 路由守卫，处理鉴权与重定向
│   │   └── index.js                  # 路由表与实例化逻辑
│   ├── stores/                       # Pinia 状态仓库
│   │   └── user.js                   # 用户登录状态与信息管理
│   ├── styles/                       # 全局样式
│   │   └── theme.css                 # 主题色与公共样式变量
│   └── utils/                        # 通用工具方法
│       ├── auth.js                   # token 与用户上下文的读写
│       ├── constants.js              # 应用常量定义
│       ├── request.js                # axios 请求封装与拦截处理
│       ├── storage.js                # localStorage/sessionStorage 封装
│       └── validate.js               # 表单与字段验证函数
└── vite.config.js                    # Vite 构建与代理配置
```

## 目录说明

### 顶层文件
- `.env` / `.env.example`：维护运行所需的后台 API 地址等变量，example 用于指导配置。
- `.gitignore`：忽略 node_modules、构建产物和本地配置。
- `README.md`：说明如何安装依赖、运行与构建小程序。
- `index.html`：Vite 的模板入口，注入打包后的脚本。
- `package.json` / `package-lock.json`：记录依赖与 npm scripts，任何库升级都需同步锁文件。
- `node_modules/`：npm 自动生成，勿手动修改。
- `vite.config.js`：负责配置路径别名、代理、构建优化等。

### src/
源代码主目录：
- `App.vue`：根组件，包含 `<router-view />`，可挂载全局布局逻辑。
- `main.js`：创建 Vue 应用，注册 Router、Pinia、全局样式等。
- `api/`：集中封装 HTTP 调用；`index.js` 创建 axios 实例并配置 token 拦截，`user.js` 提供注册/登录/获取信息等具体接口。
- `pages/`：按业务切分的 Vue 组件，目前包含 `home` 与 `auth` 两个模块，后续新增页面请在此按模块建子目录。
- `router/`：`index.js` 定义路由表、懒加载组件，`guards.js` 中实现鉴权、重定向等逻辑。
- `stores/`：Pinia 状态管理，目前只有 `user.js`，负责 token、登录态和用户信息；未来可在此新增家庭、购物等模块。
- `styles/`：全局样式与主题变量集合，`theme.css` 被 `main.js` 引入。
- `utils/`：复用工具。`auth.js` 管理 token 与用户缓存，`request.js` 统一 HTTP 错误处理，`storage.js` 封装浏览器存储，`constants.js` 保留常量，`validate.js` 提供校验器。

### pages/
- `pages/home/Home.vue`：应用首页，与后端接口交互展示概览内容。
- `pages/auth/Login.vue`：登录表单，调用 `api/user.login` 与 `useUserStore`。
- `pages/auth/Register.vue`：注册表单，复用校验逻辑并完成注册流程。

### router/
- `router/index.js`：配置 `createRouter`，导入页面组件，注册路由守卫。
- `router/guards.js`：集中放置 `beforeEach` 逻辑，检查 `useUserStore` 的 `isLoggedIn` 并统一跳转。

### stores/
- `stores/user.js`：唯一的 Pinia store，封装注册/登录/获取用户信息等流程，暴露 getter（如 `loggedIn`、`membershipType`）和 action。

### styles/
- `styles/theme.css`：维护 CSS 变量、基础排版与公共颜色，供全局导入。

### utils/
- `auth.js`：调用 `storage.js`，封装 token 持久化与用户信息缓存。
- `constants.js`：集中放置常用常量（如 token key、路由名等）。
- `request.js`：创建 axios 实例、配置请求/响应拦截器、统一错误抛出。
- `storage.js`：定义 `userInfoStorage` 等读写 helper，避免直接操作 `localStorage`。
- `validate.js`：常用校验函数，例如手机号、验证码、密码复杂度等。

## 代码组织原则

1. **组件化开发**：将业务页面拆分到 `src/pages/{module}` 下，满足单一职责；跨页面的可复用 UI 请抽离到未来的 `components/` 目录。
2. **状态管理**：Pinia Store 均放在 `src/stores`，每个 Store 专注一个上下文（当前只有 `user`），借助 `persist` 能力通过 `storage.js` 持久化关键信息。
3. **API 封装**：所有后端请求通过 `src/api` 调用 `utils/request.js` 创建的 axios 实例；统一添加 token、处理错误与重试逻辑。
4. **路由管理**：`router/index.js` 采用模块化路由定义；`guards.js` 统一进行登录校验与跳转，避免在组件内重复判断。
5. **样式管理**：全局样式统一集中在 `styles/theme.css`，组件内鼓励使用 scoped 样式并通过 CSS 变量共享主题。

## 代码示例

### API 调用示例（`src/api/user.js`）
```javascript
import request from './index'

export function login(data) {
  return request.post('/auth/login', data)
}

export function register(data) {
  return request.post('/auth/register', data)
}

export function getUserInfo() {
  return request.get('/user/info')
}
```

### 状态管理示例（`src/stores/user.js` 摘录）
```javascript
import { defineStore } from 'pinia'
import { setAuth, clearAuth, getToken, getUserInfo } from '@/utils/auth'
import { userInfoStorage } from '@/utils/storage'
import { register as registerApi, login as loginApi, getUserInfo as getUserInfoApi } from '@/api/user'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: getToken(),
    userInfo: getUserInfo(),
    isLoggedIn: !!getToken()
  }),
  actions: {
    async register(data) {
      const res = await registerApi(data)
      if (res.code === 200 && res.data) {
        this.setAuth(res.data.token, {
          user_id: res.data.user_id,
          phone: data.phone,
          nickname: data.nickname
        })
      }
      return res
    },
    async login(data) {
      const res = await loginApi(data)
      if (res.code === 200 && res.data) {
        this.setAuth(res.data.token, { user_id: res.data.user_id, phone: data.phone })
        await this.fetchUserInfo()
      }
      return res
    },
    async fetchUserInfo() {
      const res = await getUserInfoApi()
      if (res.code === 200 && res.data) {
        this.userInfo = res.data
        userInfoStorage.setUserInfo(res.data)
      }
      return res
    },
    setAuth(token, userInfo) {
      this.token = token
      this.userInfo = userInfo
      this.isLoggedIn = true
      setAuth(token, userInfo)
    },
    logout() {
      this.token = null
      this.userInfo = null
      this.isLoggedIn = false
      clearAuth()
    }
  }
})
```

### 页面结构示例（`src/pages/auth/Login.vue`）
```vue
<template>
  <form @submit.prevent="handleLogin">
    <!-- 登录表单内容 -->
  </form>
</template>

<script setup>
import { ref } from 'vue'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
const form = ref({ phone: '', password: '' })

const handleLogin = async () => {
  await userStore.login(form.value)
}
</script>
```

当新增或调整目录结构时，请同步更新本文件，确保文档与代码保持一致。
