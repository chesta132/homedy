import { FolderOpen, Terminal, Lock } from "lucide-react";
import { Link } from "react-router";
import { motion } from "framer-motion";
import { cn } from "@/lib/utils";

const QUICK_LINKS = [
  {
    name: "File Sharing",
    desc: "Manage SMB/CIFS network shares",
    href: "/dashboard/smb",
    icon: FolderOpen,
    available: true,
  },
  {
    name: "Terminal",
    desc: "Access system terminal over WebSocket",
    href: "/dashboard/terminal",
    icon: Terminal,
    available: true,
  },
];

export function DashboardPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-xl font-semibold text-[#ededed]">Dashboard</h1>
        <p className="mt-1 text-sm text-[#555555]">Welcome to Homedy</p>
      </div>

      {/* Quick access cards */}
      <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
        {QUICK_LINKS.map(({ name, desc, href, icon: Icon, available }, i) => (
          <motion.div
            key={name}
            initial={{ opacity: 0, y: 12 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: i * 0.06 }}
          >
            <Link
              to={available ? href : "#"}
              onClick={(e) => !available && e.preventDefault()}
              className={cn(
                "group flex items-start gap-4 rounded-lg border border-[#1e1e1e] bg-[#0f0f0f] p-4 transition-colors",
                available
                  ? "hover:border-[#2a2a2a] hover:bg-[#141414] cursor-pointer"
                  : "opacity-50 cursor-not-allowed"
              )}
            >
              <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-md border border-[#2a2a2a] bg-[#1a1a1a]">
                <Icon className="h-4 w-4 text-[#888888]" />
              </div>
              <div className="min-w-0">
                <div className="flex items-center gap-2">
                  <p className="text-sm font-medium text-[#ededed]">{name}</p>
                  {!available && <Lock className="h-3 w-3 text-[#444444]" />}
                </div>
                <p className="mt-0.5 text-xs text-[#555555]">{desc}</p>
              </div>
            </Link>
          </motion.div>
        ))}
      </div>

      {/* Coming soon notice */}
      <div className="rounded-lg border border-dashed border-[#1e1e1e] p-6 text-center">
        <p className="text-sm text-[#444444]">More features coming soon</p>
      </div>
    </div>
  );
}
