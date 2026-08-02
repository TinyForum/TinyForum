// hooks/useEditorCore.ts
import { useEditor, Editor } from "@tiptap/react";
import { useEffect, useCallback, useRef } from "react";
import { Extension } from "@tiptap/core";
import StarterKit from "@tiptap/starter-kit";
import Placeholder from "@tiptap/extension-placeholder";
import Link from "@tiptap/extension-link";
import Image from "@tiptap/extension-image";
import CodeBlockLowlight from "@tiptap/extension-code-block-lowlight";
import CharacterCount from "@tiptap/extension-character-count";
import { all, createLowlight } from "lowlight";

const lowlight = createLowlight(all);

const DEFAULT_EXTENSIONS = [
  StarterKit.configure({ codeBlock: false }),
  CodeBlockLowlight.configure({ lowlight }),
  Placeholder,
  Link.configure({
    openOnClick: false,
    HTMLAttributes: { target: "_blank", rel: "noopener noreferrer" },
  }),
  Image.configure({ inline: true, allowBase64: true }),
];

interface UseEditorCoreOptions {
  content: string;
  editable: boolean;
  maxLength?: number;
  placeholder?: string;
  customExtensions?: Extension[];
  onUpdate: (editor: Editor) => void;
}

export const useEditorCore = ({
  content,
  editable,
  maxLength,
  placeholder,
  customExtensions = [],
  onUpdate,
}: UseEditorCoreOptions) => {
  const extensions = [...DEFAULT_EXTENSIONS, ...customExtensions];

  if (maxLength && maxLength > 0) {
    extensions.push(CharacterCount.configure({ limit: maxLength }));
  }

  const placeholderExtension = Placeholder.configure({
    placeholder: placeholder || "写点什么...",
  });
  const placeholderIndex = extensions.findIndex(
    (ext) => ext.name === "placeholder",
  );
  if (placeholderIndex !== -1) {
    extensions[placeholderIndex] = placeholderExtension;
  } else {
    extensions.push(placeholderExtension);
  }

  // 防止 setContent 触发的 onUpdate 回环
  const isSettingRef = useRef(false);

  const editor = useEditor({
    extensions,
    content,
    editable,
    onUpdate: ({ editor }) => {
      if (isSettingRef.current) return;
      onUpdate(editor);
    },
  });

  // 同步外部内容变化到编辑器
  useEffect(() => {
    if (!editor || editor.isDestroyed) return;
    const currentHTML = editor.getHTML();
    if (content !== currentHTML) {
      isSettingRef.current = true;
      editor.commands.setContent(content);
      isSettingRef.current = false;
    }
  }, [content, editor]);

  // 同步可编辑状态
  useEffect(() => {
    if (editor) {
      editor.setEditable(editable);
    }
  }, [editor, editable]);

  // 清理编辑器实例
  useEffect(() => {
    return () => {
      if (editor && !editor.isDestroyed) {
        editor.destroy();
      }
    };
  }, [editor]);

  const getCharacterCount = useCallback(() => {
    return editor?.storage.characterCount?.characters() || 0;
  }, [editor]);

  return { editor, getCharacterCount };
};
