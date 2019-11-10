<template>
  <crud-list :globalState="globalState"
  breadcrumbCategory="Security"
  resource="security/admin/users"
  clickRoute="/admin-user-edit"
  newRoute="/new-admin-user"
  searchEnabled="true"
  includeDeletedOption="true"/>
</template>

<script>
import CrudList from '@/crud/CrudList'

export default {
  name: 'admin-user-list',
  props: ['globalState'],
  components: {
    'crud-list': CrudList
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
