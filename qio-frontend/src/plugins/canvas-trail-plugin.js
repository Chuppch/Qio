// canvas-trail-plugin.js
export default {
  install(Vue) {
    const canvas = document.createElement('canvas')
    document.body.appendChild(canvas)
    const ctx = canvas.getContext('2d')

    let mouseMoved = false
    const pointer = {
      x: 0.5 * window.innerWidth,
      y: 0.5 * window.innerHeight
    }
    const params = {
      pointsNumber: 40,
      widthFactor: 0.3,
      mouseThreshold: 0.6,
      spring: 0.4,
      friction: 0.5
    }
    const trail = new Array(params.pointsNumber).fill(null).map(() => ({
      x: pointer.x,
      y: pointer.y,
      dx: 0,
      dy: 0
    }))

    function updateMousePosition(eX, eY) {
      pointer.x = eX
      pointer.y = eY
    }

    function setupCanvas() {
      canvas.width = window.innerWidth
      canvas.height = window.innerHeight
    }

    function update(t) {
      if (!mouseMoved) {
        pointer.x = (0.5 + 0.3 * Math.cos(0.002 * t) * (Math.sin(0.005 * t))) * window.innerWidth
        pointer.y = (0.5 + 0.2 * (Math.cos(0.005 * t)) + 0.1 * Math.cos(0.01 * t)) * window.innerHeight
      }

      ctx.clearRect(0, 0, canvas.width, canvas.height)
      trail.forEach((p, pIdx) => {
        const prev = pIdx === 0 ? pointer : trail[pIdx - 1]
        const spring = pIdx === 0 ? 0.4 * params.spring : params.spring
        p.dx += (prev.x - p.x) * spring
        p.dy += (prev.y - p.y) * spring
        p.dx *= params.friction
        p.dy *= params.friction
        p.x += p.dx
        p.y += p.dy
      })

      ctx.lineCap = 'round'
      ctx.beginPath()
      ctx.moveTo(trail[0].x, trail[0].y)

      for (let i = 1; i < trail.length - 1; i++) {
        const xc = 0.5 * (trail[i].x + trail[i + 1].x)
        const yc = 0.5 * (trail[i].y + trail[i + 1].y)
        ctx.quadraticCurveTo(trail[i].x, trail[i].y, xc, yc)
        ctx.lineWidth = params.widthFactor * (params.pointsNumber - i)
        ctx.stroke()
      }
      ctx.lineTo(trail[trail.length - 1].x, trail[trail.length - 1].y)
      ctx.stroke()

      window.requestAnimationFrame(update)
    }

    setupCanvas()
    window.addEventListener('click', (e) => updateMousePosition(e.pageX, e.pageY))
    window.addEventListener('mousemove', (e) => {
      mouseMoved = true
      updateMousePosition(e.pageX, e.pageY)
    })
    window.addEventListener('touchmove', (e) => {
      mouseMoved = true
      updateMousePosition(e.targetTouches[0].pageX, e.targetTouches[0].pageY)
    })
    window.addEventListener('resize', setupCanvas)
    window.requestAnimationFrame(update)

    // Cleanup on component destroy
    Vue.mixin({
      beforeDestroy() {
        window.removeEventListener('click', (e) => updateMousePosition(e.pageX, e.pageY))
        window.removeEventListener('mousemove', (e) => {
          mouseMoved = true
          updateMousePosition(e.pageX, e.pageY)
        })
        window.removeEventListener('touchmove', (e) => {
          mouseMoved = true
          updateMousePosition(e.targetTouches[0].pageX, e.targetTouches[0].pageY)
        })
        window.removeEventListener('resize', setupCanvas)
        document.body.removeChild(canvas)
      }
    })
  }
}
