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
} from "reactflow";
import "reactflow/dist/style.css";
import {
  useNocodeMetadata,
  useValidateFlow,
  useBotActions,
} from "../hooks/bot";
import { NodeMeta, Flow } from "../noco.type";

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
    label: meta.name,
  },
});

// 简单的节点配置弹窗（可扩展）
const NodeConfigModal = ({ node, onSave, onClose }: any) => {
  const [config, setConfig] = useState(node.data.config);
  const handleSave = () => {
    onSave(node.id, config);
    onClose();
  };
  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center">
      <div className="bg-white p-4 rounded shadow-lg w-96">
        <h3>配置 {node.data.label}</h3>
        <textarea
          className="w-full border p-2 mt-2"
          rows={5}
          value={JSON.stringify(config, null, 2)}
          onChange={(e) => {
            try {
              setConfig(JSON.parse(e.target.value));
            } catch {}
          }}
        />
        <div className="flex justify-end mt-4 space-x-2">
          <button className="btn btn-sm" onClick={onClose}>
            取消
          </button>
          <button className="btn btn-sm btn-primary" onClick={handleSave}>
            保存
          </button>
        </div>
      </div>
    </div>
  );
};

export function NocoBot() {
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
  const { triggers, conditions, actions } = useMemo(() => {
    if (!metadata) return { triggers: [], conditions: [], actions: [] };
    return metadata;
  }, [metadata]);

  const onConnect = useCallback(
    (params: Connection) => setEdges((eds) => addEdge(params, eds)),
    [],
  );

  const addNodeToCanvas = (nodeMeta: NodeMeta) => {
    const newNode = toReactFlowNode(nodeMeta, {
      x: Math.random() * 400,
      y: Math.random() * 400,
    });
    setNodes((nds) => nds.concat(newNode));
  };

  const onNodeDoubleClick = (
    _event: React.MouseEvent,
    node: Node<FlowNodeData>,
  ) => {
    setSelectedNode(node);
  };

  const updateNodeConfig = (nodeId: string, config: any) => {
    setNodes((nds) =>
      nds.map((n) =>
        n.id === nodeId ? { ...n, data: { ...n.data, config } } : n,
      ),
    );
  };

  // 构建用于校验和保存的 Flow 结构
  const buildFlowPayload = (): Flow => {
    return {
      nodes: nodes.map((n) => ({
        id: n.id,
        type: n.data.nodeMeta.type,
        config: n.data.config,
      })),
      edges: edges.map((e) => ({ source: e.source, target: e.target })),
    };
  };

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
    // 校验
    const validation = await validate(flow);
    if (!validation.valid) {
      alert(`保存前校验失败：\n${validation.errors?.join("\n")}`);
      return;
    }

    // 构建创建机器人的请求体（需后端支持 nocode 类型）
    const requestData = {
      name: "未命名零代码机器人",
      description: "通过节点图创建",
      type: "nocode", // 标识零代码类型
      flow: flow, // 流程数据
    };
    const id = await createBot(requestData as any); // 强转，需根据实际 CreateBotRequest 调整
    if (id) {
      alert(`机器人创建成功！ID: ${id}`);
      // 可跳转或清空画布
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
    <div className="flex flex-col h-[600px] border rounded-lg overflow-hidden">
      <div className="flex flex-1">
        {/* 左侧节点库 */}
        <div className="w-64 bg-gray-100 p-3 overflow-y-auto border-r">
          <h3 className="font-bold mb-2">触发器</h3>
          {triggers.map((t) => (
            <div
              key={t.type}
              className="mb-2 p-2 bg-white rounded shadow cursor-move hover:bg-blue-50"
              draggable
              onDragStart={(e) => {
                e.dataTransfer.setData("application/json", JSON.stringify(t));
              }}
              onClick={() => addNodeToCanvas(t)}
            >
              {t.name}
            </div>
          ))}
          <h3 className="font-bold mt-4 mb-2">条件</h3>
          {conditions.map((c) => (
            <div
              key={c.type}
              className="mb-2 p-2 bg-white rounded shadow cursor-move"
              draggable
              onDragStart={(e) =>
                e.dataTransfer.setData("application/json", JSON.stringify(c))
              }
              onClick={() => addNodeToCanvas(c)}
            >
              {c.name}
            </div>
          ))}
          <h3 className="font-bold mt-4 mb-2">动作</h3>
          {actions.map((a) => (
            <div
              key={a.type}
              className="mb-2 p-2 bg-white rounded shadow cursor-move"
              draggable
              onDragStart={(e) =>
                e.dataTransfer.setData("application/json", JSON.stringify(a))
              }
              onClick={() => addNodeToCanvas(a)}
            >
              {a.name}
            </div>
          ))}
        </div>

        {/* 画布区域 */}
        <div className="flex-1 relative">
          <ReactFlow
            nodes={nodes}
            edges={edges}
            onNodesChange={onNodesChange}
            onEdgesChange={onEdgesChange}
            onConnect={onConnect}
            onNodeDoubleClick={onNodeDoubleClick}
            onDragOver={(e) => e.preventDefault()}
            onDrop={(e) => {
              e.preventDefault();
              const rawData = e.dataTransfer.getData("application/json");
              if (rawData) {
                const nodeMeta: NodeMeta = JSON.parse(rawData);
                const position = { x: e.clientX - 300, y: e.clientY - 200 };
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
                className="btn btn-sm btn-outline"
                onClick={handleValidate}
                disabled={validating}
              >
                {validating ? "校验中..." : "校验流程"}
              </button>
              <button
                className="btn btn-sm btn-primary"
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
