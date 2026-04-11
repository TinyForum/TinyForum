'use client';

import { locales } from '@/i18n/request';
import { useLocale } from 'next-intl';
import { useRouter, usePathname } from 'next/navigation';

export default function LanguageSwitcher() {
  const locale = useLocale();
  const router = useRouter();
  const pathname = usePathname();

  const switchLanguage = (newLocale: string) => {
    // 替换路径中的语言部分
    const newPathname = pathname.replace(`/${locale}`, `/${newLocale}`);
    router.push(newPathname);
  };

  return (
    <div className="dropdown dropdown-end">
      <label tabIndex={0} className="btn btn-ghost btn-sm">
        <span className="mr-2">🌐</span>
        {locale}
      </label>
      <ul tabIndex={0} className="dropdown-content menu p-2 shadow bg-base-100 rounded-box w-52">
        {locales.map((loc) => (
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