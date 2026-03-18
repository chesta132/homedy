# AI_CONTEXT.md — Homedy Frontend

This file is the single source of truth for any AI assistant continuing work on this codebase.
Read this before touching any file.

---

## Project Overview

**Homedy** is a self-hosted home server management dashboard. This repository is the **frontend client** only. The Go backend is a separate project (`homedy-main`).

**Stack:**

- Vite + React 19 + TypeScript (strict)
- Tailwind CSS v4 (uses `@tailwindcss/vite` plugin — no `tailwind.config.js`, config is in `src/index.css` via `@import "tailwindcss"`)
- Shadcn-style UI components (hand-written, Radix UI primitives, NOT the shadcn CLI)
- Axios for HTTP
- React Router v7
- Framer Motion for animations
- Lucide React for icons
- `clsx` + `tailwind-merge` via `cn()` util
- `sonner` for toasts
- `@xterm/xterm` + `@xterm/addon-fit` for terminal
- `dayjs` available but not yet used

**Run:**

```bash
npm install
npm run dev      # dev server on :3000, proxies /api → localhost:8080
npm run build    # production build
```

---

## Visual Design Rules

The entire app uses a **dark monochrome theme** — no colored accents, no blue/navy.

| Token          | Value                                | Usage                    |
| -------------- | ------------------------------------ | ------------------------ |
| Page bg        | `#0a0a0a`                            | Body, main backgrounds   |
| Surface        | `#111111`                            | Cards, panels            |
| Surface raised | `#1a1a1a`                            | Inputs, dropdowns        |
| Border         | `#1e1e1e` / `#2a2a2a`                | Subtle / visible borders |
| Text primary   | `#ededed`                            | Headings, active text    |
| Text secondary | `#888888`                            | Labels, descriptions     |
| Text muted     | `#555555` / `#444444`                | Placeholders, hints      |
| Primary action | `bg-white text-black`                | Main buttons             |
| Destructive    | `text-red-400 bg-red-950/30`         | Delete, errors           |
| Success        | `text-emerald-400 bg-emerald-950/30` | Connected, ok            |

**Never** introduce blue, purple, or any saturated accent color into the main UI.
The terminal xterm theme also follows this palette (see `TerminalPage.tsx → XTERM_THEME`).

---

## Directory Structure

```
src/
├── App.tsx                         # Router root — all routes defined here
├── main.tsx                        # React entry point
├── index.css                       # Tailwind v4 import + base layer
│
├── types/
│   ├── models.ts                   # Domain types: User, Share, Shares, ShareMap, SambaBool
│   ├── server.ts                   # ApiResponse envelope, ErrorResponse, AUTH_ERROR_CODES
│   └── form.ts                     # FormFields, FormConfig, ValidationRule types
│
├── lib/
│   └── utils.ts                    # cn() = clsx + twMerge
│
├── contexts/
│   └── AuthContext.tsx             # user, setUser, globalError, setGlobalError
│
├── hooks/
│   ├── useForm.ts                  # Form state + validation hook (hoshify pattern)
│   └── useAppSecret.ts             # Modal-based app_secret prompt (never stores secret)
│
├── services/
│   ├── server/
│   │   ├── ApiClient.ts            # Axios wrapper — api.auth, api.samba, api.sambaConfig
│   │   ├── ServerSuccess.ts        # Wraps successful axios response
│   │   └── ServerError.ts          # Wraps axios error response
│   ├── form-validator/
│   │   ├── FormValidator.ts        # Validates FormFields against VALIDATION_RULES
│   │   └── rules.ts                # Per-field rules matching backend validators
│   └── models/
│       └── handleError.ts          # handleError(), handleFormError()
│
├── components/
│   ├── AuthGuard.tsx               # Fetches /auth/me, shows spinner, redirects on fail
│   ├── ui/                         # Radix-based components (hand-written, not CLI)
│   │   ├── button.tsx              # variants: default | outline | ghost | destructive
│   │   ├── input.tsx               # error prop shows red border + message below
│   │   ├── label.tsx
│   │   ├── checkbox.tsx
│   │   ├── switch.tsx
│   │   ├── badge.tsx               # variants: default | outline | success | destructive
│   │   ├── card.tsx                # Card, CardHeader, CardTitle, CardDescription, CardContent
│   │   ├── dialog.tsx              # Dialog, DialogContent, DialogHeader, DialogTitle, etc.
│   │   ├── dropdown-menu.tsx       # DropdownMenu, DropdownMenuItem, etc.
│   │   ├── tabs.tsx                # Tabs, TabsList, TabsTrigger, TabsContent
│   │   ├── table.tsx               # Table, TableHeader, TableBody, TableRow, etc.
│   │   ├── avatar.tsx              # Avatar, AvatarFallback, AvatarImage
│   │   ├── separator.tsx
│   │   ├── tag-input.tsx           # Multi-value chip input (Enter/comma to add, × to remove)
│   │   ├── permission-input.tsx    # 3-digit owner/group/others input + info modal
│   │   ├── app-secret-modal.tsx    # Modal asking for APP_SECRET (reusable)
│   │   └── toaster.tsx             # Sonner wrapper — toast.success(), toast.error()
│   ├── dashboard/
│   │   ├── DashboardLayout.tsx     # Fixed sidebar + topbar + <Outlet />
│   │   ├── Sidebar.tsx             # Desktop fixed + mobile slide-in (framer-motion)
│   │   └── Topbar.tsx              # User avatar menu — Profile (navigate) + Sign Out
│   └── smb/
│       ├── ShareForm.tsx           # Create/edit share modal — all backend Share fields
│       ├── SharesTable.tsx         # Desktop table + mobile card list
│       ├── DeleteShareDialog.tsx   # Confirmation dialog before delete
│       └── SambaConfigEditor.tsx   # Global config editor — desktop inline / mobile stacked
│
└── pages/
    ├── auth/
    │   ├── SignInPage.tsx          # identifier + password + remember_me
    │   └── SignUpPage.tsx          # username + email + password + remember_me
    ├── DashboardPage.tsx           # Quick-access cards (SMB available, others locked)
    ├── SMBPage.tsx                 # File Sharing — shares tab + config tab
    └── TerminalPage.tsx            # xterm.js terminal over WebSocket
```

---

## Routing

All routes are in `src/App.tsx`:

```
/signin                → SignInPage       (public)
/signup                → SignUpPage       (public)
/dashboard             → DashboardPage   (protected, inside DashboardLayout)
/dashboard/smb         → SMBPage         (protected)
/dashboard/terminal    → TerminalPage    (protected)
/dashboard/profile     → not yet built   (navigate goes here from Topbar)
*                      → redirect /signin
```

**AuthGuard** (`components/AuthGuard.tsx`):

- Hits `GET /api/auth/me` on mount
- State machine: `"checking"` → spinner | `"ok"` → render | `"fail"` → redirect `/signin`
- AbortError (component unmount) is ignored — does NOT trigger redirect

---

## API Client

Single singleton `api` exported from `src/services/server/ApiClient.ts`.

```ts
import api from "@/services/server/ApiClient";

// Namespaced sub-clients — all share ONE axios instance (baseURL: /api)
api.auth.get<User>("/me");
api.auth.post<User>("/signin", body);
api.auth.post<User>("/signup", body);
api.auth.post("/signout");

api.samba.get<Shares>("/");
api.samba.get<Share>("/:name");
api.samba.post<Shares>("/", body);
api.samba.put<Shares>("/:name", body);
api.samba.delete<Shares>("/:name");
api.samba.post<Shares>("/backup", undefined, { params: { app_secret } });
api.samba.post<Shares>("/restore", undefined, { params: { app_secret } });

api.sambaConfig.get<ShareMap>("/", { params: { app_secret } });
api.sambaConfig.put("/", body, { params: { app_secret } });
```

**Response shape** from backend:

```json
{ "meta": { "status": "SUCCESS" }, "data": <T> }
{ "meta": { "status": "ERROR" }, "data": { "code": "...", "message": "...", "field": "..." } }
```

**Error handling**: if response code is in `AUTH_ERROR_CODES` (`UNAUTHORIZED`, `INVALID_AUTH`, `INVALID_TOKEN`), the interceptor redirects to `/signin` automatically.

**Important — never construct `new ApiClient(baseURL)` inside a constructor.** The sub-client pattern uses a single shared axios instance passed by reference. New sub-clients are only created when `prefix === ""` (root client). See the comment in `ApiClient.ts`.

---

## Backend Endpoints Reference

Backend is Go + Gin. All routes mount under no prefix (backend runs on `:8080`, frontend proxies `/api` → `http://localhost:8080`).

### Auth — `/auth/*`

| Method | Path     | Auth?  | Body                                                     | Response |
| ------ | -------- | ------ | -------------------------------------------------------- | -------- |
| POST   | /signup  | No     | `{ username, email, password, remember_me }`             | User     |
| POST   | /signin  | No     | `{ identifier, email\|username, password, remember_me }` | User     |
| POST   | /signout | Cookie | —                                                        | null     |
| GET    | /me      | Cookie | —                                                        | User     |

Auth uses **HTTP-only cookies** (access + refresh token). `withCredentials: true` is set on the axios instance.

### Samba — `/samba/*`

| Method | Path     | Auth?                | Body / Notes              | Response |
| ------ | -------- | -------------------- | ------------------------- | -------- |
| GET    | /        | Cookie               | —                         | Shares   |
| POST   | /        | Cookie               | Share body + `name` field | Shares   |
| GET    | /:name   | Cookie               | —                         | Share    |
| PUT    | /:name   | Cookie               | Share body                | Shares   |
| DELETE | /:name   | Cookie               | —                         | Shares   |
| POST   | /backup  | Cookie + ?app_secret | —                         | null     |
| POST   | /restore | Cookie + ?app_secret | —                         | Shares   |
| GET    | /config/ | Cookie + ?app_secret | —                         | ShareMap |
| PUT    | /config/ | Cookie + ?app_secret | ShareMap body             | ShareMap |

**`app_secret`** is sent as a query param: `?app_secret=<value>`. The middleware (`AppProtected`) checks it against the server env var `APP_SECRET`. Frontend **never stores** the secret — `useAppSecret` prompts via modal every time.

### WebSocket Terminal — `/ws/terminal`

| Param  | Value                                                |
| ------ | ---------------------------------------------------- |
| Query  | `?app_secret=<value>`                                |
| Auth   | Cookie (same as HTTP)                                |
| WS URL | `ws[s]://<host>/api/ws/terminal?app_secret=<secret>` |

**Messages client → server:**

- Raw bytes/string → piped directly to PTY (bash)
- JSON `{ "type": "resize", "resize": { "x": <cols>, "y": <rows> } }` → calls `pty.Setsize`

**Messages server → client:**

- Binary frames (`ArrayBuffer`) → `term.write(new Uint8Array(data))`

---

## Share Model

Matches Go `models.Share` exactly:

```ts
interface Share {
  path: string; // absolute path, e.g. "/srv/media"
  read_only: "yes" | "no";
  browsable: "yes" | "no";
  valid_users: string[]; // min 1 required
  admin_users: string[]; // min 1 required
  permissions: number[]; // [owner, group, others], each 0–7 (octal)
}
```

Share name is sent as `name` in the POST body and as `:name` URL param for PUT/DELETE. Name cannot be changed after creation. Reserved names rejected by backend: `global`, `printers`, `print$`, `config`, `backup`, `restore`.

---

## Form Validation

Pattern ported from `hoshify-client`:

```ts
const {
  form: [form, setForm],
  error: [errors, setErrors],
  validate,
} = useForm(
  { identifier: "", password: "", rememberMe: false as boolean },
  { identifier: true, password: true }, // config enables rules per field
);

validate.validateField({ identifier: value }); // on onChange
validate.validateForm(); // on submit — returns boolean
```

**Validation rules** in `services/form-validator/rules.ts` match the backend:

- `username`: regex `^[a-zA-Z0-9_]([a-zA-Z0-9_.]{1,28}[a-zA-Z0-9_]|[a-zA-Z0-9_]?)$`, not all-digits
- `password`: 8–32 chars, must have upper + lower + digit, allowed specials `@$!%*?&`
- `email`: basic format check
- `identifier`: required only

---

## Terminal Page

`src/pages/TerminalPage.tsx`

**xterm config:**

```ts
fontFamily: '"Fira Code", "Cascadia Code", monospace'
fontSize: 13, lineHeight: 1.5
cursorStyle: "block", cursorBlink: true
theme: XTERM_THEME  // monochrome, follows app palette
```

- **Stale closure fix**: modifiers are tracked via `modsRef` (ref that mirrors state). `sendKey` reads from `modsRef.current`, NOT from state. This prevents the bug where `useCallback` captures stale mod values
- Ctrl+letter: sends `\x01`–`\x1a` range
- Ctrl+arrow: sends `\x1b[1;5A` etc. (ANSI CSI modifier param)
- Alt: prepends `\x1b`
- Direct keys: Tab, Esc, Enter, arrows, Home/End, PgUp/PgDn, Del, Ins

**Connection flow:**

1. User clicks "connect" → `SecretModal` appears (inline overlay inside xterm area)
2. User enters `APP_SECRET` → `connect(secret)` called
3. WS URL: `ws[s]://<window.location.host>/api/ws/terminal?app_secret=<secret>`
4. On open: register `term.onData` handler, call `fitAddon.fit()`, send initial resize
5. `term.onResize` hook → sends `{ type: "resize", resize: { x: cols, y: rows } }` on every fit

---

## Key Patterns & Conventions

### Adding a new page

1. Create `src/pages/NewPage.tsx`
2. Add route in `src/App.tsx` inside the `<AuthGuard>` + `<DashboardLayout>` block
3. Add nav item in `src/components/dashboard/Sidebar.tsx` → `NAV_ITEMS` array (set `comingSoon: true` to lock it)

### Adding a new API endpoint

- Add call via `api.<namespace>.<method>(url, ...)` directly in the component or a service file
- If the endpoint needs `app_secret`: use `useAppSecret` hook, call `getSecret()` before the request, pass result as `{ params: { app_secret: secret } }`
- Wrap in try/catch, use `toast.error()` for failures

### Adding a new UI component

- Extends Radix primitive OR is a pure HTML component
- Always use `cn()` for className merging
- Follow color tokens above — no hardcoded non-palette colors
- Export named (not default)

### Error handling

```ts
import { handleError, handleFormError } from "@/services/models/handleError";

// Global error state
handleError(err, setGlobalError);

// Form field error (sets field error if err.field exists, otherwise global)
handleFormError(err, setErrors, setGlobalError);
```

`ServerError` exposes: `.getCode()`, `.getMessage()`, `.getField()`, `.getTitle()`, `.status`

### Dropdown navigation (IMPORTANT)

Do **NOT** use `<DropdownMenuItem asChild><Link>` — the Radix Item does not forward the Slot correctly in this codebase and causes click events to misfire.  
Use `onClick={() => navigate("/path")}` instead:

```tsx
<DropdownMenuItem onClick={() => navigate("/dashboard/profile")}>Profile</DropdownMenuItem>
```

---

## Currently Locked / Not Yet Built

Features shown in sidebar with "Soon" badge — no routes or pages exist yet:

- Chat
- DNS
- Port Forward
- PDF Converter
- Todo
- Notes
- Finance
- `/dashboard/profile` page (route navigates there but no page component exists)

---

## Known Decisions / Non-obvious Choices

| Decision                                     | Reason                                                                                                           |
| -------------------------------------------- | ---------------------------------------------------------------------------------------------------------------- |
| `app_secret` never stored                    | Security requirement — user asked on every protected action                                                      |
| `useAppSecret` uses Promise resolver pattern | Allows `await getSecret()` inline in async handlers without prop-drilling callbacks                              |
| `modsRef` mirrors `mods` state in terminal   | Prevents stale closure in `sendKey` `useCallback` — refs always have latest value                                |
| No resize JSON sent to old backend           | Backend now handles it (updated `InputToPTY`). Resize format: `{ type: "resize", resize: { x: cols, y: rows } }` |
| Sub-clients don't recurse                    | `ApiClient` only builds sub-clients when `prefix === ""` — prevents infinite constructor recursion               |
| `AuthGuard` state machine                    | `"checking"/"ok"/"fail"` instead of boolean — separates AbortError (normal unmount) from real auth failure       |
| Config editor dual layout                    | Desktop: inline `key = value ×` row. Mobile (`sm:hidden`): stacked card with Key/Value labels                    |
