import { useState } from "react";
import { Loader2, Trash2 } from "lucide-react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import api from "@/services/server/ApiClient";
import { toast } from "@/components/ui/toaster";

interface DeleteShareDialogProps {
  open: boolean;
  shareName: string;
  onClose: () => void;
  onSuccess: () => void;
}

/**
 * Confirmation dialog before deleting a share.
 */
export function DeleteShareDialog({ open, shareName, onClose, onSuccess }: DeleteShareDialogProps) {
  const [loading, setLoading] = useState(false);

  const handleDelete = async () => {
    setLoading(true);
    try {
      await api.samba.delete(`/${shareName}`);
      toast.success(`Share "${shareName}" deleted`);
      onSuccess();
      onClose();
    } catch (err: any) {
      toast.error(err?.data?.message ?? "Failed to delete share");
    } finally {
      setLoading(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={(v) => !v && onClose()}>
      <DialogContent className="sm:max-w-sm">
        <DialogHeader>
          <DialogTitle>Delete Share</DialogTitle>
          <DialogDescription>
            Are you sure you want to delete{" "}
            <span className="font-mono text-[#ededed]">"{shareName}"</span>? This action cannot be
            undone.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="outline" onClick={onClose} disabled={loading}>
            Cancel
          </Button>
          <Button
            onClick={handleDelete}
            disabled={loading}
            className="bg-red-950/80 text-red-400 hover:bg-red-900/80 border border-red-900/50"
          >
            {loading ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Deleting...
              </>
            ) : (
              <>
                <Trash2 className="mr-2 h-4 w-4" />
                Delete
              </>
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
