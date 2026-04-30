import { useTranslations } from "next-intl";
import Link from "next/link";

// 空状态组件
export function EmptyPostList({
  isAuthenticated,
}: {
  isAuthenticated: boolean;
}) {
  const t = useTranslations("Post");

  return (
    <div className="text-center py-20 text-base-content/40">
      <p className="text-lg">{t("no_posts")}</p>
      {isAuthenticated && (
        <Link href="/posts/new" className="btn btn-primary mt-4">
          {t("post_your_first_post")}
        </Link>
      )}
    </div>
  );
}
