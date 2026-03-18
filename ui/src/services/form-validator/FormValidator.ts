import type { FormConfig, FormFields } from "@/types/form";
import { VALIDATION_RULES } from "./rules";

/**
 * Validates form fields against the configured rule set.
 * Ported from hoshify-client FormValidator pattern.
 */
export class FormValidator<T extends FormFields> {
  private form: T;
  private config: FormConfig<T>;

  constructor(form: T, config: FormConfig<T>) {
    this.form = form;
    this.config = config;
  }

  validateForm() {
    return this.validateFields(this.form);
  }

  validateFields(fields: Partial<T>) {
    let hasError = false;
    const errors: Partial<Record<keyof T, string>> = {};

    for (const [key, value] of Object.entries(fields) as [keyof FormFields, any][]) {
      const rules = VALIDATION_RULES[key];
      if (!rules || value == null) continue;

      for (const rule of rules) {
        if (rule.condition(value as any, this.config as FormConfig, this.form)) {
          const message =
            typeof rule.message === "function" ? rule.message(this.config as FormConfig) : rule.message;
          errors[key as keyof T] = message;
          hasError = true;
          break;
        }
      }
    }

    return { errors, hasError };
  }
}
