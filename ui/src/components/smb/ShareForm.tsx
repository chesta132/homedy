import { useState, useEffect } from "react";
import { Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { TagInput } from "@/components/ui/tag-input";
import { PermissionInput } from "@/components/ui/permission-input";
import { Separator } from "@/components/ui/separator";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from "@/components/ui/dialog";
import type { Share, SambaBool } from "@/types/models";
import api from "@/services/server/ApiClient";
import { toast } from "@/components/ui/toaster";

type ShareFormData = {
  name: string;
  path: string;
  read_only: SambaBool;
  browsable: SambaBool;
  valid_users: string[];
  admin_users: string[];
  permissions: number[];
};

const DEFAULT_FORM: ShareFormData = {
  name: "",
  path: "",
  read_only: "no",
  browsable: "yes",
  valid_users: [],
  admin_users: [],
  permissions: [7, 7, 5],
};

type Mode = { type: "create" } | { type: "edit"; name: string; share: Share };

interface ShareFormProps {
  open: boolean;
  mode: Mode;
  onClose: () => void;
  onSuccess: () => void;
}

/**
 * Modal form for creating or editing a Samba share.
 * Fields exactly match the backend Share model.
 */
export function ShareForm({ open, mode, onClose, onSuccess }: ShareFormProps) {
  const isEdit = mode.type === "edit";
  const [form, setForm] = useState<ShareFormData>(DEFAULT_FORM);
  const [errors, setErrors] = useState<Partial<Record<keyof ShareFormData, string>>>({});
  const [loading, setLoading] = useState(false);

  // Populate form when editing
  useEffect(() => {
    if (mode.type === "edit") {
      const { name, share } = mode;
      setForm({
        name,
        path: share.path,
        read_only: share.read_only,
        browsable: share.browsable,
        valid_users: share.valid_users ?? [],
        admin_users: share.admin_users ?? [],
        permissions: share.permissions ?? [7, 7, 5],
      });
    } else {
      setForm(DEFAULT_FORM);
    }
    setErrors({});
  }, [open, mode]);

  const set = <K extends keyof ShareFormData>(key: K, value: ShareFormData[K]) => {
    setForm((prev) => ({ ...prev, [key]: value }));
    // Clear field error on change
    if (errors[key]) setErrors((prev) => ({ ...prev, [key]: undefined }));
  };

  const validate = (): boolean => {
    const next: typeof errors = {};
    if (!form.name.trim()) next.name = "Share name is required";
    else if (!/^[a-zA-Z0-9_\-]+$/.test(form.name))
      next.name = "Only letters, digits, hyphens and underscores allowed";
    if (!form.path.trim()) next.path = "Path is required";
    else if (!form.path.startsWith("/")) next.path = "Path must be an absolute path (starts with /)";
    if (form.valid_users.length === 0) next.valid_users = "At least one valid user is required";
    if (form.admin_users.length === 0) next.admin_users = "At least one admin user is required";
    if (form.permissions.length !== 3 || form.permissions.some((p) => p < 0 || p > 7))
      next.permissions = "Each permission digit must be 0–7";
    setErrors(next);
    return Object.keys(next).length === 0;
  };

  const handleSubmit = async () => {
    if (!validate()) return;
    setLoading(true);
    try {
      const body = {
        name: form.name,
        path: form.path,
        read_only: form.read_only,
        browsable: form.browsable,
        valid_users: form.valid_users,
        admin_users: form.admin_users,
        permissions: form.permissions,
      };

      if (isEdit) {
        await api.samba.put(`/${mode.name}`, body);
        toast.success("Share updated successfully");
      } else {
        await api.samba.post("/", body);
        toast.success("Share created successfully");
      }
      onSuccess();
      onClose();
    } catch (err: any) {
      const msg = err?.data?.message ?? "Something went wrong";
      toast.error(msg);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={(v) => !v && onClose()}>
      <DialogContent className="sm:max-w-lg">
        <DialogHeader>
          <DialogTitle>{isEdit ? `Edit Share: ${mode.name}` : "Create New Share"}</DialogTitle>
          <DialogDescription>
            {isEdit ? "Update share configuration" : "Configure a new SMB/CIFS network share"}
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-5 py-2">
          {/* Share name — readonly when editing */}
          <div className="space-y-1.5">
            <Label>Share Name</Label>
            <Input
              placeholder="my-share"
              value={form.name}
              onChange={(e) => set("name", e.target.value)}
              error={errors.name}
              disabled={isEdit}
              className={isEdit ? "opacity-60 cursor-not-allowed" : ""}
            />
            {isEdit && (
              <p className="text-xs text-[#444444]">Share name cannot be changed after creation</p>
            )}
          </div>

          {/* Path */}
          <div className="space-y-1.5">
            <Label>Path</Label>
            <Input
              placeholder="/srv/my-share"
              value={form.path}
              onChange={(e) => set("path", e.target.value)}
              error={errors.path}
            />
          </div>

          <Separator />

          {/* Toggles */}
          <div className="space-y-3">
            <div className="flex items-center justify-between">
              <div>
                <Label>Read Only</Label>
                <p className="text-xs text-[#444444] mt-0.5">Prevent writes to this share</p>
              </div>
              <Switch
                checked={form.read_only === "yes"}
                onCheckedChange={(v) => set("read_only", v ? "yes" : "no")}
              />
            </div>
            <div className="flex items-center justify-between">
              <div>
                <Label>Browsable</Label>
                <p className="text-xs text-[#444444] mt-0.5">Show in network browse list</p>
              </div>
              <Switch
                checked={form.browsable === "yes"}
                onCheckedChange={(v) => set("browsable", v ? "yes" : "no")}
              />
            </div>
          </div>

          <Separator />

          {/* Valid Users */}
          <div className="space-y-1.5">
            <Label>Valid Users</Label>
            <TagInput
              value={form.valid_users}
              onChange={(tags) => set("valid_users", tags)}
              placeholder="Add user..."
              error={errors.valid_users}
            />
          </div>

          {/* Admin Users */}
          <div className="space-y-1.5">
            <Label>Admin Users</Label>
            <TagInput
              value={form.admin_users}
              onChange={(tags) => set("admin_users", tags)}
              placeholder="Add admin user..."
              error={errors.admin_users}
            />
          </div>

          <Separator />

          {/* Permissions */}
          <div className="space-y-1.5">
            <Label>Permissions</Label>
            <PermissionInput
              value={form.permissions}
              onChange={(perms) => set("permissions", perms)}
              error={errors.permissions}
            />
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={onClose} disabled={loading}>
            Cancel
          </Button>
          <Button onClick={handleSubmit} disabled={loading}>
            {loading ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                {isEdit ? "Saving..." : "Creating..."}
              </>
            ) : isEdit ? (
              "Save Changes"
            ) : (
              "Create Share"
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
