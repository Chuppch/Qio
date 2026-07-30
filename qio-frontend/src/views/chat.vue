<template>
  <div style="padding-left: 20px; width: 100%; height: 100vh; background-color: #444444;">

    <div style="padding-top: 100px; text-align: left;">
      <el-input type="textarea" @keyup.native.enter="send" style="width: 300px; margin-right: 10px;"
        v-model="message"></el-input>
      <el-button @click="send">发送</el-button>
    </div>
    <div ref="chatContainer" style="overflow-y: auto; width: 500px; height: 65vh;">
      <div id="hello" class="response" style="text-align: left; "></div>
    </div>
  </div>
</template>

<script>
import useUserStore from '@/store/modules/user'
import { marked } from 'marked'
import hljs from 'highlight.js'
import 'highlight.js/styles/foundation.css'
import { MessageBox } from 'element-ui'
// import 'highlight.js/styles/monokai-sublime.css'
const rendererMD = new marked.Renderer()
marked.setOptions({
  renderer: rendererMD,
  highlight: function (code) {
    return hljs.highlightAuto(code).value
  },
  pedantic: false,
  gfm: true,
  tables: true,
  breaks: false,
  sanitize: false,
  smartLists: true,
  smartypants: false,
  xhtml: false
})

// 自定义节流函数
function throttle(func, limit) {
  let inThrottle
  return function () {
    const args = arguments
    const context = this
    if (!inThrottle) {
      func.apply(context, args)
      inThrottle = true
      setTimeout(() => (inThrottle = false), limit)
    }
  }
}

export default {
  name: 'App',
  data() {
    return {
      text: 'Hello, World!<br>',
      timeOutCount: 0,
      message: '',
      chatCount: 0,
      linked: false,
      answer: '',
      isRelogin: false
    }
  },
  created() {
    this.initWebSocket()
  },
  methods: {
    initWebSocket() {
      if ('WebSocket' in window) {
        if (!useUserStore().token) {
          this.isRelogin = true
          MessageBox.confirm('登录状态已过期,要先登录才可以体验完整功能哦~🥳', '侨缘信使', { confirmButtonText: '登录', cancelButtonText: '取消', type: 'warning' }).then(() => {
            this.isRelogin = false
            useUserStore().logOut()
            this.$router.push({
              path: '/login'
            })
          }).catch(() => {
            this.isRelogin = false
          })
          return
        }
        // 这里记得要改成你自己的ip
        if (!this.linked) {
          this.websocket = new WebSocket('ws://localhost:8080/ws/chat', useUserStore().token)
          this.websocket.onerror = this.onError
          this.websocket.onopen = this.onOpen
          this.websocket.onmessage = this.onMessage
          this.websocket.onclose = this.onClose
          this.websocketTimeout = setTimeout(() => {
            if (this.websocket.readyState !== WebSocket.OPEN) {
              this.websocket.close()
            }
          }, 2000) // 设置超时时间为3000毫秒
        }
        setTimeout(() => {
          if (!this.linked) {
            this.websocket = new WebSocket('ws://localhost:8080/ws/chat', useUserStore().token)
            // 为了方便写了这一坨屎删代码，有效解决各环境下不能调用的问题
            this.websocket.onerror = this.onError
            this.websocket.onopen = this.onOpen
            this.websocket.onmessage = this.onMessage
            this.websocket.onclose = this.onClose
            this.websocketTimeout = setTimeout(() => {
              if (this.websocket.readyState !== WebSocket.OPEN) {
                this.websocket.close()
                this.onError()
              }
            }, 2000)
          } // 设置超时时间为3000毫秒
        }, 2500)
      } else {
        alert('Not support websocket')
      }
    },
    addUnloadListener() {
      window.addEventListener('beforeunload', this.handleBeforeUnload)
    },
    removeUnloadListener() {
      window.removeEventListener('beforeunload', this.handleBeforeUnload)
    },
    handleBeforeUnload() {
      if (this.websocket) {
        this.websocket.close()
      }
    },
    onError() {
      this.timeOutCount++
      if (!this.linked && this.timeOutCount > 1) {
        this.$message.error('连接超时,请刷新页面重试!')
        this.websocket = null
      }
    },
    onOpen() {
      // this.$message.success('开始写信吧!')
      // this.send()
      // setTimeout(() => {
      //   this.websocket.send(JSON.stringify(this.chat))
      // }, 500)
      this.linked = true
      this.chatCount++
    },
    onClose() {

    },
    setMessageInnerHTML(innerHTML) {
      this.answer += innerHTML
      const div = document.getElementById('response' + this.chatCount)
      if (div != null) {
        div.innerHTML = marked(this.answer)
      } else {
        const responseDiv = document.createElement('div')
        responseDiv.innerHTML = marked(this.answer)
        responseDiv.id = 'response' + this.chatCount
        responseDiv.className = 'response'
        const lastUserDiv = document.querySelector('.user:last-of-type')
        if (lastUserDiv) {
          lastUserDiv.after(responseDiv)
        } else {
          document.getElementById('hello').after(responseDiv)
        }
      }
      this.scrollToBottom()
    },
    onMessage(event) {
      // console.log(event)
      try {
        const response = JSON.parse(event.data)
        if (response.code === 200) {
          this.$message.success(response.msg)
          this.chatCount--
          return
        } else if (response.code === 500) {
          this.$message.error(response.msg)
          this.chatCount--
          return
        }
        // console.log(response)
      } catch (e) {
        // event.data is not JSON, handle accordingly if needed
      }

      if (this.chatCount === 0) {
        document.getElementById('hello').innerHTML = event.data
      } else {
        this.setMessageInnerHTML(event.data)
      }
    },
    send: throttle(function () {
      this.chatCount++
      const userDiv = document.createElement('div')
      userDiv.id = 'user' + this.chatCount
      userDiv.className = 'user'
      userDiv.innerHTML = this.message
      if (this.chatCount === 1) {
        document.getElementById('hello').after(userDiv)
      } else {
        // console.log(this.chatCount);
        (document.getElementById('user' + this.chatCount) || document.getElementById('response' + (this.chatCount - 1))).after(userDiv)
      }
      this.websocket.send(JSON.stringify({
        message: this.message
      }))
      this.message = ''
      this.answer = ''
      this.scrollToBottom()
      // 做了一个50ms的节流,防止频繁发送,这里可以根据实际情况调整,缓解服务器压力
    }, 50),
    closeWebSocket() {
      if (this.websocket) {
        this.websocket.close()
      }
    },
    scrollToBottom() {
      this.$nextTick(() => {
        const chatContainer = this.$refs.chatContainer
        chatContainer.scrollTop = chatContainer.scrollHeight
      })
    }
  },
  beforeDestroy() {
    this.removeUnloadListener()
    this.closeWebSocket()
  }
}
</script>

<style>
.response {
  margin-top: 10px;
  text-align: left;
  color: #fff;
}

.user {
  margin-top: 10px;
  text-align: left;
  color: aqua;
}

::-webkit-scrollbar {
  opacity: 0;
  display: none;
  width: 10px;
}
</style>
