<template>
  <figure class="figure">
    <img :src="signature" class="figure-img img-thumbnail" />
    <figcaption class="figure-caption text-center">
      家長簽署
    </figcaption>
  </figure>
</template>

<script>
export default {
  props: ['signature']
}
</script>
