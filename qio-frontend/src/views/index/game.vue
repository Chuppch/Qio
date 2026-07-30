<template>
  <section>
    <div class="card">
      <div class="box">
        <div class="img-box">
          <img src="../../assets/imgss/card1.webp" alt="">
        </div>
        <div class="con-box">
          <div>
            <h2>水客智库</h2>
            <p>
              这是一个聚焦侨批知识的答题游戏。游戏中的题目涉及侨批的历史、文化、意义等多方面内容。玩家通过回答问题参与游戏，每答对一题便可获得相应猪仔钱，在积累知识的同时可获得猪仔钱，换取写信的材料哦
            </p>
          </div>
        </div>
      </div>
      <router-link to='/know'><button>开始探索</button>
        <spinner v-if="loading" :width="60" :height="60" :radius="30" :length="8" :line-width="4" color="#ffffff" />
      </router-link>
    </div>
    <div class="card">
      <div class="box">
        <div class="img-box">
          <img src="../../assets/imgss/card2.webp" alt="">
        </div>
        <div class="con-box">
          <div>
            <h2>侨途风云</h2>
            <p>
              在本游戏中，玩家将化身侨胞，沿着历史的脉络，体验充满波折的侨胞生活故事线。每次在游戏中的互动与决策都会影响故事走向。深入感受侨胞在不同境遇下的喜怒哀乐，见证他们的奋斗、坚韧与对家乡的深切思念。参加游戏可获得相应经验值，换取写信的材料哦
            </p>
          </div>
        </div>
      </div>
      <router-link to='/import'><button>开始探索</button></router-link>
    </div>
    <div class="card">
      <div class="box">
        <div class="img-box">
          <img src="../../assets/imgss/card3.webp" alt="">
        </div>
        <div class="con-box">
          <div>
            <h2>侨悦翻享</h2>
            <p>
              侨悦翻享是一款旨在通过趣味互动方式加深用户对侨批文化了解的游戏。通过翻翻乐游戏，用户将有机会探索侨批文化的独特魅力，同时享受游戏带来的乐趣。玩家的目标是通过匹配相同的侨批图片，了解背后的故事和历史信息，同时积累猪仔钱。
            </p>
          </div>
        </div>
      </div>
      <router-link to='/memory'><button>开始探索</button></router-link>
    </div>
  </section>
</template>

<script>
export default {
  name: 'IndexGame',
  data() {
    return {
      loading: false // 控制加载动画的显示
    }
  }
  //   methods: {
  //     startLoading(event) {
  //       this.loading = true
  //       setTimeout(() => {
  //         this.loading = false
  //         this.$router.push({ name: 'Know', params: { id: 1 } })
  //       }, 2000) // 2秒后跳转
  //     }
  //   }
}
</script>

<style scoped>
section {
  display: flex;
  justify-content: center;
  align-items: center;
  flex-wrap: wrap;
  transform-style: preserve-3d;
  position: relative;
  top: 360px;
  left: 50%;
  transform: translate(-50%, -50%);
}

.card {
  position: relative;
  width: 320px;
  height: 320px;
  margin: 20px;
  transform-style: preserve-3d;
  perspective: 1000px;
}

.card .box {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  transform-style: preserve-3d;
  transition: 1s ease;
}

.card .box .img-box {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}

.card .box .img-box img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.card .box .con-box {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: #8e877c;
  display: flex;
  justify-content: center;
  align-items: center;
  backface-visibility: hidden;
  transform-style: preserve-3d;
  transform: rotateY(180deg);
}

.card .box .con-box div {
  color: #fff;
  padding: 20px;
  background: linear-gradient(30deg, #c1bfa0, #868b67);
  transform: translateZ(100px);
  cursor: pointer;
}

.card .box .con-box div h2 {
  font-size: 20px;
  letter-spacing: 2px;
}

.card .box .con-box div p {
  font-size: 15px;
  text-align: justify;
  margin-top: 5px;
}

.card:hover .box {
  transform: rotateY(180deg);
}

button {
  position: absolute;
  width: 100px;
  height: 50px;
  top: 350px;
  left: 110px;
  border-radius: 50px;
  background-color: #7e3f3c;
  border: 0px;
  color: white;
  cursor: pointer;
}
</style>
