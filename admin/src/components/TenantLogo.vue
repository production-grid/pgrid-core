<template>
  <img src="./../assets/pg-logo-horiz.png"/>
</template>

<script>
export default {
  name: 'tenant-logo',
  props: ['session']
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

 img {
   height: 25px;
   margin-top: -3px;
 }

</style>
