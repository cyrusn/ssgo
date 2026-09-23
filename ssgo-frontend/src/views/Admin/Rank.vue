<template>
  <div class="my-4">
    <rank-example />
    <div
      v-if="successUpdate"
      class="mt-2 fade show alert alert-primary"
      id="successUpdate"
    >
      上載成功
    </div>
    <div
      v-else-if="errorMessage"
      class="mt-2 alert alert-danger"
      v-html="errorMessage"
    ></div>
    <div v-else class="mt-2 alert alert-success">數據正確，可以上載。</div>
    <form class="mt-2">
      <div class="form-group mb-2">
        <textarea rows="5" class="form-control" v-model="csvData"></textarea>
      </div>
      <div class="btn-group" role="group">
        <button
          class="btn btn-warning"
          type="button"
          data-bs-toggle="modal"
          data-bs-target="#exampleRank"
        >
          顯示例子
        </button>
        <button
          type="button"
          class="btn btn-primary"
          @click="update"
          :disabled="errorMessage !== ''"
        >
          上載
        </button>
      </div>
    </form>
  </div>
</template>

<script>
import Papa from 'papaparse'
import _ from 'lodash'
import { mapActions } from 'vuex'

import RankExample from '@/components/Admin/RankExample'

export default {
  components: {
    RankExample
  },
  data () {
    return {
      csvData: '',
      successUpdate: false
    }
  },
  computed: {
    jsonData () {
      const { csvData } = this
      if (csvData) {
        return Papa.parse(csvData, {
          header: true,
          encoding: 'UTF-8'
        })
      }
      return null
    },
    errorMessage () {
      const { jsonData } = this
      if (!jsonData) {
        return '請複製試算表中的名次，並貼上到以下文字方框中，完成後並按上載。<br>注意：試算表中的標題列必須為 <code>userAlias</code> 及 <code>rank</code>。'
      }

      if (jsonData.errors.length !== 0) {
        return `<li>${_(jsonData.errors)
          .map(e => e.message)
          .uniq()
          .value()
          .join('</li><li>')}</li>`
      }

      if (
        _.intersection(jsonData.meta.fields, ['userAlias', 'rank']).length !== 2
      ) {
        return '格式不正確。請確定標題分別為 <code>userAlias</code> 及 <code>rank</code>'
      }

      if (
        _.includes(
          jsonData.data.map(s => parseInt(s.rank)),
          NaN
        )
      ) {
        return '<code>rank</code> 必須為整數'
      }

      return ''
    }
  },
  methods: {
    ...mapActions('students', ['updateRanks']),
    update () {
      const vm = this
      const { jsonData, updateRanks } = this

      if (jsonData.errors.length === 0) {
        const postJSON = jsonData.data.map(o => {
          const { userAlias, rank } = o
          return { userAlias, rank: parseInt(rank) }
        })
        updateRanks(postJSON)
        this.successUpdate = true
        setTimeout(function () {
          vm.successUpdate = false
          vm.csvData = ''
        }, 3000)
      }
    }
  }
}
</script>
