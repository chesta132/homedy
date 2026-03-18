import { Outlet } from "react-router";
import { Sidebar } from "@/components/dashboard/Sidebar";
import { Topbar } from "@/components/dashboard/Topbar";

/**
 * Wraps all /dashboard/* routes with the fixed sidebar and topbar.
 * The <Outlet /> renders the active child route.
 */
export function DashboardLayout() {
  return (
    <div className="min-h-screen bg-[#0a0a0a]">
      <Sidebar />
      <Topbar />
      <main className="pt-14 lg:ml-56">
        <div className="p-4 lg:p-6">
          <Outlet />
        </div>
      </main>
    </div>
  );
}
