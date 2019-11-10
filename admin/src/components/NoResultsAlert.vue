<template>
  <div class="alert alert-outline alert-info m-3 mt-4 text-center" v-if="!hasResults">
    <h5 v-if="loading" class="mt-2">Loading...</h5>
    <h5 v-if="!loading" class="mt-2">{{effectiveCaption}}</h5>
  </div>
</template>

<script>
export default {
  name: 'no-results-alert',
  props: ['globalState', 'results', 'loading', 'caption'],
  computed: {
    effectiveCaption: function () {
      if (this.caption) {
        return this.caption
      }
      return "This list contains no data."
    },
    hasResults: function () {
      return this.results && this.results.length > 0
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
