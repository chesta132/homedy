import { useState, useEffect } from "react";
import { Plus, Loader2, Check, X } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { AppSecretModal } from "@/components/ui/app-secret-modal";
import api from "@/services/server/ApiClient";
import { toast } from "@/components/ui/toaster";
import type { ShareMap } from "@/types/models";
import { motion, AnimatePresence } from "framer-motion";
import { useAppSecret } from "@/hooks/useAppSecret";

type ConfigRow = { key: string; value: string };

/**
 * Editor for the [global] section of smb.conf.
 * Both GET and PUT /samba/config require app_secret (no caching).
 *
 * Desktop layout: | key | = | value | × |  (single row)
 * Mobile layout:  stacked card per entry   (matches mock-up style)
 */
export function SambaConfigEditor() {
  const [rows,     setRows]     = useState<ConfigRow[]>([]);
  const [original, setOriginal] = useState<ConfigRow[]>([]);
  const [loaded,   setLoaded]   = useState(false);
  const [loading,  setLoading]  = useState(false);
  const [saving,   setSaving]   = useState(false);

  const { prompting, getSecret, submitPrompt, cancelPrompt } = useAppSecret();

  // ── Load config on mount ──────────────────────────────────────────────────
  const fetchConfig = async () => {
    setLoading(true);
    try {
      const secret = await getSecret();
      if (!secret) return;

      const res = await api.sambaConfig.get<ShareMap>("/", {
        params: { app_secret: secret },
      });
      const parsed = Object.entries(res.data ?? {}).map(([key, value]) => ({ key, value }));
      setRows(parsed);
      setOriginal(parsed);
      setLoaded(true);
    } catch (err: any) {
      toast.error(err?.status === 403 ? "Invalid app secret" : "Failed to load global config");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { fetchConfig(); }, []);

  const addRow    = () => setRows((p) => [...p, { key: "", value: "" }]);
  const removeRow = (i: number) => setRows((p) => p.filter((_, idx) => idx !== i));
  const updateRow = (i: number, field: "key" | "value", val: string) =>
    setRows((p) => p.map((r, idx) => (idx === i ? { ...r, [field]: val } : r)));

  const handleCancel = () => setRows(original);

  const handleSave = async () => {
    if (rows.some((r) => r.key.trim() === "")) {
      toast.error("All keys must be non-empty. Remove blank rows first.");
      return;
    }
    const body: ShareMap = {};
    for (const { key, value } of rows) body[key.trim()] = value.trim();

    setSaving(true);
    try {
      const secret = await getSecret();
      if (!secret) return;

      await api.sambaConfig.put("/", body, { params: { app_secret: secret } });
      setOriginal(rows);
      toast.success("Global config saved");
    } catch (err: any) {
      toast.error(err?.status === 403 ? "Invalid app secret" : (err?.data?.message ?? "Failed to save config"));
    } finally {
      setSaving(false);
    }
  };

  const isDirty = JSON.stringify(rows) !== JSON.stringify(original);

  // ── Shared input className ────────────────────────────────────────────────
  const inputCls = "h-8 rounded-md border border-[#2a2a2a] bg-[#0d0d0d] px-2.5 font-mono text-xs text-[#ededed] placeholder:text-[#333333] focus:outline-none focus:border-[#3a3a3a] transition-colors w-full";

  return (
    <>
      <Card>
        <CardHeader className="flex-row items-center justify-between pb-3">
          <div>
            <CardTitle>Global Configuration</CardTitle>
            <CardDescription className="mt-0.5">
              Edit [global] section key-value pairs
            </CardDescription>
          </div>
        </CardHeader>

        <CardContent className="space-y-2 pb-4">
          {loading ? (
            <div className="flex items-center justify-center py-10">
              <Loader2 className="h-5 w-5 animate-spin text-[#444444]" />
            </div>
          ) : !loaded ? (
            <div className="flex flex-col items-center gap-3 py-8">
              <p className="text-sm text-[#444444]">
                Config requires app secret to load.
              </p>
              <Button size="sm" variant="outline" onClick={fetchConfig}>
                Load Config
              </Button>
            </div>
          ) : (
            <>
              {/* ── Rows ───────────────────────────────────────────────────── */}
              <AnimatePresence initial={false}>
                {rows.length === 0 && (
                  <p className="py-4 text-center text-sm text-[#444444]">
                    No options configured. Add one below.
                  </p>
                )}

                {rows.map((row, i) => (
                  <motion.div
                    key={i}
                    initial={{ opacity: 0, height: 0 }}
                    animate={{ opacity: 1, height: "auto" }}
                    exit={{ opacity: 0, height: 0 }}
                    transition={{ duration: 0.13 }}
                  >
                    {/* ── Desktop: single inline row ── */}
                    <div className="hidden sm:flex items-center gap-2">
                      <input
                        value={row.key}
                        onChange={(e) => updateRow(i, "key", e.target.value)}
                        placeholder="option"
                        className={inputCls}
                      />
                      <span className="text-[#333333] font-mono text-xs select-none shrink-0">=</span>
                      <input
                        value={row.value}
                        onChange={(e) => updateRow(i, "value", e.target.value)}
                        placeholder="value"
                        className={inputCls}
                      />
                      <button
                        type="button"
                        onClick={() => removeRow(i)}
                        className="h-8 w-8 shrink-0 flex items-center justify-center rounded-md text-[#444444] hover:text-red-400 hover:bg-red-950/30 transition-colors"
                      >
                        <X className="h-3.5 w-3.5" />
                      </button>
                    </div>

                    {/* ── Mobile: stacked card ── */}
                    <div className="sm:hidden rounded-lg border border-[#1e1e1e] bg-[#0d0d0d] p-3 space-y-2">
                      <div className="flex items-center justify-between gap-2">
                        <span className="text-[10px] uppercase tracking-wider text-[#444444]">Key</span>
                        <button
                          type="button"
                          onClick={() => removeRow(i)}
                          className="h-6 w-6 flex items-center justify-center rounded text-[#444444] hover:text-red-400 hover:bg-red-950/30 transition-colors"
                        >
                          <X className="h-3 w-3" />
                        </button>
                      </div>
                      <input
                        value={row.key}
                        onChange={(e) => updateRow(i, "key", e.target.value)}
                        placeholder="option name"
                        className={inputCls}
                      />
                      <span className="block text-[10px] uppercase tracking-wider text-[#444444]">Value</span>
                      <input
                        value={row.value}
                        onChange={(e) => updateRow(i, "value", e.target.value)}
                        placeholder="value"
                        className={inputCls}
                      />
                    </div>
                  </motion.div>
                ))}
              </AnimatePresence>

              {/* ── Add row ─────────────────────────────────────────────── */}
              <div className="pt-1">
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  onClick={addRow}
                  className="w-full border-dashed text-[#444444] hover:text-[#888888]"
                >
                  <Plus className="mr-1.5 h-3.5 w-3.5" />
                  Add
                </Button>
              </div>

              {/* ── Cancel / Save ────────────────────────────────────────── */}
              <div className="flex items-center justify-between gap-2 pt-3 border-t border-[#1e1e1e] mt-2">
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  onClick={handleCancel}
                  disabled={!isDirty || saving}
                  className="gap-1.5"
                >
                  <X className="h-3.5 w-3.5" />
                  Cancel
                </Button>
                <Button
                  type="button"
                  size="sm"
                  onClick={handleSave}
                  disabled={!isDirty || saving}
                  className="gap-1.5"
                >
                  {saving
                    ? <Loader2 className="h-3.5 w-3.5 animate-spin" />
                    : <Check className="h-3.5 w-3.5" />
                  }
                  Save
                </Button>
              </div>
            </>
          )}
        </CardContent>
      </Card>

      <AppSecretModal
        open={prompting}
        onSubmit={submitPrompt}
        onCancel={cancelPrompt}
      />
    </>
  );
}
