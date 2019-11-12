<template>
  <select class="form-control"
    ref="input"
    :id="elementId"
    :name="elementId"
    @change="emitValue"
    v-bind:value="internalValue"
    :required="isRequired">
    <option value="0">No</option>
    <option value="1">Yes</option>
  </select>
</template>

<script>
export default {
  name: 'boolean-input',
  props: ['value', 'field'],
  data: () => ({
    internalValue: '0'
  }),
  mounted () {
  },
  watch: {
    'value': function (val) {
      if (val) {
        this.internalValue = 1
      } else {
        this.internalValue = 0
      }
    }
  },
  computed: {
    elementId: function () {
      if (this.field) {
        return this.field.id
      } else {
        return null
      }
    },
    isRequired: function () {
      if (this.field) {
        return this.field.required
      } else {
        return this.required
      }
    },
    isMutable: function () {
      if (this.field) {
        return this.field.mutable
      } else {
        return true
      }
    }
  },
  methods: {
    emitValue: function (e) {
      this.$emit('input', e.target.value === '1')
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
