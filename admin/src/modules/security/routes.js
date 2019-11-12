import Login from '@/modules/security/Login'
import AdminUserList from '@/modules/security/AdminUserList'
import TenantList from '@/modules/security/TenantList'
import TenantEdit from '@/modules/security/TenantEdit'

export default [
  {
    path: '/users',
    name: 'AdminUserList',
    component: AdminUserList,
    props: true,
    nav: {
      index: 2,
      section: 'Security',
      caption: 'Users',
      permission: "security.perm.admin",
      allRequired: ["security.perm.admin","security.perm.user"],
      anyRequired: ["security.perm.admin","security.perm.user"]
    }
  },
  {
    path: '/groups',
    name: 'AdminUserList',
    component: AdminUserList,
    props: true,
    nav: {
      section: 'Security',
      caption: 'Groups',
      index: 3,
      permission: "security.perm.admin"
    }
  },
  {
    path: '/tenants',
    name: 'TenantList',
    component: TenantList,
    props: true,
    nav: {
      index: 1,
      section: 'Security',
      caption: '{{session.tenantPlural}}',
      permission: "security.perm.admin",
    }
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
    props: true
  },
  {
    path: '/tenant-edit',
    name: 'tenant-edit',
    component: TenantEdit,
    props: true
  },
  {
    path: '/tenant-edit/:id',
    name: 'tenant-edit',
    component: TenantEdit,
    props: true
  }
]
