import Link from "next/link";
import { useTranslations } from "next-intl";

export default function NotFound() {
  const t = useTranslations("Home");

  return (
    <div className="min-h-screen flex items-center justify-center bg-base-100">
      <div className="text-center space-y-6">
        <h1 className="text-8xl font-bold text-primary">404</h1>
        <h2 className="text-2xl font-semibold text-base-content">
          {t("page_not_found")}
        </h2>
        <p className="text-base-content/60 max-w-md mx-auto">
          {t("page_not_found_desc")}
        </p>
        <div className="flex gap-4 justify-center">
          <Link href="/" className="btn btn-primary">
            {t("back_to_home")}
          </Link>
        </div>
      </div>
    </div>
  );
}
