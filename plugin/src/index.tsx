import React from "react";
import { MyWidget } from "./Widget";
import { PluginAPI } from "./PLUGIN.type";


const pluginId = __PLUGIN_ID__;
window[`__plugin_${pluginId}__`] = async function (api: PluginAPI) {
// window.__plugin_my_plugin__ = 
  // 1. 读取管理员配置的值
  const config = api.getConfig();
  const title = (config.title as string) || "我的插件";

  // 2. 注册插槽组件
  api.registerSlot(
    "sidebar-top",
    () => React.createElement(MyWidget, {}),
    { order: 10 }, // 数字越小越靠前
  );

  // 3. 注册到帖子详情底部，并能接收 slotProps
  api.registerSlot(
    "post-detail-bottom",
    ({ postId }: { postId?: string }) =>
      React.createElement(MyWidget, { postId }),
    { order: 5 },
  );

  // 4. 监听事件
  api.on("post:view", (data) => {
    api.log("info", `用户浏览了帖子: ${JSON.stringify(data)}`);
  });

  api.on("user:login", (data) => {
    const user = api.getUser();
    api.log("info", `用户 ${user?.username} 已登录`);
  });
};
