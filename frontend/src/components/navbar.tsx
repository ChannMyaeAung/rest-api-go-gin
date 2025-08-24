"use client";
import { useTheme } from "next-themes";
import { usePathname } from "next/navigation";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { Button } from "./ui/button";
import { LogIn, LogOut, Moon, Sun } from "lucide-react";
import { useEffect, useState } from "react";
import { useAuth } from "@/contexts/AuthContext";

export function Navbar() {
  const { theme, setTheme } = useTheme();
  const pathname = usePathname();
  const router = useRouter();
  const { isAuthed, isLoading, logout } = useAuth();
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  function handleLogout() {
    logout();
    router.push("/login");
  }

  return (
    <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 px-6">
      <div className="mx-auto flex h-14 max-w-screen-2xl items-center px-4 md:px-6 lg:px-8">
        {/* Desktop Navigation Section */}
        <div className="mr-4 hidden md:flex">
          <Link href="/" className="mr-6 flex items-center space-x-2">
            <span className="hidden font-bold sm:inline-block">Events App</span>
          </Link>

          <nav className="flex items-center gap-4 text-sm lg:gap-6">
            <Link
              href="/events"
              className={`transition-colors hover:text-foreground/80 ${
                pathname?.startsWith("/events") && pathname !== "/events/new"
                  ? "text-foreground"
                  : "text-foreground/60"
              }`}
            >
              Events
            </Link>

            {/* Create Event Link - Shows immediately when auth state changes */}
            {!isLoading && isAuthed && (
              <Link
                href="/events/new"
                className={`transition-colors hover:text-foreground/80 ${
                  pathname === "/events/new"
                    ? "text-foreground"
                    : "text-foreground/60"
                }`}
              >
                Create Event
              </Link>
            )}
          </nav>
        </div>

        {/* Mobile Navigation */}
        <div className="mr-2 md:hidden">
          <Link href="/" className="font-bold">
            Events
          </Link>
        </div>

        {/* Right Side Actions */}
        <div className="flex flex-1 items-center justify-between space-x-2 md:justify-end">
          <div className="w-full flex-1 md:w-auto md:flex-none">
            {/* Future search functionality */}
          </div>

          <nav className="flex items-center gap-2">
            {/* Theme Toggle */}
            <Button
              variant="ghost"
              size="icon"
              aria-label="Toggle theme"
              onClick={() => setTheme(theme === "dark" ? "light" : "dark")}
            >
              {mounted ? (
                theme === "dark" ? (
                  <Sun className="h-[1.2rem] w-[1.2rem]" />
                ) : (
                  <Moon className="h-[1.2rem] w-[1.2rem]" />
                )
              ) : (
                <Sun className="h-[1.2rem] w-[1.2rem]" />
              )}
              <span className="sr-only">Toggle theme</span>
            </Button>

            {/* Auth Buttons - Updates immediately when auth state changes */}
            {!isLoading ? (
              !isAuthed ? (
                <Button asChild size="sm">
                  <Link href="/login">
                    <LogIn className="mr-2 h-4 w-4" />
                    Login
                  </Link>
                </Button>
              ) : (
                <Button size="sm" variant="outline" onClick={handleLogout}>
                  <LogOut className="mr-2 h-4 w-4" />
                  Logout
                </Button>
              )
            ) : (
              <Button asChild size="sm" disabled>
                <Link href="/login">
                  <LogIn className="mr-2 h-4 w-4" />
                  Login
                </Link>
              </Button>
            )}
          </nav>
        </div>
      </div>
    </header>
  );
}
