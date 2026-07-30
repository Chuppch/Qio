<template>
  <div class="memory-game">
    <h1>翻翻乐</h1>
    <div v-if="!canPlay" class="game-disabled">
      游戏已禁用，您没有剩余的翻翻乐次数。
    </div>
    <div v-else class="game-board">
      <div v-for="(card, index) in cards" :key="index" class="card"
        :class="{ flipped: card.flipped, matched: card.matched }" @click="flipCard(index)">
        <div class="card-inner">
          <div class="card-front">?</div>
          <div class="card-back">
            <img :src="card.imageUrl" :alt="'Card ' + (index + 1)" />
          </div>
        </div>
      </div>
    </div>
    <button v-if="canPlay" @click="resetGame" class="reset-button">重置游戏</button>
  </div>
</template>

<script>
import { Message } from 'element-ui'
import { getFflLimit, winFfl } from '@/api/fanfanle'
export default {
  name: 'MemoryCardGame',
  data() {
    return {
      canPlay: true,
      cards: [],
      flippedCards: [],
      canFlip: true,
      winLimit: null, // 翻翻乐的次数
      images: [
        require('../../assets/imgss/ffl-doubleFish.webp'),
        require('../../assets/imgss/ffl-fish.webp'),
        require('../../assets/imgss/ffl-heteroideusMoney.webp'),
        require('../../assets/imgss/ffl-pangxie.webp'),
        require('../../assets/imgss/ffl-rabbit.webp'),
        require('../../assets/imgss/ffl-realMoney.webp'),
        require('../../assets/imgss/ffl-rectangleMoney.webp'),
        require('../../assets/imgss/ffl-wangba.webp')
      ]
    }
  },
  methods: {
    initializeCards() {
      const duplicatedImages = [...this.images, ...this.images]
      this.cards = duplicatedImages
        .sort(() => Math.random() - 0.5)
        .map((imageUrl, index) => ({
          id: index,
          imageUrl,
          flipped: false,
          matched: false
        }))
    },
    flipCard(index) {
      if (!this.canPlay || !this.canFlip || this.cards[index].flipped || this.cards[index].matched) return

      this.cards[index].flipped = true
      this.flippedCards.push(index)

      if (this.flippedCards.length === 2) {
        this.canFlip = false
        this.checkMatch()
      }
    },
    checkMatch() {
      const [index1, index2] = this.flippedCards
      const card1 = this.cards[index1]
      const card2 = this.cards[index2]

      if (card1.imageUrl === card2.imageUrl) {
        card1.matched = true
        card2.matched = true
        this.flippedCards = []
        this.canFlip = true
        this.checkWin()
      } else {
        setTimeout(() => {
          card1.flipped = false
          card2.flipped = false
          this.flippedCards = []
          this.canFlip = true
        }, 1000)
      }
    },
    resetGame() {
      this.initializeCards()
      this.flippedCards = []
      this.canFlip = true
    },
    checkWin() {
      const allMatched = this.cards.every(card => card.matched)
      if (allMatched) {
        // console.log('游戏胜利')
        this.winFfl()
        this.getFflLimit()
      }
    },
    async winFfl() {
      const res = await winFfl()
      if (res.code === 200) {
        Message.success('恭喜你获得10猪仔钱')
      } else {
        Message.error(res.msg)
      }
    },
    async getFflLimit() {
      const res = await getFflLimit()
      if (res.code === 200) {
        this.winLimit = res.data
        // console.log(this.winLimit)
        if (this.winLimit <= 0) {
          this.canPlay = false // 设置一个标志来禁止游戏
        } else {
          this.canPlay = true // 允许游戏
        }
      } else {
        Message.error(res.msg)
      }
    }
  },
  created() {
    this.initializeCards()
    this.getFflLimit()
  }
}
</script>

<style scoped>
.memory-game {
  font-family: Arial, sans-serif;
  text-align: center;
  padding: 0 20px;
}

.game-board {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 10px;
  max-width: 450px;
  margin: 0 auto 20px;
}

.card {
  aspect-ratio: 3/4;
  perspective: 1000px;
  cursor: pointer;
}

.card-inner {
  position: relative;
  width: 100%;
  height: 100%;
  text-align: center;
  transition: transform 0.6s;
  transform-style: preserve-3d;
}

.card.flipped .card-inner {
  transform: rotateY(180deg);
}

.card-front,
.card-back {
  position: absolute;
  width: 100%;
  height: 100%;
  backface-visibility: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  font-weight: bold;
  border-radius: 8px;
}

.card-front {
  background-color: #b49285;
  color: white;
}

.card-back {
  background-color: #3498db;
  color: white;
  transform: rotateY(180deg);
}

.card.matched .card-back {
  background-color: #2ecc71;
}

.card-back img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 8px;
}

.reset-button {
  padding: 10px 20px;
  font-size: 16px;
  background-color: #a26760;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.reset-button:hover {
  background-color: #c0392b;
}

.game-disabled {
  text-align: center;
  color: rgb(172, 90, 90);
  font-size: 20px;
  margin-top: 20px;
}
</style>
