import { BrowserRouter, Routes, Route, Navigate } from "react-router";
import { AuthProvider } from "@/contexts/AuthContext";
import { AuthGuard } from "@/components/AuthGuard";
import { DashboardLayout } from "@/components/dashboard/DashboardLayout";
import { SignInPage } from "@/pages/auth/SignInPage";
import { SignUpPage } from "@/pages/auth/SignUpPage";
import { DashboardPage } from "@/pages/DashboardPage";
import { SMBPage } from "@/pages/SMBPage";
import { TerminalPage } from "@/pages/TerminalPage";
import { Toaster } from "@/components/ui/toaster";

export default function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          {/* Public auth routes */}
          <Route path="/signin" element={<SignInPage />} />
          <Route path="/signup" element={<SignUpPage />} />

          {/* Protected dashboard routes */}
          <Route element={<AuthGuard />}>
            <Route element={<DashboardLayout />}>
              <Route path="/dashboard" element={<DashboardPage />} />
              <Route path="/dashboard/smb" element={<SMBPage />} />
              <Route path="/dashboard/terminal" element={<TerminalPage />} />
            </Route>
          </Route>

          {/* Fallback */}
          <Route path="*" element={<Navigate to="/signin" replace />} />
        </Routes>
      </BrowserRouter>
      <Toaster />
    </AuthProvider>
  );
}
