import axios from 'axios'
import { Message, MessageBox } from 'element-ui'
import { getToken } from '@/utils/auth'
import cache from '@/plugins/cache'
import { tansParams } from '@/utils/qiaopi'
import errorCode from './errorCode'
import useUserStore from '@/store/modules/user'
// import router from '@/router'

// 是否显示重新登录
export const isRelogin = { show: false }
// 创建axios实例
const service = axios.create({
  method: 'get',
  baseURL: process.env.VUE_APP_API, // .env中的VUE_APP_API参数
  headers: { 'Content-Type': 'application/json;charset=utf-8' },
  timeout: 30000 // 超时时间
})

// 请求拦截器
service.interceptors.request.use((config) => {
  // 是否需要设置 token
  const isToken = (config.headers || {}).isToken === false
  // 是否需要防止数据重复提交
  const isRepeatSubmit = (config.headers || {}).repeatSubmit === false
  if (getToken() && !isToken) {
    config.headers.Authorization = getToken() // 让每个请求携带自定义token 请根据实际情况自行修改
  }
  // get请求映射params参数
  if (config.method === 'get' && config.params) {
    let url = config.url + '?' + tansParams(config.params)
    url = url.slice(0, -1)
    config.params = {}
    config.url = url
  }
  // 防止数据重复提交
  if (!isRepeatSubmit && (config.method === 'post' || config.method === 'put')) {
    const requestObj = {
      url: config.url,
      data: typeof config.data === 'object' ? JSON.stringify(config.data) : config.data,
      time: new Date().getTime()
    }
    const requestSize = Object.keys(JSON.stringify(requestObj)).length // 请求数据大小
    const limitSize = 5 * 1024 * 1024 // 限制存放数据5M
    if (requestSize >= limitSize) {
      // console.warn(`[${config.url}]: ` + '请求数据大小超出允许的5M限制，无法进行防重复提交验证。')
      return config
    }
    const sessionObj = cache.session.getJSON('sessionObj')
    if (sessionObj === undefined || sessionObj === null || sessionObj === '') {
      cache.session.setJSON('sessionObj', requestObj)
    } else {
      const surl = sessionObj.url // 请求地址
      const sdata = sessionObj.data // 请求数据
      const stime = sessionObj.time // 请求时间
      const interval = 1000 // 间隔时间(ms)，小于此时间视为重复提交
      if (sdata === requestObj.data && requestObj.time - stime < interval && surl === requestObj.url) {
        const message = '数据正在处理，请勿重复提交'
        // console.warn(`[${surl}]: ` + message)
        return Promise.reject(new Error(message))
      } else {
        cache.session.setJSON('sessionObj', requestObj)
      }
    }
  }
  // console.log(config)
  return config
}, (error) => {
  // console.log(error)
  Promise.reject(error)
}
)

let errorMessageShown = false
// 响应拦截器
service.interceptors.response.use(res => {
  // 未设置状态码则默认成功状态
  const code = res.data.code || 200
  // 获取错误信息
  const msg = errorCode[code] || res.data.msg || errorCode.default
  // 二进制数据则直接返回
  if (res.request.responseType === 'blob' || res.request.responseType === 'arraybuffer') {
    return res.data
  }

  if (code === 401) {
    if (!isRelogin.show) {
      isRelogin.show = true
      MessageBox.confirm('登录状态已过期,要先登录才可以体验完整功能哦~🥳', '侨缘信使', { confirmButtonText: '登录', cancelButtonText: '取消', type: 'warning' }).then(() => {
        isRelogin.show = false
        useUserStore().logOut().then(() => {
          location.href = '/login'
        })
      }).catch(() => {
        isRelogin.show = false
      })
    }
    // eslint-disable-next-line
    return Promise.reject('无效的会话，或者会话已过期，请重新登录。')
  } else if (code === 500) {
    Message({ message: msg, type: 'error' })
    return Promise.reject(new Error(msg))
  } else if (code === 601) {
    Message({ message: msg, type: 'warning' })
    return Promise.reject(new Error(msg))
  } else if (code !== 200) {
    // this.$notify.error({ title: msg })
    Message({ message: msg, type: 'error' })
    // eslint-disable-next-line
    return Promise.reject('error')
  } else {
    return Promise.resolve(res.data)
  }
}, error => {
  // console.log(111)
  // console.log('err' + error)
  let { message } = error

  if (error && error.response) {
    switch (error.response.status) {
      case 401:
        if (!isRelogin.show) {
          isRelogin.show = true
          MessageBox.confirm('登录状态已过期,要先登录才可以体验完整功能哦~🥳', '侨缘信使', { confirmButtonText: '登录', cancelButtonText: '取消', type: 'warning' }).then(() => {
            isRelogin.show = false
            useUserStore().logOut()
            location.href = '/login'
          }).catch(() => {
            isRelogin.show = false
          })
        }
        // router.push('/login')
        // location.href = '/login'
        message = error.response.data.msg || '未授权的访问'
        break
      case 405:
        message = '请求错误'
    }
  }

  if (message === 'Network Error') {
    message = '后端接口连接异常'
  } else if (message.includes('timeout')) {
    message = '系统接口请求超时'
  } else if (message.includes('Request failed with status code')) {
    message = '系统接口' + message.substr(message.length - 3) + '异常'
  }
  // Message({ message: message, type: 'error', duration: 5 * 1000 })

  // Throttle error messages to show only one within 1 second
  if (!errorMessageShown) {
    errorMessageShown = true
    Message({ message: message, type: 'error', duration: 5 * 1000 })
    setTimeout(() => {
      errorMessageShown = false
    }, 1000)
  }
  return Promise.reject(error)
}
)
export default service
