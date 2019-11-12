<template>
  <div class="alert alert-info" role="alert" v-if="displayAlert">
    <button type="button" class="close" data-dismiss="alert" aria-label="Close">
      <span aria-hidden="true">&times;</span>
    </button>
    <strong>{{displayAlert}}</strong>
  </div><!-- alert -->
</template>

<script>
export default {
  name: 'next-alert',
  props: ['globalState', 'alert'],
  data: () => ({
    currentAlert: null
  }),
  computed: {
    displayAlert: function () {
      if (this.alert) {
        return this.alert
      } else {
        return this.currentAlert
      }
    }
  },
  mounted () {
    this.currentAlert = this.globalState.nextAlert
    this.globalState.nextAlert = null
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>


</style>
