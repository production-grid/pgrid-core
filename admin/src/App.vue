<template>
  <div id="app">
    <div class="container">
      <router-view :session="session" class="view" @sessionInvalidated="handleSessionInvalidation"></router-view>
    </div>
  </div>
</template>

<script>
import APIService from '@/services/APIService'
import router from './router'
import ApplicationService from '@/services/ApplicationService'

export default {
  name: 'app',
  data: () => ({
    session: {}
  }),
  methods: {
    handleSessionInvalidation: function () {
      let self = this
      APIService.get('/api/security/session', function (response) {
        self.session = response.data
        router.push(ApplicationService.defaultRoutePath(self.session))
      })
    }
  },
  mounted () {
    this.handleSessionInvalidation()
  }
}
</script>

<style>

</style>
