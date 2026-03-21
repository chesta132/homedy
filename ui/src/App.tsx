import { BrowserRouter, Routes, Route, Navigate } from "react-router";
import { AuthProvider } from "@/contexts/AuthContext";
import { AuthGuard } from "@/components/AuthGuard";
import { UnauthGuard } from "@/components/UnauthGuard";
import { DashboardLayout } from "@/components/dashboard/DashboardLayout";
import { SignInPage } from "@/pages/auth/SignInPage";
import { SignUpPage } from "@/pages/auth/SignUpPage";
import { SignUpApprovalPage } from "@/pages/auth/SignUpApprovalPage";
import { SignUpReviewApprovalPage } from "@/pages/auth/SignUpReviewApprovalPage";
import { DashboardPage } from "@/pages/DashboardPage";
import { SMBPage } from "@/pages/SMBPage";
import { TerminalPage } from "@/pages/TerminalPage";
import { ConverterPage } from "@/pages/ConverterPage";
import { Toaster } from "@/components/ui/toaster";

export default function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          {/* Public auth routes — redirect to /dashboard if already logged in */}
          <Route element={<UnauthGuard />}>
            <Route path="/signin" element={<SignInPage />} />
            <Route path="/signup" element={<SignUpPage />} />
          </Route>

          {/* Approval pages — always public, no session check */}
          <Route path="/signup/approval" element={<SignUpApprovalPage />} />
          <Route path="/signup/review-approval" element={<SignUpReviewApprovalPage />} />

          {/* Protected dashboard routes */}
          <Route element={<AuthGuard />}>
            <Route element={<DashboardLayout />}>
              <Route path="/dashboard" element={<DashboardPage />} />
              <Route path="/dashboard/smb" element={<SMBPage />} />
              <Route path="/dashboard/terminal" element={<TerminalPage />} />
              <Route path="/dashboard/converter" element={<ConverterPage />} />
            </Route>
          </Route>

          {/* Fallback */}
          <Route path="*" element={<Navigate to="/dashboard" replace />} />
        </Routes>
      </BrowserRouter>
      <Toaster />
    </AuthProvider>
  );
}
