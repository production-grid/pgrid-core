<template>
    <select type="text" v-if="isMutable"
      class="form-control"
      ref="input"
      :id="elementId"
      v-bind:value="value"
      :name="elementId"
      @input="emitValue"
      :required="isRequired">
      <option v-for="zone in timeZones" :value="zone" v-bind:key="zone">{{zone}}</option>
    </select>
</template>

<script>
export default {
  name: 'time-zone-input',
  props: ['value', 'field'],
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
    timeZones: [
      'America/Adak',
      'America/Anchorage',
      'America/Boise',
      'America/Chicago',
      'America/Denver',
      'America/Detroit',
      'America/Fort_Wayne',
      'America/Indiana/Indianapolis',
      'America/Indiana/Knox',
      'America/Indiana/Marengo',
      'America/Indiana/Petersburg',
      'America/Indiana/Tell_City',
      'America/Indiana/Vevay',
      'America/Indiana/Vincennes',
      'America/Indiana/Winamac',
      'America/Indianapolis',
      'America/Juneau',
      'America/Kentucky/Louisville',
      'America/Kentucky/Monticello',
      'America/Knox_IN',
      'America/Los_Angeles',
      'America/Louisville',
      'America/Menominee',
      'America/Metlakatla',
      'America/New_York',
      'America/Nome',
      'America/North_Dakota/Beulah',
      'America/North_Dakota/Center',
      'America/North_Dakota/New_Salem',
      'America/Phoenix',
      'America/Sitka',
      'America/Yakutat'
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
