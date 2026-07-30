<template>
  <div class="contain">
    <div class="banner">
      <section>
        <h2 id="text"><img src="../../assets/imgss/bgdlogo.webp" alt="" class="logo"></h2>
        <img src="../../assets/imgs/cloud1.png" alt="" class="cloud1">
        <img src="../../assets/imgs/cloud3.png" alt="" class="cloud2">
        <img src="../../assets/imgs/cloud2.png" alt="" class="cloud3">
        <!-- <img src="../../assets/imgss/index1.webp" alt="" id="sticker1">
        <img src="../../assets/imgss/index2.webp" alt="" id="sticker2"> -->
        <img src="../../assets/imgs/paper1.png" alt="" class="paper1">
        <img src="../../assets/imgs/paper2.png" alt="" class="paper2">
        <img src="../../assets/imgs/paper3.png" alt="" class="paper3">
        <img src="../../assets/imgs/paper4.png" alt="" class="paper4">
        <img src="../../assets/imgs/paper5.png" alt="" class="paper5">
        <img src="../../assets/imgs/paper6.png" alt="" class="paper6">
        <img src="../../assets/imgs/text1.png" alt="" class="text1">
        <img src="../../assets/imgs/text2.png" alt="" class="text2">
        <img src="../../assets/imgs/text3.png" alt="" class="text3">
        <img src="../../assets/imgs/text4.png" alt="" class="text4">
        <img src="../../assets/imgss/boat.webp" alt="" class="boat">
        <img src="../../assets/imgss/boatsea.webp" alt="" class="boatsea">
      </section>
      <div class="introduce">
        <div class="left">
          <img src="../../assets/imgss/qiaopiIntroduce.webp" alt="" class="title">
          <img src="../../assets/imgss/bell.webp" alt="" class="bell" @click="toggleContent">
          <div class="imgLoader"></div>
          <div class="container">
            <div class="credit">
            </div>

            <div class="book">
              <div class="gap"></div>
              <div class="pages">
                <div class="page"></div>
                <div class="page"></div>
                <div class="page"></div>
                <div class="page"></div>
                <div class="page"></div>
                <div class="page"></div>
              </div>
              <div class="flips">
                <div class="flip flip1">
                  <div class="flip flip2">
                    <div class="flip flip3">
                      <div class="flip flip4">
                        <div class="flip flip5">
                          <div class="flip flip6">
                            <div class="flip flip7"></div>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="right">
          <img src="../../assets/imgs/paper4.png" alt="" style="position: absolute;top: 200px;left: 700px;">
          <img src="../../assets/imgs/paper2.png" alt="" style="position: absolute;top: 100px;left: 600px;">
          <img src="../../assets/imgs/tree.png" alt="" class="tree">
          <img src="../../assets/imgs/blackboat.png" alt="" class="blackboat" :style="blackboatStyle">
          <p class="first"><span class="red">侨批</span>，又称“银信”，是19世纪末至20世纪初中国东南沿海地区华侨通过民间渠道寄回国内的<span
              class="yellow">汇款和家书</span>。</p>
          <p class="second">
            这些侨批不仅承载着华侨对家乡亲人的思念和关怀，也是他们经济支持的重要方式。侨批的内容通常包括家庭近况、生活琐事、乡情问候以及汇款金额等，它们是<span
              class="yellow">华侨与家乡之间情感和经济联系的纽带</span>。</p>
          <p class="third">
            侨批的传递往往依赖于专门的<span class="blue">“水客”</span>或<span
              class="yellow">“侨批局”</span>，这些信使和机构在没有现代化通信手段的时代，扮演了至关重要的角色。侨批的往来不仅反映了当时社会的经济状况，也展现了华侨群体的生活状态和文化特色。如今，侨批已成为研究华侨历史和文化的珍贵资料，它们见证了一段段跨越山海的家国情怀。
          </p>
          <img src="../../assets/imgss/facai.webp" alt="" class="facai">
        </div>
      </div>
      <div class="history">
        <h1><img src="../../assets/imgss/title2.webp" alt=""></h1>
        <div class="banner">
          <ul>
            <li v-for="history in histories" :key="history.id" :class="{ 'active': history.active }"
              @mouseenter="toggleActive(history)" @mouseleave="toggleInactive(history)">
              <div class="pic" :style="{ backgroundImage: 'url(' + history.img + ')' }">
                <span class="math">{{ history.math }}</span>
                <span class="title">{{ history.title }}</span>
                <p v-show="history.active">{{ history.description }}</p>
              </div>
            </li>
          </ul>
        </div>
      </div>
      <div class="function">
        <h1><img src="../../assets/imgss/functionbdg.webp" alt=""></h1>
        <div class="images-container">
          <div class="top-row">
            <router-link to="/drifting"><img src="../../assets/imgss/function1.webp" alt="" class="function-img"
                title="漂流瓶"></router-link>
            <router-link to="/shop"><img src="../../assets/imgss/function2.webp" alt="" class="function-img"
                title="兑换商城"></router-link>
            <router-link to="/import"><img src="../../assets/imgss/function3.webp" alt="" class="function-img"
                title="故事线探索"></router-link>
          </div>
          <div class="bottom-row">
            <router-link to="/memory"><img src="../../assets/imgss/function4.webp" alt="" class="function-img"
                title="游戏翻翻乐"></router-link>
            <router-link to="/write"><img src="../../assets/imgss/function5.webp" alt="" class="function-img"
                title="写一封侨批"></router-link>
            <router-link to="/receive"><img src="../../assets/imgss/function6.webp" alt="" class="function-img"
                title="查看收信"></router-link>
            <router-link to="/know"><img src="../../assets/imgss/function7.webp" alt="" class="function-img"
                title="知识问答"></router-link>
          </div>
        </div>
      </div>
    </div>

    <transition name="el-fade-in-linear">
      <div v-if="isVisible" class="el-image-viewer__wrapper" :style="{ zIndex: 2002 }">
        <div class="el-image-viewer__mask" @click="closeImageViewer"></div>
        <span class="el-image-viewer__btn el-image-viewer__close" @click="closeImageViewer">
          <i class="el-icon-close"></i>
        </span>
        <div class="el-image-viewer__canvas">
          <div class="floating">
            <img :src="imageUrl" class="el-image-viewer__img" :style="{
              transform: `scale(${scale}) rotate(${rotation}deg)`,
              marginLeft: `${marginLeft}px`,
              marginTop: `${marginTop}px`,
            }" style="max-width:70%; max-height: 70%;  transition: transform 3s ease;">
          </div>
        </div>

        <div class="foot-button">
          <el-button type="warning" @click="closeImageViewer">暂不查看</el-button>
          <el-button type="primary" @click="tumpReceive">立即查看</el-button>
        </div>

        <div class="letter-title">
          <h1 class=" floating-title" style="max-width: 1000px;">好友来信</h1>
        </div>
        <div class="letter-text">
          <div class="target animationClass">您的好友给您写了一封侨批哦,快来看看吧</div>
        </div>
      </div>
    </transition>
    <div class="bottom">
      <!-- <img src="../../assets/imgss/online.webp" alt="" class="online">
        <img src="../../assets/imgss/oversees.webp" alt="" class="oversees">
        <img src="../../assets/imgss/logo.webp" alt="" class="logo1"> -->
      <p class="contact"> &copy;2024 侨缘信使 . 版权所有 . </p>
    </div>
    <el-backtop :bottom="100" :visibility-height="1200">
      <div style="height: 100%; width: 100%;text-align: center;line-height: 20px;">
        <img style="width: 40px; height: 40px;" src="../../assets/logo/backtop.png" alt="">
        <span style="font-size: 14px; color: #A52328;">回顶部</span>
      </div>
    </el-backtop>
  </div>
</template>

<script>
import useUserStore from '@/store/modules/user'
import { getMyNotReadLetter } from '@/api/letter'
export default {
  name: 'IndexIntroduce',
  data() {
    return {
      histories: [
        { id: 1, math: '01', title: '起步阶段', description: '19世纪中叶至19世纪末，随着中国东南沿海地区居民开始大规模移民海外，侨批开始出现，主要用于华侨与家人之间的通信和汇款。', img: require('@/assets/imgss/history1.webp'), active: false },
        { id: 2, math: '02', title: '发展阶段', description: '20世纪初至20世纪30年代，随着华侨数量的增加和海外华人社区的形成，侨批业务迅速发展，出现了专门从事侨批传递的“水客”和侨批局。', img: require('@/assets/imgss/history2.webp'), active: false },
        { id: 3, math: '03', title: '繁荣阶段', description: '20世纪30年代至20世纪40年代中期，侨批业务达到顶峰，成为华侨与家乡之间最主要的联系方式之一，侨批局遍布东南亚各地。', img: require('@/assets/imgss/history3.webp'), active: false },
        { id: 4, math: '04', title: '衰退阶段', description: '20世纪40年代中期至20世纪50年代，由于第二次世界大战的影响，国际邮件和汇款业务受阻侨批业务受到严重打击。', img: require('@/assets/imgss/history4.webp'), active: false },
        { id: 5, math: '05', title: '转型阶段', description: '20世纪50年代至20世纪70年代，随着战后重建和通信技术的发展，电话、电报等现代通信方式逐渐普及，侨批的传统功能开始减弱。', img: require('@/assets/imgss/history5.webp'), active: false },
        { id: 6, math: '06', title: '遗产化阶段', description: '20世纪未至今，随着现代通信技术的高度发展，侨批作为通信手段的功能已经消失，但它作为历史文化遗产的价值逐渐被认识和重视。', img: require('@/assets/imgss/history6.webp'), active: false },
        { id: 7, math: '07', title: '保护阶段', description: '21世纪初至今，侨批作为珍贵的历史文献，受到了学术界的广泛关注，多国政府和国际组织开始采取措施保护和研究侨批。', img: require('@/assets/imgss/history7.webp'), active: false }
      ],
      imageUrl: 'http://110.41.58.26:9000/qiaopi/qiaopi-images/letter/7780d647-ac9a-4ac0-9342-ad35e9e709d1.webp',
      scale: 0.1,
      rotation: 0,
      marginLeft: 0,
      marginTop: 0,
      isVisible: false,
      showFootButton: false,
      myNotReadLetter: {},
      scrollPosition: 0, // 滚动条的垂直位置
      translateYValue: 0, // 图片的下移距离，初始为0
      scaleValue: 1, // 图片的缩放比例，初始为1（原始大小）
      showContent: false, // 控制内容显示的数据属性
      translateXValue: 0
    }
  },
  computed: {
    blackboatStyle() {
      return {
        transform: `translateY(${this.translateYValue}px) translateX(${this.translateXValue}px) scale(${this.scaleValue})`
      }
    }
  },
  methods: {
    tumpReceive() {
      this.closeImageViewer()
      setTimeout(() => {
        const params = {
          showFirst: true,
          myNotReadLetter: this.myNotReadLetter
        }
        this.$router.push({
          name: 'ReceiveLetter',
          params: params
        })
      }, 1001)
    },
    getMyNotReadLetter() {
      const userStore = useUserStore()
      if (userStore.token) {
        getMyNotReadLetter().then(res => {
          if (res.data) {
            setTimeout(() => {
              this.myNotReadLetter = res.data
              this.openImageViewer(this.myNotReadLetter)
            }, 3000)
          }
        })
      }
    },
    toggleActive(history) {
      this.histories.forEach(function (h) {
        h.active = false
      })
      history.active = true
    },
    toggleInactive(history) {
      history.active = false
    },
    openImageViewer(item) {
      this.isVisible = true

      this.imageUrl = item.coverLink

      setTimeout(() => {
        this.scale = 0.9
      }, 300)
      setTimeout(() => {
        this.showFootButton = true
      }, 3000)
    },
    closeImageViewer() {
      this.scale = 0
      setTimeout(() => {
        this.isVisible = false
      }, 301) // 假设动画持续时间为300ms
      this.marginLeft = 0
      this.marginTop = 0
    },
    handleScroll() {
      const scrollPosition = window.pageYOffset
      const blackboatElement = document.querySelector('.blackboat')
      if (blackboatElement) {
        // 假设当滚动条滑动到距离页面顶部1000px时开始效果
        const startEffectAt = 1000
        // 假设当滚动条滑动到距离页面顶部1500px时效果结束
        const endEffectAt = 1500
        if (scrollPosition > startEffectAt && scrollPosition < endEffectAt) {
          // 向下滚动
          this.translateYValue = Math.min(scrollPosition, 200) // 控制向下移动的最大值
          this.scaleValue = Math.max(0.5, 1 - scrollPosition * 0.0005) // 控制缩放的最小值，假设最小缩放到0.5
          this.translateXValue = -400
        } else if (scrollPosition > endEffectAt) {
          this.translateYValue = 200 // 恢复到原来位置
          this.scaleValue = 0.5 // 恢复原始大小
          this.translateXValue = -400
        } else {
          this.translateYValue = 0 // 恢复到原来位置
          this.scaleValue = 1 // 恢复原始大小
          this.translateXValue = 0
        }
        // this.scrollPosition = scrollPosition // 更新滚动位置
      }
    },
    toggleContent() {
      this.showContent = !this.showContent // 切换内容显示状态
    }
  },
  mounted() {
    window.addEventListener('scroll', this.handleScroll)
    this.getMyNotReadLetter()
    const observer = new IntersectionObserver((entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) {
          entry.target.style.opacity = 1
          observer.unobserve(entry.target)
        }
      })
    })

    document.querySelector('.right .second') && observer.observe(document.querySelector('.right .second'))
    document.querySelector('.right .third') && observer.observe(document.querySelector('.right .third'))
  },
  beforeDestroy() {
    window.removeEventListener('scroll', this.handleScroll)
  }
}

</script>
<style scoped>
* {
  margin: 0;
  /* padding: 0; */
  list-style: none;
  box-sizing: border-box;
}

::v-deep .el-backtop {
  background-color: transparent;
  width: 50px;
  height: 75px;
  padding-top: 5px;
  border-radius: 10px;
  transition: all 0.3s;
}

.target {
  /* background: #666; */
  /* width: 108px; */
  height: 60px;
  border-radius: 60px;
  color: aliceblue;
  font-weight: bolder;
  line-height: 60px;
  text-align: center;
}

.animationClass {
  animation: 5s ease 0s infinite backwards bounce;
}

/* 主要是通过控制translateY的值来操作y轴的距离实现弹跳效果 */
@keyframes bounce {
  0% {
    transform: translateY(-64px);
    animation-timing-function: ease-in;
    opacity: 1;
  }

  24% {
    opacity: 1;
  }

  40% {
    transform: translateY(-32px);
    animation-timing-function: ease-in;
  }

  62% {
    transform: translateY(-16px);
    animation-timing-function: ease-in;
  }

  82% {
    transform: translateY(-8px);
    animation-timing-function: ease-in;
  }

  92% {
    transform: translateY(-4px);
    animation-timing-function: ease-in;
  }

  25%,
  55%,
  75%,
  90% {
    transform: translateY(0);
    animation-timing-function: ease-out;
  }

  100% {
    transform: translateY(0);
    animation-timing-function: ease-out;
    opacity: 1;
  }
}

.foot-button {
  position: fixed;
  right: 10%;
  bottom: 8%;
  transform: translate(50%, 50%);
  width: 300px;
  height: 100px;
  /* background: rgba(0, 0, 0, .4); */
  border-radius: 10px;
  display: flex;
  justify-content: center;
  align-items: center;
  opacity: 0;
  animation: fadeIn 3s ease-in-out forwards;
}

.letter-text {
  position: fixed;
  left: 20%;
  top: 50%;
  transform: translate(-50%, -50%);
  width: 300px;
  height: 100px;
  font-family: '华文中宋', sans-serif;
  font-size: 2em;
  border-radius: 10px;
  display: flex;
  justify-content: center;
  align-items: center;
  opacity: 0;
  animation: fadeIn 3s ease-in-out forwards;
}

.letter-title {
  position: fixed;
  right: 20%;
  top: 25%;
  font-family: '华文中宋', sans-serif;
  font-size: 2em;
  transform: translate(50%, -50%);
  width: 300px;
  height: 100px;
  /* background: rgba(0, 0, 0, .4); */
  color: #FFA000;
  border-radius: 10px;
  display: flex;
  justify-content: center;
  align-items: center;
  opacity: 0;
  animation: fadeIn 3s ease-in-out forwards;
}

@keyframes fadeIn {
  0% {
    opacity: 0;
  }

  100% {
    opacity: 1;
  }
}

.floating {
  max-width: 50%;
  max-height: 600px;
  animation: float 4s ease-in-out infinite;
}

.floating-title {
  max-width: 50%;
  max-height: 600px;
  animation: float 5s ease-in-out infinite;
}

@keyframes float {

  0%,
  100% {
    transform: translateY(0);
  }

  50% {
    transform: translateY(-20px);
  }
}

.contain {
  width: 100%;
  height: 100%;
  background-size: 100% 100%;
  background-position: center center;
  background-size: cover;
}

.banner {
  position: relative;
  width: 1440px;
  height: 100%;
  margin: 0 auto;
}

/* 侨批简介 */
.introduce h1 {
  position: relative;
  width: 400px;
  height: 100px;
  margin: 0 auto;
}

.introduce h1 img {
  width: 100%;
  height: 100%;
}

.introduce p {
  position: absolute;
  top: 28px;
  left: 125px;
  color: #924721;
}

section {
  position: relative;
  width: 100%;
  height: 800px;
  display: flex;
  justify-content: center;
  align-items: center;
}

section .logo {
  position: absolute;
  top: 220px;
  left: 0;
  width: 100%;
  z-index: 10;
}

/* #sticker1 {
  position: absolute;
  top: 220px;
  left: 0;
  width: 50%;
  animation: moveright 1.5s forwards;
}

@keyframes moveright {
  0% {
    left: -800px;
    top: 500px;
  }

  100% {
    top: 220px;
    left: 0;
  }
}

#sticker2 {
  position: absolute;
  top: 220px;
  right: 0;
  width: 50%;
  animation: moveleft 2s forwards;
}

@keyframes moveleft {
  0% {
    right: -800px;
    top: 500px;
  }

  100% {
    top: 220px;
    right: 0;
  }
} */

section #text {
  position: absolute;
  width: 750px;
  height: 320px;
  text-align: center;
  top: -1000px;
  animation: move 1.5s ease-in-out forwards;
}

@keyframes move {
  0% {
    width: 300px;
    top: -1000px;
  }

  100% {
    width: 750px;
    top: -100px;
  }

}

.left {
  position: relative;
  float: left;
  width: 500px;
  height: 800px;
  padding-left: 100px;
}

.left .title {
  width: 360px;
  height: 110px;
  margin-top: 150px;
}

/* 定义swing动画 */
@keyframes swing {
  0% {
    transform: rotate(0deg);
  }

  25% {
    transform: rotate(10deg);
  }

  50% {
    transform: rotate(0deg);
  }

  75% {
    transform: rotate(-10deg);
  }

  100% {
    transform: rotate(0deg);
  }
}

/* 应用动画到.bell类 */
.bell {
  position: absolute;
  top: -10px;
  left: 0;
  z-index: 0;
  width: 200px;
  height: 250px;
  animation: swing 3s linear infinite;
  cursor: pointer;
}

.left .introduce {
  width: 450px;
  height: 200px;
  margin: 0px auto;
  margin-top: 50px;
  background-color: white;
}

.right {
  position: relative;
  float: left;
  width: 900px;
  height: 900px;
  margin-top: 20px;
  /* background-image: url(../../assets/imgss/right.webp); */
}

.right .first {
  position: absolute;
  top: 350px;
  left: 230px;
  right: 120px;
  text-indent: 2em;
  z-index: 1;
  color: #5C555D;
  font-size: 35px;
  text-align: left;
  font-family: "华文隶书", "楷体", serif;
}

.right .red {
  color: #C00000;
}

.right .yellow {
  color: #B58B02;
}

.right .blue {
  color: #2E54A1;
}

.right .second {
  position: absolute;
  top: 950px;
  left: -150px;
  right: 300px;
  text-indent: 2em;
  z-index: 1;
  font-family: "华文隶书", "楷体", serif;
  color: #5C555D;
  font-size: 35px;
  text-align: left;
}

.right .third {
  position: absolute;
  top: 1200px;
  left: 100px;
  right: 200px;
  text-indent: 2em;
  z-index: 1;
  font-family: "华文隶书", "楷体", serif;
  color: #5C555D;
  font-size: 35px;
  text-align: left;
}

/* .right .third {
    position: absolute;
    top: 620px;
    left: 180px;
    right: 300px;
    z-index: 1;
    color: #000;
    font-size: 27px;
    font-family: "KaiTi", "楷体", serif;
    text-align: left;
    text-indent: 2em;
  } */
.imgLoader {
  position: fixed;
  -webkit-animation: preLoad 1s steps(1);
  animation: preLoad 1s steps(1);
  width: 1px;
  height: 1px;
}

.credit {
  position: absolute;
  top: 100%;
  left: 0px;
  font-size: 0.9em;
  text-align: left;
}

.book {
  position: relative;
  perspective: 630px;
  perspective-origin: center 50px;
  transform: scale(1.2);
  filter: drop-shadow(0px 10px 5px rgba(0, 0, 0, 0.25));
}

.page {
  width: 210px;
  height: 300px;
  background-color: #bbb;
  position: absolute;
  top: 120px;
  right: 50%;
  transform-origin: 100% 100%;
  border: solid #555 2px;
  background-size: 420px 300px;
  background-position: center;
  transform-style: preserve-3d;
  transform: rotateX(60deg) rotateY(3deg);
}

.page:nth-child(1) {
  transform: rotateX(20deg) rotateY(3deg);
}

.page:nth-child(2) {
  transform: rotateX(20deg) rotateY(4.5deg);
}

.page:nth-child(3) {
  transform: rotateX(20deg) rotateY(6deg);
  -webkit-animation: nextPage 25s infinite -24s steps(1);
  animation: nextPage 25s infinite -24s steps(1);
  background-size: 420px 300px;
  background-position: -2px -2px;
}

.page:nth-child(4) {
  transform: rotateX(20deg) rotateY(177deg);
}

.page:nth-child(5) {
  transform: rotateX(20deg) rotateY(175.5deg);
}

.page:nth-child(6) {
  transform: rotateX(20deg) rotateY(174deg);
  overflow: hidden;
}

.page:nth-child(6)::after {
  content: "";
  width: 210px;
  height: 300px;
  position: absolute;
  top: 0px;
  right: 0%;
  transform-origin: center;
  transform: rotateY(180deg);
  -webkit-animation: nextPage 25s -20s infinite steps(1);
  animation: nextPage 25s -20s infinite steps(1);
  background-size: 420px 300px;
  background-position: 100% -2px;
}

.gap {
  width: 10px;
  height: 300px;
  background: none;
  transform: rotateX(60deg);
  transform-origin: bottom;
  position: absolute;
  top: 120px;
  left: calc(50% - 5px);
}

.gap::after {
  content: "";
  position: absolute;
  bottom: 0;
  left: 50%;
  transform: translate(-50%, 50%);
  background-color: #555;
  width: 10px;
  height: 5px;
  border-radius: 50%;
}

.pages,
.flips {
  transform-style: preserve-3d;
}

.flip {
  width: 32px;
  height: 300px;
  position: absolute;
  top: 120px;
  transform-origin: 100% 100%;
  right: 100%;
  border: solid #555;
  border-width: 2px 0px;
  perspective: 4200px;
  perspective-origin: center;
  transform-style: preserve-3d;
  background-size: 420px 300px;
}

.flip::after {
  content: "";
  position: absolute;
  top: 0px;
  right: 0%;
  width: 100%;
  height: 100%;
  transform-origin: center;
  background-size: 420px 300px;
}

.flip.flip1 {
  right: 50%;
  -webkit-animation: flip1 5s infinite ease-in-out;
  animation: flip1 5s infinite ease-in-out;
  border-width: 2px 2px 2px 0;
}

.flip.flip1::after {
  -webkit-animation: nextFlip1 20s -20s infinite steps(1);
  animation: nextFlip1 20s -20s infinite steps(1);
}

.flip:not(.flip1) {
  right: calc(100% - 2px);
  top: -2px;
  transform-origin: right;
  -webkit-animation: flip2 5s ease-in-out infinite;
  animation: flip2 5s ease-in-out infinite;
}

.flip.flip2::after {
  -webkit-animation: nextFlip2 20s -20s infinite steps(1);
  animation: nextFlip2 20s -20s infinite steps(1);
}

.flip.flip3::after {
  -webkit-animation: nextFlip3 20s -20s infinite steps(1);
  animation: nextFlip3 20s -20s infinite steps(1);
}

.flip.flip4::after {
  -webkit-animation: nextFlip4 20s -20s infinite steps(1);
  animation: nextFlip4 20s -20s infinite steps(1);
}

.flip.flip5::after {
  -webkit-animation: nextFlip5 20s -20s infinite steps(1);
  animation: nextFlip5 20s -20s infinite steps(1);
}

.flip.flip6::after {
  -webkit-animation: nextFlip6 20s -20s infinite steps(1);
  animation: nextFlip6 20s -20s infinite steps(1);
}

.flip.flip7::after {
  -webkit-animation: nextFlip7 20s -20s infinite steps(1);
  animation: nextFlip7 20s -20s infinite steps(1);
}

.flip.flip7 {
  width: 30px;
  border-width: 2px 0px 2px 2px;
}

.flip.flip7::after {
  -webkit-animation: nextFlip7 20s -20s infinite steps(1);
  animation: nextFlip7 20s -20s infinite steps(1);
}

@keyframes flip1 {

  0%,
  20% {
    transform: rotateX(20deg) rotateY(6deg);
  }

  80%,
  100% {
    transform: rotateX(20deg) rotateY(174deg);
  }
}

@keyframes flip2 {

  0%,
  20% {
    transform: rotateY(0deg) translateY(0px);
  }

  50% {
    transform: rotateY(-15deg) translateY(0px);
  }
}

@keyframes preLoad {
  0% {
    background-image: url("../../assets/imgss/page1.webp");
  }

  25% {
    background-image: url("../../assets/imgss/page2.webp");
  }

  50% {
    background-image: url("../../assets/imgss/page3.webp");
  }

  75% {
    background-image: url("../../assets/imgss/page4.webp");
  }

  100% {
    display: none;
  }
}

@keyframes nextPage {
  0% {
    background-image: url("../../assets/imgss/page1.webp");
  }

  20% {
    background-image: url("../../assets/imgss/page2.webp");
  }

  40% {
    background-image: url("../../assets/imgss/page3.webp");
  }

  60% {
    background-image: url("../../assets/imgss/page4.webp");
  }

  80% {
    background-image: url("../../assets/imgss/page5.webp");
  }
}

@keyframes nextFlip1 {
  0% {
    background-image: url("../../assets/imgss/page1.webp");
    background-position: -178px -2px;
    transform: rotateY(0deg);
  }

  25% {
    background-image: url("../../assets/imgss/page2.webp");
    background-position: -210px -2px;
    transform: rotateY(180deg);
  }

  50% {
    background-image: url("../../assets/imgss/page3.webp");
    background-position: -178px -2px;
    transform: rotateY(0deg);
  }

  75% {
    background-image: url("../../assets/imgss/page4.webp");
    background-position: -210px -2px;
    transform: rotateY(180deg);
  }

  100% {
    background-image: url("../../assets/imgss/page1.webp");
    background-position: -178px -2px;
    transform: rotateY(0deg);
  }
}

@keyframes nextFlip2 {
  0% {
    background-image: url("../../assets/imgss/page1.webp");
    background-position: -148px -2px;
    transform: rotateY(0deg);
  }

  25% {
    background-image: url("../../assets/imgss/page2.webp");
    background-position: -238px -2px;
    transform: rotateY(180deg);
  }

  50% {
    background-image: url("../../assets/imgss/page3.webp");
    background-position: -148px -2px;
    transform: rotateY(0deg);
  }

  75% {
    background-image: url("../../assets/imgss/page4.webp");
    background-position: -238px -2px;
    transform: rotateY(180deg);
  }

  100% {
    background-image: url("../../assets/imgss/page5.webp");
    background-position: -148px -2px;
    transform: rotateY(0deg);
  }
}

@keyframes nextFlip3 {
  0% {
    background-image: url("../../assets/imgss/page1.webp");
    background-position: -118px -2px;
    transform: rotateY(0deg);
  }

  25% {
    background-image: url("../../assets/imgss/page2.webp");
    background-position: -268px -2px;
    transform: rotateY(180deg);
  }

  50% {
    background-image: url("../../assets/imgss/page3.webp");
    background-position: -118px -2px;
    transform: rotateY(0deg);
  }

  75% {
    background-image: url("../../assets/imgss/page4.webp");
    background-position: -268px -2px;
    transform: rotateY(180deg);
  }

  100% {
    background-image: url("../../assets/imgss/page5.webp");
    background-position: -118px -2px;
    transform: rotateY(0deg);
  }
}

@keyframes nextFlip4 {
  0% {
    background-image: url("../../assets/imgss/page1.webp");
    background-position: -88px -2px;
    transform: rotateY(0deg);
  }

  25% {
    background-image: url("../../assets/imgss/page2.webp");
    background-position: -298px -2px;
    transform: rotateY(180deg);
  }

  50% {
    background-image: url("../../assets/imgss/page3.webp");
    background-position: -88px -2px;
    transform: rotateY(0deg);
  }

  75% {
    background-image: url("../../assets/imgss/page4.webp");
    background-position: -298px -2px;
    transform: rotateY(180deg);
  }

  100% {
    background-image: url("../../assets/imgss/page5.webp");
    background-position: -88px -2px;
    transform: rotateY(0deg);
  }
}

@keyframes nextFlip5 {
  0% {
    background-image: url("../../assets/imgss/page1.webp");
    background-position: -58px -2px;
    transform: rotateY(0deg);
  }

  25% {
    background-image: url("../../assets/imgss/page2.webp");
    background-position: -328px -2px;
    transform: rotateY(180deg);
  }

  50% {
    background-image: url("../../assets/imgss/page3.webp");
    background-position: -58px -2px;
    transform: rotateY(0deg);
  }

  75% {
    background-image: url("../../assets/imgss/page4.webp");
    background-position: -328px -2px;
    transform: rotateY(180deg);
  }

  100% {
    background-image: url("../../assets/imgss/page5.webp");
    background-position: -58px -2px;
    transform: rotateY(0deg);
  }
}

@keyframes nextFlip6 {
  0% {
    background-image: url("../../assets/imgss/page1.webp");
    background-position: -28px -2px;
    transform: rotateY(0deg);
  }

  25% {
    background-image: url("../../assets/imgss/page2.webp");
    background-position: -358px -2px;
    transform: rotateY(180deg);
  }

  50% {
    background-image: url("../../assets/imgss/page3.webp");
    background-position: -28px -2px;
    transform: rotateY(0deg);
  }

  75% {
    background-image: url("../../assets/imgss/page4.webp");
    background-position: -358px -2px;
    transform: rotateY(180deg);
  }

  100% {
    background-image: url("../../assets/imgss/page5.webp");
    background-position: -28px -2px;
    transform: rotateY(0deg);
  }
}

@keyframes nextFlip7 {
  0% {
    background-image: url("../../assets/imgss/page1.webp");
    background-position: -2px -2px;
    transform: rotateY(0deg);
  }

  25% {
    background-image: url("../../assets/imgss/page2.webp");
    background-position: -388px -2px;
    transform: rotateY(180deg);
  }

  50% {
    background-image: url("../../assets/imgss/page3.webp");
    background-position: -2px -2px;
    transform: rotateY(0deg);
  }

  75% {
    background-image: url("../../assets/imgss/page4.webp");
    background-position: -388px -2px;
    transform: rotateY(180deg);
  }

  100% {
    background-image: url("../../assets/imgss/page5.webp");
    background-position: -2px -2px;
    transform: rotateY(0deg);
  }
}

.twitterLink {
  position: fixed;
  bottom: 0.5em;
  right: 0.5em;
}

.twitterLink img {
  width: 2em;
  filter: grayscale(100%);
  transition: filter 0.25s;
}

.twitterLink img:hover {
  filter: grayscale(0%);
}

/* 侨批历史 */
.history {
  width: 100%;
  margin-top: 1650px;
}

.history h1 {
  width: 400px;
  height: 100px;
  margin: 0 auto;
}

.history h1 img {
  width: 100%;
  height: 100%;
}

.history .banner {
  position: relative;
  z-index: 1;
  margin: 0 auto;
  width: 1320px;
  height: 450px;
}

.history .banner h1 {
  position: absolute;
  width: 200px;
  height: 25px;
  top: -20px;
  left: 480px;
  font-size: 28px;
}

.history .banner ul {
  position: relative;
  margin: 0;
  padding: 0;
  list-style: none;
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
}

.history .banner ul li {
  position: relative;
  flex: 0 0 141px;
  height: 380px;
  margin-right: 20px;
  margin-top: 60px;
  border-radius: 10px;
  transition: all .5s;
  cursor: pointer;
  background-repeat: no-repeat;
  background-position: center center;
  background-size: cover;
}

.history .banner ul li .pic {
  position: relative;
  width: 100%;
  height: 100%;
  border-radius: 10px;
  background-repeat: no-repeat;
  background-position: center center;
  background-size: cover;
}

.history .banner ul li::after {
  content: "";
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(360deg, #b67c1095 0%, rgba(108, 63, 255, 0) 100%);
  border-radius: 10px;
  /* 与 li 元素的 border-radius 一致 */
  transition: opacity 0.5s;
  /* 平滑过渡效果 */
  opacity: 1;
  /* 默认不显示 */
}

.history .banner ul li .pic .math {
  position: absolute;
  top: 200px;
  left: 20px;
  color: rgb(255, 255, 255);
  font-weight: bold;
  font-size: 40px;
  z-index: 10;
}

.history .banner ul li .pic .title {
  position: absolute;
  top: 260px;
  left: 20px;
  font-size: 20px;
  color: rgb(255, 255, 255);
  z-index: 10;
}

.history .banner ul li .pic p {
  display: none;
  position: absolute;
  top: 280px;
  left: 0px;
  color: white;
  font-size: 14px;
  width: 100%;
  z-index: 10;
}

.history .banner ul .active {
  flex-basis: 282px;
}

.history .banner ul .active .pic .title {
  top: 215px;
  left: 80px;
  font-size: 20px;
}

.history .banner ul .active .pic p {
  display: block;
}

.function {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  margin-top: 200px;
}

.function h1 {
  position: absolute;
  top: -140px;
  width: 400px;
  height: 100px;
}

.function h1 img {
  width: 100%;
  height: 100%;
}

.function a {
  display: inline-block;
  width: 200px;
  height: 200px;
  border-radius: 50%;
}

.images-container {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-around;
}

.top-row,
.bottom-row {
  display: flex;
  width: 100%;
  justify-content: space-around;
}

.function-img {
  width: 200px;
  /* 初始宽度 */
  height: 200px;
  /* 初始高度 */
  border-radius: 50%;
  object-fit: cover;
  transition: transform 0.3s ease;
  /* 添加一个平滑的变换效果 */
}

.function-img:hover {
  transform: scale(1.2);
  /* 鼠标悬停时放大到原来的1.2倍 */
}

.tip {
  position: absolute;
  bottom: -30px;
  /* 根据需要调整位置 */
  left: 50%;
  transform: translateX(-50%);
  background-color: #333;
  color: #fff;
  padding: 5px 10px;
  border-radius: 4px;
  visibility: hidden;
  /* 默认不可见 */
  opacity: 0;
  transition: visibility 0.3s, opacity 0.3s;
  /* 平滑过渡效果 */
}

.function-img:hover p {
  visibility: visible;
  /* 鼠标悬停时可见 */
  opacity: 1;
}

.bottom {
  position: relative;
  width: 100%;
  height: 70px;
  margin-top: 50px;
  background-image: url(../../assets/imgss/logobgd.webp);
  background-repeat: repeat-x;
  font-size: 0;
  line-height: 60px;
}

.online {
  position: absolute;
  top: 30px;
  left: 120px;
  width: 25%;
}

.oversees {
  position: absolute;
  width: 20%;
  top: 50px;
  left: 480px;
}

.logo1 {
  position: absolute;
  width: 5%;
  top: 35px;
  left: 800px;
}

.contact {
  position: absolute;
  /* top: 40px;
  right: 180px; */
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 18px;
  color: rgba(255, 255, 255, 0.6);
}

.cloud1 {
  position: absolute;
  top: 20px;
  left: 0;
  /* z-index: -1; */
  animation: cloud1move 4s forwards infinite;
}

.cloud2 {
  position: absolute;
  top: 270px;
  left: 430px;
  animation: cloud2move 4s forwards infinite;
}

.cloud3 {
  position: absolute;
  top: 120px;
  right: 150px;
  animation: cloud3move 4s forwards infinite;
}

@keyframes cloud1move {
  0% {
    left: 0px;
    top: 20px;
  }

  25% {
    left: 20px;
    top: 20px;
  }

  50% {
    top: 40px;
    left: 20px;
  }

  75% {
    top: 40px;
    left: 0px;
  }

  100% {
    left: 0px;
    top: 20px;
  }
}

@keyframes cloud2move {
  0% {
    left: 430px;
    top: 270px;
  }

  50% {
    top: 270px;
    left: 450px;
  }

  100% {
    left: 430px;
    top: 270px;
  }
}

@keyframes cloud3move {
  0% {
    right: 150px;
    top: 120px;
  }

  25% {
    right: 170px;
    top: 120px;
  }

  50% {
    right: 150px;
    top: 120px;
  }

  75% {
    right: 150px;
    top: 140px;
  }

  100% {
    right: 150px;
    top: 120px;
  }
}

.paper1 {
  position: absolute;
  top: 30px;
  left: 180px;
  animation: paper1move 4s forwards;
}

.paper2 {
  position: absolute;
  top: 400px;
  left: 300px;
  animation: paper2move 4s linear forwards;
}

.paper3 {
  position: absolute;
  top: 30px;
  left: 500px;
  animation: paper3move 4s forwards;
}

.paper4 {
  position: absolute;
  top: 100px;
  left: 700px;
  animation: paper4move 4s forwards;
}

.paper5 {
  position: absolute;
  top: 100px;
  right: 200px;
  animation: paper5move 4s forwards;
}

.paper6 {
  position: absolute;
  top: 250px;
  right: 270px;
  animation: paper6move 4s forwards;
}

@keyframes paper1move {
  0% {
    left: 180px;
    top: -100px;
  }

  50% {
    left: 200px;
    top: 0px;
  }

  100% {
    top: 30px;
    left: 180px;
  }
}

@keyframes paper2move {
  0% {
    top: -100px;
    left: 300px;
  }

  25% {
    left: 270px;
    top: 100px;
  }

  50% {
    left: 330px;
    top: 200px;
  }

  75% {
    left: 330px;
    top: 300px;
  }

  100% {
    top: 400px;
    left: 300px;
  }
}

@keyframes paper3move {
  0% {
    top: -100px;
    left: 500px;
  }

  50% {
    top: 0px;
    left: 450px;
  }

  100% {
    top: 30px;
    left: 500px;
  }
}

@keyframes paper4move {
  0% {
    top: -100px;
    left: 900px;
  }

  100% {
    top: 100px;
    left: 700px;
  }
}

@keyframes paper5move {
  0% {
    top: 100px;
    right: -100px;
  }

  100% {
    top: 100px;
    right: 200px;
  }
}

@keyframes paper6move {
  0% {
    top: 100px;
    right: -100px;
  }

  100% {
    top: 250px;
    right: 270px;
  }
}

.text1 {
  position: absolute;
  top: 400px;
  right: 300px;
  width: 480px;
  height: 60px;
}

.text2 {
  position: absolute;
  top: 450px;
  right: 50px;
  width: 600px;
  height: 60px;
}

.text3 {
  position: absolute;
  top: 500px;
  right: 270px;
  width: 600px;
  height: 60px;
}

.text4 {
  position: absolute;
  top: 550px;
  right: 30px;
  width: 700px;
  height: 60px;
}

.boat {
  position: absolute;
  bottom: 90px;
  left: -30px;
  width: 780px;
  height: 500px;
  animation: boatmove 6s infinite forwards;
}

@keyframes boatmove {

  0%,
  100% {
    bottom: 90px;
    left: -30px;
  }

  50% {
    bottom: 88px;
    left: -40px;
  }
}

.boatsea {
  position: absolute;
  bottom: -20px;
  /*left: -50px;
  width: 1600px;
  */
  width: 100vw;
  height: 230px;
  animation: seamove 8s infinite forwards;
}

@keyframes seamove {

  0%,
  100% {
    bottom: 0px;
    /* left: -50px; */
  }

  50% {
    bottom: -20px;
    /* left: -50px; */
  }
}

.tree {
  position: absolute;
  left: -30px;
  top: 0px;
  width: 700px;
  height: 400px;
}

.blackboat {
  position: absolute;
  right: 0px;
  bottom: -50px;
  z-index: 100;
  /* width: 1200px;
  height: 800px; */
  transition: transform 1.5s linear;
  /* 平滑过渡效果 */
}

.hidden {
  visibility: hidden;
}

.facai {
  position: absolute;
  top: 1150px;
  left: -500px;
  width: 400px;
}

.right .second,
.right .third {
  opacity: 0;
  transition: opacity 4s ease-in-out;
}
</style>
