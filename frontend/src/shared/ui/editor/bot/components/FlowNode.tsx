"use client";

import { memo } from "react";
import { Handle, Position, NodeProps } from "reactflow";

export interface FlowNodeData {
  nodeType: string;
  label: string;
  category: "trigger" | "control" | "variable" | "action";
  params: Record<string, unknown>;
  outputs?: { name: string; type: string; desc?: string }[];
}

const categoryConfig: Record<string, { border: string; bg: string; badge: string; badgeColor: string }> = {
  trigger:  { border: "border-green-400",  bg: "bg-green-50",  badge: "触发器", badgeColor: "bg-green-200 text-green-800" },
  control:  { border: "border-orange-400", bg: "bg-orange-50", badge: "控制",   badgeColor: "bg-orange-200 text-orange-800" },
  variable: { border: "border-purple-400", bg: "bg-purple-50", badge: "变量",   badgeColor: "bg-purple-200 text-purple-800" },
  action:   { border: "border-blue-400",   bg: "bg-blue-50",   badge: "动作",   badgeColor: "bg-blue-200 text-blue-800" },
};

export const FlowNode = memo(function FlowNode({
  id,
  data,
  selected,
}: NodeProps<FlowNodeData>) {
  const config = categoryConfig[data.category] || categoryConfig.action;
  const isIf = data.nodeType === "if";
  const isWhile = data.nodeType === "while";

  return (
    <div
      className={`border-2 rounded-lg px-3 py-2 min-w-[160px] shadow-md cursor-pointer
        ${config.border} ${config.bg}
        ${selected ? "ring-2 ring-blue-500 ring-offset-1" : ""}`}
      data-node-id={id}
    >
      {data.category !== "trigger" && (
        <Handle type="target" position={Position.Top} className="!w-3 !h-3 !bg-gray-400 !border-2 !border-white" />
      )}

      <span className={`inline-block text-xs px-1.5 py-0.5 rounded ${config.badgeColor}`}>
        {config.badge}
      </span>

      <div className="text-sm font-medium mt-1 text-gray-800">{data.label}</div>

      {/* if/else: true / false */}
      {isIf && (
        <>
          <Handle type="source" position={Position.Bottom} id="true" className="!w-3 !h-3 !bg-green-500 !border-2 !border-white" style={{ left: "30%" }} />
          <span className="absolute text-[10px] text-green-700" style={{ bottom: 0, left: "28%" }}>T</span>
          <Handle type="source" position={Position.Bottom} id="false" className="!w-3 !h-3 !bg-red-500 !border-2 !border-white" style={{ left: "68%" }} />
          <span className="absolute text-[10px] text-red-700" style={{ bottom: 0, left: "67%" }}>F</span>
        </>
      )}

      {/* while: body */}
      {isWhile && (
        <>
          <Handle type="source" position={Position.Bottom} id="body" className="!w-3 !h-3 !bg-purple-500 !border-2 !border-white" style={{ left: "30%" }} />
          <span className="absolute text-[10px] text-purple-700" style={{ bottom: 0, left: "24%" }}>body</span>
          <Handle type="source" position={Position.Bottom} id="exit" className="!w-3 !h-3 !bg-gray-500 !border-2 !border-white" style={{ left: "68%" }} />
          <span className="absolute text-[10px] text-gray-700" style={{ bottom: 0, left: "67%" }}>exit</span>
        </>
      )}

      {/* 输出变量提示 */}
      {data.outputs && data.outputs.length > 0 && (
        <div className="mt-1 flex flex-wrap gap-1">
          {data.outputs.slice(0, 3).map((o) => (
            <span key={o.name} className="text-[10px] px-1 bg-gray-200 rounded" title={o.desc}>
              {o.name}:{o.type}
            </span>
          ))}
          {data.outputs.length > 3 && (
            <span className="text-[10px] text-gray-400">+{data.outputs.length - 3}</span>
          )}
        </div>
      )}

      {/* 普通输出端口 */}
      {data.category !== "action" && !isIf && !isWhile && (
        <Handle type="source" position={Position.Bottom} className="!w-3 !h-3 !bg-gray-400 !border-2 !border-white" />
      )}
    </div>
  );
});
