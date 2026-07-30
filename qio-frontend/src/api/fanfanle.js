import request from '@/utils/request'

// 赢得翻翻乐
export function winFfl() {
  return request({
    url: '/game/winFfl',
    method: 'get'
  })
}
// 限制翻翻乐次数
export function getFflLimit() {
  return request({
    url: '/game/getFflLimit',
    method: 'get'
  })
}
