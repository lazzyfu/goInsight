export const normalizeClaimUsers = (raw) => {
  if (!raw) return '无'

  if (Array.isArray(raw)) {
    return raw.length ? raw.join(', ') : '无'
  }

  if (typeof raw === 'string') {
    try {
      const users = JSON.parse(raw)
      if (Array.isArray(users)) {
        return users.length ? users.join(', ') : '无'
      }
    } catch {
      return raw.trim() || '无'
    }
    return '无'
  }

  return '无'
}
