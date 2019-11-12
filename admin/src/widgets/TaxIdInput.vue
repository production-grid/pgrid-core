<template>
  <div class="input-group">
    <input type="text" v-if="isMutable"
      class="form-control"
      :class="{'validated-control': validated}"
      ref="input"
      :id="elementId"
      v-bind:value="value"
      :name="elementId"
      @input="emitValue"
      @keypress="handleKeyPress"
      @keyup="handleKeyUp"
      :minlength="maxValueLength"
      :maxlength="maxValueLength"
      :required="isRequired"/>
    <div class="input-group-append" v-if="addOnCaption">
      <button class="btn btn-info" @click="sendClickEvent()" type="button">{{addOnCaption}}</button>
    </div>
    <span v-if="!isMutable">{{value}}</span>
  </div>
</template>

<script>
export default {
  name: 'tax-id-input',
  props: ['value', 'required', 'taxIdType', 'field', 'addOnCaption', 'id', 'validated'],
  data: () => ({
  }),
  watch: {
    'value': function(value) {
      this.formatValue(value)
    },
    'taxIdType': function(value) {
      this.formatValue(this.value)
    }
  },
  computed: {
    maxValueLength: function () {
      switch (this.taxIdType) {
        case 'SSN':
          return 11
        case 'EIN':
          return 10
      }
    },
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
      } else if (!this.isNumeric(e.key)) {
        e.preventDefault()
      }

    },
    formatValue: function (currentString) {

      if (!currentString) {
        return
      }
      let baseString = ''
      let digitCount = 0

      switch (this.taxIdType) {
        case 'SSN':
          for (var i = 0; i < currentString.length; i++) {
            let char = currentString.charAt(i)

            if (this.isNumeric(char)) {
              switch (digitCount) {
              case 3:
                baseString = baseString + '-'
                break
              case 5:
                baseString = baseString + '-'
                break
              }
              baseString = baseString + char
              digitCount++
            }
          }
          break;
        case 'EIN':
        for (var i = 0; i < currentString.length; i++) {
          let char = currentString.charAt(i)

          if (this.isNumeric(char)) {
            switch (digitCount) {
            case 2:
              baseString = baseString + '-'
              break
            }
            baseString = baseString + char
            digitCount++
          }
        }
        break;
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

</style>
