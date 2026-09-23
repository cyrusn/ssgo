<template>
  <input type="number" class="form-control" v-model.number="localValue" @input="onInput" />
</template>

<script>
import { mapState, mapActions } from 'vuex'
import _ from 'lodash'

export default {
  props: ['code'],
  data () {
    return {
      localValue: 0
    }
  },
  created () {
    this.localValue = this.capacities[this.code] || 0
    this.debouncedUpdate = _.debounce(this.sendUpdate, 500)
  },
  watch: {
    capacities: {
      handler (newCaps) {
        this.localValue = newCaps[this.code] || 0
      },
      deep: true
    }
  },
  computed: {
    ...mapState('subject', ['capacities'])
  },
  methods: {
    ...mapActions('subject', ['updateCapacity']),
    onInput () {
      this.debouncedUpdate()
    },
    sendUpdate () {
      const { code, localValue } = this
      this.updateCapacity({
        code,
        capacity: localValue
      })
    }
  }
}
</script>
