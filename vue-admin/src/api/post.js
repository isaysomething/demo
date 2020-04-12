import request from '@/utils/request'

export function fetchList(query) {
  return request({
    url: '/posts',
    method: 'get',
    params: query
  })
}
export function fetchPost(id) {
  return request({
    url: '/posts/' + id,
    method: 'get'
  })
}

export function createPost(data) {
  return request({
    url: '/posts',
    method: 'post',
    data
  })
}

export function updatePost(id, data) {
  return request({
    url: '/posts/' + id,
    method: 'PUT',
    data
  })
}
