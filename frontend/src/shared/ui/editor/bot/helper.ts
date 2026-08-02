import { NodeMeta, ParamMeta } from "@/shared/api/types/nocode.model";

export function createDefaultParams(
  nodeMeta: NodeMeta,
): Record<string, unknown> {
  const params: Record<string, unknown> = {};
  if (!nodeMeta.params) return params;

  nodeMeta.params.forEach((meta: ParamMeta) => {
    if (meta.default !== undefined && meta.default !== null) {
      params[meta.key] = meta.default;
    } else {
      switch (meta.type) {
        case "number":
          params[meta.key] = 0;
          break;
        case "boolean":
          params[meta.key] = false;
          break;
        case "tags":
          params[meta.key] = [];
          break;
        default:
          params[meta.key] = "";
      }
    }
  });

  return params;
}
