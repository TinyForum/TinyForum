'use client';

import { routing } from '@/i18n/routing';
import { useLocale } from 'next-intl';
import { useRouter, usePathname } from 'next/navigation';

export default function LanguageSwitcher() {
  const locale = useLocale();
  const router = useRouter();
  const pathname = usePathname();

  const switchLanguage = (newLocale: string) => {
    // 替换路径中的语言部分
    // 注意：pathname 已经包含了当前的 locale
    const segments = pathname.split('/');
    // 第一个空字符串，第二个是 locale
    if (segments.length > 1 && routing.locales.includes(segments[1] as any)) {
      segments[1] = newLocale;
    } else {
      // 如果路径中没有 locale，插入 locale
      segments.splice(1, 0, newLocale);
    }
    
    const newPathname = segments.join('/') || '/';
    router.push(newPathname);
  };

  return (
    <div className="dropdown dropdown-end px-2">
      <label tabIndex={0} className="btn btn-ghost btn-sm">
        <span className="mr-2">🌐</span>
        {locale === 'zh-CN' ? '中文' : 'English'}
      </label>
      <ul tabIndex={0} className="dropdown-content menu p-2 shadow bg-base-100 rounded-box w-52">
        {/* 修复：locale 是字符串，不是数组，需要手动列出语言选项 */}
        {routing.locales.map((loc) => (
          <li key={loc}>
            <button onClick={() => switchLanguage(loc)}>
              {loc === 'zh-CN' && '🇨🇳 简体中文'}
              {loc === 'en-US' && '🇺🇸 English'}
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
}