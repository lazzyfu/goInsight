import fs from 'node:fs'
import path from 'node:path'
import { describe, expect, it } from 'vitest'

const DAS_SHELL_BATCH_FILES = [
  'src/views/das/index.vue',
  'src/views/das/favorite/DasFavorite.vue',
  'src/views/das/favorite/DasFavoriteFormModal.vue',
  'src/views/das/history/DasHistory.vue',
]

const DAS_CONSOLE_BATCH_FILES = [
  'src/views/das/console/ConsoleIndex.vue',
  'src/views/das/console/ConsoleLeft.vue',
  'src/views/das/console/ConsoleRight.vue',
  'src/views/das/console/ConsoleDbDict.vue',
  'src/views/das/console/components/PermissionHint.vue',
  'src/views/das/console/components/IconTiDB.vue',
  'src/views/das/console/components/IconMySQL.vue',
  'src/views/das/console/components/IconClickHouse.vue',
]

const BASELINE_COUNTS = {
  shell: 7,
  console: 7,
}

const readFile = (file) => fs.readFileSync(path.resolve(file), 'utf-8')
const countStaticInlineStyles = (content) => (content.match(/\sstyle="[^"]+"/g) || []).length

const countByFiles = (files) =>
  files.reduce((total, file) => total + countStaticInlineStyles(readFile(file)), 0)

describe('phase2 das page contract', () => {
  it('has all planned migration files', () => {
    const files = [...DAS_SHELL_BATCH_FILES, ...DAS_CONSOLE_BATCH_FILES]
    for (const file of files) {
      expect(fs.existsSync(path.resolve(file))).toBe(true)
    }
  })

  it('does not increase shell batch inline style baseline before migration', () => {
    expect(countByFiles(DAS_SHELL_BATCH_FILES)).toBeLessThanOrEqual(BASELINE_COUNTS.shell)
  })

  it('does not increase console batch inline style baseline before migration', () => {
    expect(countByFiles(DAS_CONSOLE_BATCH_FILES)).toBeLessThanOrEqual(BASELINE_COUNTS.console)
  })
})
