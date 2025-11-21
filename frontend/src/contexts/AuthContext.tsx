"use client";

import React, { createContext, useContext, useEffect, useState } from "react";
import Cookies from "js-cookie";
<<<<<<< HEAD
=======
import { AuthUser } from "@/lib/types";
import { api } from "@/lib/api";
>>>>>>> b2b83c2 (Added add-attendee page, menus for profile and settings)

interface AuthContextType {
  isAuthed: boolean;
  isLoading: boolean;
<<<<<<< HEAD
=======
  user: AuthUser | null;
>>>>>>> b2b83c2 (Added add-attendee page, menus for profile and settings)
  login: (token: string) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [isAuthed, setIsAuthed] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
<<<<<<< HEAD

  useEffect(() => {
    const token = Cookies.get("token");
    setIsAuthed(Boolean(token));
    setIsLoading(false);
  }, []);

  const login = (token: string) => {
    Cookies.set("token", token, { expires: 7 });
    setIsAuthed(true);
=======
  const [user, setUser] = useState<AuthUser | null>(null);

  useEffect(() => {
    const token = Cookies.get("token");
    if (!token) {
      setIsLoading(false);
      return;
    }
    setIsAuthed(true);
    async () => {
      try {
        const { data } = await api.get("/auth/me");
      } catch {
        Cookies.remove("token");
        setIsAuthed(false);
        setUser(null);
      } finally {
        setIsLoading(false);
      }
    };
    setIsLoading(false);
  }, []);

  const login = async (token: string) => {
    Cookies.set("token", token, { expires: 7 });
    setIsAuthed(true);
    setIsLoading(false);
    try {
      const { data } = await api.get("/auth/me");
      setUser(data.user);
    } catch {
      Cookies.remove("token");
      setIsAuthed(false);
      setUser(null);
      throw new Error("Unable to fetch user profile after login");
    } finally {
      setIsLoading(false);
    }
>>>>>>> b2b83c2 (Added add-attendee page, menus for profile and settings)
  };

  const logout = () => {
    Cookies.remove("token");
    setIsAuthed(false);
<<<<<<< HEAD
  };

  return (
    <AuthContext.Provider value={{ isAuthed, isLoading, login, logout }}>
=======
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ isAuthed, isLoading, user, login, logout }}>
>>>>>>> b2b83c2 (Added add-attendee page, menus for profile and settings)
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
