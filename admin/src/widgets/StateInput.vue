<template>
  <div>
    <select type="text" v-if="isMutable"
      class="form-control"
      :class="{'validated-control': validated}"
      ref="input"
      :id="elementId"
      v-bind:value="value"
      :name="elementId"
      @input="emitValue"
      :required="isRequired">
      <option v-for="state in states" :value="state" v-bind:key="state">{{state}}</option>
    </select>
    <span v-if="!isMutable">{{value}}</span>
  </div>
</template>

<script>
export default {
  name: 'state-input',
  props: ['value', 'required', 'field', 'validated', 'id'],
  data: () => ({
    states: [
      'AL',
      'AK',
      'AZ',
      'AR',
      'CA',
      'CO',
      'CT',
      'DE',
      'FL',
      'GA',
      'HI',
      'ID',
      'IL',
      'IN',
      'IA',
      'KS',
      'KY',
      'LA',
      'ME',
      'MD',
      'MA',
      'MI',
      'MN',
      'MS',
      'MO',
      'MT',
      'NE',
      'NV',
      'NH',
      'NJ',
      'NM',
      'NY',
      'NC',
      'ND',
      'OH',
      'OK',
      'OR',
      'PA',
      'PR',
      'RI',
      'SC',
      'SD',
      'TN',
      'TX',
      'UT',
      'VT',
      'VA',
      'VI',
      'WA',
      'WV',
      'WI',
      'WY'
    ]
  }),
  computed: {
    elementId: function () {
      if (this.field) {
        return this.field.id
      } else {
        return this.id
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
