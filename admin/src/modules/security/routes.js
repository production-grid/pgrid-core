import HelloWorld from '@/components/HelloWorld'
import Login from '@/modules/security/Login'

export default [
  {
    path: '/hello',
    name: 'HelloWorld',
    component: HelloWorld
  },
  {
    path: '/login',
    name: 'Login',
    component: Login
  }
]
