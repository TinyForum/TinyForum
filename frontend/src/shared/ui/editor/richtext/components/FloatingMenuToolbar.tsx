// components/FloatingMenuToolbar.tsx
import React from "react";
import { Editor } from "@tiptap/react";
import { FloatingMenu } from "@tiptap/react";
import { EditorToolbarButton } from "./EditorToolbarButton";

interface FloatingMenuToolbarProps {
  editor: Editor;
}

export const FloatingMenuToolbar: React.FC<FloatingMenuToolbarProps> = ({
  editor,
}) => (
  <FloatingMenu editor={editor} tippyOptions={{ duration: 100 }}>
    <div className="bg-base-100 shadow-lg rounded border border-base-300 p-1">
      <EditorToolbarButton
        onClick={() => editor.chain().focus().toggleHeading({ level: 2 }).run()}
      >
        H2
      </EditorToolbarButton>
      <EditorToolbarButton
        onClick={() => editor.chain().focus().toggleBulletList().run()}
      >
        列表
      </EditorToolbarButton>
    </div>
  </FloatingMenu>
);
