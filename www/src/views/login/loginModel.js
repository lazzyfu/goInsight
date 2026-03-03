export const normalizeOtpCode = (value) => String(value || '').replace(/\D/g, '').slice(0, 6)
