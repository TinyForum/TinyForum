// "use client";

// import React, { useState, useEffect, useCallback } from "react";
// import {
//   useEditor,
//   EditorContent,
//   BubbleMenu,
//   FloatingMenu,
// } from "@tiptap/react";
// import StarterKit from "@tiptap/starter-kit";
// import Placeholder from "@tiptap/extension-placeholder";
// import Link from "@tiptap/extension-link";
// import Image from "@tiptap/extension-image";
// import CodeBlockLowlight from "@tiptap/extension-code-block-lowlight";
// import CharacterCount from "@tiptap/extension-character-count";
// import { all, createLowlight } from "lowlight";
// import TurndownService from "turndown";
// import { marked } from "marked";
// import {
//   Bold,
//   Italic,
//   Strikethrough,
//   Code,
//   Heading2,
//   Heading3,
//   List,
//   ListOrdered,
//   Quote,
//   Minus,
//   Undo,
//   Redo,
//   Link as LinkIcon,
//   Image as ImageIcon,
//   Code2,
// } from "lucide-react";

// const lowlight = createLowlight(all);

// const turndownService = new TurndownService({
//   headingStyle: "atx",
//   codeBlockStyle: "fenced",
//   emDelimiter: "*",
// });
// turndownService.addRule("fencedCodeBlock", {
//   filter: "pre",
//   replacement(content, node) {
//     const codeEl = node.querySelector("code");
//     const language = codeEl?.className?.match(/language-(\w+)/)?.[1] || "";
//     return `\`\`\`${language}\n${codeEl?.textContent || content}\n\`\`\`\n\n`;
//   },
// });

// marked.setOptions({ gfm: true, breaks: true });

// export type EditorMode = "rich" | "markdown";

// interface RichTextEditorProps {
//   value?: string;
//   onChange?: (html: string) => void;
//   placeholder?: string;
//   disabled?: boolean;
//   className?: string;
//   defaultMode?: EditorMode;
//   maxLength?: number;
// }

// interface ToolbarButtonProps {
//   onClick: () => void;
//   active?: boolean;
//   disabled?: boolean;
//   children: React.ReactNode;
//   title?: string;
// }

// const ToolbarButton: React.FC<ToolbarButtonProps> = ({
//   onClick,
//   active,
//   disabled,
//   children,
//   title,
// }) => (
//   <button
//     type="button"
//     onClick={onClick}
//     disabled={disabled}
//     title={title}
//     className={`btn btn-ghost btn-xs ${active ? "btn-active text-primary" : ""}`}
//   >
//     {children}
//   </button>
// );

// export const RichTextEditor: React.FC<RichTextEditorProps> = ({
//   value = "",
//   onChange,
//   placeholder = "写点什么...",
//   disabled = false,
//   className = "",
//   defaultMode = "rich",
//   maxLength = 50000,
// }) => {
//   const [mode, setMode] = useState<EditorMode>(defaultMode);
//   const [markdownContent, setMarkdownContent] = useState("");

//   const extensions = [
//     StarterKit.configure({
//       codeBlock: false,
//     }),
//     CodeBlockLowlight.configure({ lowlight }),
//     Placeholder.configure({ placeholder }),
//     Link.configure({
//       openOnClick: false,
//       HTMLAttributes: { target: "_blank", rel: "noopener noreferrer" },
//     }),
//     Image.configure({ inline: true, allowBase64: true }),
//     ...(maxLength > 0 ? [CharacterCount.configure({ limit: maxLength })] : []),
//   ];

//   const editor = useEditor({
//     extensions,
//     content: value,
//     editable: !disabled && mode === "rich", // 初始值
//     onUpdate: ({ editor }) => {
//       const html = editor.getHTML();
//       onChange?.(html);
//       if (mode === "markdown") {
//         setMarkdownContent(turndownService.turndown(html));
//       }
//     },
//   });
//   useEffect(() => {
//     if (editor) {
//       editor.setEditable(!disabled && mode === "rich");
//     }
//   }, [editor, disabled, mode]);

//   // 同步外部 value 到编辑器
//   useEffect(() => {
//     if (editor && value !== editor.getHTML() && mode === "rich") {
//       editor.commands.setContent(value);
//     }
//   }, [value, editor, mode]);

//   // 组件卸载时销毁编辑器（必须放在条件返回之前）
//   useEffect(() => {
//     return () => {
//       if (editor && !editor.isDestroyed) {
//         editor.destroy();
//       }
//     };
//   }, [editor]);

//   // 切换到 Markdown 模式
//   const switchToMarkdown = useCallback(() => {
//     if (!editor) return;
//     const html = editor.getHTML();
//     const md = turndownService.turndown(html);
//     setMarkdownContent(md);
//     setMode("markdown");
//   }, [editor]);

//   // 切换到富文本模式
//   const switchToRich = useCallback(async () => {
//     if (!editor) return;
//     try {
//       const html = await marked.parse(markdownContent);
//       editor.commands.setContent(html);
//       onChange?.(html);
//       setMode("rich");
//     } catch (error) {
//       console.error("Markdown to HTML failed:", error);
//     }
//   }, [editor, markdownContent, onChange]);

//   // Markdown 内容变化
//   const handleMarkdownChange = useCallback(
//     (e: React.ChangeEvent<HTMLTextAreaElement>) => {
//       const newMd = e.target.value;
//       setMarkdownContent(newMd);
//       (marked.parse(newMd) as Promise<string>)
//         .then((html) => onChange?.(html))
//         .catch((err) => console.error("Markdown parse error:", err));
//     },
//     [onChange],
//   );

//   // 条件返回必须放在所有 Hooks 之后
//   if (!editor) return null;

//   const characterCount = editor.storage.characterCount?.characters() || 0;

//   return (
//     <div
//       className={`border border-base-300 rounded-xl overflow-hidden focus-within:border-primary transition-colors ${className}`}
//     >
//       {/* 工具栏区域 */}
//       <div className="flex flex-wrap items-center gap-0.5 p-2 border-b border-base-300 bg-base-200">
//         <div className="flex gap-1 mr-2 border-r border-base-300 pr-2">
//           <button
//             onClick={() => mode === "markdown" && switchToRich()}
//             className={`btn btn-xs ${mode === "rich" ? "btn-primary" : "btn-ghost"}`}
//           >
//             富文本
//           </button>
//           <button
//             onClick={() => mode === "rich" && switchToMarkdown()}
//             className={`btn btn-xs ${mode === "markdown" ? "btn-primary" : "btn-ghost"}`}
//           >
//             Markdown
//           </button>
//         </div>

//         {mode === "rich" && (
//           <>
//             <ToolbarButton
//               onClick={() => editor.chain().focus().toggleBold().run()}
//               active={editor.isActive("bold")}
//               title="加粗"
//             >
//               <Bold className="w-3.5 h-3.5" />
//             </ToolbarButton>
//             <ToolbarButton
//               onClick={() => editor.chain().focus().toggleItalic().run()}
//               active={editor.isActive("italic")}
//               title="斜体"
//             >
//               <Italic className="w-3.5 h-3.5" />
//             </ToolbarButton>
//             <ToolbarButton
//               onClick={() => editor.chain().focus().toggleStrike().run()}
//               active={editor.isActive("strike")}
//               title="删除线"
//             >
//               <Strikethrough className="w-3.5 h-3.5" />
//             </ToolbarButton>
//             <ToolbarButton
//               onClick={() => editor.chain().focus().toggleCode().run()}
//               active={editor.isActive("code")}
//               title="行内代码"
//             >
//               <Code className="w-3.5 h-3.5" />
//             </ToolbarButton>
//             <div className="divider divider-horizontal mx-0.5 h-5" />
//             <ToolbarButton
//               onClick={() =>
//                 editor.chain().focus().toggleHeading({ level: 2 }).run()
//               }
//               active={editor.isActive("heading", { level: 2 })}
//               title="标题2"
//             >
//               <Heading2 className="w-3.5 h-3.5" />
//             </ToolbarButton>
//             <ToolbarButton
//               onClick={() =>
//                 editor.chain().focus().toggleHeading({ level: 3 }).run()
//               }
//               active={editor.isActive("heading", { level: 3 })}
//               title="标题3"
//             >
//               <Heading3 className="w-3.5 h-3.5" />
//             </ToolbarButton>
//             <div className="divider divider-horizontal mx-0.5 h-5" />
//             <ToolbarButton
//               onClick={() => editor.chain().focus().toggleBulletList().run()}
//               active={editor.isActive("bulletList")}
//               title="无序列表"
//             >
//               <List className="w-3.5 h-3.5" />
//             </ToolbarButton>
//             <ToolbarButton
//               onClick={() => editor.chain().focus().toggleOrderedList().run()}
//               active={editor.isActive("orderedList")}
//               title="有序列表"
//             >
//               <ListOrdered className="w-3.5 h-3.5" />
//             </ToolbarButton>
//             <ToolbarButton
//               onClick={() => editor.chain().focus().toggleBlockquote().run()}
//               active={editor.isActive("blockquote")}
//               title="引用"
//             >
//               <Quote className="w-3.5 h-3.5" />
//             </ToolbarButton>
//             <ToolbarButton
//               onClick={() => editor.chain().focus().setHorizontalRule().run()}
//               title="分割线"
//             >
//               <Minus className="w-3.5 h-3.5" />
//             </ToolbarButton>
//             <div className="divider divider-horizontal mx-0.5 h-5" />
//             <ToolbarButton
//               onClick={() => editor.chain().focus().undo().run()}
//               disabled={!editor.can().undo()}
//               title="撤销"
//             >
//               <Undo className="w-3.5 h-3.5" />
//             </ToolbarButton>
//             <ToolbarButton
//               onClick={() => editor.chain().focus().redo().run()}
//               disabled={!editor.can().redo()}
//               title="重做"
//             >
//               <Redo className="w-3.5 h-3.5" />
//             </ToolbarButton>
//             <div className="divider divider-horizontal mx-0.5 h-5" />
//             <ToolbarButton
//               onClick={() => {
//                 const url = window.prompt("输入链接地址");
//                 if (url) editor.chain().focus().setLink({ href: url }).run();
//               }}
//               active={editor.isActive("link")}
//               title="添加链接"
//             >
//               <LinkIcon className="w-3.5 h-3.5" />
//             </ToolbarButton>
//             <ToolbarButton
//               onClick={() => {
//                 const url = window.prompt("输入图片地址");
//                 if (url) editor.chain().focus().setImage({ src: url }).run();
//               }}
//               title="插入图片"
//             >
//               <ImageIcon className="w-3.5 h-3.5" />
//             </ToolbarButton>
//             <ToolbarButton
//               onClick={() => editor.chain().focus().toggleCodeBlock().run()}
//               active={editor.isActive("codeBlock")}
//               title="代码块"
//             >
//               <Code2 className="w-3.5 h-3.5" />
//             </ToolbarButton>

//             {maxLength > 0 && (
//               <div className="ml-auto text-xs text-base-content/40 pr-2">
//                 {characterCount}/{maxLength}
//               </div>
//             )}
//           </>
//         )}

//         {mode === "markdown" && (
//           <div className="ml-auto text-xs text-base-content/40">
//             Markdown 模式
//           </div>
//         )}
//       </div>

//       {/* 编辑区域 */}
//       <div className="bg-base-100">
//         <div style={{ display: mode === "rich" ? "block" : "none" }}>
//           <EditorContent
//             editor={editor}
//             className="tiptap p-4 min-h-[300px] focus:outline-none"
//           />
//           <BubbleMenu editor={editor} tippyOptions={{ duration: 100 }}>
//             <div className="flex gap-1 bg-base-100 shadow-lg rounded border border-base-300 p-1">
//               <ToolbarButton
//                 onClick={() => editor.chain().focus().toggleBold().run()}
//               >
//                 <Bold className="w-3 h-3" />
//               </ToolbarButton>
//               <ToolbarButton
//                 onClick={() => editor.chain().focus().toggleItalic().run()}
//               >
//                 <Italic className="w-3 h-3" />
//               </ToolbarButton>
//               <ToolbarButton
//                 onClick={() => editor.chain().focus().toggleCode().run()}
//               >
//                 <Code className="w-3 h-3" />
//               </ToolbarButton>
//             </div>
//           </BubbleMenu>
//           <FloatingMenu editor={editor} tippyOptions={{ duration: 100 }}>
//             <div className="bg-base-100 shadow-lg rounded border border-base-300 p-1">
//               <ToolbarButton
//                 onClick={() =>
//                   editor.chain().focus().toggleHeading({ level: 2 }).run()
//                 }
//               >
//                 H2
//               </ToolbarButton>
//               <ToolbarButton
//                 onClick={() => editor.chain().focus().toggleBulletList().run()}
//               >
//                 列表
//               </ToolbarButton>
//             </div>
//           </FloatingMenu>
//         </div>

//         {mode === "markdown" && (
//           <textarea
//             value={markdownContent}
//             onChange={handleMarkdownChange}
//             placeholder={placeholder}
//             disabled={disabled}
//             className="w-full h-[300px] p-4 font-mono text-sm border-none focus:ring-0 bg-transparent resize-y outline-none"
//           />
//         )}
//       </div>
//     </div>
//   );
// };
