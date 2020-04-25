/** When your routing table is too long, you can split it into small modules **/

import Layout from '@/layout'

const postRouter =
{
  path: '/post',
  component: Layout,
  redirect: '/post/index',
  name: 'Post',
  meta: {
    title: 'post',
    icon: 'post'
  },
  children: [
    {
      path: 'index',
      component: () => import('@/views/post/index'),
      name: 'Post',
      meta: { title: 'post', icon: 'post', affix: true }
    },
    {
      path: 'create',
      component: () => import('@/views/post/create'),
      name: 'CreatePost',
      meta: { title: 'createPost', icon: 'edit' },
      hidden: true
    },
    {
      path: 'edit/:id(\\d+)',
      component: () => import('@/views/post/edit'),
      name: 'EditPost',
      meta: { title: 'editPost', noCache: true, activeMenu: '/post/list' },
      hidden: true
    }
  ]
}

export default postRouter
