<template>
  <div class="input-group">
    <input type="text"
      class="form-control fc-datepicker dateWidget"
      :class="{'validated-control': validated}"
      placeholder="MM/DD/YYYY"
      ref="input"
      :id="fieldId"
      v-bind:value="value"
      :name="fieldId"
      @keyup="handleKeyUp"
      @blur="onBlur"
      @change="emitValue"
      @input="emitValue"
      @keypress="handleKeyPress"
      maxlength="10"
      :required="isRequired"/>
      <span class="input-group-btn">
        <button class="btn bd bd-l-0 bg-white tx-gray-600" type="button"><i class="fa fa-calendar-alt"></i></button>
      </span>
  </div>
</template>

<script>
export default {
  name: 'date-input',
  props: ['value', 'field', 'mutable', 'required', 'id', 'validated'],
  data: () => ({
  }),
  mounted () {
    let self = this
    $('#' + self.id).datepicker({
       showOtherMonths: true,
       selectOtherMonths: true,
       onSelect: function(date) {
         self.$nextTick( function() {
           self.$emit('input', date)
         })
        }
     });
  },
  computed: {
    fieldId: function () {
      if (this.id) {
        return this.id
      } else if (this.field) {
        return this.field.id
      } else {
        return this.id
      }
    },
    isRequired: function () {
      if (this.required) {
        return this.required
      } else if (this.field) {
        return this.field.required
      } else {
        return false
      }
    },
    isMutable: function () {
      if (this.mutable) {
        return this.mutable
      } else if (this.field) {
        return this.field.mutable
      } else {
        return false
      }
    }
  },
  methods: {
    isNumeric: function (n) {
      return !isNaN(parseFloat(n)) && isFinite(n);
    },
    onBlur: function (e) {
      let ts = Date.parse(e.target.value)
      if (Number.isNaN(ts)) {
        this.$emit('validationError', this, 'Invalid Date')
      } else {
        this.emitValue(e)
      }
    },
    emitValue: function (e) {
      this.$emit('input', e.target.value)
    },
    handleKeyPress: function(e) {

      if (e.key === 'Delete') {
        return
      } else if (e.key == 'Backspace') {
        return
      } else if (!this.isNumeric(e.key)) {
        e.preventDefault()
      }

    },
    formatValue: function (currentString) {

      if (!currentString) {
        return
      }
      let baseString = ''

      for (var i = 0; i < currentString.length; i++) {
        let char = currentString.charAt(i)

        if (this.isNumeric(char)) {
          baseString = baseString + char
          switch (baseString.length) {
            case 2:
            case 5:
              baseString = baseString + '/'
              break;
          }
        }
      }

      this.$emit('input', baseString)
    },
    handleKeyUp: function (e) {
      this.formatValue(e.target.value)
    },
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
