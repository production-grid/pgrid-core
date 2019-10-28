<template>
  <div id="app">
    <div class="slim-header with-sidebar"  v-if="session.effectivePermissions">
       <div class="container-fluid">
         <div class="slim-header-left">
           <h2 class="slim-logo"><tenant-logo :session="session"/></h2>
           <a href="" id="slimSidebarMenu" class="slim-sidebar-menu"><span></span></a>
         </div><!-- slim-header-left -->
         <div class="slim-header-right">
           <div class="dropdown dropdown-c">
             <a href="#" class="logged-user" data-toggle="dropdown">
               <span>{{session.firstName}} {{session.lastName}}</span>
               <i class="fa fa-angle-down"></i>
             </a>
             <div class="dropdown-menu dropdown-menu-right">
               <nav class="nav">
                 <router-link class="nav-link" to="/change-password"><i class="icon ion-ios-gear"></i> Change Password</router-link>
                 <router-link class="nav-link" to="/logout"><i class="icon ion-forward"></i> Sign Out</router-link>
               </nav>
             </div><!-- dropdown-menu -->
           </div><!-- dropdown -->
         </div><!-- header-right -->
       </div><!-- container-fluid -->
     </div><!-- slim-header -->
     <div class="slim-body" v-if="navRoutes">
       <div class="slim-sidebar" v-if="session.effectivePermissions">
         <label class="sidebar-label">Navigation</label>
         <ul class="nav nav-sidebar" v-for="section in navRoutes" v-bind:key="section.caption">
           <li class="sidebar-nav-item with-sub">
             <a href="#" @click="makeActive(section.id)" :id="section.id" class="sidebar-nav-link" v-bind:class="{ active: isActive(section.id) }"><i :class="section.iconClass"></i> {{section.id}}</a>
             <ul class="nav sidebar-nav-sub">
               <li v-for="route in section.routes" v-bind:key="route.path" class="nav-sub-item"><router-link :to="route.path" class="nav-sub-link">{{route.nav.caption}}</router-link></li>
             </ul>
           </li>
         </ul>
       </div> <!-- slim-sidebar -->

       <div id="content">
         <router-view :session="session" class="view" @sessionInvalidated="handleSessionInvalidation"></router-view>
       </div>
      </div>
  </div> <!-- app -->
</template>

<script>
import APIService from '@/services/APIService'
import router from './router'
import ApplicationService from '@/services/ApplicationService'
import TenantLogo from '@/components/TenantLogo'

export default {
  name: 'app',
  components: {
    'tenant-logo': TenantLogo
  },
  data: () => ({
    session: {},
    navRoutes: null,
    selectedMenu : 'operationsMenu'
  }),
  methods: {
    handleSessionInvalidation: function () {
      let self = this
      APIService.get('/api/security/session', function (response) {
        self.session = response.data
        self.navRoutes = ApplicationService.visibleNavRoutes(self.session)
        router.push(ApplicationService.defaultRoutePath(self.session))
      })
    },
    makeActive: function (menu) {
      this.selectedMenu = menu
    },
    isActive: function (menu) {
      return menu === this.selectedMenu
    }
  },
  mounted () {
    this.handleSessionInvalidation()
  }
}
</script>

<style>

#content {
  padding: 20px;
}

</style>
