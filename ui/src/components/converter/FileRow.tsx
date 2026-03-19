import { FileIcon, X, Loader2 } from "lucide-react";
import { motion } from "framer-motion";
import { ConvertSelect } from "./ConvertSelect";
import { getTargets } from "@/services/converter/converter.pairs";
import type { FileEntry } from "@/services/converter/converter.types";

interface FileRowProps {
  entry: FileEntry;
  onRemove: (id: string) => void;
  onConvertToChange: (id: string, val: string) => void;
}

/** Desktop row: [icon | filename] ←justify-between→ [loader | select | x] */
export function FileRow({ entry, onRemove, onConvertToChange }: FileRowProps) {
  const targets = getTargets(entry.file.name);

  return (
    <motion.div
      layout
      initial={{ opacity: 0, y: 8 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0, x: -12 }}
      transition={{ duration: 0.15 }}
      className="flex items-center justify-between gap-3 rounded-lg border border-[#1e1e1e] bg-[#111111] px-4 py-3"
    >
      <div className="flex min-w-0 items-center gap-3">
        <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-md border border-[#2a2a2a] bg-[#1a1a1a]">
          <FileIcon className="h-3.5 w-3.5 text-[#555555]" />
        </div>
        <span className="truncate text-sm text-[#ededed]">{entry.file.name}</span>
      </div>

      <div className="flex shrink-0 items-center gap-2">
        {entry.loading && <Loader2 className="h-3.5 w-3.5 animate-spin text-[#555555]" />}
        <ConvertSelect
          options={targets}
          value={entry.convertTo}
          onChange={(v) => onConvertToChange(entry.id, v)}
          disabled={entry.loading}
        />
        <button
          onClick={() => onRemove(entry.id)}
          disabled={entry.loading}
          className="flex h-6 w-6 items-center justify-center rounded text-[#444444] transition-colors hover:bg-[#1e1e1e] hover:text-[#888888] disabled:opacity-40"
        >
          <X className="h-3.5 w-3.5" />
        </button>
      </div>
    </motion.div>
  );
}

/** Mobile card: filename+x on top row, convert-to on second row */
export function FileMobileCard({ entry, onRemove, onConvertToChange }: FileRowProps) {
  const targets = getTargets(entry.file.name);

  return (
    <motion.div
      layout
      initial={{ opacity: 0, y: 8 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0, x: -12 }}
      transition={{ duration: 0.15 }}
      className="rounded-lg border border-[#1e1e1e] bg-[#111111] px-3 py-3"
    >
      <div className="flex items-center gap-2">
        <div className="flex h-7 w-7 shrink-0 items-center justify-center rounded-md border border-[#2a2a2a] bg-[#1a1a1a]">
          <FileIcon className="h-3 w-3 text-[#555555]" />
        </div>
        <span className="flex-1 truncate text-sm text-[#ededed]">{entry.file.name}</span>
        {entry.loading && <Loader2 className="h-3.5 w-3.5 animate-spin text-[#555555]" />}
        <button
          onClick={() => onRemove(entry.id)}
          disabled={entry.loading}
          className="flex h-6 w-6 items-center justify-center rounded text-[#444444] transition-colors hover:bg-[#1e1e1e] hover:text-[#888888] disabled:opacity-40"
        >
          <X className="h-3.5 w-3.5" />
        </button>
      </div>
      <div className="mt-2 flex items-center gap-2 pl-9">
        <span className="text-xs text-[#555555]">convert to</span>
        <ConvertSelect
          options={targets}
          value={entry.convertTo}
          onChange={(v) => onConvertToChange(entry.id, v)}
          disabled={entry.loading}
        />
      </div>
    </motion.div>
  );
}
