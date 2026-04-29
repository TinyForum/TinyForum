"use client";

import { routing } from "@/i18n/routing";
import { useLocale } from "next-intl";
import { useRouter, usePathname } from "next/navigation";
import {
  Popover,
  PopoverButton,
  PopoverPanel,
  Transition,
} from "@headlessui/react";
import { Fragment } from "react";

type Locale = (typeof routing.locales)[number];

export default function LanguageSwitcher() {
  const locale = useLocale();
  const router = useRouter();
  const pathname = usePathname();

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
    <Popover className="relative inline-block">
      {({ close }: { close: () => void }) => (
        <>
          <PopoverButton
            className="btn btn-ghost btn-sm focus:outline-none"
            // 确保按钮可点击，移除可能干扰的样式
          >
            <span className="hidden sm:inline-flex items-center gap-1">
              <span>{currentFlag}</span>
              <span className="ml-1">{currentText}</span>
            </span>
            <span className="sm:hidden">{currentFlag}</span>
          </PopoverButton>

          <Transition
            as={Fragment}
            enter="transition ease-out duration-200"
            enterFrom="opacity-0 translate-y-1"
            enterTo="opacity-100 translate-y-0"
            leave="transition ease-in duration-150"
            leaveFrom="opacity-100 translate-y-0"
            leaveTo="opacity-0 translate-y-1"
          >
            <PopoverPanel
              portal
              anchor="bottom end"
              className="z-50 origin-top-right"
            >
              <div className="mt-2 w-48 rounded-lg shadow-lg ring-1 ring-black ring-opacity-5 bg-base-100">
                <div className="p-2">
                  {routing.locales.map((loc: Locale) => (
                    <button
                      key={loc}
                      onClick={() => {
                        switchLanguage(loc);
                        close();
                      }}
                      className={`
                        flex items-center gap-3 px-3 py-2 rounded-lg w-full text-left
                        ${
                          locale === loc
                            ? "bg-primary/10 text-primary font-medium"
                            : "hover:bg-base-200"
                        }
                        transition-colors duration-200
                      `}
                    >
                      <span className="text-lg">{getFlagEmoji(loc)}</span>
                      <span className="flex-1">{getDisplayText(loc)}</span>
                      {locale === loc && (
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
                      )}
                    </button>
                  ))}
                </div>
              </div>
            </PopoverPanel>
          </Transition>
        </>
      )}
    </Popover>
  );
}
