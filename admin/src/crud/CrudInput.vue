<template>
  <div :class="effectiveColumnStyle">
     <div class="form-group" v-if="!field.hideIfEmpty">
       <!-- Checkboxes go before captions -->
       <checkbox-input v-if="field.dataType == 'checkbox'" :value="value" :field="field" @input="emitValue"/>&nbsp;
       <label v-if="field.dataType == 'checkbox'" class="form-control-label font-weight-bold" v-bind:for="field.id">{{field.caption}}</label>

       <label v-else class="form-control-label font-weight-bold">{{field.caption}}: <span v-if="field.required" class="tx-danger">*</span></label>

       <!-- It's a giant pile of input widgets -->
       <address-widget v-if="field.dataType == 'address'" :value="value" :field="field" @input="emitValue" :mode="mode"/>
       <admin-roles-input v-if="field.dataType == 'adminRoles'" :value="value" :field="field" @input="emitValue"/>
       <boolean-input v-if="field.dataType == 'boolean'" :value="value" :field="field" @input="emitValue"/>
       <hour-input v-if="field.dataType == 'hour'" :value="value" :field="field" @input="emitValue" :mode="mode"/>
       <phone-input v-if="field.dataType == 'phone'" :value="value" :field="field" @input="emitValue" :mode="mode"/>
       <string-input v-if="field.dataType == 'string'" :value="value" :field="field" @input="emitValue"/>
       <textarea-input v-if="field.dataType == 'textarea'" :value="value" :field="field" @input="emitValue"/>
       <time-zone-input v-if="field.dataType == 'timeZone'" :value="value" :field="field" @input="emitValue" :mode="mode"/>
       <currency-input v-if="field.dataType == 'currency'" :value="value" :field="field" @input="emitValue" :mode="mode"/>
       <subdomain-input v-if="field.dataType == 'subdomain'" :value="value" :field="field" @input="emitValue" :mode="mode" :session="session"/>
       <tenant-type-input v-if="field.dataType == 'tenant-type'" :value="value" :field="field" @input="emitValue" :mode="mode" :session="session"/>
       <!-- End input widgets -->

       <div class="help-block" v-if="field.help">{{field.help}}</div>
     </div>
  </div>
</template>

<script>
import AddressWidget from '../widgets/AddressWidget'
import AdminRolesInput from '../widgets/AdminRolesInput'
import BooleanInput from '../widgets/BooleanInput'
import CheckboxInput from '../widgets/CheckboxInput'
import HourInput from '../widgets/HourInput'
import PhoneInput from '../widgets/PhoneInput'
import StringInput from '../widgets/StringInput'
import TextAreaInput from '../widgets/TextAreaInput'
import TimeZoneInput from '../widgets/TimeZoneInput'
import CurrencyInput from '../widgets/CurrencyInput'
import SubdomainInput from '../widgets/SubdomainInput'
import TenantTypeInput from '../widgets/TenantTypeInput'

export default {
  name: 'crud-input',
  props: ['value', 'field', 'form', 'mode', 'columnStyle', 'additionalStyles', 'session'],
  components: {
    'address-widget': AddressWidget,
    'admin-roles-input': AdminRolesInput,
    'boolean-input': BooleanInput,
    'checkbox-input': CheckboxInput,
    'hour-input': HourInput,
    'phone-input': PhoneInput,
    'string-input': StringInput,
    'textarea-input': TextAreaInput,
    'time-zone-input': TimeZoneInput,
    'currency-input': CurrencyInput,
    'subdomain-input': SubdomainInput,
    'tenant-type-input': TenantTypeInput,
  },
  data: () => ({
  }),
  mounted () {
  },
  computed: {
    effectiveColumnStyle: function () {
      if (this.field.columnStyle) {
        return this.field.columnStyle
      } else {
        return this.columnStyle
      }
    }
  },
  methods: {
    emitValue: function (value) {
      this.$emit('input', value)
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
