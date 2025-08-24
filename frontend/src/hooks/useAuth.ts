import { useEffect, useState } from "react";
import Cookies from "js-cookie";
export function useAuth() {
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

  return { isAuthed, isLoading, login, logout };
}
