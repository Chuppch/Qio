<template>
  <div class="banner">
    <div class="left">
      <img src="../../assets/imgs/exploreTitle.png" alt="" class="title">
      <div class="point">
        <router-link to="/explore"><img src="../../assets/imgs/point1.png" alt=""
            style="margin-bottom: 30px;"></router-link>
        <img src="../../assets/imgs/point2.png" alt="" @click="unlockWarning">
        <img src="../../assets/imgs/point2.png" alt="" @click="unlockWarning">
        <img src="../../assets/imgs/point2.png" alt="" @click="unlockWarning">
      </div>
    </div>
    <div class="right">
      <img src="../../assets/imgss/mengmei.webp" alt="">
    </div>
  </div>
</template>

<script>
import { MessageBox } from 'element-ui'
export default {
  name: 'QiaopiImportStory',
  methods: {
    unlockWarning() {
      MessageBox.alert('此章节待解锁', '提示', {
        confirmButtonText: '确定',
        type: 'warning'
      })
    }
  }
}
</script>

<style scoped>
.banner {
  margin: 0;
  padding: 0;
}

.left {
  position: relative;
}

.title {
  position: absolute;
  left: 0;
  top: 20px;
}

.point {
  position: absolute;
  top: 100px;
  left: 200px;
  display: flex;
  flex-direction: column;
}

.point img {
  width: 400px;
  height: 100px;
  margin-bottom: 50px;
  cursor: pointer;
}

.right {
  width: 900px;
  height: 750px;
  position: relative;
  top: 0;
  right: -550px;
}

.right img {
  position: absolute;
  right: 0;
  width: 730px;
  height: 745px;
  margin: 0;
}
</style>
