"use client";

import { useState, useCallback, useMemo } from "react";
import ReactFlow, {
  Controls,
  Background,
  MiniMap,
  ReactFlowProvider,
  ReactFlowInstance,
  useNodesState,
  useEdgesState,
  addEdge,
  Connection,
  MarkerType,
  Node,
} from "reactflow";
import "reactflow/dist/style.css";

import {
  useNocodeMetadata,
  useValidateFlow,
} from "@/features/bot/hooks/useNocodeMetadata";
import { useBotActions } from "@/features/bot/hooks/bot";
import { CreateBotRequest, TriggerNode } from "@/shared/api/types/bot.model";
import { BotTriggerType } from "@/shared/api/types/bot.model.do";
import { NodeMeta } from "@/shared/api/types/nocode.model";
import toast from "react-hot-toast";

import { useTranslations } from "next-intl";
import { NodePalette } from "./components/NodePalette";
import { FlowNode, FlowNodeData } from "./components/FlowNode";
import { ParamFormModal } from "./components/ParamFormModal";
import { generateNodeId, graphToFlow, validateDepth } from "./flowUtils";
import { createDefaultParams } from "./helper";

const nodeTypes = { flowNode: FlowNode };

// ── 主编辑器组件 ──
function FlowEditor() {
  const t = useTranslations("Bot");
  const {
    metadata,
    loading: metaLoading,
    error: metaError,
  } = useNocodeMetadata();
  const { validate, loading: validating } = useValidateFlow();
  const { createBot, loading: saving } = useBotActions();

  const [nodes, setNodes, onNodesChange] = useNodesState<FlowNodeData>([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);
  const [reactFlowInstance, setReactFlowInstance] =
    useState<ReactFlowInstance | null>(null);

  const [editingNode, setEditingNode] = useState<{
    id: string;
    meta: NodeMeta;
  } | null>(null);

  // 合并所有节点元数据（便于按 type 查找）
  const allMetas = useMemo<NodeMeta[]>(() => {
    if (!metadata) return [];
    return [
      ...metadata.triggers,
      ...metadata.control,
      ...metadata.variables,
      ...metadata.actions,
    ];
  }, [metadata]);

  // ── 连线 ──
  const onConnect = useCallback(
    (connection: Connection) => {
      setEdges((eds) =>
        addEdge(
          {
            ...connection,
            markerEnd: { type: MarkerType.ArrowClosed },
            style: { stroke: "#94a3b8", strokeWidth: 2 },
          },
          eds,
        ),
      );
    },
    [setEdges],
  );

  // ── 允许拖放 ──
  const onDragOver = useCallback((event: React.DragEvent) => {
    event.preventDefault();
    event.dataTransfer.dropEffect = "move";
  }, []);

  // ── 从画板拖入节点 ──
  const onDrop = useCallback(
    (event: React.DragEvent) => {
      event.preventDefault();

      const type = event.dataTransfer.getData("application/reactflow-type");
      const categoryStr = event.dataTransfer.getData(
        "application/reactflow-category",
      );
      const label = event.dataTransfer.getData("application/reactflow-label");

      if (!type || !categoryStr || !reactFlowInstance) return;

      const category = categoryStr as FlowNodeData["category"];
      const position = reactFlowInstance.screenToFlowPosition({
        x: event.clientX,
        y: event.clientY,
      });

      const meta = allMetas.find((m) => m.type === type);
      const newNode: Node<FlowNodeData> = {
        id: generateNodeId(),
        type: "flowNode",
        position,
        data: {
          nodeType: type,
          label,
          category,
          params: meta ? createDefaultParams(meta) : {},
        },
      };

      setNodes((nds) => [...nds, newNode]);
    },
    [reactFlowInstance, allMetas, setNodes],
  );

  // ── 双击打开配置 ──
  const onNodeDoubleClick = useCallback(
    (_event: React.MouseEvent, node: Node<FlowNodeData>) => {
      const meta = allMetas.find((m) => m.type === node.data.nodeType);
      if (meta) {
        setEditingNode({ id: node.id, meta });
      }
    },
    [allMetas],
  );

  // ── 更新节点参数 ──
  const updateNodeParams = useCallback(
    (nodeId: string, newParams: Record<string, unknown>) => {
      setNodes((nds) =>
        nds.map((n) =>
          n.id === nodeId
            ? { ...n, data: { ...n.data, params: newParams } }
            : n,
        ),
      );
    },
    [setNodes],
  );

  // ── 根据 nocode 触发器类型推导 bot 顶层字段 ──
  const deriveTriggerFields = (trigger: TriggerNode): {
    triggerType: BotTriggerType;
    eventFilter?: string;
    cronExpr?: string;
  } => {
    switch (trigger.type) {
      case "on_schedule":
        return {
          triggerType: "schedule",
          cronExpr: (trigger.params?.cron as string) || "",
        };
      case "on_new_post":
        return { triggerType: "event", eventFilter: "post.created" };
      case "on_new_comment":
        return { triggerType: "event", eventFilter: "comment.created" };
      case "on_user_register":
        return { triggerType: "event", eventFilter: "user.registered" };
      case "on_keyword":
        return { triggerType: "event", eventFilter: "post.created,comment.created" };
      default:
        return { triggerType: "manual" };
    }
  };

  // ── 构建 Flow 并保存 ──
  const handleSave = async () => {
    // 深度检查
    const depthCheck = validateDepth(nodes, edges);
    if (!depthCheck.valid) {
      toast.error(`${t("flow_validate_errors")}\n${depthCheck.errors.join("\n")}`);
      return;
    }

    const flow = graphToFlow(nodes, edges);
    if (!flow) {
      toast.error(t("add_trigger_first"));
      return;
    }

    const nonTriggerNodes = nodes.filter(
      (n) => n.data.category !== "trigger",
    );
    if (nonTriggerNodes.length === 0) {
      toast.error(t("add_at_least_one_node"));
      return;
    }

    const validation = await validate(flow);
    if (!validation.valid) {
      toast.error(
        `${t("validate_errors")}\n${validation.errors?.join("\n") || t("unknown_error_suffix")}`,
      );
      return;
    }

    const { triggerType, eventFilter, cronExpr } = deriveTriggerFields(
      flow.trigger,
    );

    const requestData: CreateBotRequest = {
      name: t("unnamed_nocode_bot"),
      version: "1.0.0",
      description: t("nocode_bot_desc"),
      type: "task",
      scriptCode: "nocode",
      triggerType,
      eventFilter,
      cronExpr,
      timeoutSec: 10,
      configValues: { flow },
    };

    const id = await createBot(requestData);
    if (id) {
      toast.success(t("create_bot_success_with_id", { id }));
      setNodes([]);
      setEdges([]);
    }
  };

  // ── 加载态 / 错误态 ──
  if (metaLoading) {
    return (
      <div className="flex items-center justify-center h-[calc(100vh-24rem)]">
        <span className="loading loading-spinner loading-md" />
        <span className="ml-3 text-gray-500">{t("loading_node_definitions")}</span>
      </div>
    );
  }

  if (metaError || !metadata) {
    return (
      <div className="flex items-center justify-center h-[calc(100vh-24rem)] text-red-500">
        {t("load_failed")}{metaError || t("unknown_error_label")}
      </div>
    );
  }

  return (
    <div className="flex flex-col h-[calc(100vh-24rem)] border rounded-lg overflow-hidden">
      <div className="flex flex-1 min-h-0">
        {/* 左侧节点库 */}
        <div className="w-64 flex-shrink-0 border-r bg-white">
          <NodePalette metadata={metadata} />
        </div>

        {/* 中间画布 */}
        <div className="flex-1">
          <ReactFlow
            nodes={nodes}
            edges={edges}
            onNodesChange={onNodesChange}
            onEdgesChange={onEdgesChange}
            onConnect={onConnect}
            onDragOver={onDragOver}
            onDrop={onDrop}
            onNodeDoubleClick={onNodeDoubleClick}
            onInit={setReactFlowInstance}
            nodeTypes={nodeTypes}
            fitView
            deleteKeyCode={["Backspace", "Delete"]}
            multiSelectionKeyCode="Shift"
            snapToGrid
            snapGrid={[20, 20]}
            defaultEdgeOptions={{
              type: "smoothstep",
              animated: true,
              style: { stroke: "#94a3b8", strokeWidth: 2 },
              markerEnd: { type: MarkerType.ArrowClosed, color: "#94a3b8" },
            }}
          >
            <Controls />
            <Background gap={20} size={1} color="#e2e8f0" />
            <MiniMap
              nodeStrokeColor="#94a3b8"
              nodeColor={(n) => {
                const d = n.data as FlowNodeData | undefined;
                if (d?.category === "trigger") return "#bbf7d0";
                if (d?.category === "control") return "#fed7aa";
                if (d?.category === "variable") return "#e9d5ff";
                if (d?.category === "action") return "#bfdbfe";
                return "#e2e8f0";
              }}
              maskColor="rgba(0,0,0,0.05)"
            />
          </ReactFlow>
        </div>
      </div>

      {/* 底部操作栏 */}
      <div className="flex-shrink-0 border-t bg-white px-4 py-3 flex items-center justify-between">
        <div className="text-sm text-gray-500">
          {t("drag_node_hint")}
        </div>
        <div className="flex gap-2">
          <button
            className="px-4 py-2 border rounded-lg text-sm hover:bg-gray-50
              disabled:opacity-50"
            disabled={validating}
          >
            {validating ? t("validating") : t("validate_flow")}
          </button>
          <button
            className="px-4 py-2 bg-blue-600 text-white rounded-lg text-sm
              hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed"
            onClick={handleSave}
            disabled={saving || validating}
          >
            {saving ? t("saving") : t("save_bot")}
          </button>
        </div>
      </div>

      {/* 参数配置弹窗 */}
      {editingNode && (
        <ParamFormModal
          title={editingNode.meta.label}
          params={
            nodes.find((n) => n.id === editingNode.id)?.data.params || {}
          }
          paramMetas={editingNode.meta.params || []}
          onSave={(newParams) => {
            updateNodeParams(editingNode.id, newParams);
            setEditingNode(null);
          }}
          onClose={() => setEditingNode(null)}
        />
      )}
    </div>
  );
}

// ── 导出（包裹 ReactFlowProvider） ──
export function BotFlowEditor() {
  return (
    <ReactFlowProvider>
      <FlowEditor />
    </ReactFlowProvider>
  );
}
