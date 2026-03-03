import { describe, expect, it } from 'vitest'
import { normalizeClaimUsers } from '../orderDetailModel'

describe('order detail model', () => {
  it('parses json claim_users string', () => {
    expect(normalizeClaimUsers('["alice","bob"]')).toBe('alice, bob')
  })

  it('returns fallback for empty values', () => {
    expect(normalizeClaimUsers('')).toBe('无')
    expect(normalizeClaimUsers(null)).toBe('无')
  })
})
