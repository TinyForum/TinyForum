"use client";

import { useTranslations } from "next-intl";
import Link from "next/link";

export function QuestionsMeClient() {
  const t = useTranslations("Questions");

  return (
    <div className="min-h-[60vh] flex items-center justify-center">
      <div className="text-center space-y-4 max-w-md">
        <h1 className="text-2xl font-bold">{t("my_questions")}</h1>
        <p className="text-base-content/60">{t("my_questions_coming_soon")}</p>
        <Link href="/questions" className="btn btn-primary">
          {t("browse_questions")}
        </Link>
      </div>
    </div>
  );
}
