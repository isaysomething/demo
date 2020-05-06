/** When your routing table is too long, you can split it into small modules **/

import Layout from '@/layout'

const postRouter =
{
  path: '/post',
  component: Layout,
  redirect: '/post/index',
  children: [
    {
      path: 'index',
      component: () => import('@/views/post/index'),
      name: 'Post',
      meta: { title: 'post', icon: 'post', affix: true, permissions: ['post:list'] }
    },
    {
      path: 'create',
      component: () => import('@/views/post/create'),
      name: 'CreatePost',
      meta: { title: 'createPost', icon: 'edit', permissions: ['post:create'] },
      hidden: true
    },
    {
      path: 'edit/:id(\\d+)',
      component: () => import('@/views/post/edit'),
      name: 'EditPost',
      meta: { title: 'editPost', noCache: true, activeMenu: '/post/list', permissions: ['post:edit'] },
      hidden: true
    }
  ]
}

export default postRouter
