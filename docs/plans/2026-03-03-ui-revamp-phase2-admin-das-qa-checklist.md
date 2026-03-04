# UI Revamp Phase 2 (Admin + DAS) QA Checklist

Date: 2026-03-03
Branch: `codex/ui-revamp-phase2-admin-das`
Scope: `www/src/views/admin/**` + `www/src/views/das/**` + phase2 contract guardrails

## 1) Automated Verification

- [x] Lint passes
  - Command: `cd www && npm run lint`
  - Result: passed (no remaining ESLint errors)
- [x] Build passes
  - Command: `cd www && npm run build`
  - Result: passed
  - Notes:
    - Sass `@import` deprecation warnings are still present in legacy SCSS entry.
    - Large chunk warnings (>500KB) remain as non-blocking pre-existing bundling concerns.

## 2) Contract Checks

- [x] Phase2 admin contract test
  - Result: list/form/org batches all enforce zero static inline styles
- [x] Phase2 das contract test
  - Result: shell/console batches both enforce zero static inline styles
- [x] Global inline style guard
  - Result: `admin + das` static inline style count = 0

## 3) Manual Route Matrix (Pending Browser Pass)

Run frontend:

```bash
cd www
npm run dev -- --host 0.0.0.0 --port 5175
```

Open: `http://localhost:5175`

Breakpoints:
- Desktop: 1440 x 900
- Tablet: 1024 x 768
- Mobile: 390 x 844

### 3.1 Admin - Permission Routes

- [ ] `/admin/perms/users`: toolbar wrap/table scroll/expanded row spacing
- [ ] `/admin/perms/roles`: search + table action spacing
- [ ] `/admin/perms/flows`: hero/table shell/tags and detail drawer
- [ ] `/admin/perms/orgs`: tree + users split panel behavior, add/edit modals

### 3.2 Admin - System Routes

- [ ] `/admin/system/environments`: toolbar and table spacing
- [ ] `/admin/system/instances`: table scroll + modal form alignment
- [ ] `/admin/system/inspect`: param list + edit modal field width
- [ ] `/admin/system/das/schemas` + `/admin/system/das/tables`: search/table/modal spacing consistency

### 3.3 DAS Routes

- [ ] `/das`: shell tabs and card body spacing
- [ ] `/das` 收藏SQL: search/table/modal behavior
- [ ] `/das` 历史SQL: search/table/detail modal cards
- [ ] `/das` 控制台: left/right pane heights, editor/result splitter, result fullscreen, dictionary modal

## 4) Known Non-blockers

1. Sass `@import` deprecation warnings (legacy SCSS import style).
2. Vite chunk-size warnings for large JS bundles.

## 5) Follow-up Suggestions

1. Migrate SCSS entry from `@import` to `@use` / `@forward`.
2. Add route-level chunk splitting for DAS console and admin heavy pages.
3. Capture final manual screenshots and attach to release note before merge to `dev`.
