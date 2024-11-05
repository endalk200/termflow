import { Inter } from "next/font/google";

import "@/styles/globals.css";

import { type Metadata } from "next";

import { TRPCReactProvider } from "@/trpc/react";
import { ThemeProvider } from "@/_components/molecules/theme-provider";

const inter = Inter({ subsets: ["latin"] });
export const metadata: Metadata = {
  title: "Termflow",
  description: "Termflow",
  icons: [{ rel: "icon", url: "/favicon.ico" }],
};

export default function RootLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="en" className={inter.className} suppressHydrationWarning>
      <body className={inter.className}>
        <TRPCReactProvider>
          <ThemeProvider
            attribute="class"
            defaultTheme="system"
            // enableSystem
            // disableTransitionOnChange
          >
            {children}
          </ThemeProvider>
        </TRPCReactProvider>
      </body>
    </html>
  );
}
