<template>
  <div>
    <div class="input-group" v-if="field.mutable">
      <input type="text"
        class="form-control"
        ref="input"
        :id="field.id"
        v-bind:value="value"
        :name="field.id"
        @input="emitValue"
        :minlength="field.min"
        :maxlength="field.max"
        :required="field.required"/>
        <div class="input-group-append">
          <span class="input-group-text">.{{session.rootTenantHost}}</span>
        </div>
    </div>
    <span v-if="!field.mutable">{{value}}</span>
  </div>
</template>

<script>
export default {
  name: 'subdomain-input',
  props: ['value', 'field', 'session'],
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
