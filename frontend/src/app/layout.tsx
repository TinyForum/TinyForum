import type { Metadata } from 'next';
import { Inter, Fira_Code } from 'next/font/google';
import './globals.css';
import Providers from '@/components/layout/Providers';
import Navbar from '@/components/layout/Navbar';

const inter = Inter({ subsets: ['latin'], variable: '--font-inter' });
const firaCode = Fira_Code({ subsets: ['latin'], variable: '--font-fira-code' });

export const metadata: Metadata = {
  title: 'BBS Forum - 技术交流社区',
  description: '一个现代化的技术交流社区',
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="zh-CN" suppressHydrationWarning>
      <body className={`${inter.variable} ${firaCode.variable} font-sans min-h-screen bg-base-200`}>
        <Providers>
          <Navbar />
          <main className="container mx-auto px-4 py-6 max-w-6xl">
            {children}
          </main>
        </Providers>
      </body>
    </html>
  );
}
