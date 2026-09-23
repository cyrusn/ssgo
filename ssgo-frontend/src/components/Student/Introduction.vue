<template>
  <div class="alert alert-primary introduction-content">
    <h5 class="alert-heading fw-bold">注意</h5>
    <hr />
    <div v-html="compiledMarkdown"></div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { marked } from 'marked'

export default {
  computed: {
    ...mapGetters(['introductionMarkdown']),
    compiledMarkdown() {
      return marked.parse(this.introductionMarkdown || '')
    }
  }
}
</script>

<style>
.introduction-content h1, 
.introduction-content h2, 
.introduction-content h3, 
.introduction-content h4, 
.introduction-content h5 {
  margin-bottom: 1rem;
  font-weight: bold;
}
.introduction-content hr {
  margin: 1rem 0;
}
.introduction-content p {
  margin-bottom: 1rem;
}
.introduction-content ul, 
.introduction-content ol {
  margin-left: 1.5rem;
  margin-bottom: 1rem;
}
</style>
