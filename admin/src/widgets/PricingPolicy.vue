<template>
  <div class="form-layout mg-t-10">
    <div class="row">
      <div class="form-group col-md-6">
        <label class="form-control-label font-weight-bold">Description: <span class="tx-danger">*</span></label>
        <input class="form-control" type="input" v-model="value.description"/>
      </div>
      <div class="form-group col-md-6">
        <label class="form-control-label font-weight-bold">Enabled:</label>
        <boolean-input v-model="value.enabled"/>
      </div>
    </div>
    <div class="row">
      <div class="form-group col-md-6">
        <label class="form-control-label font-weight-bold">Type:</label>
        <select class="form-control" v-model="value.type">
          <option value="interchange">Interchange Plus</option>
          <option value="flatrate">Flat Rate</option>
        </select>
      </div>
      <div class="form-group col-md-6">
        <label class="form-control-label font-weight-bold">Priority: <span class="tx-danger">*</span></label>
        <input class="form-control" type="number" v-model="value.priority"/>
      </div>
    </div>
    <div class="row">
      <div class="form-group col-md-6">
        <label class="form-control-label font-weight-bold">Acquirer Code:</label>
        <select class="form-control" v-model="value.acquirerCode">
          <option value="SIM">SIM</option>
          <option value="EPX">EPX</option>
          <option value="TSYS">TSYS</option>
        </select>
      </div>
    </div>
    <label class="form-control-label font-weight-bold">Pricing List</label>
    <table class="table table-striped table-bordered display nowrap">
      <thead>
        <th scope="col">Rate</th>
        <th scope="col">
          {{effectiveCostCaption}}
          <div class="text-muted rate-help" v-if="costHelp">{{costHelp}}</div>
        </th>
        <th scope="col">
          {{effectiveCurrentCaption}}
          <div class="text-muted rate-help" v-if="currentHelp">{{currentHelp}}</div>
        </th>
        <th scope="col">
          {{effectiveLimitCaption}}
          <div class="text-muted rate-help" v-if="limitHelp">{{limitHelp}}</div>
        </th>
      </thead>
      <tbody>
        <tr v-if="isFlatRate">
          <th class="rate-caption">Flat Rate:</th>
          <td><bps-input v-model="value.flatRateCost" :lt="value.flatRateLimit"/></td>
          <td><bps-input v-model="value.flatRateCurrent" :lt="value.flatRateLimit" :gt="value.flatRateCost"/></td>
          <td><bps-input v-model="value.flatRateLimit" :gt="value.flatRateCost"/></td>
        </tr>
        <tr v-if="isFlatRate">
          <th class="rate-caption">Premium Flat Rate:</th>
          <td><bps-input v-model="value.premiumFlatRateCost" :lt="value.premiumFlatRateLimit"/></td>
          <td><bps-input v-model="value.premiumFlatRateCurrent" :lt="value.premiumFlatRateLimit" :gt="value.premiumFlatRateCost"/></td>
          <td><bps-input v-model="value.premiumFlatRateLimit" :gt="value.premiumFlatRateCost"/></td>
        </tr>
        <tr v-if="isFlatRate">
          <th class="rate-caption">Keyed Flat Rate:</th>
          <td><bps-input v-model="value.keyedFlatRateCost" :lt="value.keyedFlatRateLimit"/></td>
          <td><bps-input v-model="value.keyedFlatRateCurrent" :lt="value.keyedFlatRateLimit" :gt="value.keyedFlatRateCost"/></td>
          <td><bps-input v-model="value.keyedFlatRateLimit" :gt="value.keyedFlatRateCost"/></td>
        </tr>
        <tr v-if="isInterchangePlus">
          <th class="rate-caption">Interchange Markup:</th>
          <td><bps-input v-model="value.interchangeMarkupCost" :lt="value.interchangeMarkupLimit"/></td>
          <td><bps-input v-model="value.interchangeMarkupCurrent" :lt="value.interchangeMarkupLimit" :gt="value.interchangeMarkupCost"/></td>
          <td><bps-input v-model="value.interchangeMarkupLimit" :gt="value.interchangeMarkupCost"/></td>
        </tr>
        <tr v-if="isInterchangePlus">
          <th class="rate-caption">Per Transaction Fee:</th>
          <td><currency-input v-model="value.transactionFeeCost" :lt="value.transactionFeeLimit"/></td>
          <td><currency-input v-model="value.transactionFeeCurrent" :lt="value.transactionFeeLimit" :gt="value.transactionFeeCost"/></td>
          <td><currency-input v-model="value.transactionFeeLimit" :gt="value.transactionFeeCost"/></td>
        </tr>
        <tr v-if="isInterchangePlus">
          <th class="rate-caption">Premium Interchange Markup:</th>
          <td><bps-input v-model="value.premiumInterchangeMarkupCost" :lt="value.premiumInterchangeMarkupLimit"/></td>
          <td><bps-input v-model="value.premiumInterchangeMarkupCurrent" :lt="value.premiumInterchangeMarkupLimit" :gt="value.premiumInterchangeMarkupCost"/></td>
          <td><bps-input v-model="value.premiumInterchangeMarkupLimit" :gt="value.premiumInterchangeMarkupCost"/></td>
        </tr>
        <tr v-if="isInterchangePlus">
          <th class="rate-caption">Premium Transaction Fee:</th>
          <td><currency-input v-model="value.premiumTransactionFeeCost" :lt="value.premiumTransactionFeeLimit"/></td>
          <td><currency-input v-model="value.premiumTransactionFeeCurrent" :lt="value.premiumTransactionFeeLimit" :gt="value.premiumTransactionFeeCost"/></td>
          <td><currency-input v-model="value.premiumTransactionFeeLimit" :gt="value.premiumTransactionFeeCost"/></td>
        </tr>
        <tr>
          <th class="rate-caption">Debit Transaction Fee:</th>
          <td><currency-input v-model="value.debitTransactionFeeCost" :lt="value.debitTransactionFeeLimit"/></td>
          <td><currency-input v-model="value.debitTransactionFeeCurrent" :lt="value.debitTransactionFeeLimit" :gt="value.debitTransactionFeeCost"/></td>
          <td><currency-input v-model="value.debitTransactionFeeLimit" :gt="value.debitTransactionFeeCost"/></td>
        </tr>
        <tr>
          <th class="rate-caption">EBT Transaction Fee:</th>
          <td><currency-input v-model="value.ebtTransactionFeeCost" :lt="value.ebtTransactionFeeLimit"/></td>
          <td><currency-input v-model="value.ebtTransactionFeeCurrent" :lt="value.ebtTransactionFeeLimit" :gt="value.ebtTransactionFeeCost"/></td>
          <td><currency-input v-model="value.ebtTransactionFeeLimit" :gt="value.ebtTransactionFeeCost"/></td>
        </tr>
        <tr>
          <th class="rate-caption">AVS Fee:</th>
          <td><currency-input v-model="value.avsFeeCost" :lt="value.avsFeeLimit"/></td>
          <td><currency-input v-model="value.avsFeeCurrent" :lt="value.avsFeeLimit" :gt="value.avsFeeCost"/></td>
          <td><currency-input v-model="value.avsFeeLimit" :gt="value.avsFeeCost"/></td>
        </tr>
        <tr>
          <th class="rate-caption">Batch Fee:</th>
          <td><currency-input v-model="value.batchFeeCost" :lt="value.batchFeeLimit"/></td>
          <td><currency-input v-model="value.batchFeeCurrent" :lt="value.batchFeeLimit" :gt="value.batchFeeCost"/></td>
          <td><currency-input v-model="value.batchFeeLimit" :gt="value.batchFeeCost"/></td>
        </tr>
        <tr>
          <th class="rate-caption">Voice Authorization Fee:</th>
          <td><currency-input v-model="value.voiceAuthFeeCost" :lt="value.voiceAuthFeeLimit"/></td>
          <td><currency-input v-model="value.voiceAuthFeeCurrent" :lt="value.voiceAuthFeeLimit" :gt="value.voiceAuthFeeCost"/></td>
          <td><currency-input v-model="value.voiceAuthFeeLimit" :gt="value.voiceAuthFeeCost"/></td>
        </tr>
        <!--
        <tr>
          <th class="rate-caption">Account Setup Fee:</th>
          <td><currency-input v-model="value.accountSetupFeeCost" :lt="value.accountSetupFeeLimit"/></td>
          <td><currency-input v-model="value.accountSetupFeeCurrent" :lt="value.accountSetupFeeLimit" :gt="value.accountSetupFeeCost"/></td>
          <td><currency-input v-model="value.accountSetupFeeLimit" :gt="value.accountSetupFeeCost"/></td>
        </tr>
      -->
        <tr>
          <th class="rate-caption">Monthly Fee:</th>
          <td><currency-input v-model="value.monthlyFeeCost" :lt="value.monthlyFeeLimit"/></td>
          <td><currency-input v-model="value.monthlyFeeCurrent" :lt="value.monthlyFeeLimit" :gt="value.monthlyFeeCost"/></td>
          <td><currency-input v-model="value.monthlyFeeLimit" :gt="value.monthlyFeeCost"/></td>
        </tr>
        <!--
        <tr>
          <th class="rate-caption">Annual Fee:</th>
          <td><currency-input v-model="value.annualFeeCost" :lt="value.annualFeeLimit"/></td>
          <td><currency-input v-model="value.annualFeeCurrent" :lt="value.annualFeeLimit" :gt="value.annualFeeCost"/></td>
          <td><currency-input v-model="value.annualFeeLimit" :gt="value.annualFeeCost"/></td>
        </tr>
      -->
        <tr>
          <th class="rate-caption">Chargeback Fee:</th>
          <td><currency-input v-model="value.chargebackFeeCost" :lt="value.chargebackFeeLimit"/></td>
          <td><currency-input v-model="value.chargebackFeeCurrent" :lt="value.chargebackFeeLimit" :gt="value.chargebackFeeCost"/></td>
          <td><currency-input v-model="value.chargebackFeeLimit" :gt="value.chargebackFeeCost"/></td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
import BooleanInput from '@blockchyp/crud-magic/src/widgets/BooleanInput'
import DateInput from '@blockchyp/crud-magic/src/widgets/DateInput'
import BasisPointsInput from '@blockchyp/crud-magic/src/widgets/BasisPointsInput'
import CurrencyInput from '@blockchyp/crud-magic/src/widgets/CurrencyInput'

export default {
  name: 'pricing-policy',
  props: ['value', 'costCaption', 'costHelp', 'currentCaption', 'currentHelp', 'limitCaption', 'limitHelp'],
  components: {
    'boolean-input': BooleanInput,
    'date-input': DateInput,
    'bps-input': BasisPointsInput,
    'currency-input': CurrencyInput
  },
  data: () => ({
  }),
  mounted () {
  },
  computed: {
    isInterchangePlus: function () {
      return this.value.type === 'interchange'
    },
    isFlatRate: function () {
      return this.value.type === 'flatrate'
    },
    effectiveCostCaption: function () {
      if (this.costCaption) {
        return this.costCaption
      } else {
        return "Cost"
      }
    },
    effectiveCurrentCaption: function () {
      if (this.currentCaption) {
        return this.currentCaption
      } else {
        return "Current"
      }
    },
    effectiveLimitCaption: function () {
      if (this.limitCaption) {
        return this.limitCaption
      } else {
        return "Limit"
      }
    }
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

.rate-help {
  font-weight: normal;
  font-size: 8pt;
}

.rate-caption {
  padding-top: 25px;
  font-weight: bold;
}

</style>
