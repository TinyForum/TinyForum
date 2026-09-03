// components/ModeSwitcher.tsx
import React from "react";
import { useTranslations } from "next-intl";
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
}) => {
  const t = useTranslations("Editor");

  return (
    <div className="flex gap-1 mr-2 border-r border-base-300 pr-2">
      <button
        onClick={onSwitchToRich}
        className={`btn btn-xs ${mode === "rich" ? "btn-primary" : "btn-ghost"}`}
      >
        {t("rich_text")}
      </button>
      <button
        onClick={onSwitchToMarkdown}
        className={`btn btn-xs ${mode === "markdown" ? "btn-primary" : "btn-ghost"}`}
      >
        {t("markdown_mode")}
      </button>
    </div>
  );
};
