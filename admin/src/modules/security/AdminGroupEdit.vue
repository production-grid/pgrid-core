<template>
  <crud-form
    :globalState="globalState"
    resource="security/admin/groups"
    :id="id"
    breadcrumbCategory="Security"
    listRoute="/admin-groups"
    :extraRoutes="extraRoutes"
    :session="session"
    :newForm="newForm"/>
</template>

<script>
import CrudForm from '@/crud/CrudForm'

export default {
  name: 'admin-group-edit',
  props: ['globalState', 'id', 'session'],
  components: {
    'crud-form': CrudForm
  },
  computed: {
    extraRoutes: function () {
      let routes = []
      return routes
    }
  },
  data: () => ({
    newForm: {
      id: null,
      permissions: new Array()
    }
  })
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
