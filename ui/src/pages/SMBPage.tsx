import { useState, useEffect, useCallback } from "react";
import { Plus, Download, Upload, Loader2, FolderOpen, Settings2 } from "lucide-react";
import { motion } from "framer-motion";
import { Button } from "@/components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { AppSecretModal } from "@/components/ui/app-secret-modal";
import { SharesTable, SharesCardList } from "@/components/smb/SharesTable";
import { ShareForm } from "@/components/smb/ShareForm";
import { DeleteShareDialog } from "@/components/smb/DeleteShareDialog";
import { SambaConfigEditor } from "@/components/smb/SambaConfigEditor";
import api from "@/services/server/ApiClient";
import { toast } from "@/components/ui/toaster";
import type { Share, Shares } from "@/types/models";
import { useAppSecret } from "@/hooks/useAppSecret";

type FormMode =
  | { type: "create" }
  | { type: "edit"; name: string; share: Share };

export function SMBPage() {
  const [shares, setShares] = useState<Shares>({});
  const [loadingShares, setLoadingShares] = useState(true);
  const [backingUp, setBackingUp] = useState(false);
  const [restoring, setRestoring] = useState(false);

  // Dialog state
  const [formOpen, setFormOpen] = useState(false);
  const [formMode, setFormMode] = useState<FormMode>({ type: "create" });
  const [deleteTarget, setDeleteTarget] = useState<string | null>(null);

  const { prompting, getSecret, submitPrompt, cancelPrompt } = useAppSecret();

  const fetchShares = useCallback(async () => {
    setLoadingShares(true);
    try {
      const res = await api.samba.get<Shares>("/");
      setShares(res.data ?? {});
    } catch {
      toast.error("Failed to load shares");
    } finally {
      setLoadingShares(false);
    }
  }, []);

  useEffect(() => { fetchShares(); }, [fetchShares]);

  const openCreate = () => {
    setFormMode({ type: "create" });
    setFormOpen(true);
  };

  const openEdit = (name: string, share: Share) => {
    setFormMode({ type: "edit", name, share });
    setFormOpen(true);
  };

  const handleBackup = async () => {
    setBackingUp(true);
    try {
      const secret = await getSecret();
      if (!secret) return;
      await api.samba.post("/backup", undefined, { params: { app_secret: secret } });
      toast.success("Backup completed successfully");
    } catch (err: any) {
      if (err?.status === 403) { toast.error("Invalid app secret"); }
      else toast.error("Backup failed");
    } finally {
      setBackingUp(false);
    }
  };

  const handleRestore = async () => {
    setRestoring(true);
    try {
      const secret = await getSecret();
      if (!secret) return;
      const res = await api.samba.post<Shares>("/restore", undefined, { params: { app_secret: secret } });
      setShares(res.data ?? {});
      toast.success("Restore completed successfully");
    } catch (err: any) {
      if (err?.status === 403) { toast.error("Invalid app secret"); }
      else toast.error("Restore failed");
    } finally {
      setRestoring(false);
    }
  };

  const shareCount = Object.keys(shares).length;

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      className="space-y-5"
    >
      {/* Page header */}
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-xl font-semibold text-[#ededed]">File Sharing</h1>
          <p className="mt-0.5 text-sm text-[#555555]">Manage SMB/CIFS network shares</p>
        </div>
        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={handleBackup}
            disabled={backingUp}
            className="flex-1 sm:flex-none"
          >
            {backingUp ? (
              <Loader2 className="mr-2 h-3.5 w-3.5 animate-spin" />
            ) : (
              <Download className="mr-2 h-3.5 w-3.5" />
            )}
            Backup
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={handleRestore}
            disabled={restoring}
            className="flex-1 sm:flex-none"
          >
            {restoring ? (
              <Loader2 className="mr-2 h-3.5 w-3.5 animate-spin" />
            ) : (
              <Upload className="mr-2 h-3.5 w-3.5" />
            )}
            Restore
          </Button>
        </div>
      </div>

      {/* Tabs */}
      <Tabs defaultValue="shares">
        <TabsList>
          <TabsTrigger value="shares" className="gap-2">
            <FolderOpen className="h-3.5 w-3.5" />
            Shares
          </TabsTrigger>
          <TabsTrigger value="config" className="gap-2">
            <Settings2 className="h-3.5 w-3.5" />
            Global Config
          </TabsTrigger>
        </TabsList>

        {/* ── Shares tab ── */}
        <TabsContent value="shares" className="mt-4">
          <Card>
            <CardHeader className="flex-row items-center justify-between pb-3">
              <div>
                <CardTitle>Network Shares</CardTitle>
                <CardDescription className="mt-0.5">
                  {shareCount} {shareCount === 1 ? "share" : "shares"} configured
                </CardDescription>
              </div>
              <Button size="sm" onClick={openCreate}>
                <Plus className="mr-1.5 h-3.5 w-3.5" />
                Create Share
              </Button>
            </CardHeader>

            <CardContent className="p-0">
              {loadingShares ? (
                <div className="flex items-center justify-center py-12">
                  <Loader2 className="h-5 w-5 animate-spin text-[#444444]" />
                </div>
              ) : (
                <>
                  {/* Desktop table */}
                  <div className="hidden md:block">
                    <SharesTable
                      shares={shares}
                      onEdit={openEdit}
                      onDelete={(name) => setDeleteTarget(name)}
                    />
                  </div>
                  {/* Mobile cards */}
                  <div className="md:hidden p-4">
                    <SharesCardList
                      shares={shares}
                      onEdit={openEdit}
                      onDelete={(name) => setDeleteTarget(name)}
                    />
                  </div>
                </>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        {/* ── Config tab ── */}
        <TabsContent value="config" className="mt-4">
          <SambaConfigEditor />
        </TabsContent>
      </Tabs>

      {/* Dialogs */}
      <ShareForm
        open={formOpen}
        mode={formMode}
        onClose={() => setFormOpen(false)}
        onSuccess={fetchShares}
      />

      <DeleteShareDialog
        open={deleteTarget !== null}
        shareName={deleteTarget ?? ""}
        onClose={() => setDeleteTarget(null)}
        onSuccess={fetchShares}
      />

      <AppSecretModal
        open={prompting}
        onSubmit={submitPrompt}
        onCancel={cancelPrompt}
      />
    </motion.div>
  );
}
