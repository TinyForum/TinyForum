"use client";

import { RichTextEditor } from "@/shared/ui/editor/richtext/RichTextEditor";
import { useTranslations } from "next-intl";

export interface PostEditorProps {
  content: string;
  setContent: (value: string) => void;
  placeholder?: string;
}

export function PostEditor({
  content,
  setContent,
  placeholder,
}: PostEditorProps) {
  const t = useTranslations("Post");

  return (
    <div className="border border-base-300 rounded-lg bg-base-100 shadow-sm flex flex-col h-full">
      <div className="flex-1 min-h-[400px]">
        <RichTextEditor
          key="rich-editor"
          value={content}
          onChange={setContent}
          placeholder={placeholder || t("write_post_content")}
          maxLength={20000}
          defaultMode="rich"
        />
      </div>
    </div>
  );
}
