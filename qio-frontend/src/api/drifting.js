import request from '@/utils/request'
// 获得漂流瓶
export function getDriftBottle() {
  return request({
    url: '/bottle/getDriftBottle',
    method: 'get'
  })
}
// 展示漂流瓶
export function showDriftBottle() {
  return request({
    url: '/bottle/showDriftBottle',
    method: 'get'
  })
}
// 生成漂流瓶
export function generateDriftBottle(data) {
  return request({
    url: '/bottle/generateDriftBottle',
    method: 'post',
    data: data
  })
}
// 再捡一个
export function throwBackBottle() {
  return request({
    url: '/bottle/ThrowBack',
    method: 'put'
  })
}
// 成为好友
export function sendFriendRequest(data) {
  return request({
    url: '/bottle/sendFriendRequest',
    method: 'post',
    data: data
  })
}
// 获取好友列表
export function ProcessingFriendRequests() {
  return request({
    url: '/bottle/ProcessingFriendRequests',
    method: 'get'
  })
}
// 同意好友申请
export function BecomeFriend(data) {
  return request({
    url: '/bottle/BecomeFriend',
    method: 'put',
    data: data
  })
}
