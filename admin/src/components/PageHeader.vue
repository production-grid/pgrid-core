<template>
  <div class="slim-pageheader">
    <ol class="breadcrumb slim-breadcrumb">
      <li class="breadcrumb-item" v-if="breadcrumbCategory"><a href="#">{{breadcrumbCategory}}</a></li>
      <li class="breadcrumb-item active" aria-current="page">{{pageTitle}}</li>
    </ol>
    <h6 class="slim-pagetitle">{{pageTitle}}</h6>
  </div><!-- slim-pageheader -->
</template>

<script>
export default {
  name: 'page-header',
  props: ['breadcrumbCategory', 'pageTitle'],
  data: () => ({
  }),
  mounted () {
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
