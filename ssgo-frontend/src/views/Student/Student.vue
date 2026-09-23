<template>
  <SystemTitle />
  <div v-if="isAcceptResponses">
    <div v-if="role === 'STUDENT'">
      <router-view />
    </div>
    <div v-else>
      <forbidden />
    </div>
  </div>

  <div v-else>
    <not-accept-response />
  </div>
</template>

<script>
import SystemTitle from '@/components/Shared/SystemTitle'
import Forbidden from '@/components/Shared/Forbidden.vue'
import NotAcceptResponse from '@/components/Student/NotAcceptResponse.vue'
import { mapGetters } from 'vuex'

export default {
  components: {
    SystemTitle,
    Forbidden,
    NotAcceptResponse
  },
  computed: {
    ...mapGetters(['role', 'isAcceptResponses'])
  }
}
</script>
