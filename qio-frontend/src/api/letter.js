import request from '@/utils/request'

export function generateLetter(data) {
  return request({
    url: '/letter/generateLetter',
    method: 'post',
    data: data
  })
}

export function sendLetter(data) {
  return request({
    url: '/letter/sendLetter',
    method: 'post',
    data: data
  })
}

export function getMyNotReadLetter() {
  return request({
    url: '/letter/getMyNotReadLetter',
    method: 'get'
  })
}
export function getMyReceiveLetter() {
  return request({
    url: '/letter/getMyReceiveLetter',
    method: 'get'
  })
}

export function getMySendLetter() {
  return request({
    url: '/letter/getMySendLetter',
    method: 'get'
  })
}

export function readLetter(letterId) {
  return request({
    url: '/letter/readLetter/' + letterId,
    method: 'put'
  })
}
export function getFontPaperLimit() {
  return request({
    url: '/paper/getFontPaperLimit',
    method: 'get'
  })
}
