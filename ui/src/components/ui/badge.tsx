import * as React from "react";
import { cn } from "@/lib/utils";

export interface BadgeProps extends React.HTMLAttributes<HTMLDivElement> {
  variant?: "default" | "outline" | "success" | "destructive";
}

const Badge = ({ className, variant = "default", ...props }: BadgeProps) => {
  return (
    <div
      className={cn(
        "inline-flex items-center rounded-md border px-2 py-0.5 text-xs font-medium transition-colors",
        variant === "default" && "border-transparent bg-[#2a2a2a] text-[#ededed]",
        variant === "outline" && "border-[#2a2a2a] bg-[#111111] text-[#888888]",
        variant === "success" && "border-emerald-900/50 bg-emerald-950/50 text-emerald-400",
        variant === "destructive" && "border-red-900/50 bg-red-950/50 text-red-400",
        className
      )}
      {...props}
    />
  );
};

export { Badge };
