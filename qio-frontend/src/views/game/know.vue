<template>
  <div class="banner">
    <img src="../../assets/imgs/answer.png" alt="" class="bgd" v-show="showImage">
    <p v-show="showWord">{{ paragraphs[currentParagraph] }}</p>
    <a href="#" class="nexttick" @click.prevent="nextTick" v-show="showWord">继续</a>
    <el-carousel :interval="4000" type="card" height="400px" v-show="showCarousel" :autoplay="false">
      <el-carousel-item v-for="(status, index) in lockStatus" :key="index" :class="'carousel-item-' + index"
        style="height: 400px; width: 300px; margin-left: 130px;">
        <div v-if="lockStatus[index - 1] === 1 || index === 0" @click="openpage(index)"
          style="height: 100%; width: 100%; cursor: pointer; display: flex; align-items: center; justify-content: center;">
          <h3 class="medium">{{ index + 1 }}</h3>
        </div>
        <div v-if="lockStatus[index - 1] === 0" class="locked">
          <p class="locked-text">您未解锁</p>
        </div>
      </el-carousel-item>
    </el-carousel>
  </div>
</template>

<script>
import { Message } from 'element-ui'
// import { userLoginPage, startAnswer } from '@/api/know'
import { userLoginPage, allAnswerToFront } from '@/api/know'
import CryptoJS from 'crypto-js'
export default {
  name: 'GameKnow',
  data() {
    return {
      paragraphs: [
        '在20世纪初，中国东南沿海的许多家庭依靠海外华侨寄回的侨批来维持生计。侨批不仅承载着金钱，更承载着远方亲人的思念和期望。在这个信息传递不易的时代，水客成为了连接海内外的重要纽带。',
        '你出生在一个有着悠久水客传统的家族。你的祖父和父亲都曾是水客，他们的故事和经历让你对这份工作充满了敬意和向往。今天，你将踏上你的第一次水客之旅，从马来西亚携带一批侨批返回广东潮汕。',
        '在出发前，你的祖父把你叫到一边，递给你一个装满侨批的旧皮箱，并对你说：‘这些侨批比黄金还要珍贵，它们是家人之间的桥梁。你必须确保每一份侨批都能安全送达。',
        '为了帮助你更好地完成这次任务，我会问你一些关于侨批文化的问题，每答对一个问题，你将获得一个猪仔钱积分这些积分将作为你旅程中的储备。'
      ],
      currentParagraph: 0,
      showImage: true,
      showWord: true,
      showCarousel: false,
      lockStatus: []
    }
  },
  methods: {
    nextTick() {
      if (this.currentParagraph < this.paragraphs.length - 1) {
        this.currentParagraph += 1
      } else {
        this.currentParagraph = 0
        this.showImage = false
        this.showWord = false
        this.showCarousel = true
        this.userLoginPage()
      }
    },
    async userLoginPage() {
      try {
        const res = await userLoginPage()
        if (res.code === 200) {
          this.lockStatus = res.data
          this.showCarousel = true
        } else {
          Message.error(res.msg)
        }
      } catch (error) {
        // console.error('获取题库信息失败:', error)
        Message.error('获取题库信息失败')
      }
    },
    openpage(index) {
      if (this.lockStatus[index - 1] === 1 || index === 0) {
        this.allAnswerToFront(index + 1)
      } else if (this.lockStatus[index - 1] === 0) {
        Message.warning('您还未解锁此题库')
      }
    },
    async allAnswerToFront(setId) {
      await allAnswerToFront(setId).then(res => {
        const secretKey = 'qiaopiqiaopiqiaopiqiaopiqiaopinb'

        // 示例用法
        const encryptedQuestions = res.data // 替换成你的加密字符串

        try {
          const content = this.decryptQuestions(encryptedQuestions, secretKey)
          // 这里拿去用就可以了,记得删掉log
          this.$router.push({
            path: '/question',
            query: {
              questions: content
            }
          })
          // console.log('Decrypted Questions:', content)
        } catch (error) {
          // console.error('Decryption failed:', error)
        }
      })
    },
    // 解密函数
    decrypt(encryptedData, secretKey) {
      // 解码Base64字符串
      const bytes = CryptoJS.enc.Base64.parse(encryptedData)

      // 创建密钥
      const key = CryptoJS.enc.Utf8.parse(secretKey)

      // 解密
      const decrypted = CryptoJS.AES.decrypt({ ciphertext: bytes }, key, {
        mode: CryptoJS.mode.ECB,
        padding: CryptoJS.pad.Pkcs7
      })

      // 返回解密后的字符串
      return decrypted.toString(CryptoJS.enc.Utf8)
    },

    // 解密列表中的每个字符串
    decryptList(list, secretKey) {
      return list.map(item => this.decrypt(item, secretKey))
    },

    // 解密并解析JSON字符串到对象数组
    decryptQuestions(encryptedQuestions, secretKey) {
      const decryptedJson = this.decrypt(encryptedQuestions, secretKey)
      return JSON.parse(decryptedJson)
    }
  }
}

</script>

<style scoped>
.banner {
  position: relative;
  width: 1200px;
  height: 650px;
  margin-top: 55px;
  background: url(../../assets/imgs/sea.jpg) no-repeat;
}

.bgd {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 80%;
}

p {
  position: absolute;
  top: 200px;
  left: 180px;
  right: 180px;
  font-size: 25px;
}

.nexttick {
  position: absolute;
  bottom: 150px;
  right: 300px;
  text-decoration: none;
  color: #331316;
  cursor: pointer;
  z-index: 10;
}

.el-carousel__item {
  /* background-image: url(../../assets/imgss/war.webp); */
  background-position: center center;
  background-size: cover;
}

.el-carousel__item:nth-child(2n) {
  background-color: #99a9bf;
}

.el-carousel__item:nth-child(2n+1) {
  background-color: #d3dce6;
}

.el-carousel {
  margin-top: 100px;
}

.locked {
  position: relative;
  cursor: not-allowed;
}

.locked-text {
  position: absolute;
  top: 50px;
  left: 50%;
  transform: translateX(-50%);
  color: rgb(0, 0, 0);
}

.carousel-item-0 {
  background-image: url(../../assets/imgss/pic1.webp);
  background-size: cover;
  background-position: center center;
}

.carousel-item-1 {
  background-image: url(../../assets/imgss/pic2.webp);
  background-size: cover;
  background-position: center center;
}

.carousel-item-2 {
  background-image: url(../../assets/imgss/pic3.webp);
  background-size: cover;
  background-position: center center;
}
</style>
