import { useState } from "react";
import { KeyRound } from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

interface AppSecretModalProps {
  open: boolean;
  onSubmit: (secret: string) => void;
  onCancel: () => void;
}

/**
 * Modal that asks for APP_SECRET before performing protected samba operations
 * (backup, restore, config read/write).
 */
export function AppSecretModal({ open, onSubmit, onCancel }: AppSecretModalProps) {
  const [value, setValue] = useState("");

  const handleSubmit = () => {
    if (!value.trim()) return;
    onSubmit(value.trim());
    setValue("");
  };

  return (
    <Dialog open={open} onOpenChange={(v) => !v && onCancel()}>
      <DialogContent className="sm:max-w-sm">
        <DialogHeader>
          <div className="flex items-center gap-3 mb-1">
            <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg border border-[#2a2a2a] bg-[#1a1a1a]">
              <KeyRound className="h-4 w-4 text-[#888888]" />
            </div>
            <div>
              <DialogTitle>App Secret Required</DialogTitle>
              <DialogDescription className="mt-0.5">
                This action requires the server's{" "}
                <code className="text-[#888888] text-xs">APP_SECRET</code>
              </DialogDescription>
            </div>
          </div>
        </DialogHeader>

        <div className="space-y-1.5">
          <Label htmlFor="app-secret">Secret</Label>
          <Input
            id="app-secret"
            type="password"
            autoFocus
            placeholder="Enter app secret"
            value={value}
            onChange={(e) => setValue(e.target.value)}
            onKeyDown={(e) => e.key === "Enter" && handleSubmit()}
          />
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={onCancel}>
            Cancel
          </Button>
          <Button onClick={handleSubmit} disabled={!value.trim()}>
            Confirm
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
