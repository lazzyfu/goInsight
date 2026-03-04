const ORDER_STATUS_META = {
  PENDING: { text: '待审批', color: 'gold' },
  APPROVED: { text: '已批准', color: 'blue' },
  REJECTED: { text: '已驳回', color: 'red' },
  CLAIMED: { text: '已认领', color: 'cyan' },
  EXECUTING: { text: '执行中', color: 'orange' },
  FAILED: { text: '已失败', color: 'volcano' },
  COMPLETED: { text: '已完成', color: 'green' },
  REVIEWED: { text: '已复核', color: 'green' },
  REVOKED: { text: '已撤销', color: 'default' },
}

const TASK_STATUS_META = {
  PENDING: { text: '待执行', color: 'default' },
  EXECUTING: { text: '执行中', color: 'orange' },
  PAUSED: { text: '已暂停', color: 'default' },
  FAILED: { text: '已失败', color: 'red' },
  COMPLETED: { text: '已完成', color: 'green' },
}

export const ORDER_PROGRESS_OPTIONS = [
  { value: 'PENDING', label: '待审批' },
  { value: 'APPROVED', label: '已批准' },
  { value: 'REJECTED', label: '已驳回' },
  { value: 'CLAIMED', label: '已认领' },
  { value: 'EXECUTING', label: '执行中' },
  { value: 'FAILED', label: '已失败' },
  { value: 'COMPLETED', label: '已完成' },
  { value: 'REVIEWED', label: '已复核' },
  { value: 'REVOKED', label: '已撤销' },
]

export const TASK_PROGRESS_OPTIONS = [
  { value: 'PENDING', label: '待执行' },
  { value: 'EXECUTING', label: '执行中' },
  { value: 'COMPLETED', label: '已完成' },
  { value: 'FAILED', label: '已失败' },
]

export const getOrderStatusMeta = (status) =>
  ORDER_STATUS_META[status] || { text: status, color: 'default' }

export const getTaskStatusMeta = (status) =>
  TASK_STATUS_META[status] || { text: status, color: 'default' }
