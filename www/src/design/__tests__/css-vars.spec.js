import fs from 'node:fs'
import path from 'node:path'
import { describe, expect, it } from 'vitest'

describe('css variable contract', () => {
  it('defines phase-1 required variables', () => {
    const file = path.resolve('src/assets/scss/tokens.scss')
    const content = fs.readFileSync(file, 'utf-8')
    expect(content).toContain('--gi-spacing-md: 16px;')
    expect(content).toContain('--gi-color-primary: #0f766e;')
    expect(content).toContain('--gi-radius-card: 12px;')
  })
})
