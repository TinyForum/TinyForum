import { TabType } from "@/type/admin.types";
import { Search } from "lucide-react";

// 搜索栏组件
export function AdminSearchBar({
  tab,
  keyword,
  onKeywordChange,
  onPageReset,
  t
}: {
  tab: TabType;
  keyword: string;
  onKeywordChange: (keyword: string) => void;
  onPageReset: () => void;
  t: (key: string) => string;
}) {
  const handleChange = (value: string) => {
    onKeywordChange(value);
    onPageReset();
  };

  const placeholder = tab === "users"
    ? t("search_username_or_email")
    : t("search_post_title");

  return (
    <div className="flex gap-3 mb-4">
      <div className="relative flex-1 max-w-md">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-base-content/40" />
        <input
          type="text"
          placeholder={placeholder}
          value={keyword}
          onChange={(e) => handleChange(e.target.value)}
          className="input input-bordered input-sm w-full pl-9 focus:outline-none focus:border-primary"
        />
      </div>
    </div>
  );
}