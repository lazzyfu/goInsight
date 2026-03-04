export const normalizeOtpCode = (value) => String(value || '').replace(/\D/g, '').slice(0, 6)

export const DEFAULT_AFTER_LOGIN_ROUTE = '/account/basic'

export const resolveLoginTarget = (redirect) => {
  if (typeof redirect !== 'string' || !redirect.trim()) {
    return DEFAULT_AFTER_LOGIN_ROUTE
  }

  let decodedPath = redirect
  try {
    decodedPath = decodeURIComponent(redirect)
  } catch {
    decodedPath = redirect
  }

  if (!decodedPath.startsWith('/') || decodedPath.startsWith('/login')) {
    return DEFAULT_AFTER_LOGIN_ROUTE
  }

  return decodedPath
}
