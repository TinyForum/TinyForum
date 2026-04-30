import type { Metadata } from "next";
import { Inter, Fira_Code } from "next/font/google";
import { NextIntlClientProvider } from "next-intl";
import { notFound } from "next/navigation";
import { getTranslations, getMessages } from "next-intl/server";
import "../styles/globals.css";
import Providers from "@/layout/layout/Providers";
import AuthProvider from "@/layout/providers/AuthProvider";
import Navbar from "@/shared/ui/nav/Navbar";

const inter = Inter({ subsets: ["latin"], variable: "--font-inter" });
const firaCode = Fira_Code({
  subsets: ["latin"],
  variable: "--font-fira-code",
});

// 参数类型
interface LayoutParams {
  locale: string;
}

export async function generateMetadata({
  params,
}: {
  params: Promise<LayoutParams>;
}): Promise<Metadata> {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: "Site" });
  const brandName = t("brand") || "Your Brand Name";
  const description = t("description") || "Your brand description goes here.";

  return {
    title: {
      default: brandName,
      template: `%s | ${brandName}`,
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
  params: Promise<LayoutParams>;
}) {
  const { locale } = await params;

  let messages;
  try {
    messages = await getMessages({ locale });
  } catch {
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
