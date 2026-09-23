<template>
  <div class="card shadow-sm border-0 mb-5 bg-light">
    <div class="card-header bg-primary text-white border-0 py-3 d-flex justify-content-between align-items-center">
      <h5 class="m-0 fw-bold">1. 年度選科資料庫管理 (Cohorts)</h5>
      <button class="btn btn-sm btn-light fw-bold" @click="openCreateModal">+ 新增年度</button>
    </div>
    <div class="card-body p-0">
      <div v-if="alert" class="alert alert-info alert-dismissible fade show m-3" role="alert">
        {{ alert }}
        <button type="button" class="btn-close" @click="alert = ''"></button>
      </div>
      <div v-if="error" class="alert alert-danger alert-dismissible fade show m-3" role="alert">
        {{ error }}
        <button type="button" class="btn-close" @click="error = ''"></button>
      </div>

      <div class="table-responsive">
        <table class="table table-hover align-middle mb-0">
          <thead class="table-light">
            <tr>
              <th class="ps-4">ID</th>
              <th>顯示名稱</th>
              <th class="text-center">狀態</th>
              <th class="text-end pe-4">管理操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="c in cohorts" :key="c.id || c.ID" :class="(c.isActive || c.IsActive) ? 'table-primary fw-bold' : ''">
              <td class="ps-4"><code>{{ c.id || c.ID }}</code></td>
              <td>{{ c.name || c.Name }}</td>
              <td class="text-center">
                <span v-if="c.isActive || c.IsActive" class="badge bg-success px-3 py-2">已啟用</span>
                <span v-else class="badge bg-secondary px-3 py-2">未啟用</span>
              </td>
              <td class="text-end pe-4">
                <div class="btn-group">
                  <button v-if="!(c.isActive || c.IsActive)" class="btn btn-sm btn-success" @click="setActive(c.id || c.ID)">啟用</button>
                  <button v-if="c.isActive || c.IsActive" class="btn btn-sm btn-warning" @click="deactivate(c.id || c.ID)">停用</button>
                  <button class="btn btn-sm btn-danger" @click="deleteCohort(c.id || c.ID)">刪除</button>
                </div>
              </td>
            </tr>
            <tr v-if="cohorts.length === 0">
              <td colspan="4" class="text-center text-secondary py-4">目前沒有任何年度數據庫。</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- CREATE MODAL -->
    <div v-if="showModal" class="modal-backdrop fade show" style="background-color: rgba(0,0,0,0.5);"></div>
    <div v-if="showModal" class="modal fade show d-block" tabindex="-1" style="top: 15%;">
      <div class="modal-dialog">
        <div class="modal-content border-0 shadow-lg">
          <div class="modal-header bg-primary text-white">
            <h5 class="modal-title fw-bold">新增年度</h5>
            <button type="button" class="btn-close btn-close-white" @click="showModal = false"></button>
          </div>
          <div class="modal-body">
            <div class="mb-3">
              <label class="form-label">年度 ID</label>
              <input type="text" class="form-control" v-model="newCohort.id" required />
            </div>
            <div class="mb-3">
              <label class="form-label">年度名稱</label>
              <input type="text" class="form-control" v-model="newCohort.name" required />
            </div>
          </div>
          <div class="modal-footer bg-light">
            <button type="button" class="btn btn-secondary" @click="showModal = false">取消</button>
            <button type="button" class="btn btn-primary px-4" @click="createCohort" :disabled="submitting">建立年度</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'CohortManager',
  props: {
    cohorts: { type: Array, required: true },
    jwt: { type: String, required: true },
    apiFetch: { type: Function, required: true }
  },
  data () {
    return {
      alert: '',
      error: '',
      showModal: false,
      submitting: false,
      newCohort: { id: '', name: '' }
    }
  },
  methods: {
    openCreateModal () {
      const now = new Date()
      const y = now.getFullYear()
      const m = now.getMonth() + 1
      let s, e
      if (m >= 9) { s = y; e = (y + 1) % 100 } else { s = y - 1; e = y % 100 }
      const schoolYear = `${s}-${String(e).padStart(2, '0')}`
      this.newCohort.id = schoolYear + '_Mock'
      this.newCohort.name = schoolYear + ' 模擬選科'
      this.showModal = true
    },
    async createCohort () {
      this.submitting = true
      try {
        const res = await this.apiFetch('./api/config/cohorts', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(this.newCohort)
        })
        if (!res.ok) throw new Error(await res.text())
        this.alert = '已建立'
        this.showModal = false
        this.$emit('refresh')
      } catch (e) {
        this.error = e.message
      } finally {
        this.submitting = false
      }
    },
    async setActive (id) {
      try {
        const res = await this.apiFetch(`./api/config/cohort/${id}/active`, { method: 'PUT' })
        if (!res.ok) throw new Error('失敗')
        this.alert = '已啟用'
        this.$emit('refresh')
      } catch (e) {
        this.error = e.message
      }
    },
    async deactivate (id) {
      try {
        const res = await this.apiFetch(`./api/config/cohort/${id}/deactivate`, { method: 'PUT' })
        if (!res.ok) throw new Error('失敗')
        this.alert = '已停用'
        this.$emit('refresh')
      } catch (e) {
        this.error = e.message
      }
    },
    async deleteCohort (id) {
      const input = prompt(`輸入 "${id}" 確定：`)
      if (input === id) {
        try {
          const res = await this.apiFetch(`./api/config/cohort/${id}`, { method: 'DELETE' })
          if (!res.ok) throw new Error('刪除失敗')
          this.alert = '已刪除'
          this.$emit('refresh')
        } catch (e) {
          this.error = e.message
        }
      }
    }
  }
}
</script>
