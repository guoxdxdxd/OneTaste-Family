# Frontend Miniapp 目录结构

本文档依据当前仓库 `src/frontend/miniapp` 的真实文件生成，完整列出项目包含的目录与文件，并给出简要功能说明，便于定位代码位置和职责。

## 设计系统说明

本项目采用 **温暖精致 (Warm & Refined)** 的设计风格：
- **主色调**: 暖橙色系 (#f08a5d)
- **设计理念**: 温暖但专业，适合家庭场景
- **组件规范**: BEM 命名规范，CSS 变量驱动

## 目录结构（含注释）

```
src/frontend/miniapp/
├── .env                              # 本地开发环境变量（不会提交）
├── .env.example                      # 环境变量示例
├── .gitignore                        # Git 忽略规则
├── README.md                         # 项目说明与运行指南
├── index.html                        # Vite 入口 HTML 模板
├── package.json                      # 项目依赖与 npm scripts
├── package-lock.json                 # 依赖版本锁定
├── vite.config.js                    # Vite 构建配置
├── node_modules/                     # npm 依赖目录
└── src/                              # 源代码目录
    ├── App.vue                       # 根组件
    ├── main.js                       # 应用入口，挂载 Vue 实例
    │
    ├── api/                          # HTTP API 封装
    │   ├── index.js                  # axios 实例与拦截器
    │   ├── user.js                   # 用户相关接口
    │   ├── dishes.js                 # 菜式相关接口
    │   └── menus.js                  # 菜单相关接口
    │
    ├── components/                   # 通用组件库
    │   ├── TabBar.vue                # 底部导航栏组件
    │   ├── Calendar.vue              # 日历选择组件
    │   ├── MenuCard.vue              # 菜单卡片组件
    │   ├── IngredientSelector.vue    # 食材选择器组件
    │   └── icons/                    # SVG 图标组件
    │       ├── IconMenu.vue          # 菜单图标
    │       ├── IconShopping.vue      # 购物图标
    │       ├── IconUser.vue          # 用户图标
    │       ├── IconPlus.vue          # 加号图标
    │       ├── IconCheck.vue         # 勾选图标
    │       ├── IconClose.vue         # 关闭图标
    │       ├── IconChevronLeft.vue   # 左箭头图标
    │       ├── IconChevronRight.vue  # 右箭头图标
    │       ├── IconCalendar.vue      # 日历图标
    │       ├── IconBook.vue          # 书本图标
    │       ├── IconFamily.vue        # 家庭图标
    │       ├── IconSettings.vue      # 设置图标
    │       └── IconCook.vue          # 烹饪图标
    │
    ├── pages/                        # 业务页面
    │   ├── auth/                     # 认证页面
    │   │   ├── Login.vue             # 登录页面
    │   │   └── Register.vue          # 注册页面
    │   ├── layouts/                  # 布局组件
    │   │   └── AppLayout.vue         # 主布局（含底部 TabBar）
    │   ├── menu/                     # 菜单模块
    │   │   ├── MenuHub.vue           # 菜单中心首页
    │   │   ├── CreateMenu.vue        # 创建菜单页面
    │   │   ├── EditMenu.vue          # 编辑菜单页面
    │   │   ├── DailyMenu.vue         # 每日菜单页面
    │   │   └── WeeklyMenu.vue        # 每周菜单页面
    │   ├── shopping/                 # 购物清单模块
    │   │   └── ShoppingHub.vue       # 购物清单页面
    │   ├── profile/                  # 个人中心模块
    │   │   └── ProfileHome.vue       # 个人中心页面
    │   ├── recipes/                  # 菜谱管理模块
    │   │   └── RecipeManagement.vue  # 菜谱管理页面
    │   └── invite/                   # 邀请模块
    │       └── InviteLanding.vue     # 邀请落地页
    │
    ├── router/                       # Vue Router 配置
    │   ├── index.js                  # 路由表与实例化
    │   └── guards.js                 # 路由守卫
    │
    ├── stores/                       # Pinia 状态管理
    │   ├── user.js                   # 用户状态
    │   └── family.js                 # 家庭状态
    │
    ├── styles/                       # 全局样式
    │   ├── theme.css                 # 设计系统变量（颜色、阴影、圆角等）
    │   ├── base.css                  # 基础样式重置与全局样式
    │   └── components.css            # 通用组件样式类
    │
    └── utils/                        # 工具函数
        ├── auth.js                   # token 与用户缓存管理
        ├── constants.js              # 常量定义
        ├── request.js                # axios 请求封装
        ├── storage.js                # 本地存储封装
        └── validate.js               # 表单验证函数
```

## 样式系统说明

### theme.css - 设计系统变量

包含以下变量类别：
- **色彩系统**: 主色、辅助色、功能色、中性色、背景色、文字色
- **渐变系统**: 主色渐变、背景渐变
- **阴影系统**: 多级阴影、卡片阴影、按钮阴影
- **圆角系统**: xs/sm/md/lg/xl/2xl/3xl/full
- **间距系统**: 基于 4px 的间距体系
- **字体系统**: 字体族、字号、字重、行高
- **动效系统**: 时长、缓动函数、组合过渡
- **布局系统**: 容器宽度、导航栏高度、安全区域

### base.css - 基础样式

- CSS Reset
- 全局元素样式（标题、段落、链接等）
- 滚动条美化
- 选中状态样式
- 动画关键帧定义

### components.css - 组件样式库

提供以下组件类：
- **布局**: `.page`, `.page-header`, `.page-section`
- **卡片**: `.card`, `.card--flat`, `.card--highlight`
- **按钮**: `.btn`, `.btn--primary`, `.btn--ghost`, `.btn--danger`
- **表单**: `.input`, `.textarea`, `.select`, `.form-group`
- **标签**: `.tag`, `.tag--primary`, `.tag--success`
- **头像**: `.avatar`, `.avatar-group`
- **列表**: `.list-item`
- **空状态**: `.empty-state`
- **加载**: `.loading-spinner`, `.skeleton`
- **底部弹出**: `.bottom-sheet`
- **工具类**: `.flex`, `.gap-*`, `.mb-*` 等

## 页面说明

### 认证页面
- **Login.vue**: 手机号密码登录，带装饰背景和品牌展示
- **Register.vue**: 用户注册，风格与登录页保持一致

### 菜单模块
- **MenuHub.vue**: 首页，展示今日菜单、快捷操作、日历选择
- **CreateMenu.vue**: 创建新菜单
- **EditMenu.vue**: 编辑已有菜单
- **DailyMenu.vue**: 每日菜单视图
- **WeeklyMenu.vue**: 每周菜单视图

### 购物模块
- **ShoppingHub.vue**: 购物清单，按分类展示待购食材

### 个人中心
- **ProfileHome.vue**: 用户信息、家庭管理、成员列表、邀请功能

### 菜谱管理
- **RecipeManagement.vue**: 菜式 CRUD、搜索筛选、抽屉式编辑

## 组件说明

### TabBar.vue
底部导航栏组件，支持：
- 三个 Tab（菜单、买菜、我的）
- 激活状态动画
- 毛玻璃背景效果
- 徽章显示

### Calendar.vue
日历组件，支持：
- 月份切换（限制前后3个月）
- 日期选择
- 有菜单日期标记
- 今日高亮

### IngredientSelector.vue
食材选择器，支持：
- 食材搜索
- 分类浏览
- 数量单位输入
- 多选管理

## 代码规范

### 命名规范
- 组件文件：PascalCase（如 `MenuHub.vue`）
- 样式类：BEM 命名（如 `.card__header`, `.btn--primary`）
- CSS 变量：`--color-*`, `--space-*`, `--radius-*` 等

### 组件开发规范
1. 使用 `<script setup>` 语法
2. Props 必须声明类型和默认值
3. 复杂逻辑添加注释说明
4. 样式使用 scoped，通过 CSS 变量引用设计系统

### 样式开发规范
1. 优先使用 `components.css` 中的通用类
2. 颜色、间距等必须使用 CSS 变量
3. 组件内样式使用 scoped
4. 响应式断点保持一致

## 快速开始

```bash
# 安装依赖
npm install

# 开发模式
npm run dev

# 构建生产版本
npm run build

# 预览构建结果
npm run preview
```

## 更新日志

### v2.0.0 (2026-01-31)
- 全面重构 UI 设计，采用温暖精致风格
- 新增设计系统（theme.css）
- 重构所有页面组件
- 新增图标组件库
- 优化组件样式复用

---

当新增或调整目录结构时，请同步更新本文件，确保文档与代码保持一致。
