// api/types/nocode.model.ts
export interface NodeMeta {
  type: string;
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

//-----

export interface NodeMeta {
  /** 显示名称 */
  /** 节点类型标识，如 "http_trigger", "condition_equal", "send_email" */
  type: string;
  label: string;
  /** 描述信息 */
  description?: string;
  icon?: string;
  category?: string;
  params?: ParamMeta[];
}

export interface NocodeMetadata {
  triggers: NodeMeta[];
  conditions: NodeMeta[];
  actions: NodeMeta[];
}
export interface ParamMeta {
  key: string;
  label: string;
  type: string;
  required: boolean;
  // default?: any;
  placeholder?: string;
  // options?: OptionMeta[];
}

// export interface OptionMeta {
//   label: string;
//   value: any;
// }
