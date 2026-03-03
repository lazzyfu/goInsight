# UI Revamp Phase 2 (Admin + DAS) Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Migrate `admin` and `das` modules onto the Phase 1 tokenized UI system with consistent layout, spacing, status semantics, and responsive behavior.

**Architecture:** Reuse the Phase 1 foundation (`tokens.scss`, `antdTheme`, layout shell, status/meta models), then migrate Phase 2 pages in vertical slices: first list pages and shared patterns, then modal/forms, then complex tree/console pages. Keep business API behavior unchanged and enforce migration with contract tests (inline style guard, page-pattern usage, and critical view-model tests).

**Tech Stack:** Vue 3 (`script setup`), Ant Design Vue 4, SCSS, Vitest, ESLint, Vite.

---

## Scope Snapshot

- In scope:
  - `www/src/views/admin/**`
  - `www/src/views/das/**`
  - shared UI helpers under `www/src/components` and `www/src/design`
  - QA docs under `docs/plans`
- Out of scope:
  - backend API semantics
  - permission model changes
  - query engine behavior changes

Current baseline:
- Phase 1 merged into `dev`.
- `admin + das` still have many static inline styles (`rg 'style="'` shows 53 occurrences).

## Execution Rules

- Preserve behavior first; UI consistency second; refactor extraction third.
- Migrate by route batches to keep rollback small.
- Every batch must pass: `npm run test:unit`, `npm run lint`, `npm run build`.
- Any newly touched page must not add new static `style="..."` attributes.

---

### Task 1: Extend Migration Guardrails for Phase 2

**Files:**
- Modify: `www/src/design/__tests__/inline-style-guard.spec.js`
- Create: `www/src/views/admin/__tests__/phase2-page-contract.spec.js`
- Create: `www/src/views/das/__tests__/phase2-page-contract.spec.js`

**Step 1: Write the failing test**

- Add Phase 2 route components into inline-style guard target list.
- Add contract checks for required class hooks:
  - admin list pages contain `gi-page-shell` and `gi-page-toolbar`
  - das shell pages contain `gi-page-shell` and `sql-shell` wrapper classes

**Step 2: Run test to verify it fails**

Run:

```bash
cd www
npm run test:unit -- src/design/__tests__/inline-style-guard.spec.js
```

Expected: FAIL on existing admin/das inline styles.

**Step 3: Write minimal implementation**

- Update test files only (no page refactor yet).
- Keep assertions strict but route-scoped, so failures map to migration tasks.

**Step 4: Run test to verify it passes**

- Not expected to pass until migration tasks complete.
- Record current failing files to seed Task 2-8 work list.

**Step 5: Commit**

```bash
git add www/src/design/__tests__/inline-style-guard.spec.js www/src/views/admin/__tests__/phase2-page-contract.spec.js www/src/views/das/__tests__/phase2-page-contract.spec.js
git commit -m "test: add phase2 ui migration guardrails"
```

### Task 2: Build Reusable Admin/DAS Page Patterns

**Files:**
- Create: `www/src/components/patterns/PageToolbar.vue`
- Create: `www/src/components/patterns/PageTableSection.vue`
- Create: `www/src/components/patterns/PageCardShell.vue`
- Create: `www/src/components/patterns/pagePatterns.scss`
- Modify: `www/src/assets/scss/base.scss`
- Test: `www/src/components/patterns/__tests__/page-patterns.spec.js`

**Step 1: Write the failing test**

- Component contract test:
  - `PageToolbar` renders left/right slots with wrap behavior class.
  - `PageTableSection` applies top spacing class without inline style.

**Step 2: Run test to verify it fails**

Run:

```bash
cd www
npm run test:unit -- src/components/patterns/__tests__/page-patterns.spec.js
```

Expected: FAIL because pattern components do not exist.

**Step 3: Write minimal implementation**

- Build thin wrappers around Antd primitives:
  - no new business logic
  - token-driven spacing/radius/hover styles only

**Step 4: Run test to verify it passes**

Run:

```bash
cd www
npm run test:unit -- src/components/patterns/__tests__/page-patterns.spec.js
```

Expected: PASS.

**Step 5: Commit**

```bash
git add www/src/components/patterns www/src/assets/scss/base.scss
git commit -m "feat: add reusable admin and das page patterns"
```

### Task 3: Migrate Admin List Pages (Batch A)

**Files:**
- Modify:
  - `www/src/views/admin/perms/users/UserList.vue`
  - `www/src/views/admin/perms/roles/RoleList.vue`
  - `www/src/views/admin/perms/flows/ApprovalFlowList.vue`
  - `www/src/views/admin/system/environments/EnvironmentList.vue`
  - `www/src/views/admin/system/instances/InstanceList.vue`
  - `www/src/views/admin/system/inspect/InspectParamsList.vue`
  - `www/src/views/admin/system/das/DasSchemaList.vue`
  - `www/src/views/admin/system/das/DasTableList.vue`

**Step 1: Write the failing test**

- Run guardrail tests from Task 1 (expected fail).
- Add one view-model test for list query builder extraction if missing.

**Step 2: Run test to verify it fails**

Run:

```bash
cd www
npm run test:unit -- src/design/__tests__/inline-style-guard.spec.js
```

Expected: FAIL on above files.

**Step 3: Write minimal implementation**

- Replace static inline styles (`width`, `margin-top`, `padding`) with:
  - `PageToolbar`
  - `PageTableSection`
  - scoped semantic classes using token vars
- Keep API calls and pagination logic unchanged.

**Step 4: Run test to verify it passes**

Run:

```bash
cd www
npm run test:unit -- src/design/__tests__/inline-style-guard.spec.js
npm run lint
```

Expected: relevant list-page inline style failures cleared.

**Step 5: Commit**

```bash
git add www/src/views/admin/perms/users/UserList.vue www/src/views/admin/perms/roles/RoleList.vue www/src/views/admin/perms/flows/ApprovalFlowList.vue www/src/views/admin/system/environments/EnvironmentList.vue www/src/views/admin/system/instances/InstanceList.vue www/src/views/admin/system/inspect/InspectParamsList.vue www/src/views/admin/system/das/DasSchemaList.vue www/src/views/admin/system/das/DasTableList.vue
git commit -m "refactor: migrate admin list pages to phase2 page patterns"
```

### Task 4: Migrate Admin Modal/Form Pages (Batch B)

**Files:**
- Modify:
  - `www/src/views/admin/perms/users/UserFormModal.vue`
  - `www/src/views/admin/perms/users/PasswordFormModal.vue`
  - `www/src/views/admin/perms/roles/RoleFormModal.vue`
  - `www/src/views/admin/system/environments/EnvironmentFormModal.vue`
  - `www/src/views/admin/system/instances/InstanceFormModal.vue`
  - `www/src/views/admin/system/instances/InstanceInspectParamsFormModal.vue`
  - `www/src/views/admin/system/inspect/InspectParamsFormModal.vue`
  - `www/src/views/admin/system/das/DasSchemaFormModal.vue`
  - `www/src/views/admin/system/das/DasTableFormModal.vue`
  - `www/src/views/admin/perms/orgs/BindOrgUsers.vue`
  - `www/src/views/admin/perms/orgs/OrgUsers.vue`
  - `www/src/views/admin/perms/flows/FlowBoundUsersDetail.vue`

**Step 1: Write the failing test**

- Add modal form guard test ensuring:
  - no static inline `style="width: 100%"` in touched modals
  - use `.modal-field-full-width` or shared utility class

**Step 2: Run test to verify it fails**

Run:

```bash
cd www
npm run test:unit -- src/views/admin/__tests__/phase2-page-contract.spec.js
```

Expected: FAIL on form modal files.

**Step 3: Write minimal implementation**

- Normalize form layout spacing to token scale.
- Introduce shared full-width field class.
- Keep field validation and submit behavior unchanged.

**Step 4: Run test to verify it passes**

Run:

```bash
cd www
npm run test:unit -- src/views/admin/__tests__/phase2-page-contract.spec.js
npm run lint
```

Expected: PASS.

**Step 5: Commit**

```bash
git add www/src/views/admin/perms/users/UserFormModal.vue www/src/views/admin/perms/users/PasswordFormModal.vue www/src/views/admin/perms/roles/RoleFormModal.vue www/src/views/admin/system/environments/EnvironmentFormModal.vue www/src/views/admin/system/instances/InstanceFormModal.vue www/src/views/admin/system/instances/InstanceInspectParamsFormModal.vue www/src/views/admin/system/inspect/InspectParamsFormModal.vue www/src/views/admin/system/das/DasSchemaFormModal.vue www/src/views/admin/system/das/DasTableFormModal.vue www/src/views/admin/perms/orgs/BindOrgUsers.vue www/src/views/admin/perms/orgs/OrgUsers.vue www/src/views/admin/perms/flows/FlowBoundUsersDetail.vue
git commit -m "refactor: migrate admin modal and form styles to phase2 tokens"
```

### Task 5: Migrate Admin Organization Complex Page (Batch C)

**Files:**
- Modify:
  - `www/src/views/admin/perms/orgs/OrgList.vue`
  - `www/src/views/admin/perms/orgs/AddRootOrg.vue`
  - `www/src/views/admin/perms/orgs/AddChildOrg.vue`
  - `www/src/views/admin/perms/orgs/EditOrgName.vue`

**Step 1: Write the failing test**

- Add org-page contract test:
  - hero + panels use tokenized classes
  - tree search input no inline style

**Step 2: Run test to verify it fails**

Run:

```bash
cd www
npm run test:unit -- src/views/admin/__tests__/phase2-page-contract.spec.js
```

Expected: FAIL for org page inline style and spacing.

**Step 3: Write minimal implementation**

- Keep current hero design direction.
- Replace remaining inline styles with semantic classes.
- Ensure tablet/mobile layout still readable.

**Step 4: Run test to verify it passes**

Run:

```bash
cd www
npm run test:unit -- src/views/admin/__tests__/phase2-page-contract.spec.js
```

Expected: PASS.

**Step 5: Commit**

```bash
git add www/src/views/admin/perms/orgs/OrgList.vue www/src/views/admin/perms/orgs/AddRootOrg.vue www/src/views/admin/perms/orgs/AddChildOrg.vue www/src/views/admin/perms/orgs/EditOrgName.vue
git commit -m "refactor: migrate admin org management layout to phase2 patterns"
```

### Task 6: Migrate DAS Shell + Favorites + History (Batch D)

**Files:**
- Modify:
  - `www/src/views/das/index.vue`
  - `www/src/views/das/favorite/DasFavorite.vue`
  - `www/src/views/das/favorite/DasFavoriteFormModal.vue`
  - `www/src/views/das/history/DasHistory.vue`

**Step 1: Write the failing test**

- Add DAS list/page contract test:
  - no static inline style in shell/favorite/history pages
  - use shared page toolbar/table section classes

**Step 2: Run test to verify it fails**

Run:

```bash
cd www
npm run test:unit -- src/views/das/__tests__/phase2-page-contract.spec.js
```

Expected: FAIL on the above pages.

**Step 3: Write minimal implementation**

- Replace inline body-style and spacing with classes.
- Preserve tab behavior and modal payload structure.

**Step 4: Run test to verify it passes**

Run:

```bash
cd www
npm run test:unit -- src/views/das/__tests__/phase2-page-contract.spec.js
```

Expected: PASS.

**Step 5: Commit**

```bash
git add www/src/views/das/index.vue www/src/views/das/favorite/DasFavorite.vue www/src/views/das/favorite/DasFavoriteFormModal.vue www/src/views/das/history/DasHistory.vue
git commit -m "refactor: migrate das shell, favorite and history pages"
```

### Task 7: Migrate DAS Console Workspace (Batch E)

**Files:**
- Modify:
  - `www/src/views/das/console/ConsoleIndex.vue`
  - `www/src/views/das/console/ConsoleLeft.vue`
  - `www/src/views/das/console/ConsoleRight.vue`
  - `www/src/views/das/console/ConsoleDbDict.vue`
  - `www/src/views/das/console/components/PermissionHint.vue`
  - `www/src/views/das/console/components/IconTiDB.vue`
  - `www/src/views/das/console/components/IconMySQL.vue`
  - `www/src/views/das/console/components/IconClickHouse.vue`
- Test:
  - `www/src/views/das/console/__tests__/console-layout-model.spec.js`

**Step 1: Write the failing test**

- Add console layout model test for height sync and compact breakpoint handling.
- Add inline-style contract checks for console child components.

**Step 2: Run test to verify it fails**

Run:

```bash
cd www
npm run test:unit -- src/views/das/console/__tests__/console-layout-model.spec.js
```

Expected: FAIL if helper/model does not exist yet.

**Step 3: Write minimal implementation**

- Extract pure helpers for height/viewport calculations.
- Keep SQL execution flow unchanged.
- Convert static inline styles in console pages/icons to scoped classes.

**Step 4: Run test to verify it passes**

Run:

```bash
cd www
npm run test:unit -- src/views/das/console/__tests__/console-layout-model.spec.js
npm run lint
```

Expected: PASS.

**Step 5: Commit**

```bash
git add www/src/views/das/console
git commit -m "refactor: migrate das console workspace to phase2 tokenized styles"
```

### Task 8: Final Contract Cleanup and Inline-Style Zeroing

**Files:**
- Modify: `www/src/design/__tests__/inline-style-guard.spec.js`
- Modify: `www/src/views/admin/__tests__/phase2-page-contract.spec.js`
- Modify: `www/src/views/das/__tests__/phase2-page-contract.spec.js`

**Step 1: Write the failing target**

- Update guardrail baseline to expect zero static inline style in:
  - `views/admin/**`
  - `views/das/**`
  - except explicit dynamic `:style` use cases

**Step 2: Run test to verify current state**

Run:

```bash
cd www
npm run test:unit -- src/design/__tests__/inline-style-guard.spec.js
```

Expected: either PASS or list residual files.

**Step 3: Write minimal implementation**

- Fix residual failures only.
- Avoid opportunistic refactors outside scope.

**Step 4: Run test to verify it passes**

Run:

```bash
cd www
npm run test:unit
```

Expected: all tests PASS.

**Step 5: Commit**

```bash
git add www/src/design/__tests__/inline-style-guard.spec.js www/src/views/admin/__tests__/phase2-page-contract.spec.js www/src/views/das/__tests__/phase2-page-contract.spec.js
git commit -m "test: enforce phase2 inline-style migration completion"
```

### Task 9: Verification, QA Evidence, and Release Notes

**Files:**
- Create: `docs/plans/2026-03-03-ui-revamp-phase2-admin-das-qa-checklist.md`
- Modify: `docs/plans/2026-03-03-ui-revamp-phase2-admin-das-implementation.md`

**Step 1: Define QA checklist**

- Add route and breakpoint matrix:
  - admin users/roles/orgs/flows
  - system env/instance/inspect/das
  - das console/favorite/history
  - breakpoints: 1440, 1024, 390

**Step 2: Run full verification**

Run:

```bash
cd www
npm run test:unit
npm run lint
npm run build
```

Expected: PASS (allow non-blocking warnings with note).

**Step 3: Capture manual QA evidence**

- Screenshots and issue notes for each route group.
- Record any known non-blockers.

**Step 4: Update plan with execution evidence**

- Append commit log and verification summary.

**Step 5: Commit**

```bash
git add docs/plans/2026-03-03-ui-revamp-phase2-admin-das-qa-checklist.md docs/plans/2026-03-03-ui-revamp-phase2-admin-das-implementation.md
git commit -m "docs: add phase2 admin das qa evidence"
```

---

## Risk and Mitigation

- Risk: long tail of inline style removals in modals and SVG icon wrappers.
  - Mitigation: enforce guardrail tests per batch.
- Risk: DAS console interactions regress due to layout refactor.
  - Mitigation: extract layout math helpers and test them directly.
- Risk: visual inconsistency across old and migrated pages during transition.
  - Mitigation: batch by route group and keep each commit reviewable.

## Exit Criteria

- `admin + das` pages pass zero-static-inline-style contract (except approved dynamic cases).
- All touched routes follow tokenized spacing/color/radius semantics.
- Responsive checks pass on desktop/tablet/mobile.
- Test, lint, and build all green on `dev`.

## Execution Evidence (2026-03-03)

### Completed Commits

1. `4c6db640` test: add phase2 ui migration guardrails
2. `4b725ee8` feat: add reusable admin and das page patterns
3. `f4aa826e` refactor: migrate admin list pages to phase2 page patterns
4. `e1e77813` refactor: migrate admin modal and form styles to phase2 tokens
5. `aba89db9` refactor: migrate admin org management layout to phase2 patterns
6. `1014bf23` refactor: migrate das shell, favorite and history pages
7. `dba49b49` refactor: migrate das console workspace to phase2 tokenized styles
8. `3197f3fe` test: enforce phase2 inline-style migration completion

### Verification Snapshot

- Guardrails:
  - `cd www && npm run test:unit -- src/design/__tests__/inline-style-guard.spec.js src/views/admin/__tests__/phase2-page-contract.spec.js src/views/das/__tests__/phase2-page-contract.spec.js`
  - Result: all pass
- Full unit tests:
  - `cd www && npm run test:unit`
  - Result: 12 files / 25 tests all pass
- Lint:
  - `cd www && npm run lint`
  - Result: pass
- Build:
  - `cd www && npm run build`
  - Result: pass (non-blocking Sass `@import` and chunk-size warnings remain)

### Notes

- `admin + das` static inline style attributes in Vue templates are now `0` by guardrail checks.
- Manual breakpoint walkthrough remains tracked in:
  - `docs/plans/2026-03-03-ui-revamp-phase2-admin-das-qa-checklist.md`
