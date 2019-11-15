<template>
  <crud-list :globalState="globalState"
  breadcrumbCategory="Security"
  resource="security/admin/groups"
  clickRoute="/admin-group-edit"/>
</template>

<script>
import CrudList from '@/crud/CrudList'

export default {
  name: 'admin-group-list',
  props: ['globalState'],
  components: {
    'crud-list': CrudList
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
