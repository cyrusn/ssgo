<template>
  <li class="nav-item">
    <router-link
      :to="{ name: name }"
      tag="li"
      class="nav-link"
      :class="{ active: route.name === name }"
    >
      <slot></slot>
    </router-link>
  </li>
</template>

<script>
import { mapState } from 'vuex'
export default {
  props: ['name'],
  computed: {
    ...mapState(['route'])
  }
}
</script>
