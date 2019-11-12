<template>
  <div class="form-layout mg-t-10">
    <div class="row">
      <crud-input v-for="field in fields"
        :field="field"
        v-bind:key="field.id"
        v-model="value[field.id]"
        :form=value
        :mode="mode"
        :globalState="globalState"
        :ref="field.id"
        :columnStyle="columnStyle"
        :session="session"
      />
    </div>
    <fieldset v-for="fieldset in fieldsets" v-bind:key="fieldset.id" class="form-control">
      <legend class="form-control-legend">{{fieldset.caption}}</legend>
      <div class="row">
        <crud-input v-for="field in fieldset.fields"
          :field="field"
          v-bind:key="field.id"
          v-model="value[field.id]"
          :form=value
          :mode="mode"
          :globalState="globalState"
          :ref="field.id"
          :columnStyle="columnStyle"
        />
      </div>
    </fieldset>
  </div>
</template>

<script>
import CrudInput from './CrudInput'

export default {
  name: 'crud-panel',
  props: ['globalState', 'value', 'fields', 'fieldsets', 'mode', 'columnStyle', 'session'],
  components: {
    'crud-input': CrudInput
  },
  data: () => ({
  })
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.form-control-legend {
  width: revert;
  padding: revert;
  margin-bottom: revert;
  font-size: 0.875rem;
  font-weight: 700;
}
</style>
