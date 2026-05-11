// components/Pagination.tsx
interface PaginationProps {
  totalPages: number;
  currentPage: number;
  onPageChange: (page: number) => void;
  totalItems?: number; // 可选，显示总记录数
  pageSize?: number; // 可选，配合 totalItems 显示"共 X 条"
}

export function Pagination({
  totalPages,
  currentPage,
  onPageChange,
  totalItems,
  pageSize,
}: PaginationProps) {
  // 生成页码列表（包含省略号）
  const getPageNumbers = (): (number | "ellipsis")[] => {
    const delta = 2; // 当前页前后显示的页码数量
    const range: (number | "ellipsis")[] = [];
    const left = currentPage - delta;
    const right = currentPage + delta;

    for (let i = 1; i <= totalPages; i++) {
      if (
        i === 1 || // 第一页
        i === totalPages || // 最后一页
        (i >= left && i <= right) // 当前页附近
      ) {
        range.push(i);
      } else if (range[range.length - 1] !== "ellipsis") {
        range.push("ellipsis");
      }
    }
    return range;
  };

  const pages = getPageNumbers();

  // 格式化总条数信息
  const getTotalText = () => {
    if (totalItems !== undefined && pageSize !== undefined) {
      const start = (currentPage - 1) * pageSize + 1;
      const end = Math.min(currentPage * pageSize, totalItems);
      return `第 ${start} - ${end} 条，共 ${totalItems} 条`;
    }
    return null;
  };

  const totalText = getTotalText();

  if (totalPages <= 1) return null; // 只有一页时不显示分页

  return (
    <div className="flex flex-col items-center gap-3 mt-6">
      {/* 分页按钮组 */}
      <div className="join">
        {/* 上一页 */}
        <button
          className="join-item btn btn-sm"
          disabled={currentPage === 1}
          onClick={() => onPageChange(currentPage - 1)}
          aria-label="上一页"
        >
          «
        </button>

        {/* 页码按钮（含省略号） */}
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
              aria-label={`第 ${page} 页`}
              aria-current={currentPage === page ? "page" : undefined}
            >
              {page}
            </button>
          ),
        )}

        {/* 下一页 */}
        <button
          className="join-item btn btn-sm"
          disabled={currentPage === totalPages}
          onClick={() => onPageChange(currentPage + 1)}
          aria-label="下一页"
        >
          »
        </button>
      </div>

      {/* 总条数信息（可选） */}
      {totalText && (
        <div className="text-xs text-base-content/60">{totalText}</div>
      )}
    </div>
  );
}
