import type { FieldValidations, FormFields } from "@/types/form";

/**
 * Username rules:
 *   - 3–30 chars, alphanumeric + underscores + dots
 *   - Must start/end with alphanumeric or underscore
 *   - Cannot be all digits
 * Matches backend validatorlib.ValidateUsername
 */
const usernameRegex = /^[a-zA-Z0-9_]([a-zA-Z0-9_.]{1,28}[a-zA-Z0-9_]|[a-zA-Z0-9_]?)$/;
const allDigitsRegex = /^[0-9]+$/;

/**
 * Password rules:
 *   - 8–32 chars, letters + digits + special chars (@$!%*?&)
 *   - Must have at least one uppercase, lowercase, and digit
 * Matches backend validatorlib.ValidatePassword
 */
const passwordRegex = /^[A-Za-z\d@$!%*?&]{8,32}$/;

function hasLowerUpperDigit(value: string): boolean {
  let hasLower = false;
  let hasUpper = false;
  let hasDigit = false;
  for (const c of value) {
    if (c >= "a" && c <= "z") hasLower = true;
    else if (c >= "A" && c <= "Z") hasUpper = true;
    else if (c >= "0" && c <= "9") hasDigit = true;
  }
  return hasLower && hasUpper && hasDigit;
}

export const VALIDATION_RULES: FieldValidations<FormFields> = {
  identifier: [
    {
      condition: (v, cfg) => !!cfg.identifier && (!v || v.trim() === ""),
      message: "Email or username is required",
    },
  ],
  username: [
    {
      condition: (v, cfg) => !!cfg.username && (!v || v.trim() === ""),
      message: "Username is required",
    },
    {
      condition: (v, cfg) =>
        !!cfg.username && !!v && (!usernameRegex.test(v) || allDigitsRegex.test(v)),
      message: "Username must be 3–30 chars, start/end with letter or underscore, no special chars",
    },
  ],
  email: [
    {
      condition: (v, cfg) => !!cfg.email && (!v || v.trim() === ""),
      message: "Email is required",
    },
    {
      condition: (v, cfg) => !!cfg.email && !!v && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(v),
      message: "Please enter a valid email address",
    },
  ],
  password: [
    {
      condition: (v, cfg) => !!cfg.password && (!v || v.trim() === ""),
      message: "Password is required",
    },
    {
      condition: (v, cfg) =>
        !!cfg.password && !!v && (!passwordRegex.test(v) || !hasLowerUpperDigit(v)),
      message: "Password must be 8–32 chars and include uppercase, lowercase, and a digit",
    },
  ],
};
