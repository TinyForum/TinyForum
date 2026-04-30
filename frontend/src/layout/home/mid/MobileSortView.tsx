import { SortBy } from "@/shared/type/posts.types";
import {
  MenuButton,
  Transition,
  MenuItems,
  Menu,
  MenuItem,
} from "@headlessui/react";
import { ChevronDown } from "lucide-react";
import { Fragment } from "react";

export interface MobileSortViewProps {
  sortOptions: { value: SortBy; label: string; icon: React.ReactNode }[];
  sortBy: SortBy;
  onSortChange: (value: SortBy) => void;
}

export const MobileSortView = ({
  sortOptions,
  sortBy,
  onSortChange,
}: MobileSortViewProps) => {
  const currentOption = sortOptions.find((opt) => opt.value === sortBy);
  return (
    <div className="dropdown dropdown-end">
      <Menu>
        <MenuButton className="btn btn-sm btn-primary gap-1 whitespace-nowrap">
          {currentOption?.icon}
          <span>{currentOption?.label}</span>
          <ChevronDown className="w-3 h-3" />
        </MenuButton>
        <Transition
          as={Fragment}
          enter="transition ease-out duration-100"
          enterFrom="transform opacity-0 scale-95"
          enterTo="transform opacity-100 scale-100"
          leave="transition ease-in duration-75"
          leaveFrom="transform opacity-100 scale-100"
          leaveTo="transform opacity-0 scale-95"
        >
          <MenuItems className="dropdown-content z-50 mt-2 w-48 rounded-box bg-base-100 p-2 shadow-lg ring-1 ring-base-300 focus:outline-none">
            {sortOptions.map((option) => (
              <MenuItem key={option.value}>
                {({ active }: { active: boolean }) => (
                  <button
                    className={`flex w-full items-center gap-2 rounded-btn px-4 py-2 text-sm transition-colors
                      ${active ? "bg-base-200" : ""}
                      ${
                        sortBy === option.value
                          ? "text-primary font-medium"
                          : "text-base-content"
                      }
                    `}
                    onClick={() => onSortChange(option.value)}
                  >
                    {option.icon}
                    {option.label}
                    {sortBy === option.value && (
                      <span className="ml-auto text-primary">✓</span>
                    )}
                  </button>
                )}
              </MenuItem>
            ))}
          </MenuItems>
        </Transition>
      </Menu>
    </div>
  );
};
