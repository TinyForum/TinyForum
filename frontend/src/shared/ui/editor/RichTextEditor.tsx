// RichTextEditor.tsx
import React, { useState, useEffect, useCallback } from "react";
import {
  useEditor,
  EditorContent,
  BubbleMenu,
  FloatingMenu,
} from "@tiptap/react";
import StarterKit from "@tiptap/starter-kit";
import Placeholder from "@tiptap/extension-placeholder";
import Link from "@tiptap/extension-link";
import Image from "@tiptap/extension-image";
import CodeBlockLowlight from "@tiptap/extension-code-block-lowlight";
import { all, createLowlight } from "lowlight";
import TurndownService from "turndown";
import { marked } from "marked";

// 配置 lowlight（代码高亮）
const lowlight = createLowlight(all);

// 配置 Turndown（HTML -> Markdown）
const turndownService = new TurndownService({
  headingStyle: "atx",
  codeBlockStyle: "fenced",
  emDelimiter: "*",
});
// 添加代码块规则
turndownService.addRule("fencedCodeBlock", {
  filter: "pre",
  replacement(content, node) {
    const codeEl = node.querySelector("code");
    const language = codeEl?.className?.match(/language-(\w+)/)?.[1] || "";
    return `\`\`\`${language}\n${codeEl?.textContent || content}\n\`\`\`\n\n`;
  },
});

// 配置 marked（Markdown -> HTML）
marked.setOptions({
  gfm: true,
  breaks: true,
});

export type EditorMode = "rich" | "markdown";

interface RichTextEditorProps {
  value?: string; // HTML 格式内容（富文本模式存储）
  onChange?: (html: string) => void;
  placeholder?: string;
  disabled?: boolean;
  className?: string;
  defaultMode?: EditorMode;
}

export const RichTextEditor: React.FC<RichTextEditorProps> = ({
  value = "",
  onChange,
  placeholder = "写点什么...",
  disabled = false,
  className = "",
  defaultMode = "rich",
}) => {
  const [mode, setMode] = useState<EditorMode>(defaultMode);
  const [markdownContent, setMarkdownContent] = useState("");

  // TipTap 编辑器
  const editor = useEditor({
    extensions: [
      StarterKit.configure({
        codeBlock: false, // 使用 CodeBlockLowlight 替代
      }),
      CodeBlockLowlight.configure({ lowlight }),
      Placeholder.configure({ placeholder }),
      Link.configure({
        openOnClick: false,
        HTMLAttributes: { target: "_blank", rel: "noopener noreferrer" },
      }),
      Image.configure({
        inline: true,
        allowBase64: true,
      }),
    ],
    content: value,
    editable: !disabled && mode === "rich",
    onUpdate: ({ editor }) => {
      const html = editor.getHTML();
      onChange?.(html);
      // 同步 Markdown 内容（供切换时使用）
      if (mode === "markdown") {
        setMarkdownContent(turndownService.turndown(html));
      }
    },
  });

  // 初始化 / value 变化时同步编辑器内容
  useEffect(() => {
    if (editor && value !== editor.getHTML() && mode === "rich") {
      editor.commands.setContent(value);
    }
  }, [value, editor, mode]);

  // 切换模式时，转换内容格式
  const switchToMarkdown = useCallback(() => {
    if (!editor) return;
    const html = editor.getHTML();
    const md = turndownService.turndown(html);
    setMarkdownContent(md);
    setMode("markdown");
  }, [editor]);

  const handleMarkdownChange = useCallback(
    (e: React.ChangeEvent<HTMLTextAreaElement>) => {
      const newMd = e.target.value;
      setMarkdownContent(newMd);
      (marked.parse(newMd) as Promise<string>)
        .then((html: string) => onChange?.(html))
        .catch((err: Error) => console.error("Markdown parse error:", err));
    },
    [onChange],
  );

  const switchToRich = useCallback(async () => {
    if (!editor) return;
    try {
      const html = await (marked.parse(markdownContent) as Promise<string>);
      editor.commands.setContent(html);
      onChange?.(html);
      setMode("rich");
    } catch (error) {
      console.error("Failed to convert markdown to HTML", error);
    }
  }, [editor, markdownContent, onChange]);
  // 工具栏按钮定义
  const ToolbarButton = ({
    onClick,
    active,
    disabled: btnDisabled,
    children,
  }: any) => (
    <button
      type="button"
      onClick={onClick}
      disabled={disabled || btnDisabled}
      className={`px-2 py-1 rounded text-sm font-medium transition ${
        active
          ? "bg-primary text-primary-content"
          : "bg-base-200 hover:bg-base-300 text-base-content"
      } ${disabled || btnDisabled ? "opacity-50 cursor-not-allowed" : ""}`}
    >
      {children}
    </button>
  );

  return (
    <div
      className={`border border-base-300 rounded-lg bg-base-100 ${className}`}
    >
      {/* 模式切换栏 */}
      <div className="flex items-center justify-between border-b border-base-300 p-2 bg-base-200 rounded-t-lg">
        <div className="flex gap-2">
          <button
            onClick={() => (mode === "markdown" ? switchToRich() : null)}
            className={`btn btn-xs ${mode === "rich" ? "btn-primary" : "btn-ghost"}`}
          >
            富文本
          </button>
          <button
            onClick={() => (mode === "rich" ? switchToMarkdown() : null)}
            className={`btn btn-xs ${mode === "markdown" ? "btn-primary" : "btn-ghost"}`}
          >
            Markdown
          </button>
        </div>
        {mode === "rich" && editor && (
          <div className="flex gap-1 flex-wrap">
            <ToolbarButton
              onClick={() => editor.chain().focus().toggleBold().run()}
              active={editor.isActive("bold")}
            >
              粗体
            </ToolbarButton>
            <ToolbarButton
              onClick={() => editor.chain().focus().toggleItalic().run()}
              active={editor.isActive("italic")}
            >
              斜体
            </ToolbarButton>
            <ToolbarButton
              onClick={() =>
                editor.chain().focus().toggleHeading({ level: 2 }).run()
              }
              active={editor.isActive("heading", { level: 2 })}
            >
              H2
            </ToolbarButton>
            <ToolbarButton
              onClick={() =>
                editor.chain().focus().toggleHeading({ level: 3 }).run()
              }
              active={editor.isActive("heading", { level: 3 })}
            >
              H3
            </ToolbarButton>
            <ToolbarButton
              onClick={() => editor.chain().focus().toggleBulletList().run()}
              active={editor.isActive("bulletList")}
            >
              列表
            </ToolbarButton>
            <ToolbarButton
              onClick={() => editor.chain().focus().toggleCodeBlock().run()}
              active={editor.isActive("codeBlock")}
            >
              代码块
            </ToolbarButton>
            <ToolbarButton
              onClick={() => {
                const url = window.prompt("输入链接地址");
                if (url) editor.chain().focus().setLink({ href: url }).run();
              }}
              active={editor.isActive("link")}
            >
              链接
            </ToolbarButton>
            <ToolbarButton
              onClick={() => {
                const url = window.prompt("输入图片地址");
                if (url) editor.chain().focus().setImage({ src: url }).run();
              }}
            >
              图片
            </ToolbarButton>
          </div>
        )}
      </div>

      {/* 编辑区域 */}
      <div className="p-4 min-h-[300px]">
        {mode === "rich" && (
          <>
            <EditorContent editor={editor} className="prose max-w-none" />
            {editor && (
              <>
                <BubbleMenu editor={editor} tippyOptions={{ duration: 100 }}>
                  <div className="flex gap-1 bg-base-100 shadow-lg rounded border border-base-300 p-1">
                    <ToolbarButton
                      onClick={() => editor.chain().focus().toggleBold().run()}
                    >
                      粗体
                    </ToolbarButton>
                    <ToolbarButton
                      onClick={() =>
                        editor.chain().focus().toggleItalic().run()
                      }
                    >
                      斜体
                    </ToolbarButton>
                    <ToolbarButton
                      onClick={() => editor.chain().focus().toggleCode().run()}
                    >
                      代码
                    </ToolbarButton>
                  </div>
                </BubbleMenu>
                <FloatingMenu editor={editor} tippyOptions={{ duration: 100 }}>
                  <div className="bg-base-100 shadow-lg rounded border border-base-300 p-1">
                    <ToolbarButton
                      onClick={() =>
                        editor.chain().focus().toggleHeading({ level: 2 }).run()
                      }
                    >
                      H2
                    </ToolbarButton>
                    <ToolbarButton
                      onClick={() =>
                        editor.chain().focus().toggleBulletList().run()
                      }
                    >
                      列表
                    </ToolbarButton>
                  </div>
                </FloatingMenu>
              </>
            )}
          </>
        )}

        {mode === "markdown" && (
          <textarea
            value={markdownContent}
            onChange={handleMarkdownChange}
            placeholder={placeholder}
            disabled={disabled}
            className="w-full h-[300px] p-3 font-mono text-sm border-none focus:ring-0 bg-transparent resize-y"
          />
        )}
      </div>
    </div>
  );
};
