import { useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";

// ---------- 可拖拽卡片组件（SortableItem）----------
interface SortableItemProps {
  id: string;
  typeLabel: string;
  onEdit: () => void;
  onDelete: () => void;
}

export function SortableItem({
  id,
  typeLabel,
  onEdit,
  onDelete,
}: SortableItemProps) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
  };

  return (
    <div
      ref={setNodeRef}
      style={style}
      className="bg-white border rounded p-3 flex items-center justify-between shadow-sm"
    >
      <div className="flex items-center gap-2 flex-1">
        <div
          {...attributes}
          {...listeners}
          className="cursor-grab text-gray-400 hover:text-gray-600"
        >
          ⋮⋮
        </div>
        <span className="font-medium">{typeLabel}</span>
      </div>
      <div className="flex gap-2">
        <button
          onClick={onEdit}
          className="text-blue-600 hover:text-blue-800 text-sm"
        >
          配置
        </button>
        <button
          onClick={onDelete}
          className="text-red-600 hover:text-red-800 text-sm"
        >
          删除
        </button>
      </div>
    </div>
  );
}
