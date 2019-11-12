<template>
  <crud-list :globalState="globalState"
  breadcrumbCategory="Security"
  resource="security/tenants"
  clickRoute="/tenant-edit"/>
</template>

<script>
import CrudList from '@/crud/CrudList'

export default {
  name: 'tenant-list',
  props: ['globalState'],
  components: {
    'crud-list': CrudList
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
