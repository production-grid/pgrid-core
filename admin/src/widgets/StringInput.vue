<template>
  <div>
    <input type="text" v-if="field.mutable"
      class="form-control"
      ref="input"
      :id="field.id"
      v-bind:value="value"
      :name="field.id"
      @input="emitValue"
      :minlength="field.min"
      :maxlength="field.max"
      :required="field.required"/>
    <span v-if="!field.mutable">{{value}}</span>
  </div>
</template>

<script>
export default {
  name: 'string-input',
  props: ['value', 'field'],
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
