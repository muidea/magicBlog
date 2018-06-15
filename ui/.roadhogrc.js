const path = require('path')
const { version } = require('./package.json')

export default {
  entry: "src/index.js",
  publicPath: `/${version}/`,
  outputPath: `./dist/${version}`,
  env: {
    development: {
      extraBabelPlugins: [
        "dva-hmr",
        "transform-runtime",
        ["import", { "libraryName": "antd", "style": "css" }],
      ]
    },
    production: {
      extraBabelPlugins: [
        "transform-runtime",
        ["import", { "libraryName": "antd", "style": "css" }],
      ]
    }
  },
  dllPlugin: {
    exclude: ["babel-runtime", "roadhog", "cross-env"],
    include: ["dva/router", "dva/saga", "dva/fetch"]
  },
   proxy: {
    "/api/v1/": {
      "target": "http://localhost:8866/",
      "changeOrigin": true,
      "pathRewrite": { "^/api/v1/": "/" }
    },
  }
}
