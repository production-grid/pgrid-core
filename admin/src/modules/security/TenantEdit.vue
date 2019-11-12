<template>
  <crud-form
    :globalState="globalState"
    resource="security/tenants"
    :id="id"
    breadcrumbCategory="Security"
    listRoute="/tenants"
    :extraRoutes="extraRoutes"
    :session="session"/>
</template>

<script>
import CrudForm from '@/crud/CrudForm'

export default {
  name: 'tenant-edit',
  props: ['globalState', 'id', 'session'],
  components: {
    'crud-form': CrudForm
  },
  computed: {
    extraRoutes: function () {
      let routes = []
      routes.push({
        route: '/preview-site',
        caption: 'Preview Site'
      })
      return routes
    }
  },
  data: () => ({
  })
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
