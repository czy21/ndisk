const webpackConfigPlugin = require("./webpack.config")
const {CracoAliasPlugin, configPaths} = require('react-app-rewire-alias')

module.exports = {
    eslint: {
        enable: false
    },
    plugins: [
        {
            plugin: CracoAliasPlugin,
            options: {alias: configPaths('./tsconfig.extend.json')}
        },
        {plugin: webpackConfigPlugin, options: {preText: "Will log the webpack config:"}}
    ],
};