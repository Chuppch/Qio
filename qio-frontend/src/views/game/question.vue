<template>
  <div class="banner">
    <img src="../../assets/imgs/knowledge.png" alt="" class="bgd">
    <div v-if="currentQuestion" class="question-container">
      <h2>{{ currentQuestion.content }}</h2>
      <div v-for="(option, index) in currentQuestion.options" :key="index" class="option-container">
        <input type="radio" :id="'option-' + index" :value="index" v-model="userAnswerIndex">
        <label :for="'option-' + index">{{ getLetter(index) }}: {{ option }}</label>
      </div>
      <el-button @click="goToPreviousQuestion">上一题</el-button>
      <el-button @click="goToNextQuestion" :disabled="!hasSubmittedAnswer"
        :class="{ 'btn--disabled': currentQuestionIndex !== contents.length - 1 }">
        {{ currentQuestionIndex === contents.length - 1 ? '查看获得的猪仔钱' : '下一题' }}
      </el-button>
      <el-button @click="checkAnswer">提交答案</el-button>
      <div v-if="showExplanation" class="explanation">
        <h3>正确答案： {{ currentQuestion.correctAnswer }}</h3>
        <p>解析：{{ currentQuestion.explanation }}</p>
      </div>
      <img :src="thinkImage" alt="" class="think">
      <!-- <img src="../../assets/imgs/think.png" alt="" class="think"> -->
    </div>
    <div v-if="showResults" class="results">
      <h2>答对题数：{{ pigMoney / 5 }}</h2>
      <h2>获得猪仔钱：{{ pigMoney }}</h2>
      <el-button @click="restartQuiz">再答一次</el-button>
    </div>
  </div>
</template>

<script>
import { Message } from 'element-ui'
import { submitAnswers } from '@/api/know'
import think from '@/assets/imgs/think.png'
import yes from '@/assets/imgs/yes.png'
import wrong from '@/assets/imgs/wrong.png'
export default {
  name: 'questionPage',
  data() {
    return {
      pigMoney: 0,
      questions: [],
      contents: [],
      currentQuestion: null,
      userAnswerIndex: null, // 存储用户选择的选项索引
      showExplanation: false,
      currentQuestionIndex: 0,
      correctAnswerIndex: null,
      answeredQuestions: [], // 存储所有已回答的问题
      hasAnsweredCurrentQuestion: false,
      hasSubmittedAnswer: false, // 新增状态，表示是否已提交答案
      showResults: false,
      correctCount: 0,
      thinkImage: think
    }
  },
  computed: {
    userAnswer() { // 计算属性，将索引转换为字母
      return this.userAnswerIndex !== null ? String.fromCharCode(65 + this.userAnswerIndex) : ''
    }
  },
  watch: {
    userAnswerIndex(newIndex) { // 监听索引变化，更新计算属性
      this.updateUserAnswerLetter(newIndex)
    }
  },
  created() {
    this.fetchQuestions()
  },
  methods: {
    fetchQuestions() {
      this.contents = this.$route.query.questions
      if (this.contents.length > 0) {
        this.setCurrentQuestion(0)
      }
    },
    setCurrentQuestion(index) {
      if (index >= 0 && index < this.contents.length) {
        this.currentQuestionIndex = index
        const questionData = this.contents[index]
        this.currentQuestion = {
          id: questionData.id,
          content: questionData.content,
          options: [questionData.optionA, questionData.optionB, questionData.optionC, questionData.optionD],
          correctAnswer: questionData.correctAnswer,
          explanation: questionData.explanation
        }
        this.correctAnswerIndex = this.currentQuestion.options.indexOf(questionData.correctAnswer)
        this.userAnswerIndex = null
        this.showExplanation = false
        this.hasSubmittedAnswer = false // 重置已提交答案的状态
      }
    },
    checkAnswer() {
      if (this.hasSubmittedAnswer) {
        Message.error('你已经提交过答案了！')
        return
      }
      if (this.userAnswerIndex === null) {
        Message.error('请先选择一个答案！')
        return
      }
      if (this.getLetter(this.userAnswerIndex) === this.currentQuestion.correctAnswer) {
        Message.success('正确答案！')
        this.thinkImage = yes // 答对
      } else {
        Message.error('答案错误，请看解析！')
        this.thinkImage = wrong // 答错
      }
      this.showExplanation = true
      this.saveAnswer()
      this.hasSubmittedAnswer = true // 标记为已提交答案
    },
    saveAnswer() {
      const answer = {
        questionId: this.currentQuestion.id,
        questionContent: this.currentQuestion.content,
        selectedOption: this.userAnswer
      }
      this.answeredQuestions.push(answer)
    },
    goToPreviousQuestion() {
      if (this.currentQuestionIndex > 0) {
        this.setCurrentQuestion(this.currentQuestionIndex - 1)
      } else {
        Message.info('这是第一题哦！')
      }
    },
    goToNextQuestion() {
      if (this.hasSubmittedAnswer) {
        if (this.currentQuestionIndex === this.contents.length - 1) {
          // 如果是最后一题，调用 submitAnswers 方法
          this.submitAnswers()
        } else {
          this.thinkImage = think // 重置图片为思考状态
          // 如果不是最后一题，继续下一题
          this.setCurrentQuestion(this.currentQuestionIndex + 1)
          this.hasSubmittedAnswer = false // 重置状态，允许用户在下一题提交答案
        }
      } else {
        Message.info('请先提交当前答案再进入下一题！')
      }
    },
    getLetter(index) { // 将索引转换为字母
      return String.fromCharCode(65 + index)
    },
    updateUserAnswerLetter(index) { // 更新userAnswerLetter
      this.$emit('update:userAnswer', this.getLetter(index))
    },
    async submitAnswers() {
      // const userStore = useUserStore()
      // const token = userStore.token // 从 Pinia store 获取 token
      // 构建提交的数据
      const answers = this.answeredQuestions.map((question, index) => ({
        questionId: question.questionId,
        questionContent: question.questionContent,
        selectedOption: question.selectedOption
      }))
      try {
        const response = await submitAnswers(answers)
        Message.success('答案提交成功！')
        this.handleResponse(response.data)
        // console.log(this.pigMoney)
      } catch (error) {
        // console.error('提交答案失败:', error)
        Message.error('提交答案失败')
      }
    },
    handleResponse(data) {
      this.showResults = true
      this.pigMoney = data.pigMoney
      this.currentQuestion = null
    },
    restartQuiz() {
      this.$router.push('/know')
    }
  }
}
</script>

<style scoped>
.banner {
  position: relative;
  width: 100%;
  height: 650px;
  margin-top: 55px;
  background: url(../../assets/imgs/sea.jpg) no-repeat;
  background-size: cover;
}

.bgd {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 90%;
}

.question-container {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 80%;
  max-width: 800px;
  /* Ensure content does not exceed the image */
  padding: 20px;
}

.question-container h2 {
  margin-top: 50px;
  font-size: 22px;
}

.option-container {
  display: flex;
  align-items: center;
  margin-left: 310px;
}

.option-container label {
  margin-left: 8px;
  text-align: left;
}

.explanation {
  margin-top: 5px;
}

.explanation h3 {
  margin-top: -30px;
}

.explanation p {
  margin-top: -30px;
}

.results {
  /* 定位和布局 */
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 80%;
  max-width: 800px;
  /* 限制最大宽度 */
  padding: 20px;
  box-sizing: border-box;

  /* 文本样式 */
  text-align: center;
  color: black;
  /* 文本颜色 */
  font-family: 'Arial', sans-serif;
  /* 字体 */

  /* 背景和边框 */
  background-color: rgba(255, 255, 255, 0.8);
  /* 半透明背景 */
  border-radius: 8px;
  /* 圆角边框 */
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  /* 阴影 */

  /* 动画 */
  animation: fadein 1s ease-in-out forwards;
  /* 淡入动画 */
}

.think {
  position: absolute;
  width: 186px;
  /* 或者你需要的尺寸 */
  height: 168px;
  /* 保持图片比例 */
  right: -50px;
  bottom: 50px;
}
</style>
