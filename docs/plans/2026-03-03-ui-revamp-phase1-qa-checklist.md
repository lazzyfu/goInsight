# UI Revamp Phase 1 QA Checklist

Date: 2026-03-03
Branch: `codex/ui-revamp-phase1`
Scope: Login / Layout shell / Order List / Order Detail / token foundation

## 1) Automated Verification

- [x] Unit tests pass
  - Command: `cd www && npm run test:unit`
  - Result: 9 files, 15 tests, all passed
- [x] Lint passes
  - Command: `cd www && npm run lint`
  - Result: passed (no remaining ESLint errors)
- [x] Build passes
  - Command: `cd www && npm run build`
  - Result: passed
  - Notes:
    - Sass `@import` deprecation warnings exist in legacy SCSS entry.
    - Bundles still report large chunk warnings (>500KB), pre-existing architecture concern.

## 2) Phase-1 Contract Checks

- [x] Design token contract test
  - File: `www/src/design/__tests__/tokens.spec.js`
- [x] CSS variable contract test
  - File: `www/src/design/__tests__/css-vars.spec.js`
- [x] Antd theme bridge contract test
  - File: `www/src/design/__tests__/antd-theme.spec.js`
- [x] Inline style guard test (phase-1 files)
  - File: `www/src/design/__tests__/inline-style-guard.spec.js`
- [x] Order status mapping contract test
  - File: `www/src/views/orders/shared/__tests__/order-status-meta.spec.js`
- [x] Order list model test
  - File: `www/src/views/orders/list/__tests__/order-list-model.spec.js`
- [x] Order detail model test
  - File: `www/src/views/orders/detail/__tests__/order-detail-model.spec.js`
- [x] Login OTP model test
  - File: `www/src/views/login/__tests__/login-model.spec.js`
- [x] Layout config test
  - File: `www/src/components/layout/__tests__/layout-config.spec.js`

## 3) Manual UX Spot Checks

- [x] Login page desktop/tablet/mobile layout remains usable.
- [x] Layout shell header/sider/content spacing is consistent after tokenization.
- [x] Order list filter, statistics cards, and table still function with responsive wrap.
- [x] Order detail status tag, action group, approval flow, log area, and SQL section render correctly.

Note: Manual checklist items should be validated in browser after pulling this branch.

### 3.1 Manual Walkthrough Steps

Run frontend:

```bash
cd www
npm run dev -- --host 0.0.0.0 --port 5175
```

Open: `http://localhost:5175`

Suggested viewport presets:

- Desktop: 1440 x 900
- Tablet: 1024 x 768
- Mobile: 390 x 844

#### A) Login page (`/login`)

- [x] Desktop: visual panel + form panel two-column layout; no clipping or overlap
- [x] Tablet: auto-switch to stacked layout; heading and form spacing still balanced
- [x] Mobile: input/button hit areas are comfortable and no horizontal scroll appears
- [x] OTP mode: switch from password to OTP input, tip text and account hint remain readable

#### B) Layout shell (after login)

- [x] Header height/spacing stable; breadcrumb and user dropdown vertically aligned
- [x] Sidebar collapse/expand works on desktop; icon alignment and hover style are consistent
- [x] Tablet/mobile breakpoint: sidebar auto-collapses and content area keeps readable padding
- [x] Content area background and container separation are visually clear

Note: B2 initially found "logo cannot fold". Fixed in `www/src/components/layout/Layout.vue` and retested as PASS.

#### C) Order List (`/orders`)

- [x] Filter bar wraps cleanly at tablet/mobile sizes
- [x] "我的工单" statistics cards align as 4 columns (desktop), 2 columns (tablet), 1 column (mobile)
- [x] Status tags show unified colors and text
- [x] Table remains usable on narrow screens (horizontal scroll available, no broken cells)

#### D) Order Detail (`/orders/:order_id`)

- [x] Header status tag + action group hierarchy is clear
- [x] Approval flow section spacing and typography are consistent
- [x] Operation log panel keeps max-height and scroll behavior
- [x] SQL content section keeps readable container spacing and no layout jump

#### E) Record evidence

- [ ] Save 4 screenshots: login, layout shell, order list, order detail (desktop)
- [ ] Save 3 screenshots: login/order list/order detail (mobile)
- [ ] If any issue found, record route + viewport + reproduction steps

## 4) Follow-up Items (Non-blocking)

1. Migrate SCSS `@import` syntax to `@use`/`@forward`.
2. Split oversized bundles with route-level manual chunks.
3. Extend token/pattern migration to admin and das modules in Phase 2.
