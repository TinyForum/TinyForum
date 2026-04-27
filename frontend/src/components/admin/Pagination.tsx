// 分页组件
export function Pagination({
  currentPage,
  total,
  pageSize = 20,
  onPageChange,
}: {
  currentPage: number;
  total: number;
  pageSize?: number;
  onPageChange: (page: number) => void;
}) {
  const totalPages = Math.ceil(total / pageSize);

  if (total <= pageSize) return null;

  return (
    <div className="flex justify-center p-4">
      <div className="join">
        <button
          className="join-item btn btn-sm"
          disabled={currentPage === 1}
          onClick={() => onPageChange(currentPage - 1)}
        >
          «
        </button>
        <button className="join-item btn btn-sm btn-active">
          {currentPage}
        </button>
        <button
          className="join-item btn btn-sm"
          disabled={currentPage >= totalPages}
          onClick={() => onPageChange(currentPage + 1)}
        >
          »
        </button>
      </div>
    </div>
  );
}
