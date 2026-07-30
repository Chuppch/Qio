<template>
  <div class="container">
    <div class="login-wrapper">
      <!-- 登录表单 -->
      <div class="header">登录</div>
      <div class="form-wrapper">
        <input type="text" @keyup.enter="handleCodePre" v-model="loginForm.username" name="username"
          placeholder="邮箱/用户名" class="input-item" autocomplete="off" ref="usernameInput">
        <input type="password" @keyup.enter="handleCodePre" v-model="loginForm.password" name="password"
          placeholder="密码" class="input-item">
        <router-link to="/forget" class="forget">忘记密码</router-link>
        <button @click="handleCodePre" class="btn">登录</button>
      </div>
      <p>没有账号？<router-link to="/register">去注册</router-link></p>
    </div>
    <div v-if="codeOpen" class="overlay" @click="codeOpen = false"></div>
    <div v-if="codeOpen" class="codeCheck">
      <div class="code-header">
        <span>输入验证码</span>
        <button class="close-btn" @click="codeOpen = false">×</button>
        <div class="code-content">
          <input ref="codeInput" type="text" v-model="loginForm.code" @keyup.enter="handleLogin" placeholder="验证码"
            class="inputcode" style="width: 200px; height: 50px;" />
          <div class="picbox" style="margin-left: 10px;">
            <img :src="captchaImg" alt="验证码" @click="handleGetPicCode">
          </div>
        </div>
        <button @click="handleLogin" class="checkbtn">验证</button>
      </div>
    </div>
  </div>
</template>
<script>
import { getPicCode } from '@/api/user'
import useUserStore from '@/store/modules/user'

export default {
  data() {
    return {
      codeOpen: false,
      captchaImg: '',
      loginForm: {
        username: '',
        password: '',
        code: '',
        uuid: ''
      }
    }
  },
  methods: {
    handleCheckForm() {
      if (!this.loginForm.username) {
        this.$message.error('请输入用户名')
        return false
      }
      if (!this.loginForm.password) {
        this.$message.error('请输入密码')
        return false
      }
      return true
    },
    handleCodePre() {
      if (this.handleCheckForm()) {
        this.handleGetPicCode() // 组件创建时获取验证码
        this.codeOpen = true
      }
    },
    handleGetPicCode() {
      getPicCode().then(res => {
        this.captchaImg = 'data:image/gif;base64,' + res.data.img
        this.loginForm.uuid = res.data.uuid
      })
    },
    handleLogin() {
      const userStore = useUserStore()
      this.codeOpen = false
      userStore.login(this.loginForm)
        .then(response => {
          // 测试登录成功
          console.log('Login successful:', response)
          // 跳转到首页
          this.$router.push({ name: 'Index' }) // Vue 2 中使用 this.$router 进行路由跳转
        }).catch(error => {
          // 测试登录失败
          console.log('Login failed:', error)
          this.handleGetPicCode() // 登录失败后重新获取验证码
        })
    }
  },
  watch: {
    codeOpen(newVal) {
      if (newVal) {
        this.$nextTick(() => {
          this.$refs.codeInput.focus()
        })
      }
    }
  },
  mounted() {
    // 页面加载完成后，将焦点设置到输入框
    this.$refs.usernameInput.focus()
  }
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
  /* height: 100%; */
  /* max-height: 100vh; */
  height: 850px;
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
  line-height: 200px;
}

.input-item {
  display: block;
  width: 358px;
  margin-bottom: 20px;
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
  /* 设置 placeholder 的字体 */
  font-size: 20px;
  color: #000000;
}

.email::placeholder {
  text-transform: uppercase;
  font-family: '方正姚体', sans-serif;
  /* 设置 placeholder 的字体 */
  font-size: 20px;
  color: #000000;
}

.pic {
  display: block;
  width: 50%;
  margin-bottom: 20px;
  border: 0;
  padding: 10px;
  border-bottom: 1px solid rgb(128, 125, 125);
  font-size: 15px;
  outline: none;
}

.btn {
  position: absolute;
  left: 50%;
  transform: translate(-50%, 0);
  text-align: center;
  width: 200px;
  margin-top: 60px;
  background-color: #B73E13;
  color: #fff;
  border-radius: 20px;
  border: 0;
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
  bottom: 80px;
  left: 150px;
  font-family: '方正姚体', sans-serif;
  /* 设置 placeholder 的字体 */
  font-size: 20px;
}

a {
  text-decoration-line: none;
  font-size: 18px;
}

.remember {
  float: left;
  font-size: 12px;
  vertical-align: middle;
  margin-top: 5px;
}

span {
  color: #A47348;
}

.forget {
  float: right;
  font-size: 16px;
  color: #000000;
  font-family: '方正姚体', sans-serif;
}

.picbox {
  position: absolute;
  top: 60px;
  right: 50px;
  display: inline-block;
  width: 20%;
  margin-top: 40px;
}

.picbox img {
  width: 100%;
}

.codeCheck {
  /* margin-top: 100px; */
  background-color: #fff;
  border-radius: 10px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
  width: 300px;
  height: 150px;
  padding: 50px 50px;
  position: absolute;
  left: 50%;
  top: 50%;
  z-index: 999;
  transform: translate(-50%, -50%);
}

.inputcode {
  display: block;
  width: 358px;
  padding: 20px;
  margin-top: 40px;
  box-sizing: border-box;
  font-size: 15px;
  outline: none;
  border: 1px #666 solid;
  border-radius: 10px;
  background-color: rgba(0, 0, 0, 0);
}

.checkbtn {
  position: absolute;
  top: 150px;
  left: 50%;
  transform: translate(-50%, 0);
  text-align: center;
  width: 200px;
  margin-top: 20px;
  background-color: #B73E13;
  color: #fff;
  border-radius: 20px;
  border: 0;
  display: inline-block;
  cursor: pointer;
  color: #ffffff;
  font-family: Arial;
  font-size: 17px;
  padding: 16px 31px;
  text-decoration: none;
}

.overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  /* 黑色半透明遮罩 */
  z-index: 998;
  /* 确保遮罩层在验证码弹框之下 */
}

.code-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.code-header span {
  position: absolute;
  top: 50px;
  left: 50px;
  font-size: 18px;
  font-weight: bold;
}

.close-btn {
  position: absolute;
  top: 20px;
  right: 20px;
  background: none;
  border: none;
  cursor: pointer;
  font-size: 24px;
}
</style>
