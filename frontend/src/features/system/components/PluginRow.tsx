// ── Sub-components ─────────────────────────────────────────────────────────────

import { PluginMeta } from "@/shared/type/plugin.type";
import { Switch } from "@headlessui/react";
import {
  Code2,
  ChevronUp,
  ChevronDown,
  CheckCircle2,
  PowerOff,
  Pencil,
  Trash2,
} from "lucide-react";

export function PluginRow({
  plugin,
  expanded,
  onExpand,
  onToggle,
  onEdit,
  onDelete,
}: {
  plugin: PluginMeta;
  expanded: boolean;
  onExpand: () => void;
  onToggle: () => void;
  onEdit: () => void;
  onDelete: () => void;
}) {
  return (
    <div className="flex items-center gap-3 px-4 py-3 hover:bg-base-50 transition-colors">
      {/* Icon */}
      <div className="w-9 h-9 rounded-lg bg-secondary/10 flex items-center justify-center shrink-0">
        <Code2 className="w-4 h-4 text-secondary" />
      </div>

      {/* Name + expand */}
      <button className="flex-1 text-left min-w-0 group" onClick={onExpand}>
        <div className="flex items-center gap-1.5">
          <span className="text-sm font-medium truncate">{plugin.name}</span>
          <span className="badge badge-ghost badge-xs font-mono shrink-0">
            v{plugin.version}
          </span>
          {expanded ? (
            <ChevronUp className="w-3.5 h-3.5 text-base-content/30 shrink-0" />
          ) : (
            <ChevronDown className="w-3.5 h-3.5 text-base-content/20 group-hover:text-base-content/40 shrink-0" />
          )}
        </div>
        <p className="text-xs text-base-content/40 truncate">{plugin.author}</p>
      </button>

      {/* Slots preview */}
      <div className="hidden sm:flex flex-wrap gap-1 max-w-[140px]">
        {(plugin.slots ?? []).slice(0, 2).map((s) => (
          <span key={s} className="badge badge-outline badge-xs">
            {s}
          </span>
        ))}
        {(plugin.slots ?? []).length > 2 && (
          <span className="badge badge-ghost badge-xs">
            +{(plugin.slots ?? []).length - 2}
          </span>
        )}
      </div>

      {/* Status */}
      <span
        className={`badge badge-sm shrink-0 gap-1 ${plugin.enabled ? "badge-success" : "badge-ghost"}`}
      >
        {plugin.enabled ? (
          <CheckCircle2 className="w-3 h-3" />
        ) : (
          <PowerOff className="w-3 h-3" />
        )}
        {plugin.enabled ? "启用" : "禁用"}
      </span>

      {/* Toggle */}
      <Switch
        checked={plugin.enabled}
        onChange={onToggle}
        className={`${plugin.enabled ? "bg-secondary" : "bg-base-300"} relative inline-flex h-5 w-9 shrink-0 items-center rounded-full transition-colors focus:outline-none focus:ring-2 focus:ring-secondary/30`}
      >
        <span
          className={`${plugin.enabled ? "translate-x-5" : "translate-x-1"} inline-block h-3.5 w-3.5 transform rounded-full bg-white shadow transition-transform`}
        />
      </Switch>

      {/* Action buttons */}
      <div className="flex items-center gap-0.5 shrink-0">
        <button
          onClick={onEdit}
          className="btn btn-ghost btn-xs btn-square"
          title="编辑"
        >
          <Pencil className="w-3.5 h-3.5" />
        </button>
        <button
          onClick={onDelete}
          className="btn btn-ghost btn-xs btn-square text-error hover:text-error"
          title="移除"
        >
          <Trash2 className="w-3.5 h-3.5" />
        </button>
      </div>
    </div>
  );
}
