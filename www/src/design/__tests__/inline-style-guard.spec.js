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
  it('has no style="..." attributes', () => {
    for (const file of files) {
      const content = fs.readFileSync(path.resolve(file), 'utf-8')
      expect(content).not.toMatch(/\sstyle="[^"]+"/g)
    }
  })
})
