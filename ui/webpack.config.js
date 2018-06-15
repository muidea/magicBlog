const HtmlWebpackPlugin = require('html-webpack-plugin')
const CopyWebpackPlugin = require('copy-webpack-plugin')
const UglifyJsPlugin = require('uglifyjs-webpack-plugin')
const webpack = require('webpack')

module.exports = (webpackConfig, env) => {
  const production = env === 'production'
  // FilenameHash
  webpackConfig.output.chunkFilename = '[name].[chunkhash].js'

  if (production) {
    if (webpackConfig.module) {
    // ClassnameHash
      webpackConfig.module.rules.map((item) => {
        if (String(item.test) === '/\\.less$/' || String(item.test) === '/\\.css/') {
          item.use.filter(iitem => iitem.loader === 'css')[0].options.localIdentName = '[hash:base64:5]'
        }
        return item
      })
    }
    webpackConfig.plugins.push(
      new webpack.LoaderOptionsPlugin({
        minimize: true,
        debug: false,
      }),
    )
  }

  webpackConfig.plugins = webpackConfig.plugins.concat([
    new UglifyJsPlugin(),
    new CopyWebpackPlugin([
      {
        from: 'src/public',
        to: webpackConfig.output.outputPath,
      },
    ]),
    new HtmlWebpackPlugin({
      template: `${__dirname}/src/index.ejs`,
      filename: 'index.html',
      minify: {
        collapseWhitespace: true,
      },
      hash: true,
      headScripts: ['/roadhog.dll.js'],
    }),
  ])

  // Alias
  webpackConfig.resolve.alias = {
    components: `${__dirname}/src/components`,
    utils: `${__dirname}/src/utils`,
    services: `${__dirname}/src/services`,
    models: `${__dirname}/src/models`,
    routes: `${__dirname}/src/routes`,
    themes: `${__dirname}/src/themes`,
  }

  return webpackConfig
}
