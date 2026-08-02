"use client";

import { useState, useMemo } from "react";
import { Eye, Edit } from "lucide-react";
import DOMPurify from "dompurify";
import { RichTextEditor } from "@/shared/ui/editor/richtext/RichTextEditor";

// ---------- 右侧：编辑器 + 预览组件 ----------
export interface PostEditorProps {
  content: string;
  setContent: (value: string) => void;
  placeholder?: string;
}

export function PostEditor({ content, setContent }: PostEditorProps) {
  const [mode, setMode] = useState<"edit" | "preview">("edit");
  const sanitizedHtml = useMemo(() => DOMPurify.sanitize(content), [content]);

  return (
    <div className="border border-base-300 rounded-lg bg-base-100 shadow-sm flex flex-col h-full">
      <div className="flex justify-between items-center border-b border-base-300 p-2 bg-base-200 rounded-t-lg">
        <div className="flex gap-1">
          <button
            type="button"
            onClick={() => setMode("edit")}
            className={`btn btn-xs gap-1 ${mode === "edit" ? "btn-primary" : "btn-ghost"}`}
          >
            <Edit className="w-3 h-3" />
            编辑
          </button>
          <button
            type="button"
            onClick={() => setMode("preview")}
            className={`btn btn-xs gap-1 ${mode === "preview" ? "btn-primary" : "btn-ghost"}`}
          >
            <Eye className="w-3 h-3" />
            预览
          </button>
        </div>
        {mode === "preview" && (
          <span className="text-xs text-base-content/50">纯预览模式</span>
        )}
      </div>

      <div className="flex-1 min-h-[400px]">
        {mode === "edit" ? (
          // 添加 key={mode} 强制重新创建编辑器实例，避免 DOM 冲突
          <RichTextEditor
            key="rich-editor"
            value={content}
            onChange={setContent}
            placeholder="撰写帖子内容..."
            maxLength={20000}
            defaultMode="rich"
          />
        ) : (
          <div className="prose prose-sm max-w-none p-4 overflow-auto h-full">
            {content ? (
              <div dangerouslySetInnerHTML={{ __html: sanitizedHtml }} />
            ) : (
              <p className="text-base-content/40">暂无内容</p>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
