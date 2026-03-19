/**
 * Mirror of backend converter.pairs.go ConvertPairs.
 * pdf -> xlsx is intentionally excluded per product decision.
 */
export const CONVERT_PAIRS: Record<string, string[]> = {
  html: ["md"],
  md:   ["html"],
  pdf:  ["docx", "pptx"],
  xlsx: ["pdf", "csv"],
  docx: ["pdf"],
  pptx: ["pdf"],
  csv:  ["xlsx"],
};

export const ACCEPTED_EXTS = Object.keys(CONVERT_PAIRS);

export const ACCEPTED_MIME = [
  "application/pdf",
  "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
  "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
  "application/vnd.openxmlformats-officedocument.presentationml.presentation",
  "text/html",
  "text/markdown",
  "text/csv",
  ".md",
  ".html",
  ".csv",
].join(",");

export function getExt(filename: string): string {
  return filename.split(".").pop()?.toLowerCase() ?? "";
}

export function getTargets(filename: string): string[] {
  const ext = getExt(filename);
  return CONVERT_PAIRS[ext] ?? [];
}
