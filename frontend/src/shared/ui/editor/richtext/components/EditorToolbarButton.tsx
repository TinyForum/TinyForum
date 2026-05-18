import React from "react";

interface EditorToolbarButtonProps {
  onClick: () => void;
  active?: boolean;
  disabled?: boolean;
  children: React.ReactNode;
  title?: string;
}

export const EditorToolbarButton: React.FC<EditorToolbarButtonProps> = ({
  onClick,
  active,
  disabled,
  children,
  title,
}) => (
  <button
    type="button"
    onClick={onClick}
    disabled={disabled}
    title={title}
    className={`btn btn-ghost btn-xs ${active ? "btn-active text-primary" : ""}`}
  >
    {children}
  </button>
);
