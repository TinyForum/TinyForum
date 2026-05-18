// components/ModeSwitcher.tsx
import React from "react";
import { EditorMode } from "../types/richtext.type";

interface ModeSwitcherProps {
  mode: EditorMode;
  onSwitchToRich: () => void;
  onSwitchToMarkdown: () => void;
}

export const ModeSwitcher: React.FC<ModeSwitcherProps> = ({
  mode,
  onSwitchToRich,
  onSwitchToMarkdown,
}) => (
  <div className="flex gap-1 mr-2 border-r border-base-300 pr-2">
    <button
      onClick={onSwitchToRich}
      className={`btn btn-xs ${mode === "rich" ? "btn-primary" : "btn-ghost"}`}
    >
      富文本
    </button>
    <button
      onClick={onSwitchToMarkdown}
      className={`btn btn-xs ${mode === "markdown" ? "btn-primary" : "btn-ghost"}`}
    >
      Markdown
    </button>
  </div>
);
