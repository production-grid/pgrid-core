<template>
  <div>
    <input type="text"
      class="form-control"
      ref="input"
      v-bind:value="value"
      @input="emitValue"
      :maxlength="9"
      pattern="[0-9]{9}"
      :required="required"/>
    <small class="help-text">Nine digit ACH routing number for your financial institution.</small>
  </div>
</template>

<script>
export default {
  name: 'routing-number-input',
  props: ['value', 'required'],
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
