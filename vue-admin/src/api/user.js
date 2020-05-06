import request from '@/utils/request'

export function login(data) {
  return request({
    url: '/user/login',
    method: 'post',
    data
  })
}

export function getInfo() {
  return request({
    url: '/user/info',
    method: 'get'
  })
}

export function logout() {
  return request({
    url: '/user/logout',
    method: 'post'
  })
}

export function queryUsers(query) {
  return request({
    url: '/users',
    method: 'get',
    params: query
  })
}

export function getUser(id) {
  return request({
    url: `/users/${id}`,
    method: 'get'
  })
}

export function createUser(data) {
  return request({
    url: '/users',
    method: 'post',
    data: data
  })
}

export function updateUser(id, data) {
  return request({
    url: `/users/${id}`,
    method: 'put',
    data: data
  })
}

export function deleteUser(id) {
  return request({
    url: `/users/${id}`,
    method: 'delete'
  })
}

