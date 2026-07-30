import request from '@/utils/request'

export function userLoginPage() {
  return request({
    url: '/question/userLoginPage',
    method: 'get'
  })
}
export function startAnswer(setId) {
  return request({
    url: '/question/startAnswer',
    method: 'post',
    params: { setId }
  })
}
export function submitAnswers(data) {
  return request({
    url: '/question/submitAnswers',
    method: 'post',
    data: data
  })
}
export function allAnswerToFront(setId) {
  return request({
    url: '/question/allAnswerToFront',
    method: 'get',
    params: { setId }
  })
}
