<template>
  <div class="background-image">
    <div class="envelope"><img src="../../assets/imgss/envelope1.webp" class="img-3d" @click="shakeAndMoveImage"
        ref="shakeImage"></div>
  </div>
</template>

<script>
export default {
  name: 'ImportIndex',
  methods: {
    shakeAndMoveImage() {
      const img = this.$refs.shakeImage
      if (img) {
        // 添加抖动类
        img.classList.add('shaking')
      }
      setTimeout(() => {
        this.$router.push('/index')
      }, 1000) // 1000毫秒 = 1秒
    }
  }
}
</script>

<style scoped>
.background-image {
  position: relative;
  width: 100%;
  height: 900px;
  max-width: 100%;
  /* 限制图片最大宽度 */
  transition: transform 0.3s;
  /* 设置动画过渡效果 */
  background-image: url(../../assets/imgss/backimg.webp);
  background-size: cover;
  background-position: center center;
}

.background-image::before {
  content: "";
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(to bottom, rgba(0, 0, 0, 0) 0%, rgba(0, 0, 0, 0.9) 100%);
}

.envelope {
  position: absolute;
  width: 400px;
  height: 600px;
  top: 50%;
  left: 50%;
  margin-left: -230px;
  margin-top: -250px;
}

.img-3d {
  width: 100%;
}

.envelope:hover {
  transform: rotateX(20deg) rotateY(330deg);
  cursor: pointer;
}

/* 鼠标悬停时改变Z旋转
/* 添加抖动动画的关键帧 */
@keyframes shake {
  0% {
    transform: translateY(0);
  }

  100% {
    transform: translateY(-700px);
  }
}

.shaking {
  animation: shake 1.2s cubic-bezier(.36, .07, .19, .97) both;
}
</style>
