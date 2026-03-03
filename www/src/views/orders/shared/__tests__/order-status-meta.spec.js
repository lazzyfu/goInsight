import { describe, expect, it } from 'vitest'
import { getOrderStatusMeta } from '../orderStatusMeta'

describe('order status meta', () => {
  it('returns consistent text and color for reviewed status', () => {
    expect(getOrderStatusMeta('REVIEWED')).toEqual({
      text: '已复核',
      color: 'green',
    })
  })

  it('falls back to raw text for unknown status', () => {
    expect(getOrderStatusMeta('UNKNOWN').text).toBe('UNKNOWN')
  })
})
