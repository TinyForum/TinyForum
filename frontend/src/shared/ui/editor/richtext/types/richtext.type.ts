import { Editor } from "@tiptap/react";
import { Extension } from "@tiptap/core";

export type EditorMode = "rich" | "markdown";

export interface ToolbarButtonConfig {
  id: string;
  icon: React.ReactNode;
  title?: string;
  isActive?: (editor: Editor) => boolean;
  isDisabled?: (editor: Editor) => boolean;
  // 自定义处理器，用于链接/图片等需要交互的场景
  action: (editor: Editor, customHandler?: (editor: Editor) => void) => void;
  group?: string;
}

export interface ToolbarGroupConfig {
  id: string;
  buttons: string[]; // button ids
}

export interface RichTextEditorProps {
  value?: string;
  onChange?: (html: string) => void;
  placeholder?: string;
  disabled?: boolean;
  className?: string;
  defaultMode?: EditorMode;
  maxLength?: number;
  extensions?: Extension[];
  // 修改为 Record 类型，方便通过 id 快速查找
  toolbarButtons?: Record<string, ToolbarButtonConfig>;
  toolbarGroups?: ToolbarGroupConfig[];
  onInsertLink?: (editor: Editor) => void;
  onInsertImage?: (editor: Editor) => void;
  markdownParser?: (markdown: string) => Promise<string>;
  htmlToMarkdown?: (html: string) => string;
}
