import type { Metadata } from 'next';
import { Inter, Fira_Code } from 'next/font/google';
import Providers from '@/components/layout/Providers';
import Navbar from '@/components/layout/Navbar';
import { NextIntlClientProvider } from 'next-intl';
import { getMessages } from 'next-intl/server';
import { notFound } from 'next/navigation';

import '../styles/globals.css';
const inter = Inter({ subsets: ['latin'], variable: '--font-inter' });
const firaCode = Fira_Code({ subsets: ['latin'], variable: '--font-fira-code' });

export const metadata: Metadata = {
  title: 'Mucly 论坛',
  description: '传统乐器交流社区',
};

export default async function RootLayout({
  children,
  params
}: {
  children: React.ReactNode;
  params: Promise<{ locale: string }> | { locale: string };
}) {
  // 在 Next.js 15 中，params 可能是 Promise
  const { locale } = await params;
  
  // 获取消息文件
  let messages;
  try {
    messages = (await import(`../../../messages/${locale}.json`)).default;
  } catch (error) {
    notFound();
  }

  return (
    <html lang={locale} suppressHydrationWarning>
      <body suppressHydrationWarning className={`${inter.variable} ${firaCode.variable} font-sans min-h-screen bg-base-200`}>
        <NextIntlClientProvider locale={locale} messages={messages}>
          <Providers>
            <Navbar />
            <main className="container mx-auto px-4 py-6 max-w-6xl">
              {children}
            </main>
          </Providers>
        </NextIntlClientProvider>
      </body>
    </html>
  );
}