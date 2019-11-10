import HelloWorld from '@/components/HelloWorld'
import Login from '@/modules/security/Login'
import AdminUserList from '@/modules/security/AdminUserList'

export default [
  {
    path: '/hello',
    name: 'HelloWorld',
    component: HelloWorld,
    props: true,
    nav: {
      section: 'Security',
      sectionIndex: 1,
      sectionIconClass: 'fas fa-shield-alt',
      caption: 'Hello!',
      index: 0,
      permission: "security.perm.user",
      allRequired: [],
      anyRequired: []
    }
  },
  {
    path: '/users',
    name: 'AdminUserList',
    component: AdminUserList,
    props: true,
    nav: {
      section: 'Security',
      caption: 'Users',
      index: 1,
      permission: "security.perm.admin",
      allRequired: ["security.perm.admin","security.perm.user"],
      anyRequired: ["security.perm.admin","security.perm.user"]
    }
  },
  {
    path: '/groups',
    name: 'HelloWorld',
    component: HelloWorld,
    props: true,
    nav: {
      section: 'Security',
      caption: 'Groups',
      permission: "security.perm.admin"
    }
  },
  {
    path: '/tenants',
    name: 'HelloWorld',
    component: HelloWorld,
    props: true,
    nav: {
      section: 'Security',
      caption: 'Tenants',
      permission: "security.perm.admin",
    }
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
    props: true
  }
]
