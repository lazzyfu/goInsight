import { describe, expect, it } from 'vitest'
import { antdTheme } from '../antdTheme'

describe('antd theme bridge', () => {
  it('maps shared tokens to antd component tokens', () => {
    expect(antdTheme.token.colorPrimary).toBe('#0f766e')
    expect(antdTheme.token.borderRadius).toBe(10)
    expect(antdTheme.components.Card.borderRadiusLG).toBe(12)
  })
})
