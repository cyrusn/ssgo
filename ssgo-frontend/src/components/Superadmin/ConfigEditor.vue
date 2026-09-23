<template>
  <div v-if="activeCohortID && selectedCohort" class="card shadow-sm border-0 bg-light mb-5">
    <div class="card-header bg-dark text-white border-0 py-3 d-flex justify-content-between align-items-center">
      <div>
        <h5 class="m-0 fw-bold d-inline-block me-2">2. 設定前端參數</h5>
        <small class="text-light opacity-75">年度：<code>{{ activeCohortID }}</code></small>
      </div>
      <button class="btn btn-sm btn-outline-light" @click="isCollapsed = !isCollapsed">{{ isCollapsed ? '展開設定' : '收合' }}</button>
    </div>
    <div class="card-body" v-show="!isCollapsed">
      <div v-if="alert" class="alert alert-info alert-dismissible fade show mb-4" role="alert">
        {{ alert }}
        <button type="button" class="btn-close" @click="alert = ''"></button>
      </div>
      <div v-if="error" class="alert alert-danger alert-dismissible fade show mb-4" role="alert">
        {{ error }}
        <button type="button" class="btn-close" @click="error = ''"></button>
      </div>

      <div class="mb-4 p-3 border rounded bg-white shadow-sm">
        <h6 class="fw-bold border-bottom pb-2 mb-3">複製設定</h6>
        <div class="row align-items-end">
          <div class="col-md-8">
            <label class="form-label small">從其他年度複製所有 Markdown、科目及組合設定：</label>
            <select class="form-select" v-model="targetCohortToCopy">
              <option value="">-- 請選擇來源年度 --</option>
              <option v-for="c in cohorts" :key="c.id || c.ID" :value="c.id || c.ID" v-show="(c.id || c.ID) !== activeCohortID">{{ c.name || c.Name }}</option>
            </select>
          </div>
          <div class="col-md-4">
            <button class="btn btn-outline-dark w-100" @click="copyConfigFromOtherCohort" :disabled="!targetCohortToCopy">執行複製</button>
          </div>
        </div>
      </div>

      <form @submit.prevent="saveConfig">
        <div class="mb-4 p-3 border rounded bg-white">
          <h6 class="fw-bold border-bottom pb-2 mb-3">基本系統設定</h6>
          <div class="mb-3">
            <label class="form-label fw-bold text-primary">年度數據庫顯示名稱</label>
            <input type="text" class="form-control fw-bold border-primary" v-model="selectedCohortName" required />
          </div>
          <div class="row">
            <div class="col-md-4 mb-3">
              <label class="form-label">學校名稱</label>
              <input type="text" class="form-control" v-model="editorConfig.schoolName" required />
            </div>
            <div class="col-md-4 mb-3">
              <label class="form-label">學校網址</label>
              <input type="url" class="form-control" v-model="editorConfig.schoolWebsite" required />
            </div>
            <div class="col-md-4 mb-3">
              <label class="form-label fw-bold">模擬選科 (isMock)?</label>
              <select class="form-select border-info" v-model="editorConfig.isMock">
                <option :value="true">是 (模擬選科)</option>
                <option :value="false">否 (正式選科)</option>
              </select>
            </div>
          </div>
          <div class="mb-3">
            <label class="form-label">系統標題</label>
            <input type="text" class="form-control" v-model="editorConfig.systemTitle" required />
          </div>
          <div class="mb-3">
            <label class="form-label">接受選科表提交?</label>
            <select class="form-select" v-model="editorConfig.isAcceptResponses">
              <option :value="true">開啟</option>
              <option :value="false">關閉</option>
            </select>
          </div>
        </div>

        <div class="mb-4 p-3 border rounded bg-white">
          <h6 class="fw-bold border-bottom pb-2 mb-3">學生介面內容設定 (Markdown)</h6>
          
          <div class="mb-4">
            <div class="d-flex justify-content-between align-items-center mb-1">
              <label class="form-label fw-bold m-0">1. 頁面簡介 (Introduction)</label>
              <button type="button" class="btn btn-sm btn-outline-secondary" @click="fillExample('intro')">Example</button>
            </div>
            <textarea class="form-control font-monospace" rows="10" v-model="editorConfig.introMarkdown" required></textarea>
          </div>

          <div class="mb-4">
            <div class="d-flex justify-content-between align-items-center mb-1">
              <label class="form-label fw-bold m-0">2. 使用須知 (Instruction)</label>
              <button type="button" class="btn btn-sm btn-outline-secondary" @click="fillExample('instr')">Example</button>
            </div>
            <textarea class="form-control font-monospace" rows="5" v-model="editorConfig.instrMarkdown" required></textarea>
          </div>

          <div class="mb-4">
            <div class="d-flex justify-content-between align-items-center mb-1">
              <label class="form-label fw-bold m-0">3. 截止提示 (Offline Message)</label>
              <button type="button" class="btn btn-sm btn-outline-secondary" @click="fillExample('notAccept')">Example</button>
            </div>
            <textarea class="form-control font-monospace" rows="4" v-model="editorConfig.notAcceptMarkdown" required></textarea>
          </div>

          <div class="mb-4">
            <div class="d-flex justify-content-between align-items-center mb-1">
              <label class="form-label fw-bold m-0">4. 確認提交提示 (Confirm Prompt)</label>
              <button type="button" class="btn btn-sm btn-outline-secondary" @click="fillExample('confirm')">Example</button>
            </div>
            <textarea class="form-control font-monospace" rows="4" v-model="editorConfig.confirmMarkdown" required></textarea>
          </div>
        </div>

        <div class="mb-4 p-3 border rounded bg-white">
          <h6 class="fw-bold border-bottom pb-2 mb-3">核心數據設定 (CSV)</h6>
          <div class="mb-4">
            <div class="d-flex justify-content-between align-items-center mb-1">
              <label class="form-label fw-bold m-0">5. 科目清單 (Subjects CSV)</label>
              <button type="button" class="btn btn-sm btn-outline-secondary" @click="fillExample('subjects')">Example</button>
            </div>
            <textarea class="form-control font-monospace" rows="8" v-model="subjectsCSV" required></textarea>
          </div>
          <div class="mb-4">
            <div class="d-flex justify-content-between align-items-center mb-1">
              <label class="form-label fw-bold m-0">6. 選科組合 (Combinations CSV)</label>
              <button type="button" class="btn btn-sm btn-outline-secondary" @click="fillExample('combinations')">Example</button>
            </div>
            <textarea class="form-control font-monospace" rows="8" v-model="combinationsCSV" required></textarea>
          </div>
        </div>
        <button type="submit" class="btn btn-primary btn-lg w-100 shadow-sm" :disabled="submitting">儲存所有設定及同步科目</button>
      </form>
    </div>
  </div>
</template>

<script>
import Papa from 'papaparse'

export default {
  name: 'ConfigEditor',
  props: {
    cohorts: Array,
    activeCohortID: String,
    apiFetch: Function,
    jwt: String
  },
  data () {
    return {
      alert: '',
      error: '',
      isCollapsed: true,
      selectedCohort: null,
      selectedCohortName: '',
      targetCohortToCopy: '',
      submitting: false,
      subjectsCSV: '',
      combinationsCSV: '',
      systemDefaults: null,
      editorConfig: {
        schoolName: '', schoolWebsite: '', systemTitle: '', isMock: false, isAcceptResponses: true,
        introMarkdown: '', instrMarkdown: '', notAcceptMarkdown: '', confirmMarkdown: '',
        subjectsJSON: '[]', combinationsJSON: '[]'
      }
    }
  },
  watch: {
    activeCohortID: {
      immediate: true,
      handler (newID) {
        if (newID && this.cohorts) {
          const cohort = this.cohorts.find(c => (c.id || c.ID) === newID)
          if (cohort) this.selectCohort(cohort)
        }
      }
    }
  },
  methods: {
    selectCohort (c) {
      this.selectedCohort = c
      this.selectedCohortName = c.name || c.Name
      try {
        const parsed = JSON.parse(c.config || c.Config)
        this.editorConfig = Object.assign({}, this.editorConfig, parsed)
        if (parsed.introductionMarkdown) this.editorConfig.introMarkdown = parsed.introductionMarkdown
        if (parsed.instructionMarkdown) this.editorConfig.instrMarkdown = parsed.instructionMarkdown
        this.jsonToCSV()
      } catch (e) {
        this.error = '解析失敗'
      }
    },
    async saveConfig () {
      this.submitting = true
      this.csvToJSON()
      try {
        const res = await this.apiFetch(`./api/config/cohort/${this.activeCohortID}/config`, {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ name: this.selectedCohortName, config: JSON.stringify(this.editorConfig) })
        })
        if (!res.ok) throw new Error('儲存失敗')
        this.alert = '已儲存'
        this.syncSubjects()
        this.$emit('refresh')
      } catch (e) {
        this.error = e.message
      } finally {
        this.submitting = false
      }
    },
    copyConfigFromOtherCohort () {
      const s = this.cohorts.find(c => (c.id || c.ID) === this.targetCohortToCopy)
      if (!s) return
      try {
        const sc = JSON.parse(s.config || s.Config)
        const fields = ['schoolName', 'schoolWebsite', 'systemTitle', 'isMock', 'introMarkdown', 'instrMarkdown', 'notAcceptMarkdown', 'confirmMarkdown', 'subjectsJSON', 'combinationsJSON']
        fields.forEach(f => { if (sc[f] !== undefined) this.editorConfig[f] = sc[f] })
        this.jsonToCSV()
        this.alert = '已從 ' + (s.name || s.Name) + ' 複製'
      } catch (e) {
        this.error = '複製失敗'
      }
    },
    syncSubjects () {
      try {
        const s = JSON.parse(this.editorConfig.subjectsJSON)
        const codes = s.map(x => x.code).filter(Boolean)
        if (codes.length) {
          this.apiFetch(`./api/config/cohort/${this.activeCohortID}/import/subject`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(codes)
          })
        }
      } catch (e) { console.error(e) }
    },
    jsonToCSV () {
      try {
        const s = JSON.parse(this.editorConfig.subjectsJSON || '[]')
        this.subjectsCSV = Papa.unparse(s)
        const c = JSON.parse(this.editorConfig.combinationsJSON || '[]')
        const comboRows = c.map(x => ({
          id: x.id,
          subject1: (x.subjects && x.subjects[0]) ? x.subjects[0] : '',
          subject2: (x.subjects && x.subjects[1]) ? x.subjects[1] : ''
        }))
        this.combinationsCSV = Papa.unparse(comboRows)
      } catch (e) { console.error('jsonToCSV error', e) }
    },
    csvToJSON () {
      try {
        const s = Papa.parse(this.subjectsCSV, { header: true, skipEmptyLines: true, dynamicTyping: true })
        this.editorConfig.subjectsJSON = JSON.stringify(s.data)
        const c = Papa.parse(this.combinationsCSV, { header: true, skipEmptyLines: true, dynamicTyping: true })
        const combos = c.data.map(x => ({ id: x.id, subjects: [x.subject1, x.subject2] }))
        this.editorConfig.combinationsJSON = JSON.stringify(combos)
      } catch (e) { console.error('csvToJSON error', e) }
    },
    async fillExample (type) {
      if (!this.systemDefaults) {
        const res = await this.apiFetch('./api/config/system-defaults')
        this.systemDefaults = await res.json()
      }
      const d = this.systemDefaults || {}
      const isMock = this.editorConfig.isMock
      switch (type) {
        case 'intro': this.editorConfig.introMarkdown = isMock ? d.mockIntroMarkdown : d.introMarkdown; break
        case 'instr': this.editorConfig.instrMarkdown = d.instrMarkdown; break
        case 'notAccept': this.editorConfig.notAcceptMarkdown = isMock ? d.mockNotAcceptMarkdown : d.notAcceptMarkdown; break
        case 'confirm': this.editorConfig.confirmMarkdown = d.confirmMarkdown; break
        case 'subjects': this.editorConfig.subjectsJSON = d.subjectsJSON; this.jsonToCSV(); break
        case 'combinations': this.editorConfig.combinationsJSON = d.combinationsJSON; this.jsonToCSV(); break
      }
      this.alert = '已帶入範例內容'
    }
  }
}
</script>
