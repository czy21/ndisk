const webpackConfigPlugin = require("./webpack.config")
const path = require("path");

module.exports = {
    eslint: {
        enable: false
    },
    webpack: {
        alias: {
            "@": path.resolve(__dirname, "src/")
        }
    },
    plugins: [
        {plugin: webpackConfigPlugin, options: {preText: "Will log the webpack config:"}}
    ],
};