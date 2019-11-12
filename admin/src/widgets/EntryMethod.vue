<template>
  <span>
    <i v-if="entryMethod === 'TAP'" class="fas fa-wifi fa-rotate-90"></i>
    <i v-if="entryMethod === 'CHIP'" class="fa fa-credit-card-front"></i>
    <i v-if="entryMethod === 'SWIPE'" class="fa fa-credit-card"></i>
    <i v-if="entryMethod === 'KEYED'" class="fas fa-calculator"></i>
    <i v-if="entryMethod === 'SAVED'" class="fas fa-save"></i>
    {{entryMethod}}
  </span>
</template>

<script>
export default {
  name: 'entry-method',
  props: ['entryMethod'],
  data: () => ({
  }),
  mounted () {
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
img {
  max-height: 28px;
}

</style>
