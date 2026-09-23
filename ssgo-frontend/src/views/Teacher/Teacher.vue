<template>
  <SystemTitle />
  <div v-if="includes(role)">
    <router-view />
  </div>
  <div v-else>
    <forbidden />
  </div>
</template>

<script>
import _ from 'lodash'
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
  },
  methods: {
    includes(role) {
      return _.includes(['TEACHER', 'ADMIN'], role)
    }
  }
}
</script>
