<template>
  <canvas ref="canvas"></canvas>
</template>

<script>
export default {
  name: 'CanvasAnimation',
  data() {
    return {
      pointer: {
        x: 0.5 * window.innerWidth,
        y: 0.5 * document.documentElement.clientHeight // 也可以在这里使用 document.documentElement.clientHeight
      },
      params: {
        pointsNumber: 20,
        widthFactor: 0.3,
        mouseThreshold: 0.6,
        spring: 0.4,
        friction: 0.5
      },
      trail: [],
      mouseMoved: false,
      ctx: null
    }
  },
  mounted() {
    this.initCanvas()
    window.addEventListener('click', this.updateMousePosition)
    window.addEventListener('mousemove', this.updateMousePosition)
    window.addEventListener('touchmove', this.updateMousePosition)
    window.addEventListener('resize', this.resizeCanvas)
    this.$refs.canvas.addEventListener('click', this.handleCanvasClick)
  },
  beforeDestroy() {
    window.removeEventListener('click', this.updateMousePosition)
    window.removeEventListener('mousemove', this.updateMousePosition)
    window.removeEventListener('touchmove', this.updateMousePosition)
    window.removeEventListener('resize', this.resizeCanvas)
    this.$refs.canvas.removeEventListener('click', this.handleCanvasClick)
  },
  methods: {
    handleCanvasClick(event) {
      // console.log('Canvas clicked', event)
    },
    initCanvas() {
      const canvas = this.$refs.canvas
      this.ctx = canvas.getContext('2d')
      canvas.width = window.innerWidth
      canvas.height = document.documentElement.clientHeight
      this.trail = new Array(this.params.pointsNumber)
      for (let i = 0; i < this.params.pointsNumber; i++) {
        this.trail[i] = {
          x: this.pointer.x,
          y: this.pointer.y,
          dx: 0,
          dy: 0
        }
      }
      this.animate(0)
    },
    animate(t) {
      if (!this.mouseMoved) {
        this.pointer.x = (0.5 + 0.3 * Math.cos(0.002 * t) * (Math.sin(0.005 * t))) * window.innerWidth
        this.pointer.y = (0.5 + 0.2 * (Math.cos(0.005 * t)) + 0.1 * Math.cos(0.01 * t)) * document.documentElement.clientHeight
      }
      this.ctx.clearRect(0, 0, this.$refs.canvas.width, this.$refs.canvas.height)
      this.trail.forEach((p, pIdx) => {
        const prev = pIdx === 0 ? this.pointer : this.trail[pIdx - 1]
        const spring = pIdx === 0 ? 0.4 * this.params.spring : this.params.spring
        p.dx += (prev.x - p.x) * spring
        p.dy += (prev.y - p.y) * spring
        p.dx *= this.params.friction
        p.dy *= this.params.friction
        p.x += p.dx
        p.y += p.dy
      })
      this.ctx.lineCap = 'round'
      this.ctx.beginPath()
      this.ctx.moveTo(this.trail[0].x, this.trail[0].y)
      for (let i = 1; i < this.trail.length - 1; i++) {
        const xc = 0.5 * (this.trail[i].x + this.trail[i + 1].x)
        const yc = 0.5 * (this.trail[i].y + this.trail[i + 1].y)
        this.ctx.quadraticCurveTo(this.trail[i].x, this.trail[i].y, xc, yc)
        this.ctx.lineWidth = this.params.widthFactor * (this.params.pointsNumber - i)
        this.ctx.stroke()
      }
      this.ctx.lineTo(this.trail[this.trail.length - 1].x, this.trail[this.trail.length - 1].y)
      this.ctx.stroke()
      requestAnimationFrame(this.animate)
    },
    updateMousePosition(e) {
      this.mouseMoved = true
      this.pointer.x = e.clientX
      this.pointer.y = e.clientY
    },
    resizeCanvas() {
      const canvas = this.$refs.canvas
      canvas.width = window.innerWidth
      canvas.height = document.documentElement.clientHeight
    }
  }
}
</script>

<style scoped>
canvas {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 999;
  pointer-events: none;
}
</style>
