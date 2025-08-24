"use client";

import React, { createContext, useContext, useEffect, useState } from "react";
import Cookies from "js-cookie";

interface AuthContextType {
  isAuthed: boolean;
  isLoading: boolean;
  login: (token: string) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [isAuthed, setIsAuthed] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const token = Cookies.get("token");
    setIsAuthed(Boolean(token));
    setIsLoading(false);
  }, []);

  const login = (token: string) => {
    Cookies.set("token", token, { expires: 7 });
    setIsAuthed(true);
  };

  const logout = () => {
    Cookies.remove("token");
    setIsAuthed(false);
  };

  return (
    <AuthContext.Provider value={{ isAuthed, isLoading, login, logout }}>
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
