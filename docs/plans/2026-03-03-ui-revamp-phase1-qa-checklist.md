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

- [ ] Login page desktop/tablet/mobile layout remains usable.
- [ ] Layout shell header/sider/content spacing is consistent after tokenization.
- [ ] Order list filter, statistics cards, and table still function with responsive wrap.
- [ ] Order detail status tag, action group, approval flow, log area, and SQL section render correctly.

Note: Manual checklist items should be validated in browser after pulling this branch.

## 4) Follow-up Items (Non-blocking)

1. Migrate SCSS `@import` syntax to `@use`/`@forward`.
2. Split oversized bundles with route-level manual chunks.
3. Extend token/pattern migration to admin and das modules in Phase 2.
