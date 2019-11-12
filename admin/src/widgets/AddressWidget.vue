<template>
    <div class="addressWidget pd-5" :class="{'validated-control': validated}" :id="id" :required="required">
      <div class="row">
        <div class="col-md-12">
        <input type="text" class="form-control form-control-sm" name="line1" v-model="address.line1" placeholder="Address Line 1" :readonly="!mutable">
        </div>
      </div>
      <div class="row">
        <div class="col-md-12">
        <input type="text" class="form-control form-control-sm" name="line2" v-model="address.line2" placeholder="Address Line 2"  :readonly="!mutable">
        </div>
      </div>
      <div class="row">
        <div class="col-md-5">
          <input type="text" class="form-control form-control-sm" name="city" id="city" v-model="address.city" placeholder="Hateno Village" :readonly="!mutable">
        </div>
        <div class="col-md-3">
          <select class="form-control form-control-sm" name="stateOrProvince" id="stateOrProvince" v-model="address.stateOrProvince" :readonly="!mutable">
            <option>State...</option>
            <option v-for="state in states" :value="state" v-bind:key="state">{{state}}</option>
          </select>
        </div>
        <div class="col-md-4">
          <input type="text" class="form-control form-control-sm" name="postalCode" maxlength="5" v-model="address.postalCode" placeholder="00000" :readonly="!mutable">
        </div>
      </div>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'address-widget',
  props: ['value', 'addressType', 'mutable', 'validated', 'id', 'required'],
  data: () => ({
    geocoding: false,
    badAddress: false,
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
    address () {
      if (!this.value) {
        return {}
      } else {
        return this.value
      }
    }
  },
  methods: {
    geocode: function () {
      const vm = this
      vm.geocoding = true
      axios.get('/api/locations/geocode?addr=' + vm.address.addressString)
      .then(function (response) {
        vm.geocoding = false
        vm.value = response.data
        vm.value.addressType = vm.addressType
        vm.$emit('input', response.data)
      })
    }
  }

}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

.row {
  margin-top: 2px;
}

.is-invalid {
  border: solid 1px red;
}

</style>
