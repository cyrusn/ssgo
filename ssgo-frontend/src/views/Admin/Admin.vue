<template>
  <SystemTitle />
  <div v-if="role === 'ADMIN'">
    <router-view />
  </div>
  <div v-else>
    <forbidden />
  </div>
</template>

<script>
import SystemTitle from '@/components/Shared/SystemTitle'
import Forbidden from '@/components/Shared/Forbidden.vue'
import { mapGetters } from 'vuex'

export default {
  components: {
    SystemTitle,
    Forbidden
  },
  computed: {
    ...mapGetters(['role'])
  }
}
</script>
