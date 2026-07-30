<template>
  <!-- 这里放置你的侨宝图标 -->
  <div class="qiaobao" ref="qiaobao">
    <img class="qiaobao_icon" src="../../assets/ai/qiaobao8.gif" alt="侨宝" @click="toggleDialog" />
    <transition name="trivia">
      <div class="hp_trivia_inner" v-show="openDialog" :class="{ focusin: focusin }">
      </div>
    </transition>
  </div>
</template>

<script>

export default {
  name: 'QiaoBaoAI',
  data() {
    return {
      openDialog: false,
      focusin: false
    }
  },
  methods: {
    toggleDialog() {
      this.openDialog = !this.openDialog
      if (this.openDialog) {
        setTimeout(() => {
          this.focusin = !this.focusin
        }, 800)
      } else {
        this.focusin = !this.focusin
      }
    }
  }
}
</script>

<style scoped>
.qiaobao {
  position: fixed;
  left: 3%;
  top: 80%;
  width: 10vw;
  height: 15vh;
  z-index: 1000;
  max-width: 100px;
  max-height: 150px;
  color: rgba(0, 0, 0, 0.87);
  font-family: "PingFang SC", "Microsoft YaHei", SimHei;
  font-weight: 400;
  font-size: 14px;
  line-height: 24px;
}

.qiaobao:hover {
  cursor: pointer;
}

.qiaobao .qiaobao_icon {
  width: 100px;
  height: 120px;
  z-index: 1001;
  position: absolute;
  left: 5px;
  bottom: 25%;
}

.hp_trivia_inner {
  border-radius: 2%;
  width: 0;
  height: 0;
  background-color: hsl(240 31% 95% / 1);
  position: absolute;
  bottom: 0;
  left: 5%;
}

/* 固定窗口大小的类 */
.hp_trivia_inner.focusin {
  width: 680px;
  height: 520px;
}

/* 过渡效果 */
.trivia-enter-active,
.trivia-leave-active {
  transition: all 1s ease;
}

.trivia-enter-from,
.trivia-leave-to {
  width: 0;
  height: 0;
  opacity: 0;
}

.trivia-enter-to,
.trivia-leave-from {
  width: 680px;
  height: 520px;
  opacity: 1;
}
</style>
