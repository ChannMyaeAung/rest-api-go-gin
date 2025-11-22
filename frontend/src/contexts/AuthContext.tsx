"use client";

import React, {
  createContext,
  useContext,
  useEffect,
  useState,
  useCallback,
} from "react";
import Cookies from "js-cookie";
import { AuthUser } from "@/lib/types";
import { api } from "@/lib/api";

interface AuthContextType {
  isAuthed: boolean;
  isLoading: boolean;
  user: AuthUser | null;
  login: (token: string) => Promise<void>;
  logout: () => void;
  refreshUser: () => Promise<void>;
}
const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [isAuthed, setIsAuthed] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [user, setUser] = useState<AuthUser | null>(null);

  const fetchUserProfile = useCallback(async () => {
    const { data } = await api.get<{ user: AuthUser }>("/auth/me");
    setUser(data.user);
    setIsAuthed(true);
  }, []);

  useEffect(() => {
    const token = Cookies.get("token");
    if (!token) {
      setIsLoading(false);
      return;
    }

    fetchUserProfile()
      .catch(() => {
        Cookies.remove("token");
        setUser(null);
        setIsAuthed(false);
      })
      .finally(() => {
        setIsLoading(false);
      });
  }, [fetchUserProfile]);

  const login = useCallback(
    async (token: string) => {
      Cookies.set("token", token, { expires: 7 });
      setIsLoading(true);

      try {
        await fetchUserProfile();
      } finally {
        setIsLoading(false);
      }
    },
    [fetchUserProfile]
  );

  const logout = useCallback(() => {
    Cookies.remove("token");
    setIsAuthed(false);
    setUser(null);
  }, []);

  const refreshUser = useCallback(async () => {
    try {
      await fetchUserProfile();
    } catch (error) {
      Cookies.remove("token");
      setUser(null);
      setIsAuthed(false);
      throw error;
    }
  }, [fetchUserProfile]);

  return (
    <AuthContext.Provider
      value={{ isAuthed, isLoading, user, login, logout, refreshUser }}
    >
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
