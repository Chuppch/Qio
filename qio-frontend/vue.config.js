const { defineConfig } = require('@vue/cli-service')
module.exports = defineConfig({
  transpileDependencies: true
})
// 在vue-config.js 中加入
const CompressionWebpackPlugin = require('compression-webpack-plugin')
const productionGzipExtensions = ['js', 'css']
const isProduction = process.env.NODE_ENV === 'production'

module.exports = {
  configureWebpack: config => {
    if (isProduction) {
      config.plugins.push(new CompressionWebpackPlugin({
        algorithm: 'gzip',
        test: new RegExp('\\.(' + productionGzipExtensions.join('|') + ')$'),
        threshold: 10240,
        minRatio: 0.8
      }))
    }
  }
}
module.exports = {
  devServer: {
    client: { overlay: { runtimeErrors: false } }, // 处理开发环境的全屏异常报错
    host: 'localhost',
    port: 80,
    proxy: {
      '/api': {
        // 设置代理的目标URL
        target: 'http://localhost:8080/',
        changeOrigin: true,
        pathRewrite: {
          '^/api': '' // 将请求中的/api路径替换为空
        }
      }
    }
  }
}
