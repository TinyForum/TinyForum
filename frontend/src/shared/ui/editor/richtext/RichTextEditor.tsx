// components/RichTextEditor.tsx
"use client";

import React from "react";
import { EditorContent } from "@tiptap/react";
import { RichTextEditorProps } from "./types/richtext.type";
import { BubbleMenuToolbar } from "./components/BubbleMenuToolbar";
import { EditorToolbar } from "./components/EditorToolbar";
import { FloatingMenuToolbar } from "./components/FloatingMenuToolbar";
import { MarkdownEditor } from "./components/MarkdownEditor";
import { ModeSwitcher } from "./components/ModeSwitcher";
import { useEditorCore } from "./hooks/useEditorCore.ts";
import { useMarkdownMode } from "./hooks/useMarkdownMode";

export const RichTextEditor: React.FC<RichTextEditorProps> = ({
  value = "",
  onChange,
  placeholder = "写点什么...",
  disabled = false,
  className = "",
  defaultMode = "rich",
  maxLength = 50000,
  extensions = [],
  toolbarButtons,
  toolbarGroups,
  onInsertLink,
  onInsertImage,
  markdownParser,
  htmlToMarkdown,
}) => {
  const { editor, getCharacterCount } = useEditorCore({
    content: value,
    editable: !disabled && defaultMode === "rich",
    maxLength,
    placeholder,
    customExtensions: extensions,
    onUpdate: (editor) => {
      const html = editor.getHTML();
      onChange?.(html);
    },
  });

  const {
    mode,
    markdownContent,
    switchToMarkdown,
    switchToRich,
    handleMarkdownChange,
  } = useMarkdownMode({
    editor,
    htmlToMarkdown,
    markdownParser,
    onChange,
  });

  // 当模式变化时更新编辑器的可编辑状态
  React.useEffect(() => {
    if (editor) {
      editor.setEditable(!disabled && mode === "rich");
    }
  }, [editor, disabled, mode]);

  if (!editor) return null;

  const characterCount = getCharacterCount();

  return (
    <div
      className={`border border-base-300 rounded-xl overflow-hidden focus-within:border-primary transition-colors ${className}`}
    >
      {/* 工具栏区域 */}
      <div className="flex flex-wrap items-center gap-0.5 p-2 border-b border-base-300 bg-base-200">
        <ModeSwitcher
          mode={mode}
          onSwitchToRich={switchToRich}
          onSwitchToMarkdown={switchToMarkdown}
        />

        {mode === "rich" && (
          <EditorToolbar
            editor={editor}
            buttons={toolbarButtons}
            groups={toolbarGroups}
            onInsertLink={onInsertLink}
            onInsertImage={onInsertImage}
            characterCount={characterCount}
            maxLength={maxLength}
          />
        )}

        {mode === "markdown" && (
          <div className="ml-auto text-xs text-base-content/40">
            Markdown 模式
          </div>
        )}
      </div>

      {/* 编辑区域 */}
      <div className="bg-base-100">
        <div style={{ display: mode === "rich" ? "block" : "none" }}>
          <EditorContent
            editor={editor}
            className="tiptap p-4 min-h-[300px] focus:outline-none"
          />
          <BubbleMenuToolbar editor={editor} />
          <FloatingMenuToolbar editor={editor} />
        </div>

        {mode === "markdown" && (
          <MarkdownEditor
            value={markdownContent}
            onChange={handleMarkdownChange}
            placeholder={placeholder}
            disabled={disabled}
          />
        )}
      </div>
    </div>
  );
};
