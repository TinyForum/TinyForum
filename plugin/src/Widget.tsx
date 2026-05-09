import React, { useState, useEffect } from "react";

interface Props {
  // PluginSlot 可以通过 slotProps 传入额外数据
  postId?: string;
}

export function MyWidget({ postId }: Props) {
  const [count, setCount] = useState(0);

  return (
    <div className="card bg-primary/10 border border-primary/20 p-3 text-sm">
      <p className="font-semibold text-primary">📌 我的插件</p>
      {postId && (
        <p className="text-xs text-base-content/50">当前帖子：{postId}</p>
      )}
      <button
        className="btn btn-xs btn-primary mt-2"
        onClick={() => setCount((c) => c + 1)}
      >
        点击了 {count} 次
      </button>
    </div>
  );
}
