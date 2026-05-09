// 本地插件服务端逻辑（示例）
module.exports = {
  // 插件安装时执行
  install: (context) => {
    console.log("[Local Plugin] 服务端安装，数据库可使用 context.db");
  },
  // 插件卸载时执行
  uninstall: (context) => {
    console.log("[Local Plugin] 服务端卸载");
  },
  // 启用插件
  enable: (context) => {
    console.log("[Local Plugin] 服务端启用");
  },
  // 禁用插件
  disable: (context) => {
    console.log("[Local Plugin] 服务端禁用");
  },
};
