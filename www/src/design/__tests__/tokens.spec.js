import { describe, expect, it } from 'vitest'
import { breakpoints, radiusScale, spacingScale, typographyScale } from '../tokens'

describe('design tokens', () => {
  it('exports spacing scale in 8pt rhythm', () => {
    expect(spacingScale.md).toBe(16)
    expect(spacingScale.xl).toBe(32)
  })

  it('exports responsive breakpoints', () => {
    expect(breakpoints.mobile).toBe(767)
    expect(breakpoints.desktop).toBe(1024)
  })

  it('exports typography and radius scales', () => {
    expect(typographyScale.body).toBe(14)
    expect(radiusScale.card).toBe(12)
  })
})
