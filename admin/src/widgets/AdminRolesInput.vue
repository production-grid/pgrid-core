<template>
  <div>
    <label v-for="result in availableRoles" v-bind:key="result" class="ckbox">
     <input type="checkbox" :value="result" v-model="checkedRoles" @change="emitValue"><span>{{result}}</span>
   </label>
 </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'admin-roles-input',
  props: ['value', 'field'],
  data: () => ({
    availableRoles: [],
    checkedRoles: null
  }),
  watch: {
    value: function (val) {
      this.checkedRoles = val
    }
  },
  mounted () {
    this.loadRoles()
  },
  methods: {
    emitValue: function () {
      this.$emit('input', this.checkedRoles)
    },
    loadRoles: function () {
      let self = this
      axios.get('/api/admin/roles')
        .then(function (response) {
          self.availableRoles = response.data.roles
        })
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
