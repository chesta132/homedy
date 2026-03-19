import { cn } from "@/lib/utils";
import Logo2k from "@/assets/images/icons/2k.svg?react";
import Logo16 from "@/assets/images/icons/16.svg?react";
import Logo32 from "@/assets/images/icons/32.svg?react";
import Logo48 from "@/assets/images/icons/48.svg?react";
import Logo180 from "@/assets/images/icons/180.svg?react";

const logos: Record<NonNullable<HomedyLogoProps["size"]>, React.FC<React.SVGProps<SVGSVGElement>>> = {
  "2k": Logo2k,
  "16": Logo16,
  "32": Logo32,
  "48": Logo48,
  "180": Logo180,
};

type HomedyLogoProps = {
  size?: "2k" | "16" | "32" | "48" | "180";
} & React.ComponentProps<"div">;

export const HomedyLogo = ({ size = "32", className, ...rest }: HomedyLogoProps) => {
  const LogoIcon = logos[size];
  return (
    <div {...rest} className={cn("flex items-center justify-center", className)}>
      <LogoIcon />
    </div>
  );
};
