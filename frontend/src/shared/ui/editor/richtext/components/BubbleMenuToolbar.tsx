// components/BubbleMenuToolbar.tsx
import React from "react";
import { Editor } from "@tiptap/react";
import { BubbleMenu } from "@tiptap/react";
import { EditorToolbarButton } from "./EditorToolbarButton";
import { Bold, Italic, Code } from "lucide-react";

interface BubbleMenuToolbarProps {
  editor: Editor;
}

export const BubbleMenuToolbar: React.FC<BubbleMenuToolbarProps> = ({
  editor,
}) => (
  <BubbleMenu editor={editor} tippyOptions={{ duration: 100 }}>
    <div className="flex gap-1 bg-base-100 shadow-lg rounded border border-base-300 p-1">
      <EditorToolbarButton
        onClick={() => editor.chain().focus().toggleBold().run()}
        active={editor.isActive("bold")}
      >
        <Bold className="w-3 h-3" />
      </EditorToolbarButton>
      <EditorToolbarButton
        onClick={() => editor.chain().focus().toggleItalic().run()}
        active={editor.isActive("italic")}
      >
        <Italic className="w-3 h-3" />
      </EditorToolbarButton>
      <EditorToolbarButton
        onClick={() => editor.chain().focus().toggleCode().run()}
        active={editor.isActive("code")}
      >
        <Code className="w-3 h-3" />
      </EditorToolbarButton>
    </div>
  </BubbleMenu>
);
