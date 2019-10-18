import Vue from 'vue'
import Router from 'vue-router'
import ApplicationService from '@/services/ApplicationService'
import securityRoutes from '@/modules/security/routes'

Vue.use(Router)

let routes = new Array()
routes = ApplicationService.addRoutes(routes, securityRoutes)



export default new Router({
  routes: routes
})
