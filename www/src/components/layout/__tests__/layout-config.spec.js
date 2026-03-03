import { describe, expect, it } from 'vitest'
import { layoutConfig } from '../layoutConfig'

describe('layout config', () => {
  it('defines consistent shell dimensions', () => {
    expect(layoutConfig.headerHeight).toBe(52)
    expect(layoutConfig.sidebarExpandedWidth).toBe(224)
    expect(layoutConfig.contentPadding.desktop).toBe(24)
    expect(layoutConfig.contentPadding.mobile).toBe(12)
  })
})
