// components/question/BoardSelector.tsx
'use client';

import { FolderIcon } from '@heroicons/react/24/outline';
import { Board } from '@/lib/api/types';

interface BoardSelectorProps {
  boards: Board[];
  selectedBoardId: number;
  onBoardChange: (boardId: number) => void;
  loading?: boolean;
}

export function BoardSelector({ boards, selectedBoardId, onBoardChange, loading }: BoardSelectorProps) {
  return (
    <div>
      <label className="block text-sm font-medium text-gray-700 mb-2">
        <FolderIcon className="w-4 h-4 inline mr-1" />
        选择板块 <span className="text-red-500">*</span>
      </label>
      <select
        value={selectedBoardId}
        onChange={(e) => onBoardChange(Number(e.target.value))}
        disabled={loading}
        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent outline-none disabled:bg-gray-100"
      >
        <option value={0} disabled>
          {loading ? '加载中...' : '请选择板块'}
        </option>
        {boards.map((board) => (
          <option key={board.id} value={board.id}>
            {board.name} {board.description ? `- ${board.description}` : ''}
            {board.post_count !== undefined && ` (${board.post_count} 帖子)`}
          </option>
        ))}
      </select>
      {selectedBoardId === 0 && !loading && (
        <p className="mt-1 text-sm text-red-500">请选择板块</p>
      )}
      {boards.length === 0 && !loading && (
        <p className="mt-1 text-sm text-amber-500">暂无可用板块，请联系管理员</p>
      )}
    </div>
  );
}