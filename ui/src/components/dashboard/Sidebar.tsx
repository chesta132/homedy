import { useState } from "react";
import { Link, useLocation } from "react-router";
import { cn } from "@/lib/utils";
import { Badge } from "@/components/ui/badge";
import {
  LayoutDashboard,
  FolderOpen,
  Terminal,
  MessageSquare,
  Globe,
  ArrowRightLeft,
  FileText,
  CheckSquare,
  StickyNote,
  DollarSign,
  Lock,
  Menu,
  X,
} from "lucide-react";
import { AnimatePresence, motion } from "framer-motion";

type NavItem = {
  name: string;
  href: string;
  icon: React.ComponentType<{ className?: string }>;
  comingSoon?: boolean;
};

const NAV_ITEMS: NavItem[] = [
  { name: "Dashboard", href: "/dashboard", icon: LayoutDashboard },
  { name: "File Sharing", href: "/dashboard/smb", icon: FolderOpen },
  { name: "Terminal", href: "/dashboard/terminal", icon: Terminal }, // no comingSoon — feature available
  { name: "Chat", href: "#", icon: MessageSquare, comingSoon: true },
  { name: "DNS", href: "#", icon: Globe, comingSoon: true },
  { name: "Port Forward", href: "#", icon: ArrowRightLeft, comingSoon: true },
  { name: "PDF Converter", href: "#", icon: FileText, comingSoon: true },
  { name: "Todo", href: "#", icon: CheckSquare, comingSoon: true },
  { name: "Notes", href: "#", icon: StickyNote, comingSoon: true },
  { name: "Finance", href: "#", icon: DollarSign, comingSoon: true },
];

function NavLinks({ onClose }: { onClose?: () => void }) {
  const { pathname } = useLocation();

  return (
    <nav className="flex-1 space-y-0.5 overflow-y-auto p-2">
      {NAV_ITEMS.map(({ name, href, icon: Icon, comingSoon }) => {
        const isActive = pathname === href;
        return (
          <Link
            key={name}
            to={href}
            onClick={(e) => {
              if (comingSoon) e.preventDefault();
              else onClose?.();
            }}
            className={cn(
              "flex items-center gap-3 rounded-md px-3 py-2 text-sm transition-colors",
              isActive ? "bg-[#1c1c1c] text-white" : "text-[#888888] hover:bg-[#161616] hover:text-[#ededed]",
              comingSoon && "cursor-not-allowed opacity-50 hover:bg-transparent hover:text-[#888888]"
            )}
          >
            <Icon className="h-4 w-4 shrink-0" />
            <span className="flex-1 truncate">{name}</span>
            {comingSoon && (
              <Badge variant="outline" className="h-5 shrink-0 px-1.5 text-[10px]">
                <Lock className="mr-0.5 h-2.5 w-2.5" />
                Soon
              </Badge>
            )}
          </Link>
        );
      })}
    </nav>
  );
}

function SidebarLogo() {
  return (
    <div className="flex h-14 items-center border-b border-[#1e1e1e] px-4 shrink-0">
      <Link to="/dashboard" className="flex items-center gap-2.5">
        <div className="flex h-7 w-7 items-center justify-center rounded-md bg-white">
          <span className="text-sm font-bold text-black">H</span>
        </div>
        <span className="text-sm font-semibold text-[#ededed]">Homedy</span>
      </Link>
    </div>
  );
}

/** Desktop fixed sidebar */
export function Sidebar() {
  return (
    <aside className="fixed left-0 top-0 z-40 hidden h-screen w-56 flex-col border-r border-[#1e1e1e] bg-[#0d0d0d] lg:flex">
      <SidebarLogo />
      <NavLinks />
      <div className="border-t border-[#1e1e1e] p-3 shrink-0">
        <p className="text-xs text-[#333333]">Homedy v1.0.0</p>
      </div>
    </aside>
  );
}

/** Mobile slide-in sidebar triggered by hamburger */
export function MobileSidebar() {
  const [open, setOpen] = useState(false);

  return (
    <>
      <button
        onClick={() => setOpen(true)}
        className="lg:hidden text-[#888888] hover:text-[#ededed] transition-colors p-1"
        aria-label="Open navigation"
      >
        <Menu className="h-5 w-5" />
      </button>

      <AnimatePresence>
        {open && (
          <>
            {/* Backdrop */}
            <motion.div
              key="backdrop"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              onClick={() => setOpen(false)}
              className="fixed inset-0 z-40 bg-black/60 lg:hidden"
            />

            {/* Drawer */}
            <motion.aside
              key="drawer"
              initial={{ x: "-100%" }}
              animate={{ x: 0 }}
              exit={{ x: "-100%" }}
              transition={{ type: "tween", duration: 0.22 }}
              className="fixed left-0 top-0 z-50 flex h-screen w-56 flex-col border-r border-[#1e1e1e] bg-[#0d0d0d] lg:hidden"
            >
              <div className="flex h-14 items-center justify-between border-b border-[#1e1e1e] px-4 shrink-0">
                <Link to="/dashboard" onClick={() => setOpen(false)} className="flex items-center gap-2.5">
                  <div className="flex h-7 w-7 items-center justify-center rounded-md bg-white">
                    <span className="text-sm font-bold text-black">H</span>
                  </div>
                  <span className="text-sm font-semibold text-[#ededed]">Homedy</span>
                </Link>
                <button
                  onClick={() => setOpen(false)}
                  className="text-[#555555] hover:text-[#ededed] transition-colors"
                >
                  <X className="h-4 w-4" />
                </button>
              </div>
              <NavLinks onClose={() => setOpen(false)} />
              <div className="border-t border-[#1e1e1e] p-3 shrink-0">
                <p className="text-xs text-[#333333]">Homedy v1.0.0</p>
              </div>
            </motion.aside>
          </>
        )}
      </AnimatePresence>
    </>
  );
}
