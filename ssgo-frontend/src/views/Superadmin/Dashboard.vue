<template>
  <div v-if="role === 'SUPERADMIN'" class="container py-5">
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h1 class="fw-bold m-0 text-dark">Configuration Dashboard</h1>
      <div class="d-flex align-items-center">
        <span class="me-3 fw-bold text-secondary">
          <font-awesome-icon icon="user-shield" class="me-1" />
          {{ userAlias }} (系統管理員)
        </span>
        <button class="btn btn-outline-danger" @click="onLogout">安全登出</button>
      </div>
    </div>
    <hr class="mb-5" />

    <div v-if="loading" class="text-center py-5">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
      <p class="mt-2 text-secondary">載入中...</p>
    </div>

    <div v-else>
      <CohortManager 
        :cohorts="cohorts" 
        :jwt="jwt" 
        :apiFetch="apiFetch" 
        @refresh="loadCohorts" 
      />

      <ConfigEditor 
        :cohorts="cohorts" 
        :activeCohortID="activeCohortID" 
        :jwt="jwt" 
        :apiFetch="apiFetch" 
        @refresh="loadCohorts" 
      />

      <DataImporter 
        :activeCohortID="activeCohortID" 
        :apiFetch="apiFetch" 
      />

      <UserModifier 
        :activeCohortID="activeCohortID" 
        :apiFetch="apiFetch" 
      />

      <AdminAccountManager 
        :userAlias="userAlias" 
        :apiFetch="apiFetch" 
      />
    </div>
  </div>
  <div v-else class="container py-5 text-center">
    <div class="alert alert-danger d-inline-block p-5 rounded shadow">
      <h3>權限不足</h3>
      <router-link to="/config" class="btn btn-primary mt-3">去登入</router-link>
    </div>
  </div>
</template>

<script>
import { mapState, mapGetters } from 'vuex'
import CohortManager from '@/components/Superadmin/CohortManager.vue'
import ConfigEditor from '@/components/Superadmin/ConfigEditor.vue'
import DataImporter from '@/components/Superadmin/DataImporter.vue'
import UserModifier from '@/components/Superadmin/UserModifier.vue'
import AdminAccountManager from '@/components/Superadmin/AdminAccountManager.vue'

export default {
  name: 'SuperadminDashboard',
  components: {
    CohortManager,
    ConfigEditor,
    DataImporter,
    UserModifier,
    AdminAccountManager
  },
  data () {
    return {
      cohorts: [],
      loading: true
    }
  },
  computed: {
    ...mapState(['jwt']),
    ...mapGetters(['role', 'userAlias']),
    activeCohortID () {
      const active = this.cohorts.find(c => c.isActive || c.IsActive)
      return active ? (active.id || active.ID) : null
    }
  },
  mounted () {
    if (this.role === 'SUPERADMIN') {
      this.loadCohorts()
    } else {
      this.$router.push('/config')
    }
  },
  methods: {
    async apiFetch (url, options = {}) {
      if (!options.headers) options.headers = {}
      options.headers['jwt'] = this.jwt
      
      try {
        const res = await fetch(url, options)
        if (res.status === 401) {
          console.warn('Session expired or unauthorized. Redirecting to login.')
          this.$store.commit('updateJWT', '')
          this.$router.push('/config')
          return res
        }
        return res
      } catch (e) {
        console.error('Fetch error:', e)
        throw e
      }
    },
    async loadCohorts () {
      this.loading = true
      try {
        const res = await this.apiFetch('./api/config/cohorts')
        if (res.ok) {
          const d = await res.json()
          this.cohorts = Array.isArray(d) ? d : []
        }
      } catch (e) {
        console.error('Cohorts load failed', e)
      } finally {
        this.loading = false
      }
    },
    onLogout () {
      this.$store.commit('updateJWT', '')
      window.location.reload(true)
    }
  }
}
</script>

<style scoped>
.container { max-width: 1000px; }
</style>
