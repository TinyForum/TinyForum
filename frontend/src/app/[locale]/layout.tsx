import type { Metadata } from "next";
import { Inter, Fira_Code } from "next/font/google";
import Providers from "@/components/layout/Providers";
import Navbar from "@/components/layout/Navbar";
import { NextIntlClientProvider } from "next-intl";
import { notFound } from "next/navigation";
import { getTranslations } from "next-intl/server";
import "../styles/globals.css";
import AuthProvider from "@/components/providers/AuthProvider";

const inter = Inter({ subsets: ["latin"], variable: "--font-inter" });
const firaCode = Fira_Code({
  subsets: ["latin"],
  variable: "--font-fira-code",
});

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: "site" });
  const brandName = t("brand") || "Your Brand Name";
  const description = t("description") || "Your brand description goes here.";

  return {
    title: {
      default: brandName,
      template: `%s | ${brandName}`, // 方便子页面拼接标题
    },
    description: description,
    icons: [
      { url: "/favicon.ico", sizes: "any" },
      { url: "/assets/brand/logo.svg", type: "image/svg+xml" },
    ],
  };
}

export default async function RootLayout({
  children,
  params,
}: {
  children: React.ReactNode;
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;

  let messages;
  try {
    messages = (await import(`../../messages/${locale}.json`)).default;
  } catch (error) {
    notFound();
  }

  return (
    <html lang={locale} suppressHydrationWarning>
      <body
        suppressHydrationWarning
        className={`${inter.variable} ${firaCode.variable} font-sans h-screen overflow-hidden bg-base-200`}
      >
        <AuthProvider>
          <NextIntlClientProvider locale={locale} messages={messages}>
            <Providers>
              <div className="flex flex-col h-full">
                <Navbar />
                {/* 让 main 负责滚动 */}
                <main className="flex-1 overflow-y-auto custom-scrollbar">
                  <div className="container mx-auto max-w-7xl px-4 py-6">
                    {children}
                  </div>
                </main>
              </div>
            </Providers>
          </NextIntlClientProvider>
        </AuthProvider>
      </body>
    </html>
  );
}
