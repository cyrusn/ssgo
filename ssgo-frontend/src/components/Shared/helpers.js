import moment from 'moment'
import _ from 'lodash'

export function createDownloadFile (data, filename, extname) {
  const blob = new Blob([data], {
    type: 'text/plain;charset=UTF-8'
  })
  const link = document.createElement('a')

  link.href = URL.createObjectURL(blob)

  const today = moment().format('YYYY-MM-DD_hh-mm-ss')

  link.download = `${filename}-${today}.${extname}`
  link.click()
}

export function downloadJSON (json, filename, extname) {
  const jsonData = JSON.stringify(json)
  createDownloadFile(jsonData, filename, extname)
}

// removeCombinationWithX3Subject remove combination with X3 in priorities
export function removeCombinationWithX3Subject (priorities, x3, combinations) {
  const combinationWithX3 = _(combinations)
    .filter(comb => comb.subjects.includes(x3))
    .map('id')
    .value()

  return _.difference(priorities, combinationWithX3)
}
