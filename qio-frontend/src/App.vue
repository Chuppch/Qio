<template>
  <div id="app">
    <canvas-animation></canvas-animation>
    <div v-if="isMobile" class="browserupgrade">
      <div class="logo"></div>
      <div class="tip">手机端暂未适配<br>请使用电脑端访问以获得最佳体验</div>
      <div class="tip-author"> -by 侨缘信使 </div>
      <a @click="isMobile = false">继续访问</a>
    </div>
    <div v-else class="title">
      <router-view></router-view>
      <qiaobao></qiaobao>
    </div>
  </div>
</template>

<script>
import Qiaobao from './views/ai/Qiaobao.vue'
export default {
  name: 'App',
  components: {
    Qiaobao
  },
  data() {
    return {
      isMobile: false
    }
  },
  created() {
    this.detectDevice()
  },
  methods: {
    detectDevice() {
      if (/Mobi|Android/i.test(navigator.userAgent)) {
        this.isMobile = true
      } else {
        this.isMobile = false
      }
    }
  }
}
</script>

<style>
img,
a {
  -webkit-user-drag: none;
}

.browserupgrade {
  position: absolute;
  left: 0;
  top: 0;
  width: 100%;
  height: 100%;
  background-color: #f19a2a;
  text-align: center;
  z-index: 9999;
  color: #fff;
}

.browserupgrade .logo {
  width: 70px;
  height: 70px;
  background-image: url(././assets/logo/logo.png);
  background-size: 100% 100%;
  background-position: center center;
  position: absolute;
  top: 40px;
  left: 40px;
}

.browserupgrade .tip {
  width: 80%;
  margin: 160px auto 40px;
  font-size: 24px;
  text-align: left;
}

.browserupgrade .tip-author {
  width: 80%;
  margin: 0 auto 70px;
  text-align: right;
  font-style: italic;
  font-size: 24px;
}

.browserupgrade a {
  display: block;
  width: 200px;
  border: 3px solid #fff;
  padding: 15px 10px;
  font-size: 25px;
  color: #fff;
  margin: 0 auto;
}

.title {
  text-align: center;
}

body {
  overflow-x: hidden;
  /* 隐藏水平滚动条 */
  margin: 0;
  padding: 0;
  width: 100%;
  height: 100%;
}

html,
body,
#app {
  width: 100%;
  height: 100%;
  margin: 0;
  padding: 0;
}

/* 更改滚动条默认样式 */
/* 滚动条整体轨道样式 */
::-webkit-scrollbar {
  width: 10px;
}

/* 滚动条滑块样式 */
::-webkit-scrollbar-thumb {
  background: #888;
  border-radius: 10px;
}

/* 滑动轨道按钮 */
::-webkit-scrollbar-button {
  width: 10px;
  height: 10px;
}

::-webkit-scrollbar-button:start:decrement {
  background-image: url(././assets/imgs/up.png);
  background-repeat: no-repeat;
  background-position: center;
  background-size: cover;
}

::-webkit-scrollbar-button:end:increment {
  background-image: url(././assets/imgs/down.png);
  background-repeat: no-repeat;
  background-position: center;
  background-size: cover;
}

/* 滚动条滑块悬停样式 */
::-webkit-scrollbar-thumb:hover {
  background: #555;
}

/* 滚动条轨道样式 */
::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 10px;
}

/* 滚动条轨道悬停样式 */
::-webkit-scrollbar-track:hover {
  background: #e1e1e1;
}

/* body, html {
  padding: 0;
  margin: 0;
  overscroll-behavior: none;
  overflow: hidden;
} */
</style>
