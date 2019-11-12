<template>
  <div class="text-center">
    <img src="/images/devices/8500i.png" alt="Image" v-if="terminalType == 'LUXE8000'"/>
    <img src="/images/devices/6500m.png" alt="Image" v-if="terminalType == 'LUXE6200'"/>
  </div>
</template>

<script>
export default {
  name: 'terminal-thumbnail',
  props: ['terminalType'],
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
