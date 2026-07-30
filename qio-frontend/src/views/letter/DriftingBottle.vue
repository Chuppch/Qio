<template>
  <div class="banner">
    <div class="pick" v-if="activeButton === 'smallbtn1'">
      <!-- <button class="bigbtn" @click="showDriftBottle(); select('smallbtn2')">漂流瓶</button> -->
      <button class="bigbtn" @click="showDriftBottle();">
        漂流瓶
      </button>
    </div>
    <div class="letter" v-if="activeButton === 'smallbtn2'" :style="'background-image: url(' + letter.url + ');'">
      <div class="friend-modal" v-if="isFriendModalVisible">
        <div class="modal-content">
          <span class="close" @click="toggleFriendModal">&times;</span>
          <h2>成为好友</h2>
          <form @submit.prevent="sendFriend">
            <div class="map">
              <avue-input-map :autosize="{ minRows: 1, maxRows: 4 }" placeholder="我的地址"
                v-model="friendRequest.giveAddresss"
                style="width: 100%; max-width: 500px; margin: 0 auto;"></avue-input-map>
            </div>
            <div class="form-group special">
              <label for="friendContent" style="margin-bottom: 20px;">内容:</label>
              <textarea id="friendContent" v-model="friendRequest.context"
                style="width: 100%; max-width: 500px; margin: 0 auto; height: 100px;"></textarea>
            </div>
            <button type="submit" @click="sendFriendRequest" style="position: absolute;bottom: 120px;">发送</button>
          </form>
        </div>
      </div>
      <button class="bigbtn1" @click="throwBottle">扔回海里</button>
      <button class="bigbtn2" @click="toggleFriendModal">成为好友</button>
    </div>
    <div class="write" v-if="activeButton === 'smallbtn3'">
      <div class="left">
        <h2>写漂流瓶</h2>
        <form @submit.prevent="sendLetter">
          <div class="map">
            <avue-input-map :autosize="{ minRows: 1, maxRows: 4 }" placeholder="请选择寄信地址" v-model="letter.senderAddress"
              style="width:80%;margin:0 auto;"></avue-input-map>
          </div>
          <div class="form-group special">
            <label for="letterContent">信件内容：</label>
            <textarea id="letterContent" v-model="letter.content" required style="margin-top:20px"></textarea>
          </div>
          <button type="submit" @click="generateDriftBottle" style="width: 40%;">发送</button>
        </form>
      </div>
      <div class="right" :style="{ 'background-image': `url(${bottleImageUrl})` }"></div>
    </div>
    <div class="friend" v-if="activeButton === 'smallbtn5'"></div>
    <div class="smallbtn1" @click="select('smallbtn1')" :class="activeButtonClass('smallbtn1')">捡漂流瓶</div>
    <div class="smallbtn2" @click="select('smallbtn4');" :class="activeButtonClass('smallbtn4')">
      新消息
      <div v-if="hasUnreadRequests" class="red-dot"></div>
    </div>
    <div class="smallbtn3" @click="select('smallbtn3')" :class="activeButtonClass('smallbtn3')">写漂流瓶</div>
    <!-- haha  -->
    <div class="friend" v-if="activeButton === 'smallbtn4'">
      <div class="friendquest">
        <!-- 如果有好友申请，显示每个申请 -->
        <div v-if="friendRequests && friendRequests.length > 0">
          <div v-for="request in friendRequests" :key="request.id" class="request-content">
            <div class="content">
              <img :src="request.senderAvatar" alt="Avatar" class="avatar" />
              <div class="username">{{ request.senderName }}：</div>
              <div class="context">{{ request.content }}</div>
              <div class="time">{{ new Date(request.createTime).toLocaleString() }}</div>
            </div>
            <button @click="acceptFriendRequest(request.id)" class="button1">接受</button>
            <button @click="rejectFriendRequest(request.id)" class="button2">拒绝</button>
          </div>
        </div>
        <div v-else class="no-requests">
          您还没有收到好友申请哦
        </div>
      </div>
    </div>
    <!-- haha  -->
  </div>
</template>

<script>
import { Message } from 'element-ui'
// import axios from 'axios'
// import useUserStore from '@/store/modules/user'
import { showDriftBottle, generateDriftBottle, throwBackBottle, sendFriendRequest, ProcessingFriendRequests, BecomeFriend } from '@/api/drifting'
export default {
  name: 'DriftingBottle',
  data() {
    return {
      senderAddress: [116.397455, 39.909187, '北京市东城区东华门街道天安门'],
      // giveAddresss: [116.397455, 39.909187, '北京市东城区东华门街道天安门'],
      currentBottleId: null, // 存储当前漂流瓶的 ID
      requestStatus: {},
      friendRequests: [],
      activeButton: 'smallbtn1',
      isFriendModalVisible: false, // 控制成为好友弹窗的显示
      friendRequest: { // 初始化 friendRequest 对象
        giveAddresss: {
          id: 0,
          formattedAddress: '',
          longitude: 0,
          latitude: 0,
          isDefault: 'ad'
        },
        context: '希望与你成为好友~'
      },
      letter: {
        senderAddress: {
          id: 0,
          formattedAddress: '',
          longitude: 0,
          latitude: 0,
          isDefault: 'ad'
        },
        content: ''
      },
      bottleImageUrl: 'http://110.41.58.26:9000/qiaopi/qiaopi-images/paper/plp.webp'
    }
  },
  methods: {
    select(button) {
      this.activeButton = button// 更新选中的按钮
    },
    newmessage() {
      if (this.activeButton === 4) {
        this.ProcessingFriendRequests()
      }
    },
    pickdrifting() {
      if (this.activeButton === 1) {
        this.showDriftBottle()
      }
    },
    toggleFriendModal() {
      this.isFriendModalVisible = !this.isFriendModalVisible
    },
    showWrite() {
      this.activeButton = 'smallbtn1' // 假设 showWrite 就是激活 'smallbtn3' 对应的内容
    },
    // async getDriftBottle() {
    //   const response = await getDriftBottle()
    //   if (response.code === 200) {
    //     this.letter.url = response.data.bottleUrl
    //     Message.success(response.msg)
    //   } else {
    //     Message.error(response.msg)
    //   }
    // },
    async showDriftBottle() {
      // 睡眠一秒钟
      await new Promise(resolve => setTimeout(resolve, 500))
      this.select('smallbtn2')
      const response = await showDriftBottle()
      if (response.code === 200) {
        this.letter.url = response.data
        // this.currentBottleId = response.id
        Message.success(response.msg)
        // console.log(this.currentBottleId)
        this.$forceUpdate() // 强制更新视图
      } else {
        Message.error(response.msg)
      }
    },
    // async showDriftBottleAgain() {
    //   const response = await showDriftBottle()
    //   if (response.code === 200) {
    //     this.throwBackBottle()
    //     this.letter.url = response.data
    //     // this.currentBottleId = response.id
    //     Message.success(response.msg)
    //   } else {
    //     Message.error(response.msg)
    //   }
    // },
    async generateDriftBottle() {
      // 将senderAddress数组转换为对象格式
      if (this.letter.senderAddress) {
        this.letter.senderAddress = {
          id: this.letter.senderAddress.id || 0, // 根据实际情况设置
          formattedAddress: this.letter.senderAddress[2], // 地址字符串
          longitude: this.letter.senderAddress[0], // 经度
          latitude: this.letter.senderAddress[1], // 纬度
          isDefault: this.letter.senderAddress.isDefault || 'ad'
        }
      }
      if (!this.letter.senderAddress || !this.letter.content) {
        Message.error('请填写完整的发送者地址和内容')
        return
      }
      const response = await generateDriftBottle(this.letter)
      if (response.code === 200) {
        Message.success(response.msg)
        this.bottleImageUrl = response.data
      } else {
        Message.error(response.msg)
      }
    },
    async throwBottle() {
      // 睡眠一秒钟
      await new Promise(resolve => setTimeout(resolve, 1000))
      await throwBackBottle()
      this.activeButton = 'smallbtn1'
    },
    async sendFriendRequest() {
      if (this.friendRequest.giveAddresss) {
        this.friendRequest.giveAddresss = {
          id: this.friendRequest.giveAddresss.id || 0, // 根据实际情况设置
          formattedAddress: this.friendRequest.giveAddresss[2], // 地址字符串
          longitude: this.friendRequest.giveAddresss[0], // 经度
          latitude: this.friendRequest.giveAddresss[1], // 纬度
          isDefault: this.friendRequest.giveAddresss.isDefault || 'ad'
        }
      }
      // 确保giveAddresss是对象格式
      if (this.friendRequest.giveAddresss && this.friendRequest.context) {
        const response = await sendFriendRequest(this.friendRequest)
        if (response.code === 200) {
          this.isFriendModalVisible = false
          Message.success(response.msg)
        } else {
          Message.error(response.msg)
        }
      } else {
        Message.error('请填写完整的请求地址和内容')
      }
      this.select('smallbtn1')
      // activeButton === 'smallbtn1'
    },

    async ProcessingFriendRequests() {
      const response = await ProcessingFriendRequests()
      if (response.code === 200) {
        this.friendRequests = response.data
        if (this.friendRequests && this.friendRequests.length > 0) {
          Message.success('您有新的好友申请')
        }
      } else {
        Message.error(response.msg)
      }
    },
    async acceptFriendRequest(requestId) {
      try {
        const response = await BecomeFriend({ requestId: requestId, isAccepted: 1 })
        if (response.code === 200) {
          this.requestStatus[requestId] = 'accepted'
          Message.success(response.msg)
          this.ProcessingFriendRequests()
        } else {
          Message.error(response.msg)
        }
      } catch (error) {
        // console.error('Accept friend request failed:', error)
        Message.error('接受好友请求失败')
      }
    },
    async rejectFriendRequest(requestId) {
      try {
        const response = await BecomeFriend({ requestId: requestId, isAccepted: 0 })
        if (response.code === 200) {
          this.requestStatus[requestId] = 'rejected'
          Message.success('拒绝好友请求失败')
        } else {
          Message.error(response.msg)
        }
      } catch (error) {
        // console.error('Reject friend request failed:', error)
        Message.error('拒绝好友请求失败')
      }
    }

    // try {
    //   const response = await axios.get(url, { headers });
    //   if (response.data.code === 200) {
    //     this.friendRequests = response.data.data;
    //   } else {
    //     console.error(response.data.msg);
    //   }
    // } catch (error) {
    //   console.error('获取好友申请列表失败:', error);
    // }
  },
  computed: {
    activeButtonClass() {
      return (button) => {
        return this.activeButton === button ? 'active' : ''
      }
    },
    hasUnreadRequests() {
      return this.friendRequests && this.friendRequests.some(request => request.read === false)
    }
    // formattedAddress() {
    //   return {
    //     formattedAddress: this.senderAddress[2],
    //     longitude: this.senderAddress[0],
    //     latitude: this.senderAddress[1]
    //   }
    // }
  },
  // watch: {
  //   // 监听senderAddress的变化，并更新friendRequest.giveAddresss
  //   formattedAddress: {
  //     handler(newVal) {
  //       this.friendRequest.giveAddresss = { ...newVal, isDefault: 'ad' }
  //     },
  //     immediate: true
  //   }
  // },
  mounted() {
    this.showWrite()
    // this.showDriftBottle()// 在组件挂载后获取漂流瓶图片
    this.ProcessingFriendRequests()
    // this.newmessage()
  }
}
</script>
<style scoped>
/* 样式示例 */
.active-button {
  color: #ff6347;
  /* 设置激活状态的字体颜色 */
}

.banner {
  margin-top: 40px;
  width: 1200px;
  height: 730px;
  background-color: rgb(248, 235, 211);
  display: flex;
  /* 使用 Flexbox 布局 */
  align-items: flex-start;
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
  /* 使用 Flexbox 布局 */
  align-items: flex-start;
  /* 垂直对齐子元素 */
  line-height: 40px;
}

.banner .pick {
  position: relative;
  left: 11px;
  top: 10px;
  width: 1300px;
  height: 642px;
  background-image: url(../../assets/imgs/10.jpg);
  background-position: center center;
  background-size: cover;
  border-radius: 25px;
  /* 添加这一行来设置圆角 */
}

.banner .bigbtn {
  position: absolute;
  /* 使用绝对定位 */
  bottom: 60px;
  left: 43%;
  width: 180px;
  height: 70px;
  background: rgba(255, 255, 255, 0.3);
  /* 半透明的白色背景，模拟玻璃 */
  border: 2px solid rgba(255, 255, 255, 0.5);
  /* 边框半透明，增加玻璃感 */
  border-radius: 35px;
  /* 圆角，模拟瓶子形状 */
  backdrop-filter: blur(5px);
  /* 模糊效果，提升玻璃质感 */
  box-shadow: 0px 4px 6px rgba(0, 0, 0, 0.2),
    /* 模拟瓶子的投影 */
    inset 0px 4px 6px rgba(255, 255, 255, 0.4);
  /* 内部的亮光效果，增加立体感 */
  color: #333;
  /* 文本颜色 */
  font-size: 16px;
  font-weight: bold;
  text-align: center;
  line-height: 70px;
  /* 垂直居中 */
  cursor: pointer;
  transition: transform 0.2s ease;
  /* 添加轻微缩放效果 */
  display: flex;
  justify-content: center;
  /* 水平居中 */
  align-items: center;
  /* 垂直居中 */
  font-size: 24px;
  /* 将字体大小设置为 24px */
}

.banner .bigbtn:hover {
  transform: scale(1.10);
  /* 鼠标悬停时轻微放大 */
}

.banner .letter {
  position: relative;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
  /* 组合水平和垂直居中的转换 */
  width: 950px;
  height: 600px;
  background-position: center center;
  background-size: contain;
  background-repeat: no-repeat;
}

.banner .bigbtn1 {
  position: absolute;
  /* 使用绝对定位 */
  top: 380px;
  left: 730px;
  width: 180px;
  height: 70px;
  background: rgba(255, 255, 255, 0.3);
  /* 半透明的白色背景，模拟玻璃 */
  border: 2px solid rgba(255, 255, 255, 0.5);
  /* 边框半透明，增加玻璃感 */
  border-radius: 35px;
  /* 圆角，模拟瓶子形状 */
  backdrop-filter: blur(5px);
  /* 模糊效果，提升玻璃质感 */
  box-shadow: 0px 4px 6px rgba(0, 0, 0, 0.2),
    /* 模拟瓶子的投影 */
    inset 0px 4px 6px rgba(255, 255, 255, 0.4);
  /* 内部的亮光效果，增加立体感 */
  color: #333;
  /* 文本颜色 */
  font-size: 16px;
  font-weight: bold;
  text-align: center;
  line-height: 70px;
  /* 垂直居中 */
  cursor: pointer;
  transition: transform 0.2s ease;
  /* 添加轻微缩放效果 */
  display: flex;
  justify-content: center;
  /* 水平居中 */
  align-items: center;
  /* 垂直居中 */
  font-size: 24px;
  /* 将字体大小设置为 24px */
}

.banner .bigbtn2 {
  position: absolute;
  /* 使用绝对定位 */
  top: 470px;
  left: 730px;
  width: 150px;
  height: 70px;
  background: rgba(255, 255, 255, 0.3);
  /* 半透明的白色背景，模拟玻璃 */
  border: 2px solid rgba(255, 255, 255, 0.5);
  /* 边框半透明，增加玻璃感 */
  border-radius: 35px;
  /* 圆角，模拟瓶子形状 */
  backdrop-filter: blur(5px);
  /* 模糊效果，提升玻璃质感 */
  box-shadow: 0px 4px 6px rgba(0, 0, 0, 0.2),
    /* 模拟瓶子的投影 */
    inset 0px 4px 6px rgba(255, 255, 255, 0.4);
  /* 内部的亮光效果，增加立体感 */
  color: #333;
  /* 文本颜色 */
  font-size: 16px;
  font-weight: bold;
  text-align: center;
  line-height: 70px;
  /* 垂直居中 */
  cursor: pointer;
  transition: transform 0.2s ease;
  /* 添加轻微缩放效果 */
  display: flex;
  justify-content: center;
  /* 水平居中 */
  align-items: center;
  /* 垂直居中 */
  font-size: 24px;
  /* 将字体大小设置为 24px */
}

.banner .smallbtn1,
.banner .smallbtn2,
.banner .smallbtn3 {
  position: absolute;
  width: 120px;
  height: 50px;
  background-color: rgba(199, 43, 43, 0.3);
  /* 默认的半透明红色背景 */
  border: 2px solid rgba(255, 255, 255, 0.5);
  /* 默认的半透明白色边框 */
  border-radius: 25px;
  /* 圆角效果 */
  backdrop-filter: blur(10px);
  /* 模糊效果，增加玻璃感 */
  color: #fff;
  /* 白色文本 */
  font-size: 18px;
  font-weight: bold;
  text-align: center;
  display: flex;
  justify-content: center;
  /* 水平居中 */
  align-items: center;
  /* 垂直居中 */
  right: 50px;
  cursor: pointer;
  transition: background-color 0.3s ease,
    border-color 0.3s ease,
    transform 0.3s ease,
    /* 添加平滑过渡效果 */
    box-shadow 0.3s ease;
  /* 添加阴影过渡效果 */
}

.banner .smallbtn1 {
  top: 50px;
}

.banner .smallbtn2 {
  top: 120px;
}

.banner .smallbtn3 {
  top: 190px;
}

/* 激活状态的样式 */
.banner .smallbtn1.active,
.banner .smallbtn2.active,
.banner .smallbtn3.active {
  background-color: rgba(199, 43, 43, 0.8);
  /* 激活时的背景颜色 */
  border-color: rgba(255, 255, 255, 1);
  /* 激活时的边框颜色 */
  transform: scale(1.1);
  /* 激活时略微放大按钮 */
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
  /* 激活时的阴影效果 */
}

/* 写漂流瓶 */
.write .left {
  position: absolute;
  width: 420px;
  top: 130px;
  left: 5%;
  padding: 20px;
  box-sizing: border-box;
  background-color: rgba(240, 222, 200, 0.88);
  /* 更和谐的背景色 */
  border-radius: 25px;
  box-shadow: 0 6px 18px rgba(0, 0, 0, 0.15);
  /* 柔和的阴影效果 */
  font-family: Arial, sans-serif;
  height: 400px;
  transition: all 0.3s ease;
}

.write .left h2 {
  color: rgb(92, 58, 42);
  /* 更柔和的深色 */
  font-size: 26px;
  margin-bottom: 20px;
  text-align: center;
}

.write .left .map avue-input-map {
  width: 100%;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.write .left .form-group label {
  color: rgb(92, 58, 42);
  font-weight: bold;
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
}

.write .left .form-group textarea {
  width: 100%;
  padding: 12px;
  font-size: 14px;
  border-radius: 8px;
  border: 1px solid rgba(200, 200, 200, 0.5);
  /* 柔和的边框颜色 */
  resize: vertical;
  box-shadow: inset 0 2px 5px rgba(0, 0, 0, 0.05);
  /* 内阴影增加层次感 */
}

.write .left button {
  display: block;
  width: 100%;
  padding: 14px;
  font-size: 16px;
  color: #fff;
  background-image: linear-gradient(135deg, rgba(98, 62, 45, 0.95), rgba(79, 50, 35, 0.95));
  /* 高级渐变色 */
  border: none;
  border-radius: 25px;
  cursor: pointer;
  transition: all 0.3s ease;
  margin-top: 100px;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.15);
}

.write .left button:hover {
  background-image: linear-gradient(135deg, rgba(79, 50, 35, 0.95), rgba(98, 62, 45, 0.95));
  /* 悬停时反向渐变 */
  box-shadow: 0 6px 14px rgba(0, 0, 0, 0.2);
}

.write .right {
  position: absolute;
  left: 52%;
  top: 18px;
  width: 440px;
  height: 642px;
  background-position: center;
  background-size: cover;
  border-radius: 25px;
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.2);
  /* 加深阴影效果 */
  transition: all 0.3s ease;
}

.write .right:hover {
  transform: scale(1.02);
  /* 悬停时略微放大效果 */
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.25);
  /* 增加悬停阴影 */
}

/* end */

h2 {
  height: 30px;
}

.form-group label {
  display: block;
  height: 35px;
}

.form-group input,
.form-group select,
.form-group textarea {
  width: 80%;
  padding: 8px;
  border-radius: 4px;
  resize: none;
  outline: none;
  border: 0;
}

button {
  padding: 10px 20px;
  background-color: #73705d;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  margin-right: 10px;
}

button:hover {
  background-color: #0056b3;
}

/* 模态框背景 */
.friend-modal {
  width: 100%;
  height: 600px;
  display: flex;
  justify-content: center;
  align-items: center;
  text-align: center;
  position: fixed;
  top: 0;
  left: 0;
  z-index: 10000;
  transition: opacity 0.3s ease;
}

/* 模态框内容 */
.modal-content {
  border-radius: 10px;
  width: 35%;
  /* 调整宽度 */
  height: 70%;
  /* 调整高度 */
  max-height: 400px;
  /* 最大高度 */
  background-color: #fff;
  padding: 30px;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.2);
  text-align: center;
  transition: transform 0.3s ease-in-out;
  z-index: 10001;
  /* 确保内容位于背景之上 */
  display: flex;
  flex-direction: column;
  align-items: stretch;
  /* 让子元素扩展至填充整个容器 */
  overflow-y: auto;
  /* 内容超出时显示滚动条 */
}

/* 关闭按钮 */
.close {
  position: absolute;
  /* 绝对定位，使按钮相对于最近的非 static 定位祖先元素定位 */
  top: 70px;
  /* 设置按钮距离顶部的距离为70px */
  right: 290px;
  /* 设置按钮距离右侧的距离为290px */
  font-size: 24px;
  /* 设置字体大小为24px */
  color: #333;
  /* 设置字体颜色为深灰色 */
  cursor: pointer;
  /* 设置鼠标指针为手指形状，表示可点击 */
  transition: color 0.3s ease;
  /* 设置颜色变化的过渡效果，0.3秒内平滑过渡 */
  background-color: transparent;
  /* 设置背景为透明 */
  border: none;
  /* 移除边框 */
  padding: 0;
  /* 移除内边距 */
}

.close:hover {
  color: #f44336;
  /* 鼠标悬停时变红 */
}

/* 标题 */
h2 {
  font-family: 'Roboto', sans-serif;
  font-weight: 600;
  color: #333;
  margin-bottom: 20px;
}

/* 表单 */
form {
  display: flex;
  flex-direction: column;
  flex-grow: 1;
  /* 确保表单内容可以扩展 */
  align-items: center;
  /* 居中对齐表单内容 */
}

/* 地图输入框 */
.map {
  margin-bottom: 20px;
  width: 100%;
  max-width: 500px;
  margin: 0 auto;
}

.avue-input-map {
  width: 100%;
  padding: 12px;
  border-radius: 5px;
  border: 1px solid #ccc;
  font-size: 14px;
  background-color: #fff;
  box-sizing: border-box;
  transition: border-color 0.3s ease;
}

.avue-input-map:focus {
  border-color: #007bff;
  outline: none;
}

/* 文本区域 */
.form-group.special {
  width: 100%;
  max-width: 500px;
  margin: 0 auto;
}

label {
  font-size: 16px;
  color: #333;
  margin-bottom: 10px;
}

textarea {
  width: 100%;
  height: 100px;
  padding: 12px;
  border-radius: 5px;
  border: 1px solid #ccc;
  font-size: 14px;
  background-color: #fff;
  resize: none;
  box-sizing: border-box;
  transition: border-color 0.3s ease;
}

textarea:focus {
  border-color: #007bff;
  outline: none;
}

/* 提交按钮 */
button {
  padding: 12px 24px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  font-size: 16px;
  transition: background-color 0.3s ease;
  margin-top: 20px;
}

button:hover {
  background-color: #0056b3;
  /* 鼠标悬停时变暗 */
}

/* 媒体查询，用于小屏幕设备 */
@media (max-width: 600px) {
  .modal-content {
    width: 95%;
    /* 在更小的屏幕上进一步缩放 */
  }

  button {
    padding: 10px 20px;
    font-size: 14px;
  }
}

button {
  padding: 10px 20px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  font-size: 16px;
  transition: background-color 0.3s ease;
  margin-top: 20px;
}

button:hover {
  background-color: #0056b3;
  /* Darker blue on hover */
}

.form-group {
  height: 80px;
}

.form-group label {
  display: block;
  margin-bottom: 0px;
}

.form-group input,
.form-group select,
.form-group textarea {
  width: 60%;
  border-radius: 4px;
  resize: none;
  outline: none;
  border: 1px solid #ccc;
}

button {
  padding: 10px 20px;
  background-color: #73705d;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  margin-right: 10px;
  margin-top: 50px;
}

button:hover {
  background-color: #35471a;
}

.smallbtn1.active,
.smallbtn2.active,
.smallbtn3.active {
  background-color: #73705d;
  /* 激活状态下的背景颜色 */
  color: white;
}

.red-dot {
  position: absolute;
  top: 0;
  right: 0;
  width: 10px;
  height: 10px;
  background-color: red;
  border-radius: 50%;
}

.request-content {
  /* display: inline-block; */
  position: relative;
}

.friendquest {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 60vw;
  /* 设置宽度为视口宽度的 60% */
  height: 80vh;
  /* 设置高度为视口高度的 80% */
  max-width: 750px;
  /* 限制最大宽度 */
  max-height: 500px;
  /* 限制最大高度 */
  padding: 25px;
  background: linear-gradient(135deg, #f3f4f6, #e5e7eb);
  border-radius: 12px;
  box-shadow: 0px 4px 15px rgba(0, 0, 0, 0.2);
  overflow-y: auto;
}

.request-content {
  position: relative;
  width: 700px;
  padding: 15px;
  background-color: #ffffff;
  border-radius: 8px;
  box-shadow: 0px 2px 8px rgba(0, 0, 0, 0.15);
  display: flex;
  align-items: center;
  margin-bottom: 20px;
  transition: transform 0.2s ease;
}

.request-content:hover {
  transform: translateY(-5px);
}

.content {
  flex-grow: 1;
  display: flex;
  align-items: center;
  padding-left: 60px;
  position: relative;
}

.avatar {
  position: absolute;
  left: 5px;
  width: 45px;
  height: 45px;
  border-radius: 50%;
  box-shadow: 0px 2px 4px rgba(0, 0, 0, 0.2);
}

.username {
  font-weight: 600;
  color: #333;
  font-size: 16px;
  margin-right: 20px;
}

.context {
  color: #555;
  font-size: 15px;
}

.time {
  font-size: 13px;
  color: #999;
  position: absolute;
  bottom: 5px;
  right: 0;
}

.button1,
.button2 {
  padding: 8px 16px;
  font-size: 14px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: background-color 0.3s ease, transform 0.2s ease;
}

.button1 {
  background-color: #34d399;
  color: white;
  margin-right: 10px;
}

.button1:hover {
  background-color: #10b981;
  transform: scale(1.05);
}

.button2 {
  background-color: #f87171;
  color: white;
}

.button2:hover {
  background-color: #ef4444;
  transform: scale(1.05);
}

.no-requests {
  text-align: center;
  padding: 20px;
  font-size: 16px;
  color: #888;
}
</style>
