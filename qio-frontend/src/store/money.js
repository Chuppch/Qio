// store.js
import Vue from 'vue'
import Vuex from 'vuex'
import { getUserMoney } from '@/api/user' // 确保导入了API调用函数

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    money: 0 // 用户的余额
  },
  getters: {
    getMoney: state => state.money // 获取用户的余额
  },
  mutations: {
    SET_MONEY(state, money) {
      state.money = money // mutation来更新用户的余额
    }
  },
  actions: {
    async fetchUserMoney({ commit }) {
      const response = await getUserMoney() // 调用API获取用户余额
      if (response && response.data) {
        commit('SET_MONEY', response.money) // 假设返回的数据中money是用户余额
      }
    }
  }
})
