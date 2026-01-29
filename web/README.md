# 衣搭库 - Vue 前端

基于《衣搭库产品需求文档》与原有静态原型（`index.html` / `css` / `js`）迁移的 **Vue 3** 前端项目。

## 技术栈

- **Vue 3**（Composition API + `<script setup>`）
- **Vue Router 4**（Hash 模式）
- **Vite 5**

## 项目结构

```
web/
├── index.html          # Vite 入口
├── package.json
├── vite.config.js
├── src/
│   ├── main.js         # 应用入口
│   ├── App.vue         # 根组件（头部 + 主内容 + 底部导航）
│   ├── router/         # 路由配置
│   ├── views/          # 页面：首页、衣橱、推荐、上传、收藏、设置
│   ├── components/     # 公共组件：AppHeader、BottomNav
│   ├── composables/    # 逻辑复用：useUploadCache（IndexedDB 图片缓存）
│   ├── utils/          # 工具：toast
│   └── styles/         # 全局样式（沿用原 UI 设计规范）
└── public/
```

## 开发与构建

```bash
# 安装依赖
npm install

# 本地开发（默认 http://localhost:5173）
npm run dev

# 构建生产
npm run build

# 预览构建结果
npm run preview
```

## 包含的视图

1. **首页**：天气卡片、最近上传、每日推荐入口、快捷场合 Chip
2. **衣橱**：筛选栏、网格/列表切换、衣物卡片、上传入口
3. **推荐**：场合 + 天气选择、穿搭卡片、收藏
4. **上传**：点击/拖拽上传、本地缓存（IndexedDB）、已缓存列表、识别结果占位
5. **收藏**：空状态占位
6. **设置**：占位页

## 设计说明

视觉与交互规范见项目根目录 `docs/UI设计说明.md`。样式变量与组件类名与原静态原型一致，便于后续对接接口与扩展功能。
