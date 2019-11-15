<template>
  <div class="row">
    <div class="col-md-4" v-for="perm in availablePerms" v-bind:key="perm.key">
      <label  class="ckbox">
       <input type="checkbox" :value="perm.key" v-model="checkedPerms" @change="emitValue"><span>{{perm.description}}</span>
     </label>
   </div>
 </div>
</template>

<script>
import APIService from '@/services/APIService'

export default {
  name: 'permissions-input',
  props: ['value', 'field', 'permSource'],
  data: () => ({
    availablePerms: [],
    checkedPerms: []
  }),
  watch: {
    value: function (val) {
      this.checkedPerms = val
    }
  },
  mounted () {
    this.loadPerms()
  },
  methods: {
    emitValue: function () {
      this.$emit('input', this.checkedPerms)
    },
    loadPerms: function () {
      let self = this
      APIService.get(self.permSource, function (response) {
        self.availablePerms = response.data.visibleResults
      })
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
