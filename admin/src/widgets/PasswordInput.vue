<template>
  <div>
    <input type="password"
     v-bind:class="{ 'is-invalid': !passwordValid }"
      name="password"
      placeholder="*************"
      class="form-control"
      ref="input"
      v-bind:value="value"
      @input="emitValue"/>
    <div class="invalid-feedback" v-if="!passwordValid">
      <ul>
        <li v-if="!hasDigit">Must include a number.</li>
        <li v-if="!hasLetter">Must include a letter.</li>
        <li v-if="!minPasswordLength">Must be at least seven characters long.</li>
        <li v-if="invalidCharacters">Contains invalid character.</li>
      </ul>
    </div>
  </div>
</template>

<script>
export default {
  name: 'password-input',
  props: ['value', 'field'],
  data: () => ({
    minPasswordLength: false,
    hasDigit: false,
    hasLetter: false,
    invalidCharacters: false,
    passwordValid: true
  }),
  mounted () {
  },
  watch: {
    value: function (val) {
      this.minPasswordLength = val.length >= 7
      this.hasDigit = this.containsNumber(val)
      this.hasLetter = this.containsLetter(val)
      this.invalidCharacters = this.constainsInvalidCharacter(val)
      this.passwordValid = this.minPasswordLength && this.hasDigit && this.hasLetter && !this.invalidCharacters
    }
  },
  methods: {
    containsLetter: function (val) {
      return /[A-Za-z]/.test(val)
    },
    containsNumber: function (val) {
      return /\d/.test(val)
    },
    constainsInvalidCharacter: function (val) {
      return /\s/.test(val)
    },
    emitValue: function (e) {
      this.$emit('input', e.target.value)
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
