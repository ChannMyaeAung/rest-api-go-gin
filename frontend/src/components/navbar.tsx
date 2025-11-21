"use client";
import { useTheme } from "next-themes";
import { usePathname } from "next/navigation";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { Button } from "./ui/button";
<<<<<<< HEAD
import { LogIn, LogOut, Moon, Sun } from "lucide-react";
import { useEffect, useState } from "react";
import { useAuth } from "@/contexts/AuthContext";
=======
import { LogIn, LogOut, Moon, Settings, Sun } from "lucide-react";
import { useEffect, useState } from "react";
import { useAuth } from "@/contexts/AuthContext";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuTrigger,
} from "./ui/dropdown-menu";
import { Avatar, AvatarFallback, AvatarImage } from "./ui/avatar";
import { is } from "zod/v4/locales";
>>>>>>> b2b83c2 (Added add-attendee page, menus for profile and settings)

export function Navbar() {
  const { theme, setTheme } = useTheme();
  const pathname = usePathname();
<<<<<<< HEAD
  const router = useRouter();
  const { isAuthed, isLoading, logout } = useAuth();
=======
  const { isAuthed, isLoading, logout, user } = useAuth();
>>>>>>> b2b83c2 (Added add-attendee page, menus for profile and settings)
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  function handleLogout() {
    logout();
<<<<<<< HEAD
    router.push("/login");
  }

  return (
    <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 px-6">
      <div className="mx-auto flex h-14 max-w-screen-2xl items-center px-4 md:px-6 lg:px-8">
=======
  }

  return (
    <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-backdrop-filter:bg-background/60 px-6">
      <div className="mx-auto flex h-14 max-w-screen-2xl items-center md:px-6 lg:px-8">
>>>>>>> b2b83c2 (Added add-attendee page, menus for profile and settings)
        {/* Desktop Navigation Section */}
        <div className="mr-4 hidden md:flex">
          <Link href="/" className="mr-6 flex items-center space-x-2">
            <span className="hidden font-bold sm:inline-block">Events App</span>
          </Link>

          <nav className="flex items-center gap-4 text-sm lg:gap-6">
            {isAuthed && (
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
            )}
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
<<<<<<< HEAD
                <Button size="sm" variant="outline" onClick={handleLogout}>
                  <LogOut className="mr-2 h-4 w-4" />
                  Logout
                </Button>
=======
                <DropdownMenu>
                  <DropdownMenuTrigger>
                    <Avatar>
                      <AvatarImage
                        src={
                          "https://www.htgtrading.co.uk/wp-content/uploads/2016/03/no-user-image-square.jpg"
                        }
                        alt="User Avatar"
                      />
                      <AvatarFallback>
                        <span className="sr-only">User Menu</span>
                        <AvatarImage
                          src={
                            "https://www.htgtrading.co.uk/wp-content/uploads/2016/03/no-user-image-square.jpg"
                          }
                          alt="User Avatar"
                        />
                      </AvatarFallback>
                    </Avatar>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent>
                    <DropdownMenuLabel className="font-normal">
                      <div className="flex flex-col">
                        <span className="text-sm font-semibold">
                          User #{user?.id ?? "-"}
                        </span>
                        <span className="text-xs text-muted-foreground">
                          {user?.email ?? "No email available"}
                        </span>
                      </div>
                    </DropdownMenuLabel>
                    <DropdownMenuGroup>
                      <DropdownMenuItem>
                        <Button size="sm" variant="ghost">
                          <Link
                            href="/settings"
                            className="flex items-center gap-1"
                          >
                            <Settings className="mr-2 h-4 w-4" />
                            Settings
                          </Link>
                        </Button>
                      </DropdownMenuItem>
                      <DropdownMenuItem>
                        <Button
                          size="sm"
                          variant="ghost"
                          onClick={handleLogout}
                        >
                          <Link
                            href={"/login"}
                            className="flex items-center gap-1"
                          >
                            <LogOut className="mr-2 h-4 w-4 text-destructive" />
                            <span className="text-destructive">Log out</span>
                          </Link>
                        </Button>
                      </DropdownMenuItem>
                    </DropdownMenuGroup>
                  </DropdownMenuContent>
                </DropdownMenu>
>>>>>>> b2b83c2 (Added add-attendee page, menus for profile and settings)
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
