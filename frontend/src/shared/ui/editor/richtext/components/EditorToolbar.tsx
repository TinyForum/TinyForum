// components/EditorToolbar.tsx
import React from "react";
import { Editor } from "@tiptap/react";
import { EditorToolbarButton } from "./EditorToolbarButton";
import {
  Bold,
  Italic,
  Strikethrough,
  Code,
  Heading2,
  Heading3,
  List,
  ListOrdered,
  Quote,
  Minus,
  Undo,
  Redo,
  Link as LinkIcon,
  Image as ImageIcon,
  Code2,
} from "lucide-react";
import {
  ToolbarButtonConfig,
  ToolbarGroupConfig,
} from "../types/richtext.type";

// 默认按钮配置（修正 action 签名）
export const DEFAULT_TOOLBAR_BUTTONS: Record<string, ToolbarButtonConfig> = {
  bold: {
    id: "bold",
    icon: <Bold className="w-3.5 h-3.5" />,
    title: "加粗",
    isActive: (editor) => editor.isActive("bold"),
    action: (editor) => editor.chain().focus().toggleBold().run(),
  },
  italic: {
    id: "italic",
    icon: <Italic className="w-3.5 h-3.5" />,
    title: "斜体",
    isActive: (editor) => editor.isActive("italic"),
    action: (editor) => editor.chain().focus().toggleItalic().run(),
  },
  strike: {
    id: "strike",
    icon: <Strikethrough className="w-3.5 h-3.5" />,
    title: "删除线",
    isActive: (editor) => editor.isActive("strike"),
    action: (editor) => editor.chain().focus().toggleStrike().run(),
  },
  inlineCode: {
    id: "inlineCode",
    icon: <Code className="w-3.5 h-3.5" />,
    title: "行内代码",
    isActive: (editor) => editor.isActive("code"),
    action: (editor) => editor.chain().focus().toggleCode().run(),
  },
  heading2: {
    id: "heading2",
    icon: <Heading2 className="w-3.5 h-3.5" />,
    title: "标题2",
    isActive: (editor) => editor.isActive("heading", { level: 2 }),
    action: (editor) =>
      editor.chain().focus().toggleHeading({ level: 2 }).run(),
  },
  heading3: {
    id: "heading3",
    icon: <Heading3 className="w-3.5 h-3.5" />,
    title: "标题3",
    isActive: (editor) => editor.isActive("heading", { level: 3 }),
    action: (editor) =>
      editor.chain().focus().toggleHeading({ level: 3 }).run(),
  },
  bulletList: {
    id: "bulletList",
    icon: <List className="w-3.5 h-3.5" />,
    title: "无序列表",
    isActive: (editor) => editor.isActive("bulletList"),
    action: (editor) => editor.chain().focus().toggleBulletList().run(),
  },
  orderedList: {
    id: "orderedList",
    icon: <ListOrdered className="w-3.5 h-3.5" />,
    title: "有序列表",
    isActive: (editor) => editor.isActive("orderedList"),
    action: (editor) => editor.chain().focus().toggleOrderedList().run(),
  },
  blockquote: {
    id: "blockquote",
    icon: <Quote className="w-3.5 h-3.5" />,
    title: "引用",
    isActive: (editor) => editor.isActive("blockquote"),
    action: (editor) => editor.chain().focus().toggleBlockquote().run(),
  },
  horizontalRule: {
    id: "horizontalRule",
    icon: <Minus className="w-3.5 h-3.5" />,
    title: "分割线",
    action: (editor) => editor.chain().focus().setHorizontalRule().run(),
  },
  undo: {
    id: "undo",
    icon: <Undo className="w-3.5 h-3.5" />,
    title: "撤销",
    isDisabled: (editor) => !editor.can().undo(),
    action: (editor) => editor.chain().focus().undo().run(),
  },
  redo: {
    id: "redo",
    icon: <Redo className="w-3.5 h-3.5" />,
    title: "重做",
    isDisabled: (editor) => !editor.can().redo(),
    action: (editor) => editor.chain().focus().redo().run(),
  },
  link: {
    id: "link",
    icon: <LinkIcon className="w-3.5 h-3.5" />,
    title: "添加链接",
    isActive: (editor) => editor.isActive("link"),
    // 支持第二个参数作为自定义处理器
    action: (editor, customHandler) => {
      if (customHandler && typeof customHandler === "function") {
        customHandler(editor);
      } else {
        const url = window.prompt("输入链接地址");
        if (url) editor.chain().focus().setLink({ href: url }).run();
      }
    },
  },
  image: {
    id: "image",
    icon: <ImageIcon className="w-3.5 h-3.5" />,
    title: "插入图片",
    action: (editor, customHandler) => {
      if (customHandler && typeof customHandler === "function") {
        customHandler(editor);
      } else {
        const url = window.prompt("输入图片地址");
        if (url) editor.chain().focus().setImage({ src: url }).run();
      }
    },
  },
  codeBlock: {
    id: "codeBlock",
    icon: <Code2 className="w-3.5 h-3.5" />,
    title: "代码块",
    isActive: (editor) => editor.isActive("codeBlock"),
    action: (editor) => editor.chain().focus().toggleCodeBlock().run(),
  },
};

// 默认分组
export const DEFAULT_TOOLBAR_GROUPS: ToolbarGroupConfig[] = [
  { id: "text-style", buttons: ["bold", "italic", "strike", "inlineCode"] },
  { id: "headings", buttons: ["heading2", "heading3"] },
  {
    id: "lists",
    buttons: ["bulletList", "orderedList", "blockquote", "horizontalRule"],
  },
  { id: "history", buttons: ["undo", "redo"] },
  { id: "insert", buttons: ["link", "image", "codeBlock"] },
];

interface EditorToolbarProps {
  editor: Editor;
  buttons?: Record<string, ToolbarButtonConfig>;
  groups?: ToolbarGroupConfig[];
  onInsertLink?: (editor: Editor) => void;
  onInsertImage?: (editor: Editor) => void;
  characterCount?: number;
  maxLength?: number;
}

export const EditorToolbar: React.FC<EditorToolbarProps> = ({
  editor,
  buttons = DEFAULT_TOOLBAR_BUTTONS,
  groups = DEFAULT_TOOLBAR_GROUPS,
  onInsertLink,
  onInsertImage,
  characterCount = 0,
  maxLength,
}) => {
  // 获取按钮实际执行函数（处理自定义回调）
  const getButtonAction = (button: ToolbarButtonConfig) => {
    if (button.id === "link" && onInsertLink) {
      return () => button.action(editor, onInsertLink);
    }
    if (button.id === "image" && onInsertImage) {
      return () => button.action(editor, onInsertImage);
    }
    return () => button.action(editor);
  };

  return (
    <div className="flex flex-wrap items-center gap-0.5 p-2 border-b border-base-300 bg-base-200">
      {groups.map((group) => (
        <div
          key={group.id}
          className="flex gap-1 mr-2 border-r border-base-300 pr-2 last:border-r-0"
        >
          {group.buttons.map((buttonId) => {
            const button = buttons[buttonId];
            if (!button) return null;

            return (
              <EditorToolbarButton
                key={button.id}
                onClick={getButtonAction(button)}
                active={button.isActive?.(editor)}
                disabled={button.isDisabled?.(editor)}
                title={button.title}
              >
                {button.icon}
              </EditorToolbarButton>
            );
          })}
        </div>
      ))}

      {maxLength && maxLength > 0 && (
        <div className="ml-auto text-xs text-base-content/40 pr-2">
          {characterCount}/{maxLength}
        </div>
      )}
    </div>
  );
};
