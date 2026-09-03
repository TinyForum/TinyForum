// components/Pagination.tsx
"use client";

import { useTranslations } from "next-intl";

interface PaginationProps {
  totalPages: number;
  currentPage: number;
  onPageChange: (page: number) => void;
  totalItems?: number;
  pageSize?: number;
}

export function Pagination({
  totalPages,
  currentPage,
  onPageChange,
  totalItems,
  pageSize,
}: PaginationProps) {
  const t = useTranslations("Common");

  const getPageNumbers = (): (number | "ellipsis")[] => {
    const delta = 2;
    const range: (number | "ellipsis")[] = [];
    const left = currentPage - delta;
    const right = currentPage + delta;

    for (let i = 1; i <= totalPages; i++) {
      if (
        i === 1 ||
        i === totalPages ||
        (i >= left && i <= right)
      ) {
        range.push(i);
      } else if (range[range.length - 1] !== "ellipsis") {
        range.push("ellipsis");
      }
    }
    return range;
  };

  const pages = getPageNumbers();

  const getTotalText = () => {
    if (totalItems !== undefined && pageSize !== undefined) {
      const start = (currentPage - 1) * pageSize + 1;
      const end = Math.min(currentPage * pageSize, totalItems);
      return t("pagination_info", { start, end, total: totalItems });
    }
    return null;
  };

  const totalText = getTotalText();

  if (totalPages <= 1) return null;

  return (
    <div className="flex flex-col items-center gap-3 mt-6">
      <div className="join">
        <button
          className="join-item btn btn-sm"
          disabled={currentPage === 1}
          onClick={() => onPageChange(currentPage - 1)}
          aria-label={t("prev_page")}
        >
          «
        </button>

        {pages.map((page, idx) =>
          page === "ellipsis" ? (
            <button
              key={`ellipsis-${idx}`}
              className="join-item btn btn-sm btn-disabled"
              aria-hidden="true"
            >
              ...
            </button>
          ) : (
            <button
              key={page}
              className={`join-item btn btn-sm ${
                currentPage === page ? "btn-active btn-primary" : ""
              }`}
              onClick={() => onPageChange(page)}
              aria-label={t("page_label", { page })}
              aria-current={currentPage === page ? "page" : undefined}
            >
              {page}
            </button>
          ),
        )}

        <button
          className="join-item btn btn-sm"
          disabled={currentPage === totalPages}
          onClick={() => onPageChange(currentPage + 1)}
          aria-label={t("next_page")}
        >
          »
        </button>
      </div>

      {totalText && (
        <div className="text-xs text-base-content/60">{totalText}</div>
      )}
    </div>
  );
}
