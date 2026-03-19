import { useRef, useState } from "react";
import { Upload } from "lucide-react";
import { cn } from "@/lib/utils";
import { toast } from "@/components/ui/toaster";
import { Button } from "@/components/ui/button";
import { ACCEPTED_EXTS, ACCEPTED_MIME, getExt } from "@/services/converter/converter.pairs";
import { ConvertSelect } from "./ConvertSelect";

interface UploadRowProps {
  onFiles: (files: File[]) => void;
}

function filterFiles(files: File[]): File[] {
  const valid = files.filter((f) => ACCEPTED_EXTS.includes(getExt(f.name)));
  if (valid.length < files.length) toast.error("Some files were skipped — unsupported format");
  return valid;
}

/** Desktop: drag-and-drop row with disabled convert select on the right */
export function UploadRow({ onFiles }: UploadRowProps) {
  const inputRef = useRef<HTMLInputElement>(null);
  const [dragging, setDragging] = useState(false);

  return (
    <div
      onDragOver={(e) => { e.preventDefault(); setDragging(true); }}
      onDragLeave={() => setDragging(false)}
      onDrop={(e) => {
        e.preventDefault();
        setDragging(false);
        const valid = filterFiles(Array.from(e.dataTransfer.files));
        if (valid.length) onFiles(valid);
      }}
      className={cn(
        "flex items-center justify-between gap-3 rounded-lg border border-dashed px-4 py-3 transition-colors cursor-pointer",
        dragging
          ? "border-[#3a3a3a] bg-[#161616]"
          : "border-[#1e1e1e] bg-[#0d0d0d] hover:border-[#2a2a2a] hover:bg-[#111111]"
      )}
      onClick={() => inputRef.current?.click()}
    >
      <input
        ref={inputRef}
        type="file"
        multiple
        accept={ACCEPTED_MIME}
        className="hidden"
        onChange={(e) => {
          if (e.target.files) {
            const valid = filterFiles(Array.from(e.target.files));
            if (valid.length) onFiles(valid);
          }
          e.target.value = "";
        }}
      />
      <div className="flex items-center gap-3">
        <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-md border border-[#2a2a2a] bg-[#1a1a1a]">
          <Upload className="h-3.5 w-3.5 text-[#555555]" />
        </div>
        <span className="text-sm text-[#444444]">Drop files or click to upload</span>
      </div>
      <ConvertSelect options={[]} value="" onChange={() => {}} disabled />
    </div>
  );
}

/** Mobile: full-width upload button */
export function UploadButton({ onFiles }: UploadRowProps) {
  const inputRef = useRef<HTMLInputElement>(null);

  return (
    <>
      <input
        ref={inputRef}
        type="file"
        multiple
        accept={ACCEPTED_MIME}
        className="hidden"
        onChange={(e) => {
          if (e.target.files) {
            const valid = filterFiles(Array.from(e.target.files));
            if (valid.length) onFiles(valid);
          }
          e.target.value = "";
        }}
      />
      <Button variant="outline" className="w-full gap-2" onClick={() => inputRef.current?.click()}>
        <Upload className="h-4 w-4" />
        Upload Files
      </Button>
    </>
  );
}
