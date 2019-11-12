<template>
  <input type="checkbox"
    class="ml-2"
    ref="input"
    :id="field.id"
    v-bind:checked="value"
    :name="field.id"
    @click="emitValue"
    :required="field.required"/>
</template>

<script>
export default {
  name: 'boolean-input',
  props: ['value', 'field'],
  data: () => ({
  }),
  mounted () {
  },
  methods: {
    emitValue: function (e) {
      this.$emit('input', e.target.checked)
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
