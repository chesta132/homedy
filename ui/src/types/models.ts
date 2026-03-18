// Samba share boolean values as used by the backend
export type SambaBool = "yes" | "no";

// A single SMB share entry matching backend models.Share
export interface Share {
  path: string;
  read_only: SambaBool;
  browsable: SambaBool;
  valid_users: string[];
  admin_users: string[];
  /** Linux-style file permissions as [owner, group, others] where each digit is 0-7 */
  permissions: number[];
}

// A map of share name → share config as returned by the backend
export type Shares = Record<string, Share>;

// Global samba config key-value pairs ([global] section)
export type ShareMap = Record<string, string>;

// User model returned on auth endpoints
export interface User {
  id: string;
  username: string;
  email: string;
  created_at: string;
  updated_at: string;
}
