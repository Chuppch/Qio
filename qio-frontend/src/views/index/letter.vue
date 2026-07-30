<template>
  <div class="banner">
    <div class="box" v-for="item in items" :key="item.title">
      <div class="image-container">
        <img :src="item.img" alt="">
      </div>
      <div class="content">
        <!-- <h3>{{ item.title }}</h3> -->
        <p v-show="item.active">
          {{ item.description }}
        </p>
        <router-link :to="item.link" tag="button" class="btn">
          {{ item.btnText }}
        </router-link>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'IndexLetter',
  data() {
    return {
      items: [
        {
          title: '书海情航',
          description: '在书信的海洋中航行，寻找情感的归宿。在这些泛黄的纸张上，每一行文字都承载着深情的寄托，像是一艘艘小船，载着思念和故事，穿越时间的波涛，抵达心灵的港湾。',
          img: require('@/assets/imgss/write.webp'),
          active: true,
          btnText: '我要写信',
          link: '/write'
        },
        {
          title: '书简候音',
          description: '书简中等候着远方的声音。那些被墨水封印的字句，静静躺在信封中，仿佛在低语，诉说着远方的牵挂和期盼。它们在静默中等待，直到被打开的那一刻，释放出所有被压抑的情感和消息。',
          img: require('@/assets/imgss/piaoliu.webp'),
          active: true,
          btnText: '查看收信',
          link: '/receive'
        },
        {
          title: '信风远航',
          description: '借助风力，信件远航到未知的地方。每一封信都是一次勇敢的航行，它们随风起航，越过千山万水，带着希望和梦想，去往未知的目的地。无论路途多远，它们总能找到自己的方向，将心意送达期待的人手中。',
          img: require('@/assets/imgss/send.webp'),
          active: true,
          btnText: '写漂流信',
          link: '/drifting'
        }
      ]
    }
  }
}
</script>

<style scoped>
.banner {
  width: 1000px;
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  margin: 100px auto;
}

.box {
  position: relative;
  width: 300px;
  height: 400px;
  /* margin-top: 150px; */
  box-shadow: 0 15px 15px rgba(0, 0, 0, 0.5);
}

/* .banner:before {
  content: "";
  width: 100%;
  display: block;
  position: absolute;
  background: url(https://www.taoyuantudigong.org.tw/main/wp-content/themes/project-theme/src/img/bannerbg.png) 0 0 / auto 100% repeat-x;
  height: 100%;
  top: 0;
} */

.image-container {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}

.image-container img {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: opacity 0.5s;
}

.box:hover .image-container img {
  opacity: 0;
}

.content {
  position: absolute;
  bottom: 20px;
  left: 10%;
  width: 80%;
  height: 60px;
  transition: all 0.5s;
  overflow: hidden;
  padding: 15px;
  box-sizing: border-box;
}

.box:hover .content {
  width: 100%;
  height: 100%;
  bottom: 0;
  left: 0;
  background-color: #f0efed;
  border-radius: 25px;
  cursor: pointer;
}

.box:hover {
  box-shadow: 0px rgba(0, 0, 0, 0);
}

.box:hover .btn {
  opacity: 1;
  transition-delay: 0.1s;
}

.content p {
  position: absolute;
  top: 40%;
  left: 10px;
  transform: translateY(-50%);
  padding: 0;
  opacity: 0;
  line-height: 1.6em;
  transition: opacity 0s;
  text-align: center;
  font-size: 20px;
}

.content p a {
  color: #222;
  text-decoration: none;
}

.box:hover .content p {
  opacity: 1;
  transition-delay: 0.1s;
}

.btn {
  display: inline-block;
  position: absolute;
  bottom: 50px;
  left: 50px;
  width: 200px;
  height: 50px;
  border-radius: 50px;
  background-color: #8b4d38;
  border: 0px;
  color: white;
  opacity: 0;

}
</style>
