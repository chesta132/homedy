import * as React from "react";
import { cn } from "@/lib/utils";

export interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  error?: string;
}

const Input = React.forwardRef<HTMLInputElement, InputProps>(
  ({ className, type, error, ...props }, ref) => {
    return (
      <div className="w-full">
        <input
          type={type}
          className={cn(
            "flex h-9 w-full rounded-md border border-[#2a2a2a] bg-[#1a1a1a] px-3 py-1 text-sm text-[#ededed] placeholder:text-[#555555] transition-colors",
            "focus:outline-none focus:border-[#3a3a3a]",
            "disabled:cursor-not-allowed disabled:opacity-50",
            error && "border-red-500/60 focus:border-red-500/60",
            className
          )}
          ref={ref}
          {...props}
        />
        {error && <p className="mt-1 text-xs text-red-400">{error}</p>}
      </div>
    );
  }
);
Input.displayName = "Input";

export { Input };
