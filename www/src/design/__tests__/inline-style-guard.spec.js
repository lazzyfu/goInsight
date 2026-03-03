import fs from 'node:fs'
import path from 'node:path'
import { describe, expect, it } from 'vitest'

const PHASE1_ZERO_INLINE_STYLE_FILES = [
  'src/components/layout/Layout.vue',
  'src/views/orders/list/OrderList.vue',
  'src/views/orders/detail/OrderDetail.vue',
  'src/views/orders/detail/HeaderExtra.vue',
]

const PHASE2_ADMIN_DAS_INLINE_STYLE_BASELINE = 50

const readFile = (file) => fs.readFileSync(path.resolve(file), 'utf-8')

const countStaticInlineStyles = (content) => {
  const matches = content.match(/\sstyle="[^"]+"/g)
  return matches ? matches.length : 0
}

const collectVueFiles = (dir) => {
  const targetDir = path.resolve(dir)
  if (!fs.existsSync(targetDir)) return []

  return fs.readdirSync(targetDir, { withFileTypes: true }).flatMap((entry) => {
    const resolvedPath = path.join(targetDir, entry.name)
    if (entry.isDirectory()) {
      return collectVueFiles(resolvedPath)
    }
    if (entry.isFile() && entry.name.endsWith('.vue')) {
      return [resolvedPath]
    }
    return []
  })
}

describe('ui inline style guardrails', () => {
  it('has no style="..." attributes', () => {
    for (const file of PHASE1_ZERO_INLINE_STYLE_FILES) {
      const content = readFile(file)
      expect(content).not.toMatch(/\sstyle="[^"]+"/g)
    }
  })

  it('does not increase static inline styles in admin and das during phase2 migration', () => {
    const adminFiles = collectVueFiles('src/views/admin')
    const dasFiles = collectVueFiles('src/views/das')
    const allFiles = [...adminFiles, ...dasFiles]

    const totalInlineStyles = allFiles.reduce((total, file) => {
      return total + countStaticInlineStyles(readFile(file))
    }, 0)

    expect(totalInlineStyles).toBeLessThanOrEqual(PHASE2_ADMIN_DAS_INLINE_STYLE_BASELINE)
  })
})
