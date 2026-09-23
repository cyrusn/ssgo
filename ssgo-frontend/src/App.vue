<template>
  <div id="app">
    <!-- Offline view if no active cohort configuration is configured in Master DB -->
    <div v-if="!configLoaded && !isSuperadminRoute" class="container d-flex flex-column align-items-center justify-content-center min-vh-100 text-center">
      <div class="card shadow p-5 border-0 bg-light" style="max-width: 550px;">
        <div class="text-danger mb-4">
          <font-awesome-icon icon="exclamation-triangle" size="4x" />
        </div>
        <h2 class="fw-bold mb-3">系統尚未啟用</h2>
        <p class="text-secondary mb-4 fs-5">
          目前系統沒有設定任何「已啟用 (Active)」的選科年度。<br />
          如果您是學生或老師，請聯絡系統管理員處理。
        </p>
        <div class="border-top pt-4">
          <p class="text-secondary small">如果您是系統管理員 (super_admin)，請點擊下方按鈕登入設定選科年度：</p>
          <router-link to="/config" class="btn btn-primary px-4 py-2">
            <font-awesome-icon icon="sign-in-alt" class="me-2" />系統管理員登入
          </router-link>
        </div>
      </div>
    </div>

    <!-- Normal App View -->
    <div v-else>
      <router-view name="navbar" />
      <div class="container pt-4">
        <router-view />
      </div>
    </div>
  </div>
</template>

<script>
import { mapState } from 'vuex'

export default {
  computed: {
    ...mapState(['config']),
    configLoaded() {
      return this.config !== null && Object.keys(this.config).length > 0
    },
    isSuperadminRoute() {
      const path = this.$route && this.$route.path ? this.$route.path : ''
      const name = this.$route && this.$route.name ? this.$route.name : ''
      return path.startsWith('/config') || name.startsWith('config') || path.startsWith('/superadmin')
    }
  }
}
</script>
