<template>
  <div class="card shadow-sm border-0 mb-5 bg-light">
    <div class="card-header bg-danger text-white border-0 py-3">
      <h5 class="m-0 fw-bold">5. 系統管理員帳號管理</h5>
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
        <div class="col-md-6 border-end">
          <h6 class="fw-bold mb-3 border-bottom pb-2">現有管理員</h6>
          <ul class="list-group list-group-flush">
            <li v-for="admin in superadmins" :key="admin" class="list-group-item d-flex justify-content-between align-items-center">
              <code>{{ admin }}</code>
              <button v-if="admin === 'root' && userAlias === 'root'" class="btn btn-sm btn-warning" @click="resetRoot">恢復預設</button>
            </li>
          </ul>
        </div>
        <div class="col-md-6">
          <h6 class="fw-bold mb-3 border-bottom pb-2">新增管理員</h6>
          <form @submit.prevent="createAdmin">
            <input type="text" class="form-control mb-2" v-model="newAdmin.username" placeholder="用戶名" required />
            <input type="password" class="form-control mb-2" v-model="newAdmin.password" placeholder="密碼" required />
            <button type="submit" class="btn btn-danger w-100 shadow-sm" :disabled="submitting">建立管理員</button>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'AdminAccountManager',
  props: {
    userAlias: String,
    apiFetch: Function
  },
  data () {
    return {
      alert: '',
      error: '',
      submitting: false,
      superadmins: [],
      newAdmin: { username: '', password: '' }
    }
  },
  mounted () {
    this.loadAdmins()
  },
  methods: {
    async loadAdmins () {
      try {
        const res = await this.apiFetch('./api/config/manage/admins')
        this.superadmins = await res.json()
      } catch (e) {
        console.error('Admins load failed', e)
      }
    },
    async createAdmin () {
      this.submitting = true
      try {
        const res = await this.apiFetch('./api/config/manage/admins', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(this.newAdmin)
        })
        if (!res.ok) throw new Error('建立失敗')
        this.alert = '已建立'
        this.newAdmin = { username: '', password: '' }
        this.loadAdmins()
      } catch (e) {
        this.error = e.message
      } finally {
        this.submitting = false
      }
    },
    async resetRoot () {
      if (!confirm('確定恢復 root 密碼為預設？')) return
      try {
        const res = await this.apiFetch('./api/config/manage/root/reset', { method: 'POST' })
        if (!res.ok) throw new Error('恢復失敗')
        this.alert = 'root 密碼已恢復'
      } catch (e) {
        this.error = e.message
      }
    }
  }
}
</script>
