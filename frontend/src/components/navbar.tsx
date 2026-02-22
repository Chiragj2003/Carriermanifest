"use client";

import { useState } from "react";
import Link from "next/link";
import { useAuth } from "@/lib/auth-context";
import { Button } from "@/components/ui/button";
import { useRouter } from "next/navigation";

export function Navbar() {
  const { user, logout, isAdmin } = useAuth();
  const router = useRouter();
  const [menuOpen, setMenuOpen] = useState(false);

  const handleLogout = () => {
    logout();
    router.push("/login");
  };

  return (
    <nav className="border-b bg-white/80 backdrop-blur-sm sticky top-0 z-50">
      <div className="container mx-auto px-4 h-16 flex items-center justify-between">
        <Link href={user ? "/dashboard" : "/"} className="flex items-center gap-2">
          <div className="h-8 w-8 rounded-lg bg-primary flex items-center justify-center">
            <span className="text-white font-bold text-sm">CM</span>
          </div>
          <span className="font-bold text-xl text-foreground">
            Career<span className="text-primary">Manifest</span>
          </span>
        </Link>

        {/* Desktop Nav */}
        <div className="hidden md:flex items-center gap-4">
          {user ? (
            <>
              <Link href="/dashboard" className="text-sm text-muted-foreground hover:text-foreground transition-colors">
                Dashboard
              </Link>
              <Link href="/assessment" className="text-sm text-muted-foreground hover:text-foreground transition-colors">
                New Assessment
              </Link>
              {isAdmin && (
                <Link href="/admin" className="text-sm text-muted-foreground hover:text-foreground transition-colors">
                  Admin
                </Link>
              )}
              <div className="flex items-center gap-3 ml-4">
                <span className="text-sm text-muted-foreground">{user.name}</span>
                <Button variant="outline" size="sm" onClick={handleLogout}>
                  Logout
                </Button>
              </div>
            </>
          ) : (
            <>
              <Link href="/login">
                <Button variant="ghost" size="sm">Login</Button>
              </Link>
              <Link href="/register">
                <Button size="sm">Get Started</Button>
              </Link>
            </>
          )}
        </div>

        {/* Mobile Hamburger */}
        <button
          className="md:hidden flex flex-col gap-1.5 p-2"
          onClick={() => setMenuOpen(!menuOpen)}
          aria-label="Toggle menu"
        >
          <span className={`block w-6 h-0.5 bg-foreground transition-all duration-300 ${menuOpen ? "rotate-45 translate-y-2" : ""}`} />
          <span className={`block w-6 h-0.5 bg-foreground transition-all duration-300 ${menuOpen ? "opacity-0" : ""}`} />
          <span className={`block w-6 h-0.5 bg-foreground transition-all duration-300 ${menuOpen ? "-rotate-45 -translate-y-2" : ""}`} />
        </button>
      </div>

      {/* Mobile Menu Dropdown */}
      {menuOpen && (
        <div className="md:hidden border-t bg-white px-4 py-4 space-y-3 animate-in slide-in-from-top-2">
          {user ? (
            <>
              <div className="text-sm font-medium text-foreground pb-2 border-b">
                Hi, {user.name}
              </div>
              <Link
                href="/dashboard"
                className="block text-sm text-muted-foreground hover:text-foreground py-2"
                onClick={() => setMenuOpen(false)}
              >
                📊 Dashboard
              </Link>
              <Link
                href="/assessment"
                className="block text-sm text-muted-foreground hover:text-foreground py-2"
                onClick={() => setMenuOpen(false)}
              >
                📝 New Assessment
              </Link>
              {isAdmin && (
                <Link
                  href="/admin"
                  className="block text-sm text-muted-foreground hover:text-foreground py-2"
                  onClick={() => setMenuOpen(false)}
                >
                  ⚙️ Admin
                </Link>
              )}
              <Button variant="outline" size="sm" className="w-full mt-2" onClick={handleLogout}>
                Logout
              </Button>
            </>
          ) : (
            <div className="flex flex-col gap-2">
              <Link href="/login" onClick={() => setMenuOpen(false)}>
                <Button variant="ghost" className="w-full">Login</Button>
              </Link>
              <Link href="/register" onClick={() => setMenuOpen(false)}>
                <Button className="w-full">Get Started</Button>
              </Link>
            </div>
          )}
        </div>
      )}
    </nav>
  );
}
