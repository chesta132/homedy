import { useState } from "react";
import { Link, useNavigate } from "react-router";
import { Eye, EyeOff, Loader2 } from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import useForm from "@/hooks/useForm";
import { handleFormError } from "@/services/models/handleError";
import api from "@/services/server/ApiClient";
import { useAuth } from "@/contexts/AuthContext";
import { HomedyLogo } from "@/components/ui/logo";
import { Mail } from "lucide-react";

export function SignUpPage() {
  const navigate = useNavigate();
  const { setGlobalError } = useAuth();
  const [showPassword, setShowPassword] = useState(false);
  const [loading, setLoading] = useState(false);
  const [submitted, setSubmitted] = useState(false);
  const [submittedEmail, setSubmittedEmail] = useState("");

  const {
    form: [form, setForm],
    error: [errors, setErrors],
    validate,
  } = useForm({ username: "", email: "", password: "" }, { username: true, email: true, password: true });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!validate.validateForm()) return;
    setLoading(true);
    try {
      await api.auth.post("/signup", {
        username: form.username,
        email: form.email,
        password: form.password,
      });
      setSubmittedEmail(form.email);
      setSubmitted(true);
    } catch (err) {
      handleFormError(err, setErrors as any, setGlobalError);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-[#0a0a0a] p-4">
      <AnimatePresence mode="wait">
        {submitted ? (
          <motion.div
            key="pending"
            initial={{ opacity: 0, y: 16 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -16 }}
            transition={{ duration: 0.3 }}
            className="w-full max-w-sm text-center"
          >
            <div className="mb-6 flex justify-center">
              <div className="flex h-16 w-16 items-center justify-center rounded-full border border-[#2a2a2a] bg-[#111111]">
                <Mail className="h-7 w-7 text-[#888888]" />
              </div>
            </div>
            <h1 className="text-2xl font-semibold text-[#ededed]">Check your inbox</h1>
            <p className="mt-2 text-sm text-[#666666]">
              Your request has been submitted. The owner will receive an approval email shortly.
            </p>
            <div className="mt-4 rounded-lg border border-[#1e1e1e] bg-[#111111] px-4 py-3">
              <p className="text-xs text-[#555555]">Registered as</p>
              <p className="mt-0.5 text-sm font-medium text-[#ededed]">{submittedEmail}</p>
            </div>
            <p className="mt-6 text-xs text-[#444444]">
              You'll receive an email at this address once the owner reviews your request.
            </p>
            <div className="mt-6 space-y-3">
              <Link to="/signup/approval">
                <Button className="w-full h-10">Check request status</Button>
              </Link>
              <p className="text-center text-sm text-[#555555]">
                Already approved?{" "}
                <Link to="/signin" className="text-[#ededed] underline underline-offset-4 hover:text-white">
                  Sign in
                </Link>
              </p>
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
            {/* Logo */}
            <div className="mb-8 text-center">
              <HomedyLogo size="48" />
              <h1 className="text-2xl font-semibold text-[#ededed]">Create an account</h1>
              <p className="mt-1.5 text-sm text-[#666666]">Get started with Homedy</p>
            </div>

            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="space-y-1.5">
                <Label htmlFor="username">Username</Label>
                <Input
                  id="username"
                  type="text"
                  placeholder="admin"
                  value={form.username ?? ""}
                  onChange={(e) => validate.validateField({ username: e.target.value })}
                  error={errors.username}
                  autoComplete="username"
                />
              </div>

              <div className="space-y-1.5">
                <Label htmlFor="email">Email</Label>
                <Input
                  id="email"
                  type="email"
                  placeholder="admin@homedy.local"
                  value={form.email ?? ""}
                  onChange={(e) => validate.validateField({ email: e.target.value })}
                  error={errors.email}
                  autoComplete="email"
                />
              </div>

              <div className="space-y-1.5">
                <Label htmlFor="password">Password</Label>
                <div className="relative">
                  <Input
                    id="password"
                    type={showPassword ? "text" : "password"}
                    placeholder="Create a secure password"
                    value={form.password ?? ""}
                    onChange={(e) => validate.validateField({ password: e.target.value })}
                    error={errors.password}
                    className="pr-10"
                    autoComplete="new-password"
                  />
                  <button
                    type="button"
                    onClick={() => setShowPassword((v) => !v)}
                    className="absolute right-3 top-[9px] text-[#555555] hover:text-[#888888] transition-colors"
                    tabIndex={-1}
                  >
                    {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                  </button>
                </div>
                <p className="text-xs text-[#444444]">8–32 chars, must include uppercase, lowercase, and a digit</p>
              </div>

              <Button type="submit" disabled={loading} className="w-full h-10">
                {loading ? (
                  <>
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    Creating account...
                  </>
                ) : (
                  "Sign Up"
                )}
              </Button>
            </form>

            <p className="mt-6 text-center text-sm text-[#555555]">
              Already have an account?{" "}
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
