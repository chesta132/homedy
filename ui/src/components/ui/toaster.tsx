import { Toaster as SonnerToaster, toast } from "sonner";

export { toast };

export function Toaster() {
  return (
    <SonnerToaster
      position="bottom-right"
      toastOptions={{
        style: {
          background: "#1a1a1a",
          border: "1px solid #2a2a2a",
          color: "#ededed",
        },
      }}
    />
  );
}
