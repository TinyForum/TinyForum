// components/question/BoardSelector.tsx
"use client";

import { Board } from "@/lib/api";
import { HashtagIcon } from "@heroicons/react/24/outline";

interface BoardSelectorProps {
  boards: Board[];
  selectedBoardId: number;
  onBoardChange: (id: number) => void;
  loading?: boolean;
}

export function BoardSelector({
  boards,
  selectedBoardId,
  onBoardChange,
  loading,
}: BoardSelectorProps) {
  if (loading) {
    return (
      <div className="select select-bordered w-full bg-base-100">
        <option>加载中...</option>
      </div>
    );
  }

  return (
    <div className="form-control">
      <select
        value={selectedBoardId || ""}
        onChange={(e) => onBoardChange(Number(e.target.value))}
        className="select select-bordered w-full bg-base-100 focus:select-primary transition-colors"
        required
      >
        <option value="" disabled>
          请选择板块
        </option>
        {boards.map((board) => (
          <option key={board.id} value={board.id}>
            {board.name} {board.description ? `- ${board.description}` : ""}
          </option>
        ))}
      </select>
      <label className="label">
        <span className="label-text-alt text-base-content/50">
          <HashtagIcon className="w-3 h-3 inline mr-1" />
          选择正确的板块能让问题更快得到回答
        </span>
      </label>
    </div>
  );
}
