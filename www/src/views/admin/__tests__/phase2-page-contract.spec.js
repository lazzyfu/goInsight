import fs from 'node:fs'
import path from 'node:path'
import { describe, expect, it } from 'vitest'

const ADMIN_LIST_BATCH_FILES = [
  'src/views/admin/perms/users/UserList.vue',
  'src/views/admin/perms/roles/RoleList.vue',
  'src/views/admin/perms/flows/ApprovalFlowList.vue',
  'src/views/admin/system/environments/EnvironmentList.vue',
  'src/views/admin/system/instances/InstanceList.vue',
  'src/views/admin/system/inspect/InspectParamsList.vue',
  'src/views/admin/system/das/DasSchemaList.vue',
  'src/views/admin/system/das/DasTableList.vue',
]

const ADMIN_FORM_BATCH_FILES = [
  'src/views/admin/perms/users/UserFormModal.vue',
  'src/views/admin/perms/users/PasswordFormModal.vue',
  'src/views/admin/perms/roles/RoleFormModal.vue',
  'src/views/admin/system/environments/EnvironmentFormModal.vue',
  'src/views/admin/system/instances/InstanceFormModal.vue',
  'src/views/admin/system/instances/InstanceInspectParamsFormModal.vue',
  'src/views/admin/system/inspect/InspectParamsFormModal.vue',
  'src/views/admin/system/das/DasSchemaFormModal.vue',
  'src/views/admin/system/das/DasTableFormModal.vue',
  'src/views/admin/perms/orgs/BindOrgUsers.vue',
  'src/views/admin/perms/orgs/OrgUsers.vue',
  'src/views/admin/perms/flows/FlowBoundUsersDetail.vue',
]

const ADMIN_ORG_BATCH_FILES = [
  'src/views/admin/perms/orgs/OrgList.vue',
  'src/views/admin/perms/orgs/AddRootOrg.vue',
  'src/views/admin/perms/orgs/AddChildOrg.vue',
  'src/views/admin/perms/orgs/EditOrgName.vue',
]

const BASELINE_COUNTS = {
  list: 0,
  form: 0,
  org: 1,
}

const readFile = (file) => fs.readFileSync(path.resolve(file), 'utf-8')
const countStaticInlineStyles = (content) => (content.match(/\sstyle="[^"]+"/g) || []).length

const countByFiles = (files) =>
  files.reduce((total, file) => total + countStaticInlineStyles(readFile(file)), 0)

describe('phase2 admin page contract', () => {
  it('has all planned migration files', () => {
    const files = [...ADMIN_LIST_BATCH_FILES, ...ADMIN_FORM_BATCH_FILES, ...ADMIN_ORG_BATCH_FILES]
    for (const file of files) {
      expect(fs.existsSync(path.resolve(file))).toBe(true)
    }
  })

  it('removes static inline styles from list batch pages', () => {
    expect(countByFiles(ADMIN_LIST_BATCH_FILES)).toBeLessThanOrEqual(BASELINE_COUNTS.list)
  })

  it('removes static inline styles from form and modal batch pages', () => {
    expect(countByFiles(ADMIN_FORM_BATCH_FILES)).toBeLessThanOrEqual(BASELINE_COUNTS.form)
  })

  it('does not increase org batch inline style baseline before migration', () => {
    expect(countByFiles(ADMIN_ORG_BATCH_FILES)).toBeLessThanOrEqual(BASELINE_COUNTS.org)
  })
})
