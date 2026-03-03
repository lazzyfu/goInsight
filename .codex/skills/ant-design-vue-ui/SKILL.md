---
name: ant-design-vue-ui
description: Generate and refactor Vue 3 interfaces with Ant Design Vue using consistent layout, spacing, semantic colors, and component conventions. Use when building or polishing admin pages (list, detail, edit, dashboard, settings), especially for Ant Design Vue forms, tables, navigation, and status feedback in Composition API with script setup.
---

# Ant Design Vue UI

Use this skill to produce predictable, clean Ant Design Vue screens that stay close to the design system and avoid one-off styling.

## Required Stack

- Use Vue 3.
- Use Ant Design Vue.
- Use Composition API with `script setup`.
- Prefer Ant Design Vue built-in components before writing custom UI.

## Layout Rules

Use the standard admin shell:

```
Layout
 |- Sider (Menu)
 |- Header
 `- Content
```

- Use `Layout`, `Layout.Sider`, `Layout.Header`, and `Layout.Content`.
- Keep sidebar width between `200px` and `240px`.
- Wrap primary page content in `a-card`.
- Keep page padding between `16px` and `24px`.
- Use the 24-column grid with `a-row` and `a-col`.
- Use `gutter={16}` or `gutter={24}` for grid spacing.

Spacing tokens:

- Small: `8px`
- Medium: `16px`
- Large: `24px`
- Section: `32px`

## Color System

Use only these semantic tokens:

- Primary and Info: `#1677ff`
- Success: `#52c41a`
- Warning: `#faad14`
- Error: `#ff4d4f`
- Page background: `#f5f7fa`
- Card background: `#ffffff`
- Primary text: `rgba(0,0,0,0.88)`
- Secondary text: `rgba(0,0,0,0.65)`
- Border: `#d9d9d9`

Apply these rules:

- Use Ant Design semantic color usage instead of inventing ad-hoc colors.
- Use primary color for important actions.
- Use tag/status colors for operational states.

## Component Conventions

Buttons:

- Use `a-button type="primary"` for the primary action.
- Use default `a-button` for secondary actions.
- Use `a-button danger` for destructive actions.
- Keep no more than two primary buttons on a page.

Forms:

- Use `a-form` for all edit/filter data entry.
- Always provide labels and validation rules.
- Align labels consistently.
- Use vertical layout for complex forms.
- Prefer `a-input`, `a-select`, `a-date-picker`, `a-switch`, `a-radio-group`, and `a-checkbox-group`.

Tables:

- Use `a-table` for list data.
- Always set `rowKey`.
- Enable pagination by default.
- Keep columns focused; avoid over-dense tables.
- Keep render logic simple.
- Add sorting only when it helps the user complete tasks faster.

Cards and navigation:

- Use `a-card` to group important content.
- Add card titles for major sections.
- Avoid nested cards unless there is no cleaner structure.
- Use `a-menu` inside sider navigation.
- Use `a-breadcrumb` at the top of content pages.
- Use `a-tabs` for multi-view pages.
- Use `a-dropdown` for contextual action groups.

Feedback and status:

- Use `message.success()` for lightweight success feedback.
- Use `notification.error()` for prominent failures.
- Use `modal.confirm()` for irreversible operations.
- Use `a-spin` to indicate loading states.
- Use `a-tag` for status display with these colors:
  - Success: `green`
  - Processing: `blue`
  - Warning: `orange`
  - Error: `red`
  - Default: `gray`

## Page Blueprints

List page:

1. Place search form at top.
2. Place table in the center area.
3. Place pagination at the bottom-right.

Detail page:

1. Use `a-descriptions` for key-value details.
2. Add `a-tabs` only when content naturally splits into sections.

Edit page:

1. Use form-first structure.
2. Place submit and cancel actions in a clear action row.

Dashboard page:

1. Use cards for modules.
2. Use `a-statistic` for headline metrics.
3. Add charts in separate cards.

Settings page:

1. Use tabs to separate domains.
2. Keep each tab focused around one form set.

## Interaction Placement

- Place global and primary actions at the top-right.
- Place search/filter controls at the top.
- Keep data table or main content in the center.
- Group related buttons together with clear priority order.

## Generation Decision Map

- Map "list data" requests to `a-table`.
- Map "edit data" requests to `a-form`.
- Map "show detailed info" requests to `a-descriptions`.
- Map "show metrics" requests to `a-statistic`.
- Map "group content" requests to `a-card`.
- Map "split multiple sections" requests to `a-tabs`.

## Code Style Guardrails

- Use Composition API state and computed values.
- Use `script setup` syntax.
- Avoid inline styles; prefer Ant Design props and scoped class names only when needed.
- Use Ant Design icons for action cues.
- Keep UI simple, hierarchical, and consistent across pages.
