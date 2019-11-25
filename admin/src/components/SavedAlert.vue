<template>
  <div class="alert alert-outline alert-success" role="alert" v-if="visible">
   <button type="button" class="close" data-dismiss="alert" aria-label="Close">
     <span aria-hidden="true">&times;</span>
   </button>
   <strong>Done.</strong> Your changes have been saved.
 </div><!-- alert -->
</template>

<script>
export default {
  name: 'saved-alert',
  props: ['globalState', 'visible']
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
