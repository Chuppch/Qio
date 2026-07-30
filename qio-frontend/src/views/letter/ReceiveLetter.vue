<template>
  <div class="banner">

    <div class="content">
      <el-row class="row-bg" justify="space-around">
        <el-empty v-if="isEmpty" style="height: 100%;" description="没有收到侨批哦,快去结交新好友吧!"></el-empty>
        <el-col :span="6" v-for="item in myReceiveList" :key="item.id">
          <el-badge value="未读" style="margin-right: 40px;" :hidden="item.readStatus === 1">
            <el-skeleton style="width: 250px ;margin:0 auto" :loading="loading" animated>
              <template slot="template">
                <el-skeleton-item variant="image" style="width: 250px; height: 250px;" />
                <div style="padding: 2px 10px 0">
                  <el-skeleton-item variant="text" style="margin-right: 16px; " />
                </div>
              </template>
              <template>
                <!-- 启用懒加载 -->
                <el-image ref="firstImage" class="image-hover" v-if="item.coverLink !== undefined"
                  style="width: 250px; height: 250px; cursor: pointer;transition: transform .3s ease; " lazy
                  @click="openImageViewerPre($event, item)" :src="item.coverLink">
                  <div slot="placeholder" class="image-slot" style=" background:transparent">
                    <i class="el-icon-loading" style="line-height: 250px; color: #A52328; font-size: 30px;"></i>
                  </div>
                  <div slot="error" class="image-slot" style="width: 100%; height: 100%; color: #A52328;">
                    <i class="el-icon-picture-outline" style="line-height: 250px; "></i>
                    加载失败
                  </div>
                </el-image>
                <p style="color: #A52328; margin: -15px 0 0; ">
                  {{ formatDate(item.createTime) }},
                  <span
                    style="margin: -20px 0 -14px;max-width: 100px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; display: inline-block;">
                    {{ item.senderName }}</span>记
                </p>
              </template>
            </el-skeleton>
          </el-badge>

        </el-col>
      </el-row>
    </div>

    <transition name="el-zoom-in-center">
      <div v-if="isVisible" class="el-image-viewer__wrapper" :style="{ zIndex: 2002 }">
        <div class="el-image-viewer__mask" @click="closeImageViewer"></div>
        <span class="el-image-viewer__btn el-image-viewer__close" @click="closeImageViewer">
          <i class="el-icon-close"></i>
        </span>
        <div class="el-image-viewer__btn el-image-viewer__actions">
          <div class="el-image-viewer__actions__inner">
            <div><i class="el-icon-zoom-out" @click="zoomOut"></i></div>
            <div><i class="el-icon-zoom-in" @click="zoomIn"></i></div>
            <div><i class="el-image-viewer__actions__divider"></i></div>
            <div><i class="el-icon-full-screen" :class="{ 'toggle-icon': isFullscreen }" @click="toggleFullScreen"></i>
            </div>
            <div><i class="el-image-viewer__actions__divider"></i></div>
            <div><i class="el-icon-refresh-left" @click="rotateLeft"></i></div>
            <div><i class="el-icon-refresh-right" @click="rotateRight"></i></div>
          </div>
        </div>
        <div class="el-image-viewer__canvas" @mousedown="startDrag" @mousemove="dragImage" @mouseup="endDrag"
          @wheel.prevent="zoomWithWheel">
          <img :src="imageUrl" class="el-image-viewer__img" :class="{ fullscreen: isFullscreen }" :style="{
            transform: `scale(${scale}) rotate(${rotation}deg)`,
            marginLeft: `${marginLeft}px`,
            marginTop: `${marginTop}px`,
          }" style="max-width: 100%; max-height: 100%;  transition: transform 0.3s ease;">
        </div>
        <!-- 侧边栏 -->
        <div class="aside">
          <el-card class="box-card">
            <div class="card_item">
              <i class="el-icon-s-custom"></i>
              收信人：{{ letterVo.recipientName }}
            </div>
            <div class="card_item">
              <i class="el-icon-map-location"></i>
              收信人地址：{{ letterVo.recipientAddress.formattedAddress }}
            </div>
            <div class="card_item">
              <i class="el-icon-position"></i>
              阅读状态：{{ getLetterStatusLabel(letterVo.readStatus) }}
            </div>
            <div class="card_item">
              <i class="el-icon-date"></i>
              写信时间：{{ formatDateTime(letterVo.createTime) }}
            </div>
            <div class="card_item">
              <i class="el-icon-message"></i>
              寄信人邮箱：{{ letterVo.senderEmail }}
            </div>
            <div class="card_item">
              <i class="el-icon-s-home"></i>
              寄信人地址：{{ letterVo.senderAddress.formattedAddress }}
            </div>
            <div class="card_item">
              <i class="el-icon-user"></i>
              寄信人：{{ letterVo.senderName }}
            </div>
            <div class="card_item" v-show="letterVo.piggyMoney > 0">
              <i class="el-icon-money"></i>
              携带猪仔钱：{{ letterVo.piggyMoney }}
            </div>
            <div style="margin-left: 20px;">
              <i class="el-icon-document-add"></i>
              <el-button type="text" @click="writeMore">写一封回信</el-button>
              <!-- <el-divider direction="vertical"></el-divider>
                <i class="el-icon-delete"></i>
                <el-button type="text">忍痛</el-button> -->
            </div>
          </el-card>

          <div class="aside_close" @click="handleShowDetail">
            <i v-if="showLetterDetail" class="el-icon-arrow-left"></i>
            <i v-else class="el-icon-arrow-right"></i>
          </div>
        </div>
      </div>
    </transition>
    <!-- 使用transition组件来添加动画 -->
    <transition name="el-fade-in-linear">
      <div v-if="showMask" class="el-image-viewer__wrapper" :style="{ zIndex: 999 }">
        <div class="el-image-viewer__mask" @click="closeImageViewer"></div>
      </div>
    </transition>

  </div>

</template>

<script>
import { getMyReceiveLetter, readLetter } from '@/api/letter'

export default {
  name: 'ReceiveLetter',
  data() {
    return {
      letterStatusMap: {
        0: '未读',
        1: '已读'
      },
      // 查课侨批大图组件
      showProgress: true,
      showLetterDetail: true,
      showMask: false,
      isVisible: false,
      letterVo: {
        letterId: '',
        senderName: '',
        recipientName: '',
        senderAddress: {
          formattedAddress: '',
          longitude: '',
          latitude: ''
        },
        recipientAddress: {
          formattedAddress: '',
          longitude: '',
          latitude: ''
        },
        recipientEmail: '',
        letterLink: '',
        coverLink: '',
        deliveryProgress: '',
        createTime: '',
        expectedDeliveryTime: '',
        status: ''
      },
      coverUrl: 'http://110.41.58.26:9000/qiaopi/qiaopi-images/cover/823eec55-0d12-41bd-a013-483892daf166.png',
      imageUrl: 'http://110.41.58.26:9000/qiaopi/qiaopi-images/letter/7780d647-ac9a-4ac0-9342-ad35e9e709d1.png',
      scale: 0,
      rotation: 0,
      marginLeft: 0,
      marginTop: 0,
      isDragging: false,
      initialX: 0,
      initialY: 0,
      startX: 0,
      startY: 0,
      isFullscreen: false,
      // 信件图片列表
      myReceiveList: ['', '', '', '', '', '', '', ''],
      isEmpty: false,

      // 图片加载状态
      loading: false
    }
  },
  methods: {
    writeMore() {
      this.closeImageViewer()
      setTimeout(() => {
        this.$router.push('/write')
      }, 2300)
    },
    handleShowDetail() {
      this.showLetterDetail = !this.showLetterDetail
      const aside = document.querySelector('.aside')
      if (this.showLetterDetail) {
        aside.style.left = '0'
      } else {
        aside.style.left = '-270px'
        // aside.classList.toggle('aside_close')
      }
    },
    formatDateTime(dateString) {
      const date = new Date(dateString)
      const year = date.getFullYear()
      const month = String(date.getMonth() + 1).padStart(2, '0')
      const day = String(date.getDate()).padStart(2, '0')
      const hours = String(date.getHours()).padStart(2, '0')
      const minutes = String(date.getMinutes()).padStart(2, '0')
      const seconds = String(date.getSeconds()).padStart(2, '0')
      return `${year}年${month}月${day}日 ${hours}:${minutes}:${seconds}`
    },
    formatDate(dateString) {
      const date = new Date(dateString)
      const year = date.getFullYear()
      const month = String(date.getMonth() + 1).padStart(2, '0')
      const day = String(date.getDate()).padStart(2, '0')
      return `${year}年${month}月${day}日`
    },
    // 查看大图函数开始
    openImageViewerPre(event, item) {
      // console.log(event)
      this.showMask = true
      const imageElement = document.createElement('img')
      imageElement.setAttribute('style', 'width: 250px; height:250px; position: absolute; top: 0; left: 0; z-index: 1000; ')
      imageElement.setAttribute('v-if', 'item.coverLink !== undefined')
      imageElement.setAttribute('src', item.coverLink)

      // 将新元素插入到当前点击的图片元素之上
      const targetElement = event.target
      targetElement.parentNode.insertBefore(imageElement, targetElement.nextSibling)

      // 设置新元素的样式
      imageElement.style.position = 'absolute'
      imageElement.style.top = '0'
      imageElement.style.left = '0'
      imageElement.style.zIndex = '1000'

      // 获取当前点击图片的位置
      const currentImageRect = targetElement.getBoundingClientRect()
      imageElement.style.left = `${currentImageRect.left}px`
      imageElement.style.top = `${currentImageRect.top}px`

      // 处理点击关闭事件
      // imageElement.addEventListener('click', () => {
      //   targetElement.parentNode.removeChild(imageElement)
      // })

      // 移到 .content 层中
      const contentElement = document.querySelector('.content')
      contentElement.appendChild(imageElement)

      // 计算当前图片在 .content 层中的位置
      const contentRect = contentElement.getBoundingClientRect()
      const newTop = currentImageRect.top - contentRect.top + 23
      const newLeft = currentImageRect.left - contentRect.left + 20

      // 设置新的位置
      imageElement.style.left = `${newLeft}px`
      imageElement.style.top = `${newTop}px`

      // 计算居中的位置
      const contentWidth = contentElement.offsetWidth
      const imageWidth = imageElement.offsetWidth
      const leftPosition = (contentWidth - imageWidth) / 2
      imageElement.style.setProperty('--pre-left', `${newLeft}px`)
      imageElement.style.setProperty('--pre-top', `${newTop}px`)
      imageElement.style.setProperty('--pre-transform', 'translateY(0)')
      // 设置目标位置
      imageElement.style.setProperty('--target-left', `${leftPosition}px`)
      imageElement.style.setProperty('--target-top', '50%')
      imageElement.style.setProperty('--target-transform', 'translateY(-50%)')

      // 触发动画
      imageElement.classList.add('animate-move-center')

      // 动画结束后清除动画类
      imageElement.addEventListener('animationend', () => {
        this.openImageViewer(item)
        // 监听 isVisible 变化
        this.$watch('isVisible', (newValue, oldValue) => {
          if (!newValue && imageElement.parentNode) {
            imageElement.classList.add('animate-back-to-original')
            imageElement.classList.remove('animate-move-center')
            this.showMask = false
            imageElement.addEventListener('animationend', () => {
              imageElement.classList.remove('animate-back-to-original')
              imageElement.parentNode.removeChild(imageElement)
            }, { once: true })
          }
        })
      }, { once: true })
    },
    openImageViewer(item) {
      this.isVisible = true

      this.imageUrl = item.letterLink
      // 设置 letterVo
      this.letterVo = item
      if (item.readStatus === 0) {
        this.readLetter(item.id)
      }
      setTimeout(() => {
        if (!this.showLetterDetail) {
          this.handleShowDetail()
        }
        this.scale = 0.9
      }, 100)
    },
    readLetter(letterId) {
      readLetter(letterId).then(res => {
        // console.log(res)
        if (res.code === 200) {
          this.letterVo.readStatus = 1
          this.getMyReceiveLetter()
        }
      })
    },
    closeImageViewer() {
      // this.isVisible = false
      this.isFullscreen = false

      if (this.showLetterDetail) {
        this.handleShowDetail()
      }
      this.scale = 0
      setTimeout(() => {
        this.isVisible = false
      }, 301) // 假设动画持续时间为300ms
      this.marginLeft = 0
      this.marginTop = 0
    },
    zoomIn() {
      this.scale += 0.1
      if (this.scale > 3) {
        this.scale = 3
      }
      // 如果缩放比例在 1.4 到 1.5 之间，且显示详情，则触发隐藏详情
      if (this.scale >= 1.4 && this.scale <= 1.5 && this.showLetterDetail) {
        this.handleShowDetail()
      }
      if (this.scale >= 1.2 && this.showProgress) {
        this.handleShowProgress()
      }
    },
    zoomOut() {
      this.scale -= 0.1
      if (this.scale < 0.1) {
        this.scale = 0.1
      }
      if (this.scale <= 1.2 && !this.showProgress) {
        this.handleShowProgress()
      }
    },
    toggleFullScreen() {
      this.isFullscreen = !this.isFullscreen
      // 非全屏状态下的设置
      this.scale = 1
      this.marginLeft = 0
      this.marginTop = 0
      if (this.showLetterDetail) {
        this.handleShowDetail()
      }
      if (this.showProgress) {
        this.handleShowProgress()
      }
    },
    rotateLeft() {
      this.rotation -= 90
    },
    rotateRight() {
      this.rotation += 90
    },
    startDrag(event) {
      event.preventDefault()
      // const imgElement = this.$el.querySelector('.el-image-viewer__img')
      // const imgRect = imgElement.getBoundingClientRect()

      this.startX = event.clientX
      this.startY = event.clientY

      this.initialX = this.marginLeft
      this.initialY = this.marginTop

      this.isDragging = true

      window.addEventListener('mousemove', this.dragImage)
      window.addEventListener('mouseup', this.endDrag)
    },
    dragImage(event) {
      if (!this.isDragging) return

      event.preventDefault()
      const newX = this.initialX + event.clientX - this.startX
      const newY = this.initialY + event.clientY - this.startY

      this.marginLeft = newX
      this.marginTop = newY
    },
    endDrag() {
      this.isDragging = false
      window.removeEventListener('mousemove', this.dragImage)
      window.removeEventListener('mouseup', this.endDrag)
    },
    zoomWithWheel(event) {
      event.preventDefault()
      const delta = -event.deltaY / 5000 // 调整缩放速度
      this.scale += delta
      if (this.scale < 0.1) {
        this.scale = 0.1
      }
      if (this.scale > 3) {
        this.scale = 3
      }
    },
    // 查看大图函数结束
    getLetterStatusLabel(statusId) {
      return this.letterStatusMap[statusId] || '未知状态'
    },
    getMyReceiveLetter() {
      this.loading = true
      getMyReceiveLetter().then(res => {
        this.myReceiveList = res.data
        this.loading = false
        if (this.myReceiveList.length === 0) {
          this.isEmpty = true
        }
        // 从路由获取
        if (this.$route.params.showFirst) {
          // console.log(this.$route.params)
          this.loadFirstImage(this.$route.params.myNotReadLetter)
          this.$route.params.showFirst = false
        }
      })
    },
    loadFirstImage(item) {
      setTimeout(() => {
        // console.log(item)
        this.showMask = true
        const imageElement = document.createElement('img')
        imageElement.setAttribute('style', 'width: 250px; height: 250px; position: absolute; top: 0; left: 0; z-index: 1000')
        imageElement.setAttribute('src', item.coverLink)

        // 获取第一个图片元素
        const firstImageElement = document.querySelector('.content .row-bg .el-col .el-image')
        if (!firstImageElement) {
          // console.error('No image element found')
          this.showMask = false
          return
        }

        // 将新元素插入到 .content 层中
        const contentElement = document.querySelector('.content')
        contentElement.appendChild(imageElement)

        // 获取第一个图片元素的位置
        const currentImageRect = firstImageElement.getBoundingClientRect()

        // 设置新元素的样式
        imageElement.style.position = 'absolute'
        imageElement.style.top = '0'
        imageElement.style.left = '0'
        imageElement.style.zIndex = '1000'

        // 计算当前图片在 .content 层中的位置
        const contentRect = contentElement.getBoundingClientRect()
        const newTop = currentImageRect.top - contentRect.top + 10
        const newLeft = currentImageRect.left - contentRect.left + 10

        // 设置新的位置
        imageElement.style.left = `${newLeft}px`
        imageElement.style.top = `${newTop}px`

        // 计算居中的位置
        const contentWidth = contentElement.offsetWidth
        const imageWidth = imageElement.offsetWidth
        const leftPosition = (contentWidth - imageWidth) / 2
        imageElement.style.setProperty('--pre-left', `${newLeft}px`)
        imageElement.style.setProperty('--pre-top', `${newTop}px`)
        imageElement.style.setProperty('--pre-transform', 'translateY(0)')
        // 设置目标位置
        imageElement.style.setProperty('--target-left', `${leftPosition}px`)
        imageElement.style.setProperty('--target-top', '50%')
        imageElement.style.setProperty('--target-transform', 'translateY(-50%)')

        // 触发动画
        imageElement.classList.add('animate-move-center')

        // 动画结束后清除动画类
        imageElement.addEventListener('animationend', () => {
          this.openImageViewer(item)

          // 监听 isVisible 变化
          this.$watch('isVisible', (newValue, oldValue) => {
            if (!newValue && imageElement.parentNode) {
              imageElement.classList.add('animate-back-to-original')
              imageElement.classList.remove('animate-move-center')
              this.showMask = false
              imageElement.addEventListener('animationend', () => {
                imageElement.classList.remove('animate-back-to-original')
                imageElement.parentNode.removeChild(imageElement)
              }, { once: true })
            }
          })
        }, { once: true })
      }, 500)
    }
  },
  mounted() {
    this.getMyReceiveLetter()
  }
}
</script>

<style scoped>
/* 滑动轨道按钮 */
::-webkit-scrollbar-button {
  width: 10px;
  height: 10px;
  display: none;
}

.row-bg {
  width: 95%;
  height: 600px;
  position: relative;
  padding-top: 20px;
  padding-right: 20px;
  margin: 40px auto 0 auto;
  flex-wrap: wrap;
  border-radius: 20px;
  /* 明确允许换行 */
  overflow: hidden;
  overflow-y: scroll;
  /* z-index: 10; */
  background-color: rgba(255, 255, 255, 0.5);
}

.row-bg::-webkit-scrollbar-track {
  background-color: transparent !important;
}

.content {
  flex-grow: 1;
  /* 占据剩余空间 */
  padding: 20px;
  /* 添加边框 */
}

.banner {
  position: relative;
  margin-top: 40px;
  width: 90%;
  height: 680px;
  border-radius: 20px;
  /* height: 1200px; */
  background-color: transparent;
  background: url('@/assets/imgss/yellowbackground.png') 0 0 / 400px auto repeat, #f9f9f9;
  /* display: flex; */
  /* 使用 Flexbox 布局 */
  align-items: flex-start;
  /* 垂直对齐子元素 */
  line-height: 40px;
}

::v-deep .el-dialog {
  margin-top: 200px !important;
  border-radius: 10px;
  background-color: rgba(222, 201, 162, 1) !important;
}

::v-deep input::-webkit-input-placeholder {
  color: #666;
  font-size: 15px;
}

::v-deep .el-form-item {
  margin-bottom: 20px;
}

.demo-table-expand {
  font-size: 0;
}

.demo-table-expand>>>.el-form-item__content {
  margin-left: -5px !important;
}

.demo-table-expand label {
  width: 90px;
  color: #99a9bf;
}

.dialog-footer {
  text-align: center;
}

.fullscreen {
  transition: .3s ease;
  max-width: none !important;
  max-height: none !important;
  margin: 0 !important;
}

.el-icon-full-screen.toggle-icon::before {
  content: "\e7c6";
}

.el-image-viewer__actions__inner {
  display: block
}

.el-image-viewer__actions {
  left: 94%;
  bottom: 4%;
  transform: translateX(-12%);
  height: 210px;
  width: 40px;
  padding: 8px 9px;
  background-color: #606266;
  border-color: #fff;
  border-radius: 25px;
}

/* 定义动画效果 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 2s;
  /* 过渡属性和持续时间 */
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  /* 初始和结束时的透明度 */
}

/* 分割线· */
/* 定义动画效果 */
.zoom-in-enter-active,
.zoom-in-leave-active {
  transition: opacity .3s;
  /* 过渡属性和持续时间 */
}

.zoom-in-from,
.zoom-in-to {
  opacity: 0;
  /* 初始和结束时的透明度 */
}

.image-hover:hover,
.image-hover:focus,
.image-hover:active {
  transform: scale(1.1);
}

.aside {
  position: fixed;
  left: 0%;
  top: 0;
  height: 100%;
  display: flex;
  background-color: rgba(0, 0, 0, .7);
  transition: 1s ease-in-out;
  border-radius: 0 10px 10px 0;

}

.box-card {
  width: 280px;
  flex: 6;
  background: transparent;
  color: #fff;
  font-size: 12px;
  margin-top: 100px;
  border: none;
  display: flex;
  align-items: center;
  margin-bottom: 10px;
}

.aside_close {
  flex: 1;
  /* position: absolute;
  right: 0;
  top: 0; */
  width: 35px;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  line-height: 30px;
  text-align: center;
  cursor: pointer;
  color: #fff;
  transition: 1s ease-in-out;
  /* background: rgba(0, 0, 0, .7); */
}

.card_item {
  text-align: left;
}

.progress {
  position: fixed;
  left: 50%;
  bottom: 8%;
  transform: translate(-50%, 50%);
  width: 500px;
  height: 100px;
  background: rgba(0, 0, 0, .7);
  border-radius: 10px;
  transition: 1s ease-in-out;
  display: flex;
  justify-content: center;
  align-items: center;
  opacity: 1;
}

/* 平移动画 */
/* 移至common.css */
</style>
