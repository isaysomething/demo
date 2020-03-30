import request from '@/utils/request'

export function captcha() {
  return request({
    url: '/captcha',
    method: 'post'
  })
}

export function checkCaptcha(id, captcha) {
  return request({
    url: '/check-captcha',
    method: 'post',
    data: {
      id: id,
      captcha: captcha
    }
  })
}
