<template>
  <div class="row">
    <div class="col-md-4" v-for="group in availableGroups" v-bind:key="group.id">
      <label  class="ckbox">
       <input type="checkbox" :value="group.id" v-model="checkedGroups" @change="emitValue"><span>{{group.name}}</span>
     </label>
   </div>
 </div>
</template>

<script>
import APIService from '@/services/APIService'

export default {
  name: 'group-input',
  props: ['value', 'field', 'groupSource'],
  data: () => ({
    availableGroups: [],
    checkedGroups: []
  }),
  watch: {
    value: function (val) {
      this.checkedGroups = val
    }
  },
  mounted () {
    this.loadGroups()
  },
  methods: {
    emitValue: function () {
      this.$emit('input', this.checkedGroups)
    },
    loadGroups: function () {
      let self = this
      APIService.get(self.groupSource, function (response) {
        self.availableGroups = response.data.visibleResults
      })
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
