import { describe, expect, it } from 'vitest'
import { normalizeOtpCode } from '../loginModel'

describe('login otp model', () => {
  it('keeps only 6 digits', () => {
    expect(normalizeOtpCode('a1b2c3d4')).toBe('1234')
    expect(normalizeOtpCode('123456789')).toBe('123456')
  })

  it('handles empty values', () => {
    expect(normalizeOtpCode()).toBe('')
  })
})
