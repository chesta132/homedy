import { useNavigate } from "react-router";
import { LogOut, User } from "lucide-react";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Link } from "react-router";
import { MobileSidebar } from "./Sidebar";
import { useAuth } from "@/contexts/AuthContext";
import api from "@/services/server/ApiClient";
import { HomedyLogo } from "../ui/logo";

export function Topbar() {
  const { user } = useAuth();
  const navigate = useNavigate();

  const handleSignOut = async () => {
    try {
      await api.auth.post("/signout");
    } catch {
      // Best-effort — redirect regardless
    }
    navigate("/signin");
  };

  const initials = user?.username ? user.username.slice(0, 2).toUpperCase() : "??";

  return (
    <header className="fixed left-0 right-0 top-0 z-30 h-14 border-b border-[#1e1e1e] bg-[#0a0a0a]/80 backdrop-blur-sm lg:left-56">
      <div className="flex h-full items-center justify-between px-4 lg:px-6">
        {/* Left — mobile menu + logo */}
        <div className="flex items-center gap-3">
          <MobileSidebar />
          <Link to="/dashboard" className="flex items-center gap-2 lg:hidden">
            <HomedyLogo />
            <span className="text-sm font-semibold text-[#ededed]">Homedy</span>
          </Link>
        </div>

        {/* Right — user menu */}
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <button className="flex items-center gap-2 rounded-md px-2 py-1 text-[#888888] hover:bg-[#1a1a1a] hover:text-[#ededed] transition-colors focus:outline-none">
              <Avatar className="h-6 w-6">
                <AvatarFallback className="text-[10px]">{initials}</AvatarFallback>
              </Avatar>
              <span className="hidden text-sm sm:inline">{user?.username ?? "..."}</span>
            </button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-44">
            <DropdownMenuLabel>My Account</DropdownMenuLabel>
            <DropdownMenuSeparator />
            {/* Use onClick + navigate instead of asChild+Link to avoid event conflict */}
            {/* TODO: add /dashboard/profile page and api endpoint */}
            <DropdownMenuItem onClick={() => navigate("/dashboard/profile")}>
              <User className="mr-2 h-4 w-4" />
              Profile
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem onClick={handleSignOut} className="text-red-400 focus:bg-red-950/30 focus:text-red-400">
              <LogOut className="mr-2 h-4 w-4" />
              Sign Out
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </header>
  );
}
