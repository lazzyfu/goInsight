# GoInsight UI Revamp Phase 1 Design

Date: 2026-03-03
Owner: Codex + @fuzongfei
Status: Approved (brainstorming validated)

## 1. 背景与问题

当前前端基于 Vue3 + Ant Design Vue，已经具备基础可用性，但存在以下持续性问题：

- 样式来源分散：`scoped` 样式、行内样式、历史 SCSS 并存，缺少统一入口。
- 排版与间距不统一：页面节奏不稳定，不同模块视觉密度差异大。
- 视觉层级弱：标题、操作区、信息区、反馈区权重不稳定。
- 响应式策略缺位：移动端与平板适配主要依赖自然换行，细节体验不足。

## 2. 目标与范围

## 2.1 目标

- 建立可复用、可演进的全站 UI 基础规范（Token + Layout + Page Pattern）。
- 第一阶段完成高频路径页面改造并沉淀范式，作为后续模块迁移基线。
- 达成全端可用且细节一致（移动端/平板/桌面/宽屏）。

## 2.2 第一阶段覆盖范围（已确认）

- 登录页
- 主框架（侧边栏、顶栏、内容容器）
- 工单列表
- 工单详情

## 2.3 改造深度（已确认）

- 中改：允许统一 Token、组件样式与局部信息布局调整。
- 允许替换页面内行内样式与硬编码样式。

## 2.4 非目标（Phase 1）

- 不在本阶段重构所有 admin/das 子页面。
- 不引入全新 UI 组件库，不替换 Ant Design Vue。
- 不做业务接口语义变更。

## 3. 视觉方向（已确认）

- 风格：现代轻量（Modern Lightweight）
- 关键词：清爽留白、信息清晰、可控层级、弱装饰强可读
- 设计约束：
  - 保持企业后台可信感，不做炫技动效
  - 优先保证信息效率和可维护性

## 4. 方案评估与选型

已评估三种路径：

1. Token 先行 + 外壳统一 + 页面对齐（推荐）
2. 组件层先行 + 页面渐进替换
3. 页面点修 + 最小全局约束

最终选择方案 1，原因：

- 一致性收益最大，避免后续反复返工。
- 第一阶段产出可直接作为 admin/das 迁移模板。
- 能同时解决视觉、结构、响应式三类问题。

## 5. 信息架构与样式治理

采用四层结构：

`Design Token -> Layout Shell -> Page Pattern -> Business Module`

## 5.1 Design Token 层

建立语义化变量集（优先 CSS 变量，必要时映射到 Antd theme token）：

- Color: brand/text/bg/border/status
- Spacing: 4/8/12/16/24/32
- Radius: 6/8/10/12
- Shadow: sm/md/lg
- Typography: 12/14/16/20/24 + 行高规范
- Breakpoint: mobile/tablet/desktop/wide
- Motion: 120/200/280ms

约束：

- 页面样式必须优先使用 Token，不允许散点色值与散点间距。
- 行内样式仅允许动态计算场景，静态样式必须迁移到样式文件。

## 5.2 Layout Shell 层

统一后台外壳行为与样式：

- 侧栏：宽度、菜单项高度、active/hover 语义统一
- 顶栏：高度、边界线、面包屑与用户区对齐规则统一
- 内容区：背景层与容器层分离，避免贴边白屏感

登录页独立壳层，但共享同一 Token 体系。

## 5.3 Page Pattern 层

定义页面骨架模式：

- `page-shell`
- `page-header`
- `page-toolbar`
- `page-section`
- `page-grid`

统一卡片、工具条、数据区、空态区的节奏与间距。

## 5.4 Business Module 层

业务页面仅实现业务逻辑和数据映射，尽量复用 Page Pattern 与语义类。

## 6. 三个先交付件（ant-design-vue guardrail）

## 6.1 Style Brief

- Brand tone: calm and efficient
- Palette: neutral base + restrained teal/blue accents
- Density: medium density with clearer whitespace rhythm
- Priority: readability > decoration > novelty

## 6.2 Component Mapping Table

| 页面 | 区块 | 组件模式 | 备注 |
| --- | --- | --- | --- |
| Login | 表单区 | AuthCard / FormField / PrimaryAction | OTP 流程与默认登录统一节奏 |
| Layout | 导航区 | AppSider / AppHeader / ContentShell | 折叠与移动端抽屉行为一致 |
| OrderList | 筛选区 | FilterBar | 控件高度、间距、折行规则统一 |
| OrderList | 概览区 | MetricCard | 状态数字与标签语义色统一 |
| OrderList | 数据区 | DataTableShell | 表格 hover/空态/滚动一致 |
| OrderDetail | 标题区 | DetailHero | 关键信息与操作按钮分层 |
| OrderDetail | 审批/日志/内容 | InfoSection | 标题、内边距、边框样式一致 |

## 6.3 UI Acceptance Checklist

- 间距是否只使用刻度值（4/8/12/16/24/32）？
- 字体层级是否只使用规范档位（12/14/16/20/24）？
- 状态色是否来源于统一状态映射？
- 目标页面是否移除硬编码行内样式？
- mobile/tablet/desktop/wide 四断点是否都可用且无布局破坏？
- 关键路径（登录、工单筛选、详情操作）是否流畅？

## 7. 页面级设计细节

## 7.1 登录页

- Desktop: 双栏（视觉叙事 + 表单聚焦）
- Tablet: 上下布局，视觉区缩高
- Mobile: 单栏聚焦表单，保证触控友好尺寸

交互要求：

- OTP 状态切换清晰可见
- 主按钮强视觉锚点
- 入场动效轻量且可降级

## 7.2 主框架

- 侧栏折叠态图标对齐统一
- 顶栏高度与间距固定化
- 内容区默认使用容器化布局，减少贴边阅读负担

## 7.3 工单列表

- 标题栏 -> 筛选条 -> 概览卡 -> 表格 的固定结构
- `FilterBar` 统一控件高与折行规则
- `MetricCard` 统一数字和标签视觉优先级
- 表格状态标签统一映射

移动端策略：

- 筛选区分行
- 概览卡 1-2 列自适应
- 表格关键列优先 + 横向滚动

## 7.4 工单详情

- Header 区强调核心状态和主要动作
- 审批流、日志、SQL 内容采用一致 Section 视觉
- 操作按钮主次分明，危险操作保持二次确认

移动端策略：

- Header 改纵向信息流
- 操作按钮分组展示，避免一行拥挤

## 8. 数据流与状态映射

- 页面状态拆分：query state / ui state / domain state
- 统一 `statusMap`：集中定义 text/color/priority/icon
- 页面消费语义状态，不直接写具体色值

## 9. 错误处理与空态规范

- 网络错误：统一提示 + 区块重试入口
- 空数据：统一空态文案和视觉模板
- 长文本：统一截断与换行策略，避免破坏布局

## 10. 响应式与可访问性

断点建议：

- mobile: <= 767
- tablet: 768 - 1023
- desktop: 1024 - 1439
- wide: >= 1440

可用性要求：

- 可点击区符合触屏最小尺寸
- 文字与背景对比满足可读性
- 焦点态可见，键盘导航不丢失

## 11. 验证与验收

验证维度：

- 视觉回归（四断点截图对比）
- 交互回归（登录、筛选、分页、详情操作）
- 样式审计（行内样式与硬编码清理）
- 性能与稳定性（避免高成本特效）

Phase 1 通过标准：

- 目标页面达到统一样式节奏
- 关键路径体验可用且一致
- 为 admin/das 提供可直接复用的模板

## 12. 风险与缓解

- 风险：历史样式耦合导致局部回归
  - 缓解：分页面渐进替换 + 快速回归清单
- 风险：业务逻辑与样式重构交织
  - 缓解：严格限定为 UI 重构，逻辑最小扰动
- 风险：移动端细节遗漏
  - 缓解：强制四断点验收和关键路径真机检查

## 13. 交付清单（Phase 1）

- 全局 Token 规范与落地
- Layout Shell 统一改造
- 登录页、工单列表、工单详情视觉重构
- 响应式适配（mobile/tablet/desktop/wide）
- UI 验收清单执行记录

## 14. 执行证据（2026-03-03）

已完成的关键实现提交（`codex/ui-revamp-phase1`）：

- `54b65594` style: add global token and base stylesheet entry
- `c2c811b5` style: bridge design tokens into antd theme
- `34ab954d` refactor: centralize order status metadata
- `8fcd084b` style: remove static inline styles from phase1 pages
- `e2cb302e` feat: unify responsive layout shell pattern
- `fc1ee6be` refactor: apply page pattern to order list
- `706ca0e6` refactor: modernize order detail structure and sections
- `aa4e5f8c` refactor: polish responsive login with tokenized styles

验证命令（已执行）：

- `cd www && npm run lint` -> PASS
- `cd www && npm run build` -> PASS（含历史遗留告警：Sass @import deprecation、chunk size）

详细检查记录见：`docs/plans/2026-03-03-ui-revamp-phase1-qa-checklist.md`
