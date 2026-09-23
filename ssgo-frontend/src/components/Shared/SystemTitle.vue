<template>
  <div>
    <h1 class="text-center">
      {{ systemTitle }}
    </h1>
    <hr />
  </div>
</template>

<script>
import { mapGetters } from 'vuex'

export default {
  computed: {
    ...mapGetters(['systemTitle'])
  }
}
</script>
