<template>
  <div>
    <h1>Access Forbidden!</h1>
    <h2>Error 403</h2>
  </div>
</template>

<script>
export default {
  created () {
    this.$router.push('/')
  }
}
</script>
