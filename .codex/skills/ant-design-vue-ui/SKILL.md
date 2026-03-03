Name: ant-design-vue-ui
Purpose: Generate consistent UI using Ant Design Vue for enterprise applications

Framework
- Vue 3
- Ant Design Vue
- Composition API
- script setup

General UI Principles
- Follow Ant Design Design System
- Use built-in components first
- Avoid custom styles unless necessary
- Keep UI simple and consistent
- Maintain spacing and layout hierarchy

Layout System

Use standard admin layout:

Layout
 ├── Sider (Menu)
 ├── Header
 └── Content

Rules
- Use Layout component
- Sidebar width 200-240px
- Content wrapped in Card
- Page padding 16-24px

Grid System
- Use 24 grid layout
- Use Row / Col
- Maintain spacing (gutter 16 or 24)

Spacing Rules
Small: 8px
Medium: 16px
Large: 24px
Section: 32px

---

Color System

Primary Color
# 1677ff

Success
# 52c41a

Warning
# faad14

Error
# ff4d4f

Info
# 1677ff

Background
Page: #f5f7fa
Card: #ffffff

Text Color
Primary text: rgba(0,0,0,0.88)
Secondary text: rgba(0,0,0,0.65)

Border
# d9d9d9

Rules
- Use Ant Design semantic colors
- Do not invent new colors
- Important action uses primary color
- Status uses Tag colors

---

Component Usage Rules

Buttons

Primary action
a-button type="primary"

Secondary
a-button

Danger
a-button danger

Avoid more than 2 primary buttons on a page.

---

Forms

Use
a-form

Rules
- Always include label
- Use validation rules
- Align labels
- Complex forms use vertical layout

Common Inputs
- Input
- Select
- DatePicker
- Switch
- Radio
- Checkbox

---

Tables

Use
a-table

Rules
- Always set rowKey
- Use pagination
- Avoid too many columns
- Prefer simple render logic
- Support sorting if needed

Standard Table Layout

Card
 ├── Search Form
 └── Table

---

Cards

Use
a-card

Rules
- Wrap important content
- Use title
- Avoid nested cards

---

Navigation

Use

Menu
Breadcrumb
Tabs
Dropdown

Rules
- Menu in sider
- Breadcrumb at page top
- Tabs for multi-view pages

---

Feedback Components

Use

Message
Modal
Notification
Alert
Spin

Rules

Success
message.success()

Error
notification.error()

Confirm
modal.confirm()

Loading
spin

---

Status Display

Use Tag

Success
green

Processing
blue

Warning
orange

Error
red

Default
gray

---

Page Structure Rules

List Page

Search Form
Table
Pagination

Detail Page

Descriptions
Tabs if needed

Edit Page

Form
Submit Button

Dashboard Page

Card
Statistic
Charts

Settings Page

Tabs
Forms

---

Interaction Rules

- Actions on top-right
- Search on top
- Table in center
- Pagination bottom-right
- Buttons grouped logically

---

Code Style

- Composition API
- script setup
- Reactive state
- Avoid inline styles
- Use Ant Design icons

---

When generating UI

List data → Table
Edit data → Form
Show info → Descriptions
Show metrics → Statistic
Group content → Card
Multi sections → Tabs
