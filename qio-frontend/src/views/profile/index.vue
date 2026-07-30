<template>
  <div class="container">
    <div class="banner">
      <h1 class="profile"><img src="../../assets/imgss/profile.webp" alt=""></h1>
      <div class="data">

        <div class="pic" @click="getAvatarList">
          <img :src="user.avatar" alt="">
          <div class="change-text">
            更换头像
          </div>
        </div>
        <!-- 头像列表弹窗 -->
        <el-dialog title="选择头像" :visible.sync="showAvatarModal" width="50%">
          <div class="avatar-list" style="height: 300px;">
            <div class="avatar-item" v-for="avatar in avatars" :key="avatar.id" @click="selectAvatar(avatar)">
              <img :src="avatar.url" alt="" style="width:100px;height:100px;">
            </div>
          </div>
        </el-dialog>
        <div class="information">
          <div class="name">{{ user.nickname }}<el-link icon="el-icon-edit" class="edit"
              @click="showEditNicknameDialog">编辑</el-link></div>
          <!-- 修改昵称弹窗 -->
          <el-dialog title="修改昵称" :visible.sync="editNicknameDialogVisible" width="30%"
            @close="editNicknameDialogVisible = false">
            <el-form ref="nicknameForm" label-width="80px">
              <el-form-item label="昵称">
                <el-input v-model="newNickname" autocomplete="off"></el-input>
              </el-form-item>
            </el-form>
            <span slot="footer" class="dialog-footer">
              <el-button @click="editNicknameDialogVisible = false" class="cancel">取消</el-button>
              <el-button type="primary" @click="updateNickname" style="background-color:#888B69;">确定</el-button>
            </span>
          </el-dialog>
          <div class="userid">用户名:<span>{{ user.username }}</span><el-link icon="el-icon-edit" class="edit"
              @click="showEditUsernameDialog">编辑</el-link></div>
          <el-dialog title="修改用户名" :visible.sync="editUsernameDialogVisible" width="30%"
            @close="editUsernameDialogVisible = false">
            <el-form ref="usernameForm" label-width="80px">
              <el-form-item label="用户名">
                <el-input v-model="newUsername" autocomplete="off"></el-input>
              </el-form-item>
            </el-form>
            <span slot="footer" class="dialog-footer">
              <el-button @click="editUsernameDialogVisible = false" class="cancel">取消</el-button>
              <el-button type="primary" @click="updateUsername" style="background-color:#888B69;">确定</el-button>
            </span>
          </el-dialog>
          <div class="sex">性别:<span>{{ user.sex }}</span><el-link icon="el-icon-edit" class="edit"
              @click="showEditSexDialog">编辑</el-link></div>
          <!-- 修改性别弹窗 -->
          <el-dialog title="修改性别" :visible.sync="editSexDialogVisible" width="30%"
            @close="editSexDialogVisible = false">
            <el-form ref="sexForm" label-width="80px">
              <el-form-item label="性别">
                <el-select v-model="newSex" placeholder="请选择">
                  <el-option label="男" value="男"></el-option>
                  <el-option label="女" value="女"></el-option>
                </el-select>
              </el-form-item>
            </el-form>
            <span slot="footer" class="dialog-footer">
              <el-button @click="editSexDialogVisible = false" class="cancel">取消</el-button>
              <el-button type="primary" @click="updateSex" style="background-color:#888B69;">确定</el-button>
            </span>
          </el-dialog>
          <div class="email">邮箱:<span>{{ user.email }}</span></div>
          <el-dialog title="我的地址" :visible.sync="showAddressesModal" width="720px" append-to-body>
            <div class="AddressDialog">
              <div class="addressContent">
                <div class="UserAddress beautify-scroll-bar ">
                  <div class="myAddressItem" v-for="address in addresses" :key="address.id">
                    <div class="defaultTip" v-if="address.isDefault == 'true'">默认</div>
                    <div class="setDefaultTip" @click="setUserDefaultAddress(address.id)" v-else>设为默认</div>
                    <i class="el-icon-delete closeIcon" v-if="address.isDefault !== 'true'"
                      @click="deleteMyAddress(address.id)"></i>
                    <div class="recipientInfo">
                      <div class="recipient">{{ user.nickname }}</div>
                      <div class="addressCountry">{{ getCountyNameById(address.countryId) }}</div>
                    </div>
                    <div class="detailAddress f-els-2">{{ address.formattedAddress }}</div>
                  </div>
                </div>
              </div>
            </div>
          </el-dialog>
          <el-dialog title="我的好友" :visible.sync="showFriendsModal" width="520px" append-to-body>
            <div class="friends-container beautify-scroll-bar">
              <div class="friend" v-for="friend in friends" :key="friend.id">
                <div style="float: left;">
                  <el-avatar :src="friend.avatar" shape="square" fit="fill"></el-avatar>
                </div>
                <div style="float: left; margin-left: 10px;">
                  <div class="friendname editFriendName" v-if="!friend.isEditing">
                    {{ friend.remark || friend.name }}
                    <i v-if="friend.sex === '男'" style="color:cornflowerblue; font-size: 16px;"
                      class="el-icon-male"></i>
                    <i v-if="friend.sex === '女'" style="color: pink; font-size: 16px;" class="el-icon-female"></i>
                    <i v-if="!isEditingAnyFriend" class="el-icon-edit-outline editRemarkIcon"
                      @click="editFriendRemark(friend.id)"></i>
                  </div>
                  <div class="friendname" v-else>
                    <el-input v-model="newRemark" style="height: 20px;width: 150px;" maxlength="10"
                      show-word-limit></el-input>
                    <el-button @click="saveFriendRemark(friend.id)" size="mini"
                      style="height: 20px; line-height: 7px; border-radius: 5px;  padding: 0 auto !important; margin-left: 10px; background-color:  #ff6200; color: #fff;"
                      round>提交</el-button>
                    <el-button @click="cancelEdit(friend.id)" size="mini"
                      style="height: 20px; line-height: 7px; border-radius: 5px; padding: 0 auto !important;"
                      round>取消</el-button>
                  </div>
                  <div class="friendemail" v-show="friend.remark"><span
                      style="color: #666; font-size: 14px;">昵称：</span>{{ friend.name }}
                  </div>
                  <div class="friendemail"><span style="color: #666; font-size: 14px;">邮箱号：</span>{{ friend.email }}
                  </div>
                </div>
                <el-popover placement="right" trigger="click" width="250">
                  <div class="AddressDialog">
                    <div class="addressTip">好友地址</div>
                    <div class="addressContent">
                      <div class="UserAddress beautify-scroll-bar ">
                        <div class="friendAddressItem" v-for="address in friend.addresses" :key="address.id">
                          <div class="defaultTip" v-if="address.isDefault == 'true'">默认</div>
                          <div class="setDefaultTip" @click="setFriendDefaultAddress(friend.id, address.id)" v-else>设为默认
                          </div>
                          <i class="el-icon-delete closeIcon" v-if="address.isDefault !== 'true'"
                            @click="deleteFriendAddress(friend.id, address.id)"></i>
                          <div class="recipientInfo">
                            <div class="recipient">{{ friend.remark || friend.name }}</div>
                            <div class="addressCountry">{{ getCountyNameById(address.countryId) }}</div>
                          </div>
                          <div class="detailAddress f-els-2">{{ address.formattedAddress }}</div>
                        </div>

                      </div>
                    </div>
                  </div>
                  <i slot="reference" v-show="friend.addresses.length > 0" class="el-icon-s-home closeIcon"></i>
                </el-popover>
              </div>
            </div>
          </el-dialog>

        </div>
        <a @click="showAddressesModal = true" class="address">我的地址</a>
        <a @click="toggleFriendsModal" class="friends">我的好友</a>
        <a class="code" @click="showEditPasswordDialog">修改密码</a>
        <el-dialog title="修改密码" :visible.sync="showPasswordModal" width="30%" @close="showPasswordModal = false">
          <el-form ref="passwordForm" :model="passwordForm" label-width="80px">
            <el-form-item label="原密码">
              <el-input type="password" v-model="passwordForm.oldPassword" autocomplete="off"></el-input>
            </el-form-item>
            <el-form-item label="新密码">
              <el-input type="password" v-model="passwordForm.newPassword" autocomplete="off"></el-input>
            </el-form-item>
            <el-form-item label="确认密码">
              <el-input type="password" v-model="passwordForm.confirmPassword" autocomplete="off"></el-input>
            </el-form-item>
          </el-form>
          <span slot="footer" class="dialog-footer">
            <el-button @click="showPasswordModal = false" class="cancel">取消</el-button>
            <el-button type="primary" @click="updatePassword" style="background-color:#888B69;">确定</el-button>
          </span>
        </el-dialog>
      </div>
      <div class="store">
        <img src="../../assets/imgss/store.webp" alt="仓库" class="bgd">
        <img src="../../assets/imgss/color.webp" alt="墨水" class="ink" @click="showInkImage">
        <img src="../../assets/imgss/pen.webp" alt="笔" class="pen" @click="showWordImage">
        <img src="../../assets/imgss/envelope.webp" alt="信纸" class="envelope" @click="showEnvelopeImage">
        <div v-if="showInkImageModal" class="image-modal">
          <div class="colorbox">
            <div class="box" v-for="color in inkColors" :key="color.id"
              :style="{ backgroundImage: 'url(' + color.previewImage + ')' }">
            </div>
          </div>
          <span class="close" @click="showInkImageModal = false">&times;</span>
          <img src="../../assets/imgss/mood.webp" alt="墨水图片" class="modal-image" />
        </div>
        <!-- 字体展示台 -->
        <div v-if="showWordModal" class="image-modal">
          <div class="penboxs">
            <div class="penbox" v-for="font in wordFonts" :key="font.id"
              :style="{ backgroundImage: 'url(' + font.previewImage + ')' }">
            </div>
          </div>
          <span class="close" @click="showWordModal = false">&times;</span>
          <img src="../../assets/imgss/mood.webp" alt="墨水图片" class="modal-image" />
        </div>
        <!-- 信纸展示台 -->
        <div v-if="showEnvelopeModal" class="image-modal">
          <div class="penboxs">
            <div class="envelopebox" v-for="paper in paperPreviews" :key="paper.id"
              :style="{ backgroundImage: 'url(' + paper.previewImage + ')' }">
            </div>
          </div>
          <span class="close" @click="showEnvelopeModal = false">&times;</span>
          <img src="../../assets/imgss/mood.webp" alt="墨水图片" class="modal-image" />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { Message } from 'element-ui'
import { updateFriendRemark, getUserFriends, setFriendDefaultAddress, deleteFriendAddress, setUserDefaultAddress, deleteUserAddress, getMyAddress, updateUsername, updateNickname, updateSex, getAvatarList, updatePassword } from '@/api/user'
import axios from 'axios'
import useUserStore from '@/store/modules/user'

export default {
  name: 'UserProfile',
  data() {
    return {
      showPasswordModal: false, // 控制修改密码弹窗的显示
      passwordForm: {
        oldPassword: '',
        newPassword: '',
        confirmPassword: ''
      },
      countries: [],
      avatars: [], // 存储头像列表
      showAvatarModal: false, // 控制头像列表弹窗的显示
      editSexDialogVisible: false, // 控制性别编辑弹窗显示的状态
      editNicknameDialogVisible: false, // 控制昵称编辑弹窗显示的状态
      editUsernameDialogVisible: false,
      showAddressesModal: false, // 控制地址列表弹窗显示的状态
      addresses: [],
      showFriendsModal: false,
      friends: [],
      isEditingAnyFriend: false,
      wordFonts: [], // 存储字体信息
      inkColors: [], // 存储墨水颜色的 hexCode
      paperPreviews: [], // 存储信纸预览图信息
      showInkImageModal: false,
      showWordModal: false,
      showEnvelopeModal: false,
      user: {},
      newRemark: '',
      newUsername: '',
      newNickname: '',
      newSex: '',
      money: '',
      repository: {
        fonts: [],
        papers: []
      },
      editModalVisible: false // 控制弹窗显示的状态
    }
  },
  async created() {
    try {
      const userStore = useUserStore()
      const token = userStore.token // 从 Pinia store 获取 token

      // 获取用户信息
      const userResponse = await this.getUserInfo()
      if (userResponse && userResponse.data.code === 200) {
        this.user = userResponse.data.data
      }

      // 获取用户的钱数
      const moneyResponse = await this.getUserMoney(token)
      if (moneyResponse && moneyResponse.data.code === 200) {
        this.money = moneyResponse.data.data.money
      }

      // 获取用户仓库信息
      const repoResponse = await axios.get('/api/user/getUserRepository', {
        headers: { Authorization: token }
      })
      if (repoResponse.data.code === 200 && repoResponse.data.data) {
        this.repository = repoResponse.data.data
        this.inkColors = repoResponse.data.data.fontColors // 存储墨水颜色信息
        this.wordFonts = repoResponse.data.data.fonts // 存储字体信息
        this.paperPreviews = repoResponse.data.data.papers // 存储信纸预览图信息
      }

      // 获取好友列表
      await this.getMyFriends(token)

      // 获取好友地址
      // await this.updateFriendsAddresses(token)

      // 获取我的地址
      await this.fetchAddresses(token)
    } catch (error) {
      // console.error('请求错误:', error)
    }
  },
  methods: {
    loadRouterParams() {
      const command = this.$route.query.command
      if (command) {
        if (command === 'friends') {
          this.showFriendsModal = true
        } else if (command === 'address') {
          this.showAddressesModal = true
        } else if (command === 'cangku') {
          this.showEnvelopeImage()
        } else if (command === 'password') {
          this.showEditPasswordDialog()
        }
      }
      setTimeout(() => {
        this.$router.push({ query: null }) // 清除 URL 中的参数
      }, 500)
    },
    editFriendRemark(friendId) {
      this.isEditingAnyFriend = true
      const friend = this.friends.find(f => f.id === friendId)
      if (friend) {
        friend.isEditing = true
        this.newRemark = friend.remark || friend.name
        // this.dialogKey += 1 // 强制重新渲染
        this.$forceUpdate() // 强制刷新组件
      }
    },
    saveFriendRemark(friendId) {
      const friend = this.friends.find(f => f.id === friendId)
      if (friend) {
        updateFriendRemark(friendId, this.newRemark).then(res => {
          friend.remark = this.newRemark
          getUserFriends().then(res => {
            this.friends = res.data
          })
          this.$message.success(res.msg)
        })
        friend.isEditing = false
        this.$forceUpdate() // 强制刷新组件
      }
      this.isEditingAnyFriend = false
    },
    cancelEdit(friendId) {
      const friend = this.friends.find(f => f.id === friendId)
      if (friend) {
        friend.isEditing = false
        this.dialogKey += 1 // 强制重新渲染
        friend.id = friendId
        // this.dialogKey += 1 // 强制重新渲染
        this.$forceUpdate() // 强制刷新组件
      }
      this.isEditingAnyFriend = false
    },
    setFriendDefaultAddress(friendId, addressId) {
      setFriendDefaultAddress(friendId, addressId).then(res => {
        getUserFriends().then(res => {
          this.friends = res.data
        })
        this.$message.success(res.msg)
      })
    },
    deleteFriendAddress(friendId, addressId) {
      deleteFriendAddress(friendId, addressId).then(res => {
        getUserFriends().then(res => {
          this.friends = res.data
        })
        this.$message.success(res.msg)
      })
    },
    setUserDefaultAddress(id) {
      setUserDefaultAddress(id).then(res => {
        this.fetchAddresses()
        this.$message.success(res.msg)
      })
    },
    deleteMyAddress(id) {
      deleteUserAddress(id).then(res => {
        this.fetchAddresses()
        this.$message.success(res.msg)
      })
    },
    getCountyNameById(id) {
      // console.log(id)
      for (const country of this.countries) {
        if (country.id === id) {
          return country.countryName
        }
      }
      return '中国'
    },
    showEditSexDialog() {
      this.editSexDialogVisible = true
    },
    toggleFriendsModal() {
      this.showFriendsModal = !this.showFriendsModal
    },
    showInkImage() {
      this.showInkImageModal = true // 显示墨水图片
    },
    showWordImage() {
      this.showWordModal = true // 显示墨水图片
    },
    showEnvelopeImage() {
      this.showEnvelopeModal = true // 显示墨水图片
    },
    async getMyFriends(token) {
      try {
        const response = await axios.get('/api/user/getMyFriends', {
          headers: {
            Authorization: token
          }
        })
        if (response.data.code === 200) {
          this.friends = response.data.data
        } else {
          // console.error('获取好友列表失败:', response.data.msg)
        }
      } catch (error) {
        // console.error('获取好友列表错误:', error)
      }
    },
    async getUserInfo() {
      const userStore = useUserStore()
      userStore.getUserInfo().then(res => {
        this.user = res.data
      })
      // const response = await axios.get('/api/user/getUserInfo', {
      //   headers: {
      //     Authorization: token
      //   }
      // })
      // return response
      // const userStore = useUserStore()
    },
    async getUserMoney(token) {
      const response = await axios.get('/api/user/getUserMoney', {
        headers: {
          Authorization: token
        }
      })
      return response
    },
    async getUserRepository(token) {
      const response = await axios.get('/api/user/getUserRepository', {
        headers: {
          Authorization: token
        }
      })
      return response
    },
    showEditModal() {
      this.editModalVisible = true
    },
    hideEditModal() {
      this.editModalVisible = false
    },
    showEditNicknameDialog() {
      this.editNicknameDialogVisible = true
    },
    async updateNickname() {
      await updateNickname(this.newNickname).then(res => {
        this.editNicknameDialogVisible = false
        Message.success('修改昵称成功')
        // this.user.nickname = res.data.nickname // 更新用户名
        this.getUserInfo()
      })
    },
    async updateSex() {
      await updateSex(this.newSex).then(res => {
        this.editSexDialogVisible = false
        Message.success('修改性别成功')
        // this.user.sex = res.data.sex // 更新用户名
        this.getUserInfo()
      })
    },
    fetchAddresses() {
      // 获取我的地址
      getMyAddress().then(res => {
        this.addresses = res.data
      })
    },
    // async getFriendAddress(friendId, token) {
    //   try {
    //     const response = await axios.get('/api/user/getFriendAddress', {
    //       params: { friendId },
    //       headers: { Authorization: token }
    //     })
    //     if (response.data.code === 200 && response.data.data.length > 0) {
    //       return response.data.data[0].formattedAddress
    //     } else {
    //       // console.error('获取好友地址失败:', response.data.msg)
    //       return null
    //     }
    //   } catch (error) {
    //     // console.error('获取好友地址错误:', error)
    //     return null
    //   }
    // },
    // async updateFriendsAddresses(token) {
    //   for (const friend of this.friends) {
    //     const address = await this.getFriendAddress(friend.id, token)
    //     if (address) {
    //       this.$set(friend, 'address', address)
    //     }
    //   }
    // },
    showEditUsernameDialog() {
      this.editUsernameDialogVisible = true
    },
    // 更新用户名
    async updateUsername() {
      await updateUsername(this.newUsername).then(res => {
        this.editUsernameDialogVisible = false
        Message.success('修改用户名成功')
        // this.user.username = res.data.username // 更新用户名
        this.getUserInfo()
      })
    },
    async getAvatarList() {
      await getAvatarList().then(res => {
        this.showAvatarModal = true // 显示头像列表弹窗
        this.avatars = res.data // 存储头像列表
      })
    },
    async selectAvatar(avatar) {
      // 更新用户头像 ID
      this.user.avatarId = avatar.id
      // 发送请求到后端，更新用户信息
      await this.updateUserInfo()
      // 关闭头像列表弹窗
      this.showAvatarModal = false
      // 立即更新页面上的头像显示
      // this.user.avatar = avatar.url
      this.getUserInfo()
    },
    async updateUserInfo() {
      try {
        const token = useUserStore().token // 从 Pinia store 获取 token
        const response = await axios.put('/api/user/updateUserInfo', {
          nickname: this.user.nickname,
          avatarId: this.user.avatarId,
          sex: this.user.sex
        }, {
          headers: { Authorization: token }
        })
        if (response.data.code === 200) {
          Message.success(response.data.msg)
        } else {
          Message.error(response.data.msg || '修改用户信息失败')
        }
      } catch (error) {
        // console.error('更新用户信息错误:', error)
        Message.error('更新用户信息时发生错误')
      }
    },
    showEditPasswordDialog() {
      this.showPasswordModal = true
    },
    async updatePassword() {
      if (this.passwordForm.newPassword !== this.passwordForm.confirmPassword) {
        Message.error('新密码和确认密码不一致')
        return
      }
      const userStore = useUserStore()
      const response = await updatePassword(
        this.passwordForm.oldPassword,
        this.passwordForm.newPassword,
        this.passwordForm.confirmPassword,
        userStore.token
      )
      if (response.code === 200) {
        this.showPasswordModal = false
        Message.success(response.msg)
      } else {
        Message.error(response.data.msg || '修改密码失败')
      }
    }
  },
  mounted() {
    this.loadRouterParams()
  },
  beforeRouteUpdate(to, from, next) {
    this.loadRouterParams()
    next()
  }
}
</script>

<style scoped>
/* .container {
  position: relative;
  top: -60px;
  padding-top: 120px;
  box-sizing: border-box;
  min-height: 100vh;
  width: 100%;
  background-image: url(../../assets/imgss/background.jpg);
  background-size: 100% 100%;
  background-position: top center;
} */
.UserAddress {
  align-self: flex-start;
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  max-height: 300px;
  overflow: auto;
  width: 100%
}

.beautify-scroll-bar {
  scrollbar-color: #e6e9eb;
  scrollbar-width: 4px
}

.beautify-scroll-bar::-webkit-scrollbar-thumb {
  background: #e6e9eb;
  border-radius: 4px;
  height: 4px;
  width: 4px
}

.UserAddress.scrollHide {
  scrollbar-width: none
}

/* 滑动轨道按钮 */
::-webkit-scrollbar-button {
  width: 10px;
  height: 10px;
  display: none;
}

.UserAddress .myAddressItem {
  background-color: #f3f6f8;
  border: 1px solid transparent;
  border-radius: 8px;
  cursor: pointer;
  font-family: PingFangSC-Regular;
  font-size: 16px;
  font-weight: 400;
  /* height: 80px; */
  margin-right: 12px;
  padding: 12px 12px 9px;
  position: relative;
  width: 300px;
}

.closeIcon {
  cursor: pointer;
  display: inline-block;
  font-size: 18px;
  line-height: 18px;
  position: absolute;
  right: 10px;
  top: 34px
}

.UserAddress .myAddressItem:nth-of-type(2n) {
  margin-right: 0
}

.UserAddress .myAddressItem:nth-of-type(n+3) {
  margin-top: 12px
}

.UserAddress .myAddressItem.isSelected,
.UserAddress .myAddressItem:hover {
  background-color: #fff;
  border: 1px solid #ff6200
}

.UserAddress .friendAddressItem {
  background-color: #f3f6f8;
  border: 1px solid transparent;
  border-radius: 8px;
  cursor: pointer;
  font-family: PingFangSC-Regular;
  font-size: 12px;
  font-weight: 400;
  /* height: 80px; */
  margin-right: 12px;
  padding: 12px 12px 9px;
  position: relative;
  width: 200px;
}

.UserAddress .friendAddressItem:nth-of-type(2n) {
  margin-top: 12px
}

.UserAddress .friendAddressItem:nth-of-type(n+3) {
  margin-top: 12px
}

.UserAddress .friendAddressItem.isSelected,
.UserAddress .friendAddressItem:hover {
  background-color: #fff;
  border: 1px solid #ff6200
}

.UserAddress .defaultTip {
  background-color: #ff6200;
  border-radius: 4px;
  color: #fff;
  font-size: 12px;
  line-height: 12px;
  padding: 4px;
  position: absolute;
  right: 4px;
  top: 4px
}

.UserAddress .setDefaultTip {
  background-color: rgba(0, 0, 0, .2);
  border-radius: 4px;
  color: #fff;
  font-size: 12px;
  line-height: 12px;
  padding: 4px;
  position: absolute;
  right: 4px;
  top: 4px
}

.UserAddress .detailAddress {
  color: #50607a;
  height: 36px;
  line-height: 18px
}

.UserAddress .recipientInfo {
  display: flex;
  flex-direction: row;
  line-height: 12px;
  margin-bottom: 8px
}

.UserAddress .recipientInfo .recipient {
  font-family: PingFangSC-Semibold;
  font-weight: 600;
  margin-right: 8px
}

.UserAddress .recipientInfo .addressCountry {
  color: #50607a;
  font-family: PingFangSC-Regular;
  font-weight: 400
}

.beautify-scroll-bar::-webkit-scrollbar {
  height: 4px;
  margin-right: 4px;
  width: 4px
}

.f-els-1 {
  white-space: nowrap;
  word-break: keep-all
}

.f-els-1,
.f-els-2 {
  overflow: hidden;
  text-overflow: ellipsis
}

.f-els-2 {
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  display: -webkit-box
}

.AddressDialog {
  display: flex;
  flex-direction: column;
  height: 100%;
  position: relative
}

.addressTip {
  color: #11192d;
  font-family: PingFangSC-Semibold;
  font-size: 14px;
  font-weight: 600;
  line-height: 14px;
  margin-bottom: 12px;
  margin-left: 4px
}

.addressContent {
  flex: 1;
  align-items: center;
  display: flex;
  flex-direction: column;
  justify-content: center
}

.address-aside {
  float: right;
  background-color: #A52328;
  color: #fff;
  width: 40px;
  border-radius: 8px;
  border: 1px solid #ff6200;
  text-align: center;
  font-size: 20px;
  line-height: 40px;
}

::v-deep .el-dialog {
  margin-top: 200px !important;
  border-radius: 10px;
  /* background-color: rgba(222, 201, 162, 1) !important; */
  background: 0 0 / 400px auto repeat, #f9f9f9;
}

.banner {
  position: relative;
  margin: 0 auto;
  width: 1300px;
  height: 800px;
}

.banner .profile {
  position: absolute;
  top: -20px;
  left: 35%;
  transform: translateY(-50%);
  width: 400px;
  height: 100px;
}

.banner h1 img {
  width: 100%;
  height: 100%;
}

.data {
  position: relative;
  top: 60px;
  left: 50px;
  width: 350px;
  height: 640px;
  padding-top: 30px;
  box-sizing: border-box;
  margin-top: 60px;
  background-image: url(../../assets/imgss/profilebgd.webp);
  background-position: center center;
  background-size: cover;
}

.data .pic {
  position: absolute;
  top: 50px;
  left: 50%;
  border-radius: 5px;
  transform: translate(-50%, 0);
  width: 200px;
  height: 200px;
  cursor: pointer;
}

.pic:hover .change-text {
  opacity: 1;
}

.change-text {
  position: absolute;
  top: 50%;
  left: 50%;
  width: 200px;
  height: 200px;
  line-height: 200px;
  transform: translate(-50%, -50%);
  background-color: rgba(0, 0, 0, 0.5);
  color: white;
  padding: 5px;
  border-radius: 5px;
  opacity: 0;
  transition: opacity 0.3s ease;
  z-index: 1;
}

.data .pic img {
  width: 100%;
  height: 100%;
  border-radius: 5px;
}

.data .name {
  position: absolute;
  top: 10px;
  left: 20px;
  font-size: 42px;
  font-weight: bold;
  font-family: '华文新魏', sans-serif;
}

.name:hover .edit {
  display: inline-block;
}

.data .userid {
  position: absolute;
  top: 70px;
  left: 20px;
  font-size: 20px;
  font-weight: bold;
  font-family: '华文中宋', sans-serif;
}

.userid span {
  font-size: 14px;
  font-weight: normal;
  margin-left: 2px;
  font-family: Arial, sans-serif;
}

.userid:hover .edit {
  display: inline-block;
}

.edit {
  display: none;
}

.data .sex {
  position: absolute;
  top: 110px;
  left: 20px;
  font-size: 20px;
  font-weight: bold;
  font-family: '华文中宋', sans-serif;
}

.sex:hover .edit {
  display: inline-block;
}

.data .email {
  position: absolute;
  top: 150px;
  left: 20px;
  font-size: 20px;
  font-weight: bold;
  font-family: '华文中宋', sans-serif;
}

.email span {
  /* position: absolute;
  left: 65px;
  top: 8px; */
  font-size: 16px;
  font-weight: normal;
  font-family: Arial, sans-serif;
}

.sex span {
  font-size: 16px;
  font-weight: normal;
  font-family: Arial, sans-serif;
}

.store {
  position: absolute;
  left: 396px;
  top: 60px;
  width: 850px;
  height: 640px;
  background-image: url(../../assets/imgss/frame.webp);
  background-position: center center;
  background-size: cover;
}

.store .bgd {
  position: absolute;
  top: 35px;
  left: 40px;
  width: 768px;
  height: 578px;
}

.data .information {
  position: relative;
  top: 240px;
  left: 25px;
  width: 300px;
  height: 300px;
  background-color: #ADA587;
}

.ink {
  position: absolute;
  bottom: 180px;
  left: 200px;
  cursor: pointer;
  animation: float 2.7s ease-in-out infinite;
}

.ink:hover {
  animation: big 0.5s linear forwards;
}

.pen:hover {
  animation: big 0.5s linear forwards;
}

.envelope:hover {
  animation: big 0.5s linear forwards;
}

.pen {
  position: absolute;
  bottom: 180px;
  left: 330px;
  cursor: pointer;
  animation: float 3s ease-in-out infinite;
}

.envelope {
  position: absolute;
  width: 100px;
  height: 70px;
  bottom: 180px;
  left: 450px;
  cursor: pointer;
  animation: float 2.5s ease-in-out infinite;
}

@keyframes big {

  0% {
    transform: scale(1);
  }

  100% {
    transform: scale(1.2);
  }
}

@keyframes float {

  0%,
  100% {
    transform: translateX(0);
  }

  50% {
    transform: translateX(-5px) translateY(-1px);
  }
}

.data .showEdit {
  position: absolute;
  bottom: 85px;
  left: 50px;
  cursor: pointer;
}

.edit-modal {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 100%;
  max-width: 400px;
  background-image: url(../../assets/imgss/profilebgd.webp);
  border-radius: 25px;
  z-index: 1000;
  padding: 20px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
}

ul {
  list-style: none;
  /* 移除列表项的默认样式 */
  padding: 0;
  /* 移除默认的内边距 */
  list-style: none;
  /* 移除列表项的默认样式 */
  padding: 0;
  /* 移除默认的内边距 */
}

li {
  margin-bottom: 10px;
  /* 列表项之间的间距 */
  margin-bottom: 10px;
  /* 列表项之间的间距 */
}

.close {
  position: absolute;
  top: -40px;
  right: -10px;
  cursor: pointer;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
}

.form-group input,
.form-group select {
  width: 80%;
  padding: 8px;
  border: 1px solid rgb(194, 179, 153);
  border-radius: 4px;
}

.save {
  width: 100px;
  height: 50px;
  background-color: #c0c0a8;
  border: 0;
  border-radius: 25px;
  width: 100px;
  height: 50px;
  background-color: #c0c0a8;
  border: 0;
  border-radius: 25px;
}

.save:hover {
  color: white;
  background-color: #6f6f5e;
  ;
  background-color: #6f6f5e;
  ;
  cursor: pointer;
}

.image-modal {
  position: fixed;
  top: 50%;
  left: 50%;
  width: 500px;
  height: 500px;
  transform: translate(-50%, -50%);
  z-index: 1000;
}

.image-modal span {
  position: absolute;
  top: 20px;
  right: 40px;
}

.modal-image {
  width: 100%;
  width: 100%;
}

.colorbox {
  position: absolute;
  top: 50px;
  left: 80px;
  z-index: 1;
}

.box {
  float: left;
  width: 105px;
  height: 50px;
  background-size: cover;
  border: 1px solid #ccc;
  /* 可选的边框样式 */
  margin-right: 12px;
  margin-bottom: 5px;
  border-radius: 5px;
}

.penboxs {
  position: absolute;
  top: 40px;
  left: 40px;
  z-index: 1;
}

.penbox {
  float: left;
  width: 130px;
  height: 60px;
  background-size: cover;
  margin-right: 12px;
  margin-bottom: 5px;
  border-radius: 5px;
}

.envelopebox {
  float: left;
  width: 185px;
  height: 120px;
  border: 1px solid #000000;
  background-position: center center;
  background-size: cover;
  margin-right: 12px;
  margin-bottom: 5px;
  border-radius: 5px;
  z-index: 1;
}

.modal-content {
  position: relative;
}

.close {
  position: absolute;
  top: -40px;
  right: -10px;
  cursor: pointer;
  z-index: 1;
}

.friends-container {
  align-self: flex-start;
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  max-height: 300px;
  overflow: auto;
  width: 100%
}

.friends {
  position: absolute;
  top: 500px;
  left: 230px;
  cursor: pointer;
  color: #565648;
}

.friend {
  position: relative;
  width: 80%;
  min-height: 50px;
  max-height: 70px;
  text-align: left;
  margin-left: 30px;
  margin-bottom: 10px;
  padding: 10px;
  padding-top: 20px;
  /* background-color: #c0c0a8; */
  border-radius: 5px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
}

.friendname ::v-deep .el-input__inner {
  font-size: 12px;
  height: 20px !important;
  margin-bottom: 2px;
}

.editFriendName .editRemarkIcon {
  display: none;
}

.editFriendName:hover .editRemarkIcon {
  display: inline-block;
}

.friendname,
.friendsex,
.friendemail {
  font-size: 15px;
  margin-bottom: 5px;
}

.address {
  position: absolute;
  top: 500px;
  left: 140px;
  font-size: 14px;
  cursor: pointer;
  color: #565648;
}

.data a {
  font-size: 15px;
  font-size: 15px;
}

.cancel:hover {
  border: #6f6f5e 1px solid;
  background-color: #e2ddcb;
  color: #6f6f5e;
}

.el-link:hover {
  color: #6f6f5e;
}

.code {
  position: absolute;
  top: 500px;
  left: 60px;
  font-size: 14px;
  cursor: pointer;
  color: #565648;
}

.avatar-list {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  justify-content: center;
}

.avatar-item {
  cursor: pointer;
}
</style>
