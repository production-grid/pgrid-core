<template>
  <div>
    <textarea v-if="field.mutable"
      class="form-control"
      ref="input"
      :id="field.id"
      v-bind:value="value"
      :name="field.id"
      @input="emitValue"
      :required="field.required">
    </textarea>
    <span v-if="!field.mutable">{{value}}</span>
  </div>
</template>

<script>
export default {
  name: 'textarea-input',
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
