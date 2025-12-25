const regPassword = /^.*(?=.{7,})(?=.*\d)(?=.*[A-Z])(?=.*[a-z])(?=.*[!@#$%^&*? ]).*$/
const regEmail = /^\w+([-+.]?\w+)*@\w+([-.]?\w+)*\.\w+([-.]?)/;
const regPhone = /^1[3-9]\d{9}$/;;

export { regEmail, regPassword, regPhone };

