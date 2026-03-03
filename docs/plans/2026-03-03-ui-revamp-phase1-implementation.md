# UI Revamp Phase 1 Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Deliver a modern-lightweight, fully responsive UI baseline for Login, Layout shell, Order List, and Order Detail while establishing reusable design tokens and page patterns.

**Architecture:** Build a token-first UI foundation (`design tokens -> theme bridge -> shell patterns -> page modules`) and migrate the phase-1 pages to consume semantic classes/utilities instead of inline or hardcoded styles. Keep business API behavior unchanged, and use contract tests for token consistency, status mapping, and inline-style guardrails.

**Tech Stack:** Vue 3 (`script setup`), Ant Design Vue 4, Vite 7, SCSS, Vitest + Vue Test Utils + jsdom, ESLint.

---

Implementation skills to apply during execution: `@test-driven-development`, `@ant-design-vue-ui`, `@verification-before-completion`.

### Task 1: Add Frontend Test Harness and Token Contracts

**Files:**
- Modify: `www/package.json`
- Create: `www/vitest.config.js`
- Create: `www/src/test/setup.js`
- Create: `www/src/design/tokens.js`
- Test: `www/src/design/__tests__/tokens.spec.js`

**Step 1: Write the failing test**

```js
import { describe, expect, it } from 'vitest'
import { breakpoints, radiusScale, spacingScale, typographyScale } from '../tokens'

describe('design tokens', () => {
  it('exports spacing scale in 8pt rhythm', () => {
    expect(spacingScale.md).toBe(16)
    expect(spacingScale.xl).toBe(32)
  })

  it('exports responsive breakpoints', () => {
    expect(breakpoints.mobile).toBe(767)
    expect(breakpoints.desktop).toBe(1024)
  })

  it('exports typography and radius scales', () => {
    expect(typographyScale.body).toBe(14)
    expect(radiusScale.card).toBe(12)
  })
})
```

**Step 2: Run test to verify it fails**

Run: `cd www && npm run test:unit -- src/design/__tests__/tokens.spec.js`

Expected: FAIL with module/script errors (`test:unit` missing and/or `../tokens` missing).

**Step 3: Write minimal implementation**

```js
// www/src/design/tokens.js
export const spacingScale = { xs: 4, sm: 8, md: 16, lg: 24, xl: 32 }
export const radiusScale = { sm: 6, md: 8, lg: 10, card: 12 }
export const typographyScale = { caption: 12, body: 14, title: 20, hero: 24 }
export const breakpoints = { mobile: 767, tablet: 768, desktop: 1024, wide: 1440 }
```

Also add:
- `test:unit` script and Vitest dev dependencies in `www/package.json`
- base Vitest config in `www/vitest.config.js`
- `www/src/test/setup.js`

**Step 4: Run test to verify it passes**

Run: `cd www && npm run test:unit -- src/design/__tests__/tokens.spec.js`

Expected: PASS.

**Step 5: Commit**

```bash
git add www/package.json www/vitest.config.js www/src/test/setup.js www/src/design/tokens.js www/src/design/__tests__/tokens.spec.js
git commit -m "test: add frontend token contracts with vitest"
```

### Task 2: Build Global SCSS Token Layer and App-Wide Style Entry

**Files:**
- Create: `www/src/assets/scss/tokens.scss`
- Create: `www/src/assets/scss/base.scss`
- Modify: `www/src/assets/scss/index.scss`
- Modify: `www/src/main.js`
- Test: `www/src/design/__tests__/css-vars.spec.js`

**Step 1: Write the failing test**

```js
import fs from 'node:fs'
import path from 'node:path'
import { describe, expect, it } from 'vitest'

describe('css variable contract', () => {
  it('defines phase-1 required variables', () => {
    const file = path.resolve('src/assets/scss/tokens.scss')
    const content = fs.readFileSync(file, 'utf-8')
    expect(content).toContain('--gi-spacing-md: 16px;')
    expect(content).toContain('--gi-color-primary: #0f766e;')
    expect(content).toContain('--gi-radius-card: 12px;')
  })
})
```

**Step 2: Run test to verify it fails**

Run: `cd www && npm run test:unit -- src/design/__tests__/css-vars.spec.js`

Expected: FAIL because `tokens.scss` does not exist yet.

**Step 3: Write minimal implementation**

```scss
/* www/src/assets/scss/tokens.scss */
:root {
  --gi-color-primary: #0f766e;
  --gi-color-bg-page: #f4f7f8;
  --gi-spacing-sm: 8px;
  --gi-spacing-md: 16px;
  --gi-spacing-lg: 24px;
  --gi-radius-card: 12px;
}
```

```scss
/* www/src/assets/scss/index.scss */
@import 'tokens', 'base', 'layout', 'form', 'flow', 'screen';
```

```js
// www/src/main.js
import '@/assets/scss/index.scss'
```

**Step 4: Run test to verify it passes**

Run: `cd www && npm run test:unit -- src/design/__tests__/css-vars.spec.js`

Expected: PASS.

**Step 5: Commit**

```bash
git add www/src/assets/scss/tokens.scss www/src/assets/scss/base.scss www/src/assets/scss/index.scss www/src/main.js www/src/design/__tests__/css-vars.spec.js
git commit -m "style: add global token and base stylesheet entry"
```

### Task 3: Add Ant Design Theme Bridge from Shared Tokens

**Files:**
- Create: `www/src/design/antdTheme.js`
- Modify: `www/src/App.vue`
- Test: `www/src/design/__tests__/antd-theme.spec.js`

**Step 1: Write the failing test**

```js
import { describe, expect, it } from 'vitest'
import { antdTheme } from '../antdTheme'

describe('antd theme bridge', () => {
  it('maps shared tokens to antd component tokens', () => {
    expect(antdTheme.token.colorPrimary).toBe('#0f766e')
    expect(antdTheme.token.borderRadius).toBe(10)
    expect(antdTheme.components.Card.borderRadiusLG).toBe(12)
  })
})
```

**Step 2: Run test to verify it fails**

Run: `cd www && npm run test:unit -- src/design/__tests__/antd-theme.spec.js`

Expected: FAIL because `antdTheme.js` does not exist.

**Step 3: Write minimal implementation**

```js
// www/src/design/antdTheme.js
export const antdTheme = {
  token: {
    colorPrimary: '#0f766e',
    borderRadius: 10,
    fontSize: 14,
  },
  components: {
    Card: { borderRadiusLG: 12 },
    Button: { borderRadius: 10, controlHeight: 40 },
    Input: { controlHeight: 40 },
  },
}
```

Then wrap root router output in `a-config-provider` in `www/src/App.vue`, using `:theme="antdTheme"`.

**Step 4: Run test to verify it passes**

Run: `cd www && npm run test:unit -- src/design/__tests__/antd-theme.spec.js`

Expected: PASS.

**Step 5: Commit**

```bash
git add www/src/design/antdTheme.js www/src/App.vue www/src/design/__tests__/antd-theme.spec.js
git commit -m "style: bridge design tokens into antd theme"
```

### Task 4: Centralize Order Status Metadata

**Files:**
- Create: `www/src/views/orders/shared/orderStatusMeta.js`
- Modify: `www/src/views/orders/list/OrderList.vue`
- Modify: `www/src/views/orders/detail/OrderDetail.vue`
- Modify: `www/src/views/orders/tasks/TaskList.vue`
- Test: `www/src/views/orders/shared/__tests__/order-status-meta.spec.js`

**Step 1: Write the failing test**

```js
import { describe, expect, it } from 'vitest'
import { getOrderStatusMeta } from '../orderStatusMeta'

describe('order status meta', () => {
  it('returns consistent text and color for reviewed status', () => {
    expect(getOrderStatusMeta('REVIEWED')).toEqual({
      text: '已复核',
      color: 'green',
    })
  })

  it('falls back to raw text for unknown status', () => {
    expect(getOrderStatusMeta('UNKNOWN').text).toBe('UNKNOWN')
  })
})
```

**Step 2: Run test to verify it fails**

Run: `cd www && npm run test:unit -- src/views/orders/shared/__tests__/order-status-meta.spec.js`

Expected: FAIL because `orderStatusMeta.js` does not exist.

**Step 3: Write minimal implementation**

```js
// www/src/views/orders/shared/orderStatusMeta.js
const STATUS_META = {
  PENDING: { text: '待审批', color: 'gold' },
  APPROVED: { text: '已批准', color: 'blue' },
  REJECTED: { text: '已驳回', color: 'red' },
  CLAIMED: { text: '已认领', color: 'cyan' },
  EXECUTING: { text: '执行中', color: 'orange' },
  COMPLETED: { text: '已完成', color: 'green' },
  FAILED: { text: '已失败', color: 'volcano' },
  REVIEWED: { text: '已复核', color: 'green' },
  REVOKED: { text: '已撤销', color: 'default' },
}

export const getOrderStatusMeta = (status) => STATUS_META[status] || { text: status, color: 'default' }
export const ORDER_PROGRESS_OPTIONS = Object.keys(STATUS_META)
```

Refactor affected views to import this shared map and remove duplicated inline maps.

**Step 4: Run test to verify it passes**

Run: `cd www && npm run test:unit -- src/views/orders/shared/__tests__/order-status-meta.spec.js`

Expected: PASS.

**Step 5: Commit**

```bash
git add www/src/views/orders/shared/orderStatusMeta.js www/src/views/orders/shared/__tests__/order-status-meta.spec.js www/src/views/orders/list/OrderList.vue www/src/views/orders/detail/OrderDetail.vue www/src/views/orders/tasks/TaskList.vue
git commit -m "refactor: centralize order status metadata"
```

### Task 5: Add Inline-Style Guard for Phase-1 Pages

**Files:**
- Create: `www/src/design/__tests__/inline-style-guard.spec.js`
- Modify: `www/src/components/layout/Layout.vue`
- Modify: `www/src/views/orders/detail/OrderDetail.vue`
- Modify: `www/src/views/orders/detail/HeaderExtra.vue`
- Modify: `www/src/views/orders/list/OrderList.vue`

**Step 1: Write the failing test**

```js
import fs from 'node:fs'
import path from 'node:path'
import { describe, expect, it } from 'vitest'

const files = [
  'src/components/layout/Layout.vue',
  'src/views/orders/list/OrderList.vue',
  'src/views/orders/detail/OrderDetail.vue',
  'src/views/orders/detail/HeaderExtra.vue',
]

describe('phase-1 files should avoid static inline styles', () => {
  it('has no style=\"...\" attributes', () => {
    for (const file of files) {
      const content = fs.readFileSync(path.resolve(file), 'utf-8')
      expect(content).not.toMatch(/\sstyle=\"[^\"]+\"/g)
    }
  })
})
```

**Step 2: Run test to verify it fails**

Run: `cd www && npm run test:unit -- src/design/__tests__/inline-style-guard.spec.js`

Expected: FAIL because `OrderDetail.vue` and/or `HeaderExtra.vue` currently contain static inline style attributes.

**Step 3: Write minimal implementation**

- Replace static `style="..."` blocks with semantic classes in template.
- Move corresponding rules into scoped SCSS sections (or shared SCSS module if reused).
- Keep dynamic `:style` only where data-driven rendering is required.

Example conversion:

```vue
<!-- before -->
<a-card size="small" title="审批流" style="margin-top: 12px">

<!-- after -->
<a-card size="small" title="审批流" class="detail-section-card">
```

```scss
.detail-section-card {
  margin-top: var(--gi-spacing-sm);
}
```

**Step 4: Run test to verify it passes**

Run: `cd www && npm run test:unit -- src/design/__tests__/inline-style-guard.spec.js`

Expected: PASS.

**Step 5: Commit**

```bash
git add www/src/design/__tests__/inline-style-guard.spec.js www/src/components/layout/Layout.vue www/src/views/orders/detail/OrderDetail.vue www/src/views/orders/detail/HeaderExtra.vue www/src/views/orders/list/OrderList.vue
git commit -m "style: remove static inline styles from phase1 pages"
```

### Task 6: Refactor Layout Shell Into Reusable Responsive Pattern

**Files:**
- Create: `www/src/components/layout/layoutConfig.js`
- Modify: `www/src/components/layout/Layout.vue`
- Test: `www/src/components/layout/__tests__/layout-config.spec.js`

**Step 1: Write the failing test**

```js
import { describe, expect, it } from 'vitest'
import { layoutConfig } from '../layoutConfig'

describe('layout config', () => {
  it('defines consistent shell dimensions', () => {
    expect(layoutConfig.headerHeight).toBe(52)
    expect(layoutConfig.sidebarExpandedWidth).toBe(224)
    expect(layoutConfig.contentPadding.desktop).toBe(24)
    expect(layoutConfig.contentPadding.mobile).toBe(12)
  })
})
```

**Step 2: Run test to verify it fails**

Run: `cd www && npm run test:unit -- src/components/layout/__tests__/layout-config.spec.js`

Expected: FAIL because `layoutConfig.js` does not exist.

**Step 3: Write minimal implementation**

```js
// www/src/components/layout/layoutConfig.js
export const layoutConfig = {
  headerHeight: 52,
  sidebarExpandedWidth: 224,
  sidebarCollapsedWidth: 72,
  contentPadding: { wide: 32, desktop: 24, tablet: 16, mobile: 12 },
}
```

Then update `Layout.vue` to consume these values via CSS vars/classes and unify:
- sider width and menu rhythm
- header structure and spacing
- content shell background + container separation

**Step 4: Run test to verify it passes**

Run: `cd www && npm run test:unit -- src/components/layout/__tests__/layout-config.spec.js`

Expected: PASS.

**Step 5: Commit**

```bash
git add www/src/components/layout/layoutConfig.js www/src/components/layout/__tests__/layout-config.spec.js www/src/components/layout/Layout.vue
git commit -m "feat: unify responsive layout shell pattern"
```

### Task 7: Refactor Order List to Page Pattern + View Model

**Files:**
- Create: `www/src/views/orders/list/orderListModel.js`
- Modify: `www/src/views/orders/list/OrderList.vue`
- Test: `www/src/views/orders/list/__tests__/order-list-model.spec.js`

**Step 1: Write the failing test**

```js
import { describe, expect, it } from 'vitest'
import { buildOrderQuery, summarizeMyOrders } from '../orderListModel'

describe('order list model', () => {
  it('builds api query from ui state', () => {
    expect(buildOrderQuery({ page: 2, pageSize: 20, search: 'ddl', progress: 'PENDING', onlyMine: true })).toEqual({
      page: 2,
      page_size: 20,
      is_page: true,
      search: 'ddl',
      progress: 'PENDING',
      only_my_orders: true,
    })
  })

  it('summarizes my orders', () => {
    const rows = [{ progress: 'PENDING' }, { progress: 'EXECUTING' }, { progress: 'FAILED' }]
    expect(summarizeMyOrders(rows)).toEqual({ total: 3, pending: 1, executing: 1, failed: 1 })
  })
})
```

**Step 2: Run test to verify it fails**

Run: `cd www && npm run test:unit -- src/views/orders/list/__tests__/order-list-model.spec.js`

Expected: FAIL because `orderListModel.js` does not exist.

**Step 3: Write minimal implementation**

```js
// www/src/views/orders/list/orderListModel.js
export const buildOrderQuery = ({ page, pageSize, search, progress, onlyMine }) => ({
  page,
  page_size: pageSize,
  is_page: true,
  search,
  progress,
  only_my_orders: onlyMine,
})

export const summarizeMyOrders = (rows) => ({
  total: rows.length,
  pending: rows.filter((r) => r.progress === 'PENDING').length,
  executing: rows.filter((r) => r.progress === 'EXECUTING').length,
  failed: rows.filter((r) => r.progress === 'FAILED').length,
})
```

Refactor `OrderList.vue` to:
- use shared status meta + model helpers
- adopt `page-toolbar`, `overview-grid`, `table-shell` semantic classes
- apply responsive behavior for mobile/tablet breakpoints

**Step 4: Run test to verify it passes**

Run: `cd www && npm run test:unit -- src/views/orders/list/__tests__/order-list-model.spec.js`

Expected: PASS.

**Step 5: Commit**

```bash
git add www/src/views/orders/list/orderListModel.js www/src/views/orders/list/__tests__/order-list-model.spec.js www/src/views/orders/list/OrderList.vue
git commit -m "refactor: apply page pattern to order list"
```

### Task 8: Refactor Order Detail Sections + Detail Model

**Files:**
- Create: `www/src/views/orders/detail/orderDetailModel.js`
- Modify: `www/src/views/orders/detail/OrderDetail.vue`
- Modify: `www/src/views/orders/detail/HeaderContent.vue`
- Modify: `www/src/views/orders/detail/HeaderExtra.vue`
- Modify: `www/src/views/orders/detail/ApprovalSteps.vue`
- Test: `www/src/views/orders/detail/__tests__/order-detail-model.spec.js`

**Step 1: Write the failing test**

```js
import { describe, expect, it } from 'vitest'
import { normalizeClaimUsers } from '../orderDetailModel'

describe('order detail model', () => {
  it('parses json claim_users string', () => {
    expect(normalizeClaimUsers('[\"alice\",\"bob\"]')).toBe('alice, bob')
  })

  it('returns fallback for empty values', () => {
    expect(normalizeClaimUsers('')).toBe('无')
    expect(normalizeClaimUsers(null)).toBe('无')
  })
})
```

**Step 2: Run test to verify it fails**

Run: `cd www && npm run test:unit -- src/views/orders/detail/__tests__/order-detail-model.spec.js`

Expected: FAIL because `orderDetailModel.js` does not exist.

**Step 3: Write minimal implementation**

```js
// www/src/views/orders/detail/orderDetailModel.js
export const normalizeClaimUsers = (raw) => {
  if (!raw) return '无'
  if (Array.isArray(raw)) return raw.length ? raw.join(', ') : '无'
  if (typeof raw === 'string') {
    try {
      const users = JSON.parse(raw)
      return Array.isArray(users) && users.length ? users.join(', ') : (raw || '无')
    } catch {
      return raw || '无'
    }
  }
  return '无'
}
```

Refactor detail files to:
- replace inline style cards with semantic section classes
- align action buttons hierarchy
- standardize log container and approval section spacing
- use model helper in `HeaderContent.vue`

**Step 4: Run test to verify it passes**

Run: `cd www && npm run test:unit -- src/views/orders/detail/__tests__/order-detail-model.spec.js`

Expected: PASS.

**Step 5: Commit**

```bash
git add www/src/views/orders/detail/orderDetailModel.js www/src/views/orders/detail/__tests__/order-detail-model.spec.js www/src/views/orders/detail/OrderDetail.vue www/src/views/orders/detail/HeaderContent.vue www/src/views/orders/detail/HeaderExtra.vue www/src/views/orders/detail/ApprovalSteps.vue
git commit -m "refactor: modernize order detail structure and sections"
```

### Task 9: Refactor Login Page with Shared Tokenized Behaviors

**Files:**
- Create: `www/src/views/login/loginModel.js`
- Modify: `www/src/views/login/Login.vue`
- Test: `www/src/views/login/__tests__/login-model.spec.js`

**Step 1: Write the failing test**

```js
import { describe, expect, it } from 'vitest'
import { normalizeOtpCode } from '../loginModel'

describe('login otp model', () => {
  it('keeps only 6 digits', () => {
    expect(normalizeOtpCode('a1b2c3d4')).toBe('1234')
    expect(normalizeOtpCode('123456789')).toBe('123456')
  })

  it('handles empty values', () => {
    expect(normalizeOtpCode()).toBe('')
  })
})
```

**Step 2: Run test to verify it fails**

Run: `cd www && npm run test:unit -- src/views/login/__tests__/login-model.spec.js`

Expected: FAIL because `loginModel.js` does not exist.

**Step 3: Write minimal implementation**

```js
// www/src/views/login/loginModel.js
export const normalizeOtpCode = (value) => String(value || '').replace(/\D/g, '').slice(0, 6)
```

Refactor `Login.vue` to:
- consume `normalizeOtpCode` helper
- replace local hardcoded values with token vars
- finalize desktop/tablet/mobile visual rhythm and spacing consistency

**Step 4: Run test to verify it passes**

Run: `cd www && npm run test:unit -- src/views/login/__tests__/login-model.spec.js`

Expected: PASS.

**Step 5: Commit**

```bash
git add www/src/views/login/loginModel.js www/src/views/login/__tests__/login-model.spec.js www/src/views/login/Login.vue
git commit -m "refactor: polish responsive login with tokenized styles"
```

### Task 10: Full Verification and Release-Ready Checklist

**Files:**
- Create: `docs/plans/2026-03-03-ui-revamp-phase1-qa-checklist.md`
- Modify: `docs/plans/2026-03-03-ui-revamp-phase1-design.md` (append executed evidence section)

**Step 1: Write the failing verification script/checklist target**

Create checklist entries with strict pass criteria:
- unit tests all pass
- lint passes
- build passes
- no static inline styles in phase-1 files
- screenshot checks done for mobile/tablet/desktop/wide

**Step 2: Run verification commands and capture failures**

Run:

```bash
cd www && npm run test:unit
cd www && npm run lint
cd www && npm run build
```

Expected: Any failure is documented as blocking item in checklist.

**Step 3: Apply minimal fixes until all commands pass**

- Fix lint/type/build regressions introduced during refactor.
- Re-run failing command immediately after each fix.

**Step 4: Re-run full verification to green**

Run:

```bash
cd www && npm run test:unit && npm run lint && npm run build
```

Expected: Full PASS, plus manual breakpoint checks complete.

**Step 5: Commit**

```bash
git add docs/plans/2026-03-03-ui-revamp-phase1-qa-checklist.md docs/plans/2026-03-03-ui-revamp-phase1-design.md
git commit -m "docs: record phase1 UI verification evidence"
```

