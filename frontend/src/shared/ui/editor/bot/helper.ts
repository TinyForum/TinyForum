// ---------- 辅助函数 ----------
// function _getNodeLabel(nodeMeta: NodeMeta): string {
//   console.log("获取节点标签: ", nodeMeta);
//   return nodeMeta.label;
// }

import { NodeMeta } from "@/features/bot/noco.type";

/**
 * 创建默认节点参数的函数
 * @param nodeMeta - 节点的元数据信息，包含schema等信息
 * @returns 返回一个包含默认参数的Record对象
 */
export function createDefaultParams(
  nodeMeta: NodeMeta,
): Record<string, unknown> {
  console.log("创建默认节点参数: ", nodeMeta); // 输出创建默认参数时的节点元数据信息
  // 可从 nodeMeta.schema 生成默认值，此处简单返回空对象
  return {}; // 返回一个空对象作为默认参数
}
