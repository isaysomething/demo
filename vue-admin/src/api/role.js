import request from '@/utils/request'

export function queryRoles() {
  return request({
    url: '/roles',
    method: 'get'
  })
}

export function createRole(data) {
  return request({
    url: '/roles',
    method: 'post',
    data
  })
}

export function updateRole(name, data) {
  return request({
    url: `/roles/${name}`,
    method: 'put',
    data
  })
}

export function deleteRole(name) {
  return request({
    url: `/roles/${name}`,
    method: 'delete'
  })
}

export function getRole(name) {
  return request({
    url: `/roles/${name}`,
    method: 'get'
  })
}

export function queryPermissions() {
  return request({
    url: '/permissions',
    method: 'get'
  })
}

export function queryPermissionGroups() {
  return request({
    url: '/permission-groups',
    method: 'get'
  })
}
