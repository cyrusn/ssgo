<template>
  <div class="row mb-4">
    <div class="col-6 row">
      <div v-for="g in groups" :key="g.id" class="col-6">
        <h4 class="text-center">選修{{ g }}組</h4>
        <div v-for="s in subjectsGrouped[g]" :key="s.code" class="mb-2">
          <div class="input-group">
            <div class="input-group-prepend">
              <span class="input-group-text">{{ s.slug }} </span>
            </div>
            <input-capacity :code="s.code" />
          </div>
        </div>
      </div>
    </div>
    <div class="col-6">
      <allocation-result />
    </div>
  </div>
</template>

<script>
import InputCapacity from '@/components/Admin/InputCapacity'
import AllocationResult from '@/components/Admin/AllocationResult'
import _ from 'lodash'

import { mapState, mapActions, mapGetters } from 'vuex'

export default {
  mounted() {
    this.list()
    this.getStudents()
  },
  data() {
    return {
      groups: ['1', '2']
    }
  },
  components: {
    InputCapacity,
    AllocationResult
  },
  computed: {
    ...mapState('subject', ['capacities']),
    ...mapGetters(['subjects']),
    subjectsGrouped() {
      return _.groupBy(this.subjects, 'group')
    }
  },
  methods: {
    ...mapActions('subject', ['list']),
    ...mapActions('students', { getStudents: 'get' })
  }
}
</script>
