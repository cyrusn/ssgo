<template>
  <div class="alert alert-danger not-accept-content" v-html="compiledMarkdown"></div>
</template>

<script>
import { mapGetters } from 'vuex'
import { marked } from 'marked'

export default {
  computed: {
    ...mapGetters(['notAcceptMarkdown']),
    compiledMarkdown() {
      return marked.parse(this.notAcceptMarkdown || '')
    }
  }
}
</script>

<style>
.not-accept-content h1, 
.not-accept-content h2, 
.not-accept-content h3, 
.not-accept-content h4, 
.not-accept-content h5 {
  margin-bottom: 1rem;
  font-weight: bold;
}
</style>
