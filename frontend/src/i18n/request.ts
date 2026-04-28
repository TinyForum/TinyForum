// src/i18n/request.ts
import { getRequestConfig } from "next-intl/server";
import { routing } from "./routing";

type Locale = typeof routing.locales[number];

export default getRequestConfig(async ({ requestLocale }) => {
  // 从请求中获取 locale
  let locale = await requestLocale;

  console.log("i18n request locale:", locale);

  // 确保 locale 有效
  if (!locale || !routing.locales.includes(locale as Locale)) {
    locale = routing.defaultLocale;
  }

  try {
    const messages = (await import(`../messages/${locale}.json`)).default;

    return {
      locale,
      messages,
      timeZone: "Asia/Shanghai" as const,
      formats: {
        dateTime: {
          short: {
            day: "numeric" as const,
            month: "short" as const,
            year: "numeric" as const,
          },
        },
      },
    };
  } catch (error: unknown) {
    console.error(`Failed to load messages for locale: ${locale}`, error);
    // 尝试加载默认语言
    const defaultMessages = (
      await import(`../messages/${routing.defaultLocale}.json`)
    ).default;
    return {
      locale: routing.defaultLocale,
      messages: defaultMessages,
      timeZone: "Asia/Shanghai" as const,
    };
  }
});