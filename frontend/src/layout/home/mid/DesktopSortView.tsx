import { SortBy } from "@/shared/type/posts.types";

export interface DesktopSortViewProps {
  sortOptions: { value: SortBy; label: string; icon: React.ReactNode }[];
  sortBy: SortBy;
  onSortChange: (value: SortBy) => void;
}

export const DesktopSortView = ({
  sortOptions,
  sortBy,
  onSortChange,
}: DesktopSortViewProps) => (
  <div className="join">
    {sortOptions.map((option) => (
      <button
        key={option.value}
        className={`btn btn-sm join-item transition-all duration-200 ${
          sortBy === option.value
            ? "btn-primary shadow-md"
            : "btn-ghost hover:bg-base-200"
        }`}
        onClick={() => onSortChange(option.value)}
      >
        {option.icon}
        <span className="hidden sm:inline">{option.label}</span>
      </button>
    ))}
  </div>
);
