import { useSearchParams } from "react-router";
import { motion } from "framer-motion";
import { CheckCircle, XCircle } from "lucide-react";
import { HomedyLogo } from "@/components/ui/logo";

export function SignUpReviewApprovalPage() {
  const [params] = useSearchParams();
  const username = params.get("username") ?? "";
  const email = params.get("email") ?? "";
  const action = params.get("action");

  const isValid = Boolean(username && email && action);
  const isApproved = action === "approve";

  return (
    <div className="flex min-h-screen items-center justify-center bg-[#0a0a0a] p-4">
      <motion.div
        initial={{ opacity: 0, y: 16 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.3 }}
        className="w-full max-w-sm text-center"
      >
        <div className="mb-6 flex justify-center">
          <HomedyLogo size="48" />
        </div>

        {!isValid ? (
          <>
            <div className="mb-5 flex justify-center">
              <div className="flex h-16 w-16 items-center justify-center rounded-full border border-[#2a2a2a] bg-[#111111]">
                <XCircle className="h-7 w-7 text-red-400" />
              </div>
            </div>
            <h1 className="text-2xl font-semibold text-[#ededed]">Invalid approval link</h1>
            <p className="mt-2 text-sm text-[#666666]">
              This approval link is invalid or has already been used.
            </p>
          </>
        ) : isApproved ? (
          <>
            <div className="mb-5 flex justify-center">
              <div className="flex h-16 w-16 items-center justify-center rounded-full border border-emerald-900/40 bg-emerald-950/30">
                <CheckCircle className="h-7 w-7 text-emerald-400" />
              </div>
            </div>
            <h1 className="text-2xl font-semibold text-[#ededed]">Account approved</h1>
            <p className="mt-2 text-sm text-[#666666]">
              <span className="text-[#ededed] font-medium">{username}</span>'s account has been approved.
            </p>
            <div className="mt-4 rounded-lg border border-[#1e1e1e] bg-[#111111] px-4 py-3 text-left space-y-2">
              <div>
                <p className="text-xs text-[#555555]">Username</p>
                <p className="text-sm font-medium text-[#ededed]">{username}</p>
              </div>
              <div>
                <p className="text-xs text-[#555555]">Email</p>
                <p className="text-sm font-medium text-[#ededed]">{email}</p>
              </div>
            </div>
          </>
        ) : (
          <>
            <div className="mb-5 flex justify-center">
              <div className="flex h-16 w-16 items-center justify-center rounded-full border border-red-900/40 bg-red-950/30">
                <XCircle className="h-7 w-7 text-red-400" />
              </div>
            </div>
            <h1 className="text-2xl font-semibold text-[#ededed]">Request denied</h1>
            <p className="mt-2 text-sm text-[#666666]">
              <span className="text-[#ededed] font-medium">{username}</span>'s account request has been denied.
            </p>
            <div className="mt-4 rounded-lg border border-[#1e1e1e] bg-[#111111] px-4 py-3 text-left space-y-2">
              <div>
                <p className="text-xs text-[#555555]">Username</p>
                <p className="text-sm font-medium text-[#ededed]">{username}</p>
              </div>
              <div>
                <p className="text-xs text-[#555555]">Email</p>
                <p className="text-sm font-medium text-[#ededed]">{email}</p>
              </div>
            </div>
          </>
        )}
      </motion.div>
    </div>
  );
}

