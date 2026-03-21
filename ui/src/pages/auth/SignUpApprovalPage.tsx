import { useState } from "react";
import { Link } from "react-router";
import { Loader2, Mail, CheckCircle, XCircle, Clock } from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { HomedyLogo } from "@/components/ui/logo";
import api from "@/services/server/ApiClient";

type ApprovalStatus = "pending" | "approved" | "denied";

interface ApprovalStatusResponse {
  username?: string;
  email: string;
  status: ApprovalStatus;
}

export function SignUpApprovalPage() {
  const [email, setEmail] = useState("");
  const [emailError, setEmailError] = useState("");
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<ApprovalStatusResponse | null>(null);

  const handleCheck = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!email.trim()) {
      setEmailError("Email is required");
      return;
    }
    setEmailError("");
    setLoading(true);
    try {
      const res = await api.auth.get<ApprovalStatusResponse>("/signup/approval-status", {
        params: { email },
      });
      setResult(res.data);
    } catch {
      setEmailError("Failed to check status. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  const handleReset = () => {
    setResult(null);
    setEmail("");
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-[#0a0a0a] p-4">
      <AnimatePresence mode="wait">
        {result ? (
          <motion.div
            key="result"
            initial={{ opacity: 0, y: 16 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -16 }}
            transition={{ duration: 0.3 }}
            className="w-full max-w-sm text-center"
          >
            <div className="mb-6 flex justify-center">
              <HomedyLogo size="48" />
            </div>

            {result.status === "approved" && (
              <>
                <div className="mb-5 flex justify-center">
                  <div className="flex h-16 w-16 items-center justify-center rounded-full border border-emerald-900/40 bg-emerald-950/30">
                    <CheckCircle className="h-7 w-7 text-emerald-400" />
                  </div>
                </div>
                <h1 className="text-2xl font-semibold text-[#ededed]">Account approved</h1>
                <p className="mt-2 text-sm text-[#666666]">Your account is ready. You can sign in now.</p>
              </>
            )}

            {result.status === "pending" && (
              <>
                <div className="mb-5 flex justify-center">
                  <div className="flex h-16 w-16 items-center justify-center rounded-full border border-[#2a2a2a] bg-[#111111]">
                    <Clock className="h-7 w-7 text-[#888888]" />
                  </div>
                </div>
                <h1 className="text-2xl font-semibold text-[#ededed]">Still pending</h1>
                <p className="mt-2 text-sm text-[#666666]">Your request hasn't been reviewed yet. Check back later.</p>
              </>
            )}

            {result.status === "denied" && (
              <>
                <div className="mb-5 flex justify-center">
                  <div className="flex h-16 w-16 items-center justify-center rounded-full border border-red-900/40 bg-red-950/30">
                    <XCircle className="h-7 w-7 text-red-400" />
                  </div>
                </div>
                <h1 className="text-2xl font-semibold text-[#ededed]">Request denied</h1>
                <p className="mt-2 text-sm text-[#666666]">
                  Your account request was denied. If you think this was a mistake, contact the owner.
                </p>
              </>
            )}

            <div className="mt-4 rounded-lg border border-[#1e1e1e] bg-[#111111] px-4 py-3 text-left space-y-2">
              {result.username && (
                <div>
                  <p className="text-xs text-[#555555]">Username</p>
                  <p className="text-sm font-medium text-[#ededed]">{result.username}</p>
                </div>
              )}
              <div>
                <p className="text-xs text-[#555555]">Email</p>
                <p className="text-sm font-medium text-[#ededed]">{result.email}</p>
              </div>
            </div>

            <div className="mt-6 space-y-2">
              {result.status === "approved" && (
                <Link to="/signin">
                  <Button className="w-full h-10">Sign In</Button>
                </Link>
              )}
              <Button variant="outline" className="w-full h-10" onClick={handleReset}>
                Check another email
              </Button>
            </div>
          </motion.div>
        ) : (
          <motion.div
            key="form"
            initial={{ opacity: 0, y: 16 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -16 }}
            transition={{ duration: 0.3 }}
            className="w-full max-w-sm"
          >
            <div className="mb-8 text-center">
              <div className="mb-4 flex justify-center">
                <div className="flex h-14 w-14 items-center justify-center rounded-full border border-[#2a2a2a] bg-[#111111]">
                  <Mail className="h-6 w-6 text-[#888888]" />
                </div>
              </div>
              <h1 className="text-2xl font-semibold text-[#ededed]">Check request status</h1>
              <p className="mt-1.5 text-sm text-[#666666]">Enter the email you used to sign up</p>
            </div>

            <form onSubmit={handleCheck} className="space-y-4">
              <div className="space-y-1.5">
                <Label htmlFor="email">Email</Label>
                <Input
                  id="email"
                  type="email"
                  placeholder="admin@homedy.local"
                  value={email}
                  onChange={(e) => {
                    setEmail(e.target.value);
                    if (emailError) setEmailError("");
                  }}
                  error={emailError}
                  autoComplete="email"
                />
              </div>

              <Button type="submit" disabled={loading} className="w-full h-10">
                {loading ? (
                  <>
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    Checking...
                  </>
                ) : (
                  "Check Status"
                )}
              </Button>
            </form>

            <p className="mt-6 text-center text-sm text-[#555555]">
              Already approved?{" "}
              <Link to="/signin" className="text-[#ededed] underline underline-offset-4 hover:text-white">
                Sign in
              </Link>
            </p>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
