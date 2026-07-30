import request from '@/utils/request'
// 获取纸张商城列表
export function fetchPaperList() {
  return request({
    url: '/paper/list',
    method: 'get'
  })
}
// 获取字体商城列表
export function fetchFontList() {
  return request({
    url: '/font/list',
    method: 'get'
  })
}
// 获取墨水列表
export function fetchColorList() {
  return request({
    url: '/font/listColor',
    method: 'get'
  })
}
// 买纸
export function buyPaper(paperId) {
  return request({
    url: '/paper/buyPaper',
    method: 'post',
    params: { paperId }
  })
}
// 买字体
export function buyFont(fontId) {
  return request({
    url: '/font/buyFont',
    method: 'post',
    params: { fontId }
  })
}
// 买字体颜色
export function buyFontColor(fontColorId) {
  return request({
    url: '/font/buyFontColor',
    method: 'post',
    params: { fontColorId }
  })
}
// 获取功能卡商城列表
export function fetchCardList() {
  return request({
    url: '/card/list',
    method: 'get'
  })
}
// 买功能卡
export function buyCard(cardId) {
  return request({
    url: '/card/buyCard',
    method: 'post',
    params: { cardId }
  })
}
