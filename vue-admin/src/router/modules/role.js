/** When your routing table is too long, you can split it into small modules**/

import Layout from '@/layout'

const roleRouter =
{
  path: '/role',
  component: Layout,
  redirect: '/role/index',
  children: [
    {
      path: 'index',
      component: () => import('@/views/role/index'),
      name: 'Role',
      meta: { title: 'role', icon: 'role', affix: true, permissions: ['role:list'] }
    },
    {
      path: 'create',
      component: () => import('@/views/role/create'),
      name: 'CreateRole',
      meta: { title: 'createRole', icon: 'edit', permissions: ['role:create'] },
      hidden: true
    },
    {
      path: 'edit/:id(\\w+)',
      component: () => import('@/views/role/edit'),
      name: 'EditRole',
      meta: { title: 'editRole', noCache: true, activeMenu: '/role/list', permissions: ['role:edit'] },
      hidden: true
    }
  ]
}

export default roleRouter
