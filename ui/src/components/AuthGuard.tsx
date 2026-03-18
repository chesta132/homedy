import { useEffect, useState } from "react";
import { Navigate, Outlet } from "react-router";
import { Loader2 } from "lucide-react";
import api from "@/services/server/ApiClient";
import { useAuth } from "@/contexts/AuthContext";
import type { User } from "@/types/models";

/**
 * Protects dashboard routes by verifying the session via GET /auth/me.
 *
 * State machine:
 *   "checking" → spinner shown, waiting for /me response
 *   "ok"       → user loaded, render children
 *   "fail"     → /me failed (not abort), redirect to /signin
 */
export function AuthGuard() {
  const { user, setUser } = useAuth();
  // If user is already in context (e.g. just signed in), skip the check
  const [status, setStatus] = useState<"checking" | "ok" | "fail">(
    user ? "ok" : "checking"
  );

  useEffect(() => {
    // Already resolved — nothing to do
    if (user) {
      setStatus("ok");
      return;
    }

    const controller = new AbortController();

    api.auth
      .get<User>("/me", { signal: controller.signal })
      .then((res) => {
        setUser(res.data);
        setStatus("ok");
      })
      .catch((err) => {
        // AbortError fires when the component unmounts — not a real failure
        if (err?.name === "AbortError" || err?.code === "ERR_CANCELED") return;
        // Any real error (401, network, etc.) → send to signin
        setStatus("fail");
      });

    return () => controller.abort();
  // Run once on mount only
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  if (status === "checking") {
    return (
      <div className="flex min-h-screen items-center justify-center bg-[#0a0a0a]">
        <Loader2 className="h-5 w-5 animate-spin text-[#333333]" />
      </div>
    );
  }

  if (status === "fail") return <Navigate to="/signin" replace />;

  return <Outlet />;
}
