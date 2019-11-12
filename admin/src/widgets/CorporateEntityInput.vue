<template>
  <div>
    <select type="text" v-if="isMutable"
      class="form-control"
      ref="input"
      :id="elementId"
      v-bind:value="value"
      :name="elementId"
      @input="emitValue"
      :required="isRequired">
      <option v-for="type in options" :value="type.code" v-bind:key="type.code">{{type.description}}</option>
    </select>
    <span v-if="!isMutable">{{value}}</span>
  </div>
</template>

<script>
export default {
  name: 'corporate-entity-input',
  props: ['value', 'required', 'field'],
  data: () => ({
    options: [
      {
        code: "SOLE_PROPRIETORSHIP",
        description: "Sole Proprietorship"
      },
      {
        code: "CORPORATION",
        description: "Corporation"
      },
      {
        code: "LLC",
        description: "LLC (Limited Liability Corporation)"
      },
      {
        code: "PARTNERSHIP",
        description: "Partnership"
      },
      {
        code: "NON_PROFIT",
        description: "Non Profit"
      }
    ]
  }),
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
