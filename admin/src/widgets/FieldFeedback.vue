<template>
  <div>
    <div class="error-feedback" v-if="hasErrors">
      <ul>
        <li v-for="error in errorMap[errorKey]">{{error}}</li>
      </ul>
    </div>
    <small class="help-text" v-if="!hasErrors">{{helpText}}</small>
  </div>
</template>

<script>
export default {
  name: 'field-feedback',
  props: ['errorMap', 'errorKey', 'helpText'],
  data: () => ({

  }),
  computed: {
    hasErrors: function () {
      return this.errorMessages && this.errorMessages.length > 0
    },
    errorMessages: function () {
      return this.errorMap[this.errorKey]
    }
  },
  mounted () {
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

.error-feedback {
  color: red;
  font-size: 8pt;
}

</style>
