<template>
  <div>
    <div class="input-group">
    <span class="form-control is-disabled" v-if="!isMutable">{{value}}</span>
    <input type="text" v-if="isMutable"
      class="form-control"
      :class="{'validated-control': validated, 'is-invalid': !valid}"
      ref="input"
      :id="elementId"
      v-bind:value="value"
      :name="elementId"
      :readonly="readOnly"
      @input="emitValue"
      @keypress="handleKeyPress"
      @keyup="handleKeyUp"
      @blur="validate"
      :required="isRequired"/>
    <div class="input-group-append">
      <span class="input-group-text" v-if="!percentMode">bps</span>
      <span class="input-group-text" v-if="percentMode">%</span>
    </div>
  </div>
  <div class="error-feedback" v-if="!valid">
    <ul>
      <li v-for="error in errors">{{error}}</li>
    </ul>
  </div>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'bps-input',
  props: ['value', 'required', 'field', 'addOnCaption', 'id', 'validated', 'disabled', 'lt', 'gt', 'readOnly', 'percentMode'],
  data: () => ({
    valid: true,
    errors: []
  }),
  watch: {
    'value': function(value) {
      this.formatValue(value)
    }
  },
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
      if (this.disabled) {
        return false
      } else if (this.field) {
        return this.field.mutable
      } else {
        return true
      }
    }
  },
  methods: {
    sendClickEvent: function () {
      this.$emit('addOnClicked', '')
    },
    isNumeric: function (n) {
      return !isNaN(parseFloat(n)) && isFinite(n);
    },
    handleKeyPress: function(e) {

      if (e.key === 'Delete') {
        return
      } else if (e.key == 'Backspace') {
        return
      } else if (!this.isNumeric(e.key) && e.key !== '.' && e.key !== '%') {
        e.preventDefault()
      }

    },
    validate: function (e) {

      this.$emit('blur', e)

      if (this.percentMode) {
        return
      }

      let self = this
      self.errors = []

      if (this.value) {
        let url = '/api/bps?q=' + encodeURI(this.value)

        axios.get(url)
          .then(function(response) {
            if (response.data.success) {
              self.valid = true
              self.$emit('input', response.data.formattedValue)
              if (self.lt || self.gt) {
                let currentValue = parseFloat(response.data.formattedValue)
                if (self.lt) {
                  let ltValue = parseFloat(self.lt)
                  if (currentValue > ltValue) {
                    self.valid = false
                    self.errors.push("Must be less than " + self.lt + " bps.")
                  }
                }
                if (self.gt) {
                  let gtValue = parseFloat(self.gt)
                  if (currentValue < gtValue) {
                    self.valid = false
                    self.errors.push("Must be greater than " + self.gt + " bps.")
                  }
                }
              }
            } else {
              self.valid = false
            }
          })
          .catch(function (error) {
            alert(error)
          })
      }
    },
    formatValue: function (currentString) {

      if (!currentString) {
        return
      }
      let baseString = ''
      let digitCount = 0

      for (var i = 0; i < currentString.length; i++) {
        let char = currentString.charAt(i)

        if (this.isNumeric(char) || char === '.' || char === '%') {
          baseString = baseString + char
          digitCount++
        }
      }
      this.$emit('input', baseString)
    },
    handleKeyUp: function (e) {
      this.formatValue(e.target.value)
    },
    emitValue: function (e) {
      //this.$emit('input', e.target.value)
    }
  },
  mounted () {
    this.formatValue(this.value)
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

span.is-disabled {
  background-color: #EEEEEE;
}

.error-feedback {
  color: red;
  font-size: 8pt;
}

</style>
