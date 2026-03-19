import { useState, useCallback } from "react";
import { Layers, FileOutput, Loader2 } from "lucide-react";
import { AnimatePresence, motion } from "framer-motion";
import { Button } from "@/components/ui/button";
import { FileRow, FileMobileCard } from "@/components/converter/FileRow";
import { UploadRow, UploadButton } from "@/components/converter/UploadRow";
import { runSingle, runBatch } from "@/services/converter/converter.service";
import { getTargets, ACCEPTED_EXTS } from "@/services/converter/converter.pairs";
import type { FileEntry } from "@/services/converter/converter.types";

export function ConverterPage() {
  const [entries, setEntries] = useState<FileEntry[]>([]);
  const [batchLoading, setBatchLoading] = useState(false);

  const addFiles = useCallback((files: File[]) => {
    setEntries((prev) => [
      ...prev,
      ...files.map((f): FileEntry => ({
        id: `${Date.now()}-${Math.random()}`,
        file: f,
        convertTo: getTargets(f.name)[0] ?? "",
        loading: false,
      })),
    ]);
  }, []);

  const removeEntry = useCallback((id: string) => {
    setEntries((prev) => prev.filter((e) => e.id !== id));
  }, []);

  const setConvertTo = useCallback((id: string, val: string) => {
    setEntries((prev) => prev.map((e) => (e.id === id ? { ...e, convertTo: val } : e)));
  }, []);

  const setEntryLoading = useCallback((id: string, loading: boolean) => {
    setEntries((prev) => prev.map((e) => (e.id === id ? { ...e, loading } : e)));
  }, []);

  const handleSingle = useCallback(async () => {
    await runSingle(entries, setEntryLoading);
  }, [entries, setEntryLoading]);

  const handleBatch = useCallback(async () => {
    setBatchLoading(true);
    try {
      await runBatch(entries);
    } finally {
      setBatchLoading(false);
    }
  }, [entries]);

  const anyLoading = batchLoading || entries.some((e) => e.loading);
  const hasEntries = entries.length > 0;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-xl font-semibold text-[#ededed]">Converter</h1>
        <p className="mt-1 text-sm text-[#555555]">Convert files between formats</p>
      </div>

      {/* Desktop layout */}
      <div className="hidden sm:block space-y-2">
        <AnimatePresence mode="popLayout">
          {entries.map((entry) => (
            <FileRow
              key={entry.id}
              entry={entry}
              onRemove={removeEntry}
              onConvertToChange={setConvertTo}
            />
          ))}
        </AnimatePresence>
        <UploadRow onFiles={addFiles} />
      </div>

      {/* Mobile layout */}
      <div className="flex sm:hidden flex-col gap-2">
        <AnimatePresence mode="popLayout">
          {entries.map((entry) => (
            <FileMobileCard
              key={entry.id}
              entry={entry}
              onRemove={removeEntry}
              onConvertToChange={setConvertTo}
            />
          ))}
        </AnimatePresence>
        <UploadButton onFiles={addFiles} />
      </div>

      {/* Action bar */}
      <div className="flex items-center justify-between gap-3">
        <Button
          variant="outline"
          className="gap-2"
          disabled={!hasEntries || anyLoading}
          onClick={handleBatch}
        >
          {batchLoading ? (
            <Loader2 className="h-4 w-4 animate-spin" />
          ) : (
            <Layers className="h-4 w-4" />
          )}
          Batch
        </Button>

        <Button
          variant="default"
          className="gap-2"
          disabled={!hasEntries || anyLoading}
          onClick={handleSingle}
        >
          {entries.some((e) => e.loading) ? (
            <Loader2 className="h-4 w-4 animate-spin" />
          ) : (
            <FileOutput className="h-4 w-4" />
          )}
          Single
        </Button>
      </div>

      {/* Empty state */}
      {!hasEntries && (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          className="rounded-lg border border-dashed border-[#1e1e1e] p-8 text-center"
        >
          <p className="text-sm text-[#333333]">
            Supported: {ACCEPTED_EXTS.map((e) => `.${e}`).join(", ")}
          </p>
        </motion.div>
      )}
    </div>
  );
}
