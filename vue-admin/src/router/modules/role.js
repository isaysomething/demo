/** When your routing table is too long, you can split it into small modules**/

import Layout from '@/layout'

const roleRouter =
{
  path: '/role',
  component: Layout,
  redirect: '/role/index',
  meta: {
    title: 'role',
    icon: 'role'
  },
  children: [
    {
      path: 'index',
      component: () => import('@/views/role/index'),
      name: 'Role',
      meta: { title: 'role', icon: 'role', affix: true }
    },
    {
      path: 'create',
      component: () => import('@/views/role/create'),
      name: 'CreateRole',
      meta: { title: 'createRole', icon: 'edit' },
      hidden: true
    },
    {
      path: 'edit/:id(\\w+)',
      component: () => import('@/views/role/edit'),
      name: 'EditRole',
      meta: { title: 'editRole', noCache: true, activeMenu: '/role/list' },
      hidden: true
    }
  ]
}

export default roleRouter
