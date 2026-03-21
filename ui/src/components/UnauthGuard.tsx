import { useEffect, useState } from "react";
import { Navigate, Outlet } from "react-router";
import { Loader2 } from "lucide-react";
import api from "@/services/server/ApiClient";
import { useAuth } from "@/contexts/AuthContext";
import type { User } from "@/types/models";

/**
 * Blocks auth pages (/signin, /signup) from logged-in users.
 *
 * State machine:
 *   "checking" → spinner shown, waiting for /me response
 *   "authed"   → user is logged in, redirect to /dashboard
 *   "guest"    → not logged in, render children (auth pages)
 */
export function UnauthGuard() {
  const { user, setUser } = useAuth();
  const [status, setStatus] = useState<"checking" | "authed" | "guest">(
    user ? "authed" : "checking"
  );

  useEffect(() => {
    if (user) {
      setStatus("authed");
      return;
    }

    const controller = new AbortController();

    api.auth
      .get<User>("/me", { signal: controller.signal })
      .then((res) => {
        setUser(res.data);
        setStatus("authed");
      })
      .catch((err) => {
        if (err?.name === "AbortError" || err?.code === "ERR_CANCELED") return;
        setStatus("guest");
      });

    return () => controller.abort();
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  if (status === "checking") {
    return (
      <div className="flex min-h-screen items-center justify-center bg-[#0a0a0a]">
        <Loader2 className="h-5 w-5 animate-spin text-[#333333]" />
      </div>
    );
  }

  if (status === "authed") return <Navigate to="/dashboard" replace />;

  return <Outlet />;
}
