export const buildOrderQuery = ({ page, pageSize, search, progress, onlyMine }) => ({
  page,
  page_size: pageSize,
  is_page: true,
  search,
  progress,
  only_my_orders: onlyMine,
})

export const summarizeMyOrders = (rows = []) => ({
  total: rows.length,
  pending: rows.filter((item) => item.progress === 'PENDING').length,
  executing: rows.filter((item) => item.progress === 'EXECUTING').length,
  failed: rows.filter((item) => item.progress === 'FAILED').length,
})
