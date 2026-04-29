"use client";

import { routing } from "@/i18n/routing";
import { useLocale } from "next-intl";
import { useRouter, usePathname } from "next/navigation";
import { useState } from "react";

type Locale = (typeof routing.locales)[number];

export default function LanguageSwitcher() {
  const locale = useLocale();
  const router = useRouter();
  const pathname = usePathname();
  const [isOpen, setIsOpen] = useState(false);

  const switchLanguage = (newLocale: string) => {
    const segments = pathname.split("/");

    if (
      segments.length > 1 &&
      routing.locales.includes(segments[1] as Locale)
    ) {
      segments[1] = newLocale;
    } else {
      segments.splice(1, 0, newLocale);
    }

    const newPathname = segments.join("/") || "/";
    router.push(newPathname);
    setIsOpen(false); // 选择后关闭下拉菜单
  };

  const getDisplayText = (loc: string) => {
    if (loc === "zh-CN") return "简体中文";
    if (loc === "en-US") return "English";
    return loc;
  };

  const getFlagEmoji = (loc: string) => {
    if (loc === "zh-CN") return "🇨🇳";
    if (loc === "en-US") return "🇺🇸";
    return "🌐";
  };

  const currentFlag = getFlagEmoji(locale);
  const currentText =
    locale === "zh-CN" ? "中文" : locale === "en-US" ? "EN" : locale;

  return (
    <div className="dropdown dropdown-end">
      {/* 桌面端触发器 */}
      <label
        tabIndex={0}
        className="btn btn-ghost btn-sm hidden sm:flex"
        onClick={() => setIsOpen(!isOpen)}
      >
        <span>{currentFlag}</span>
        <span className="ml-1">{currentText}</span>
      </label>

      {/* 移动端触发器 - 仅图标 */}
      <label
        tabIndex={0}
        className="btn btn-ghost btn-sm btn-square sm:hidden"
        onClick={() => setIsOpen(!isOpen)}
      >
        <span>{currentFlag}</span>
      </label>

      {/* 下拉菜单 */}
      <ul
        tabIndex={0}
        className={`
          dropdown-content menu p-2 shadow-lg bg-base-100 rounded-box 
          ${isOpen ? "block" : "hidden"}
          min-w-[140px] sm:min-w-[160px]
          z-50
        `}
      >
        {routing.locales.map((loc: Locale) => (
          <li key={loc}>
            <button
              onClick={() => switchLanguage(loc)}
              className={`
                flex items-center gap-3 px-3 py-2 rounded-lg w-full
                ${locale === loc ? "bg-primary/10 text-primary font-medium" : "hover:bg-base-200"}
                transition-colors duration-200
              `}
            >
              <span className="text-lg">{getFlagEmoji(loc)}</span>
              <span>{getDisplayText(loc)}</span>
              {locale === loc && (
                <span className="ml-auto">
                  <svg
                    className="w-4 h-4"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M5 13l4 4L19 7"
                    />
                  </svg>
                </span>
              )}
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
}
