module.exports = {
  outputDir: "site",
  devServer: {
    proxy: {
      "/api": {
        target: "http://localhost:6612",
        changeOrigin: true,
      },
    },
  },
  publicPath: process.env.NODE_ENV === "production" ? "/ss/" : "/",
};
