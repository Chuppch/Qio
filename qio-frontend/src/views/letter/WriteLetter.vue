<template>
  <div class="banner">

    <div class="selectPage" @click="showSend" :class="{ 'active': currentView === 'send' }">
      <el-select v-model="currentView" style="width: 90px;" @change="handleChangePage">
        <el-option v-for="item in pageList" :key="item.value" :label="item.name" :value="item.value" />
      </el-select>
    </div>
    <div v-if="currentView === 'write'">
      <div class="left">
        <div class="c-toggle" @click="handleClickTab">
          <button class="c-toggle__item" :class="{ 'c-toggle__item--checked': activeTab === 'information' }"
            @click.stop="activeTab = 'information'">
            信件信息
          </button>
          <button class="c-toggle__item" :class="{ 'c-toggle__item--checked': activeTab === 'letter' }"
            @click.stop="activeTab = 'letter'">
            开始撰写
          </button>
        </div>

        <el-form :model="letterGen" :inline="true" label-width="65px" :rules="rules" ref="letterGen"
          :hide-required-asterisk="true">
          <div v-show="activeTab === 'information'">
            <el-form-item label="信纸" prop="paperId" style="margin-left: -20px;">
              <el-select v-model="letterGen.paperId" placeholder="请选择纸张" @change="handleChange"
                style="width: 200px; height: 10px;margin-top:10px;margin-left: 10px" height="10">
                <el-option v-for="item in repository.papers" :key="item.id" :label="item.name" :value="item.id"
                  :disabled="letterGen.letterType !== item.type" />
              </el-select>
            </el-form-item>

            <el-form-item label="字体" prop="fontId" style="margin-left: -20px; ">
              <el-select v-model="letterGen.fontId" placeholder="请选择字体" @change="handleChange"
                style="width: 200px; height: 10px;margin-top:10px;margin-left: 10px" height="10">
                <el-option v-for="item in repository.fonts" :key="item.id" :label="item.name" :value="item.id" />
              </el-select>
            </el-form-item>
            <el-form-item label="颜色" prop="fontColorId" style="margin-left: -20px; ">
              <el-select v-model="letterGen.fontColorId" placeholder="请选择字体颜色" @change="handleChange"
                style=" width: 200px; height: 10px;margin-top:10px;margin-left: 10px" height="10">
                <el-option v-for=" item in repository.fontColors" :key="item.id" :label="item.description"
                  :value="item.id" />
              </el-select>
            </el-form-item>
            <el-form-item label="款式" prop="letterType" style="margin-left: -20px; ">
              <el-select v-model="letterGen.letterType" placeholder="请选择信件款式" @change="handleLetterTypeChange"
                style=" width: 200px; height: 10px;margin-top:10px;margin-left: 10px" height="10">
                <el-option v-for=" item in letterTypes" :key="item.id" :label="item.label" :value="item.id" />
              </el-select>
            </el-form-item>
            <el-form-item label=" 寄信人" prop="senderName" style="margin-left: -20px; " label-width="75px">
              <el-input :maxlength="15" show-word-limit v-model="letterGen.senderName" placeholder="请输入寄信人姓名"
                @input="handleChange" style="margin-top:10px;height: 10px; width: 200px;" />
            </el-form-item>
            <el-form-item label="收信人" prop="recipientName" style="margin-left: -20px;" label-width="75px">
              <!-- <el-input v-model="letterGen.recipientName" placeholder="请输入收信人" @input="handleCheckFriend"
                style="margin-top:10px;" /> -->
              <el-select v-model="selectedFriendId" filterable allow-create default-first-option placeholder="请输入收信人姓名"
                no-data-text="没有好友哦，直接输入名字吧" style="margin-top:10px; height: 10px;width: 200px;"
                @change="handleChangeFriend">
                <el-option v-for="item in friends" :key="item.id" :label="item.remark || item.name" :value="item.id">
                  <span v-if="item.remark" style="float: left">{{ item.name + '(' + item.remark + ')' }}</span>
                  <span v-else style="float: left">{{ item.name }}</span>
                  <span style="float: right; color: #8492a6; font-size: 13px">{{ item.email }}</span>
                </el-option>
              </el-select>
            </el-form-item>

          </div>
          <div v-show="activeTab === 'letter'" class="text">
            <el-form-item label="信的内容：" prop="letterContent" style=" text-align: center;" label-width="100px">
              <el-input v-model="letterGen.letterContent" :maxlength="currentMaxlength" show-word-limit
                placeholder="请开始写信吧" required type="textarea" :autosize="{ minRows: 20, maxRows: 30 }"
                style="width: 300px; margin-left: -20px;" resize="none" @input="handleChange"></el-input>
            </el-form-item>
          </div>
        </el-form>
      </div>
      <div class="right">
        <el-skeleton style="width: 100%; margin:0 auto" :loading="loading" animated>
          <template slot="template">

            <div class="waibu" style="margin-top: -13px;">
              <div class="scene">
                <div class="objects">
                  <div class="square"></div>
                  <div class="circle"></div>
                  <div class="triangle"></div>
                </div>
                <div class="wizard">
                  <div class="body"></div>
                  <div class="right-arm">
                    <div class="right-hand"></div>
                  </div>
                  <div class="left-arm">
                    <div class="left-hand"></div>
                  </div>
                  <div class="head">
                    <div class="beard"></div>
                    <div class="face">
                      <div class="adds"></div>
                    </div>
                    <div class="hat">
                      <div class="hat-of-the-hat"></div>
                      <div class="four-point-star --first"></div>
                      <div class="four-point-star --second"></div>
                      <div class="four-point-star --third"></div>
                    </div>
                  </div>
                </div>
              </div>
              <div class="letterProgress"></div>
              <div class="noise"></div>
            </div>
            <!-- <el-skeleton-item variant="image" style="width: 100%; height: 500px;" /> -->
          </template>
          <template>
            <!-- <img crossorigin="anonymous" :src="backImageUrl" style="max-width: 100%; max-height: 640px;" /> -->
            <el-image crossorigin="anonymous" :src="letterUrl"
              style="max-width: 100%; max-height: 640px;margin-left: 100px;">
              <div slot="placeholder" style="width: 100% ;height: 100%;">
                <img crossorigin="anonymous" :src="backImageUrl"
                  style="max-width: 100%; max-height: 640px; margin-top: 13px;" />
                <!-- <i class="el-icon-loading" style="line-height: 250px; color: #A52328; font-size: 30px;"></i> -->
              </div>
              <div slot="error" class="image-slot" style="width: 100%; height: 100%; color: #A52328;">
                <i class="el-icon-picture-outline" style="line-height: 250px; "></i>
                加载失败
              </div>
            </el-image>

          </template>
        </el-skeleton>
      </div>
      <button class="panda-button" @click="handleShowConfirmDialog">
        <div class="panda_face">
          <svg x="0px" y="0px" width="100%" height="100%" viewBox="0 0 91.7 71.3">
            <path class="panda_icon panda_white"
              d="M86.6,20.4v15.3h5.1v20.4h-5.1v5.1h-5.1v5H71.3v5.1H20.4v-5.1H10.2v-5H5.1v-5.1H0V35.7h5.1V20.4L25.5,5.1h40.7L86.6,20.4z" />
            <path class="panda_icon pink" d="M40.7,45.8h10.2v5.1H40.7V45.8z M10.2,25.4H5.1v-5.1H0V10.2h5.1V5.1h5.1V0h10.2v5.1h5.1v5.1h-5.1v5.1h-5.1v5h-5.1V25.4z
            M30.5,45.8h-5v5.1H15.3V35.6h5.1v-5.1h10.1V45.8z M61.1,56H56v5.1H35.6V56H56v-5.1h5.1V56z M76.4,50.9H66.2v-5.1h-5.1V30.5h10.2
          v5.1h5.1V50.9z M91.7,20.3h-5.1v5.1h-5.1v-5.1h-5.1v-5h-5.1v-5.1h-5.1V5.1h5.1V0h10.2v5.1h5.1v5.1h5.1V20.3z" />
          </svg>
        </div>
        <span>发给好友!</span>
      </button>
      <!-- 确认信件信息 -->
      <el-dialog title="确认信件信息" :visible.sync="showConfirmDialog" width="650px" append-to-body>
        <el-form :model="letter" :inline="true" label-width="100px" :rules="rules" ref="letter" hide-required-asterisk>
          <el-form-item label="收信人邮箱" prop="recipientEmail">
            <el-input v-model="letter.recipientEmail" style="width: 200px;" placeholder="请输入收信人邮箱"
              @input="getFriendAddressesByEmail"></el-input>
          </el-form-item>

          <el-form-item label="收信人地址" prop="recipientAddress">
            <el-select v-model="recipientCountryId" filterable default-first-option placeholder="请选择收信人所在国家"
              style="margin-right:10px; height: 10px;width: 200px; float: left;"
              @change="handleChangeRecipientCountryId">
              <el-option v-for="item in countries" :key="item.id" :label="item.countryName" :value="item.id">
                <span style="float: left">{{ item.countryName }}</span>
                <span style="float: right; color: #8492a6; font-size: 13px">{{ item.countryNameEnglish }}</span>
              </el-option>
            </el-select>
            <avue-input-map @input="selectedFriendAddressId = null" v-show="recipientCountryId === 1" :params="params"
              :autosize="{ minRows: 1, maxRows: 3 }" placeholder="请选择收信地址" v-model="recipientAddress"
              style="float: left; margin-top: 3px;"></avue-input-map>
            <el-input v-show="recipientCountryId !== 1" @input="selectedFriendAddressId = null"
              v-model="letter.recipientAddress.formattedAddress" placeholder="请输入详细地址" required type="textarea"
              :autosize="{ minRows: 1, maxRows: 3 }"
              style="width: 200px; min-height: 35px !important; margin-top: 3px; float: left;" resize="none" />
            <el-popover placement="right" v-model="showFriendAddresses" width="250">
              <div class="AddressDialog">
                <i class="el-icon-close closeIcon" @click="showFriendAddresses = false"></i>
                <div class="addressTip">选择收信人地址</div>
                <div class="addressContent">
                  <div class="UserAddress beautify-scroll-bar ">
                    <div class="userAddressItem" :class="{ isSelected: selectedFriendAddressId === address.id }"
                      @click="selectFriendAddress(address.id)" v-for="address in friendAddresses" :key="address.id">
                      <div class="defaultTip" v-if="address.isDefault == 'true'">默认</div>
                      <div class="recipientInfo">
                        <div class="recipient">{{ letterGen.recipientName }}</div>
                        <div class="addressCountry">{{ getCountyNameById(address.countryId) }}</div>
                      </div>
                      <div class="detailAddress f-els-2">{{ address.formattedAddress }}</div>
                    </div>

                  </div>
                </div>
              </div>
              <i slot="reference" v-show="showFriendAddressBtn" class="el-icon-s-home address-aside"></i>
            </el-popover>
          </el-form-item>
          <el-form-item label="我的地址" prop="senderAddress">
            <el-select v-model="senderCountryId" filterable default-first-option placeholder="请选择寄信人所在国家"
              style="margin-right:10px; height: 10px;width: 200px; float: left;" @change="handleChangeSenderCountryId">
              <el-option v-for="item in countries" :key="item.id" :label="item.countryName" :value="item.id">
                <span style="float: left">{{ item.countryName }}</span>
                <span style="float: right; color: #8492a6; font-size: 13px">{{ item.countryNameEnglish }}</span>
              </el-option>
            </el-select>
            <avue-input-map v-show="senderCountryId === 1" :params="params" :autosize="{ minRows: 1, maxRows: 3 }"
              placeholder="请选择寄信地址" v-model="senderAddress" style="float: left; margin-top: 3px;"></avue-input-map>
            <el-input v-show="senderCountryId !== 1" v-model="letter.senderAddress.formattedAddress"
              placeholder="请输入详细地址" required type="textarea" :autosize="{ minRows: 1, maxRows: 3 }"
              style="width: 200px; min-height: 35px !important; margin-top: 3px;  float: left;" resize="none" />
            <el-popover placement="right" v-model="showMyAddresses" width="250">
              <div class="AddressDialog">
                <i class="el-icon-close closeIcon" @click="showMyAddresses = false"></i>
                <div class="addressTip">选择我的地址</div>
                <div class="addressContent">
                  <div class="UserAddress beautify-scroll-bar ">
                    <div class="userAddressItem " :class="{ isSelected: selectedMyAddressId === address.id }"
                      @click="selectMyAddress(address.id)" v-for="address in myAddresses" :key="address.id">
                      <div class="defaultTip" v-if="address.isDefault == 'true'">默认</div>
                      <div class="recipientInfo">
                        <div class="recipient">{{ letterGen.senderName }}</div>
                        <div class="addressCountry">{{ getCountyNameById(address.countryId) }}</div>
                      </div>
                      <div class="detailAddress f-els-2">{{ address.formattedAddress }}</div>
                    </div>

                  </div>
                </div>
              </div>
              <i slot="reference" v-show="showMyAddressBtn" class="el-icon-s-home address-aside"></i>
            </el-popover>
          </el-form-item>
          <el-form-item label="携带猪仔钱" prop="piggyMoney">
            <el-input-number v-model="letter.piggyMoney" style="width: 200px;" controls-position="right"
              @change="handleChange" :min="0" :max="500"></el-input-number>
          </el-form-item>
        </el-form>

        <template #footer>
          <div class="dialog-footer">
            <el-button type="primary" @click="submitPre" :loading="sendLoading">发 送</el-button>
            <el-button @click="cancel">取 消</el-button>
          </div>
        </template>
      </el-dialog>

    </div>
    <div v-else-if="currentView === 'send'">

      <div class="content" v-if="showContent">
        <!-- <p>
          <label style="margin-right: 16px;">切换 Loading</label>
          <el-switch v-model="loading" />
        </p> -->
        <el-row class="row-bg" justify="space-around">
          <el-empty v-if="isEmpty" style="height: 100%;" description="没有写过侨批哦,快去写一封吧!"></el-empty>
          <el-col :span="6" v-for="item in mySendList" :key="item.id">
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
              <div><i class="el-icon-full-screen" :class="{ 'toggle-icon': isFullscreen }"
                  @click="toggleFullScreen"></i>
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
                <i class="el-icon-message"></i>
                收信人邮箱：{{ letterVo.recipientEmail }}
              </div>
              <div class="card_item">
                <i class="el-icon-map-location"></i>
                收信人地址：{{ letterVo.recipientAddress.formattedAddress }}
              </div>
              <div class="card_item">
                <i class="el-icon-date"></i>
                写信时间：{{ formatDateTime(letterVo.createTime) }}
              </div>
              <div class="card_item">
                <i class="el-icon-s-home"></i>
                我的地址：{{ letterVo.senderAddress.formattedAddress }}
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
                <el-button type="text" @click="writeMore">再写一封</el-button>
                <el-divider v-if="this.letterVo.status === 2" direction="vertical"></el-divider>

                <el-popover placement="bottom" width="200" v-model="isUseCard">
                  <el-select no-data-text="没有功能卡哦" v-model="useCardDto.cardId" placeholder="请选择功能卡">
                    <el-option v-for="item in myFunctionCards" :key="item.id" :label="item.cardName" :value="item.id">
                      <span style="float: left">{{ item.cardName }}</span>
                      <span style="float: right; color: #8492a6; font-size: 13px">{{ item.number }}</span>
                    </el-option>
                  </el-select>
                  <div style="text-align: right; margin: 0;margin-top: 10px;">
                    <el-button size="mini" type="text" @click="isUseCard = false">取消</el-button>
                    <el-button type="primary" size="mini" @click="useCard">确定</el-button>
                  </div>

                  <el-button v-show="this.letterVo.status === 2" slot="reference" type="text"><i
                      class="el-icon-postcard" style="color: #fff;"></i>使用卡片</el-button>
                </el-popover>

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
          <div class="progress">
            <div style="width: 100%; margin-top: 10px ;margin-left: 55px">
              <el-progress ref="progress" :text-inside="true" :stroke-width="24"
                :percentage="letterVo.deliveryProgress / 100" style="height: 45px; width: 80%; line-height: 60px;">
              </el-progress>
            </div>
            <div style="color: #fff;">
              <i class="el-icon-position"></i>
              信件状态：{{ getLetterStatusLabel(letterVo.status) }}
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

  </div>
</template>

<script>
import { getCountries, getUserRepository, getUserFriends, getMyAddress, getMyFunctionCard } from '@/api/user'
import { getFontPaperLimit, sendLetter, getMySendLetter } from '@/api/letter'
import { useCard } from '@/api/card'
import useUserStore from '@/store/modules/user'
import { MessageBox } from 'element-ui'
// 自定义节流函数
function throttle(func, limit) {
  let inThrottle
  return function () {
    const args = arguments
    const context = this
    if (!inThrottle) {
      func.apply(context, args)
      inThrottle = true
      setTimeout(() => (inThrottle = false), limit)
    }
  }
}

export default {
  name: 'WriteLetter',
  data() {
    // eslint-disable-next-line
    var validateAddress = (rule, value, callback) => {
      this.setAddressByCountryId()
      // 定义市级以上城市列表,以及省、市、自治区、特别行政区等关键字列表
      const cities = [
        '北京', '上海', '天津', '重庆',
        '广州', '深圳', '成都', '武汉', '南京', '杭州', '西安', '沈阳', '长沙', '济南', '青岛', '大连', '苏州', '宁波', '厦门', '福州',
        '哈尔滨', '长春', '合肥', '南昌', '郑州', '昆明', '贵阳', '南宁', '海口', '太原', '石家庄', '呼和浩特', '乌鲁木齐', '银川', '兰州', '西宁', '拉萨',
        '香港', '澳门', '台湾', '河北', '山西', '黑龙江', '吉林', '辽宁', '江苏', '浙江', '安徽', '福建', '江西', '山东', '河南', '湖北', '湖南', '广东', '海南', '四川', '贵州', '云南', '陕西', '甘肃', '青海', '台湾', '内蒙古', '广西', '西藏', '宁夏', '新疆'
      ]

      // 定义需要包含的关键字列表
      const requiredKeywords = ['省', '市', '自治区', '特别行政区']
      const name = value.formattedAddress

      if (value.length === 0 || name === '') {
        callback(new Error('请输入详细地址'))
      }
      if (value.countryId !== 1 && value.countryId !== null) {
        callback()
      }
      // 检查 name 是否包含市级以上城市之一
      if (!cities.some(city => name.includes(city))) {
        // throw new Error('地址必须包含市级以上城市之一')
        callback(new Error('请再次点击以获得详细地址'))
      }

      // 检查 name 是否包含需要包含的关键字之一
      if (!requiredKeywords.some(keyword => name.includes(keyword))) {
        callback(new Error('请再次点击以获得详细地址'))
        // throw new Error('地址必须包含省、市、自治区或特别行政区')
      }

      callback()
    }
    return {
      activeTab: 'information', // 默认激活 "For rent" 标签
      // websocket
      websocket: null,
      sendLoading: false,
      isUseCard: false,
      useCardDto: {
        cardId: '',
        letterId: ''
      },
      letterStatusMap: {
        0: '未发送',
        1: '已发送',
        2: '传递中',
        3: '已送达'
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
          countryId: '',
          formattedAddress: '',
          longitude: '',
          latitude: ''
        },
        recipientAddress: {
          countryId: '',
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
        status: '',
        letterType: 1,
        piggyMoney: ''
      },
      coverUrl: '',
      imageUrl: '',
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
      mySendList: ['', '', '', '', '', '', '', ''],
      isEmpty: false,
      myFunctionCards: {},
      // 图片加载状态
      loading: false,
      // 地图组件
      params: {
        zoom: 15, // 地图显示的缩放级别范围，在PC上，默认范围[3,18]，取值范围[3-18]；在移动设备上，默认范围[3-19]，取值范围[3-19]
        zoomEnable: true, // 地图是否可通过鼠标拖拽平移。
        dragEnable: true // 地图是否可缩放
      },
      pageList: [{
        name: '写信',
        value: 'write'
      }, {
        name: '忆往昔',
        value: 'send'
      }],
      rules: {
        senderName: [
          { required: true, message: '请输入寄信人', trigger: 'blur' },
          { min: 1, max: 15, message: '最多输入15个字符', trigger: 'blur' }
        ],
        recipientName: [
          { required: true, message: '请输入收信人', trigger: 'blur' },
          { min: 1, max: 15, message: '最多输入15个字符', trigger: 'blur' }
        ],
        fontId: [
          { required: true, message: '请选择字体', trigger: 'change' }
        ],
        paperId: [
          { required: true, message: '请选择纸张', trigger: 'change' }
        ],
        fontColorId: [
          { required: true, message: '请选择颜色', trigger: 'change' }
        ],
        letterContent: [
          { required: true, message: '信件内容不能为空', trigger: 'blur' },
          { min: 1, message: '请至少输入 1 个字符', trigger: 'blur' }
        ],
        recipientEmail: [
          { required: true, message: '请输入收信人邮箱', trigger: 'blur' },
          { pattern: /^([a-zA-Z0-9]+[-_]?)+@[a-zA-Z0-9]+\.[a-z]+$/, message: '请输入正确的邮箱格式', trigger: 'blur' }
        ],
        senderAddress: [
          { required: true, validator: validateAddress }
        ],
        recipientAddress: [
          { required: true, validator: validateAddress }
        ],
        letterType: [
          { required: true, message: '请选择信件款式', trigger: 'change' }
        ]
      },
      letterTypes: [{
        label: '侨批',
        value: '侨批',
        default: true,
        id: 1
      }, {
        label: '普通信件',
        value: '普通信件',
        id: 2
      }],
      letterGen: {
        fontId: '',
        fontColorId: '',
        paperId: '',
        letterContent: ' ',
        senderName: useUserStore().name,
        recipientName: '',
        letterType: 1
      },
      backImageUrl: '',
      letterUrl: '',
      letter: {
        letterContent: '',
        senderName: '',
        recipientName: '',
        senderAddress: {
          countryId: '',
          formattedAddress: '',
          longitude: '',
          latitude: ''
        },
        recipientAddress: {
          countryId: '',
          formattedAddress: '',
          longitude: '',
          latitude: ''
        },
        recipientEmail: '',
        recipientUserId: '',
        signetId: '',
        letterType: '1',
        piggyMoney: ''
      },
      repository: {
        fonts: [],
        papers: [],
        fontColors: [],
        signets: []
      },
      countries: [],
      recipientCountryId: 1,
      senderCountryId: 1,
      myAddresses: [],
      friendAddresses: [],
      showMyAddresses: false,
      showFriendAddresses: false,
      showMyAddressBtn: false,
      showFriendAddressBtn: false,
      selectedFriendAddressId: '',
      selectedMyAddressId: '',
      showConfirmDialog: false, // 控制弹窗的显示
      senderAddress: [116.397455, 39.909187, ''],
      recipientAddress: [116.397455, 39.909187, ''],
      friends: [],
      selectedFriendId: '',
      showContent: false,
      currentView: '',
      newLayout: false, // 控制表单布局
      showTip: true,
      disabled: false,
      fontPaperLimitList: [],
      timeOutCount: 0,
      linked: false
    }
  },

  methods: {
    getFontPaperLimit() {
      getFontPaperLimit().then(response => {
        this.fontPaperLimitList = response.data
      })
    },
    handleClickTab(event) {
      // 阻止事件冒泡，防止点击按钮时触发外层 div 的点击事件
      event.stopPropagation()
    },
    loadFirstImage(item) {
      this.showSend()
      setTimeout(() => {
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
      }, 2500)
    },
    writeMore() {
      this.closeImageViewer()
      setTimeout(() => {
        this.currentView = 'write'
      }, 2300)
    },
    handleShowProgress() {
      this.showProgress = !this.showProgress
      const progress = document.querySelector('.progress')
      if (this.showProgress) {
        progress.style.opacity = '1'
      } else {
        progress.style.opacity = '0'
      }
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
      this.letterVo.deliveryProgress = (parseInt(item.deliveryProgress))
      // console.log(this.letterVo)

      setTimeout(() => {
        if (!this.showLetterDetail) {
          this.handleShowDetail()
        }
        if (!this.showProgress) {
          this.handleShowProgress()
        }
        this.scale = 0.9
        // 更新进度
        this.updateInnerTextBackground()
      }, 100)
    },
    closeImageViewer() {
      // this.isVisible = false
      this.isFullscreen = false

      if (this.showLetterDetail) {
        this.handleShowDetail()
      }
      if (this.showProgress) {
        this.handleShowProgress()
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
      // 如果缩放比例在 1.4 到 1.5 之间，且显示详情，则触发隐藏详情
      if (this.scale >= 1.4 && this.scale <= 1.5 && this.showLetterDetail) {
        this.handleShowDetail()
      }
      if (this.scale >= 1.2 && this.showProgress) {
        this.handleShowProgress()
      }
      if (this.scale <= 1.2 && !this.showProgress) {
        this.handleShowProgress()
      }
    },
    initWebSocket() {
      if ('WebSocket' in window) {
        // 这里记得要改成你自己的ip
        if (!this.isLinked) {
          this.websocket = new WebSocket('ws://localhost:8080/ws/letterGen', useUserStore().token)
          this.websocket.onerror = this.onError
          this.websocket.onopen = this.onOpen
          this.websocket.onmessage = this.onMessage
          this.websocket.onclose = this.onClose
          this.websocketTimeout = setTimeout(() => {
            if (this.websocket.readyState !== WebSocket.OPEN) {
              this.websocket.close()
            }
          }, 3000) // 设置超时时间为3000毫秒
        }
      } else {
        alert('Not support websocket')
      }
    },
    addUnloadListener() {
      window.addEventListener('beforeunload', this.handleBeforeUnload)
    },
    removeUnloadListener() {
      window.removeEventListener('beforeunload', this.handleBeforeUnload)
    },
    handleBeforeUnload() {
      if (this.websocket) {
        this.websocket.close()
      }
    },
    onError() {
      if (!this.linked) {
        // this.$message.error('连接超时,请刷新页面重试!')
        this.websocket = null
      }
    },
    onOpen() {
      this.linked = true
      setTimeout(() => {
        this.websocket.send(JSON.stringify(this.letterGen))
      }, 500)
    },
    onClose() {
      this.timeOutCount++
      if (this.timeOutCount > 5) {
        this.$notify.error({
          title: '错误',
          message: '哎呀，亲爱的朋友们，服务器小憩了一下，很快就会回来啦！请稍后再试，或者刷新一下页面。侨宝在这里等你们哦！🌟📶',
          offset: 100
        })
      } else {
        this.websocket = null
        this.linked = false
        this.initWebSocket()
      }
      // this.$message.warning('连接已断开')
    },
    onMessage(event) {
      if (event.data === 'success') {
        this.loading = false
      } else {
        this.letterUrl = 'data:image/jpg;base64,' + event.data
      }
    },
    handleLetterTypeChange() {
      this.letterGen.paperId = this.repository.papers.find(paper => paper.type === this.letterGen.letterType).id
      this.handleChange()
    },
    // websocket
    handleChange() {
      const content = this.letterGen.letterContent
      if (/[a-zA-Z0-9]/.test(content) && this.letterGen.letterType === 1) {
        this.letterGen.letterType = 2
        const firstPaperWithType2 = this.repository.papers.find(paper => paper.type === 2)
        if (firstPaperWithType2) {
          this.letterGen.paperId = firstPaperWithType2.id
        } else {
          this.letterGen.paperId = 4
        }
        this.$message({
          message: '检测到您的信件内容中包含英文或数字，已自动切换为普通信件',
          type: 'warning'
        })
      }
      this.showTip = true
      if (this.websocket == null) {
        this.initWebSocket()
      } else {
        this.send('letterGen')
      }
    },
    send: throttle(function () {
      this.$refs.letterGen.validate((valid) => {
        if (valid) {
          this.loading = true
          if (this.websocket && this.letterGen) {
            this.backImageUrl = this.letterUrl
            this.websocket.send(JSON.stringify(this.letterGen))
            this.loading = false
          }
        } else {
          if (this.letterGen.recipientName === '' || this.letterGen.senderName === '') {
            this.activeTab = 'information'
          } else {
            this.activeTab = 'letter'
          }
          return false
        }
      })
      // 做了一个50ms的节流,防止频繁发送,这里可以根据实际情况调整,缓解服务器压力
    }, 50),
    closeWebSocket() {
      if (this.websocket) {
        this.websocket.close()
      }
    },
    // 查看大图函数结束
    getLetterStatusLabel(statusId) {
      if (this.letterVo.status === 2 && this.letterVo.speedRate !== '1' && this.letterVo.speedRate !== null) {
        return `${this.letterVo.speedRate}倍加速` + this.letterStatusMap[statusId]
      }
      return this.letterStatusMap[statusId] || '未知状态'
    },
    getMySendLetter() {
      this.loading = true
      getMySendLetter().then(res => {
        this.mySendList = res.data
        this.loading = false
        // setTimeout(() => {
        //   this.loading = false
        // }, 1000)
        if (this.mySendList.length === 0) {
          this.isEmpty = true
        } else {
          this.isEmpty = false
        }
      })
    },
    getMyFunctionCard() {
      getMyFunctionCard().then(res => {
        this.myFunctionCards = res.data
      })
    },
    // 使用卡片
    useCard() {
      this.useCardDto.letterId = this.letterVo.id
      if (this.letterVo.status !== 2) {
        this.$notify({
          title: '使用失败',
          message: '只有传递中的侨批才能使用功能卡哦!😅',
          type: 'warning',
          offset: 100
        })
        return
      }
      if (this.letterVo.speedRate !== '1' && this.letterVo.speedRate !== null) {
        const card = this.myFunctionCards.find(card => card.id === this.useCardDto.cardId)
        // console.log(card)
        if (card && card.cardType === 1) {
          MessageBox.confirm(`当前已处于${this.letterVo.speedRate}倍加速中，继续使用加速卡将覆盖原来加速效果！是否继续使用？`, '侨缘信使', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }).then(() => {
            useCard(this.useCardDto).then(res => {
              Object.assign(this.letterVo, res.data)
              this.$notify({
                title: '使用成功',
                message: '侨批很快就要送达了!🥳',
                type: 'success',
                offset: 100
              })
              this.getMyFunctionCard()
            })
          })
          this.isUseCard = false
          return
        }
      }
      if (this.useCardDto.cardId !== '') {
        useCard(this.useCardDto).then(res => {
          Object.assign(this.letterVo, res.data)
          this.$notify({
            title: '使用成功',
            message: '侨批很快就要送达了!🥳',
            type: 'success',
            offset: 100
          })
          this.getMyFunctionCard()
        })
      }
      this.isUseCard = false
    },
    handleChangePage(val) {
      // console.log(val)
      // 切换页面
      if (val === 'write') {
        this.showWrite()
      } else {
        this.showSend()
      }
    },
    setAddressByCountryId() {
      if (this.letter.recipientAddress.countryId === 1 || this.letter.recipientAddress.countryId === '' || this.letter.recipientAddress.countryId === null) {
        this.letter.recipientAddress.countryId = 1
        this.letter.recipientAddress.longitude = this.recipientAddress[0]
        this.letter.recipientAddress.latitude = this.recipientAddress[1]
        this.letter.recipientAddress.formattedAddress = this.recipientAddress[2]
      }
      if (this.letter.senderAddress.countryId === 1 || this.letter.senderAddress.countryId === '' || this.letter.senderAddress.countryId === null) {
        this.letter.senderAddress.countryId = 1
        this.letter.senderAddress.longitude = this.senderAddress[0]
        this.letter.senderAddress.latitude = this.senderAddress[1]
        this.letter.senderAddress.formattedAddress = this.senderAddress[2]
      }
    },
    submitPre() {
      this.sendLoading = true
      // console.log(this.letter)
      this.setAddressByCountryId()
      this.letter.letterContent = this.letterGen.letterContent
      this.letter.senderName = this.letterGen.senderName
      this.letter.letterType = this.letterGen.letterType
      this.letter.recipientName = this.letterGen.recipientName
      this.submitForm('letter')
    },
    /** 提交按钮 */
    submitForm(letter) {
      this.$refs[letter].validate((valid) => {
        if (valid) {
          // 提交表单
          // console.log('提交表单:', this.letter)
          sendLetter(this.letter).then(res => {
            this.$notify({
              title: '发送成功',
              message: '侨批在传递中了哦!🥳',
              type: 'success',
              offset: 100
            })
            this.showConfirmDialog = false
            this.loadFirstImage(res.data)
          })
          setTimeout(() => {
            this.sendLoading = false
          }, 1000)
        } else {
          this.sendLoading = false
          return false
        }
      })
    },
    cancel() {
      this.showConfirmDialog = false
    },
    getMyAddress() {
      // 获取我的地址
      getMyAddress().then(res => {
        this.myAddresses = res.data
      })
    },
    handleChangeRecipientCountryId(val) {
      if ((val !== 1 && this.letter.senderAddress.countryId !== 1) && (val !== 2 && this.letter.senderAddress.countryId !== 2)) {
        this.$message.warning('寄信人和收信人不能都不在中国')
        this.recipientCountryId = 1
      } else {
        if (val === 1) {
          this.recipientAddress = [116.397455, 39.909187, '']
        } else if (this.letter.recipientAddress.countryId === 1) {
          this.letter.recipientAddress.formattedAddress = ''
        }
        this.letter.recipientAddress.countryId = val
        this.selectedFriendAddressId = ''
      }
    },
    handleChangeSenderCountryId(val) {
      if ((val !== 1 && this.letter.recipientAddress.countryId !== 1) && (val !== 2 && this.letter.recipientAddress.countryId !== 2)) {
        this.$message.warning('寄信人和收信人不能都不在中国')
        this.senderCountryId = 1
      } else {
        if (val === 1) {
          this.senderAddress = [116.397455, 39.909187, '']
        } else if (this.letter.senderAddress.countryId === 1) {
          this.letter.senderAddress.formattedAddress = ''
        }
        this.letter.senderAddress.countryId = val
        this.selectedMyAddressId = ''
      }
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
    handleChangeFriend(selectedFriendId) {
      // console.log(selectedFriendId)
      for (const friend of this.friends) {
        if (friend.id === selectedFriendId) {
          this.letterGen.recipientName = friend.remark || friend.nickname
          this.letter.recipientUserId = friend.userId
          this.letter.recipientEmail = friend.email
          this.handleChange()
          return
        }
      }
      // 如果不是好友的情况,这个名字就是这个id
      this.letterGen.recipientName = selectedFriendId
      this.handleChange()
    },
    selectFriendAddress(addressId) {
      for (const address of this.friendAddresses) {
        if (address.id === addressId) {
          this.letter.recipientAddress = { ...address }
          this.selectedFriendAddressId = addressId
          this.showFriendAddresses = false
          this.recipientCountryId = address.countryId || 1
          if (this.recipientCountryId === 1) {
            this.recipientAddress = [address.longitude, address.latitude, address.formattedAddress]
          }
          break
        }
      }
    },
    selectMyAddress(addressId) {
      for (const address of this.myAddresses) {
        if (address.id === addressId) {
          this.letter.senderAddress = { ...address }
          this.selectedMyAddressId = addressId
          this.showMyAddresses = false
          this.senderCountryId = address.countryId || 1
          if (this.senderCountryId === 1) {
            this.senderAddress = [address.longitude, address.latitude, address.formattedAddress]
          }
          break
        }
      }
    },
    getFriendAddressesByEmail() {
      if (this.letter.recipientEmail !== '') {
        let friendAddresses = []
        for (const friend of this.friends) {
          if (friend.email === this.letter.recipientEmail) {
            this.letter.recipientUserId = friend.userId
            friendAddresses = friend.addresses
            break
          }
        }
        if (friendAddresses.length > 0) {
          this.friendAddresses = friendAddresses
          for (const address of friendAddresses) {
            if (address.isDefault === 'true') {
              this.selectedFriendAddressId = address.id
              this.letter.recipientAddress = { ...address }
              this.recipientCountryId = address.countryId || 1
              if (this.recipientCountryId === 1) {
                this.recipientAddress = [address.longitude, address.latitude, address.formattedAddress]
              }
              break
            }
          }
          this.showFriendAddressBtn = true
        } else {
          this.showFriendAddressBtn = false
        }
      }
    },
    handleShowConfirmDialog() {
      this.$refs.letterGen.validate((valid) => {
        if (valid) {
          this.showConfirmDialog = true

          this.getFriendAddressesByEmail()
          if (this.myAddresses.length > 0) {
            for (const address of this.myAddresses) {
              if (address.isDefault === 'true') {
                this.selectedMyAddressId = address.id
                this.letter.senderAddress = address
                this.senderCountryId = address.countryId || 1
                if (this.senderCountryId === 1) {
                  this.senderAddress = [address.longitude, address.latitude, address.formattedAddress]
                }
                break
              }
            }
            this.showMyAddressBtn = true
          } else {
            this.showMyAddressBtn = false
          }
        } else {
          this.$message.info('请填写完整信息')
        }
      })
    },
    getCountries() {
      // 获取国家列表
      getCountries().then(res => {
        this.countries = res.data
      })
    },
    getUserRepository() {
      this.loading = true
      // 获取用户仓库
      getUserRepository().then(res => {
        this.repository = res.data
        this.letterGen.fontId = this.repository.fonts[0].id
        this.letterGen.fontColorId = this.repository.fontColors[0].id
        this.letterGen.paperId = this.repository.papers[0].id
        this.letterUrl = this.repository.papers[0].previewImage
        this.loading = false
      })
    },
    getUserFriends() {
      // 获取用户好友列表
      getUserFriends().then(res => {
        this.friends = res.data
      })
    },
    showWrite() {
      this.showContent = true
      this.currentView = 'write'
      this.initWebSocket()
    },
    showSend() {
      this.getMySendLetter()
      this.getMyFunctionCard()
      this.showContent = true
      this.currentView = 'send'
    },
    updateInnerTextBackground() {
      if (this.isVisible) {
        const innerTextElement = this.$refs.progress.$el.querySelector('.el-progress-bar__innerText')
        if (innerTextElement) {
          innerTextElement.style.backgroundImage = `url(${this.innerTextBackground})`
          if (this.letterVo.deliveryProgress > 1000) {
            innerTextElement.style.margin = '-3px 3px'
          } else {
            innerTextElement.style.margin = '-3px -36px'
          }
        }
      }
    }
  },
  mounted() {
    this.showWrite()// 页面加载时自动触发
    // this.showSend()
    this.getUserRepository()
    this.getUserFriends()
    this.getMyAddress()
    this.getCountries()
    this.getFontPaperLimit()
  },
  beforeDestroy() {
    this.removeUnloadListener()
    this.closeWebSocket()
  },
  computed: {
    isLinked() {
      return this.linked && this.websocket
    },
    innerTextBackground() {
      if (this.letterVo.deliveryProgress / 100 <= 40) {
        return require('@/assets/imgs/car1.png')
      } else if (this.letterVo.deliveryProgress / 100 <= 80) {
        return require('@/assets/imgs/car2.png')
      } else {
        return require('@/assets/imgs/car3.png')
      }
    },
    currentMaxlength() {
      const { paperId, fontId } = this.letterGen
      const limitItem = this.fontPaperLimitList.find(
        item => item.paperId === paperId && item.fontId === fontId
      )
      return limitItem ? limitItem.fitNumber : 300 // 默认值为300
    }
  },
  watch: {
    'letterVo.deliveryProgress': {
      handler() {
        this.updateInnerTextBackground()
      },
      immediate: true
    }
  }
}
</script>

<style scoped>
.dialog-address {
  background: #fff;
  border: 1px solid rgba(0, 17, 51, .05);
  border-radius: 12px;
  box-shadow: 0 0 6px 0 rgba(0, 17, 51, .1);
  padding: 16px;
}

.UserAddress {
  align-self: flex-start;
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  max-height: 176px;
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

.UserAddress .userAddressItem {
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

.closeIcon {
  cursor: pointer;
  display: inline-block;
  font-size: 18px;
  line-height: 18px;
  position: absolute;
  right: 0;
  top: -4px
}

.UserAddress .userAddressItem:nth-of-type(2n) {
  margin-top: 12px
}

.UserAddress .userAddressItem:nth-of-type(n+3) {
  margin-top: 12px
}

.UserAddress .userAddressItem.isSelected,
.UserAddress .userAddressItem:hover {
  background-color: #fff;
  border: 1px solid #ff6200
}

.UserAddress .defaultTip {
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

.row-bg {
  width: 95%;
  height: 600px;
  position: relative;
  padding-top: 20px;
  padding-left: 20px;
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
  /* 使用 Flexbox 布局 */
  align-items: flex-start;
  /* 垂直对齐子元素 */
  line-height: 40px;
}

.selectPage {
  position: absolute;
  z-index: 5;
  /*提高层级 */
  top: 20px;
  right: 50px;
  border: 0;
  cursor: pointer;
  /* 将鼠标指针设置为手型 */
  outline: none;
  /* 移除焦点时的边框 */
}

.left {
  position: absolute;
  width: 520px;
  top: 10px;
  left: 80px;
  padding: 20px;
  box-sizing: border-box;
  /* background-color: #f9f9f9; */
  background: url('@/assets/imgss/yellowbackground.png') 0 0 / 400px auto repeat, #f9f9f9;
  border-radius: 25px;
  text-align: left;
}

.right {
  position: absolute;
  /* background-color: rgb(222, 201, 162); */
  padding: 20px;
  border-radius: 25px;
  right: 150px;
  top: 115px;
  width: 700px;
  height: 450px;
  display: flex;
  justify-content: center;
  align-items: center;
}

.text {
  margin-left: -25px;
}

.text label {
  display: block;
  margin-bottom: 5px;
}

.text textarea {
  width: 350px;
  height: 400px;
  padding: 8px;
  border-radius: 4px;
  resize: none;
  outline: none;
  border: 0;
}

.left ::v-deep .el-input__inner {
  height: 35px;
  border: 1px solid #1296db;
  /* background-color: rgba(0, 0, 0, 0.1) !important; */
  color: #666;
}

::v-deep .el-input__inner,
::v-deep .el-textarea__inner {
  resize: none;
  min-height: 35px !important;
  min-width: 200px;
  height: 35px;
  border: 1px solid #1296db;
  /* background-color: rgba(0, 0, 0, 0.1) !important; */
  color: #666;
}

.selectPage ::v-deep .el-input__inner,
.selectPage ::v-deep .el-textarea__inner {
  height: 35px;
  border-radius: 4px;
  min-width: 0;
  border: 1px solid #1296db;
  /* background-color: rgba(222, 201, 162, 0.6) !important; */
  color: #b00f0f;
}

::v-deep .el-dialog {
  margin-top: 200px !important;
  border-radius: 10px;
  /* background-color: rgba(222, 201, 162, 1) !important; */
  background: 0 0 / 400px auto repeat, #f9f9f9;
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

/* .progress {
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
} */

/* 平移动画 */
/* 移至common.css */
/* @charset "UTF-8"; */

.waibu {
  display: -webkit-box;
  display: flex;
  -webkit-box-align: center;
  align-items: center;
  -webkit-box-pack: center;
  justify-content: center;
  -webkit-box-orient: vertical;
  -webkit-box-direction: normal;
  flex-direction: column;
  flex-wrap: wrap;
  width: 700px;
  height: 450px;
  overflow: hidden;
  border-radius: 10px;
  background-color: rgb(222, 201, 162);

}

.scene {
  display: -webkit-box;
  display: flex
}

.wizard {
  position: relative;
  width: 190px;
  height: 240px
}

.body {
  position: absolute;
  bottom: 0;
  left: 68px;
  height: 100px;
  width: 60px;
  background: #3f64ce
}

.body::after {
  content: "";
  position: absolute;
  bottom: 0;
  left: 20px;
  height: 100px;
  width: 60px;
  background: #3f64ce;
  -webkit-transform: skewX(14deg);
  transform: skewX(14deg)
}

.right-arm {
  position: absolute;
  bottom: 74px;
  left: 110px;
  height: 44px;
  width: 90px;
  background: #3f64ce;
  border-radius: 22px;
  -webkit-transform-origin: 16px 22px;
  transform-origin: 16px 22px;
  -webkit-transform: rotate(70deg);
  transform: rotate(70deg);
  -webkit-animation: right_arm 3s ease-in-out infinite;
  animation: right_arm 3s ease-in-out infinite
}

.right-arm .right-hand {
  position: absolute;
  right: 8px;
  bottom: 8px;
  width: 30px;
  height: 30px;
  border-radius: 50%;
  background: #f1c5b4;
  -webkit-transform-origin: center center;
  transform-origin: center center;
  -webkit-transform: rotate(-40deg);
  transform: rotate(-40deg);
  -webkit-animation: right_hand 3s ease-in-out infinite;
  animation: right_hand 3s ease-in-out infinite
}

.right-arm .right-hand::after {
  content: "";
  position: absolute;
  right: 0px;
  top: -8px;
  width: 15px;
  height: 30px;
  border-radius: 10px;
  background: #f1c5b4;
  -webkit-transform: translateY(16px);
  transform: translateY(16px);
  -webkit-animation: right_finger 3s ease-in-out infinite;
  animation: right_finger 3s ease-in-out infinite
}

.left-arm {
  position: absolute;
  bottom: 74px;
  left: 26px;
  height: 44px;
  width: 70px;
  background: #3f64ce;
  border-bottom-left-radius: 8px;
  -webkit-transform-origin: 60px 26px;
  transform-origin: 60px 26px;
  -webkit-transform: rotate(-70deg);
  transform: rotate(-70deg);
  -webkit-animation: left_arm 3s ease-in-out infinite;
  animation: left_arm 3s ease-in-out infinite
}

.left-arm .left-hand {
  position: absolute;
  left: -18px;
  top: 0;
  width: 18px;
  height: 30px;
  border-top-left-radius: 35px;
  border-bottom-left-radius: 35px;
  background: #f1c5b4
}

.left-arm .left-hand::after {
  content: "";
  position: absolute;
  right: 0;
  top: 0;
  width: 30px;
  height: 15px;
  border-radius: 20px;
  background: #f1c5b4;
  -webkit-transform-origin: right bottom;
  transform-origin: right bottom;
  -webkit-transform: scaleX(0);
  transform: scaleX(0);
  -webkit-animation: left_finger 3s ease-in-out infinite;
  animation: left_finger 3s ease-in-out infinite
}

.head {
  position: absolute;
  top: 0;
  left: 14px;
  width: 160px;
  height: 210px;
  -webkit-transform-origin: center center;
  transform-origin: center center;
  -webkit-transform: rotate(-3deg);
  transform: rotate(-3deg);
  -webkit-animation: head 3s ease-in-out infinite;
  animation: head 3s ease-in-out infinite
}

.head .beard {
  position: absolute;
  bottom: 0;
  left: 38px;
  height: 106px;
  width: 80px;
  border-bottom-right-radius: 55%;
  background: #ffffff
}

.head .beard::after {
  content: "";
  position: absolute;
  top: 16px;
  left: -20px;
  width: 40px;
  height: 20px;
  border-radius: 20px;
  background: #ffffff
}

.head .face {
  position: absolute;
  bottom: 76px;
  left: 38px;
  height: 30px;
  width: 60px;
  background: #f1c5b4
}

.head .face::before {
  content: "";
  position: absolute;
  top: 0px;
  left: 40px;
  width: 20px;
  height: 40px;
  border-bottom-right-radius: 20px;
  border-bottom-left-radius: 20px;
  background: #f1c5b4
}

.head .face::after {
  content: "";
  position: absolute;
  top: 16px;
  left: -20px;
  width: 50px;
  height: 20px;
  border-radius: 20px;
  border-bottom-right-radius: 0px;
  background: #ffffff
}

.head .face .adds {
  position: absolute;
  top: 0px;
  left: -20px;
  width: 40px;
  height: 20px;
  border-radius: 20px;
  background: #f1c5b4
}

.head .face .adds::after {
  content: "";
  position: absolute;
  top: 5px;
  left: 80px;
  width: 15px;
  height: 20px;
  border-bottom-right-radius: 20px;
  border-top-right-radius: 20px;
  background: #f1c5b4
}

.head .hat {
  position: absolute;
  bottom: 106px;
  left: 0;
  width: 160px;
  height: 20px;
  border-radius: 20px;
  background: #3f64ce
}

.head .hat::before {
  content: "";
  position: absolute;
  top: -70px;
  left: 50%;
  -webkit-transform: translatex(-50%);
  transform: translatex(-50%);
  width: 0;
  height: 0;
  border-style: solid;
  border-width: 0 34px 70px 50px;
  border-color: transparent transparent #3f64ce transparent
}

.head .hat::after {
  content: "";
  position: absolute;
  top: 0;
  left: 0;
  width: 160px;
  height: 20px;
  background: #3f64ce;
  border-radius: 20px
}

.head .hat .hat-of-the-hat {
  position: absolute;
  bottom: 78px;
  left: 79px;
  width: 0;
  height: 0;
  border-style: solid;
  border-width: 0 25px 25px 19px;
  border-color: transparent transparent #3f64ce transparent
}

.head .hat .hat-of-the-hat::after {
  content: "";
  position: absolute;
  top: 6px;
  left: -4px;
  width: 35px;
  height: 10px;
  border-radius: 10px;
  border-bottom-left-radius: 0px;
  background: #3f64ce;
  -webkit-transform: rotate(40deg);
  transform: rotate(40deg)
}

.head .hat .four-point-star {
  position: absolute;
  width: 12px;
  height: 12px
}

.head .hat .four-point-star::after,
.head .hat .four-point-star::before {
  content: "";
  position: absolute;
  background: #ffffff;
  display: block;
  left: 0;
  width: 141.4213%;
  top: 0;
  bottom: 0;
  border-radius: 10%;
  -webkit-transform: rotate(66.66deg) skewX(45deg);
  transform: rotate(66.66deg) skewX(45deg)
}

.head .hat .four-point-star::after {
  -webkit-transform: rotate(156.66deg) skew(45deg);
  transform: rotate(156.66deg) skew(45deg)
}

.head .hat .four-point-star.--first {
  bottom: 28px;
  left: 46px
}

.head .hat .four-point-star.--second {
  bottom: 40px;
  left: 80px
}

.head .hat .four-point-star.--third {
  bottom: 15px;
  left: 108px
}

@-webkit-keyframes right_arm {
  0% {
    -webkit-transform: rotate(70deg);
    transform: rotate(70deg)
  }

  10% {
    -webkit-transform: rotate(8deg);
    transform: rotate(8deg)
  }

  15% {
    -webkit-transform: rotate(20deg);
    transform: rotate(20deg)
  }

  20% {
    -webkit-transform: rotate(10deg);
    transform: rotate(10deg)
  }

  25% {
    -webkit-transform: rotate(26deg);
    transform: rotate(26deg)
  }

  30% {
    -webkit-transform: rotate(10deg);
    transform: rotate(10deg)
  }

  35% {
    -webkit-transform: rotate(28deg);
    transform: rotate(28deg)
  }

  40% {
    -webkit-transform: rotate(9deg);
    transform: rotate(9deg)
  }

  45% {
    -webkit-transform: rotate(28deg);
    transform: rotate(28deg)
  }

  50% {
    -webkit-transform: rotate(8deg);
    transform: rotate(8deg)
  }

  58% {
    -webkit-transform: rotate(74deg);
    transform: rotate(74deg)
  }

  62% {
    -webkit-transform: rotate(70deg);
    transform: rotate(70deg)
  }
}

@keyframes right_arm {
  0% {
    -webkit-transform: rotate(70deg);
    transform: rotate(70deg)
  }

  10% {
    -webkit-transform: rotate(8deg);
    transform: rotate(8deg)
  }

  15% {
    -webkit-transform: rotate(20deg);
    transform: rotate(20deg)
  }

  20% {
    -webkit-transform: rotate(10deg);
    transform: rotate(10deg)
  }

  25% {
    -webkit-transform: rotate(26deg);
    transform: rotate(26deg)
  }

  30% {
    -webkit-transform: rotate(10deg);
    transform: rotate(10deg)
  }

  35% {
    -webkit-transform: rotate(28deg);
    transform: rotate(28deg)
  }

  40% {
    -webkit-transform: rotate(9deg);
    transform: rotate(9deg)
  }

  45% {
    -webkit-transform: rotate(28deg);
    transform: rotate(28deg)
  }

  50% {
    -webkit-transform: rotate(8deg);
    transform: rotate(8deg)
  }

  58% {
    -webkit-transform: rotate(74deg);
    transform: rotate(74deg)
  }

  62% {
    -webkit-transform: rotate(70deg);
    transform: rotate(70deg)
  }
}

@-webkit-keyframes left_arm {
  0% {
    -webkit-transform: rotate(-70deg);
    transform: rotate(-70deg)
  }

  10% {
    -webkit-transform: rotate(6deg);
    transform: rotate(6deg)
  }

  15% {
    -webkit-transform: rotate(-18deg);
    transform: rotate(-18deg)
  }

  20% {
    -webkit-transform: rotate(5deg);
    transform: rotate(5deg)
  }

  25% {
    -webkit-transform: rotate(-18deg);
    transform: rotate(-18deg)
  }

  30% {
    -webkit-transform: rotate(5deg);
    transform: rotate(5deg)
  }

  35% {
    -webkit-transform: rotate(-17deg);
    transform: rotate(-17deg)
  }

  40% {
    -webkit-transform: rotate(5deg);
    transform: rotate(5deg)
  }

  45% {
    -webkit-transform: rotate(-18deg);
    transform: rotate(-18deg)
  }

  50% {
    -webkit-transform: rotate(6deg);
    transform: rotate(6deg)
  }

  58% {
    -webkit-transform: rotate(-74deg);
    transform: rotate(-74deg)
  }

  62% {
    -webkit-transform: rotate(-70deg);
    transform: rotate(-70deg)
  }
}

@keyframes left_arm {
  0% {
    -webkit-transform: rotate(-70deg);
    transform: rotate(-70deg)
  }

  10% {
    -webkit-transform: rotate(6deg);
    transform: rotate(6deg)
  }

  15% {
    -webkit-transform: rotate(-18deg);
    transform: rotate(-18deg)
  }

  20% {
    -webkit-transform: rotate(5deg);
    transform: rotate(5deg)
  }

  25% {
    -webkit-transform: rotate(-18deg);
    transform: rotate(-18deg)
  }

  30% {
    -webkit-transform: rotate(5deg);
    transform: rotate(5deg)
  }

  35% {
    -webkit-transform: rotate(-17deg);
    transform: rotate(-17deg)
  }

  40% {
    -webkit-transform: rotate(5deg);
    transform: rotate(5deg)
  }

  45% {
    -webkit-transform: rotate(-18deg);
    transform: rotate(-18deg)
  }

  50% {
    -webkit-transform: rotate(6deg);
    transform: rotate(6deg)
  }

  58% {
    -webkit-transform: rotate(-74deg);
    transform: rotate(-74deg)
  }

  62% {
    -webkit-transform: rotate(-70deg);
    transform: rotate(-70deg)
  }
}

@-webkit-keyframes right_hand {
  0% {
    -webkit-transform: rotate(-40deg);
    transform: rotate(-40deg)
  }

  10% {
    -webkit-transform: rotate(-20deg);
    transform: rotate(-20deg)
  }

  15% {
    -webkit-transform: rotate(-5deg);
    transform: rotate(-5deg)
  }

  20% {
    -webkit-transform: rotate(-60deg);
    transform: rotate(-60deg)
  }

  25% {
    -webkit-transform: rotate(0deg);
    transform: rotate(0deg)
  }

  30% {
    -webkit-transform: rotate(-60deg);
    transform: rotate(-60deg)
  }

  35% {
    -webkit-transform: rotate(0deg);
    transform: rotate(0deg)
  }

  40% {
    -webkit-transform: rotate(-40deg);
    transform: rotate(-40deg)
  }

  45% {
    -webkit-transform: rotate(-60deg);
    transform: rotate(-60deg)
  }

  50% {
    -webkit-transform: rotate(10deg);
    transform: rotate(10deg)
  }

  60% {
    -webkit-transform: rotate(-40deg);
    transform: rotate(-40deg)
  }
}

@keyframes right_hand {
  0% {
    -webkit-transform: rotate(-40deg);
    transform: rotate(-40deg)
  }

  10% {
    -webkit-transform: rotate(-20deg);
    transform: rotate(-20deg)
  }

  15% {
    -webkit-transform: rotate(-5deg);
    transform: rotate(-5deg)
  }

  20% {
    -webkit-transform: rotate(-60deg);
    transform: rotate(-60deg)
  }

  25% {
    -webkit-transform: rotate(0deg);
    transform: rotate(0deg)
  }

  30% {
    -webkit-transform: rotate(-60deg);
    transform: rotate(-60deg)
  }

  35% {
    -webkit-transform: rotate(0deg);
    transform: rotate(0deg)
  }

  40% {
    -webkit-transform: rotate(-40deg);
    transform: rotate(-40deg)
  }

  45% {
    -webkit-transform: rotate(-60deg);
    transform: rotate(-60deg)
  }

  50% {
    -webkit-transform: rotate(10deg);
    transform: rotate(10deg)
  }

  60% {
    -webkit-transform: rotate(-40deg);
    transform: rotate(-40deg)
  }
}

@-webkit-keyframes right_finger {
  0% {
    -webkit-transform: translateY(16px);
    transform: translateY(16px)
  }

  10% {
    -webkit-transform: none;
    transform: none
  }

  50% {
    -webkit-transform: none;
    transform: none
  }

  60% {
    -webkit-transform: translateY(16px);
    transform: translateY(16px)
  }
}

@keyframes right_finger {
  0% {
    -webkit-transform: translateY(16px);
    transform: translateY(16px)
  }

  10% {
    -webkit-transform: none;
    transform: none
  }

  50% {
    -webkit-transform: none;
    transform: none
  }

  60% {
    -webkit-transform: translateY(16px);
    transform: translateY(16px)
  }
}

@-webkit-keyframes left_finger {
  0% {
    -webkit-transform: scaleX(0);
    transform: scaleX(0)
  }

  10% {
    -webkit-transform: scaleX(1) rotate(6deg);
    transform: scaleX(1) rotate(6deg)
  }

  15% {
    -webkit-transform: scaleX(1) rotate(0deg);
    transform: scaleX(1) rotate(0deg)
  }

  20% {
    -webkit-transform: scaleX(1) rotate(8deg);
    transform: scaleX(1) rotate(8deg)
  }

  25% {
    -webkit-transform: scaleX(1) rotate(0deg);
    transform: scaleX(1) rotate(0deg)
  }

  30% {
    -webkit-transform: scaleX(1) rotate(7deg);
    transform: scaleX(1) rotate(7deg)
  }

  35% {
    -webkit-transform: scaleX(1) rotate(0deg);
    transform: scaleX(1) rotate(0deg)
  }

  40% {
    -webkit-transform: scaleX(1) rotate(5deg);
    transform: scaleX(1) rotate(5deg)
  }

  45% {
    -webkit-transform: scaleX(1) rotate(0deg);
    transform: scaleX(1) rotate(0deg)
  }

  50% {
    -webkit-transform: scaleX(1) rotate(6deg);
    transform: scaleX(1) rotate(6deg)
  }

  58% {
    -webkit-transform: scaleX(0);
    transform: scaleX(0)
  }
}

@keyframes left_finger {
  0% {
    -webkit-transform: scaleX(0);
    transform: scaleX(0)
  }

  10% {
    -webkit-transform: scaleX(1) rotate(6deg);
    transform: scaleX(1) rotate(6deg)
  }

  15% {
    -webkit-transform: scaleX(1) rotate(0deg);
    transform: scaleX(1) rotate(0deg)
  }

  20% {
    -webkit-transform: scaleX(1) rotate(8deg);
    transform: scaleX(1) rotate(8deg)
  }

  25% {
    -webkit-transform: scaleX(1) rotate(0deg);
    transform: scaleX(1) rotate(0deg)
  }

  30% {
    -webkit-transform: scaleX(1) rotate(7deg);
    transform: scaleX(1) rotate(7deg)
  }

  35% {
    -webkit-transform: scaleX(1) rotate(0deg);
    transform: scaleX(1) rotate(0deg)
  }

  40% {
    -webkit-transform: scaleX(1) rotate(5deg);
    transform: scaleX(1) rotate(5deg)
  }

  45% {
    -webkit-transform: scaleX(1) rotate(0deg);
    transform: scaleX(1) rotate(0deg)
  }

  50% {
    -webkit-transform: scaleX(1) rotate(6deg);
    transform: scaleX(1) rotate(6deg)
  }

  58% {
    -webkit-transform: scaleX(0);
    transform: scaleX(0)
  }
}

@-webkit-keyframes head {
  0% {
    -webkit-transform: rotate(-3deg);
    transform: rotate(-3deg)
  }

  10% {
    -webkit-transform: translatex(10px) rotate(7deg);
    transform: translatex(10px) rotate(7deg)
  }

  50% {
    -webkit-transform: translatex(0px) rotate(0deg);
    transform: translatex(0px) rotate(0deg)
  }

  56% {
    -webkit-transform: rotate(-3deg);
    transform: rotate(-3deg)
  }
}

@keyframes head {
  0% {
    -webkit-transform: rotate(-3deg);
    transform: rotate(-3deg)
  }

  10% {
    -webkit-transform: translatex(10px) rotate(7deg);
    transform: translatex(10px) rotate(7deg)
  }

  50% {
    -webkit-transform: translatex(0px) rotate(0deg);
    transform: translatex(0px) rotate(0deg)
  }

  56% {
    -webkit-transform: rotate(-3deg);
    transform: rotate(-3deg)
  }
}

.objects {
  position: relative;
  width: 200px;
  height: 240px
}

.square {
  position: absolute;
  bottom: -60px;
  left: -5px;
  width: 120px;
  height: 120px;
  border-radius: 50%;
  -webkit-transform: rotate(-360deg);
  transform: rotate(-360deg);
  -webkit-animation: path_square 3s ease-in-out infinite;
  animation: path_square 3s ease-in-out infinite
}

.square::after {
  content: "";
  position: absolute;
  top: 10px;
  left: 0;
  width: 50px;
  height: 50px;
  background: #9ab3f5
}

.circle {
  position: absolute;
  bottom: 10px;
  left: 0;
  width: 100px;
  height: 100px;
  border-radius: 50%;
  -webkit-transform: rotate(-360deg);
  transform: rotate(-360deg);
  -webkit-animation: path_circle 3s ease-in-out infinite;
  animation: path_circle 3s ease-in-out infinite
}

.circle::after {
  content: "";
  position: absolute;
  bottom: -20px;
  left: 25px;
  width: 50px;
  height: 50px;
  border-radius: 50%;
  background: #c56183
}

.triangle {
  position: absolute;
  bottom: -62px;
  left: -20px;
  width: 110px;
  height: 110px;
  border-radius: 50%;
  -webkit-transform: rotate(-360deg);
  transform: rotate(-360deg);
  -webkit-animation: path_triangle 3s ease-in-out infinite;
  animation: path_triangle 3s ease-in-out infinite
}

.triangle::after {
  content: "";
  position: absolute;
  top: 0;
  right: -20px;
  width: 0;
  height: 0;
  border-style: solid;
  border-width: 0 28px 48px 28px;
  border-color: transparent transparent #89beb3 transparent
}

@-webkit-keyframes path_circle {
  0% {
    -webkit-transform: translateY(0);
    transform: translateY(0)
  }

  10% {
    -webkit-transform: translateY(-100px) rotate(-5deg);
    transform: translateY(-100px) rotate(-5deg)
  }

  55% {
    -webkit-transform: translateY(-100px) rotate(-360deg);
    transform: translateY(-100px) rotate(-360deg)
  }

  58% {
    -webkit-transform: translateY(-100px) rotate(-360deg);
    transform: translateY(-100px) rotate(-360deg)
  }

  63% {
    -webkit-transform: rotate(-360deg);
    transform: rotate(-360deg)
  }
}

@keyframes path_circle {
  0% {
    -webkit-transform: translateY(0);
    transform: translateY(0)
  }

  10% {
    -webkit-transform: translateY(-100px) rotate(-5deg);
    transform: translateY(-100px) rotate(-5deg)
  }

  55% {
    -webkit-transform: translateY(-100px) rotate(-360deg);
    transform: translateY(-100px) rotate(-360deg)
  }

  58% {
    -webkit-transform: translateY(-100px) rotate(-360deg);
    transform: translateY(-100px) rotate(-360deg)
  }

  63% {
    -webkit-transform: rotate(-360deg);
    transform: rotate(-360deg)
  }
}

@-webkit-keyframes path_square {
  0% {
    -webkit-transform: translateY(0);
    transform: translateY(0)
  }

  10% {
    -webkit-transform: translateY(-155px) translatex(-15px) rotate(10deg);
    transform: translateY(-155px) translatex(-15px) rotate(10deg)
  }

  55% {
    -webkit-transform: translateY(-155px) translatex(-15px) rotate(-350deg);
    transform: translateY(-155px) translatex(-15px) rotate(-350deg)
  }

  57% {
    -webkit-transform: translateY(-155px) translatex(-15px) rotate(-350deg);
    transform: translateY(-155px) translatex(-15px) rotate(-350deg)
  }

  63% {
    -webkit-transform: rotate(-360deg);
    transform: rotate(-360deg)
  }
}

@keyframes path_square {
  0% {
    -webkit-transform: translateY(0);
    transform: translateY(0)
  }

  10% {
    -webkit-transform: translateY(-155px) translatex(-15px) rotate(10deg);
    transform: translateY(-155px) translatex(-15px) rotate(10deg)
  }

  55% {
    -webkit-transform: translateY(-155px) translatex(-15px) rotate(-350deg);
    transform: translateY(-155px) translatex(-15px) rotate(-350deg)
  }

  57% {
    -webkit-transform: translateY(-155px) translatex(-15px) rotate(-350deg);
    transform: translateY(-155px) translatex(-15px) rotate(-350deg)
  }

  63% {
    -webkit-transform: rotate(-360deg);
    transform: rotate(-360deg)
  }
}

@-webkit-keyframes path_triangle {
  0% {
    -webkit-transform: translateY(0);
    transform: translateY(0)
  }

  10% {
    -webkit-transform: translateY(-172px) translatex(10px) rotate(-10deg);
    transform: translateY(-172px) translatex(10px) rotate(-10deg)
  }

  55% {
    -webkit-transform: translateY(-172px) translatex(10px) rotate(-365deg);
    transform: translateY(-172px) translatex(10px) rotate(-365deg)
  }

  58% {
    -webkit-transform: translateY(-172px) translatex(10px) rotate(-365deg);
    transform: translateY(-172px) translatex(10px) rotate(-365deg)
  }

  63% {
    -webkit-transform: rotate(-360deg);
    transform: rotate(-360deg)
  }
}

@keyframes path_triangle {
  0% {
    -webkit-transform: translateY(0);
    transform: translateY(0)
  }

  10% {
    -webkit-transform: translateY(-172px) translatex(10px) rotate(-10deg);
    transform: translateY(-172px) translatex(10px) rotate(-10deg)
  }

  55% {
    -webkit-transform: translateY(-172px) translatex(10px) rotate(-365deg);
    transform: translateY(-172px) translatex(10px) rotate(-365deg)
  }

  58% {
    -webkit-transform: translateY(-172px) translatex(10px) rotate(-365deg);
    transform: translateY(-172px) translatex(10px) rotate(-365deg)
  }

  63% {
    -webkit-transform: rotate(-360deg);
    transform: rotate(-360deg)
  }
}

.letterProgress {
  position: relative;
  margin-top: 60px;
  width: 400px;
  height: 20px;
  background: #eee
}

.letterProgress::after {
  content: "";
  position: absolute;
  top: 0;
  left: 0;
  width: 0;
  height: 100%;
  background: #637373;
  -webkit-animation: letterProgress 1.4s linear infinite;
  animation: letterProgress 1.4s linear infinite
}

@-webkit-keyframes letterProgress {
  from {
    width: 0
  }

  to {
    width: 100%
  }
}

@keyframes letterProgress {
  from {
    width: 0
  }

  to {
    width: 100%
  }
}

.noise {
  position: absolute;
  top: 50%;
  left: 50%;
  -webkit-transform: translate(-50%, -50%);
  transform: translate(-50%, -50%);
  width: 100%;
  height: 100%;
  /* z-index: 1000; */
  /* background-color: rgb(222, 201, 162); */
  /* background-image: url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADIAAAAyCAMAAAAp4XiDAAAAUVBMVEWFhYWDg4N3d3dtbW17e3t1dXWBgYGHh4d5eXlzc3OLi4ubm5uVlZWPj4+NjY19fX2JiYl/f39ra2uRkZGZmZlpaWmXl5dvb29xcXGTk5NnZ2c8TV1mAAAAG3RSTlNAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEAvEOwtAAAFVklEQVR4XpWWB67c2BUFb3g557T/hRo9/WUMZHlgr4Bg8Z4qQgQJlHI4A8SzFVrapvmTF9O7dmYRFZ60YiBhJRCgh1FYhiLAmdvX0CzTOpNE77ME0Zty/nWWzchDtiqrmQDeuv3powQ5ta2eN0FY0InkqDD73lT9c9lEzwUNqgFHs9VQce3TVClFCQrSTfOiYkVJQBmpbq2L6iZavPnAPcoU0dSw0SUTqz/GtrGuXfbyyBniKykOWQWGqwwMA7QiYAxi+IlPdqo+hYHnUt5ZPfnsHJyNiDtnpJyayNBkF6cWoYGAMY92U2hXHF/C1M8uP/ZtYdiuj26UdAdQQSXQErwSOMzt/XWRWAz5GuSBIkwG1H3FabJ2OsUOUhGC6tK4EMtJO0ttC6IBD3kM0ve0tJwMdSfjZo+EEISaeTr9P3wYrGjXqyC1krcKdhMpxEnt5JetoulscpyzhXN5FRpuPHvbeQaKxFAEB6EN+cYN6xD7RYGpXpNndMmZgM5Dcs3YSNFDHUo2LGfZuukSWyUYirJAdYbF3MfqEKmjM+I2EfhA94iG3L7uKrR+GdWD73ydlIB+6hgref1QTlmgmbM3/LeX5GI1Ux1RWpgxpLuZ2+I+IjzZ8wqE4nilvQdkUdfhzI5QDWy+kw5Wgg2pGpeEVeCCA7b85BO3F9DzxB3cdqvBzWcmzbyMiqhzuYqtHRVG2y4x+KOlnyqla8AoWWpuBoYRxzXrfKuILl6SfiWCbjxoZJUaCBj1CjH7GIaDbc9kqBY3W/Rgjda1iqQcOJu2WW+76pZC9QG7M00dffe9hNnseupFL53r8F7YHSwJWUKP2q+k7RdsxyOB11n0xtOvnW4irMMFNV4H0uqwS5ExsmP9AxbDTc9JwgneAT5vTiUSm1E7BSflSt3bfa1tv8Di3R8n3Af7MNWzs49hmauE2wP+ttrq+AsWpFG2awvsuOqbipWHgtuvuaAE+A1Z/7gC9hesnr+7wqCwG8c5yAg3AL1fm8T9AZtp/bbJGwl1pNrE7RuOX7PeMRUERVaPpEs+yqeoSmuOlokqw49pgomjLeh7icHNlG19yjs6XXOMedYm5xH2YxpV2tc0Ro2jJfxC50ApuxGob7lMsxfTbeUv07TyYxpeLucEH1gNd4IKH2LAg5TdVhlCafZvpskfncCfx8pOhJzd76bJWeYFnFciwcYfubRc12Ip/ppIhA1/mSZ/RxjFDrJC5xifFjJpY2Xl5zXdguFqYyTR1zSp1Y9p+tktDYYSNflcxI0iyO4TPBdlRcpeqjK/piF5bklq77VSEaA+z8qmJTFzIWiitbnzR794USKBUaT0NTEsVjZqLaFVqJoPN9ODG70IPbfBHKK+/q/AWR0tJzYHRULOa4MP+W/HfGadZUbfw177G7j/OGbIs8TahLyynl4X4RinF793Oz+BU0saXtUHrVBFT/DnA3ctNPoGbs4hRIjTok8i+algT1lTHi4SxFvONKNrgQFAq2/gFnWMXgwffgYMJpiKYkmW3tTg3ZQ9Jq+f8XN+A5eeUKHWvJWJ2sgJ1Sop+wwhqFVijqWaJhwtD8MNlSBeWNNWTa5Z5kPZw5+LbVT99wqTdx29lMUH4OIG/D86ruKEauBjvH5xy6um/Sfj7ei6UUVk4AIl3MyD4MSSTOFgSwsH/QJWaQ5as7ZcmgBZkzjjU1UrQ74ci1gWBCSGHtuV1H2mhSnO3Wp/3fEV5a+4wz//6qy8JxjZsmxxy5+4w9CDNJY09T072iKG0EnOS0arEYgXqYnXcYHwjTtUNAcMelOd4xpkoqiTYICWFq0JSiPfPDQdnt+4/wuqcXY47QILbgAAAABJRU5ErkJggg==) */
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
  justify-content: center;
  align-items: center;
  opacity: 1;
}

::v-deep .el-progress-bar__inner {
  position: absolute;
  height: 100%;
  background: linear-gradient(90deg, #ffa800 0%, #ffdd86 100%);
}

::v-deep .el-progress-bar__innerText {
  display: block;
  width: 35px;
  height: 30px !important;
  float: right;
  background-image: url("../../assets/imgs/car1.png");
  background-size: 100% 100%;
  margin: -3px 3px;
  position: absolute;
  right: 0;
  font-size: 0;
}

/* 熊猫按钮 */

span {
  margin-left: 12px;
}

.panda-button {
  width: 180px;
  padding: 10px 14px;
  display: flex;
  justify-content: center;
  align-items: center;
  border: solid rgb(131, 88, 59);
  border-width: 2px 2px 7px 2px;
  border-radius: 10px;
  background-color: rgb(83, 39, 9);
  font-family: 'Press Start 2P', cursive;
  font-size: 15px;
  color: white;
  transition: 0.3s;
  cursor: pointer;
  position: absolute;
  bottom: 10px;
  right: 10px;
}

.panda-button:hover {
  border-width: 2px 2px 4px 2px;
  margin-top: 3px;
}

.panda-button:focus {
  outline: none;
  box-shadow: 0 0 0 4px hotpink;
}

.panda_face {
  position: relative;
  height: 25px;
  width: auto;
}

.panda_face {
  opacity: 1;
}

.panda_face svg {
  height: 100%;
  width: auto;
}

.panda_white {
  fill: white;
}

.pink {
  fill: hotpink;
}

.panda-button:hover .panda_face {
  animation: bop infinite 0.5s;
}

@keyframes bop {
  0% {
    top: 0;
    transform: scale(1);
  }

  50% {
    top: -5px;
    transform: scale(0.95, 1.05);
  }

  90% {
    transform: scale(1.1, 0.9);
  }

  100% {
    top: 0;
    transform: scale(1);
  }
}

.panda-button.animate svg {
  animation: spin infinite linear 0.5s;
}

@keyframes spin {
  0% {
    transform: rotate(0);
  }

  100% {
    transform: rotate(-720deg);
  }
}

@keyframes fadeout {
  0% {
    opacity: 1;
  }

  90% {
    opacity: 0;
  }

  100% {
    opacity: 0;
    display: none;
  }
}

.panda-button:hover .fade_in,
.panda_face.fade_in {
  animation: fade_in forwards 0.5s;
}

@keyframes fade_in {
  0% {
    opacity: 0;
  }

  100% {
    opacity: 1;
  }
}

.panda-button.animate {
  border-color: rgb(58, 28, 9);
  transform-origin: bottom right;
  animation: shake ease forwards 0.5s;
}

@keyframes shake {
  0% {
    transform-origin: bottom right;
    transform: translate(0, 0) rotate(0deg);
  }

  40% {
    transform: translate(5px, 0) rotate(45deg);
  }

  70% {
    transform-origin: bottom right;
    transform: translate(5px, 0) rotate(45deg);
  }

  85% {
    transform-origin: bottom left;
    transform: translate(0, 0) rotate(-20deg);
  }

  95% {
    transform: translate(0, 0) rotate(0deg);
  }

  100% {
    transform-origin: bottom right;
    transform: translate(0, 0) rotate(0deg);
  }
}

/* 熊猫按钮结束 */

.c-toggle {
  display: inline-flex;
  align-items: center;
  margin-left: 70px;
  margin-bottom: 50px;
}

.c-toggle__item {
  cursor: pointer;
  border: 1px solid #1a78d4;
  color: #1a78d4;
  display: inline-block;
  min-width: 10.4rem;
  padding: 0.64rem 0.8rem;
  opacity: 0.75;
  transition: all 200ms ease-out 0s;
}

.c-toggle__item:hover {
  background: rgba(26, 120, 212, 0.15);
}

.c-toggle__item:active {
  box-shadow: inset 0px 0px 0.4rem #1a78d4;
}

.c-toggle__item:first-of-type {
  border-top-left-radius: 0.64rem;
  border-bottom-left-radius: 0.64rem;
}

.c-toggle__item:last-of-type {
  border-top-right-radius: 0.64rem;
  border-bottom-right-radius: 0.64rem;
}

.c-toggle__item:not(:first-of-type) {
  border-left: 0;
}

.c-toggle__item--checked {
  background: rgba(26, 120, 212, 0.15);
  font-weight: 700;
  opacity: 1;
}
</style>
