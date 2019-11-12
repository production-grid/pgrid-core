<template>
  <div>
    <select type="text" v-if="field.mutable"
      class="form-control"
      ref="input"
      :id="field.id"
      v-bind:value="value"
      :name="field.id"
      @input="emitValue"
      :required="field.required">
      <option v-for="hour in hours" :value="hour" v-bind:key="hour">{{hour}}</option>
    </select>
    <span v-if="!field.mutable">{{value}}</span>
  </div>
</template>

<script>
export default {
  name: 'hour-input',
  props: ['value', 'field'],
  data: () => ({
    hours: [
      '12:00 AM',
      '1:00 AM',
      '2:00 AM',
      '3:00 AM',
      '4:00 AM',
      '5:00 AM',
      '6:00 AM',
      '7:00 AM',
      '8:00 AM',
      '9:00 AM',
      '10:00 AM',
      '11:00 AM',
      '12:00 PM',
      '1:00 PM',
      '2:00 PM',
      '3:00 PM',
      '4:00 PM',
      '5:00 PM',
      '6:00 PM',
      '7:00 PM',
      '8:00 PM',
      '9:00 PM',
      '10:00 PM',
      '11:00 PM',
    ]
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
