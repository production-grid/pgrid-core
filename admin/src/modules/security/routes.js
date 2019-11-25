import Login from '@/modules/security/Login'
import AdminUserList from '@/modules/security/AdminUserList'
import AdminUserEdit from '@/modules/security/AdminUserEdit'
import TenantList from '@/modules/security/TenantList'
import TenantEdit from '@/modules/security/TenantEdit'
import AdminGroupList from '@/modules/security/AdminGroupList'
import AdminGroupEdit from '@/modules/security/AdminGroupEdit'

export default [
  {
    path: '/admin-users',
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
    path: '/admin-groups',
    name: 'AdminGroupList',
    component: AdminGroupList,
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
  },
  {
    path: '/admin-group-edit',
    name: 'admin-group-edit',
    component: AdminGroupEdit,
    props: true
  },
  {
    path: '/admin-group-edit/:id',
    name: 'admin-group-edit',
    component: AdminGroupEdit,
    props: true
  },
  {
    path: '/admin-user-edit',
    name: 'admin-user-edit',
    component: AdminUserEdit,
    props: true
  },
  {
    path: '/admin-user-edit/:id',
    name: 'admin-user-edit',
    component: AdminUserEdit,
    props: true
  }
]
