// hooks/useMarkdownMode.ts
import { useState, useCallback, useEffect } from "react";
import { Editor } from "@tiptap/react";
import TurndownService from "turndown";
import { marked } from "marked";

const defaultTurndown = new TurndownService({
  headingStyle: "atx",
  codeBlockStyle: "fenced",
  emDelimiter: "*",
});
defaultTurndown.addRule("fencedCodeBlock", {
  filter: "pre",
  replacement(content, node) {
    const codeEl = node.querySelector("code");
    const language = codeEl?.className?.match(/language-(\w+)/)?.[1] || "";
    return `\`\`\`${language}\n${codeEl?.textContent || content}\n\`\`\`\n\n`;
  },
});

marked.setOptions({ gfm: true, breaks: true });

interface UseMarkdownModeOptions {
  editor: Editor | null;
  htmlToMarkdown?: (html: string) => string;
  markdownParser?: (markdown: string) => Promise<string>;
  onChange?: (html: string) => void;
}

export const useMarkdownMode = ({
  editor,
  htmlToMarkdown = (html) => defaultTurndown.turndown(html),
  markdownParser = (markdown) => marked.parse(markdown) as Promise<string>,
  onChange,
}: UseMarkdownModeOptions) => {
  const [mode, setMode] = useState<"rich" | "markdown">("rich");
  const [markdownContent, setMarkdownContent] = useState("");

  // 切换到 Markdown 模式时，将当前 HTML 转为 Markdown
  const switchToMarkdown = useCallback(() => {
    if (!editor) return;
    const html = editor.getHTML();
    const md = htmlToMarkdown(html);
    setMarkdownContent(md);
    setMode("markdown");
  }, [editor, htmlToMarkdown]);

  // 切换到富文本模式时，将 Markdown 转为 HTML 并更新编辑器
  const switchToRich = useCallback(async () => {
    if (!editor) return;
    try {
      const html = await markdownParser(markdownContent);
      editor.commands.setContent(html);
      onChange?.(html);
      setMode("rich");
    } catch (error) {
      console.error("Markdown to HTML conversion failed:", error);
    }
  }, [editor, markdownContent, markdownParser, onChange]);

  const handleMarkdownChange = useCallback(
    async (newMarkdown: string) => {
      setMarkdownContent(newMarkdown);
      try {
        const html = await markdownParser(newMarkdown);
        onChange?.(html);
      } catch (error) {
        console.error("Markdown parse error:", error);
      }
    },
    [markdownParser, onChange],
  );

  // 当外部编辑器内容变化时，如果当前在 Markdown 模式，更新 markdownContent
  useEffect(() => {
    if (editor && mode === "markdown") {
      const html = editor.getHTML();
      const md = htmlToMarkdown(html);
      setMarkdownContent(md);
    }
  }, [editor, mode, htmlToMarkdown]);

  return {
    mode,
    markdownContent,
    switchToMarkdown,
    switchToRich,
    handleMarkdownChange,
  };
};
