window.__plugin_my_plugin__ = function(api) {
  api.log("info", "插件已加载");

  // 注册一个组件到侧边栏顶部
  api.registerSlot("sidebar-top", function MyWidget() {
    return React.createElement(
      "div",
      { style: { padding: "12px", background: "#f0f9ff", borderRadius: "8px" } },
      "👋 Hello from My Plugin!"
    );
  });
};