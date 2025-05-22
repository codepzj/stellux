import { Inter as FontSans } from "next/font/google";
import "@/global.css";
import { Metadata, Viewport } from "next";
import { Link } from "@heroui/link";
import clsx from "clsx";

import { Providers } from "./providers";
import Navbar from "@/components/Navbar";
import { siteConfig } from "@/config/site";

// 字体配置
const fontSans = FontSans({
  subsets: ["latin"],
  variable: "--font-sans",
});

// 页面元数据
export const metadata: Metadata = {
  title: {
    default: siteConfig.name,
    template: `%s - ${siteConfig.name}`,
  },
  description: siteConfig.description,
  icons: {
    icon: "/favicon.ico",
  },
};

// 响应式视口主题色设置
export const viewport: Viewport = {
  themeColor: [
    { media: "(prefers-color-scheme: light)", color: "white" },
    { media: "(prefers-color-scheme: dark)", color: "black" },
  ],
};

// 布局组件
export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html suppressHydrationWarning lang="en" className={`${fontSans.variable}`}>
      <head />
      <body className={clsx("min-h-screen bg-background font-sans antialiased")}>
        <Providers themeProps={{ attribute: "class", defaultTheme: "dark" }}>
          <div className="relative flex flex-col min-h-screen">
            <Navbar />
            <main className="container mx-auto max-w-6xl pt-16 px-6 flex-1 mb-32">
              {children}
            </main>
            <footer className="w-full flex items-center justify-center py-3">
              <Link
                isExternal
                className="flex items-center gap-1 text-current"
                href="https://github.com/codepzj/stellux"
                title="基于开源项目stellux构建"
              >
                <span className="text-default-600">Powered by</span>
                <p className="text-primary">Stellux</p>
              </Link>
            </footer>
          </div>
        </Providers>
      </body>
    </html>
  );
}
