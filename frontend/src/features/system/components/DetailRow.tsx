import { Puzzle, Plus } from "lucide-react";

export function DetailRow({ label, value }: { label: string; value: string }) {
  return (
    <div>
      <span className="text-xs text-base-content/40 uppercase tracking-wide">
        {label}
      </span>
      <p className="mt-0.5 text-base-content/70 text-xs">{value}</p>
    </div>
  );
}

export function EmptyState({
  keyword,
  onInstall,
}: {
  keyword: string;
  onInstall: () => void;
}) {
  return (
    <div className="flex flex-col items-center justify-center h-40 gap-3 text-base-content/40">
      <Puzzle className="w-8 h-8" />
      <p className="text-sm">
        {keyword ? `没有找到"${keyword}"相关插件` : "尚未安装任何插件"}
      </p>
      {!keyword && (
        <button onClick={onInstall} className="btn btn-ghost btn-xs gap-1">
          <Plus className="w-3 h-3" /> 安装第一个插件
        </button>
      )}
    </div>
  );
}
