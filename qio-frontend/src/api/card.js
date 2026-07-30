import request from '@/utils/request'

export function useCard(data) {
  // useCard function is used to use a card
  return request({
    url: '/card/useCard',
    method: 'put',
    data: data
  })
}
