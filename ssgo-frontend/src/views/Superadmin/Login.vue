<template>
  <div class="container d-flex align-items-center justify-content-center min-vh-100">
    <div class="card shadow p-4 border-0 bg-light" style="width: 100%; max-width: 450px;">
      <h2 class="text-center fw-bold mb-4">系統管理員登入</h2>
      <hr />
      
      <div v-if="error" class="alert alert-danger" role="alert">
        {{ error }}
      </div>

      <form @submit.prevent="onLogin">
        <div class="form-group mb-3">
          <label class="form-label">管理員用戶名 (Username)</label>
          <input
            type="text"
            class="form-control form-control-lg"
            v-model="username"
            required
            autofocus
          />
        </div>
        
        <div class="form-group mb-4">
          <label class="form-label">密碼 (Password)</label>
          <input
            type="password"
            class="form-control form-control-lg"
            v-model="password"
            required
            autocomplete="current-password"
          />
        </div>

        <button type="submit" class="btn btn-primary btn-lg w-100" :disabled="loading">
          <font-awesome-icon icon="sign-in-alt" class="me-2" />
          {{ loading ? '登入中...' : '登入系統' }}
        </button>
      </form>
      
      <div class="text-center mt-4">
        <router-link to="/" class="text-secondary small">返回首頁</router-link>
      </div>
    </div>
  </div>
</template>

<script>
import { mapActions } from 'vuex'

export default {
  name: 'SuperadminLogin',
  data() {
    return {
      username: '',
      password: '',
      loading: false,
      error: ''
    }
  },
  methods: {
    ...mapActions(['fetchConfig']),
    onLogin() {
      const { username, password } = this
      if (!username || !password) return

      this.loading = true
      this.error = ''

      fetch('./api/config/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password })
      })
        .then(async (res) => {
          if (!res.ok) {
            if (res.status === 401) throw new Error('用戶名或密碼錯誤')
            throw new Error('登入失敗，請稍後再試')
          }
          return res.text()
        })
        .then((token) => {
          this.$store.commit('updateJWT', token)
          // Attempt to load config so standard pages work afterward
          return this.fetchConfig().catch(() => {})
        })
        .then(() => {
          this.$router.push('/config/dashboard')
        })
        .catch((err) => {
          this.error = err.message
        })
        .finally(() => {
          this.loading = false
        })
    }
  }
}
</script>
