import { getUserInfo, login } from '@/api/user'
import { getToken, setToken, removeToken } from '@/utils/auth'
import defAva from '@/assets/default-avatar.webp'
import { defineStore } from 'pinia'

// eslint-disable-next-line
const useUserStore = defineStore(
  'user',
  {
    state: () => ({
      token: getToken(),
      id: '',
      name: '',
      avatar: ''
      // money: ''
    }),
    actions: {
      // 登录
      login(userInfo) {
        const username = userInfo.username.trim()
        const password = userInfo.password
        const code = userInfo.code
        const uuid = userInfo.uuid
        return new Promise((resolve, reject) => {
          login(username, password, code, uuid).then(res => {
            const data = res.data
            setToken(data.token)
            this.token = data.token
            // console.log(data)
            const avatar = (data.avatar === '' || data.avatar == null) ? defAva : data.avatar
            this.id = data.id
            this.name = data.nickname
            this.avatar = avatar

            // this.money = data.money
            resolve()
          }).catch(error => {
            reject(error)
          })
        })
      },
      // 获取用户信息
      getUserInfo() {
        return new Promise((resolve, reject) => {
          getUserInfo().then(res => {
            const data = res.data
            const avatar = (data.avatar === '' || data.avatar == null) ? defAva : data.avatar
            this.id = data.id
            this.name = data.nickname
            this.avatar = avatar
            // this.money = data.money
            resolve(res)
          }).catch(error => {
            reject(error)
          })
        })
      },
      // 退出系统
      logOut() {
        this.token = ''
        removeToken()
        // return new Promise((resolve, reject) => {
        //   logout(this.token).then(() => {
        //   this.token = ''
        //   removeToken()
        //   resolve()
        //   }).catch(error => {
        //     reject(error)
        //   })
        // })
      }
    }
  })

export default useUserStore
