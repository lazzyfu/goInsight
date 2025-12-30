/**
 * 设置 localStorage 添加对时间的控制，hour单位为小时
 * @param key 保存在storage的key
 * @param data 需存储的数据
 * @param hour null时存sessionStorage(key,value)，即关闭浏览器过期
 * hour=0时，使用localStorage，即永不过期
 * hour>0时localStorage添加时间控制
 */

export function setStorage(key, data, hour) {
  let newData = data
  if (typeof data === 'object') {
    newData = JSON.stringify(data)
  }
  /*if (!data) {
    return
  }*/
  if (hour === 0) {
    window.localStorage.setItem(key, newData)
  } else if (hour && hour > 0) {
    const now = new Date()
    const valueDate = JSON.stringify({
      __value: data,
      __time: now.setSeconds(now.getSeconds() + hour * 3600)
    })
    window.localStorage.setItem(key, valueDate)
  } else {
    window.sessionStorage.setItem(key, newData)
  }
}

/**
 * 获取storage
 * @param key 保存时的key
 * @param hour 如果保存时使用了时间，则需要传true
 * @return 返回保存的值，过期后返回false,其他异常或不存在返回undefined
 */
export const getStorage = (key, hour) => {
  let data
  if (hour) {
    data = window.localStorage.getItem(key)
    try {
      data = JSON.parse(data)
      if (typeof data === 'object' && data.__time) {
        if (!data.__value) {
          data = undefined
        }
        // 使用了时间的
        // 在当前时间后，表示没过期
        if (new Date().getTime() < data.__time) {
          data = data.__value
        } else {
          // 过期了
          data = false
        }
      }
    } catch {
      /* empty */
      data = undefined
    }
  } else {
    //保存时没传时间的，存在session里
    data = window.sessionStorage.getItem(key)
  }
  try {
    return JSON.parse(data)
  } catch {
    return data
  }
}
/**
 * 移除storage
 * @param key 要移除的key
 * @param hour set时使用了hour，移除时则传true
 */
export const removeStorage = (key, hour) => {
  if (hour) {
    window.localStorage.removeItem(key)
  } else {
    window.sessionStorage.removeItem(key)
  }
}
