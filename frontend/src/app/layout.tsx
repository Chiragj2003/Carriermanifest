/**
 * Root Layout — Wraps the entire app with providers and shared UI.
 *
 * Providers: ThemeProvider (dark/light), AuthProvider (JWT auth)
 * Shared UI: Navbar (top), Footer (bottom)
 */
import type { Metadata, Viewport } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import { AuthProvider } from "@/lib/auth-context";
import { ThemeProvider } from "@/lib/theme-context";
import { Navbar } from "@/components/navbar";

const inter = Inter({ subsets: ["latin"] });

export const viewport: Viewport = {
  width: "device-width",
  initialScale: 1,
  maximumScale: 1,
};

export const metadata: Metadata = {
  title: "CareerManifest - AI-Powered Career Decision Platform",
  description: "AI-Powered Career Decision Platform for Indian Students. Get personalized career recommendations based on your academic background, financial situation, and personality.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className={inter.className}>
        <ThemeProvider>
          <AuthProvider>
            <div className="min-h-screen flex flex-col">
              <Navbar />
              <main className="flex-1">{children}</main>
              <footer className="border-t py-8 text-center text-sm text-muted-foreground bg-card">
                <div className="container mx-auto px-4">
                  <p className="font-medium">&copy; {new Date().getFullYear()} CareerManifest</p>
                  <p className="mt-1 text-xs">Built for Indian students, by India 🇮🇳</p>
                </div>
              </footer>
            </div>
          </AuthProvider>
        </ThemeProvider>
      </body>
    </html>
  );
}
