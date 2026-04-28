"use client";

import { routing } from "@/i18n/routing";
import { useLocale } from "next-intl";
import { useRouter, usePathname } from "next/navigation";

type Locale = (typeof routing.locales)[number];

export default function LanguageSwitcher() {
  const locale = useLocale();
  const router = useRouter();
  const pathname = usePathname();

  const switchLanguage = (newLocale: string) => {
    // 替换路径中的语言部分
    // 注意：pathname 已经包含了当前的 locale
    const segments = pathname.split("/");
    // 第一个空字符串，第二个是 locale
    if (
      segments.length > 1 &&
      routing.locales.includes(segments[1] as Locale)
    ) {
      segments[1] = newLocale;
    } else {
      // 如果路径中没有 locale，插入 locale
      segments.splice(1, 0, newLocale);
    }

    const newPathname = segments.join("/") || "/";
    router.push(newPathname);
  };

  // 获取显示文本
  const getDisplayText = (loc: string) => {
    if (loc === "zh-CN") return "🇨🇳 简体中文";
    if (loc === "en-US") return "🇺🇸 English";
    return loc;
  };

  // 获取当前语言显示
  const getCurrentDisplay = () => {
    if (locale === "zh-CN") return "中文";
    if (locale === "en-US") return "English";
    return locale;
  };

  return (
    <div className="dropdown dropdown-end px-2">
      <label tabIndex={0} className="btn btn-ghost btn-sm">
        <span className="mr-2">🌐</span>
        {getCurrentDisplay()}
      </label>
      <ul
        tabIndex={0}
        className="dropdown-content menu p-2 shadow bg-base-100 rounded-box w-52"
      >
        {routing.locales.map((loc: Locale) => (
          <li key={loc}>
            <button onClick={() => switchLanguage(loc)}>
              {getDisplayText(loc)}
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
}
