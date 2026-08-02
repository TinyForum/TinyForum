import { Node, Edge } from "reactflow";
import {
  Flow,
  TriggerNode,
  CondType,
  FlowStep,
} from "@/shared/api/types/bot.model";
import { FlowNodeData } from "./components/FlowNode";

let _counter = 0;
export function generateNodeId(): string {
  return `node_${Date.now()}_${++_counter}`;
}

const MAX_DEPTH = 5;

function getOutgoingEdges(srcId: string, edges: Edge[]): Record<string, Edge[]> {
  const g: Record<string, Edge[]> = {};
  edges.filter((e) => e.source === srcId).forEach((e) => {
    const h = e.sourceHandle || "__default";
    if (!g[h]) g[h] = [];
    g[h].push(e);
  });
  return g;
}

function getTargetIds(eds: Edge[]): string[] {
  return eds.map((e) => e.target);
}

function nodeToStep(
  id: string,
  nodeMap: Map<string, Node<FlowNodeData>>,
  edges: Edge[],
  visited: Set<string>,
  depth: number,
): FlowStep | null {
  if (visited.has(id)) return null;
  visited.add(id);
  const node = nodeMap.get(id);
  if (!node || node.data.category === "trigger") return null;

  const outEdges = getOutgoingEdges(id, edges);
  const step: FlowStep = {
    id: generateNodeId(),
    type: node.data.nodeType,
    label: node.data.label,
    params: node.data.params || {},
  };

  // if 节点
  if (node.data.nodeType === "if") {
    if (depth >= MAX_DEPTH) return step;
    const trueIds = getTargetIds(outEdges["true"] || []);
    const falseIds = getTargetIds(outEdges["false"] || []);
    step.branch = {
      condition: buildCondFromParams(node.data.params),
      true: trueIds
        .map((tid) => nodeToStep(tid, nodeMap, edges, visited, depth + 1))
        .filter(Boolean) as FlowStep[],
      false: falseIds.length > 0
        ? (falseIds
            .map((fid) => nodeToStep(fid, nodeMap, edges, visited, depth + 1))
            .filter(Boolean) as FlowStep[])
        : undefined,
    };
    // 标记子节点已访问
    [...trueIds, ...falseIds].forEach((t) => visited.add(t));
    return step;
  }

  // while 节点
  if (node.data.nodeType === "while") {
    if (depth >= MAX_DEPTH) return step;
    const bodyIds = getTargetIds(outEdges["body"] || []);
    const exitIds = getTargetIds(outEdges["exit"] || []);
    step.loop = {
      condition: buildCondFromParams(node.data.params),
      body: bodyIds
        .map((bid) => nodeToStep(bid, nodeMap, edges, visited, depth + 1))
        .filter(Boolean) as FlowStep[],
      max_iter: Number(node.data.params?.max_iter) || 100,
    };
    bodyIds.forEach((b) => visited.add(b));
    if (exitIds.length > 0) return step;
    return step;
  }

  // 普通节点：不在此标记子节点 visited，由 call site 按序遍历
  return step;
}

function buildCondFromParams(params: Record<string, unknown>): {
  type: CondType;
  params: Record<string, unknown>;
} {
  const field = String(params?.condition_field || "post_id");
  const op = String(params?.condition_op || "equals");
  const value = String(params?.condition_value || "");

  const opMap: Record<string, CondType> = {
    equals: "field_equals",
    not_equals: "field_not_equals",
    contains: "field_contains",
    greater_than: "field_greater_than",
    less_than: "field_less_than",
    is_empty: "field_is_empty",
    is_not_empty: "field_is_not_empty",
  };

  return {
    type: opMap[op] || "field_equals",
    params: { field, value },
  };
}

export function graphToFlow(
  nodes: Node<FlowNodeData>[],
  edges: Edge[],
): Flow | null {
  const triggerNodes = nodes.filter((n) => n.data.category === "trigger");
  if (triggerNodes.length === 0) return null;

  const nodeMap = new Map(nodes.map((n) => [n.id, n]));
  const visited = new Set<string>();

  const trigger: TriggerNode = {
    type: triggerNodes[0].data.nodeType as TriggerNode["type"],
    params: triggerNodes[0].data.params || {},
  };

  // 从 trigger 出发收集所有步骤
  const steps: FlowStep[] = [];
  const triggerOut = getOutgoingEdges(triggerNodes[0].id, edges);
  const firstIds = getTargetIds(triggerOut["__default"] || []);

  if (firstIds.length > 0) {
    let currentId: string | undefined = firstIds[0];
    while (currentId && !visited.has(currentId)) {
      const step = nodeToStep(currentId, nodeMap, edges, visited, 0);
      if (step) steps.push(step);

      // 找下一个未访问的节点
      const node = nodeMap.get(currentId);
      if (!node) break;
      const outE = getOutgoingEdges(currentId, edges);
      const nextIds = getTargetIds(outE["__default"] || []);
      // 跳过 if/while 节点的子节点（已由 branch/loop 处理）
      if (node.data.nodeType === "if" || node.data.nodeType === "while") {
        const exitIds = getTargetIds(outE["exit"] || []);
        currentId = exitIds.length > 0 ? exitIds[0] : undefined;
      } else {
        currentId = nextIds.length > 0 && !visited.has(nextIds[0])
          ? nextIds[0]
          : undefined;
      }
    }
  }

  // 安全网：确保所有非 trigger 节点都被收集（兜底）
  if (steps.length === 0) {
    for (const node of nodes) {
      if (node.data.category === "trigger") continue;
      steps.push({
        id: generateNodeId(),
        type: node.data.nodeType,
        params: node.data.params || {},
      });
    }
  }

  return {
    version: "1",
    trigger,
    steps,
  };
}

export function validateDepth(
  nodes: Node<FlowNodeData>[],
  edges: Edge[],
): { valid: boolean; errors: string[] } {
  const errors: string[] = [];
  const nodeMap = new Map(nodes.map((n) => [n.id, n]));

  function check(id: string, depth: number, path: string[]) {
    if (depth > MAX_DEPTH) {
      errors.push(
        `嵌套深度超限 (${depth} > ${MAX_DEPTH}): ${path.join(" → ")}`,
      );
      return;
    }
    path.push(id);
    const outE = getOutgoingEdges(id, edges);
    const node = nodeMap.get(id);
    if (node?.data.nodeType === "if") {
      getTargetIds(outE["true"] || []).forEach((t) => check(t, depth + 1, [...path]));
      getTargetIds(outE["false"] || []).forEach((t) => check(t, depth + 1, [...path]));
    } else if (node?.data.nodeType === "while") {
      getTargetIds(outE["body"] || []).forEach((t) => check(t, depth + 1, [...path]));
    }
    getTargetIds(outE["__default"] || []).forEach((t) => check(t, depth, [...path]));
  }

  const triggers = nodes.filter((n) => n.data.category === "trigger");
  if (triggers.length === 0) {
    errors.push("缺少触发器节点");
    return { valid: false, errors };
  }

  const outE = getOutgoingEdges(triggers[0].id, edges);
  getTargetIds(outE["__default"] || []).forEach((t) => check(t, 0, [triggers[0].id]));

  return { valid: errors.length === 0, errors };
}
