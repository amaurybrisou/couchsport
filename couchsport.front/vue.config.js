const path = require('path')
module.exports = {
  transpileDependencies: ['vuetify'],
  runtimeCompiler: true,
  devServer: {
    port: 8081,
    proxy: {
      '^/api': {
        target: process.env.VUE_APP_BACKEND_API_URL || 'http://localhost:9000',
        changeOrigin: true,
        secure: false,
        ws: true,
        logLevel: 'debug'
      },
      '^/uploads': {
        target:
          process.env.VUE_APP_BACKEND_STATIC_URL || 'http://localhost:9001',
        changeOrigin: true,
        secure: false,
        ws: true,
        logLevel: 'debug'
      }
    }
  },
  chainWebpack: (config) => {
    config.entry('app').add('./src/main.js').end()

    config.resolve.alias
      .set('@', path.resolve(__dirname, 'src/'))
      .set('~', path.resolve(__dirname, 'node_modules'))
      .set('pages', path.resolve(__dirname, 'src/pages'))
      .set('store', path.resolve(__dirname, 'src/store'))
      .set('repos', path.resolve(__dirname, 'src/repositories'))
      .set('css', path.resolve(__dirname, 'assets/css'))
      .set('components', path.resolve(__dirname, 'src/components'))
      .set('mixins', path.resolve(__dirname, 'src/mixins'))
      .set('plugins', path.resolve(__dirname, 'src/plugins'))
      .set('static', path.resolve(__dirname, 'static'))

    config.resolve.extensions.add('.js').add('.vue').add('.json').add('.css')

    config.plugin('html').tap((args) => {
      args[0].template = path.resolve(__dirname, 'index.html')
      args[0].title = 'Couchsport'
      args[0].meta = {
        viewport:
          'width=device-width, initial-scale=0.5, maximum-scale=1, minimal-ui'
      }
      return args
    })
  }
}
