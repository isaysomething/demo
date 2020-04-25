/** When your routing table is too long, you can split it into small modules **/

import Layout from '@/layout'

const userRouter =
{
  path: '/user',
  component: Layout,
  redirect: '/user/index',
  meta: {
    title: 'user',
    icon: 'user'
  },
  children: [
    {
      path: 'user',
      component: () => import('@/views/user/index'),
      name: 'User',
      meta: { title: 'user', icon: 'user', affix: true }
    },
    {
      path: 'create',
      component: () => import('@/views/user/create'),
      name: 'CreateUser',
      meta: { title: 'createUser', icon: 'add' },
      hidden: true
    },
    {
      path: 'edit/:id(\\d+)',
      component: () => import('@/views/user/edit'),
      name: 'EditUser',
      meta: { title: 'editUser', noCache: true, activeMenu: '/user/list' },
      hidden: true
    }
  ]
}

export default userRouter
