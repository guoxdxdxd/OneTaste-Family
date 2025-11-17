# Frontend Miniapp 目录结构

本文档描述了小程序（Vue 3）的详细目录结构。

## 目录结构

```
src/frontend/miniapp/
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
│   │   │   ├── Card.vue          # 卡片组件，通用卡片UI组件
│   │   │   ├── Tag.vue           # 标签组件，通用标签UI组件
│   │   │   └── Empty.vue         # 空状态组件，显示空数据状态
│   │   │
│   │   └── business/     # 业务组件
│   │       ├── RecipeCard.vue        # 食谱卡片组件，显示食谱信息卡片
│   │       ├── IngredientCard.vue   # 食材卡片组件，显示食材信息卡片
│   │       ├── ShoppingItem.vue      # 购物清单项组件，显示购物清单单项
│   │       └── FamilyMember.vue     # 家庭成员组件，显示家庭成员信息
│   │
│   ├── pages/             # 页面组件
│   │   ├── home/         # 首页
│   │   │   └── Home.vue              # 首页组件，应用首页展示
│   │   │
│   │   ├── recipe/       # 食谱相关
│   │   │   ├── RecipeList.vue       # 食谱列表页面，显示所有食谱列表
│   │   │   ├── RecipeDetail.vue     # 食谱详情页面，显示食谱详细信息
│   │   │   └── RecipeRecommend.vue # 食谱推荐页面，AI推荐食谱
│   │   │
│   │   ├── shopping/     # 购物清单
│   │   │   ├── ShoppingList.vue     # 购物清单列表页面，显示所有购物清单
│   │   │   └── ShoppingDetail.vue   # 购物清单详情页面，显示购物清单详情
│   │   │
│   │   ├── family/       # 家庭管理
│   │   │   ├── FamilyHome.vue       # 家庭首页，家庭信息概览
│   │   │   ├── FamilyMembers.vue    # 家庭成员页面，管理家庭成员
│   │   │   └── FamilySettings.vue   # 家庭设置页面，家庭设置管理
│   │   │
│   │   ├── profile/      # 个人中心
│   │   │   ├── Profile.vue          # 个人资料页面，查看和编辑个人信息
│   │   │   ├── Settings.vue         # 设置页面，个人设置管理
│   │   │   └── About.vue            # 关于页面，应用信息说明
│   │   │
│   │   └── auth/         # 认证相关
│   │       ├── Login.vue            # 登录页面，用户登录
│   │       └── Register.vue        # 注册页面，用户注册
│   │
│   ├── router/            # 路由配置
│   │   ├── index.js      # 路由主文件，定义所有路由规则
│   │   └── guards.js     # 路由守卫，处理路由权限和登录验证
│   │
│   ├── stores/            # 状态管理（Pinia）
│   │   ├── user.js       # 用户状态管理，管理用户信息和登录状态
│   │   ├── family.js     # 家庭状态管理，管理家庭相关状态
│   │   ├── recipe.js     # 食谱状态管理，管理食谱相关状态
│   │   └── shopping.js   # 购物清单状态管理，管理购物清单相关状态
│   │
│   ├── api/               # API接口
│   │   ├── index.js      # API配置文件，HTTP请求基础配置
│   │   ├── user.js       # 用户API接口，用户相关的API调用
│   │   ├── family.js     # 家庭API接口，家庭相关的API调用
│   │   ├── recipe.js     # 食谱API接口，食谱相关的API调用
│   │   ├── ingredient.js # 食材API接口，食材相关的API调用
│   │   └── shopping.js   # 购物清单API接口，购物清单相关的API调用
│   │
│   ├── utils/             # 工具函数
│   │   ├── request.js    # HTTP请求封装，axios封装和拦截器
│   │   ├── auth.js       # 认证工具，token管理和认证相关函数
│   │   ├── storage.js    # 本地存储工具，localStorage和sessionStorage封装
│   │   ├── format.js     # 格式化工具，日期、金额等格式化函数
│   │   ├── validate.js   # 验证工具，表单验证和数据验证函数
│   │   ├── constants.js  # 常量定义，项目常量配置
│   │   └── wechat.js     # 微信小程序工具，微信API封装
│   │
│   ├── styles/            # 样式文件
│   │   ├── variables.scss # SCSS变量文件，定义颜色、尺寸等变量
│   │   ├── mixins.scss   # SCSS混入文件，定义可复用的样式混入
│   │   ├── common.scss   # 通用样式文件，全局通用样式
│   │   └── reset.scss    # 样式重置文件，重置浏览器默认样式
│   │
│   ├── composables/       # 组合式函数
│   │   ├── useAuth.js    # 认证相关组合函数，认证逻辑封装
│   │   ├── useFamily.js  # 家庭相关组合函数，家庭管理逻辑封装
│   │   ├── useRecipe.js  # 食谱相关组合函数，食谱操作逻辑封装
│   │   └── useShopping.js # 购物清单相关组合函数，购物清单逻辑封装
│   │
│   ├── directives/        # 自定义指令
│   │   └── permission.js # 权限指令，控制元素显示权限
│   │
│   ├── plugins/           # 插件
│   │   └── ui-library.js # UI库配置，UI组件库初始化
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
- 小程序中需要注意图片大小优化

### src/components/
组件目录：
- **common/**: 通用UI组件
- **business/**: 业务组件，如食谱卡片、购物清单项等

### src/pages/
页面组件目录：
- 小程序的主要页面
- 按功能模块组织
- 每个页面对应一个路由

### src/router/
路由配置：
- 小程序路由定义
- 路由守卫（登录验证、权限控制）
- 页面跳转配置

### src/stores/
状态管理（Pinia）：
- 全局状态管理
- 用户信息、家庭信息、购物清单等状态
- 支持持久化存储

### src/api/
API接口封装：
- 统一的HTTP请求封装
- 按业务模块组织
- 请求拦截和响应处理
- 错误统一处理

### src/utils/
工具函数：
- HTTP请求工具
- 认证工具（Token管理）
- 本地存储工具
- 格式化工具（日期、金额等）
- 验证工具（表单验证）
- 微信小程序相关工具（如获取用户信息）

### src/styles/
样式文件：
- SCSS变量和混入
- 全局样式
- 小程序样式规范

### src/composables/
组合式函数：
- 可复用的业务逻辑
- 如家庭管理、食谱推荐、购物清单等

### src/directives/
自定义指令：
- Vue自定义指令
- 如权限控制指令

### src/plugins/
插件配置：
- 第三方库的初始化
- UI组件库配置

### public/
公共资源目录：
- 不会被构建工具处理的静态文件

## 文件命名规范

### Vue组件
- 使用PascalCase：`RecipeCard.vue`
- 页面组件：`Home.vue`, `RecipeList.vue`
- 组件名与文件名保持一致

### JavaScript文件
- 使用camelCase：`userService.js`
- 工具文件：`formatDate.js`, `validateForm.js`
- 常量文件：`constants.js`

### 样式文件
- 使用kebab-case：`recipe-card.scss`
- 或使用功能命名：`common.scss`, `variables.scss`

## 小程序特殊说明

### 页面配置
每个页面可能需要独立的配置文件：
- `pages/home/home.json` - 页面配置
- `app.json` - 全局配置

### 小程序API
- 使用微信小程序API时需要封装
- 注意小程序的生命周期
- 处理小程序特有的功能（如分享、扫码等）

### 性能优化
- 图片懒加载
- 列表虚拟滚动
- 代码分包加载
- 减少包体积

## 代码组织原则

1. **组件化开发**：
   - 单一职责原则
   - 可复用性
   - 组件通信规范化

2. **状态管理**：
   - 全局状态使用Pinia
   - 支持持久化（localStorage）
   - 避免过度使用全局状态

3. **API封装**：
   - 统一的请求拦截器
   - Token自动添加
   - 错误统一处理
   - 请求重试机制

4. **路由管理**：
   - 路由懒加载
   - 路由守卫（登录验证）
   - 页面权限控制

5. **样式管理**：
   - 使用SCSS变量统一主题
   - 响应式设计
   - 小程序样式规范

## 开发规范

### 组件结构
```vue
<template>
  <!-- 模板内容 -->
</template>

<script setup>
// 导入依赖
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'

// 定义props
const props = defineProps({
  // ...
})

// 响应式数据
const loading = ref(false)

// 生命周期
onMounted(() => {
  // 初始化逻辑
})

// 方法
const handleAction = () => {
  // ...
}
</script>

<style scoped>
/* 样式 */
</style>
```

### API调用示例
```javascript
// api/recipe.js
import request from '@/utils/request'

export function getRecipeList(params) {
  return request.get('/api/recipes', { params })
}

export function getRecipeById(id) {
  return request.get(`/api/recipes/${id}`)
}

export function recommendRecipe(data) {
  return request.post('/api/recipes/recommend', data)
}
```

### 状态管理示例
```javascript
// stores/family.js
import { defineStore } from 'pinia'

export const useFamilyStore = defineStore('family', {
  state: () => ({
    currentFamily: null,
    members: []
  }),
  
  getters: {
    memberCount: (state) => state.members.length
  },
  
  actions: {
    async loadFamily(id) {
      // 加载家庭信息
    },
    
    async addMember(member) {
      // 添加成员
    }
  },
  
  persist: {
    enabled: true,
    strategies: [
      {
        key: 'family',
        storage: localStorage
      }
    ]
  }
})
```

## 相关文档

- [API接口文档](./接口文档.md)
- [数据库设计文档](./数据库设计.md)
- [部署文档](./部署文档.md)


