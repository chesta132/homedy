import { createContext, useContext, useState, type ReactNode } from "react";
import type { User } from "@/types/models";
import type { StateErrorServer } from "@/types/server";

type AuthContextValue = {
  user: User | null;
  setUser: React.Dispatch<React.SetStateAction<User | null>>;
  globalError: StateErrorServer | null;
  setGlobalError: React.Dispatch<React.SetStateAction<StateErrorServer | null>>;
};

const AuthContext = createContext<AuthContextValue | null>(null);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<User | null>(null);
  const [globalError, setGlobalError] = useState<StateErrorServer | null>(null);

  return (
    <AuthContext value={{ user, setUser, globalError, setGlobalError }}>
      {children}
    </AuthContext>
  );
};

export const useAuth = () => {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used within AuthProvider");
  return ctx;
};
