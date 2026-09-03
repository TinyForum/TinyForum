// components/nav/SearchBar.tsx
import { Search } from "lucide-react";
import { useTranslations } from "next-intl";

interface SearchBarProps {
  keyword: string;
  onKeywordChange: (keyword: string) => void;
  placeholder?: string;
}

export default function SearchBar({
  keyword,
  onKeywordChange,
  placeholder,
}: SearchBarProps) {
  const t = useTranslations("Moderator");
  // TODO: 真实数据
  return (
    <div className="relative">
      <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-base-content/40" />
      <input
        type="text"
        value={keyword}
        onChange={(e) => onKeywordChange(e.target.value)}
        placeholder={placeholder || t("search")}
        className="input input-bordered w-full pl-10"
      />
    </div>
  );
}
