<template>
  <div class="mt-4">
    <div class="progress" v-if="processing">
      <div class="progress-bar progress-bar-striped progress-bar-animated" role="progressbar" style="width: 100%" aria-valuenow="100" aria-valuemin="0" aria-valuemax="100"></div>
    </div>
    <div v-if="!processing">
      <button type="button" class="btn btn-warning" data-test-id="canceButton" @click="cancelForm()" v-if="listRoute">CANCEL</button>
      <button type="button" class="btn btn-danger mg-l-5" v-if="isDeletable" @click="sendDelete()">DELETE</button>
      <button type="submit" class="btn btn-primary mg-l-5 float-right" data-test-id="submitButton">{{effectiveSaveCaption}}</button>
    </div>
  </div>
</template>

<script>
import router from '../router'

export default {
  name: 'form-button-panel',
  props: ['globalState', 'processing', 'listRoute', 'isDeletable', 'saveCaption'],
  computed: {
    effectiveSaveCaption: function () {
      if (this.saveCaption) {
        return this.saveCaption
      }
      return 'Save'
    }
  },
  methods: {
    cancelForm: function () {
      router.push(this.listRoute)
    },
    sendDelete: function () {
      this.$emit('delete')
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
