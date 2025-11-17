# Frontend Admin 目录结构

本文档描述了管理后台（Vue 3）的详细目录结构。

## 目录结构

```
src/frontend/admin/
├── src/                    # 源代码目录
│   ├── assets/            # 静态资源
│   │   ├── images/       # 图片资源目录，存放项目使用的图片文件
│   │   ├── icons/        # 图标资源目录，存放图标文件
│   │   └── fonts/        # 字体文件目录，存放自定义字体文件
│   │
│   ├── components/        # 公共组件
│   │   ├── common/       # 通用组件
│   │   │   ├── Button.vue        # 按钮组件，通用按钮UI组件
│   │   │   ├── Input.vue         # 输入框组件，通用输入框UI组件
│   │   │   ├── Table.vue         # 表格组件，通用数据表格组件
│   │   │   ├── Dialog.vue        # 对话框组件，通用弹窗组件
│   │   │   └── Pagination.vue    # 分页组件，通用分页UI组件
│   │   │
│   │   ├── layout/       # 布局组件
│   │   │   ├── Header.vue        # 头部组件，页面顶部导航栏
│   │   │   ├── Sidebar.vue      # 侧边栏组件，左侧导航菜单
│   │   │   ├── Footer.vue       # 底部组件，页面底部信息
│   │   │   └── Breadcrumb.vue   # 面包屑组件，页面路径导航
│   │   │
│   │   └── business/     # 业务组件
│   │       ├── UserCard.vue         # 用户卡片组件，显示用户信息
│   │       ├── RecipeCard.vue       # 食谱卡片组件，显示食谱信息
│   │       └── StatisticsCard.vue  # 统计卡片组件，显示统计数据
│   │
│   ├── views/             # 页面组件
│   │   ├── dashboard/    # 仪表盘
│   │   │   └── Dashboard.vue       # 仪表盘页面，数据概览和统计
│   │   │
│   │   ├── user/         # 用户管理
│   │   │   ├── UserList.vue        # 用户列表页面，显示所有用户
│   │   │   ├── UserDetail.vue      # 用户详情页面，显示用户详细信息
│   │   │   └── UserEdit.vue        # 用户编辑页面，编辑用户信息
│   │   │
│   │   ├── family/       # 家庭管理
│   │   │   ├── FamilyList.vue      # 家庭列表页面，显示所有家庭
│   │   │   └── FamilyDetail.vue    # 家庭详情页面，显示家庭详细信息
│   │   │
│   │   ├── recipe/       # 食谱管理
│   │   │   ├── RecipeList.vue      # 食谱列表页面，显示所有食谱
│   │   │   ├── RecipeDetail.vue    # 食谱详情页面，显示食谱详细信息
│   │   │   └── RecipeEdit.vue      # 食谱编辑页面，编辑食谱信息
│   │   │
│   │   └── system/       # 系统设置
│   │       ├── Settings.vue        # 系统设置页面，系统配置管理
│   │       └── Logs.vue            # 日志页面，查看系统日志
│   │
│   ├── router/            # 路由配置
│   │   ├── index.js      # 路由主文件，定义所有路由规则
│   │   └── guards.js     # 路由守卫，处理路由权限和登录验证
│   │
│   ├── stores/            # 状态管理（Pinia）
│   │   ├── user.js       # 用户状态管理，管理用户信息和登录状态
│   │   ├── app.js        # 应用状态管理，管理全局应用状态
│   │   ├── recipe.js     # 食谱状态管理，管理食谱相关状态
│   │   └── family.js     # 家庭状态管理，管理家庭相关状态
│   │
│   ├── api/               # API接口
│   │   ├── index.js      # API配置文件，HTTP请求基础配置
│   │   ├── user.js       # 用户API接口，用户相关的API调用
│   │   ├── family.js     # 家庭API接口，家庭相关的API调用
│   │   ├── recipe.js     # 食谱API接口，食谱相关的API调用
│   │   └── ingredient.js # 食材API接口，食材相关的API调用
│   │
│   ├── utils/             # 工具函数
│   │   ├── request.js    # HTTP请求封装，axios封装和拦截器
│   │   ├── auth.js       # 认证工具，token管理和认证相关函数
│   │   ├── storage.js    # 本地存储工具，localStorage和sessionStorage封装
│   │   ├── format.js     # 格式化工具，日期、金额等格式化函数
│   │   ├── validate.js   # 验证工具，表单验证和数据验证函数
│   │   └── constants.js # 常量定义，项目常量配置
│   │
│   ├── styles/            # 样式文件
│   │   ├── variables.scss # SCSS变量文件，定义颜色、尺寸等变量
│   │   ├── mixins.scss   # SCSS混入文件，定义可复用的样式混入
│   │   ├── common.scss   # 通用样式文件，全局通用样式
│   │   └── reset.scss    # 样式重置文件，重置浏览器默认样式
│   │
│   ├── composables/       # 组合式函数
│   │   ├── useAuth.js    # 认证相关组合函数，认证逻辑封装
│   │   ├── useTable.js   # 表格相关组合函数，表格操作逻辑封装
│   │   └── useForm.js    # 表单相关组合函数，表单处理逻辑封装
│   │
│   ├── directives/        # 自定义指令
│   │   ├── permission.js # 权限指令，控制元素显示权限
│   │   └── loading.js    # 加载指令，显示加载状态
│   │
│   ├── plugins/           # 插件
│   │   └── element-plus.js # Element Plus插件配置，UI库初始化
│   │
│   ├── App.vue           # 根组件，应用根组件
│   └── main.js           # 入口文件，应用启动入口
│
├── public/                # 公共资源
│   ├── favicon.ico       # 网站图标，浏览器标签页图标
│   └── index.html        # HTML模板，应用HTML入口模板
│
├── .env                   # 环境变量文件（不提交到版本控制），包含敏感配置
├── .env.development       # 开发环境变量文件，开发环境配置
├── .env.production        # 生产环境变量文件，生产环境配置
├── .env.example           # 环境变量示例文件，环境变量模板
│
├── package.json           # 项目依赖配置文件，定义项目依赖和脚本
├── package-lock.json      # 依赖锁定文件，锁定依赖版本
├── vite.config.js         # Vite配置文件，构建工具配置
├── .gitignore             # Git忽略文件，定义不提交到版本控制的文件
├── .eslintrc.js           # ESLint配置文件，代码检查规则
├── .prettierrc            # Prettier配置文件，代码格式化规则
└── README.md              # 项目说明文档，项目说明和使用指南
```

## 目录说明

### src/assets/
静态资源目录：
- 图片、图标、字体等静态文件
- 这些文件会被Vite处理并打包

### src/components/
组件目录，按功能分类：
- **common/**: 通用UI组件，可复用的基础组件
- **layout/**: 布局组件，用于页面结构
- **business/**: 业务组件，特定业务场景的组件

### src/views/
页面组件目录：
- 每个路由对应的页面组件
- 按功能模块组织

### src/router/
路由配置：
- 路由定义
- 路由守卫（权限控制、登录验证）

### src/stores/
状态管理（Pinia）：
- 全局状态定义
- 状态操作方法
- 按模块划分

### src/api/
API接口封装：
- 统一的HTTP请求封装
- 按业务模块组织API方法
- 请求拦截和响应处理

### src/utils/
工具函数：
- 通用工具函数
- 业务相关的工具函数
- 常量定义

### src/styles/
样式文件：
- SCSS变量和混入
- 全局样式
- 样式重置

### src/composables/
组合式函数（Composition API）：
- 可复用的逻辑封装
- 类似React Hooks的概念

### src/directives/
自定义指令：
- Vue自定义指令
- 如权限控制、加载状态等

### src/plugins/
插件配置：
- 第三方库的初始化配置
- 如Element Plus、Ant Design Vue等

### public/
公共资源目录：
- 不会被Vite处理的静态文件
- 直接复制到构建输出目录

## 文件命名规范

### Vue组件
- 使用PascalCase：`UserProfile.vue`
- 页面组件：`UserList.vue`, `UserDetail.vue`
- 组件名与文件名保持一致

### JavaScript文件
- 使用camelCase：`userService.js`
- 工具文件：`formatDate.js`, `validateForm.js`
- 常量文件：`constants.js`

### 样式文件
- 使用kebab-case：`user-profile.scss`
- 或使用功能命名：`common.scss`, `variables.scss`

## 代码组织原则

1. **组件化开发**：
   - 单一职责原则
   - 可复用性
   - 组件通信规范化

2. **状态管理**：
   - 全局状态使用Pinia
   - 局部状态使用组件内部状态
   - 避免过度使用全局状态

3. **API封装**：
   - 统一的请求拦截器
   - 统一的错误处理
   - 统一的响应格式

4. **路由管理**：
   - 路由懒加载
   - 路由守卫统一管理
   - 路由元信息配置

5. **样式管理**：
   - 使用SCSS变量统一主题
   - 组件样式作用域化
   - 避免全局样式污染

## 开发规范

### 组件结构
```vue
<template>
  <!-- 模板内容 -->
</template>

<script setup>
// 导入依赖
import { ref, computed } from 'vue'

// 定义props
const props = defineProps({
  // ...
})

// 定义emits
const emit = defineEmits(['update'])

// 响应式数据
const count = ref(0)

// 计算属性
const doubleCount = computed(() => count.value * 2)

// 方法
const handleClick = () => {
  // ...
}
</script>

<style scoped>
/* 样式 */
</style>
```

### API调用示例
```javascript
// api/user.js
import request from '@/utils/request'

export function getUserList(params) {
  return request.get('/api/users', { params })
}

export function getUserById(id) {
  return request.get(`/api/users/${id}`)
}
```

### 状态管理示例
```javascript
// stores/user.js
import { defineStore } from 'pinia'

export const useUserStore = defineStore('user', {
  state: () => ({
    userInfo: null,
    token: ''
  }),
  
  actions: {
    async login(credentials) {
      // 登录逻辑
    },
    
    logout() {
      // 登出逻辑
    }
  }
})
```

## 相关文档

- [API接口文档](./接口文档.md)
- [数据库设计文档](./数据库设计.md)
- [部署文档](./部署文档.md)


