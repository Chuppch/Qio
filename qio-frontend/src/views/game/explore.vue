<template>
  <div class="game-container">
    <div class="content-wrapper" :style="{ backgroundImage: `url(${currentScene.backgroundImage})` }">
      <div class="story-content" @click="showAllText">
        <p>
          <span
            v-for="(char, index) in currentScene.text"
            :key="index"
            :class="{ 'visible': index < displayedTextLength }"
          >{{ char }}</span>
        </p>
      </div>
      <div class="choices" v-if="showChoices">
        <button
          v-for="(choice, index) in currentScene.choices"
          :key="index"
          @click="makeChoice(choice.next)"
          class="choice-button"
        >
          {{ choice.text }}
        </button>
      </div>
      <button v-if="!showChoices && !gameOver" @click="nextScene" class="next-button">下一段</button>
      <button v-if="gameOver" @click="restartGame" class="restart-button">重新开始</button>
    </div>
  </div>
</template>

<script>
export default {
  name: 'QiaopiStoryGame',
  data() {
    return {
      currentSceneId: 'start',
      displayedTextLength: 0,
      showChoices: false,
      gameOver: false,
      scenes: {
        start: {
          text: '梦梅的家乡有一口废弃的井，因为曾有人投井自杀而被封闭。他小时候常偷看井里想象着井下的神秘世界，甚至幻想井是通往大海和外部世界的窗户。银溪村靠近韩江，村民们通过韩江出海去往远方，被称为"番客"。梦梅的家族，溪前，曾经显赫一时，但现在已经衰落，与溪后的家族关系疏远。',
          next: 'familyHistory',
          backgroundImage: require('../../assets/imgss/1.webp')
        },
        familyHistory: {
          text: '梦梅的曾祖父郑鸿顺(九爷)在马六甲发了财，与堂兄弟郑鸿利一起在银溪村建了两座大宅，时光里和平安里，分别代表溪前和溪后。溪前家族多才情，但人丁不旺，而溪后家族则务实且人丁兴旺。梦梅的祖父郑集炎和另一位溪前家族成员在同一天被人放火烧死，这使得溪前家族进一步衰落。',
          next: 'fatherStory',
          backgroundImage: require('../../assets/imgss/2.webp')
        },
        fatherStory: {
          text: '梦梅的父亲阿女，名字女性化，是因为家族希望他能长寿。阿女有三好：好茶、好客、好石，但这些爱好都需要金钱维持。梦梅的哥哥梦龙(复生)曾有机会重振家声，他留学日本，加入同盟会，参与革命，但最终在一次刺杀行动中牺牲。',
          next: 'mengmeiStory',
          backgroundImage: require('../../assets/imgss/3.webp')
        },
        mengmeiStory: {
          text: '梦梅自己则被家族的期望和历史的重负所困扰。他有四个儿子，但其中一个被交换给了溪后家族，另一个则被过继给了已故的哥哥。梦梅的生活充满了矛盾和挣扎，他既想继承家族的传统，又想摆脱过去的阴影。他被老祖(他的祖母)寄予厚望，希望他能重振家族，但梦梅自己却感到迷茫和无力。',
          next: 'ThinkingStory',
          backgroundImage: require('../../assets/imgss/4.webp')
        },
        ThinkingStory: {
          text: '他深知，每一步选择都不仅关乎个人的命运，更牵连着整个家族的未来。他不得不在心灵的挣扎中寻求一丝解脱，希望能够找到一条既能尊重家族期望又能实现自我价值的道路。',
          next: 'decision',
          backgroundImage: require('../../assets/imgss/5.webp')
        },
        decision: {
          text: '此时，宋万昌来到了梦梅的家中，提出让他成为义子并接管批局的提议。',
          choices: [
            { text: '接受', next: 'accept' },
            { text: '拒绝', next: 'refuse' }
          ],
          backgroundImage: require('../../assets/imgss/6.webp')
        },
        accept: {
          text: '梦梅无法承担这样的重任，担心这样做可能会辜负宋万昌的信任。此外，梦梅还担心自己的家族历史和个人的能力限制。在接受之后，批局开始出现一些纰漏，信誉渐渐出了问题，梦梅跟随世界的大起大落，过完辛苦无聊的一生。',
          choices: [
            { text: '重新开始', next: 'import' }
          ],
          backgroundImage: require('../../assets/imgss/7.webp')
        },
        refuse: {
          text: '在1916年秋。一个月明星稀、狗吠声声的后半夜，大厝林立，小屋歪斜，半是豪华半是破败的村庄留给了老人、孥仔和姿娘们，梦梅在这个时辰悄悄离开的乘坐洋船前往暹罗（今泰国）。',
          next: 'things',
          backgroundImage: require('../../assets/imgss/8.webp')
        },
        things: {
          text: '梦梅和其他人一样，身上都少不了几样东西：预示平安顺利的一枚顺治铜钱，代表故土的一把井心泥，够吃半个月甚至更久的一袋子甜粿，一条用途广泛的水布，等等。',
          next: 'transportation',
          backgroundImage: require('../../assets/imgss/9.webp')
        },
        transportation: {
          text: '前往东南亚的途径有：靠季风往还的老式木船，便宜但危险，许多过番人在海中消亡。还有需要船票的洋船。洋船用大功率的柴油机驱动，有指南针指示方向，更快也更安全。',
          choices: [
            { text: '老式木船', next: 'moodboat' },
            { text: '洋船', next: 'technologyboat' }
          ],
          backgroundImage: require('../../assets/imgss/10.webp')
        },
        moodboat: {
          text: '梦梅乘坐的老式木船刚好遇上台风，船翻人亡。',
          choices: [
            { text: '重新开始', next: 'import' }
          ],
          backgroundImage: require('../../assets/imgss/11.webp')
        },
        technologyboat: {
          text: '前往码头，发现停在港湾的大船，不是洋船，而是老式的竞争力已经大大降低的红头船，三桅五层的大型木船，虽然土气，但气势雄伟。船头的红色大显斑驳，红色部分变得很暗很暗了。吃水线和水面之间的船体上寄生着很多絮状、藤状和壳状生物，有些在阳光下痛苦蠕动。',
          next: 'reading1',
          backgroundImage: require('../../assets/imgss/12.webp')
        },
        reading1: {
          text: '梦梅隐约听见有人在诵读出海祭文：灵效瀛海 每著神功',
          next: 'reading2',
          backgroundImage: require('../../assets/imgss/13.webp')
        },
        reading2: {
          text: '谋生异域 取道重溟',
          next: 'reading3',
          backgroundImage: require('../../assets/imgss/14.webp')
        },
        reading3: {
          text: '妈祖保佑 安涉利行......',
          next: 'ending',
          backgroundImage: require('../../assets/imgss/15.webp')
        }
      }
    }
  },
  computed: {
    currentScene() {
      return this.scenes[this.currentSceneId]
    }
  },
  methods: {
    makeChoice(nextSceneId) {
      if (nextSceneId === 'import') {
        // 直接跳转到'/import'页面
        this.$router.push('/import')
      } else {
        // 否则，更新当前场景并显示文本
        this.currentSceneId = nextSceneId
        this.displayedTextLength = 0
        this.showChoices = false
        this.typeText()
      }
    },
    restartGame() {
    // 如果当前场景是'import'，则跳转到'/import'页面
      if (this.currentSceneId === 'import') {
        this.$router.push('/import')
      } else {
      // 否则，重置游戏状态
        this.currentSceneId = 'start'
        this.displayedTextLength = 0
        this.showChoices = false
        this.gameOver = false
      }
    },
    typeText() {
      const text = this.currentScene.text
      const timer = setInterval(() => {
        if (this.displayedTextLength < text.length) {
          this.displayedTextLength++
        } else {
          clearInterval(timer)
          if (this.currentScene.choices && this.currentScene.choices.length > 0) {
            this.showChoices = true
          }
        }
      }, 50)
    },
    nextScene() {
      if (this.currentScene.next) {
        this.makeChoice(this.currentScene.next)
      } else {
        this.gameOver = false
      }
    },
    showAllText() {
      this.displayedTextLength = this.currentScene.text.length
    }
  },
  mounted() {
    this.typeText()
  }
}
</script>

<style scoped>
@font-face {
    font-family: 'MyFont'; /* 这里定义一个字体名称 */
    src: url('../../assets/fonts/字魂逍遥行书.ttf') format('truetype'); /* 指定字体文件的路径和格式 */
}
.game-container {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  font-family: '楷体', 'STSong', serif;
}

.content-wrapper {
  position: relative;
  background-color: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(10px);
  padding: 20px;
  border-radius: 10px;
  width: 90%;
  height: 685px;
  margin: 20px;
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
   animation: fadeInOut 3s ease-in-out forwards;
}

.story-content {
  position: absolute;
  left: 50%;
  bottom: 20px;
  transform: translateX(-50%);
  width: 70%;
  background-color: rgba(255, 255, 255, 0.8);
  border: 1px solid #8B4513;
  border-radius: 5px;
  padding: 25px;
  margin-bottom: 20px;
  min-height: 100px;
}
.story-content p {
  text-align: left;
  font-size: 25px;
  font-weight: bold;
}

.story-content span {
  opacity: 0;
  transition: opacity 0.1s ease-in-out;
}

.story-content span.visible {
  opacity: 1;
}

.choices {
  position: absolute;
  width: 30%;
  top: 50%;
  left: 50%;
  transform: translate(-50%,-50%);
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.choice-button, .restart-button {
  background-color: #8B4513;
  color: white;
  border: none;
  padding: 10px 15px;
  border-radius: 5px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.restart-button {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%,-50%);
}

.choice-button:hover, .restart-button:hover {
  background-color: #A0522D;
}

.next-button {
  position: absolute;
  bottom: 20px;
  right: 280px;
  background-color: #8B4513;
  color: white;
  border: none;
  padding: 10px 15px;
  border-radius: 5px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.next-button:hover {
  background-color: #A0522D;
}
@keyframes fadeInOut {
  0%{
    opacity: 0;
  }
  100% {
    opacity: 1;
  }
}

@media (max-width: 600px) {
  .content-wrapper {
    width: 95%;
    padding: 15px;
  }
}

@media (max-height: 600px) {
  .game-container {
    align-items: flex-start;
  }
}
</style>
