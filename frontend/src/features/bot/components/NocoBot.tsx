// features/bot/components/NocoBot.tsx
import { useState, useCallback, useMemo, useRef } from "react";
import ReactFlow, {
  addEdge,
  Background,
  Connection,
  Edge,
  Node,
  useNodesState,
  useEdgesState,
  Controls,
  MiniMap,
  Panel,
  useReactFlow,
  ReactFlowProvider,
} from "reactflow";
import "reactflow/dist/style.css";

import { NodeMeta, Flow, NocodeMetadata } from "../noco.type";
import { useBotActions } from "../hooks/bot";
import { useNocodeMetadata, useValidateFlow } from "../hooks/useNocodeMetadata";

// 定义画布节点结构（扩展 React Flow 节点）
interface FlowNodeData {
  nodeMeta: NodeMeta; // 元数据
  config: Record<string, any>; // 用户配置
  label: string;
}

// 将后端 NodeMeta 转换为 React Flow 节点
const toReactFlowNode = (
  meta: NodeMeta,
  position: { x: number; y: number },
): Node<FlowNodeData> => ({
  id: `${meta.type}_${Date.now()}_${Math.random()}`,
  type: "default",
  position,
  data: {
    nodeMeta: meta,
    config: {},
    label: meta.label,
  },
});

// 节点配置弹窗组件
const NodeConfigModal = ({ node, onSave, onClose }: any) => {
  const [config, setConfig] = useState(node.data.config);
  const handleSave = () => {
    onSave(node.id, config);
    onClose();
  };
  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white p-4 rounded shadow-lg w-96">
        <h3 className="font-bold text-lg mb-2">配置 {node.data.label}</h3>
        <textarea
          className="w-full border p-2 mt-2 font-mono text-sm"
          rows={6}
          value={JSON.stringify(config, null, 2)}
          onChange={(e) => {
            try {
              setConfig(JSON.parse(e.target.value));
            } catch {
              // 如果 JSON 解析失败，保留原文本
            }
          }}
        />
        <div className="flex justify-end mt-4 space-x-2">
          <button
            className="px-3 py-1 border rounded hover:bg-gray-100"
            onClick={onClose}
          >
            取消
          </button>
          <button
            className="px-3 py-1 bg-blue-600 text-white rounded hover:bg-blue-700"
            onClick={handleSave}
          >
            保存
          </button>
        </div>
      </div>
    </div>
  );
};

interface CollapsibleSectionProps {
  title: string;
  children: React.ReactNode;
  defaultOpen?: boolean;
}

function CollapsibleSection({
  title,
  children,
  defaultOpen = true,
}: CollapsibleSectionProps) {
  const [isOpen, setIsOpen] = useState(defaultOpen);
  return (
    <div className="mb-4">
      <div
        className="flex items-center justify-between cursor-pointer py-1 select-none"
        onClick={() => setIsOpen(!isOpen)}
      >
        <h3 className="font-bold text-gray-800">{title}</h3>
        <span className="text-gray-500 text-sm">{isOpen ? "▼" : "▶"}</span>
      </div>
      {isOpen && <div className="mt-2 space-y-2">{children}</div>}
    </div>
  );
}

// 内部组件，必须被 ReactFlowProvider 包裹，才能使用 useReactFlow
function NocoBotInner() {
  const { screenToFlowPosition, getViewport } = useReactFlow();
  const containerRef = useRef<HTMLDivElement>(null);

  const {
    metadata,
    loading: metaLoading,
    error: metaError,
  } = useNocodeMetadata();
  const { validate, loading: validating } = useValidateFlow();
  const { createBot, loading: saving } = useBotActions();
  const [nodes, setNodes, onNodesChange] = useNodesState<FlowNodeData>([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);
  const [selectedNode, setSelectedNode] = useState<Node<FlowNodeData> | null>(
    null,
  );
  const [savingError, setSavingError] = useState<string | null>(null);

  // 获取分类节点用于拖拽面板
  const { triggers, conditions, actions } = useMemo<NocodeMetadata>(() => {
    if (!metadata) return { triggers: [], conditions: [], actions: [] };
    return metadata;
  }, [metadata]);

  const onConnect = useCallback(
    (params: Connection) => setEdges((eds) => addEdge(params, eds)),
    [setEdges],
  );

  // 点击左侧节点列表 —— 添加到画布中心（修复问题2）
  const addNodeToCanvas = useCallback(
    (nodeMeta: NodeMeta) => {
      if (containerRef.current) {
        const rect = containerRef.current.getBoundingClientRect();
        const centerScreen = {
          x: rect.left + rect.width / 2,
          y: rect.top + rect.height / 2,
        };
        const centerFlow = screenToFlowPosition(centerScreen);
        const newNode = toReactFlowNode(nodeMeta, centerFlow);
        setNodes((nds) => nds.concat(newNode));
      } else {
        // 后备：使用 viewport 中心
        const { x, y, zoom } = getViewport();
        const centerFlow = { x: x + 250 / zoom, y: y + 200 / zoom };
        const newNode = toReactFlowNode(nodeMeta, centerFlow);
        setNodes((nds) => nds.concat(newNode));
      }
    },
    [screenToFlowPosition, getViewport, setNodes],
  );

  const onNodeDoubleClick = (
    _event: React.MouseEvent,
    node: Node<FlowNodeData>,
  ) => {
    setSelectedNode(node);
  };

  const updateNodeConfig = useCallback(
    (nodeId: string, config: any) => {
      setNodes((nds) =>
        nds.map((n) =>
          n.id === nodeId ? { ...n, data: { ...n.data, config } } : n,
        ),
      );
    },
    [setNodes],
  );

  // 构建用于校验和保存的 Flow 结构
  const buildFlowPayload = useCallback((): Flow => {
    return {
      nodes: nodes.map((n) => ({
        id: n.id,
        type: n.data.nodeMeta.type,
        config: n.data.config,
      })),
      edges: edges.map((e) => ({ source: e.source, target: e.target })),
    };
  }, [nodes, edges]);

  const handleValidate = async () => {
    const flow = buildFlowPayload();
    const result = await validate(flow);
    if (!result.valid) {
      alert(`校验失败：\n${result.errors?.join("\n") || "未知错误"}`);
    } else {
      alert("流程校验通过！");
    }
    return result.valid;
  };

  const handleSave = async () => {
    setSavingError(null);
    const flow = buildFlowPayload();
    const validation = await validate(flow);
    if (!validation.valid) {
      alert(`保存前校验失败：\n${validation.errors?.join("\n")}`);
      return;
    }

    const requestData = {
      name: "未命名零代码机器人",
      description: "通过节点图创建",
      type: "nocode",
      flow: flow,
    };
    const id = await createBot(requestData as any);
    if (id) {
      alert(`机器人创建成功！ID: ${id}`);
      setNodes([]);
      setEdges([]);
    } else {
      setSavingError("保存失败，请重试");
    }
  };

  if (metaLoading) {
    return <div className="p-4 text-center">加载节点定义中...</div>;
  }
  if (metaError) {
    return (
      <div className="p-4 text-center text-red-500">加载失败：{metaError}</div>
    );
  }

  return (
    <div className="flex flex-col h-[calc(100vh-24rem)] border rounded-lg overflow-hidden">
      <div className="flex flex-1 min-h-0" ref={containerRef}>
        {/* 左侧节点库 - 独立滚动 */}
        <div className="w-64 bg-gray-100 p-3 overflow-y-auto border-r h-full">
          <CollapsibleSection title="触发器" defaultOpen={true}>
            {triggers.map((t) => (
              <div
                key={t.type}
                className="p-2 bg-white rounded shadow cursor-move hover:bg-blue-50"
                draggable
                onDragStart={(e) =>
                  e.dataTransfer.setData("application/json", JSON.stringify(t))
                }
                onClick={() => addNodeToCanvas(t)}
              >
                {t.label}
              </div>
            ))}
          </CollapsibleSection>

          <CollapsibleSection title="条件" defaultOpen={true}>
            {conditions.map((c) => (
              <div
                key={c.type}
                className="p-2 bg-white rounded shadow cursor-move hover:bg-blue-50"
                draggable
                onDragStart={(e) =>
                  e.dataTransfer.setData("application/json", JSON.stringify(c))
                }
                onClick={() => addNodeToCanvas(c)}
              >
                {c.label}
              </div>
            ))}
          </CollapsibleSection>

          <CollapsibleSection title="动作" defaultOpen={true}>
            {actions.map((a) => (
              <div
                key={a.type}
                className="p-2 bg-white rounded shadow cursor-move hover:bg-blue-50"
                draggable
                onDragStart={(e) =>
                  e.dataTransfer.setData("application/json", JSON.stringify(a))
                }
                onClick={() => addNodeToCanvas(a)}
              >
                {a.label}
              </div>
            ))}
          </CollapsibleSection>
        </div>

        {/* 右侧画布区域 - 独立内部滚动/平移 */}
        <div className="flex-1 relative h-full">
          <ReactFlow
            nodes={nodes}
            edges={edges}
            onNodesChange={onNodesChange}
            onEdgesChange={onEdgesChange}
            onConnect={onConnect}
            onNodeDoubleClick={onNodeDoubleClick}
            onDragOver={(e) => e.preventDefault()}
            className="h-full bg-base-100"
            onDrop={(e) => {
              e.preventDefault();
              const rawData = e.dataTransfer.getData("application/json");
              if (rawData) {
                const nodeMeta: NodeMeta = JSON.parse(rawData);
                const position = screenToFlowPosition({
                  x: e.clientX,
                  y: e.clientY,
                });
                const newNode = toReactFlowNode(nodeMeta, position);
                setNodes((nds) => nds.concat(newNode));
              }
            }}
            fitView
          >
            <Background />
            <Controls />
            <MiniMap />
            <Panel position="top-right" className="space-x-2">
              <button
                className="px-3 py-1 border rounded text-sm hover:bg-gray-50"
                onClick={handleValidate}
                disabled={validating}
              >
                {validating ? "校验中..." : "校验流程"}
              </button>
              <button
                className="px-3 py-1 bg-blue-600 text-white rounded text-sm hover:bg-blue-700 disabled:bg-gray-400"
                onClick={handleSave}
                disabled={saving}
              >
                {saving ? "保存中..." : "保存机器人"}
              </button>
            </Panel>
          </ReactFlow>
        </div>
      </div>

      {savingError && (
        <div className="text-red-500 text-center p-2">{savingError}</div>
      )}
      {selectedNode && (
        <NodeConfigModal
          node={selectedNode}
          onSave={updateNodeConfig}
          onClose={() => setSelectedNode(null)}
        />
      )}
    </div>
  );
}

// 对外暴露的组件，用 ReactFlowProvider 包裹，提供 useReactFlow 所需的上下文
export function NocoBot() {
  return (
    <ReactFlowProvider>
      <NocoBotInner />
    </ReactFlowProvider>
  );
}
