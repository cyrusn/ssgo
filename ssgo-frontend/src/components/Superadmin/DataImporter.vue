<template>
  <div v-if="activeCohortID" class="card shadow-sm border-0 mb-5 bg-light">
    <div class="card-header bg-info text-white border-0 py-3 d-flex justify-content-between align-items-center">
      <h5 class="m-0 fw-bold">3. 批次匯入數據 (CSV)</h5>
    </div>
    <div class="card-body">
      <div v-if="alert" class="alert alert-info alert-dismissible fade show mb-4" role="alert">
        {{ alert }}
        <button type="button" class="btn-close" @click="alert = ''"></button>
      </div>
      <div v-if="error" class="alert alert-danger alert-dismissible fade show mb-4" role="alert">
        {{ error }}
        <button type="button" class="btn-close" @click="error = ''"></button>
      </div>

      <div class="row">
        <div class="col-md-4 mb-3">
          <label class="form-label">數據類型</label>
          <select class="form-select" v-model="importer.dataType" required>
            <option value="">-- 請選擇匯入類型 --</option>
            <option value="student">學生資料 (Student CSV)</option>
            <option value="teacher">教師資料 (Teacher CSV)</option>
            <option value="admin">管理員資料 (Admin CSV)</option>
          </select>
        </div>
        <div v-if="importer.dataType" class="col-md-8 mb-3">
          <div class="d-flex justify-content-between align-items-center mb-1">
            <label class="form-label m-0 fw-bold">貼上 CSV 或 JSON</label>
            <button class="btn btn-sm btn-outline-secondary" @click="fillSampleData">Example</button>
          </div>
          <textarea class="form-control font-monospace small" style="height:180px;" v-model="importer.rawJSON" placeholder="請在此處貼上..."></textarea>
        </div>
      </div>
      <button class="btn btn-info text-white w-100 fw-bold shadow-sm" :disabled="!importer.dataType || !importer.rawJSON || importing" @click="importData">匯入數據</button>
    </div>
  </div>
</template>

<script>
import Papa from 'papaparse'

export default {
  name: 'DataImporter',
  props: {
    activeCohortID: String,
    apiFetch: Function
  },
  data () {
    return {
      alert: '',
      error: '',
      importing: false,
      importer: { dataType: '', rawJSON: '' }
    }
  },
  methods: {
    async importData () {
      this.importing = true
      this.error = ''
      this.alert = ''
      try {
        const data = Papa.parse(this.importer.rawJSON, { header: true, skipEmptyLines: true, dynamicTyping: true }).data
        const res = await this.apiFetch(`./api/config/cohort/${this.activeCohortID}/import/${this.importer.dataType}`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(data)
        })
        if (!res.ok) {
          const msg = await res.text()
          throw new Error(msg || '匯入失敗')
        }
        this.alert = '數據已匯入'
        this.importer.rawJSON = ''
      } catch (e) {
        this.error = e.message
      } finally {
        this.importing = false
      }
    },
    fillSampleData () {
      const t = this.importer.dataType
      const now = new Date()
      const y = now.getFullYear()
      const m = now.getMonth() + 1
      let s, e
      if (m >= 9) { s = y; e = (y + 1) % 100 } else { s = y - 1; e = y % 100 }
      const schoolYear = `${s}-${String(e).padStart(2, '0')}`
      
      if (t === 'student') this.importer.rawJSON = `userAlias,name,cname,password,classCode,classNo\ns${schoolYear.replace('-', '')}01,CHAN TAI MAN,陳大文,student_password,4A,1`
      else if (t === 'teacher') this.importer.rawJSON = `userAlias,name,cname,password\nt01,WONG SIU MING,黃小明,teacher_password`
      else if (t === 'admin') this.importer.rawJSON = `userAlias,name,cname,password\nadmin01,ADMIN USER,管理員,admin_password`
    }
  }
}
</script>
