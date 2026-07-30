<template>
  <div class="container">
    <div class="login-wrapper">
      <div class="header">找回密码</div>
      <div class="form-wrapper">
        <input v-model="email" type="text" name="username" placeholder="邮箱号" class="input-item" autocomplete="off"
          ref="usernameInput" @keyup.enter="sendCode">
        <span class="alert" v-if="isMistake">请输入正确的邮箱号</span>
        <input v-model="password" type="password" name="password" placeholder="新密码:长度必须为6-20个字符" class="input-item"
          @keyup.enter="sendCode">
        <input v-model="confirmPassword" type="password" name="repassword" placeholder="再次确认密码" class="input-item"
          @keyup.enter="sendCode">
        <span class="isSame" v-if="isSame">两次密码不一致请重新输入</span>
        <span class="isSame" v-if="isEmpty">请输入长度为6-20个字符的密码</span>
        <input type="text" v-model="code" name="repassword" placeholder="邮箱验证码" class="email" @keyup.enter="register">
        <button @click="sendCode" v-if="!loading" :disabled="second !== totalSecond" class="code">
          {{ second === totalSecond ? '获取验证码' : second + '秒后重新发送' }}
        </button>
        <button @click="sendCode" v-else :disabled="second !== totalSecond" class="code">
          <i class="el-icon-loading"></i>
        </button>
        <button @click="register" class="btn">重置密码</button>
        <p>已有账号？<router-link to="/login">去登录</router-link></p>
      </div>
    </div>
  </div>
</template>

<script>
import { Message } from 'element-ui'
import axios from 'axios'

export default {
  name: 'ForgetCode',
  data() {
    return {
      isEmpty: false,
      loading: false,
      isMistake: false,
      totalSecond: 60, // 总秒数
      second: 60, // 当前秒数
      timer: null, // 定时器 id
      email: '', // 邮箱号
      password: '', // 密码
      confirmPassword: '', // 再次确认密码
      code: '',
      uuid: ' ',
      responseMessage: null,
      isSame: false
    }
  },
  methods: {
    // 校验邮箱是否合法
    validFn() {
      if (!/^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$/.test(this.email)) {
        this.isMistake = true
        return false
      } else {
        this.isMistake = false
        return true
      }
    },
    validatePassword() {
      if (!this.password || !this.confirmPassword) {
        this.isEmpty = true
        return false
      }
      if (this.password.length < 6 || this.password.length > 20) {
        this.isEmpty = true
        return false
      }
      if (this.password !== this.confirmPassword) {
        this.isSame = true
        return false
      } else {
        this.isSame = false
        return true
      }
    },
    destroyed() {
      clearInterval(this.timer)
    },
    // 开始验证码倒计时
    startCountdown() {
      this.second = this.totalSecond
      if (!this.timer && this.second === this.totalSecond) {
        this.timer = setInterval(() => {
          this.second--
          if (this.second <= 0) {
            clearInterval(this.timer)
            this.timer = null
            this.second = this.totalSecond
          }
        }, 1000)
      }
    },
    async sendCode() {
      try {
        if (!this.validFn()) {
          throw new Error('无效的邮箱地址')
        }
        if (!this.validatePassword()) {
          throw new Error('请检查密码')
        }
        this.loading = true
        const response = await axios.get('/api/user/sendResetPasswordCode', {
          params: { email: this.email }
        })
        if (response.data.code === 200) {
          this.responseMessage = response.data.msg || '验证码已发送，请检查您的邮箱。'
          Message.success(this.responseMessage)
          this.startCountdown() // 开始倒计时
        } else {
          this.responseMessage = response.data.msg
          Message.error(this.responseMessage)
        }
        this.loading = false
      } catch (error) {
        console.error('Error:', error)
        this.responseMessage = error.message || '发送验证码失败'
        Message.error(this.responseMessage)
        this.loading = false
      }
    },
    // 注册（重置密码）
    async register() {
      try {
        if (!this.validFn() || !this.validatePassword()) {
          throw new Error('请检查邮箱地址和密码')
        }
        if (!this.code) {
          throw new Error('请输入验证码')
        }
        const response = await axios.post('/api/user/resetPasswordByEmail', {
          username: this.email,
          password: this.password,
          code: this.code,
          confirmPassword: this.confirmPassword
        })
        if (response.data.code === 200) {
          this.responseMessage = response.data.msg || '重置成功'
          Message.success(this.responseMessage)
          this.$router.push('/login')
        } else {
          this.responseMessage = response.data.msg
          Message.error(this.responseMessage)
        }
      } catch (error) {
        console.error('Error:', error)
        this.responseMessage = error.message || '重置失败'
        Message.error(this.responseMessage)
      }
    }
  },
  mounted() {
    // 页面加载完成后，将焦点设置到输入框
    this.$refs.usernameInput.focus()
  }
  // computed: {
  //   // 计算属性，用于判断表单是否有效
  //   isFormInvalid() {
  //     return !this.validFn() || !this.validatePassword() || this.isMistake || this.isSame
  //   }
  // }
}
</script>

<style scoped>
* {
  margin: 0;
  padding: 0;
}

html {
  height: 100%;
}

body {
  height: 100%;
}

.container {
  position: relative;
  height: 50rem;
  width: 100%;
  background-image: url(../../assets/imgss/loginbgd.webp);
  background-position: center center;
  background-size: cover;
  background-repeat: no-repeat;
}

.login-wrapper {
  background-color: rgba(255, 255, 255, 0.6);
  width: 358px;
  height: 588px;
  border-radius: 80px;
  padding: 0 50px;
  position: absolute;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
}

.header {
  font-size: 38px;
  font-weight: bold;
  text-align: center;
  line-height: 160px;
}

.input-item {
  display: block;
  width: 358px;
  margin-bottom: 10px;
  border: 0;
  padding: 20px;
  box-sizing: border-box;
  border-bottom: 2px solid rgb(0, 0, 0);
  font-size: 15px;
  outline: none;
  background-color: rgba(0, 0, 0, 0);
}

.email {
  display: block;
  width: 160px;
  margin-bottom: 10px;
  border: 0;
  padding: 20px;
  box-sizing: border-box;
  border-bottom: 2px solid rgb(0, 0, 0);
  font-size: 15px;
  outline: none;
  background-color: rgba(0, 0, 0, 0);
}

.input-item::placeholder {
  text-transform: uppercase;
  font-family: '方正姚体', sans-serif;
  font-size: 20px;
  color: #000000;
}

.email::placeholder {
  text-transform: uppercase;
  font-family: '方正姚体', sans-serif;
  font-size: 20px;
  color: #000000;
}

.btn {
  text-align: center;
  position: absolute;
  left: 50%;
  transform: translate(-50%, 0);
  text-align: center;
  width: 200px;
  margin-top: 20px;
  background-color: #B73E13;
  color: #fff;
  border-radius: 20px;
  border: 1px solid #ffffff;
  display: inline-block;
  cursor: pointer;
  color: #ffffff;
  font-family: Arial;
  font-size: 17px;
  padding: 16px 31px;
  text-decoration: none;
}

.msg {
  text-align: center;
  line-height: 88px;
}

p {
  position: absolute;
  bottom: 30px;
  left: 150px;
  font-family: '方正姚体', sans-serif;
  font-size: 20px;
}

a {
  text-decoration-line: none;
  font-size: 18px;
}

.code {
  position: absolute;
  top: 390px;
  right: 90px;
  width: 110px;
  height: 40px;
  border: 0px;
  background-color: #B73E13;
  color: #fff;
  border-radius: 15px;
}

.alert {
  position: absolute;
  top: 220px;
  left: 60px;
  color: rgba(168, 4, 4, 0.873);
  font-size: 10px;
}

.isSame {
  position: absolute;
  top: 360px;
  left: 60px;
  color: rgba(168, 4, 4, 0.873);
  font-size: 10px;
}
</style>
