---
name: ant-design-vue
description: Use when building, refactoring, or reviewing Vue 3 interfaces with Ant Design Vue 4.x in this repository, especially for layout, spacing, radius, button consistency, forms, tables, navigation, feedback, and enterprise admin page styling decisions.
---

# Ant Design Vue（GoInsight 统一版）

用于当前仓库的前端 UI 规范技能。目标是输出可维护、可复用、视觉一致的企业后台界面方案，并减少一次性样式和组件漂移。

## 适用范围

- 技术栈：Vue 3 + Ant Design Vue 4.x（当前项目使用 4.2.6）。
- 语法约束：使用 `script setup`，不使用 `lang="ts"`。
- 页面范围：登录、列表、详情、编辑、仪表盘、设置等全前端页面。
- 任务类型：新页面设计、旧页面重构、组件替换、样式统一、代码评审建议。

## 触发信号（命中任一即使用本技能）

- 用户提到：`Ant Design Vue`、`AntDV`、`ant-design-vue`。
- 用户要求：按钮/表单/表格/弹窗/菜单/分页等组件实现或改造。
- 用户提到：布局、间距、圆角、主题、视觉统一、组件规范。

## 输出方式（无强制门禁）

默认先给“规范化 UI 方案”，再给实现建议。  
方案中至少包含：

1. 页面结构与布局建议
2. 组件映射建议
3. 与全局规范一致性的检查点

## 全局强制规范

以以下文件作为唯一真源：

- `/Users/zongfei.fu/Downloads/CodeX/goInsight/www/src/design/tokens.js`
- `/Users/zongfei.fu/Downloads/CodeX/goInsight/www/src/design/antdTheme.js`

### 1) 间距规范（强制）

- 仅使用：`4 / 8 / 12 / 16 / 24 / 32`
- 禁止散点间距（如 10、14、18、22 等无来源值）

### 2) 圆角规范（强制）

- 仅使用：`6 / 8 / 10 / 12`
- 禁止组件内手写临时圆角值

### 3) 按钮规范（强制）

- 默认统一：`size="middle"`
- 禁止手动设置按钮颜色、圆角、高度（包含 inline style）
- 仅允许使用 AntDV 语义属性控制按钮表现：`type`、`danger`、`ghost`、`block`、`size`
- 登录页按钮也遵守统一 `middle`

## 布局规范

推荐后台壳层结构：

```
Layout
 |- Sider (Menu)
 |- Header
 `- Content
```

- 使用 `Layout` / `Layout.Sider` / `Layout.Header` / `Layout.Content`
- 侧边栏宽度建议 `200px ~ 240px`
- 主内容区优先用 `a-card` 承载
- 页面内边距建议 `16px ~ 24px`
- 使用 24 栅格：`a-row` + `a-col`
- 栅格 `gutter` 优先 `16` 或 `24`

## 语义颜色规范

优先使用 Ant Design 语义色，不随意造色。基础语义：

- Primary / Info：`#1677ff`
- Success：`#52c41a`
- Warning：`#faad14`
- Error：`#ff4d4f`
- Border：`#d9d9d9`
- 主文本：`rgba(0,0,0,0.88)`
- 次文本：`rgba(0,0,0,0.65)`

规则：

- 重要动作用主色
- 状态展示优先 `a-tag` 语义色
- 禁止页面内零散定义“看起来差不多”的新颜色

## 组件规范

### 按钮

- 主操作：`a-button type="primary"`
- 次操作：默认 `a-button`
- 危险操作：`a-button danger`
- 同一视觉区域最多两个主按钮

### 表单

- 统一使用 `a-form`
- 必须有 label 和基础校验规则
- 标签对齐方式同页保持一致
- 复杂表单优先纵向布局
- 输入组件优先：`a-input` / `a-select` / `a-date-picker` / `a-switch` / `a-radio-group` / `a-checkbox-group`

### 表格

- 统一使用 `a-table`
- 必须设置 `rowKey`
- 默认启用分页
- 列信息聚焦，避免密度过高
- 排序仅在显著提升任务效率时启用

### 导航与信息组织

- 导航：`a-menu` / `a-breadcrumb` / `a-tabs`
- 分组：`a-card`
- 操作组：`a-dropdown`
- 详情展示：`a-descriptions`

### 反馈

- 轻提示：`message.success()` 等
- 显著错误：`notification.error()`
- 不可逆动作：`modal.confirm()`
- 加载反馈：`a-spin`

## 页面蓝图

### 列表页

1. 顶部筛选区（`a-form`）
2. 中部数据区（`a-table`）
3. 底部分页（`a-pagination` 或表格内建）

### 详情页

1. 核心信息区用 `a-descriptions`
2. 自然分区内容再使用 `a-tabs`

### 编辑页

1. 以表单为主结构
2. 提交与取消按钮放在清晰的操作行

### 仪表盘

1. 关键指标用 `a-statistic`
2. 分模块卡片化展示

### 设置页

1. 按领域拆分 `a-tabs`
2. 每个 tab 内聚焦一组表单能力

## 交互位置约定

- 全局主操作：右上或标题区右侧
- 筛选/检索：页面顶部
- 主数据区：中部连续区域
- 同类按钮：紧邻分组并有明确优先级

## 常见反模式与纠偏

| 反模式 | 纠偏方式 |
| --- | --- |
| 按钮写 inline style 改颜色和圆角 | 改用全局 token + AntDV 语义属性 |
| 页面出现 14px/18px 这类散点间距 | 回退到 4/8/12/16/24/32 |
| 组件圆角写成 5px/9px/14px | 回退到 6/8/10/12 |
| 登录页按钮单独用 `large` | 统一为 `middle` |
| 同页混用多个组件库 | 优先统一为 AntDV 组件 |
| 使用 `<script setup lang="ts">` | 改为 JS 版 `script setup` |

## 快速检查清单

- 是否仅使用项目 token 定义的间距与圆角？
- 按钮是否全部 `size="middle"`？
- 是否存在按钮手写视觉样式（颜色/高度/圆角/inline style）？
- 页面是否使用 AntDV 语义组件而不是自造同类组件？
- 是否保持布局层级清晰、操作优先级明确？
