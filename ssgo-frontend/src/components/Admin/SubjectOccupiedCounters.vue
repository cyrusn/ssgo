<template>
  <div class="row">
    <div v-for="g in groups" :key="g.id" class="col-6">
      <h4>選修{{ g }}組</h4>
      <ul class="list-group">
        <li v-for="s in subjectsGrouped[g]" :key="s.code" class="list-group-item py-1">
          <span>{{ s.slug }}</span>
          <span class="ms-2 badge bg-secondary rounded-pill"
            >{{ counters[s.code] }} / {{ capacities[s.code] }}</span
          >
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
import { mapState, mapGetters } from 'vuex'
import _ from 'lodash'

export default {
  data () {
    return {
      groups: [1, 2]
    }
  },
  computed: {
    ...mapState('subject', ['capacities']),
    ...mapGetters(['subjects']),
    subjectsGrouped () {
      return _.groupBy(this.subjects, 'group')
    }
  },
  props: ['counters']
}
</script>
