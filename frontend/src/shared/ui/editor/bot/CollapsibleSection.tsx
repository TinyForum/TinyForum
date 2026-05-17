import { useState } from "react";

// ---------- 折叠面板组件 ----------
interface CollapsibleSectionProps {
  title: string;
  children: React.ReactNode;
  defaultOpen?: boolean;
}

export function CollapsibleSection({
  title,
  children,
  defaultOpen = true,
}: CollapsibleSectionProps) {
  const [isOpen, setIsOpen] = useState(defaultOpen);
  return (
    <div className="mb-4">
      <div
        className="flex items-center justify-between cursor-pointer py-1 select-none"
        onClick={() => setIsOpen(!isOpen)}
      >
        <h3 className="font-bold text-gray-800">{title}</h3>
        <span className="text-gray-500 text-sm">{isOpen ? "▼" : "▶"}</span>
      </div>
      {isOpen && <div className="mt-2 space-y-2">{children}</div>}
    </div>
  );
}
