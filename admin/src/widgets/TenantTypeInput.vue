<template>
    <select type="text" v-if="isMutable"
      class="form-control"
      ref="input"
      :id="elementId"
      v-bind:value="value"
      :name="elementId"
      @input="emitValue"
      :required="isRequired">
      <option v-for="tenantType in session.tenantTypes" :value="tenantType.singular" v-bind:key="tenantType.singular">{{tenantType.singular}}</option>
    </select>
</template>

<script>
export default {
  name: 'tenant-type-input',
  props: ['value', 'field', 'session'],
  computed: {
    elementId: function () {
      if (this.field) {
        return this.field.id
      } else {
        return null
      }
    },
    isMutable: function () {
      if (this.field) {
        return this.field.mutable
      } else {
        return true
      }
    },
    isRequired: function () {
      if (this.field) {
        return this.field.required
      } else {
        return false
      }
    }

  },
  data: () => ({
  }),
  mounted () {
  },
  methods: {
    emitValue: function (e) {
      this.$emit('input', e.target.value)
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
