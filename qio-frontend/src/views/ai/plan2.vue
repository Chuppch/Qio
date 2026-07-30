<template>
  <!-- 这里放置你的侨宝图标 -->
  <div class="qiaobao" ref="qiaobao">
    <img class="qiaobao_icon" src="../../assets/ai/qiaobao8.gif" alt="侨宝" @click="openDialog = !openDialog" />
    <transition name="trivia">
      <div class="hp_trivia_inner" v-show="openDialog">
      </div>
    </transition>
  </div>
</template>

<script>
export default {
  name: 'QiaoBaoAI',
  data() {
    return {
      openDialog: false
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
  /* letter-spacing: 0px; */
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
  width: 680px;
  height: 520px;
  position: absolute;
  bottom: 0;
  left: 5%;
  transform-origin: 0 100%;
  background-color: hsl(240 31% 95% / 1);

}

/* 过渡效果 */
.trivia-enter-active,
.trivia-leave-active {
  transition: all 1s ease;
}

.trivia-enter-from,
.trivia-leave-to {
  transform: scale(0);
  opacity: 0;
}

.trivia-enter-to,
.trivia-leave-from {
  transform: scale(1);
  opacity: 1;
}
</style>
