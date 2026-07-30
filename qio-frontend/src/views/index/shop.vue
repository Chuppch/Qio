<template>
  <div class="banner">
    <img src="../../assets/imgss/storebgd.webp" alt="">
    <div class="tag1" :class="{'tag-active': activeTag === 'paper'}" @click="showDetails('paper')"><img
        src="../../assets/imgss/tag1.webp" alt=""><span>信纸</span></div>
    <div class="tag2" :class="{'tag-active': activeTag === 'ink'}" @click="showDetails('ink')"><img
        src="../../assets/imgss/tag2.webp" alt=""><span>墨水</span></div>
    <div class="tag3" :class="{'tag-active': activeTag === 'font'}" @click="showDetails('font')"><img
        src="../../assets/imgss/tag3.webp" alt=""><span>字体</span></div>
    <div class="tag4" :class="{'tag-active': activeTag === 'card'}" @click="showDetails('card')"><img
        src="../../assets/imgss/tag4.webp" alt=""><span>功能卡</span></div>
    <div class="details" v-if="showDetailsArea">
      <!-- 根据当前选中的类别显示不同的内容 -->
      <el-collapse-transition>
        <div v-if="selectedCategory === 'paper'">
          <h2>信纸商城</h2>
          <div v-for="paper in paperList" :key="paper.id" class="paper-item">
            <div class="paperGoods">
              <img :src="paper.previewImage" alt="" style="width:200px;height:140px;">
              <h3>{{ paper.name }}</h3>
              <p><img src="../../assets/imgss/pigmoney.webp" alt=""
                  style="width:30px;height:20px;vertical-align: middle;"> {{ paper.price }}</p>
              <button v-if="paper.own" @click="buyPaper(paper.id)" disabled="true">已拥有</button>
              <button v-if="!paper.own" @click="buyPaper(paper.id)">兑换</button>
            </div>
          </div>
        </div>
      </el-collapse-transition>
      <el-collapse-transition>
        <div v-if="selectedCategory === 'ink'">
          <h2>墨水商城</h2>
          <div v-for="color in colorList" :key="color.id" class="paper-item">
            <div class="colorGoods">
              <img :src="color.previewImage" alt="" style="width:200px;height:140px;">
              <h3>{{ color.description }}</h3>
              <p><img src="../../assets/imgss/pigmoney.webp" alt=""
                  style="width:30px;height:20px;vertical-align: middle;"> {{ color.price }}</p>
              <button v-if="color.own" @click="buyFontColor(color.id)" disabled="true">已拥有</button>
              <button v-if="!color.own" @click="buyFontColor(color.id)">兑换</button>
            </div>
          </div>
        </div>
      </el-collapse-transition>
      <el-collapse-transition>
        <div v-if="selectedCategory === 'font'">
          <h2>字体商城</h2>
          <div v-for="font in fontList" :key="font.id" class="paper-item">
            <div class="fontsGoods">
              <img :src="font.previewImage" alt="" style="width:100px;height:70px;">
              <h3>{{ font.name }}</h3>
              <p><img src="../../assets/imgss/pigmoney.webp" alt=""
                  style="width:30px;height:20px;vertical-align: middle;"> {{ font.price }}</p>
              <button v-if="font.own" @click="buyFont(font.id)" disabled="true">已拥有</button>
              <button v-if="!font.own" @click="buyFont(font.id)">兑换</button>
            </div>
          </div>
        </div>
      </el-collapse-transition>
      <el-collapse-transition>
        <div v-if="selectedCategory === 'card'">
          <h2>功能卡商城</h2>
          <div v-for="card in cardList" :key="card.id" class="card-item">
            <div class="cardGoods">
              <img :src="card.cardPreviewLink" alt="" style="width:100px;height:100px;">
              <h3>{{ card.cardName }}</h3>
              <p><img src="../../assets/imgss/pigmoney.webp" alt=""
                  style="width:30px;height:20px;vertical-align: middle;"> {{ card.price }}</p>
              <p class="number">{{ card.number }}</p><button v-if="!card.own" @click="buyCard(card.id)">+</button>
            </div>
          </div>
        </div>
      </el-collapse-transition>
    </div>
  </div>
</template>

<script>
import { fetchPaperList, fetchFontList, fetchColorList, buyPaper, buyFont, buyFontColor, fetchCardList, buyCard } from '@/api/shop'
import { getUserMoney } from '@/api/user'
// import axios from 'axios'
export default {
  name: 'QiaopiShop',
  data() {
    return {
      selectedCategory: 'paper',
      showDetailsArea: true,
      paperList: [],
      fontList: [],
      colorList: [],
      cardList: [], // 新增功能卡列表
      ownedItems: {},
      activeTag: 'paper' // 控制动画的变量
    }
  },
  methods: {
    showDetails(category) {
      this.selectedCategory = category
      this.showDetailsArea = true
      if (category === 'paper') {
        this.fetchPaperList()
      } else if (category === 'font') {
        this.fetchFontList()
      } else if (category === 'ink') {
        this.fetchColorList() // 新增获取字体颜色列表的方法
      } else if (category === 'card') {
        this.fetchCardList() // 新增获取功能卡列表的方法
      }
      this.activeTag = category // 更新 activeTag 来控制动画
    },
    async fetchPaperList() {
      try {
        const response = await fetchPaperList()
        this.paperList = response.data
        // console.log(this.paperList) // 检查数据是否正确加载
      } catch (error) {
        // console.error('获取纸张商城列表失败:', error)
      }
    },
    async fetchFontList() {
      try {
        const response = await fetchFontList()
        this.fontList = response.data
      } catch (error) {
        // console.error('获取纸张商城列表失败:', error)
      }
    },
    async fetchColorList() {
      try {
        const response = await fetchColorList()
        this.colorList = response.data
      } catch (error) {
        // console.error('获取纸张商城列表失败:', error)
      }
    },
    async fetchCardList() {
      try {
        const response = await fetchCardList()
        this.cardList = response.data
      } catch (error) {
        // console.error('获取纸张商城列表失败:', error)
      }
    },
    async buyPaper(paperId) {
      try {
        const response = await buyPaper(paperId)
        const data = response
        if (data.code === 200) {
          this.$message({
            message: data.msg,
            type: 'success'
          })
          // 重新加载纸张列表
          this.fetchPaperList()
          getUserMoney()
        } else {
          this.$message({
            message: data.msg,
            type: 'error'
          })
        }
      } catch (error) {
        // console.error('购买纸张失败:', error)
      }
    },
    async buyFont(fontId) {
      try {
        const response = await buyFont(fontId)
        const data = response
        if (data.code === 200) {
          this.$message({
            message: data.msg,
            type: 'success'
          })
          // 重新加载字体列表
          this.fetchFontList()
          getUserMoney()
        } else {
          this.$message({
            message: data.msg,
            type: 'error'
          })
        }
      } catch (error) {
        // console.error('购买字体失败:', error)
      }
    },
    async buyFontColor(fontColorId) {
      try {
        const response = await buyFontColor(fontColorId)
        const data = response
        if (data.code === 200) {
          this.$message({
            message: data.msg,
            type: 'success'
          })
          // 重新加载字体颜色列表
          this.fetchColorList()
          getUserMoney()
        } else {
          this.$message({
            message: data.msg,
            type: 'error'
          })
        }
      } catch (error) {
        // console.error('购买字体颜色失败:', error)
      }
    },
    async buyCard(CardId) {
      try {
        const response = await buyCard(CardId)
        const data = response
        if (data.code === 200) {
          this.$message({
            message: data.msg,
            type: 'success'
          })
          // 重新加载功能卡列表
          this.fetchCardList()
          getUserMoney()
        } else {
          this.$message({
            message: data.msg,
            type: 'error'
          })
        }
      } catch (error) {
        // console.error('购买功能卡失败:', error)
      }
    }
  },
  mounted() {
    this.fetchPaperList()
  }
}
</script>

<style scoped>
.banner {
  position: relative;
  margin-top: 40px;
  width: 1200px;
  height: 700px;
  /* background-image: url(../../assets/imgs/storebgd.png); */
  background-position: center center;
  background-size: cover;
  z-index: 2;
}
.banner img {
  /* position: absolute;
  top: 0;
  height: 0;
  z-index: 2; */
  width: 100%;
  height: 100%;
}

.tag1,
.tag2,
.tag3,
.tag4 {
  position: absolute;
  left: -30px;
  /* transform: translateX(50%); */
  width: 100px;
  height: 80px;
  cursor: pointer;
  z-index: -1;
}
@keyframes move {
  0%{
    left: -30px;
  }
  100%{
    left: -89px;
  }
}
.tag-active {
  animation: move 0.5s forwards;
}

.tag1 {
  top: 20px;
}

.tag2 {
  top: 120px;
}

.tag3 {
  top: 220px;
}

.tag4 {
  top: 320px;
}

.tag1 img,
.tag2 img,
.tag3 img,
.tag4 img {
  width: 100%;
  height: 100%;
}

.tag1 span,
.tag2 span,
.tag3 span,
.tag4 span {
  position: absolute;
  top: 20px;
  left: 50px;
  z-index: 100;
  color: white;
}

.details {
  position: absolute;
  left: 0;
  top: 0;
  width: 1200px;
  height: 100%;
  padding: 20px;
  box-sizing: border-box;
  pointer-events: none;
}

.details h2 {
  color: rgba(103, 12, 12, 0.255);
}

.paperGoods h3 {
  position: absolute;
  top: 110px;
  left: 60px;
  font-size: 16px;
}

.paperGoods p {
  margin-top: -10px;
  /* vertical-align: top; */
}

.paperGoods {
  position: relative;
  /* position: absolute;
  top: 0;
  left: 0; */
  float: left;
  width: 200px;
  margin-left: 70px;
}

.paperGoods button {
  position: absolute;
  width: 70px;
  height: 30px;
  top: 200px;
  left: 50%;
  transform: translateX(-50%);
  background-color: #956b3c;
  border: 0;
  color: rgb(235, 221, 190);
}

.fontsGoods {
  position: relative;
  /* position: absolute;
  top: 0;
  left: 0; */
  float: left;
  width: 200px;
  margin-left: 25px;
  margin-bottom: 20px;
}

.fontsGoods h3 {
  position: absolute;
  top: 50px;
  left: 50%;
  transform: translateX(-50%);
  font-size: 14px;
}

.fontsGoods p {
  margin-top: -10px;
  /* vertical-align: top; */
}

.fontsGoods button {
  position: absolute;
  width: 70px;
  height: 30px;
  top: 130px;
  left: 50%;
  transform: translateX(-50%);
  background-color: #956b3c;
  border: 0;
  color: rgb(235, 221, 190);
}

.colorGoods {
  position: relative;
  /* position: absolute;
  top: 0;
  left: 0; */
  float: left;
  width: 200px;
  margin-left: 70px;
  margin-bottom:20px;
}

.colorGoods h3 {
  position: absolute;
  top: 110px;
  left: 50%;
  transform: translateX(-50%);
  font-size: 14px;
}

.colorGoods p {
  margin-top: -10px;
  /* vertical-align: top; */
}

.colorGoods button {
  position: absolute;
  width: 70px;
  height: 30px;
  top: 200px;
  left: 50%;
  transform: translateX(-50%);
  background-color: #956b3c;
  border: 0;
  color: rgb(235, 221, 190);
}

.cardGoods {
  position: relative;
  /* position: absolute;
  top: 0;
  left: 0; */
  float: left;
  width: 200px;
  height: 170px;
  margin-left: 25px;
  margin-bottom: 10px;
}

.cardGoods h3 {
  position: absolute;
  top: 70px;
  left: 50%;
  transform: translateX(-50%);
  font-size: 14px;
}

.cardGoods p {
  margin-top: -20px;
}

/* .cardGoods .number {

} */
.cardGoods .number {
  position: relative;
  top: -30px;
  left: -20px;
}

.cardGoods button {
  position: absolute;
  width: 30px;
  height: 30px;
  top: 145px;
  left: 60%;
  transform: translateX(-50%);
  background-color: #956b3c;
  border: 0;
  color: rgb(235, 221, 190);
}
button {
  cursor: pointer;
}
</style>
