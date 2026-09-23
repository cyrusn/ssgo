<template>
  <img v-if="src" :src="src" class="img-thumbnail zoom" width="50" />
</template>

<script>
export default {
  props: ['src']
}
</script>
