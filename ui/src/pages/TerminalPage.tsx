import { useEffect, useRef, useState, useCallback } from "react";
import { Terminal } from "@xterm/xterm";
import { FitAddon } from "@xterm/addon-fit";
import { motion, AnimatePresence } from "framer-motion";
import { Loader2, KeyRound } from "lucide-react";
import { cn } from "@/lib/utils";
import "@xterm/xterm/css/xterm.css";

// ── Resize helper — matches backend terminalMessage schema ───────────────────
// Backend parses: { type: "resize", resize: { x: cols, y: rows } }
function sendResize(ws: WebSocket, cols: number, rows: number) {
  if (ws.readyState !== WebSocket.OPEN) return;
  ws.send(JSON.stringify({ type: "resize", resize: { x: cols, y: rows } }));
}

// ── WS URL ───────────────────────────────────────────────────────────────────
// Derive from current page origin so it works in dev (proxied) and prod
function getWsUrl(secret: string): string {
  const proto = window.location.protocol === "https:" ? "wss:" : "ws:";
  const host = window.location.host;
  return `${proto}//${host}/api/ws/terminal?app_secret=${encodeURIComponent(secret)}`;
}

// ── xterm theme — follows app palette (#0a0a0a base, white accents) ──────────
const XTERM_THEME = {
  background:          "#0a0a0a",
  foreground:          "#e0e0e0",
  cursor:              "#ffffff",
  cursorAccent:        "#0a0a0a",
  selectionBackground: "rgba(255,255,255,0.15)",
  black:               "#1a1a1a",  brightBlack:   "#3a3a3a",
  red:                 "#f87171",  brightRed:     "#fca5a5",
  green:               "#86efac",  brightGreen:   "#bbf7d0",
  yellow:              "#fde047",  brightYellow:  "#fef08a",
  blue:                "#93c5fd",  brightBlue:    "#bfdbfe",
  magenta:             "#d8b4fe",  brightMagenta: "#e9d5ff",
  cyan:                "#67e8f9",  brightCyan:    "#a5f3fc",
  white:               "#e5e5e5",  brightWhite:   "#ffffff",
};

// ── Helper key model ─────────────────────────────────────────────────────────
// "modifier" keys (Ctrl, Shift, Alt) toggle on/off and combine with next key.
// "direct" keys send immediately.
type DirectKey = { type: "direct"; label: string; data: string; wide?: boolean };
type ModifierKey = { type: "modifier"; label: string; mod: "ctrl" | "shift" | "alt" };
type HelperKey = DirectKey | ModifierKey;

// Modifier → escape prefix table used when building combined sequences
// Ctrl+key:  send \x01-\x1a range
// Shift+key: send uppercase (terminal handles most cases natively, but helper sends it)
// Arrow / special key sequences with modifiers use ANSI CSI param ;modifier
const MODIFIER_PARAM: Record<string, number> = {
  shift: 2, alt: 3, "ctrl+shift": 6, ctrl: 5,
};

const HELPER_KEYS: HelperKey[] = [
  // ── Modifier toggles (Termux-style) ───────────────────────────────────────
  { type: "modifier", label: "CTRL", mod: "ctrl" },
  { type: "modifier", label: "ALT",  mod: "alt"  },
  // ── Arrow cluster ─────────────────────────────────────────────────────────
  { type: "direct", label: "↑", data: "\x1b[A" },
  { type: "direct", label: "↓", data: "\x1b[B" },
  { type: "direct", label: "←", data: "\x1b[D" },
  { type: "direct", label: "→", data: "\x1b[C" },
  // ── Common control shortcuts ──────────────────────────────────────────────
  { type: "direct", label: "Tab",  data: "\t" },
  { type: "direct", label: "Esc",  data: "\x1b" },
  { type: "direct", label: "↵",    data: "\r",   wide: false },
  // ── Navigation ────────────────────────────────────────────────────────────
  { type: "direct", label: "Home",  data: "\x1b[H" },
  { type: "direct", label: "End",   data: "\x1b[F" },
  { type: "direct", label: "PgUp",  data: "\x1b[5~" },
  { type: "direct", label: "PgDn",  data: "\x1b[6~" },
  { type: "direct", label: "Del",   data: "\x1b[3~" },
  { type: "direct", label: "Ins",   data: "\x1b[2~" },
];

// Ctrl+letter → \x01-\x1a
function ctrlCode(ch: string): string {
  const code = ch.toUpperCase().charCodeAt(0) - 64;
  if (code >= 1 && code <= 26) return String.fromCharCode(code);
  return ch;
}

type ConnStatus = "idle" | "connecting" | "connected" | "error";
type Mods = { ctrl: boolean; shift: boolean; alt: boolean };

// ── App-secret modal ─────────────────────────────────────────────────────────
function SecretModal({ onConnect }: { onConnect: (s: string) => void }) {
  const [secret, setSecret] = useState("");
  const submit = () => { if (secret.trim()) onConnect(secret.trim()); };
  return (
    <motion.div
      initial={{ opacity: 0, scale: 0.96 }}
      animate={{ opacity: 1, scale: 1 }}
      exit={{ opacity: 0, scale: 0.96 }}
      className="absolute inset-0 z-10 flex flex-col items-center justify-center bg-[#0a0a0a]/95 backdrop-blur-sm rounded-b-md"
    >
      <div className="w-full max-w-xs space-y-4 px-6">
        <div className="flex flex-col items-center gap-2">
          <div className="flex h-10 w-10 items-center justify-center rounded-lg border border-[#2a2a2a] bg-[#111111]">
            <KeyRound className="h-4 w-4 text-[#888888]" />
          </div>
          <p className="text-sm font-medium text-[#ededed]">App Secret Required</p>
          <p className="text-center text-xs text-[#555555]">
            Enter the <code className="text-[#888888]">APP_SECRET</code> configured on the server
          </p>
        </div>
        <input
          type="password"
          autoFocus
          value={secret}
          onChange={(e) => setSecret(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && submit()}
          placeholder="app secret"
          className={cn(
            "w-full rounded-md border border-[#2a2a2a] bg-[#1a1a1a] px-3 py-2",
            "font-mono text-sm text-[#ededed] placeholder:text-[#444444]",
            "focus:outline-none focus:border-[#3a3a3a] transition-colors"
          )}
        />
        <button
          onClick={submit}
          disabled={!secret.trim()}
          className="w-full rounded-md bg-white py-2 text-sm font-medium text-black transition hover:bg-[#e0e0e0] disabled:opacity-40 disabled:cursor-not-allowed"
        >
          Connect
        </button>
      </div>
    </motion.div>
  );
}

// ── Main component ────────────────────────────────────────────────────────────
export function TerminalPage() {
  const containerRef  = useRef<HTMLDivElement>(null);
  const termRef       = useRef<Terminal | null>(null);
  const fitRef        = useRef<FitAddon | null>(null);
  const wsRef         = useRef<WebSocket | null>(null);
  const inputHandlerRef = useRef<{ dispose(): void } | null>(null);

  const [status,      setStatus]      = useState<ConnStatus>("idle");
  const [showSecret,  setShowSecret]  = useState(false);
  const [mods,        setMods]        = useState<Mods>({ ctrl: false, shift: false, alt: false });
  // Ref mirrors mods state so sendKey always reads fresh values without needing
  // to re-create the callback (avoids stale closure bug with useCallback)
  const modsRef = useRef<Mods>({ ctrl: false, shift: false, alt: false });

  // ── Init xterm once ────────────────────────────────────────────────────────
  useEffect(() => {
    const term = new Terminal({
      fontFamily: '"Fira Code", "Cascadia Code", monospace',
      fontSize: 13,
      lineHeight: 1.5,
      cursorBlink: true,
      cursorStyle: "block",
      theme: XTERM_THEME,
      allowProposedApi: true,
    });

    const fit = new FitAddon();
    term.loadAddon(fit);
    term.open(containerRef.current!);
    requestAnimationFrame(() => fit.fit());

    const onResize = () => fit.fit();
    window.addEventListener("resize", onResize);

    termRef.current = term;
    fitRef.current  = fit;

    // When xterm dimensions change (window resize / fit), notify the PTY
    term.onResize(({ cols, rows }) => {
      if (wsRef.current?.readyState === WebSocket.OPEN) {
        sendResize(wsRef.current, cols, rows);
      }
    });

    return () => {
      window.removeEventListener("resize", onResize);
      term.dispose();
    };
  }, []);

  // ── Connect after secret entered ──────────────────────────────────────────
  const connect = useCallback((secret: string) => {
    if (wsRef.current) return;
    setStatus("connecting");
    setShowSecret(false);

    const url = getWsUrl(secret);
    const ws  = new WebSocket(url);
    ws.binaryType = "arraybuffer";
    wsRef.current = ws;

    ws.onopen = () => {
      setStatus("connected");

      // Dispose stale handler before registering new one (prevents double-send on reconnect)
      inputHandlerRef.current?.dispose();
      inputHandlerRef.current = termRef.current!.onData((data) => {
        if (ws.readyState === WebSocket.OPEN) ws.send(data);
      });

      fitRef.current?.fit();

      // Send initial PTY size after fit resolves
      const term = termRef.current!;
      sendResize(ws, term.cols, term.rows);

      termRef.current?.focus();
    };

    ws.onmessage = (e) => {
      if (e.data instanceof ArrayBuffer) {
        termRef.current?.write(new Uint8Array(e.data));
      } else {
        termRef.current?.write(e.data as string);
      }
    };

    ws.onclose = (e) => {
      // code 1003 / 1008 = server rejected (wrong secret or forbidden)
      if (e.code === 1003 || e.code === 1008 || e.code === 4003) {
        termRef.current?.writeln("\r\n\x1b[31mConnection rejected — invalid app secret.\x1b[0m");
      }
      setStatus(e.wasClean ? "idle" : "error");
      wsRef.current = null;
    };

    ws.onerror = () => setStatus("error");
  }, []);

  const disconnect = useCallback(() => {
    wsRef.current?.close(1000, "user disconnect");
  }, []);

  // ── Toggle modifier key ───────────────────────────────────────────────────
  const toggleMod = (mod: keyof Mods) => {
    setMods((prev) => {
      const next = { ...prev, [mod]: !prev[mod] };
      modsRef.current = next;
      return next;
    });
  };

  // sendKey reads modsRef (not state) so it always sees the latest modifiers
  // without needing to be recreated on every state change
  const sendKey = useCallback((key: DirectKey) => {
    if (wsRef.current?.readyState !== WebSocket.OPEN) return;

    const m = modsRef.current;
    let data = key.data;

    if (m.ctrl) {
      if (data.length === 1) {
        // Single char: compute control code, e.g. C→\x03, D→\x04
        data = ctrlCode(data);
      } else if (data.startsWith("\x1b[") && !data.includes(";")) {
        // Arrow/special: \x1b[A → \x1b[1;5A (Ctrl), \x1b[1;6A (Ctrl+Shift)
        const p = m.shift ? MODIFIER_PARAM["ctrl+shift"] : MODIFIER_PARAM["ctrl"];
        data = data.replace("\x1b[", `\x1b[1;${p}`);
      }
    } else if (m.shift && data.startsWith("\x1b[") && !data.includes(";")) {
      // Shift+arrow: \x1b[A → \x1b[1;2A
      data = data.replace("\x1b[", `\x1b[1;${MODIFIER_PARAM["shift"]}`);
    }

    // Alt: prepend ESC (works independently or combined with shift)
    if (m.alt && !m.ctrl) {
      data = "\x1b" + data;
    }

    wsRef.current.send(data);
    termRef.current?.focus();

    // Auto-clear all modifiers after use (Termux behaviour)
    const cleared: Mods = { ctrl: false, shift: false, alt: false };
    modsRef.current = cleared;
    setMods(cleared);
  }, []); // stable — no deps needed since we read from refs

  // ── Status badge text/colour ──────────────────────────────────────────────
  const badgeClass: Record<ConnStatus, string> = {
    idle:       "border-[#2a2a2a] text-[#444444]",
    connecting: "border-[#3a3a3a] text-[#666666]",
    connected:  "border-emerald-800/60 bg-emerald-950/30 text-emerald-400",
    error:      "border-red-800/60 bg-red-950/30 text-red-400",
  };
  const badgeLabel: Record<ConnStatus, string> = {
    idle: "disconnected", connecting: "connecting…",
    connected: "connected", error: "error",
  };

  const isConnected  = status === "connected";
  const isConnecting = status === "connecting";

  return (
    <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} className="flex flex-col gap-0">

      {/* ── Title bar ──────────────────────────────────────────────────────── */}
      <div className="flex h-10 items-center gap-3 rounded-t-lg border border-[#1e1e1e] bg-[#0f0f0f] px-3.5">
        {/* macOS traffic lights */}
        <div className="flex gap-1.5 shrink-0">
          <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
          <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
          <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
        </div>

        <div className="flex flex-1 items-center justify-center gap-2">
          <span className="text-[10px] font-bold uppercase tracking-widest text-[#444444]">
            terminal
          </span>
          <div className={cn(
            "flex items-center gap-1.5 rounded-full border px-2 py-0.5 text-[10px] transition-all",
            badgeClass[status]
          )}>
            <div className={cn(
              "h-1.5 w-1.5 rounded-full bg-current",
              isConnected && "animate-pulse"
            )} />
            {badgeLabel[status]}
          </div>
        </div>

        {/* Right actions */}
        <div className="flex items-center gap-1.5">
          {isConnected ? (
            <button
              onClick={disconnect}
              className="rounded px-2 py-0.5 text-[10px] text-[#555555] hover:text-red-400 transition-colors"
            >
              disconnect
            </button>
          ) : (
            <button
              onClick={() => setShowSecret(true)}
              disabled={isConnecting}
              className="rounded px-2 py-0.5 text-[10px] text-[#555555] hover:text-[#ededed] disabled:opacity-40 transition-colors"
            >
              connect
            </button>
          )}
          <button
            onClick={() => termRef.current?.clear()}
            className="rounded px-2 py-0.5 text-[10px] text-[#444444] hover:text-[#888888] transition-colors"
          >
            clear
          </button>
        </div>
      </div>

      {/* ── xterm area ─────────────────────────────────────────────────────── */}
      <div className="relative border-x border-[#1e1e1e] bg-[#0a0a0a]">
        <div ref={containerRef} className="h-[420px] w-full p-1" />

        {/* Overlays */}
        <AnimatePresence>
          {/* Secret input modal */}
          {showSecret && (
            <SecretModal key="secret" onConnect={connect} />
          )}

          {/* Idle / error overlay when not connected and no modal */}
          {!showSecret && !isConnected && !isConnecting && (
            <motion.div
              key="idle"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              className="absolute inset-0 flex flex-col items-center justify-center gap-3 bg-[#0a0a0a]/90"
            >
              <p className="text-xs text-[#444444]">
                {status === "error" ? "connection failed — " : "not connected — "}
                <button
                  onClick={() => setShowSecret(true)}
                  className="text-[#888888] underline underline-offset-2 hover:text-[#ededed]"
                >
                  {status === "error" ? "retry" : "connect"}
                </button>
              </p>
            </motion.div>
          )}

          {/* Connecting spinner */}
          {isConnecting && (
            <motion.div
              key="connecting"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              className="absolute inset-0 flex items-center justify-center bg-[#0a0a0a]/80"
            >
              <Loader2 className="h-4 w-4 animate-spin text-[#333333]" />
            </motion.div>
          )}
        </AnimatePresence>
      </div>

    </motion.div>
  );
}
