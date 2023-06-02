const { defineConfig } = require("@vue/cli-service");
module.exports = defineConfig({
  transpileDependencies: true,
  outputDir: "../assets/_dash/",
  publicPath: "/_dash/",
});
