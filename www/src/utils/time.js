import Moment from "moment"

const formatDate = (date) => {
    return Moment(date).format('YYYY-MM-DD HH:mm')
}

const formatDateToY = (date) => {
    return Moment(date).format('YYYY-MM-DD')
}

const formatDateToYMDS = (date) => {
    return Moment(date).format('YYYY-MM-DD HH:mm:ss')
}

const formatDateValue = (date) => {
    return Moment(date)
}

const formatTimestampToDate = (date) => {
    return Moment(date * 1000).format('YYYY-MM-DD HH:mm:ss')
}

const timeDifference = (startTime, endTime) => {
    const dateDiff = Moment(endTime).diff(startTime);
    if (dateDiff >= 1) {
        let result = ''
        const leave1 = dateDiff % (24 * 3600 * 1000)
        const hours = Math.floor(leave1 / (3600 * 1000))
        const leave2 = leave1 % (3600 * 1000)
        const minutes = Math.floor(leave2 / (60 * 1000))
        const leave3 = leave2 % (60 * 1000)
        const seconds = Math.round(leave3 / 1000)

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
}

const formatTime = (millis) => {
    const totalSeconds = Math.floor(millis / 1000)
    const minutes = Math.floor(totalSeconds / 60)
    const seconds = totalSeconds % 60

    const formattedMinutes = String(minutes).padStart(2, '0')
    const formattedSeconds = String(seconds).padStart(2, '0')

    return `${formattedMinutes}:${formattedSeconds}`
}


export default {
    formatDate,
    formatDateToY,
    formatDateToYMDS,
    formatDateValue,
    formatTimestampToDate,
    timeDifference,
    formatTime,
}