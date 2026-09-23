<template>
  <div class="alert alert-warning instruction-content">
    <h5 class="alert-heading fw-bold">使用須知</h5>
    <hr />
    <div v-html="compiledMarkdown"></div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { marked } from 'marked'

export default {
  computed: {
    ...mapGetters(['instructionMarkdown']),
    compiledMarkdown() {
      return marked.parse(this.instructionMarkdown || '')
    }
  }
}
</script>

<style>
.instruction-content h1, 
.instruction-content h2, 
.instruction-content h3, 
.instruction-content h4, 
.instruction-content h5 {
  margin-bottom: 1rem;
  font-weight: bold;
}
.instruction-content hr {
  margin: 1rem 0;
}
.instruction-content p {
  margin-bottom: 1rem;
}
.instruction-content ul, 
.instruction-content ol {
  margin-left: 1.5rem;
  margin-bottom: 1rem;
}
</style>
