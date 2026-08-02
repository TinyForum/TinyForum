"use client";

import { NocodeMetadata, NodeMeta } from "@/shared/api/types/nocode.model";
import { CollapsibleSection } from "./CollapsibleSection";
import { FlowNodeData } from "./FlowNode";

interface NodePaletteProps {
  metadata: NocodeMetadata;
}

function onDragStart(
  event: React.DragEvent,
  nodeMeta: NodeMeta,
  category: FlowNodeData["category"],
) {
  event.dataTransfer.setData("application/reactflow-type", nodeMeta.type);
  event.dataTransfer.setData("application/reactflow-category", category);
  event.dataTransfer.setData("application/reactflow-label", nodeMeta.label);
  event.dataTransfer.setData(
    "application/reactflow-outputs",
    JSON.stringify(nodeMeta.outputs || []),
  );
  event.dataTransfer.effectAllowed = "move";
}

function PaletteItem({
  nodeMeta,
  category,
}: {
  nodeMeta: NodeMeta;
  category: FlowNodeData["category"];
}) {
  return (
    <div
      draggable
      onDragStart={(e) => onDragStart(e, nodeMeta, category)}
      className="p-2 bg-white rounded shadow cursor-grab active:cursor-grabbing
        hover:bg-blue-50 hover:shadow-md transition-all
        border border-gray-200 text-sm"
      title={nodeMeta.description}
    >
      <div className="font-medium text-gray-800 truncate">{nodeMeta.label}</div>
      {nodeMeta.description && (
        <div className="text-xs text-gray-500 mt-0.5 line-clamp-2">
          {nodeMeta.description}
        </div>
      )}
      {nodeMeta.outputs && nodeMeta.outputs.length > 0 && (
        <div className="mt-1 flex flex-wrap gap-0.5">
          {nodeMeta.outputs.map((o) => (
            <span key={o.name} className="text-[10px] px-1 bg-gray-100 rounded">
              {o.name}
            </span>
          ))}
          {/* {nodeMeta.outputs.slice(0, 2).map((o) => (
            <span key={o.name} className="text-[10px] px-1 bg-gray-100 rounded">{o.name}</span>
          ))}
          {nodeMeta.outputs.length > 2 && (
            <span className="text-[10px] text-gray-400">+{nodeMeta.outputs.length - 2}</span>
          )} */}
        </div>
      )}
    </div>
  );
}

const sectionColors: Record<string, string> = {
  trigger: "border-l-green-400",
  control: "border-l-orange-400",
  variable: "border-l-purple-400",
  action: "border-l-blue-400",
};

export function NodePalette({ metadata }: NodePaletteProps) {
  return (
    <div className="h-full overflow-y-auto p-3 space-y-3 bg-gray-50">
      <div className="text-sm font-bold text-gray-600 mb-2">拖拽节点到画布</div>

      <Section
        color="trigger"
        title="触发器"
        items={metadata.triggers}
        category="trigger"
      />
      <Section
        color="control"
        title="控制"
        items={metadata.control}
        category="control"
      />
      <Section
        color="variable"
        title="变量"
        items={metadata.variables}
        category="variable"
      />
      <Section
        color="action"
        title="动作"
        items={metadata.actions}
        category="action"
      />
    </div>
  );
}

function Section({
  color,
  title,
  items,
  category,
}: {
  color: string;
  title: string;
  items: NodeMeta[];
  category: FlowNodeData["category"];
}) {
  if (!items || items.length === 0) return null;
  return (
    <div className={`border-l-2 pl-2 ${sectionColors[color]}`}>
      <CollapsibleSection title={title} defaultOpen>
        {items.map((m) => (
          <PaletteItem key={m.type} nodeMeta={m} category={category} />
        ))}
      </CollapsibleSection>
    </div>
  );
}
