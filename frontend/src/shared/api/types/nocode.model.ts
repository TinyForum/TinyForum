// api/types/nocode.model.ts

export interface VarOutput {
  name: string;
  type: string; // string | number | boolean | object
  desc?: string;
}

export interface NodeMeta {
  type: string;
  label: string;
  description?: string;
  icon?: string;
  category?: string;
  params?: ParamMeta[];
  outputs?: VarOutput[];
}

export interface NocodeMetadata {
  triggers: NodeMeta[];
  control: NodeMeta[];
  variables: NodeMeta[];
  actions: NodeMeta[];
}

export interface ParamMeta {
  key: string;
  label: string;
  type: string;
  required: boolean;
  default?: unknown;
  placeholder?: string;
  options?: OptionMeta[];
}

export interface OptionMeta {
  label: string;
  value: unknown;
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
