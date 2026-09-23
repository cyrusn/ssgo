<template>
  <div v-if="activeCohortID" class="card shadow-sm border-0 bg-light mb-5">
    <div class="card-header bg-secondary text-white border-0 py-3">
      <h5 class="m-0 fw-bold">4. 修改個別用戶資料</h5>
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

      <form @submit.prevent="updateUser">
        <div class="row">
          <div class="col-md-3 mb-3">
            <label class="form-label">用戶登入帳號</label>
            <div class="input-group">
              <input type="text" class="form-control" v-model="userModifier.userAlias" required placeholder="帳號" @blur="fetchDetails" />
              <button type="button" class="btn btn-outline-secondary" @click="fetchDetails">查詢</button>
            </div>
          </div>
          <div class="col-md-3 mb-3">
            <label class="form-label">英文姓名</label>
            <input type="text" class="form-control" v-model="userModifier.name" required />
          </div>
          <div class="col-md-3 mb-3">
            <label class="form-label">中文姓名</label>
            <input type="text" class="form-control" v-model="userModifier.cname" required />
          </div>
          <div class="col-md-3 mb-3">
            <label class="form-label">用戶權限 (Role)</label>
            <select class="form-select" v-model="userModifier.role" :disabled="userModifier.role === 'STUDENT'">
              <option value="STUDENT" v-if="userModifier.role === 'STUDENT'">STUDENT (學生)</option>
              <option value="TEACHER">TEACHER (教師)</option>
              <option value="ADMIN">ADMIN (管理員)</option>
            </select>
          </div>
        </div>
        <div class="row" v-if="userModifier.role === 'STUDENT'">
          <div class="col-md-6 mb-3">
            <label class="form-label">班別 (Class Code)</label>
            <input type="text" class="form-control" v-model="userModifier.classCode" />
          </div>
          <div class="col-md-6 mb-3">
            <label class="form-label">班號 (Class No)</label>
            <input type="number" class="form-control" v-model.number="userModifier.classNo" />
          </div>
        </div>
        <div class="row">
          <div class="col-md-12 mb-3">
            <label class="form-label">重設密碼 (不修改請留空)</label>
            <input type="password" class="form-control" v-model="userModifier.password" placeholder="輸入新密碼以重設" />
          </div>
        </div>
        <div class="d-flex gap-2">
          <button type="submit" class="btn btn-secondary fw-bold w-100 shadow-sm" :disabled="!userModifier.userAlias || submitting">更新用戶資料</button>
          <button type="button" class="btn btn-outline-danger fw-bold w-100 shadow-sm" :disabled="!userModifier.userAlias || submitting" @click="deleteUser">刪除用戶</button>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
export default {
  name: 'UserModifier',
  props: {
    activeCohortID: String,
    apiFetch: Function
  },
  data () {
    return {
      alert: '',
      error: '',
      submitting: false,
      userModifier: { userAlias: '', name: '', cname: '', password: '', classCode: '', classNo: 0, role: '' }
    }
  },
  methods: {
    async fetchDetails () {
      if (!this.userModifier.userAlias) return
      try {
        const res = await this.apiFetch(`./api/config/cohort/${this.activeCohortID}/user/${this.userModifier.userAlias}`)
        if (!res.ok) throw new Error('用戶不存在')
        const d = await res.json()
        this.userModifier.name = d.name
        this.userModifier.cname = d.cname
        this.userModifier.role = d.role
        this.userModifier.classCode = d.classCode
        this.userModifier.classNo = d.classNo
        this.userModifier.password = ''
        this.error = ''
      } catch (e) {
        this.error = e.message
      }
    },
    async updateUser () {
      this.submitting = true
      try {
        const res = await this.apiFetch(`./api/config/cohort/${this.activeCohortID}/user/${this.userModifier.userAlias}`, {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(this.userModifier)
        })
        if (!res.ok) throw new Error('更新失敗')
        this.alert = '已更新'
        this.userModifier.password = ''
      } catch (e) {
        this.error = e.message
      } finally {
        this.submitting = false
      }
    },
    async deleteUser () {
      if (!confirm(`確定刪除用戶 ${this.userModifier.userAlias}？`)) return
      this.submitting = true
      try {
        const res = await this.apiFetch(`./api/config/cohort/${this.activeCohortID}/user/${this.userModifier.userAlias}`, { method: 'DELETE' })
        if (!res.ok) throw new Error('刪除失敗')
        this.alert = '用戶已刪除'
        this.userModifier = { userAlias: '', name: '', cname: '', password: '', classCode: '', classNo: 0, role: '' }
      } catch (e) {
        this.error = e.message
      } finally {
        this.submitting = false
      }
    }
  }
}
</script>
