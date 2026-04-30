// 分页组件
export function Pagination({
  totalPages,
  currentPage,
  onPageChange,
}: {
  totalPages: number;
  currentPage: number;
  onPageChange: (page: number) => void;
}) {
  const getPageNumbers = (): number[] => {
    const pages: number[] = [];
    const maxVisible = 7;
    const start = Math.max(1, currentPage - Math.floor(maxVisible / 2));
    const end = Math.min(totalPages, start + maxVisible - 1);

    for (let i = start; i <= end; i++) {
      pages.push(i);
    }
    return pages;
  };

  return (
    <div className="flex justify-center mt-6">
      <div className="join">
        <button
          className="join-item btn btn-sm"
          disabled={currentPage === 1}
          onClick={() => onPageChange(currentPage - 1)}
          aria-label="Previous page"
        >
          «
        </button>
        {getPageNumbers().map((p: number) => (
          <button
            key={p}
            className={`join-item btn btn-sm ${currentPage === p ? "btn-active btn-primary" : ""}`}
            onClick={() => onPageChange(p)}
            aria-label={`Go to page ${p}`}
            aria-current={currentPage === p ? "page" : undefined}
          >
            {p}
          </button>
        ))}
        <button
          className="join-item btn btn-sm"
          disabled={currentPage === totalPages}
          onClick={() => onPageChange(currentPage + 1)}
          aria-label="Next page"
        >
          »
        </button>
      </div>
    </div>
  );
}
