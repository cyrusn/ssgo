<template>
  <span>{{ formattedDatetime }}</span>
</template>

<script>
import moment from 'moment'

moment.locale('zh-hk', {
  longDateFormat: {
    LT: 'Ah時m分',
    LTS: 'Ah時m分s秒',
    l: 'M月D日',
    lll: 'YYYY年M月D日Ah時m分',
    llll: 'YYYY年M月D日Ah時m分s秒'
  },
  meridiem: function (hour, minute) {
    const hm = hour * 100 + minute
    if (hm < 600) {
      return '凌晨'
    } else if (hm < 900) {
      return '早上'
    } else if (hm < 1130) {
      return '上午'
    } else if (hm === 1200) {
      return '正午'
    } else if (hm < 1230) {
      return '中午'
    } else if (hm < 1800) {
      return '下午'
    } else {
      return '晚上'
    }
  }
})

const ChineseNumbers = [
  '零',
  '一',
  '二',
  '三',
  '四',
  '五',
  '六',
  '七',
  '八',
  '九'
]
const ChineseNumbersEmptyZero = [
  '',
  '一',
  '二',
  '三',
  '四',
  '五',
  '六',
  '七',
  '八',
  '九'
]
const CHINESE_TEN = '十'

function cYear (year) {
  return (
    [].map
      .call(year.toString(), letter => ChineseNumbers[parseInt(letter, 10)])
      .join('') + '年'
  )
}

function cMonth (month) {
  const unitsDigit = (month + 1) % 10
  return (
    (month + 1 >= 10 ? CHINESE_TEN : '') +
    ChineseNumbersEmptyZero[unitsDigit] +
    '月'
  )
}

function cDate (day) {
  const tensDigit = Math.floor(day / 10)
  const unitsDigit = day % 10
  let result = ''

  if (tensDigit > 0) {
    result += (tensDigit > 1 ? ChineseNumbers[tensDigit] : '') + CHINESE_TEN
  }
  result += ChineseNumbersEmptyZero[unitsDigit]
  result += '日'
  return result
}

export default {
  props: ['datetime', 'format'],
  computed: {
    formattedDatetime () {
      const { datetime, format } = this
      const d = moment(datetime)

      switch (format) {
        case 'L':
          return `${cMonth(d.month())}${cDate(d.date())}`
        case 'LL':
          return `${cYear(d.year())}${cMonth(d.month())}${cDate(d.date())}`
        case 'LLL':
          return `${cMonth(d.month())}${cDate(d.date())}${d.format('LT')}`
        case 'LLLL':
          return `${cYear(d.year())}${cMonth(d.month())}${cDate(
            d.date()
          )}${d.format('LTS')}`
        case 'YYYY':
          return `${cYear(d.year())}`
        default:
          return d.format(format)
      }
    }
  }
}
</script>
