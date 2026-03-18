import { useState } from "react";
import { Link, useNavigate } from "react-router";
import { Eye, EyeOff, Loader2 } from "lucide-react";
import { motion } from "framer-motion";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Checkbox } from "@/components/ui/checkbox";
import useForm from "@/hooks/useForm";
import { handleFormError } from "@/services/models/handleError";
import api from "@/services/server/ApiClient";
import { useAuth } from "@/contexts/AuthContext";
import type { User } from "@/types/models";

export function SignUpPage() {
  const navigate = useNavigate();
  const { setUser, setGlobalError } = useAuth();
  const [showPassword, setShowPassword] = useState(false);
  const [loading, setLoading] = useState(false);

  const { form: [form, setForm], error: [errors, setErrors], validate } = useForm(
    { username: "", email: "", password: "", rememberMe: false as boolean },
    { username: true, email: true, password: true }
  );

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!validate.validateForm()) return;
    setLoading(true);
    try {
      const res = await api.auth.post<User>("/signup", {
        username: form.username,
        email: form.email,
        password: form.password,
        remember_me: form.rememberMe,
      });
      setUser(res.data);
      navigate("/dashboard");
    } catch (err) {
      handleFormError(err, setErrors as any, setGlobalError);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-[#0a0a0a] p-4">
      <motion.div
        initial={{ opacity: 0, y: 16 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.3 }}
        className="w-full max-w-sm"
      >
        {/* Logo */}
        <div className="mb-8 text-center">
          <div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-xl bg-white">
            <span className="text-xl font-bold text-black">H</span>
          </div>
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
            <p className="text-xs text-[#444444]">
              8–32 chars, must include uppercase, lowercase, and a digit
            </p>
          </div>

          <div className="flex items-center gap-2">
            <Checkbox
              id="rememberMe"
              checked={form.rememberMe ?? false}
              onCheckedChange={(v) => setForm((prev) => ({ ...prev, rememberMe: !!v }))}
            />
            <Label htmlFor="rememberMe" className="text-sm font-normal cursor-pointer">
              Remember me
            </Label>
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
    </div>
  );
}
