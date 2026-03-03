import { describe, expect, it } from 'vitest'
import { buildOrderQuery, summarizeMyOrders } from '../orderListModel'

describe('order list model', () => {
  it('builds api query from ui state', () => {
    expect(
      buildOrderQuery({
        page: 2,
        pageSize: 20,
        search: 'ddl',
        progress: 'PENDING',
        onlyMine: true,
      }),
    ).toEqual({
      page: 2,
      page_size: 20,
      is_page: true,
      search: 'ddl',
      progress: 'PENDING',
      only_my_orders: true,
    })
  })

  it('summarizes my orders', () => {
    const rows = [{ progress: 'PENDING' }, { progress: 'EXECUTING' }, { progress: 'FAILED' }]
    expect(summarizeMyOrders(rows)).toEqual({
      total: 3,
      pending: 1,
      executing: 1,
      failed: 1,
    })
  })
})
