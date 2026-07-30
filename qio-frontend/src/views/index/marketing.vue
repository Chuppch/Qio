<template>
  <div class="banner">
    <div class="content" data-spm="item" style="">
      <div class="contentInner " v-infinite-scroll="list" infinite-scroll-immediate="false"
        infinite-scroll-disabled="disabledLoad" style="transform: translate3d(-4px, 0px, 0px);">
        <a class="doubleCardbanner--_6NpK_ey" v-for="item in data" :key="item.id" :href="item.link" target="_blank"
          data-spm="1" data-appeared="false" data-has-appeared="true" data-before-current-y="-414.22222900390625"
          style="min-height: 370px;" data-has-disappeared="true">
          <div class="doubleCard--gO3Bz6bu">
            <div class="mainPicAndDesc--Q5PYrWux">
              <div class=" mainPicbanner--qRLTAeii" data-name="itemExp" data-aplus-ae="x1_29da46f"
                data-spm-anchor-id="a21n57.1.1.i0.4330523cXUTEuN">
                <el-image :src="item.image" height="240" width="240" lazy class="mainPic--Ds3X7I8z">
                  <div slot="placeholder" class="mainPic--Ds3X7I8z">
                    <img :src="item.image" height="240" width="240" class="mainPic--Ds3X7I8z">
                  </div>
                </el-image>

                <!-- 广告标签 -->
                <img class="mainP4pPic--jbnK3QAX"
                  src="https://img.alicdn.com/imgextra/i3/O1CN01njPHBL1iDAvy7mi2J_!!6000000004378-2-tps-42-27.png">
              </div>
              <div class="descbanner--Ta96FeyX normalMod--c15U6g3b">
                <div class="title--qJ7Xg_90 "><img class="iconPic--pZyhTNNV"
                    src="https://img.alicdn.com/imgextra/i1/O1CN01IGJRXF1fome45FX8K_!!6000000004054-2-tps-71-42.png"
                    style="padding-left: 0px;"><span class="">{{ item.name }}</span></div>
              </div>
              <div class="container--dWqrHwI7 newContainer--oW1_6Dcr " style="height: 22px; margin-top: 4px;">
                <div class="abstractbanner--whLX5va5" style="margin-left: 8px;">
                  <div class="descBox--RunOO4S3" style="height: 22px; width: auto;"><span class="text--eAiSCa_r"
                      style="color: rgb(122, 122, 122);">{{ item.marketing }} </span></div>
                  <!-- <div class="descBox--RunOO4S3" style="height: 22px; width: auto;"><span class="text--eAiSCa_r"
                      style="color: rgb(122, 122, 122);">卡通动漫</span></div>
                  <div class="descBox--RunOO4S3" style="height: 22px; width: auto;"><span class="text--eAiSCa_r"
                      style="color: rgb(122, 122, 122);">卡通</span></div> -->
                </div>
              </div>
              <div class="pricebanner--dBtPZ2K1 "><span class="unit--D3KGoZe2">¥</span>
                <div style="margin-right: 8px;"><span class="priceInt--yqqZMJ5a">{{ item.price }}</span>
                </div>
              </div>
            </div>
          </div>
        </a>
        <p v-if="loading" style="width: 100%; color: #ff6200; margin-top: -2px; font-size: 18px;"><i
            class="el-icon-loading"></i>加载中...
        </p>
        <p v-if="noMore" style="width: 100%; color: #ff6200; margin-top: -2px; font-size: 18px;"><i
            class="el-icon-shopping-cart-2"></i>没有更多了</p>
      </div>
    </div>
    <!-- <vue-lazyload v-if="data.length">
      <img v-lazy="item.image" height="240" width="240" class="mainPic--Ds3X7I8z">
    </vue-lazyload> -->
  </div>
</template>

<script>
import { list } from '@/api/marketing'

export default {
  name: 'QiaopiMarketing',
  data() {
    return {
      data: [],
      page: 1,
      limit: 20,
      noMore: false,
      loading: false
    }
  },
  computed: {
    disabledLoad() {
      return this.loading || this.noMore
    }
  },
  methods: {
    list() {
      this.loading = true
      list(
        this.page,
        this.limit
      ).then(res => {
        // console.log(res.data.records.length)
        if (res.data.records.length === 0) {
          this.loading = false
          this.noMore = true
          return
        }
        this.data = [...this.data, ...res.data.records]
        this.page++
        this.loading = false
      })
    }
  },
  mounted() {
    this.list()
  }
}
</script>

<style scoped>
.banner {
  position: relative;
  margin-top: 40px;
  width: 85%;
  height: 85vh;
  background-position: center center;
  background-size: cover;
  z-index: 1;
}

/* 滑动轨道按钮 */
::-webkit-scrollbar-button {
  width: 10px;
  height: 10px;
  display: none;
}

.content {
  background-color: #fff;
  width: 100%;
  /* height: 100%; */
  height: 80vh;
  font-family: Roboto, "Helvetica Neue", Helvetica, Tahoma, Arial, "PingFang SC", "Microsoft YaHei";
  font-size: 12px;
  line-height: 1.28571;
  color: #111111;
  box-sizing: border-box;
  border-radius: 10px;
  padding: 20px;
}

.contentInner {
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  justify-content: space-between;
  transform: translate3d(-4px, 0px, 0px);
  overflow-y: scroll;
  height: 75vh;
  /* Adjust height to account for padding */
}

a:link {
  color: #2481e6;
}

.doubleCardbanner--_6NpK_ey,
.doubleCardbannerAdapt--mEcC7olq {
  border: 1px solid hsla(0, 0%, 100%, 0);
  border-radius: 12px;
  box-sizing: border-box;
  display: block;
  height: auto;
  justify-content: center;
  margin-bottom: 12px;
  overflow: hidden;
  padding: 3px;
  text-decoration: none !important;
  transition: width .5s, height .5s, background-color .5s, transform .5s;
  width: 248px;
}

a {
  text-decoration: none;
  background-color: transparent;
}

.doubleCard--gO3Bz6bu {
  display: flex;
  flex-direction: column;
  height: auto;
}

.doubleCardbanner--_6NpK_ey:hover {
  border: 1px solid #ff6200;
  border-radius: 12px;
}

.mainPicAndDesc--Q5PYrWux {
  align-items: center;
  color: #11192d;
  display: flex;
  flex-direction: column;
  text-decoration: none;
}

.mainPicAdaptbanner--V_ayd2hD,
.mainPicbanner--qRLTAeii {
  border-radius: 8px;
  flex-shrink: 0;
  height: 240px;
  overflow: hidden;
  position: relative;
  width: 240px;
}

.mainPic--Ds3X7I8z {
  height: 240px;
  /* -o-object-fit: cover; */
  /* object-fit: cover; */
  width: 240px;
}

.mainP4pPic--jbnK3QAX {
  bottom: 0;
  height: 18px;
  position: absolute;
  right: 0;
  width: auto;
}

img {
  border-style: none;
}

.normalMod--c15U6g3b {
  padding-left: 8px;
  padding-right: 8px;
}

.descbanner--Ta96FeyX {
  display: flex;
  flex-direction: row;
  margin-top: 8px;
  text-decoration: none;
  width: 100%;
}

.title--qJ7Xg_90 {
  -webkit-box-orient: vertical;
  word-wrap: break-word;
  -webkit-line-clamp: 2;
  /* 兼容旧版Webkit内核浏览器 */
  line-clamp: 2;
  /* 标准属性，提高兼容性 */
  color: #11192d;
  display: -webkit-box;
  font-family: PingFangSC-Medium;
  font-size: 16px;
  font-weight: 500;
  height: 48px;
  line-height: 24px;
  overflow: hidden;
  text-overflow: ellipsis;
  vertical-align: middle;
  width: 100%;
  word-break: break-all;
}

.iconPic--pZyhTNNV {
  padding-left: 0px;
  height: 16px;
  margin-bottom: 3px;
  padding-right: 6px;
  vertical-align: text-bottom;
}

.newContainer--oW1_6Dcr {
  height: 22px;
  margin-top: 4px;
  align-items: center;
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  justify-content: flex-start;
  overflow: hidden;
}

.container--dWqrHwI7 {
  display: flex;
  height: 48px;
  padding-right: 8px;
  width: 100%;
}

.abstractbanner--whLX5va5,
.secondKillInfobanner--sf5GHTt6 {
  align-items: center;
  display: flex;
  flex-direction: row;
  justify-content: flex-start;
  overflow: hidden;
  width: 100%;
}

.abstractbanner--whLX5va5 {
  flex-wrap: wrap;
  font-family: PingFang SC;
  max-height: 22px;
  margin-left: 8px;
}

.descBox--RunOO4S3 {
  color: #11192d;
  display: flex;
  flex-direction: column;
  font-size: 14px;
  justify-content: space-between;
  justify-content: center;
  line-height: 22px;
  margin-right: 8px;
  width: 56px;
  height: 22px;
  width: auto;
}

.text--eAiSCa_r {
  font-size: 14px;
  line-height: 22px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: rgb(122, 122, 122);

}

.pricebanner--dBtPZ2K1 {
  align-items: flex-end;
  box-sizing: border-box;
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  height: 36px;
  max-height: 36px;
  overflow: hidden;
  padding-left: 8px;
  padding-right: 8px;
  width: 100%;
}

.unit--D3KGoZe2 {
  color: #ff6200;
  font-family: AlibabaSans102v1TaoBao-Bold, Alibaba-Sans102, PingFangSC-Regular;
  font-size: 16px;
  font-weight: 700;
  line-height: 16px;
  margin-bottom: 7px;
}

.price--GwpwcSLr,
.priceFloat--XpixvyQ1,
.priceInt--yqqZMJ5a {
  color: #ff6200;
  font-family: AlibabaSans102v1TaoBao-Bold, Alibaba-Sans102, PingFangSC-Regular;
  font-size: 24px;
  font-weight: 700;
  line-height: 36px;
  vertical-align: middle;
}

.price--GwpwcSLr,
.priceFloat--XpixvyQ1,
.priceInt--yqqZMJ5a {
  color: #ff6200;
  font-family: AlibabaSans102v1TaoBao-Bold, Alibaba-Sans102, PingFangSC-Regular;
  font-size: 24px;
  font-weight: 700;
  line-height: 36px;
  vertical-align: middle;
}

.subIconbanner--Vl8zAdQn {
  align-items: flex-start;
  align-items: center;
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  height: 16px;
  margin-left: 8px;
  margin-top: 4px;
  min-height: 16px;
  overflow: hidden;
}
</style>
