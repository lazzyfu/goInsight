import Moment from "moment"

const filters =  {
  fmtTime(val) {
    const dt = new Date(val)
    const y = dt.getFullYear()
    const m = (dt.getMonth() + 1 + '').padStart(2, '0')
    const d = (dt.getDate() + '').padStart(2, '0')
    const hh = (dt.getHours() + '').padStart(2, '0')
    const mm = (dt.getMinutes() + '').padStart(2, '0')
    const ss = (dt.getSeconds() + '').padStart(2, '0')

    return `${y}-${m}-${d} ${hh}:${mm}:${ss}`
  },
  
  fmtIsContainer(val) {
    return ['是', '否'][val]
  },

  ecsMemory(value) {
    if (value === 0) return '0'
    var k = 1024,
    sizes = ['MB', 'GB', 'TB'],
    i = Math.floor(Math.log(value) / Math.log(k))
    return (value / Math.pow(k, i)).toPrecision(3) + ' ' + sizes[i]
  },

  addZero: function(value){
    if (value === "" || !value){
       return "Unknown"
    } else if (typeof value === "string") {
      return parseFloat(value).toFixed(2)
    } else {
      return value.toFixed(2)
    }
  },

  sizeType: function (value){
    if (value === "" || !value) {
      return "Unknown"
    } else if (value === 0) return "0"
    const k = 1024
    const sizes = ["B", "Ki", "Mi", "Gi", "Ti", "Pi"]
    const i = Math.floor(Math.log(value) / Math.log(k))
    return (value / Math.pow(k, i)).toPrecision(3) + " " + sizes[i]
  },

  formatMilliSecond: function (value) {
    const date = new Date()
    date.setTime(value)
    return date.getMinutes() + ":" + date.getSeconds()
  },

  getCostInMillis: function (nanoSeconds) {
    return nanoSeconds / 1000000.0
  },

  formatDate: function (date) {
    return Moment(date).format('YYYY-MM-DD HH:mm')
  },

  formatDateToY: function (date) {
    return Moment(date).format('YYYY-MM-DD')
  },

  formatDateToYMDS: function (date) {
    return Moment(date).format('YYYY-MM-DD HH:mm:ss')
  },

  formatDateValue: function (date) {
    return Moment(date)
  },

  formatTimestampToDate: function (date) {
    return Moment(date * 1000).format('YYYY-MM-DD HH:mm:ss')
  },

  timeDifference: function(startTime, endTime) {
    const dateDiff = Moment(endTime).diff(startTime);
    if (dateDiff >= 1) {
        let result = ''
        const leave1 = dateDiff % (24 * 3600 * 1000);
        const hours = Math.floor(leave1 / (3600 * 1000));
        const leave2 = leave1 % (3600 * 1000);
        const minutes = Math.floor(leave2 / (60 * 1000));
        const leave3 = leave2 % (60 * 1000);
        const seconds = Math.round(leave3 / 1000);
        if (hours > 0) {
            result = result + hours + '时'
        }
        if (minutes > 0) {
            result = result + minutes + '分'
        }
        if (seconds > 0) {
            result = result + seconds + '秒'
        }
        return result
    } else {
        return '-'
    }
  },

  formatTime: function(millis) {
    const totalSeconds = Math.floor(millis / 1000)
    const minutes = Math.floor(totalSeconds / 60)
    const seconds = totalSeconds % 60
    const formattedMinutes = String(minutes).padStart(2, '0')
    const formattedSeconds = String(seconds).padStart(2, '0')
    return `${formattedMinutes}:${formattedSeconds}`
  }

}


export default filters