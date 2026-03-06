const regPassword = /^.*(?=.{7,})(?=.*\d)(?=.*[A-Z])(?=.*[a-z])(?=.*[!@#$%^&*? ]).*$/
const regEmail = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
const regPhone = /^1[3-9]\d{9}$/

export { regEmail, regPassword, regPhone }
