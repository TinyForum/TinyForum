// api/types/nocode.model.ts
export interface NodeMeta {
  type: string;
  name: string;
  description?: string;
  inputs?: Record<string, unknown>;
  outputs?: Record<string, unknown>;
}

export interface NocodeMetadata {
  triggers: NodeMeta[];
  conditions: NodeMeta[];
  actions: NodeMeta[];
}

export interface FlowNode {
  id: string;
  type: string;
  config?: Record<string, unknown>;
}

export interface FlowEdge {
  source: string;
  target: string;
}

export interface Flow {
  nodes: FlowNode[];
  edges: FlowEdge[];
}

// 用于 validate 请求
export interface ValidateFlowRequest {
  flow: Flow;
}
