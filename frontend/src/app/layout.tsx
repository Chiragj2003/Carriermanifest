import type { Metadata, Viewport } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import { AuthProvider } from "@/lib/auth-context";
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
    <html lang="en">
      <body className={inter.className}>
        <AuthProvider>
          <div className="min-h-screen flex flex-col">
            <Navbar />
            <main className="flex-1">{children}</main>
            <footer className="border-t py-6 text-center text-sm text-muted-foreground">
              <p>&copy; {new Date().getFullYear()} CareerManifest. Built for Indian students, by India.</p>
            </footer>
          </div>
        </AuthProvider>
      </body>
    </html>
  );
}
