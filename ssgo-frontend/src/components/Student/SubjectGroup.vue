<template>
  <li class="list-group-item flex-fill">
    <span
      v-if="name === 'available'"
      class="me-2 badge bg-primary text-light rounded-pill"
    >
      <font-awesome-icon icon="arrows-alt" />
    </span>
    <span
      v-else
      class="badge text-light rounded-pill rounded-pill"
      :class="isConfirmed ? 'bg-success ' : 'bg-danger'"
      >{{ index + 1 }}</span
    >
    <span
      v-for="(s, n) in element.subjects"
      class="ms-2 p-1 rounded font-weight-bolder"
      :class="
        highlightedSubjects[s]
          ? n == 0
            ? 'bg-warning'
            : 'bg-info'
          : 'text-dark '
      "
      :key="s"
      @click="toggleHighlightedSubject(s)"
      :title="getSubject(s).cname"
    >
      {{ getSubject(s).slug }}
    </span>
  </li>
</template>

<script>
import { mapState, mapMutations, mapGetters } from 'vuex'
import _ from 'lodash'

export default {
  props: ['name', 'element', 'index', 'movable', 'isConfirmed'],

  computed: {
    ...mapState('subject', ['highlightedSubjects']),
    ...mapGetters(['subjects'])
  },
  methods: {
    ...mapMutations('subject', ['toggleHighlightedSubject']),

    getSubject (code) {
      return _.find(this.subjects, { code }) || { slug: code, cname: code }
    }
  }
}
</script>

<style>
.dragArea {
  min-height: 50px;
}

/* Overrides list-group-item from Bootstrap */
.list-group-item {
  /* height: 30px; */
  padding: 3px 15px !important;
}
</style>
