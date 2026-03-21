import { useState } from "react";
import { useSearchParams } from "react-router";
import { motion } from "framer-motion";
import { CheckCircle, XCircle, Loader2, ShieldCheck } from "lucide-react";
import { HomedyLogo } from "@/components/ui/logo";
import { Button } from "@/components/ui/button";
import { AppSecretModal } from "@/components/ui/app-secret-modal";
import { useAppSecret } from "@/hooks/useAppSecret";
import api from "@/services/server/ApiClient";
import { toast } from "sonner";
import type { User } from "@/types/models";

type ReviewState = "pending" | "loading" | "approved" | "denied" | "error";

export function SignUpReviewApprovalPage() {
  const [params] = useSearchParams();
  const identifier = params.get("identifier") ?? "";
  const action = params.get("action");

  const { prompting, getSecret, submitPrompt, cancelPrompt } = useAppSecret();

  const [state, setState] = useState<ReviewState>("pending");
  const [user, setUser] = useState<User | null>(null);
  const [errorMessage, setErrorMessage] = useState("");

  const isValidParams = Boolean(identifier && (action === "approve" || action === "deny"));
  const isApproveAction = action === "approve";

  const handleSubmit = async () => {
    const secret = await getSecret();
    if (!secret) return;

    setState("loading");

    try {
      const res = await api.auth.patch<User>("/signup/approval", { identifier, action }, {
        headers: { "X-APP-SECRET": secret },
      });
      setUser(res.data);
      setState(isApproveAction ? "approved" : "denied");
    } catch (err: any) {
      const message: string =
        err?.getMessage?.() ?? err?.message ?? "Something went wrong";
      setErrorMessage(message);
      setState("error");
      toast.error(message);
    }
  };

  if (!isValidParams) {
    return <ReviewLayout><InvalidView /></ReviewLayout>;
  }

  if (state === "approved" && user) {
    return <ReviewLayout><ApprovedView user={user} /></ReviewLayout>;
  }
  if (state === "denied" && user) {
    return <ReviewLayout><DeniedView user={user} /></ReviewLayout>;
  }
  if (state === "error") {
    return (
      <ReviewLayout>
        <ErrorView message={errorMessage} onRetry={() => setState("pending")} />
      </ReviewLayout>
    );
  }

  return (
    <ReviewLayout>
      <AppSecretModal open={prompting} onSubmit={submitPrompt} onCancel={cancelPrompt} />

      <div className="mb-5 flex justify-center">
        <div
          className={`flex h-16 w-16 items-center justify-center rounded-full border ${
            isApproveAction
              ? "border-emerald-900/40 bg-emerald-950/30"
              : "border-red-900/40 bg-red-950/30"
          }`}
        >
          {state === "loading" ? (
            <Loader2 className="h-7 w-7 animate-spin text-[#888888]" />
          ) : isApproveAction ? (
            <ShieldCheck className="h-7 w-7 text-emerald-400" />
          ) : (
            <XCircle className="h-7 w-7 text-red-400" />
          )}
        </div>
      </div>

      <h1 className="text-2xl font-semibold text-[#ededed]">
        {isApproveAction ? "Approve account request?" : "Deny account request?"}
      </h1>
      <p className="mt-2 text-sm text-[#666666]">
        {isApproveAction
          ? "The user will be able to sign in after this is confirmed."
          : "The user's account request will be permanently rejected."}
      </p>

      <Button
        className="mt-6 w-full"
        variant={isApproveAction ? "default" : "destructive"}
        disabled={state === "loading"}
        onClick={handleSubmit}
      >
        {state === "loading" ? (
          <>
            <Loader2 className="mr-2 h-4 w-4 animate-spin" />
            Processing...
          </>
        ) : isApproveAction ? (
          "Approve"
        ) : (
          "Deny"
        )}
      </Button>
    </ReviewLayout>
  );
}

function ReviewLayout({ children }: { children: React.ReactNode }) {
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
        {children}
      </motion.div>
    </div>
  );
}

function InvalidView() {
  return (
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
  );
}

function ApprovedView({ user }: { user: User }) {
  return (
    <>
      <div className="mb-5 flex justify-center">
        <div className="flex h-16 w-16 items-center justify-center rounded-full border border-emerald-900/40 bg-emerald-950/30">
          <CheckCircle className="h-7 w-7 text-emerald-400" />
        </div>
      </div>
      <h1 className="text-2xl font-semibold text-[#ededed]">Account approved</h1>
      <p className="mt-2 text-sm text-[#666666]">
        <span className="font-medium text-[#ededed]">{user.username}</span>'s account has been
        approved. They can now sign in to Homedy.
      </p>
      <UserCard user={user} />
    </>
  );
}

function DeniedView({ user }: { user: User }) {
  return (
    <>
      <div className="mb-5 flex justify-center">
        <div className="flex h-16 w-16 items-center justify-center rounded-full border border-red-900/40 bg-red-950/30">
          <XCircle className="h-7 w-7 text-red-400" />
        </div>
      </div>
      <h1 className="text-2xl font-semibold text-[#ededed]">Request denied</h1>
      <p className="mt-2 text-sm text-[#666666]">
        <span className="font-medium text-[#ededed]">{user.username}</span>'s account request has
        been denied.
      </p>
      <UserCard user={user} />
    </>
  );
}

function ErrorView({ message, onRetry }: { message: string; onRetry: () => void }) {
  return (
    <>
      <div className="mb-5 flex justify-center">
        <div className="flex h-16 w-16 items-center justify-center rounded-full border border-red-900/40 bg-red-950/30">
          <XCircle className="h-7 w-7 text-red-400" />
        </div>
      </div>
      <h1 className="text-2xl font-semibold text-[#ededed]">Action failed</h1>
      <p className="mt-2 text-sm text-[#666666]">{message}</p>
      <Button variant="outline" className="mt-5 w-full" onClick={onRetry}>
        Try again
      </Button>
    </>
  );
}

function UserCard({ user }: { user: User }) {
  return (
    <div className="mt-4 rounded-lg border border-[#1e1e1e] bg-[#111111] px-4 py-3 text-left space-y-2">
      <div>
        <p className="text-xs text-[#555555]">Username</p>
        <p className="text-sm font-medium text-[#ededed]">{user.username}</p>
      </div>
      <div>
        <p className="text-xs text-[#555555]">Email</p>
        <p className="text-sm font-medium text-[#ededed]">{user.email}</p>
      </div>
    </div>
  );
}
