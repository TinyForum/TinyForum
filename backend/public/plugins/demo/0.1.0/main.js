// my-plugin.js（打包后上传到 CDN）
window.__plugin_my_plugin__ = function (api) {
  // 注册一个 React 组件到侧边栏顶部插槽
  api.registerSlot("sidebar-top", function MyWidget() {
    return React.createElement(
      "div",
      {
        className: "card bg-base-100 border border-base-300 p-3 text-sm",
      },
      "👋 Hello from My Plugin!",
    );
  });

  // 监听事件
  api.on("post:view", function (data) {
    api.log("info", "User viewed a post: " + JSON.stringify(data));
  });
};
