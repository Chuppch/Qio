import request from '@/utils/request'
export function list(page, limit) {
  return request({
    url: '/marketing/list',
    method: 'get',
    params: { page, limit }
  })
}
