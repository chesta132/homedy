// Fields available for form validation
export type FormFields = Partial<{
  identifier: string;
  username: string;
  email: string;
  password: string;
  rememberMe: boolean;
}>;

// Per-field config to toggle which rules are active
export type FormConfig<T extends FormFields = FormFields> = {
  [K in keyof T]?: boolean;
};

export type ValidationRule<T> = {
  condition: (value: T, config: FormConfig, form: FormFields) => boolean;
  message: string | ((config: FormConfig) => string);
};

export type FieldValidations<T extends FormFields> = {
  [K in keyof T]?: ValidationRule<T[K]>[];
};
