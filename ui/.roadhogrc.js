export default {
  entry: "src/index.js",
  env: {
    development: {
      extraBabelPlugins: [
        "dva-hmr",
        "transform-runtime",
        ["import", { "libraryName": "antd", "style": "css" }],
        ["module-resolver", {
          "alias": {
            "routes": `${__dirname}/src/routes`,
            "models": `${__dirname}/src/models`,
            "services": `${__dirname}/src/services`,
            "utils": `${__dirname}/src/utils`
          }
        }]
      ]
    },
    production: {
      extraBabelPlugins: [
        "transform-runtime",
        ["import", { "libraryName": "antd", "style": "css" }],
        ["module-resolver", {
          "alias": {
            "routes": `${__dirname}/src/routes`,
            "models": `${__dirname}/src/models`,
            "services": `${__dirname}/src/services`,
            "utils": `${__dirname}/src/utils`
          }
        }]
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
